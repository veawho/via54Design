// via54Design — ComfyUI Workflow Template Loader
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadWorkflowTemplate loads a single workflow template by ID.
func LoadWorkflowTemplate(id, baseDir string) (*WorkflowTemplate, error) {
	candidates := []string{
		filepath.Join(baseDir, "templates", "workflows", id+".yaml"),
		filepath.Join(baseDir, "templates", "workflows", id+".yml"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var t WorkflowTemplate
		if err := yaml.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("parse workflow %s: %w", id, err)
		}
		if t.ID == "" {
			return nil, fmt.Errorf("workflow %s has empty id", id)
		}
		// Load skeleton JSON if available
		skeletonPath := filepath.Join(baseDir, "templates", "workflows", id+".skeleton.json")
		if skData, skErr := os.ReadFile(skeletonPath); skErr == nil {
			var skeleton map[string]interface{}
			if err := json.Unmarshal(skData, &skeleton); err == nil {
				t.Skeleton = skeleton
			}
		}
		if t.Skeleton == nil {
			// Fallback: generate skeleton from YAML nodes list
			t.Skeleton = generateSkeleton(&t)
		}
		return &t, nil
	}
	return nil, fmt.Errorf("workflow template %q not found in templates/workflows/", id)
}

// ListWorkflowTemplates scans the workflows directory and returns all template IDs.
func ListWorkflowTemplates(baseDir string) ([]string, error) {
	dir := filepath.Join(baseDir, "templates", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read workflows directory: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "registry.yaml" || name == "registry.yml" {
			continue
		}
		if filepath.Ext(name) == ".yaml" || filepath.Ext(name) == ".yml" {
			id := name[:len(name)-len(filepath.Ext(name))]
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// LoadRegistry loads the workflow registry file.
func LoadRegistry(baseDir string) (*WorkflowRegistry, error) {
	candidates := []string{
		filepath.Join(baseDir, "templates", "workflows", "registry.yaml"),
		filepath.Join(baseDir, "templates", "workflows", "registry.yml"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var reg WorkflowRegistry
		if err := yaml.Unmarshal(data, &reg); err != nil {
			return nil, fmt.Errorf("parse registry: %w", err)
		}
		return &reg, nil
	}
	return nil, fmt.Errorf("workflow registry not found in templates/workflows/")
}

// generateSkeleton builds a minimal node graph from the YAML template's Nodes map.
// This is the fallback when no .skeleton.json exists.
func generateSkeleton(t *WorkflowTemplate) map[string]interface{} {
	skeleton := make(map[string]interface{})
	// Build txt2img skeleton
	if t.Type == "txt2img" || t.Type == "txt2vid" {
		// checkpoint
		skeleton["1"] = map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{"ckpt_name": t.Model},
		}
		// positive prompt
		skeleton["2"] = map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{"text": "__PROMPT__", "clip": []interface{}{"1", 1}},
		}
		// negative prompt
		skeleton["3"] = map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{"text": "__NEGATIVE__", "clip": []interface{}{"1", 1}},
		}
		// latent
		skeleton["4"] = map[string]interface{}{
			"class_type": "EmptyLatentImage",
			"inputs": map[string]interface{}{"width": 1024, "height": 1024, "batch_size": 1},
		}
		// sampler
		skeleton["5"] = map[string]interface{}{
			"class_type": "KSampler",
			"inputs": map[string]interface{}{
				"seed": 0, "steps": 30, "cfg": 7.5, "sampler_name": "euler",
				"scheduler": "normal", "denoise": 1.0,
				"model": []interface{}{"1", 0}, "positive": []interface{}{"2", 0},
				"negative": []interface{}{"3", 0}, "latent_image": []interface{}{"4", 0},
			},
		}
		// vae decode
		skeleton["6"] = map[string]interface{}{
			"class_type": "VAEDecode",
			"inputs": map[string]interface{}{"samples": []interface{}{"5", 0}, "vae": []interface{}{"1", 2}},
		}
		// save
		skeleton["7"] = map[string]interface{}{
			"class_type": "SaveImage",
			"inputs": map[string]interface{}{"images": []interface{}{"6", 0}, "filename_prefix": "via54"},
		}
	}
	if t.Type == "img2img" {
		skeleton["1"] = map[string]interface{}{"class_type": "LoadImage", "inputs": map[string]interface{}{"image": ""}}
		skeleton["2"] = map[string]interface{}{"class_type": "CheckpointLoaderSimple", "inputs": map[string]interface{}{"ckpt_name": t.Model}}
		skeleton["3"] = map[string]interface{}{"class_type": "VAEEncode", "inputs": map[string]interface{}{"pixels": []interface{}{"1", 0}, "vae": []interface{}{"2", 2}}}
		skeleton["4"] = map[string]interface{}{"class_type": "CLIPTextEncode", "inputs": map[string]interface{}{"text": "__PROMPT__", "clip": []interface{}{"2", 1}}}
		skeleton["5"] = map[string]interface{}{"class_type": "CLIPTextEncode", "inputs": map[string]interface{}{"text": "__NEGATIVE__", "clip": []interface{}{"2", 1}}}
		skeleton["6"] = map[string]interface{}{"class_type": "KSampler", "inputs": map[string]interface{}{
			"seed": 0, "steps": 30, "cfg": 7.5, "sampler_name": "euler",
			"scheduler": "normal", "denoise": 0.6,
			"model": []interface{}{"2", 0}, "positive": []interface{}{"4", 0},
			"negative": []interface{}{"5", 0}, "latent_image": []interface{}{"3", 0},
		}}
		skeleton["7"] = map[string]interface{}{"class_type": "VAEDecode", "inputs": map[string]interface{}{"samples": []interface{}{"6", 0}, "vae": []interface{}{"2", 2}}}
		skeleton["8"] = map[string]interface{}{"class_type": "SaveImage", "inputs": map[string]interface{}{"images": []interface{}{"7", 0}, "filename_prefix": "via54"}}
	}
	return skeleton
}
