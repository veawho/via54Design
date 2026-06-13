// SPDX-License-Identifier: AGPL-3.0-only
//
// i2i.go — v2.3 image-to-image bilingual prompt engine.
//
// Inputs:  scene (free text, user's language), ref_image (path), platform,
//          max_chars (target prompt budget).
// Outputs: I2IResult{ZH, EN, RefDescription, Sections, Debug} — ZH is the
//          formatted bilingual view the user reads, EN is the final-prompt
//          text pasted into the target image-gen tool.
//
// Pipeline:
//
//   1. Load platform YAML template
//   2. Optional: ref image → visual description (caller-provided; we just
//      consume a pre-extracted string so the engine stays LLM-free)
//   3. For each section, decide value:
//      - if it has a user override → use it
//      - elif ref image is present → write a ref-anchored placeholder the
//        LLM will fill from the ref description
//      - elif YAML has a default → use it
//      - else → leave a "（LLM填充）" marker
//   4. For free-text fields (subject/secondary/environment) where the user
//      provided scene text, LLM-translate the scene text into the target
//      language using a single hermes -z call (batched, one per field).
//   5. Render the final prompt within max_chars:
//      -EN side: comma-joined, weight syntax (group:1.2), negatives, params
//      -ZH side: formatted block per category, with EN tail in parens
//   6. fillToMax: if final < max_chars and there are unfilled placeholders,
//      ask the LLM to expand the most-important ones. Cap at max_chars - 10.
//
// Engine NEVER embeds an LLM. Translation is done by the caller (feishu bot
// or inbox_watcher) via hermes -z pipes. This keeps via54Design hermes-
// independent at the binary level (Hermes is a deployment convenience, not
// a hard requirement).
package prompt

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// I2IRequest is the input to GenerateI2I.
type I2IRequest struct {
	Scene          string             // user text in any language (Chinese typical)
	Platform       string             // target platform id
	RefImage       string             // path to ref image; empty = no ref
	RefDescription string             // pre-extracted visual description (caller fills via LLaVA/mmx)
	MaxChars       int                // target prompt budget; 0 = platform default
	Overrides      map[string]string  // user-edit overrides keyed by section.ID
	BaseDir        string             // repo base dir for template loading
}

// I2IResult is the output of GenerateI2I.
type I2IResult struct {
	Platform       string             // resolved platform
	Model          string             // platform display name (zh)
	FinalEN        string             // ⭐ the English prompt to paste into image-gen tool
	FormattedZH    string             // ⭐ the formatted bilingual view the user reads
	RefDescription string             // extracted visual description (if ref image was supplied)
	RefImage       string             // echo back the ref path
	MaxChars       int                // budget used
	FinalChars     int                // actual length of FinalEN
	Sections       map[string]string  // resolved per-section values (EN rendering)
	SectionsZH     map[string]string  // resolved per-section values (ZH rendering)
	Negative       []string           // final negative list
	Params         map[string]string  // platform params (--ar 16:9 etc)
	Unfilled       []string           // section IDs that still have "（LLM填充）" markers
	Debug          string             // human-readable debug summary
}

// MaxChars returns the platform-specific prompt budget.
// 0 or negative falls back to a sensible default.
func MaxCharsForPlatform(platform string) int {
	switch platform {
	case "midjourney":
		return 4500 // V6 raised to 4500 chars; older v5 4000
	case "jimeng":
		return 1500
	case "gemini":
		return 8000
	case "dalle3":
		return 4000
	case "flux":
		return 4000
	case "sd3", "stable_diffusion":
		return 2000 // SDXL token limit ≈ 77*4=308 words; 2000 chars safe
	case "kling", "pika", "sora", "veo", "seedance", "minimax", "video_camera", "video_generic", "video_keyframe":
		return 2500
	case "ideogram", "recraft":
		return 2500
	}
	return 2000
}

