// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

// via54Design — Storyboard → Video Pipeline
// 1. Single image → video opening/scene prompt
// 2. Full storyboard → narrative scaffold → video generation prompts
// Replaces scripts/storyboard2video.py — zero external dependencies.
package vision

import (
	"fmt"
	"os"
	"strings"
)

// ─── Narrative Model Definitions ──────────────────────────────────────────

// NarrativeBeat defines a single beat within a narrative model.
type NarrativeBeat struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Mood   string  `json:"mood"`
}

// NarrativeModel defines a narrative structure with named beats.
type NarrativeModel struct {
	Name  string          `json:"name"`
	Beats []NarrativeBeat `json:"beats"`
}

var narrativeModels = map[string]NarrativeModel{
	"three-act": {
		Name: "三幕结构",
		Beats: []NarrativeBeat{
			{ID: "setup", Name: "铺垫", Weight: 0.20, Mood: "平静"},
			{ID: "inciting", Name: "激励事件", Weight: 0.10, Mood: "好奇"},
			{ID: "rising", Name: "上升行动", Weight: 0.35, Mood: "紧张"},
			{ID: "climax", Name: "高潮", Weight: 0.15, Mood: "激动"},
			{ID: "resolution", Name: "结局", Weight: 0.20, Mood: "释然"},
		},
	},
	"heros-journey": {
		Name: "英雄之旅",
		Beats: []NarrativeBeat{
			{ID: "ordinary", Name: "平凡世界", Weight: 0.10, Mood: "平静"},
			{ID: "call", Name: "冒险召唤", Weight: 0.10, Mood: "好奇"},
			{ID: "threshold", Name: "跨越门槛", Weight: 0.15, Mood: "紧张"},
			{ID: "trials", Name: "考验/盟友/敌人", Weight: 0.30, Mood: "紧张"},
			{ID: "ordeal", Name: "最大考验", Weight: 0.10, Mood: "激动"},
			{ID: "reward", Name: "奖赏", Weight: 0.10, Mood: "喜悦"},
			{ID: "return", Name: "携宝而归", Weight: 0.15, Mood: "释然"},
		},
	},
	"problem-solution": {
		Name: "问题→解决方案",
		Beats: []NarrativeBeat{
			{ID: "problem", Name: "问题呈现", Weight: 0.20, Mood: "压抑"},
			{ID: "pain", Name: "痛点放大", Weight: 0.20, Mood: "紧张"},
			{ID: "discovery", Name: "发现方案", Weight: 0.15, Mood: "好奇"},
			{ID: "solution", Name: "方案展示", Weight: 0.30, Mood: "自信"},
			{ID: "result", Name: "成果愿景", Weight: 0.15, Mood: "希望"},
		},
	},
}

// ─── Types ────────────────────────────────────────────────────────────────

// BeatInfo is the runtime struct for a narrative beat with computed values.
type BeatInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Mood          string `json:"mood"`
	StartTime     int    `json:"start_time"`
	Duration      int    `json:"duration"`
	Voiceover     string `json:"voiceover"`
	VisualContext string `json:"visual_context"`
	ImageHint     string `json:"image_hint"`
	ImageCount    int    `json:"image_count"`
	Translation   string `json:"translation"`
}

// VideoPromptInfo is a ComfyUI-compatible video generation prompt.
type VideoPromptInfo struct {
	Scene    int    `json:"scene"`
	BeatID   string `json:"beat_id"`
	Prompt   string `json:"prompt"`
	Negative string `json:"negative"`
	Mood     string `json:"mood"`
	Duration int    `json:"duration"`
	Workflow string `json:"workflow"`
}

// NarrativeScaffold holds the complete narrative structure.
type NarrativeScaffold struct {
	Model         string            `json:"model"`
	ModelName     string            `json:"model_name"`
	TotalDuration int               `json:"total_duration"`
	Beats         []BeatInfo        `json:"beats"`
	VideoPrompts  []VideoPromptInfo `json:"video_prompts"`
	NarrativeText string            `json:"narrative_text"`
}

// ─── Public API ───────────────────────────────────────────────────────────

// AnalyzeStoryboard analyzes multiple storyboard images.
func AnalyzeStoryboard(imagePaths []string, userDesc string) map[string]interface{} {
	frames := make([]map[string]interface{}, 0, len(imagePaths))
	for i, path := range imagePaths {
		analysis := AnalyzeImageToMap(path)
		if err, hasErr := analysis["error"]; hasErr && err != nil {
			analysis = map[string]interface{}{
				"index": i,
				"error": fmt.Sprintf("file not found: %s", path),
			}
		} else {
			prompt := BuildPromptFromAnalysisMap(analysis, userDesc)
			analysis["generated_prompt"] = prompt
			analysis["index"] = i
		}
		frames = append(frames, analysis)
	}

	return map[string]interface{}{
		"total_frames":       len(imagePaths),
		"frames":             frames,
		"narrative_scaffold": nil,
		"video_prompts":      nil,
	}
}

