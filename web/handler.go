// via54Design — Web UI HTTP Handlers (全功能版)
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: MIT

package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/veawho/via54Design/internal/vision"
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
	mux.HandleFunc("/api/upload", handleUpload)
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/img2prompt", handleImg2Prompt)
	mux.HandleFunc("/api/regen", handleRegen)
	mux.HandleFunc("/api/storyboard", handleStoryboard)
	mux.HandleFunc("/api/video-prompt", handleVideoPrompt)
	mux.HandleFunc("/api/story2ppt", handleStory2PPT)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir()))))
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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

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

// ═══════════════════════════════════════════════════════
// Image Upload + Analysis + Prompt + Regeneration
// ═══════════════════════════════════════════════════════

func uploadDir() string {
	d := filepath.Join(baseDir, "web", "uploads")
	os.MkdirAll(d, 0755)
	return d
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("handleUpload: parse multipart: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("parse: %v", err)})
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("handleUpload: form file: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("file: %v", err)})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" && ext != ".gif" {
		json.NewEncoder(w).Encode(map[string]string{"error": "unsupported format, use png/jpg/webp/gif"})
		return
	}

	filename := fmt.Sprintf("img_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir(), filename)
	out, err := os.Create(dst)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("create: %v", err)})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("write file: %v", err)})
		return
	}

	url := fmt.Sprintf("/uploads/%s", filename)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url": url, "path": dst, "filename": filename, "size": header.Size,
	})
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	var req map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

	imgPath, _ := req["path"].(string)
	if imgPath == "" {
		log.Printf("handleAnalyze: path required")
		json.NewEncoder(w).Encode(map[string]string{"error": "path required"})
		return
	}
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(baseDir, imgPath)
	}
	// Path traversal prevention
	cleanPath := filepath.Clean(imgPath)
	cleanBase := filepath.Clean(baseDir)
	if cleanPath != cleanBase && !strings.HasPrefix(cleanPath, cleanBase+string(filepath.Separator)) {
		log.Printf("handleAnalyze: path traversal blocked: %s", imgPath)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path"})
		return
	}

	result := vision.AnalyzeImageToMap(imgPath)
	if e, ok := result["error"]; ok && e != nil {
		errStr, _ := e.(string)
		log.Printf("handleAnalyze: analysis error: %s", errStr)
		json.NewEncoder(w).Encode(result)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func handleImg2Prompt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	var req map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

	imgPath, _ := req["path"].(string)
	userDesc, _ := req["desc"].(string)
	platform, _ := req["platform"].(string)
	if platform == "" {
		platform = "flux"
	}
	if imgPath == "" {
		log.Printf("handleImg2Prompt: path required")
		json.NewEncoder(w).Encode(map[string]string{"error": "path required"})
		return
	}
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(baseDir, imgPath)
	}
	// Path traversal prevention
	cleanPath := filepath.Clean(imgPath)
	cleanBase := filepath.Clean(baseDir)
	if cleanPath != cleanBase && !strings.HasPrefix(cleanPath, cleanBase+string(filepath.Separator)) {
		log.Printf("handleImg2Prompt: path traversal blocked: %s", imgPath)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path"})
		return
	}

	analysis := vision.AnalyzeImageToMap(imgPath)
	if e, ok := analysis["error"]; ok && e != nil {
		log.Printf("handleImg2Prompt: analysis error: %v", e)
		json.NewEncoder(w).Encode(analysis)
		return
	}

	// Generate the base prompt from image analysis
	scene := vision.BuildPromptFromAnalysisMap(analysis, userDesc)

	cli := filepath.Join(baseDir, "via54.exe")
	args := []string{"prompt", "--scene", scene, "--platform", platform, "--format", "markdown"}
	cmd2 := exec.Command(cli, args...)
	cmd2.Dir = baseDir
	promptOut, err := cmd2.Output()
	promptText := string(promptOut)
	if err != nil {
		promptText = scene
	}

	analysis["platform"] = platform
	analysis["final_prompt"] = promptText
	analysis["raw_analysis"] = analysis["generated_prompt"]
	delete(analysis, "generated_prompt")

	json.NewEncoder(w).Encode(analysis)
}

func handleRegen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	var req map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

	prompt, _ := req["prompt"].(string)
	negative, _ := req["negative"].(string)
	workflowID, _ := req["workflow"].(string)
	if workflowID == "" {
		workflowID = "sdxl_txt2img"
	}
	if prompt == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "prompt required"})
		return
	}

	payload := map[string]interface{}{
		"prompt": prompt, "negative_prompt": negative,
		"steps": 30, "cfg_scale": 7.5, "width": 1024, "height": 1024,
		"sampler_name": "Euler", "save_images": true,
	}

	// Graceful degradation: always returns a result, even when Forge is down
	// Core function: generate the prompt + instructions regardless of backend
	info := fmt.Sprintf("✅ Prompt 已就绪 (%d chars)\n", len(prompt))
	info += "📤 提交: via54 forge --workflow " + workflowID + " --prompt \"...\" --send\n"
	info += "📋 或复制下方 JSON 手动 POST\n\n"

	payloadBytes, _ := json.MarshalIndent(payload, "", "  ")
	info += string(payloadBytes)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("http://localhost:7860/sdapi/v1/txt2img", "application/json", bytes.NewReader(payloadBytes))

	if err != nil {
		// Forge not available - not an error, just return the prompt info
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":       true,
			"submitted":     false,
			"workflow":      workflowID,
			"message":       "后端（Forge）未连接，提示词已就绪可手动提交",
			"hint":          "启动 Forge 后重试，或复制下方 payload 手动 POST",
			"prompt_info":   info,
			"api_payload":   payload,
			"api_endpoint":  "http://localhost:7860/sdapi/v1/txt2img",
			"offline_mode":  true,
		})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var forgeResult map[string]interface{}
	json.Unmarshal(respBody, &forgeResult)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true, "submitted": true, "workflow": workflowID,
		"forge_response": forgeResult,
		"images": forgeResult["images"],
		"info":   forgeResult["info"],
	})
}