// GenerateI2I is the public entry point.
//
//   s, err := prompt.GenerateI2I(prompt.I2IRequest{
//       Scene:    "温馨家庭客厅, 蓝底清爽风",
//       Platform: "jimeng",
//       RefImage: "/Users/david/Desktop/ref.png",
//       RefDescription: "modern Scandinavian living room, soft daylight, beige sofa, plants",
//       MaxChars: 1500,
//       BaseDir:  "/Users/david/Desktop/developments/via54Design",
//   })
func GenerateI2I(req I2IRequest) (*I2IResult, error) {
	if req.Platform == "" {
		return nil, fmt.Errorf("platform is required")
	}
	if req.BaseDir == "" {
		return nil, fmt.Errorf("baseDir is required")
	}

	// 1. Load platform YAML
	tmpl := loadTemplate(req.Platform, req.BaseDir)
	if tmpl == nil {
		// Fallback: treat scene as-is, no YAML structure
		return genericI2I(req), nil
	}

	// 2. Determine max chars
	if req.MaxChars <= 0 {
		req.MaxChars = MaxCharsForPlatform(req.Platform)
	}

	// 3. Walk sections, resolve values
	isVideo := isVideoPlatform(req.Platform)
	resolvedEN := make(map[string]string)
	resolvedZH := make(map[string]string)
	sectionWeights := make(map[string]float64)
	sectionCategories := make(map[string]string)
	labels := make(map[string]string)
	var unfilled []string

	for _, sec := range tmpl.Sections {
		// Skip video-only sections for image platforms
		if sec.VideoOnly && !isVideo {
			continue
		}
		labels[sec.ID] = sec.Label
		sectionCategories[sec.ID] = sec.Category
		if sec.Weight > 0 {
			sectionWeights[sec.ID] = sec.Weight
		}

		val := ""
		// Override wins
		if v, ok := req.Overrides[sec.ID]; ok && v != "" {
			val = v
		} else if sec.Default != "" {
			if sec.Default == "{{seed}}" {
				val = req.Scene
			} else {
				val = sec.Default
			}
		}

		// If still empty, this is a fill-by-LLM field
		if val == "" {
			val = fmt.Sprintf("（LLM填充：%s）", sec.Hint)
			unfilled = append(unfilled, sec.ID)
		}

		// For free-text subject/secondary/environment, the user gave us
		// scene text. We don't auto-inject scene into Subject unconditionally;
		// the caller (feishu bot) is expected to set Overrides["subject"]
		// from the user's intent. If they didn't, leave it as LLM-fill.
		// This is intentional — silent injection produces worse prompts than
		// explicit LLM translation.

		b := ParseBilingual(val)
		resolvedEN[sec.ID] = b.RenderEN()
		resolvedZH[sec.ID] = b.RenderZH()
	}

	// 4. Inject ref image info into "open" subject/secondary/environment
	//    IF the user didn't override them — this is the "lock ref into subject"
	//    behaviour the user asked for. The visual description (caller-extracted)
	//    becomes the anchor that we ask the LLM to "lock" against.
	//
	// Ref-lock semantics (v2.3):
	//   - subject  → "preserve visual identity of ref: <ref_desc>"
	//   - secondary → "supporting elements consistent with ref style"
	//   - environment → "setting/atmosphere from ref, expanded with scene"
	if req.RefImage != "" && req.RefDescription != "" {
		refBaseEN := fmt.Sprintf("preserve visual identity of reference image (%s): %s",
			filepath.Base(req.RefImage), req.RefDescription)
		refBaseZH := fmt.Sprintf("保留参考图(%s)的视觉身份: %s",
			filepath.Base(req.RefImage), req.RefDescription)
		envEN := fmt.Sprintf("%s; setting: %s",
			refBaseEN, req.Scene)
		envZH := fmt.Sprintf("%s; 场景: %s",
			refBaseZH, req.Scene)
		secondaryEN := fmt.Sprintf("supporting elements consistent with the reference image visual style: %s", req.RefDescription)
		secondaryZH := fmt.Sprintf("与参考图视觉风格一致的辅助元素: %s", req.RefDescription)

		if v, ok := resolvedEN["subject"]; ok && strings.HasPrefix(v, "（LLM填充") {
			resolvedEN["subject"] = refBaseEN
			resolvedZH["subject"] = refBaseZH
		}
		if v, ok := resolvedEN["secondary"]; ok && strings.HasPrefix(v, "（LLM填充") {
			resolvedEN["secondary"] = secondaryEN
			resolvedZH["secondary"] = secondaryZH
		}
		if v, ok := resolvedEN["environment"]; ok && strings.HasPrefix(v, "（LLM填充") {
			resolvedEN["environment"] = envEN
			resolvedZH["environment"] = envZH
		}
		// Strip from unfilled if we filled them
		cleaned := unfilled[:0]
		for _, id := range unfilled {
			if id != "subject" && id != "secondary" && id != "environment" {
				cleaned = append(cleaned, id)
			}
		}
		unfilled = cleaned
	}

	// 5. Build final prompts
	negative := NegativeBank[req.Platform]
	if negative == nil && len(tmpl.Negative) > 0 {
		negative = tmpl.Negative
	}
	finalEN := buildWeightedPromptEN(tmpl, resolvedEN, sectionWeights, negative)
	finalZH := buildFormattedZH(tmpl, resolvedEN, resolvedZH, sectionWeights, sectionCategories, labels, req.Platform, finalEN)

	// 6. fillToMax — if there's room, ask the LLM to expand unfilled
	//    markers. We don't call the LLM here; we surface what the caller
	//    should ask the LLM to do, then trim/pad FinalEN to the budget.
	if len(finalEN) < req.MaxChars && len(unfilled) > 0 {
		// Mark for caller expansion. The caller (inbox_watcher) will
		// call hermes -z with the unfilled hints, get back text, splice it
		// in via Overrides, re-call GenerateI2I.
		// We simply surface the budget gap.
	}

	// 7. Trim if over budget (defensive)
	if req.MaxChars > 0 && len(finalEN) > req.MaxChars {
		// Drop the lowest-priority fields first. For now: truncate.
		// (Caller can pre-edit Overrides to remove fields.)
		// We use a soft trim that keeps the parameter/negative tail intact.
		cutoff := req.MaxChars - len(" --no "+strings.Join(negative, ", ")) - 50
		if cutoff > 0 {
			idx := strings.LastIndex(finalEN[:cutoff], ", ")
			if idx > 0 {
				finalEN = finalEN[:idx] + " --no " + strings.Join(negative, ", ")
			}
		}
	}

	return &I2IResult{
		Platform:       req.Platform,
		Model:          tmpl.Name["zh"],
		FinalEN:        finalEN,
		FormattedZH:    finalZH,
		RefDescription: req.RefDescription,
		RefImage:       req.RefImage,
		MaxChars:       req.MaxChars,
		FinalChars:     len(finalEN),
		Sections:       resolvedEN,
		SectionsZH:     resolvedZH,
		Negative:       negative,
		Params:         tmpl.Params,
		Unfilled:       unfilled,
	}, nil
}

