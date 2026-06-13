// SPDX-License-Identifier: AGPL-3.0-only
// 16维度提示词生成器 — 参考: designer-image (10维度) + CLIP Interrogator
package prompt

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func GeneratePrompt(scene string, platform string, refImage string, baseDir string) (*PromptScaffold, error) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		return generateGeneric(scene, platform, refImage), nil
	}

	s := &PromptScaffold{
		Seed: scene, Platform: platform, Model: tmpl.Name["zh"],
		Fields: make(map[string]string), Weights: make(map[string]float64),
		Negative: NegativeBank[platform], Params: tmpl.Params, RefImage: refImage,
	}
	// Fallback: if NegativeBank has no entry for this platform, use YAML template's negative
	if s.Negative == nil && len(tmpl.Negative) > 0 {
		s.Negative = tmpl.Negative
	}
	isVideo := isVideoPlatform(platform)
	for _, sec := range tmpl.Sections {
		if sec.VideoOnly && !isVideo {
			continue
		}
		val := sec.Default
		if val == "{{seed}}" {
			val = s.Seed
		}
		if val == "" {
			val = fmt.Sprintf("（LLM填充：%s）", sec.Hint)
		}
		s.Fields[sec.ID] = val
		if sec.Weight > 0 {
			s.Weights[sec.ID] = sec.Weight
		}
	}
	if refImage != "" {
		injectReferenceImage(s, refImage)
	}
	s.FinalPrompt = buildWeightedPrompt(tmpl, s.Fields, s.Weights, s.Negative)
	s.Expanded = buildPromptDebug(s, tmpl)
	return s, nil
}

// isVideoPlatform returns true if the platform supports video dimensions.
func isVideoPlatform(platform string) bool {
	videoPlatforms := map[string]bool{
		"kling": true, "veo": true, "sora": true, "pika": true,
		"seedance": true, "video_generic": true, "video_camera": true, "video_keyframe": true,
	}
	return videoPlatforms[platform]
}

func buildWeightedPrompt(tmpl *PromptTemplate, fields map[string]string, weights map[string]float64, negative []string) string {
	categories := []string{"subject", "style", "comp", "lighting", "color", "env", "detail", "quality"}
	categorized := make(map[string][]string)
	for _, sec := range tmpl.Sections {
		v, ok := fields[sec.ID]
		if !ok || v == "" || strings.HasPrefix(v, "（LLM") {
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
	if vals, ok := categorized["other"]; ok {
		parts = append(parts, vals...)
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

func applyWeights(value string, weight float64) string {
	if weight == 0 || weight == 1.0 {
		return value
	}
	return fmt.Sprintf("(%s:%.1f)", value, weight)
}

func injectReferenceImage(s *PromptScaffold, refPath string) {
	for k := range s.Fields {
		if strings.HasPrefix(s.Fields[k], "（LLM填充") {
			s.Fields[k] = fmt.Sprintf("（参考图 %s 提取: %s）", filepath.Base(refPath), s.Fields[k])
		}
	}
}

func buildPromptDebug(s *PromptScaffold, tmpl *PromptTemplate) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📋 平台: %s | 维度: %d\n", s.Platform, len(s.Fields)))
	for _, sec := range tmpl.Sections {
		val := s.Fields[sec.ID]
		mark := "✅"
		if strings.HasPrefix(val, "（LLM") {
			mark = "⏳"
		}
		w := s.Weights[sec.ID]
		ws := ""
		if w > 0 && w != 1.0 {
			ws = fmt.Sprintf(" ×%.1f", w)
		}
		b.WriteString(fmt.Sprintf("  %s [%s] %s%s\n", mark, sec.ID, val[:min(len(val), 60)], ws))
	}
	b.WriteString(fmt.Sprintf("\n最终 prompt:\n%s\n", s.FinalPrompt))
	return b.String()
}

func (s *PromptScaffold) UpdateField(id, value string)           { s.Fields[id] = value }
func (s *PromptScaffold) UpdateWeight(id string, weight float64) { s.Weights[id] = weight }
func (s *PromptScaffold) Regenerate(platform string, baseDir string) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		s.FinalPrompt = s.Seed
		return
	}
	s.FinalPrompt = buildWeightedPrompt(tmpl, s.Fields, s.Weights, s.Negative)
	s.Expanded = buildPromptDebug(s, tmpl)
}

func generateGeneric(scene, platform, refImage string) *PromptScaffold {
	return &PromptScaffold{
		Seed: scene, Platform: platform, Model: "通用",
		Fields:   map[string]string{"scene": scene},
		Negative: NegativeBank[platform], FinalPrompt: scene, RefImage: refImage,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