// BuildNarrativeScaffold builds a narrative scaffold from analyzed storyboard images.
func BuildNarrativeScaffold(imageAnalyses []map[string]interface{}, modelID string,
	duration int, userDesc string) *NarrativeScaffold {

	model, ok := narrativeModels[modelID]
	if !ok {
		model = narrativeModels["three-act"]
	}
	totalBeats := len(model.Beats)

	// Distribute images across beats
	remaining := make([]map[string]interface{}, len(imageAnalyses))
	copy(remaining, imageAnalyses)
	imgsPerBeat := len(remaining) / totalBeats
	if imgsPerBeat < 1 {
		imgsPerBeat = 1
	}

	beats := make([]BeatInfo, 0, totalBeats)
	curTime := 0

	for i, beat := range model.Beats {
		beatDur := int(float64(duration) * beat.Weight)
		if beatDur < 3 {
			beatDur = 3
		}
		if i == totalBeats-1 {
			beatDur = duration - curTime
			if beatDur < 3 {
				beatDur = 3
			}
		}

		// Assign images to this beat
		nImages := imgsPerBeat
		if nImages > len(remaining) {
			nImages = len(remaining)
		}
		if i == totalBeats-1 {
			nImages = len(remaining) // last beat gets all remaining
		}
		beatImages := remaining[:nImages]
		remaining = remaining[nImages:]

		// Build visual context from images
		visualContext := ""
		if len(beatImages) > 0 {
			colors := make(map[string]bool)
			stylesSet := make(map[string]bool)
			for _, img := range beatImages {
				if dc, ok := img["dominant_colors"].([]interface{}); ok {
					for j, ci := range dc {
						if j >= 2 {
							break
						}
						if cm, ok := ci.(map[string]interface{}); ok {
							if h, ok := cm["hex"].(string); ok {
								colors[h] = true
							}
						}
					}
				}
				if ss, ok := img["suggested_styles"].([]interface{}); ok {
					for _, s := range ss {
						if str, ok := s.(string); ok {
							stylesSet[str] = true
						}
					}
				}
			}
			parts := []string{}
			if len(colors) > 0 {
				hexes := make([]string, 0, 4)
				for c := range colors {
					hexes = append(hexes, c)
					if len(hexes) >= 4 {
						break
					}
				}
				parts = append(parts, fmt.Sprintf("palette: %s", strings.Join(hexes, ", ")))
			}
			if len(stylesSet) > 0 {
				styles := make([]string, 0, 3)
				for s := range stylesSet {
					styles = append(styles, s)
					if len(styles) >= 3 {
						break
					}
				}
				parts = append(parts, fmt.Sprintf("style: %s", strings.Join(styles, ", ")))
			}
			visualContext = strings.Join(parts, " | ")
		}

		// Image hint from first image's generated prompt
		beatPrompt := ""
		if len(beatImages) > 0 {
			if gp, ok := beatImages[0]["generated_prompt"].(string); ok {
				if len(gp) > 120 {
					beatPrompt = gp[:120]
				} else {
					beatPrompt = gp
				}
			}
		}

		voiceover := fmt.Sprintf("(%s场景)", beat.Name)
		if userDesc != "" {
			trunc := userDesc
			if len(trunc) > 80 {
				trunc = trunc[:80]
			}
			voiceover += " " + trunc
		}

		beats = append(beats, BeatInfo{
			ID:            beat.ID,
			Name:          beat.Name,
			Mood:          beat.Mood,
			StartTime:     curTime,
			Duration:      beatDur,
			Voiceover:     voiceover,
			VisualContext: visualContext,
			ImageHint:     beatPrompt,
			ImageCount:    len(beatImages),
			Translation:   fmt.Sprintf("%s — %s", beat.Name, beat.Mood),
		})

		curTime += beatDur
	}

	// Build video prompts
	videoPrompts := buildVideoPrompts(beats, userDesc)

	// Build narrative text
	narrativeLines := make([]string, 0, len(beats))
	for _, b := range beats {
		line := fmt.Sprintf("## %s (%ds-%ds)\n情绪: %s | %s\n旁白: %s",
			b.Name, b.StartTime, b.StartTime+b.Duration, b.Mood, b.VisualContext, b.Voiceover)
		narrativeLines = append(narrativeLines, line)
	}

	return &NarrativeScaffold{
		Model:         modelID,
		ModelName:     model.Name,
		TotalDuration: duration,
		Beats:         beats,
		VideoPrompts:  videoPrompts,
		NarrativeText: strings.Join(narrativeLines, "\n\n"),
	}
}