// genericI2I is the no-template fallback.
func genericI2I(req I2IRequest) *I2IResult {
	maxChars := req.MaxChars
	if maxChars <= 0 {
		maxChars = MaxCharsForPlatform(req.Platform)
	}
	finalEN := req.Scene
	if maxChars > 0 && len(finalEN) > maxChars {
		finalEN = finalEN[:maxChars-3] + "..."
	}
	return &I2IResult{
		Platform:    req.Platform,
		Model:       "通用",
		FinalEN:     finalEN,
		FormattedZH: "## 🎨 通用提示词\n\n```\n" + finalEN + "\n```\n",
		RefImage:    req.RefImage,
		MaxChars:    maxChars,
		FinalChars:  len(finalEN),
		Sections:    map[string]string{"scene": finalEN},
		SectionsZH:  map[string]string{"scene": req.Scene},
	}
}

// buildWeightedPromptEN produces the comma-joined, weight-suffixed, param-
// appended, negative-appended final English prompt.
//
//   (group:1.2), other, ... --ar 16:9 --v 6 --no blurry, low quality, ...
func buildWeightedPromptEN(tmpl *PromptTemplate, fields map[string]string, weights map[string]float64, negative []string) string {
	categories := []string{"subject", "style", "comp", "lighting", "color", "env", "detail", "quality", "motion", "other"}
	categorized := make(map[string][]string)
	for _, sec := range tmpl.Sections {
		v, ok := fields[sec.ID]
		if !ok || v == "" {
			continue
		}
		w := weights[sec.ID]
		if w == 0 {
			w = 1.0
		}
		if w != 1.0 && sec.Weighted {
			v = applyWeights(v, w)
		}
		cat := sec.Category
		if cat == "" {
			cat = "other"
		}
		categorized[cat] = append(categorized[cat], v)
	}
	var parts []string
	for _, cat := range categories {
		if vals, ok := categorized[cat]; ok {
			parts = append(parts, vals...)
		}
	}
	prompt := strings.Join(parts, ", ")
	var paramStr []string
	keys := make([]string, 0, len(tmpl.Params))
	for k := range tmpl.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		paramStr = append(paramStr, fmt.Sprintf("--%s %s", k, tmpl.Params[k]))
	}
	if len(paramStr) > 0 {
		prompt += " " + strings.Join(paramStr, " ")
	}
	if len(negative) > 0 {
		prompt += " --no " + strings.Join(negative, ", ")
	}
	return prompt
}

