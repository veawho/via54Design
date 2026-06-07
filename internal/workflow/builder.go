// via54Design — ComfyUI Workflow Builder
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
// overrides can contain keys like "steps", "cfg", "seed", "width", "height",
// "sampler", "scheduler", "denoise", "frames".
func BuildWorkflow(tmpl *WorkflowTemplate, prompt, negativePrompt string, overrides map[string]interface{}) (*BuildResult, error) {
	// Merge template params with overrides
	params := make(map[string]interface{})
	for k, v := range tmpl.Params {
		params[k] = v
	}
	for k, v := range overrides {
		params[k] = v
	}

	// Ensure seed is set
	if seed, ok := params["seed"]; ok {
		if s, ok := seed.(int); ok && s == -1 {
			params["seed"] = rand.Intn(1<<31 - 1)
		}
	}

	nodes := make(map[string]*ComfyUINode)
	nextID := 1

	// Helper to generate sequential node IDs
	nodeID := func() string {
		id := strconv.Itoa(nextID)
		nextID++
		return id
	}

	// Track node IDs for connections
	var (
		checkpointID string
		positiveID   string
		negativeID   string
		samplerID    string
		latentID     string
		vaeDecodeID  string
		saveImageID  string
	)

	switch tmpl.Type {
	case "txt2img":
		// 1. CheckpointLoaderSimple
		checkpointID = nodeID()
		nodes[checkpointID] = &ComfyUINode{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": tmpl.Model,
			},
		}

		// 2. CLIPTextEncode (positive prompt)
		positiveID = nodeID()
		nodes[positiveID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": prompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 3. CLIPTextEncode (negative prompt)
		negativeID = nodeID()
		nodes[negativeID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": negativePrompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 4. EmptyLatentImage
		latentID = nodeID()
		width := toInt(params["width"], 1024)
		height := toInt(params["height"], 1024)
		batchSize := 1
		if bs, ok := params["batch_size"]; ok {
			batchSize = toInt(bs, 1)
		}
		nodes[latentID] = &ComfyUINode{
			ClassType: "EmptyLatentImage",
			Inputs: map[string]interface{}{
				"width":      width,
				"height":     height,
				"batch_size": batchSize,
			},
		}

		// 5. KSampler
		samplerID = nodeID()
		nodes[samplerID] = &ComfyUINode{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":          toInt(params["seed"], 0),
				"steps":         toInt(params["steps"], 30),
				"cfg":           toFloat(params["cfg"], 7.5),
				"sampler_name":  toString(params["sampler"], "euler"),
				"scheduler":     toString(params["scheduler"], "normal"),
				"denoise":       toFloat(params["denoise"], 1.0),
				"model":         []interface{}{checkpointID, 0},
				"positive":      []interface{}{positiveID, 0},
				"negative":      []interface{}{negativeID, 0},
				"latent_image":  []interface{}{latentID, 0},
			},
		}

		// 6. VAEDecode
		vaeDecodeID = nodeID()
		nodes[vaeDecodeID] = &ComfyUINode{
			ClassType: "VAEDecode",
			Inputs: map[string]interface{}{
				"samples": []interface{}{samplerID, 0},
				"vae":     []interface{}{checkpointID, 2},
			},
		}

		// 7. SaveImage
		saveImageID = nodeID()
		nodes[saveImageID] = &ComfyUINode{
			ClassType: "SaveImage",
			Inputs: map[string]interface{}{
				"images":     []interface{}{vaeDecodeID, 0},
				"filename_prefix": "via54",
			},
		}

	case "img2img":
		// 1. LoadImage
		loadImageID := nodeID()
		nodes[loadImageID] = &ComfyUINode{
			ClassType: "LoadImage",
			Inputs: map[string]interface{}{
				"image": "", // user must provide
			},
		}

		// 2. CheckpointLoaderSimple
		checkpointID = nodeID()
		nodes[checkpointID] = &ComfyUINode{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": tmpl.Model,
			},
		}

		// 3. VAEEncode (for img2img input)
		vaeEncodeID := nodeID()
		nodes[vaeEncodeID] = &ComfyUINode{
			ClassType: "VAEEncode",
			Inputs: map[string]interface{}{
				"pixels": []interface{}{loadImageID, 0},
				"vae":    []interface{}{checkpointID, 2},
			},
		}

		// 4. CLIPTextEncode (positive prompt)
		positiveID = nodeID()
		nodes[positiveID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": prompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 5. CLIPTextEncode (negative prompt)
		negativeID = nodeID()
		nodes[negativeID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": negativePrompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 6. KSampler
		samplerID = nodeID()
		nodes[samplerID] = &ComfyUINode{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":          toInt(params["seed"], 0),
				"steps":         toInt(params["steps"], 30),
				"cfg":           toFloat(params["cfg"], 7.5),
				"sampler_name":  toString(params["sampler"], "euler"),
				"scheduler":     toString(params["scheduler"], "normal"),
				"denoise":       toFloat(params["denoise"], 0.6),
				"model":         []interface{}{checkpointID, 0},
				"positive":      []interface{}{positiveID, 0},
				"negative":      []interface{}{negativeID, 0},
				"latent_image":  []interface{}{vaeEncodeID, 0},
			},
		}

		// 7. VAEDecode
		vaeDecodeID = nodeID()
		nodes[vaeDecodeID] = &ComfyUINode{
			ClassType: "VAEDecode",
			Inputs: map[string]interface{}{
				"samples": []interface{}{samplerID, 0},
				"vae":     []interface{}{checkpointID, 2},
			},
		}

		// 8. SaveImage
		saveImageID = nodeID()
		nodes[saveImageID] = &ComfyUINode{
			ClassType: "SaveImage",
			Inputs: map[string]interface{}{
				"images":          []interface{}{vaeDecodeID, 0},
				"filename_prefix": "via54",
			},
		}
		_ = vaeEncodeID
		_ = loadImageID

	case "txt2vid":
		// 1. CheckpointLoaderSimple
		checkpointID = nodeID()
		nodes[checkpointID] = &ComfyUINode{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": tmpl.Model,
			},
		}

		// 2. CLIPTextEncode (positive prompt)
		positiveID = nodeID()
		nodes[positiveID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": prompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 3. CLIPTextEncode (negative prompt)
		negativeID = nodeID()
		nodes[negativeID] = &ComfyUINode{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": negativePrompt,
				"clip": []interface{}{checkpointID, 1},
			},
		}

		// 4. EmptyLatentImage (with batch_size = frames)
		latentID = nodeID()
		frames := toInt(params["frames"], 16)
		width := toInt(params["width"], 512)
		height := toInt(params["height"], 512)
		nodes[latentID] = &ComfyUINode{
			ClassType: "EmptyLatentImage",
			Inputs: map[string]interface{}{
				"width":      width,
				"height":     height,
				"batch_size": frames,
			},
		}

		// 5. KSampler
		samplerID = nodeID()
		nodes[samplerID] = &ComfyUINode{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":          toInt(params["seed"], 0),
				"steps":         toInt(params["steps"], 25),
				"cfg":           toFloat(params["cfg"], 7.0),
				"sampler_name":  toString(params["sampler"], "euler"),
				"scheduler":     toString(params["scheduler"], "normal"),
				"denoise":       toFloat(params["denoise"], 1.0),
				"model":         []interface{}{checkpointID, 0},
				"positive":      []interface{}{positiveID, 0},
				"negative":      []interface{}{negativeID, 0},
				"latent_image":  []interface{}{latentID, 0},
			},
		}

		// 6. VAEDecode
		vaeDecodeID = nodeID()
		nodes[vaeDecodeID] = &ComfyUINode{
			ClassType: "VAEDecode",
			Inputs: map[string]interface{}{
				"samples": []interface{}{samplerID, 0},
				"vae":     []interface{}{checkpointID, 2},
			},
		}

		// 7. SaveImage (video frames will be saved as frame sequence)
		saveImageID = nodeID()
		nodes[saveImageID] = &ComfyUINode{
			ClassType: "SaveImage",
			Inputs: map[string]interface{}{
				"images":          []interface{}{vaeDecodeID, 0},
				"filename_prefix": "via54_video",
			},
		}

	default:
		return nil, fmt.Errorf("unsupported workflow type: %s", tmpl.Type)
	}

	// Build the JSON as a map (ComfyUI expects string keys for node IDs)
	wfMap := make(map[string]*ComfyUINode)
	// Sort keys for deterministic output
	var keys []string
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		wfMap[k] = nodes[k]
	}

	jsonData, err := json.MarshalIndent(wfMap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal workflow: %w", err)
	}

	// Count injected prompts
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

func toString(v interface{}, defaultVal string) string {
	if s, ok := v.(string); ok {
		return s
	}
	return defaultVal
}