// buildVideoPrompts builds ComfyUI-compatible video prompts from beats.
func buildVideoPrompts(beats []BeatInfo, userDesc string) []VideoPromptInfo {
	prompts := make([]VideoPromptInfo, 0, len(beats))
	for i, beat := range beats {
		prefix := "Scene "
		if i == 0 {
			prefix = "Opening: "
		} else {
			prefix = fmt.Sprintf("Scene %d: ", i+1)
		}

		desc := userDesc
		if desc == "" {
			if len(beat.Voiceover) > 60 {
				desc = beat.Voiceover[:60]
			} else {
				desc = beat.Voiceover
			}
		}
		if len(desc) > 60 {
			desc = desc[:60]
		}

		promptText := fmt.Sprintf("%s%s — %s", prefix, beat.Name, desc)
		if beat.ImageHint != "" {
			hint := beat.ImageHint
			if len(hint) > 80 {
				hint = hint[:80]
			}
			promptText += fmt.Sprintf(". Visual reference: %s", hint)
		}

		workflow := "sdxl_txt2img"
		if beat.Duration > 5 {
			workflow = "animatediff_txt2vid"
		}

		prompts = append(prompts, VideoPromptInfo{
			Scene:    i + 1,
			BeatID:   beat.ID,
			Prompt:   promptText,
			Negative: "blurry, low quality, distorted, ugly",
			Mood:     beat.Mood,
			Duration: beat.Duration,
			Workflow: workflow,
		})
	}
	return prompts
}

// ProcessStoryboard is the main entry point: analyze images, build narrative, output structured result.
func ProcessStoryboard(imagePaths []string, model string, duration int, desc string, singleMode bool) map[string]interface{} {
	// Validate paths
	validPaths := make([]string, 0, len(imagePaths))
	for _, p := range imagePaths {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			validPaths = append(validPaths, p)
		}
	}

	if len(validPaths) == 0 {
		return map[string]interface{}{
			"error":         "no valid image files found",
			"paths_checked": imagePaths,
		}
	}

	// Analyze all images
	result := AnalyzeStoryboard(validPaths, desc)

	if singleMode || len(validPaths) == 1 {
		// Single image → generate opening scene prompt
		frames, _ := result["frames"].([]interface{})
		if len(frames) > 0 {
			img, _ := frames[0].(map[string]interface{})
			prompt := desc
			if gp, ok := img["generated_prompt"].(string); ok {
				prompt = gp
			}

			// Build visual context
			hexes := make([]string, 0)
			if dc, ok := img["dominant_colors"].([]interface{}); ok {
				for i, ci := range dc {
					if i >= 5 {
						break
					}
					if cm, ok := ci.(map[string]interface{}); ok {
						if h, ok := cm["hex"].(string); ok {
							hexes = append(hexes, h)
						}
					}
				}
			}
			brightness := ""
			if bl, ok := img["brightness_label"].(string); ok {
				brightness = bl
			}
			moods := []string{}
			if sm, ok := img["suggested_moods"].([]interface{}); ok {
				for _, m := range sm {
					if str, ok := m.(string); ok {
						moods = append(moods, str)
					}
				}
			}

			result["mode"] = "single_image"
			result["opening_prompt"] = prompt
			result["video_prompt"] = map[string]interface{}{
				"scene":    1,
				"prompt":   fmt.Sprintf("Opening scene: %s", truncateStr(prompt, 200)),
				"negative": "blurry, low quality",
				"workflow": "sdxl_txt2img",
				"duration": 5,
				"visual_context": map[string]interface{}{
					"colors":     hexes,
					"brightness": brightness,
					"moods":      moods,
				},
			}
		}
	} else {
		// Multiple images → full storyboard narrative
		frames, _ := result["frames"].([]interface{})
		analyses := make([]map[string]interface{}, 0, len(frames))
		for _, f := range frames {
			if m, ok := f.(map[string]interface{}); ok {
				analyses = append(analyses, m)
			}
		}

		scaffold := BuildNarrativeScaffold(analyses, model, duration, desc)
		result["mode"] = "storyboard"
		result["narrative_scaffold"] = scaffold
		result["video_prompts"] = scaffold.VideoPrompts
	}

	result["config"] = map[string]interface{}{
		"model":       model,
		"duration":    duration,
		"image_count": len(validPaths),
	}

	return result
}

// ─── Helpers ──────────────────────────────────────────────────────────────

// truncateStr truncates a string to maxLen characters.
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
