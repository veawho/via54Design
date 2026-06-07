// via54Design — ComfyUI Workflow Tests
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWorkflowTemplate(t *testing.T) {
	// Find the project base directory
	baseDir := findBaseDir(t)
	t.Logf("baseDir: %s", baseDir)

	tmpl, err := LoadWorkflowTemplate("sdxl_txt2img", baseDir)
	if err != nil {
		t.Fatalf("failed to load sdxl_txt2img: %v", err)
	}
	if tmpl.ID != "sdxl_txt2img" {
		t.Errorf("expected id sdxl_txt2img, got %s", tmpl.ID)
	}
	if tmpl.Type != "txt2img" {
		t.Errorf("expected type txt2img, got %s", tmpl.Type)
	}
}

func TestLoadWorkflowTemplate_NotFound(t *testing.T) {
	baseDir := findBaseDir(t)
	_, err := LoadWorkflowTemplate("nonexistent", baseDir)
	if err == nil {
		t.Fatal("expected error for nonexistent template")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestListWorkflowTemplates(t *testing.T) {
	baseDir := findBaseDir(t)
	ids, err := ListWorkflowTemplates(baseDir)
	if err != nil {
		t.Fatalf("ListWorkflowTemplates failed: %v", err)
	}
	if len(ids) == 0 {
		t.Fatal("expected at least 1 workflow template")
	}
	t.Logf("found %d workflow templates: %v", len(ids), ids)

	// Should contain expected templates
	expected := []string{"sdxl_txt2img", "sdxl_img2img", "flux_dev_txt2img", "sd15_txt2img", "animatediff_txt2vid", "wan_txt2vid"}
	for _, e := range expected {
		found := false
		for _, id := range ids {
			if id == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected template %q not found in list", e)
		}
	}
}

func TestLoadRegistry(t *testing.T) {
	baseDir := findBaseDir(t)
	reg, err := LoadRegistry(baseDir)
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}
	if reg.Version != 2 {
		t.Errorf("expected version 2, got %d", reg.Version)
	}
	if len(reg.Workflows) == 0 {
		t.Fatal("expected at least 1 registry entry")
	}
}

func TestBuildWorkflow_Txt2Img(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("sdxl_txt2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	result, err := BuildWorkflow(tmpl, "a cat in a hat", "ugly, blurry", nil, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	// Validate JSON output
	var nodes map[string]interface{}
	if err := json.Unmarshal(result.JSON, &nodes); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	// Should have 7 nodes for txt2img
	if len(nodes) != 7 {
		t.Errorf("expected 7 nodes, got %d", len(nodes))
	}

	// Check class types
	classes := make(map[string]bool)
	for _, n := range nodes {
		node := n.(map[string]interface{})
		classes[node["class_type"].(string)] = true
	}

	expectedClasses := []string{"CheckpointLoaderSimple", "CLIPTextEncode", "EmptyLatentImage", "KSampler", "VAEDecode", "SaveImage"}
	for _, c := range expectedClasses {
		if !classes[c] {
			t.Errorf("missing node class: %s", c)
		}
	}

	// Verify prompt was injected
	promptCount := 0
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if node["class_type"] == "CLIPTextEncode" {
			inputs := node["inputs"].(map[string]interface{})
			text := inputs["text"].(string)
			if text == "a cat in a hat" || text == "ugly, blurry" {
				promptCount++
			}
		}
	}
	if promptCount < 2 {
		t.Errorf("expected at least 2 CLIPTextEncode with injected prompts, got %d", promptCount)
	}
}

func TestBuildWorkflow_Img2Img(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("sdxl_img2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	result, err := BuildWorkflow(tmpl, "a cat", "ugly", nil, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	var nodes map[string]interface{}
	if err := json.Unmarshal(result.JSON, &nodes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// img2img should have 8 nodes (LoadImage + VAEEncode + standard txt2img nodes)
	if len(nodes) != 8 {
		t.Errorf("expected 8 nodes for img2img, got %d", len(nodes))
	}

	classes := make(map[string]bool)
	for _, n := range nodes {
		node := n.(map[string]interface{})
		classes[node["class_type"].(string)] = true
	}

	if !classes["LoadImage"] {
		t.Error("img2img missing LoadImage node")
	}
	if !classes["VAEEncode"] {
		t.Error("img2img missing VAEEncode node")
	}
}

func TestBuildWorkflow_Txt2Vid(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("animatediff_txt2vid", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	result, err := BuildWorkflow(tmpl, "a cat walking", "ugly", nil, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	var nodes map[string]interface{}
	if err := json.Unmarshal(result.JSON, &nodes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// txt2vid should have 7 nodes like txt2img
	if len(nodes) != 7 {
		t.Errorf("expected 7 nodes for txt2vid, got %d", len(nodes))
	}
}

func TestBuildWorkflow_Overrides(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("sd15_txt2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	overrides := map[string]interface{}{
		"steps":    40,
		"cfg":      10.0,
		"seed":     12345,
		"width":    768,
		"height":   768,
	}
	result, err := BuildWorkflow(tmpl, "test", "", overrides, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	var nodes map[string]interface{}
	if err := json.Unmarshal(result.JSON, &nodes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Find KSampler and check overrides
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if node["class_type"] == "KSampler" {
			inputs := node["inputs"].(map[string]interface{})
			if inputs["steps"].(float64) != 40 {
				t.Errorf("expected steps=40, got %v", inputs["steps"])
			}
			if inputs["cfg"].(float64) != 10.0 {
				t.Errorf("expected cfg=10.0, got %v", inputs["cfg"])
			}
		}
	}
}

func TestBuildWorkflow_FluxDev(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("flux_dev_txt2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	if tmpl.Model != "flux_dev.safetensors" {
		t.Errorf("expected flux_dev.safetensors, got %s", tmpl.Model)
	}

	result, err := BuildWorkflow(tmpl, "test", "", map[string]interface{}{
		"guidance": 3.5,
	}, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	var nodes map[string]interface{}
	if err := json.Unmarshal(result.JSON, &nodes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(nodes) != 7 {
		t.Errorf("expected 7 nodes, got %d", len(nodes))
	}
}

func TestBuildWorkflow_DeterministicSeed(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("sdxl_txt2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	// Using a fixed seed should give deterministic output
	overrides := map[string]interface{}{"seed": 42}
	result1, _ := BuildWorkflow(tmpl, "cat", "", overrides, nil)
	result2, _ := BuildWorkflow(tmpl, "cat", "", overrides, nil)

	if string(result1.JSON) != string(result2.JSON) {
		t.Error("same inputs should produce identical JSON")
	}
}

func TestNegativePrompt(t *testing.T) {
	baseDir := findBaseDir(t)
	tmpl, err := LoadWorkflowTemplate("sdxl_txt2img", baseDir)
	if err != nil {
		t.Fatalf("load template: %v", err)
	}

	result, err := BuildWorkflow(tmpl, "cat", "bad quality, ugly", nil, nil)
	if err != nil {
		t.Fatalf("BuildWorkflow failed: %v", err)
	}

	var nodes map[string]interface{}
	json.Unmarshal(result.JSON, &nodes)

	// Find negative CLIPTextEncode
	negFound := false
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if node["class_type"] == "CLIPTextEncode" {
			inputs := node["inputs"].(map[string]interface{})
			if inputs["text"] == "bad quality, ugly" {
				negFound = true
				break
			}
		}
	}
	if !negFound {
		t.Error("negative prompt not found in output")
	}
}

// findBaseDir walks up from common paths to find the project root.
func findBaseDir(t *testing.T) string {
	t.Helper()
	candidates := []string{
		".",
		"..",
		"../..",
		os.Getenv("HOME") + "/AppData/Local/Temp/via54Design",
		"C:/Users/via54/AppData/Local/Temp/via54Design",
		"/c/Users/via54/AppData/Local/Temp/via54Design",
	}
	for _, d := range candidates {
		abs, _ := filepath.Abs(d)
		if _, err := os.Stat(filepath.Join(abs, "templates", "workflows", "sdxl_txt2img.yaml")); err == nil {
			return abs
		}
	}
	t.Fatal("could not find project base directory")
	return ""
}
