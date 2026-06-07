// via54Design — Web UI HTTP Handlers (全功能版)
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: MIT

package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/veawho/via54Design/internal/workflow"
)

//go:embed templates/index.html
var embeddedFiles embed.FS

var baseDir string

func init() {
	candidates := []string{".", "..", "../.."}
	for _, c := range candidates {
		path := filepath.Join(c, "templates", "workflows")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)
			baseDir = abs
			return
		}
	}
	baseDir, _ = os.Getwd()
}

func Handler(bd string) http.Handler {
	if bd != "" {
		baseDir = bd
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/templates", handleTemplates)
	mux.HandleFunc("/api/build", handleBuild)
	mux.HandleFunc("/api/prompt", handlePrompt)
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/narrate", handleNarrate)
	mux.HandleFunc("/api/export", handleExport)
	mux.HandleFunc("/api/media", handleMedia)
	return mux
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := embeddedFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok", "version": "v0.6.0",
	})
}

func handleTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reg, err := workflow.LoadRegistry(baseDir)
	if err != nil {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	json.NewEncoder(w).Encode(reg.Workflows)
}

func handleBuild(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	workflowID, _ := req["workflow_id"].(string)
	prompt, _ := req["prompt"].(string)
	neg, _ := req["negative"].(string)
	format, _ := req["format"].(string)
	if format == "" {
		format = "comfyui"
	}

	if workflowID == "" || prompt == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "workflow_id and prompt required"})
		return
	}

	overrides := map[string]interface{}{}
	if v, ok := req["steps"].(float64); ok {
		overrides["steps"] = int(v)
	}
	if v, ok := req["cfg"].(float64); ok {
		overrides["cfg"] = v
	}
	if v, ok := req["seed"].(float64); ok {
		overrides["seed"] = int(v)
	}
	if v, ok := req["width"].(float64); ok {
		overrides["width"] = int(v)
	}
	if v, ok := req["height"].(float64); ok {
		overrides["height"] = int(v)
	}

	if format == "forge" {
		fp := map[string]interface{}{
			"prompt": prompt, "negative_prompt": neg,
			"steps": 30, "cfg_scale": 7.5, "width": 1024, "height": 1024,
			"sampler_name": "Euler", "save_images": true,
		}
		if v, ok := overrides["steps"].(int); ok {
			fp["steps"] = v
		}
		if v, ok := overrides["cfg"].(float64); ok {
			fp["cfg_scale"] = v
		}
		if v, ok := overrides["seed"].(int); ok {
			fp["seed"] = v
		}
		res := map[string]interface{}{
			"format": "forge_a1111", "template": workflowID,
			"api_payload": fp, "api_endpoint": "http://localhost:7860/sdapi/v1/txt2img",
		}
		j, _ := json.MarshalIndent(res, "", "  ")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"json": string(j), "nodes": 1, "format": "forge",
		})
		return
	}

	// ComfyUI mode
	tmpl, err := workflow.LoadWorkflowTemplate(workflowID, baseDir)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var kfs []workflow.Keyframe
	result, err := workflow.BuildWorkflow(tmpl, prompt, neg, overrides, kfs)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"json": string(result.JSON), "nodes": len(strings.Split(string(result.JSON), "\n")) / 3,
		"template": result.TemplateID, "format": "comfyui",
	})
}

func handlePrompt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	scene, _ := req["scene"].(string)
	platform, _ := req["platform"].(string)
	outFormat, _ := req["format"].(string)
	if scene == "" || platform == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "scene and platform required"})
		return
	}
	if outFormat == "" {
		outFormat = "markdown"
	}

	// Execute via54 prompt CLI
	exe := filepath.Join(baseDir, "via54.exe")
	args := []string{"prompt", "--scene", scene, "--platform", platform, "--format", outFormat}
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("exec error: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"output":   string(out),
		"platform": platform,
		"format":   outFormat,
		"length":   len(out),
	})
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	layout, _ := req["layout"].(string)
	color, _ := req["color"].(string)
	font, _ := req["font"].(string)
	title, _ := req["title"].(string)
	presentation, _ := req["presentation"].(bool)
	if layout == "" {
		layout = "hero-split-16-9"
	}
	if color == "" {
		color = "ink-wash"
	}
	if font == "" {
		font = "ming-hei-editorial"
	}
	if title == "" {
		title = "via54Design"
	}

	exe := filepath.Join(baseDir, "via54.exe")
	args := []string{
		"generate", "--layout", layout, "--color", color,
		"--font", font, "--title", title,
	}
	if presentation {
		args = append(args, "--presentation")
	}
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("exec error: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"html": string(out), "layout": layout, "color": color,
		"font": font, "title": title, "length": len(out),
	})
}

func handleNarrate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	seed, _ := req["seed"].(string)
	model, _ := req["model"].(string)
	outFormat, _ := req["format"].(string)
	duration := 30
	if v, ok := req["duration"].(float64); ok {
		duration = int(v)
	}
	if seed == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "seed required"})
		return
	}
	if model == "" {
		model = "three-act"
	}
	if outFormat == "" {
		outFormat = "markdown"
	}

	exe := filepath.Join(baseDir, "via54.exe")
	args := []string{
		"narrate", "--seed", seed, "--model", model,
		"--duration", strconv.Itoa(duration), "--format", outFormat,
	}
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("exec error: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"output": string(out), "model": model,
		"duration": duration, "length": len(out),
	})
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	expType, _ := req["type"].(string)
	source, _ := req["source"].(string)
	output, _ := req["output"].(string)
	if expType == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "export type required"})
		return
	}
	if output == "" {
		output = fmt.Sprintf("output.%s", expType)
	}

	exe := filepath.Join(baseDir, "via54.exe")
	args := []string{"export", expType}
	if source != "" {
		args = append(args, source)
	}
	args = append(args, "--output", output)
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("exec error: %v", err),
			"output": string(out),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": string(out), "type": expType,
		"output": output, "success": true,
	})
}

func handleMedia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var req map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	action, _ := req["action"].(string)
	source, _ := req["source"].(string)
	target, _ := req["target"].(string)
	paramsStr, _ := req["params"].(string)
	if action == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "action required"})
		return
	}

	exe := filepath.Join(baseDir, "via54.exe")
	args := []string{"media", action}
	if source != "" {
		args = append(args, source)
	}
	if target != "" {
		args = append(args, target)
	}
	if paramsStr != "" {
		args = append(args, paramsStr)
	}
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("exec error: %v", err), "output": string(out),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": string(out), "action": action, "success": true,
	})
}
