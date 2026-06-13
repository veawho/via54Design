// via54Design — ComfyUI Workflow Builder (数据驱动 v2)
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package workflow

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

// BuildWorkflow converts a WorkflowTemplate + prompts into ComfyUI API JSON.
// The template skeleton is a valid ComfyUI JSON with __PROMPT__ / __NEGATIVE__ placeholders.
// overrides: steps, cfg, seed, width, height, sampler, scheduler, denoise, frames.
// keyframes: optional time-varying prompts for video (each has Frame + Prompt).
func BuildWorkflow(tmpl *WorkflowTemplate, prompt, negativePrompt string, overrides map[string]interface{}, keyframes []Keyframe) (*BuildResult, error) {
	if tmpl == nil || tmpl.Skeleton == nil {
		return nil, fmt.Errorf("workflow template %q has no skeleton", tmpl.ID)
	}

	// Deep-copy the skeleton
	skJSON, err := json.Marshal(tmpl.Skeleton)
	if err != nil {
		return nil, fmt.Errorf("marshal skeleton: %w", err)
	}
	var wfMap map[string]interface{}
	if err := json.Unmarshal(skJSON, &wfMap); err != nil {
		return nil, fmt.Errorf("unmarshal skeleton: %w", err)
	}

	// Merge params with overrides
	params := make(map[string]interface{})
	for k, v := range tmpl.Params {
		params[k] = v
	}
	for k, v := range overrides {
		params[k] = v
	}
	if s, ok := params["seed"]; ok {
		if si, ok := s.(int); ok && si == -1 {
			params["seed"] = rand.Intn(1<<31 - 1)
		}
	}
	if _, ok := params["width"]; !ok {
		params["width"] = 1024
	}
	if _, ok := params["height"]; !ok {
		params["height"] = 1024
	}

	// First pass: collect sorted keys
	var oldKeys []string
	for k := range wfMap {
		oldKeys = append(oldKeys, k)
	}
	sort.Strings(oldKeys)

	// Renumber nodes sequentially
	nodeIDMap := make(map[string]string) // old → new
	nextID := 1
	newWf := make(map[string]interface{})

	for _, oldKey := range oldKeys {
		rawNode := wfMap[oldKey].(map[string]interface{})
		classType, _ := rawNode["class_type"].(string)
		inputs, _ := rawNode["inputs"].(map[string]interface{})
		newInputs := make(map[string]interface{})

		if inputs != nil {
			for ik, iv := range inputs {
				switch v := iv.(type) {
				case string:
					switch v {
					case "__PROMPT__":
						newInputs[ik] = prompt
					case "__NEGATIVE__":
						newInputs[ik] = negativePrompt
					case "__MODEL__":
						newInputs[ik] = tmpl.Model
					default:
						newInputs[ik] = iv
					}
				case []interface{}:
					newInputs[ik] = iv // keep connection ref as-is
				default:
					newInputs[ik] = iv
				}
			}

			// Apply overrides to KSampler
			if classType == "KSampler" || classType == "KSamplerAdvanced" {
				if v, ok := params["seed"]; ok {
					newInputs["seed"] = toInt(v, 0)
				}
				if v, ok := params["steps"]; ok {
					newInputs["steps"] = toInt(v, 30)
				}
				if v, ok := params["cfg"]; ok {
					newInputs["cfg"] = toFloat(v, 7.5)
				}
				if v, ok := params["sampler"]; ok {
					newInputs["sampler_name"] = v
				}
				if v, ok := params["scheduler"]; ok {
					newInputs["scheduler"] = v
				}
				if v, ok := params["denoise"]; ok {
					newInputs["denoise"] = toFloat(v, 1.0)
				}
			}

			// Apply overrides to EmptyLatentImage
			if classType == "EmptyLatentImage" {
				if v, ok := params["width"]; ok {
					newInputs["width"] = toInt(v, 1024)
				}
				if v, ok := params["height"]; ok {
					newInputs["height"] = toInt(v, 1024)
				}
				if v, ok := params["frames"]; ok {
					newInputs["batch_size"] = toInt(v, 1)
				}
			}
		}

		node := map[string]interface{}{
			"class_type": classType,
			"inputs":     newInputs,
		}
		newKey := strconv.Itoa(nextID)
		nodeIDMap[oldKey] = newKey
		newWf[newKey] = node
		nextID++
	}

	// Second pass: remap connection references
	for _, rawNode := range newWf {
		node := rawNode.(map[string]interface{})
		inputs, _ := node["inputs"].(map[string]interface{})
		if inputs == nil {
			continue
		}
		for ik, iv := range inputs {
			arr, ok := iv.([]interface{})
			if !ok || len(arr) != 2 {
				continue
			}
			refID, ok := arr[0].(string)
			if !ok {
				continue
			}
			if newID, exists := nodeIDMap[refID]; exists {
				inputs[ik] = []interface{}{newID, arr[1]}
			}
		}
	}

	// Handle keyframe scheduling (prompt travel for video)
	kfCount := 0
	if len(keyframes) > 0 {
		// Sort keyframes by frame
		sort.Slice(keyframes, func(i, j int) bool { return keyframes[i].Frame < keyframes[j].Frame })

		// Create a BatchPromptSchedule node
		// Format: "0" → "prompt at frame 0", "8" → "prompt at frame 8"
		scheduleStr := ""
		for i, kf := range keyframes {
			if i > 0 {
				scheduleStr += "\n"
			}
			scheduleStr += fmt.Sprintf("\"%d\"", kf.Frame) + ": \"" + kf.Prompt + "\","
		}
		// Fallback: if negative prompt node exists, use it for pre_text
		preText := ""
		if negativePrompt != "" {
			preText = negativePrompt
		}

		kfNode := map[string]interface{}{
			"class_type": "BatchPromptSchedule",
			"inputs": map[string]interface{}{
				"text":          scheduleStr,
				"pre_text":      preText,
				"appendix":      "",
				"start_frame":   0,
				"end_frame":     toInt(params["frames"], 16),
				"frame_rate":    24,
				"interpolation": "linear",
			},
		}
		kfKey := strconv.Itoa(nextID)
		newWf[kfKey] = kfNode
		nextID++
		kfCount = len(keyframes)
	}

	jsonData, err := json.MarshalIndent(newWf, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal workflow: %w", err)
	}

	injected := 0
	if prompt != "" {
		injected++
	}
	if negativePrompt != "" {
		injected++
	}

	return &BuildResult{
		JSON:       jsonData,
		Injected:   injected,
		TemplateID: tmpl.ID,
		Keyframes:  kfCount,
	}, nil
}

// ─── Type Helpers ───

func toInt(v interface{}, defaultVal int) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	case int64:
		return int(val)
	}
	return defaultVal
}

func toFloat(v interface{}, defaultVal float64) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}