// ═══════════════════════════════════════════════════════
// Storyboard → Video Pipeline
// ═══════════════════════════════════════════════════════

func handleStoryboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	// Accept multipart (multiple images) or JSON (paths)
	var paths []string
	model := "three-act"
	duration := 30
	desc := ""
	singleMode := false

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("parse: %v", err)})
			return
		}
		model = r.FormValue("model")
		if model == "" { model = "three-act" }
		if d, err := strconv.Atoi(r.FormValue("duration")); err == nil && d > 0 {
			duration = d
		}
		desc = r.FormValue("desc")
		singleMode = r.FormValue("single") == "true"

		files := r.MultipartForm.File["images"]
		for _, fh := range files {
			f, err := fh.Open()
			if err != nil { continue }
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" { f.Close(); continue }
			filename := fmt.Sprintf("sb_%d_%s", time.Now().UnixNano(), fh.Filename)
			dst := filepath.Join(uploadDir(), filename)
			out, err := os.Create(dst)
			if err != nil { f.Close(); continue }
			if _, err := io.Copy(out, f); err != nil { out.Close(); f.Close(); continue }
			out.Close()
			f.Close()
			paths = append(paths, dst)
		}
	} else {
		var req map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
		if err := json.Unmarshal(body, &req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
		}
		if p, ok := req["paths"].([]interface{}); ok {
			for _, pp := range p {
				if s, ok := pp.(string); ok {
					if !filepath.IsAbs(s) { s = filepath.Join(baseDir, s) }
					paths = append(paths, s)
				}
			}
		}
		if m, ok := req["model"].(string); ok { model = m }
		if d, ok := req["duration"].(float64); ok && d > 0 { duration = int(d) }
		if d, ok := req["desc"].(string); ok { desc = d }
		singleMode, _ = req["single"].(bool)
	}

	if len(paths) == 0 {
		log.Printf("handleStoryboard: no images provided")
		json.NewEncoder(w).Encode(map[string]string{"error": "at least 1 image required"})
		return
	}

	result := vision.ProcessStoryboard(paths, model, duration, desc, singleMode)
	json.NewEncoder(w).Encode(result)
}

func handleVideoPrompt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}
	var req map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
	}

	imgPath, _ := req["path"].(string)
	desc, _ := req["desc"].(string)
	workflow, _ := req["workflow"].(string)
	if workflow == "" { workflow = "animatediff_txt2vid" }

	if imgPath == "" {
		log.Printf("handleVideoPrompt: path required")
		json.NewEncoder(w).Encode(map[string]string{"error": "path required"})
		return
	}
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(baseDir, imgPath)
	}
	// Path traversal prevention
	cleanPath := filepath.Clean(imgPath)
	cleanBase := filepath.Clean(baseDir)
	if cleanPath != cleanBase && !strings.HasPrefix(cleanPath, cleanBase+string(filepath.Separator)) {
		log.Printf("handleVideoPrompt: path traversal blocked: %s", imgPath)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path"})
		return
	}

	// Single image → video prompt
	result := vision.ProcessStoryboard([]string{imgPath}, "three-act", 10, desc, true)

	// Add recommended workflow
	if vp, ok := result["video_prompt"].(map[string]interface{}); ok {
		vp["recommended_workflow"] = workflow
	}

	json.NewEncoder(w).Encode(result)
}

// ═══════════════════════════════════════════════════════
// Story → PPT (Doc/Image/PPTX/TXT → Presentation framework)
// ═══════════════════════════════════════════════════════

func handleStory2PPT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"error": "use POST"})
		return
	}

	var docPath string
	userPrompt := ""

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("parse: %v", err)})
			return
		}
		userPrompt = r.FormValue("prompt")
		file, header, err := r.FormFile("file")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("file: %v", err)})
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		filename := fmt.Sprintf("s2p_%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join(uploadDir(), filename)
		out, err := os.Create(dst)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("create: %v", err)})
			file.Close()
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("write file: %v", err)})
			return
		}
		docPath = dst
	} else {
		var req map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": "read body: "+err.Error()}); return }
		if err := json.Unmarshal(body, &req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: "+err.Error()}); return
		}
		docPath, _ = req["path"].(string)
		userPrompt, _ = req["prompt"].(string)
		if !filepath.IsAbs(docPath) {
			docPath = filepath.Join(baseDir, docPath)
		}
	}

	if docPath == "" {
		log.Printf("handleStory2PPT: file required")
		json.NewEncoder(w).Encode(map[string]string{"error": "file required"})
		return
	}

	result := vision.Story2PPT(docPath, userPrompt)
	json.NewEncoder(w).Encode(result)
}