// buildFormattedZH produces the formatted bilingual block the user reads in
// the Feishu reply. The structure is:
//
//   ## 🎨 主体
//   **主体**: 摄影写实 (Photorealistic) · 现代 (Modern)
//   **辅助**: ...
//   **艺术风格**: 摄影写实 (Photorealistic)
//   ...
//   **参数**: --ar 16:9 --v 6
//
//   ## 📋 最终英文提示词 (用于粘贴)
//   (subject, ...photorealistic..., --ar 16:9)
//
//   --ar 16:9 --v 6 --no blurry, ...
func buildFormattedZH(tmpl *PromptTemplate, fieldsEN, fieldsZH map[string]string, weights map[string]float64, categories map[string]string, labels map[string]string, platform, finalEN string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## 🎨 %s 提示词\n\n", tmpl.Name["zh"]))

	// Group sections by category
	byCategory := make(map[string][]PromptSection)
	order := []string{}
	seen := make(map[string]bool)
	for _, sec := range tmpl.Sections {
		cat := sec.Category
		if cat == "" {
			cat = "other"
		}
		if !seen[cat] {
			seen[cat] = true
			order = append(order, cat)
		}
		byCategory[cat] = append(byCategory[cat], sec)
	}
	categoryNamesZH := map[string]string{
		"subject": "主体", "style": "风格", "comp": "构图",
		"lighting": "光线", "color": "色彩", "env": "环境",
		"detail": "细节", "quality": "质量", "motion": "运动",
		"other": "其他",
	}
	for _, cat := range order {
		b.WriteString(fmt.Sprintf("### %s\n", categoryNamesZH[cat]))
		for _, sec := range byCategory[cat] {
			zh := fieldsZH[sec.ID]
			en := fieldsEN[sec.ID]
			if zh == "" && en == "" {
				continue
			}
			display := zh
			if zh == en {
				display = zh
			} else if zh != "" && en != "" {
				display = fmt.Sprintf("%s (%s)", zh, en)
			}
			if w, ok := weights[sec.ID]; ok && w != 1.0 {
				display = fmt.Sprintf("%s ×%.1f", display, w)
			}
			label := sec.Label
			if label == "" {
				label = sec.ID
			}
			b.WriteString(fmt.Sprintf("- **%s**: %s\n", label, display))
		}
		b.WriteString("\n")
	}

	// Params
	if len(tmpl.Params) > 0 {
		b.WriteString("### 参数\n")
		keys := make([]string, 0, len(tmpl.Params))
		for k := range tmpl.Params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			b.WriteString(fmt.Sprintf("- `--%s %s`\n", k, tmpl.Params[k]))
		}
		b.WriteString("\n")
	}

	// Negative
	neg := NegativeBank[platform]
	if neg == nil && len(tmpl.Negative) > 0 {
		neg = tmpl.Negative
	}
	if len(neg) > 0 {
		b.WriteString("### 负面词\n")
		for _, n := range neg {
			b.WriteString(fmt.Sprintf("- %s\n", n))
		}
		b.WriteString("\n")
	}

	// Final English prompt — this is what the user copies after confirming
	b.WriteString("### 📋 最终英文提示词 (确认后可复制粘贴到 ")
	b.WriteString(tmpl.Name["zh"])
	b.WriteString(")\n```\n")
	b.WriteString(finalEN)
	b.WriteString("\n```\n")
	return b.String()
}
