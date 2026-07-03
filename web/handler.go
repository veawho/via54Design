// via54Design — Web UI HTTP Handlers (全功能版)
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: MIT

package web

import (
	"bytes"
	"embed"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	media "github.com/veawho/via54Design/internal/media"
	vt "github.com/veawho/via54Design/internal/template"
	"github.com/veawho/via54Design/internal/vision"
	"github.com/veawho/via54Design/internal/workflow"
)

//go:embed templates/index.html
//go:embed templates/pane_design.html
//go:embed templates/pane_prompt.html
//go:embed templates/pane_present.html
//go:embed templates/pane_video.html
//go:embed templates/pane_forge.html
//go:embed templates/pane_reimagine.html
//go:embed templates/pane_tools.html
//go:embed templates/phases45.html
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

	// ── HTMX endpoints (return HTML fragments, no JS) ──
	mux.HandleFunc("/api/local-fonts", handleLocalFonts)
	mux.HandleFunc("/api/htmx/pattern", handleHTMXPattern)
	mux.HandleFunc("/api/htmx/quality", handleHTMXQuality)
	mux.HandleFunc("/api/htmx/list", handleHTMXList)
	mux.HandleFunc("/api/htmx/ai", handleHTMXAI)
	mux.HandleFunc("/api/htmx/trace", handleHTMXTrace)
	mux.HandleFunc("/api/htmx/status", handleHTMXStatus)
	mux.HandleFunc("/api/htmx/pane", handleHTMXPane)
	mux.HandleFunc("/api/htmx/generate", handleHTMXGenerate)
	mux.HandleFunc("/api/htmx/prompt", handleHTMXPrompt)
	mux.HandleFunc("/api/htmx/narrate", handleHTMXNarrate)
	mux.HandleFunc("/api/htmx/upload", handleHTMXUpload)
	mux.HandleFunc("/api/htmx/regen", handleHTMXRegen)
	mux.HandleFunc("/api/htmx/story2ppt", handleHTMXStory2PPT)
	mux.HandleFunc("/api/htmx/forge-status", handleHTMXForgeStatus)
	mux.HandleFunc("/api/htmx/reimagine", handleHTMXReimagine)
	mux.HandleFunc("/api/htmx/download", handleHTMXDownload)
	mux.HandleFunc("/api/htmx/preview", handleHTMXPreview)
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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
	exe := selfPath()
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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

	exe := selfPath()
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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

	exe := selfPath()
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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

	exe := selfPath()
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
			"error":  fmt.Sprintf("exec error: %v", err),
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	action, _ := req["action"].(string)
	source, _ := req["source"].(string)
	target, _ := req["target"].(string)
	paramsStr, _ := req["params"].(string)
	if action == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "action required"})
		return
	}

	exe := selfPath()
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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

	cli := selfPath()
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
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
			"success":      true,
			"submitted":    false,
			"workflow":     workflowID,
			"message":      "后端（Forge）未连接，提示词已就绪可手动提交",
			"hint":         "启动 Forge 后重试，或复制下方 payload 手动 POST",
			"prompt_info":  info,
			"api_payload":  payload,
			"api_endpoint": "http://localhost:7860/sdapi/v1/txt2img",
			"offline_mode": true,
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
		"images":         forgeResult["images"],
		"info":           forgeResult["info"],
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
		if model == "" {
			model = "three-act"
		}
		if d, err := strconv.Atoi(r.FormValue("duration")); err == nil && d > 0 {
			duration = d
		}
		desc = r.FormValue("desc")
		singleMode = r.FormValue("single") == "true"

		files := r.MultipartForm.File["images"]
		for _, fh := range files {
			f, err := fh.Open()
			if err != nil {
				continue
			}
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
				f.Close()
				continue
			}
			filename := fmt.Sprintf("sb_%d_%s", time.Now().UnixNano(), fh.Filename)
			dst := filepath.Join(uploadDir(), filename)
			out, err := os.Create(dst)
			if err != nil {
				f.Close()
				continue
			}
			if _, err := io.Copy(out, f); err != nil {
				out.Close()
				f.Close()
				continue
			}
			out.Close()
			f.Close()
			paths = append(paths, dst)
		}
	} else {
		var req map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
			return
		}
		if p, ok := req["paths"].([]interface{}); ok {
			for _, pp := range p {
				if s, ok := pp.(string); ok {
					if !filepath.IsAbs(s) {
						s = filepath.Join(baseDir, s)
					}
					paths = append(paths, s)
				}
			}
		}
		if m, ok := req["model"].(string); ok {
			model = m
		}
		if d, ok := req["duration"].(float64); ok && d > 0 {
			duration = int(d)
		}
		if d, ok := req["desc"].(string); ok {
			desc = d
		}
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
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	imgPath, _ := req["path"].(string)
	desc, _ := req["desc"].(string)
	workflow, _ := req["workflow"].(string)
	if workflow == "" {
		workflow = "animatediff_txt2vid"
	}

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
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "read body: " + err.Error()})
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
			return
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

// ═══════════════════════════════════════════════════════
// HTMX Handlers — return HTML fragments, no JS required
// ═══════════════════════════════════════════════════════

// selfPath returns the path to the current executable (cross-platform)
func selfPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "via54"
	}
	return exe
}

func htmxWrite(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func htmxError(w http.ResponseWriter, msg string) {
	htmxWrite(w, `<div class="output-area" style="color:var(--accent2)">❌ `+msg+`</div>`)
}

func handleHTMXStatus(w http.ResponseWriter, r *http.Request) {
	// Engine status
	engineDot, engineText := "dot off", "未连接"
	if resp, err := http.Get("http://localhost:8642/health"); err == nil {
		defer resp.Body.Close()
		engineDot, engineText = "dot ok", "✓ v0.6.0"
		_ = resp.Body.Close()
	}

	// Forge status
	forgeDot, forgeText := "dot off", "未连接"
	if resp, err := http.Get("http://localhost:7860/sdapi/v1/sd-models"); err == nil {
		defer resp.Body.Close()
		forgeDot, forgeText = "dot ok", "✓ 已连接"
	}

	// Templates
	tplDot, tplText := "dot off", "?"
	if _, err := workflow.LoadRegistry(baseDir); err == nil {
		tplDot, tplText = "dot ok", "✓ templates"
	}

	htmxWrite(w, fmt.Sprintf(`
<div class="status-grid">
  <span class="status-chip"><span class="dot %s"></span><span class="label">引擎</span><span class="val">%s</span></span>
  <span class="status-chip"><span class="dot %s"></span><span class="label">模板</span><span class="val">%s</span></span>
  <span class="status-chip"><span class="dot ok"></span><span class="label">二进制</span><span class="val">via54.exe</span></span>
  <span class="status-chip"><span class="dot %s"></span><span class="label">Forge</span><span class="val">%s</span></span>
</div>`, engineDot, engineText, tplDot, tplText, forgeDot, forgeText))
}

func handleHTMXPane(w http.ResponseWriter, r *http.Request) {
	intent := r.URL.Query().Get("intent")
	names := map[string]string{
		"design":    "pane_design.html",
		"prompt":    "pane_prompt.html",
		"present":   "pane_present.html",
		"video":     "pane_video.html",
		"forge":     "pane_forge.html",
		"reimagine": "pane_reimagine.html",
		"tools":     "pane_tools.html",
	}
	file, ok := names[intent]
	if !ok {
		file = "pane_design.html"
	}
	data, err := embeddedFiles.ReadFile("templates/" + file)
	if err != nil {
		htmxError(w, "pane not found: "+file)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func handleHTMXGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		htmxError(w, err.Error())
		return
	}
	exePath := selfPath()
	if _, statErr := os.Stat(exePath); statErr != nil {
		htmxError(w, fmt.Sprintf("selfPath binary not found: %s - %v", exePath, statErr))
		return
	}

	title := r.FormValue("title")
	if title == "" {
		title = "via54Design"
	}
	mode := r.FormValue("mode")
	layout := r.FormValue("layout")
	if layout == "" {
		layout = "hero-split-16-9"
	}
	color := r.FormValue("color")
	if color == "" {
		color = "ink-wash"
	}
	font := r.FormValue("font")
	if font == "" {
		font = "display-sans-bold"
	}
	pres := mode == "presentation"

	// If the chosen font is a local system font rather than preset ID, fallback to display-sans-bold for compiler
	originalFont := font
	isPreset := font == "display-sans-bold" || font == "cormorant-elegant" || font == "mono-terminal" || font == "ming-hei-editorial" || font == "system-utility"
	if !isPreset {
		font = "display-sans-bold"
	}

	exe := selfPath()
	args := []string{"generate", "--layout", layout, "--color", color, "--font", font, "--title", title}
	if pres {
		args = append(args, "--presentation")
	}
	cmd := exec.Command(exe, args...)
	cmd.Dir = baseDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		htmxError(w, fmt.Sprintf("生成失败: %v\n<pre style='font-size:11px'>%s</pre>", err, string(out)))
		return
	}
	// Post-processing to inject system font override if necessary
	if !isPreset && originalFont != "" {
		outPath := filepath.Join(baseDir, "output.html")
		htmlBytes, err := os.ReadFile(outPath)
		if err == nil {
			html := string(htmlBytes)
			// Replace display font, body font, and general styling fonts with local family name
			html = strings.ReplaceAll(html, "'Archivo Black', 'Anton', 'Manrope', sans-serif", fmt.Sprintf("'%s', sans-serif", originalFont))
			html = strings.ReplaceAll(html, "'Inter', -apple-system, 'Helvetica Neue', sans-serif", fmt.Sprintf("'%s', sans-serif", originalFont))
			html = strings.ReplaceAll(html, "'Outfit', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans', 'Microsoft YaHei', 'Meiryo', 'Noto Sans SC', sans-serif", originalFont)
			_ = os.WriteFile(outPath, []byte(html), 0644)
		}
	}

	htmxWrite(w, fmt.Sprintf(`<div class="output-area">✅ 已生成 (%d bytes)<br><br><a href="/api/htmx/download?name=%s" class="btn-small">📥 下载 HTML</a></div>`, len(out), title))
}

func handleHTMXDownload(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "output"
	}
	path := filepath.Join(baseDir, "output.html")
	data, err := os.ReadFile(path)
	if err != nil {
		htmxError(w, "文件未生成")
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.html\"", name))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func handleHTMXPreview(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(baseDir, "output.html")
	data, err := os.ReadFile(path)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<html><body style="background:#050508;color:#737068;font-family:sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;padding:20px;text-align:center"><div><div style="font-size:32px;margin-bottom:12px">🖥️</div><div style="font-size:14px;font-weight:600">via54Design Sandbox Preview</div><div style="font-size:12px;margin-top:6px;opacity:0.8">No preview generated yet. Enter parameters on the left and click Generate HTML.</div></div></body></html>`))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}


func handleHTMXPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		htmxError(w, err.Error())
		return
	}

	scene := r.FormValue("scene")
	platform := r.FormValue("platform")
	if scene == "" {
		htmxError(w, "请输入场景描述")
		return
	}
	if platform == "" {
		platform = "midjourney"
	}

	exe := selfPath()
	cmd := exec.Command(exe, "prompt", "--scene", scene, "--platform", platform, "--format", "markdown")
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		htmxError(w, fmt.Sprintf("生成失败: %v", err))
		return
	}
	escaped := strings.ReplaceAll(string(out), "<", "&lt;")
	escaped = strings.ReplaceAll(escaped, ">", "&gt;")
	htmxWrite(w, `<div class="output-area"><pre style="white-space:pre-wrap">`+escaped+`</pre></div>`)
}

func handleHTMXNarrate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		htmxError(w, err.Error())
		return
	}

	seed := r.FormValue("seed")
	if seed == "" {
		htmxError(w, "请输入故事种子")
		return
	}
	model := r.FormValue("model")
	if model == "" {
		model = "three-act"
	}
	durationStr := r.FormValue("duration")
	duration := 30
	if d, err := strconv.Atoi(durationStr); err == nil && d > 0 {
		duration = d
	}

	exe := selfPath()
	cmd := exec.Command(exe, "narrate", "--seed", seed, "--model", model, "--duration", strconv.Itoa(duration), "--format", "markdown")
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		htmxError(w, fmt.Sprintf("生成失败: %v", err))
		return
	}
	escaped := strings.ReplaceAll(string(out), "<", "&lt;")
	escaped = strings.ReplaceAll(escaped, ">", "&gt;")
	htmxWrite(w, `<div class="output-area"><pre style="white-space:pre-wrap">`+escaped+`</pre></div>`)
}

func handleHTMXUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		htmxError(w, err.Error())
		return
	}

	field := r.FormValue("_field")
	if field == "" {
		field = "image"
	}
	file, header, err := r.FormFile(field)
	if err != nil {
		// Fallback for standard layout upload field names used across templates
		for _, fallbackField := range []string{"file", "screenshot", "image"} {
			if f, h, e := r.FormFile(fallbackField); e == nil {
				file = f
				header = h
				err = nil
				break
			}
		}
	}
	if err != nil {
		htmxError(w, "文件未上传: "+err.Error())
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("img_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir(), filename)
	out, err := os.Create(dst)
	if err != nil {
		htmxError(w, "保存失败")
		return
	}
	defer out.Close()
	io.Copy(out, file)

	isImg := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" || ext == ".gif"
	previewHTML := ""
	if isImg {
		url := "/uploads/" + filename
		previewHTML = fmt.Sprintf(`<img src="%s" style="width:40px;height:40px;object-fit:cover;border-radius:4px">`, url)
	} else {
		previewHTML = `<span style="font-size:24px;margin-right:4px">📄</span>`
	}

	htmxWrite(w, fmt.Sprintf(`<input type="hidden" name="_path" value="%s">
<div style="display:flex;align-items:center;gap:8px;margin-top:6px">
  %s
  <span style="font-size:12px;color:var(--text-secondary)">%s (%d bytes)</span>
</div>`, dst, previewHTML, header.Filename, header.Size))
}

func handleHTMXRegen(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		htmxError(w, err.Error())
		return
	}

	prompt := r.FormValue("prompt")
	workflowID := r.FormValue("workflow")
	if workflowID == "" {
		workflowID = "sdxl_txt2img"
	}
	if prompt == "" {
		htmxError(w, "请先生成提示词")
		return
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"prompt": prompt, "negative_prompt": r.FormValue("negative"),
		"steps": 30, "cfg_scale": 7.5, "width": 1024, "height": 1024,
	})

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post("http://localhost:7860/sdapi/v1/txt2img", "application/json", bytes.NewReader(payload))
	if err != nil {
		htmxWrite(w, fmt.Sprintf(`<div class="output-area">⚠️ Forge 未运行<br><br>提示词就绪 (%d chars)<br><code style="font-size:11px">via54 forge --workflow %s --prompt "..." --send</code></div>`, len(prompt), workflowID))
		return
	}
	defer resp.Body.Close()
	htmxWrite(w, `<div class="output-area" style="color:var(--accent3)">✅ 已提交到 Forge! (工作流: `+workflowID+`)</div>`)
}

func handleHTMXStory2PPT(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		htmxError(w, err.Error())
		return
	}

	path := r.FormValue("_path")
	seed := r.FormValue("seed")

	out := ""
	if path != "" {
		result := vision.Story2PPT(path, seed)
		if slides, ok := result["slides"].([]interface{}); ok {
			for i, s := range slides {
				if m, ok := s.(map[string]interface{}); ok {
					title, _ := m["title"].(string)
					out += fmt.Sprintf("<li>%d. %s</li>\n", i+1, title)
				}
			}
		}
		if out == "" {
			if errStr, ok := result["error"].(string); ok {
				htmxError(w, errStr)
				return
			}
			out = "<li>分析完成</li>"
		}
		htmxWrite(w, `<div class="output-area"><ol>`+out+`</ol></div>`)
	} else if seed != "" {
		exe := selfPath()
		cmd := exec.Command(exe, "narrate", "--seed", seed, "--model", "three-act", "--duration", "30", "--format", "markdown")
		cmd.Dir = baseDir
		result, err := cmd.Output()
		if err != nil {
			htmxError(w, err.Error())
			return
		}
		escaped := strings.ReplaceAll(string(result), "<", "&lt;")
		htmxWrite(w, `<div class="output-area"><pre style="white-space:pre-wrap">`+escaped+`</pre></div>`)
	} else {
		htmxError(w, "请上传文件或输入故事种子")
	}
}

func handleHTMXStoryboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		htmxError(w, err.Error())
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		htmxError(w, "请上传至少一张故事板图片")
		return
	}

	model := r.FormValue("model")
	if model == "" {
		model = "three-act"
	}
	duration := 30
	if d, err := strconv.Atoi(r.FormValue("duration")); err == nil && d > 0 {
		duration = d
	}

	var paths []string
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		filename := fmt.Sprintf("sb_%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join(uploadDir(), filename)
		out, _ := os.Create(dst)
		io.Copy(out, f)
		out.Close()
		f.Close()
		paths = append(paths, dst)
	}

	result := vision.ProcessStoryboard(paths, model, duration, "", len(files) == 1)

	var html string
	if scaffold, ok := result["narrative_scaffold"].(map[string]interface{}); ok {
		name, _ := scaffold["model_name"].(string)
		td, _ := scaffold["total_duration"].(float64)
		html = fmt.Sprintf("<div class='output-area'><strong>📖 %s | %ds</strong><br>", name, int(td))
		if beats, ok := scaffold["beats"].([]interface{}); ok {
			for _, b := range beats {
				if m, ok := b.(map[string]interface{}); ok {
					n, _ := m["name"].(string)
					st, _ := m["start_time"].(float64)
					d, _ := m["duration"].(float64)
					mo, _ := m["mood"].(string)
					html += fmt.Sprintf("  <br>%s (%ds-%ds) [%s]", n, int(st), int(st+d), mo)
				}
			}
		}
		html += "</div>"
	} else {
		html = fmt.Sprintf("<div class='output-area'>%s</div>", result)
	}
	htmxWrite(w, html)
}

func handleHTMXForgeStatus(w http.ResponseWriter, r *http.Request) {
	dot, text := "dot off", "❌ 未运行"
	if resp, err := http.Get("http://localhost:7860/sdapi/v1/sd-models"); err == nil {
		defer resp.Body.Close()
		dot, text = "dot ok", "✅ 已连接"
	}
	htmxWrite(w, fmt.Sprintf(`<span class="status-chip"><span class="dot %s"></span>%s</span>`, dot, text))
}

// handleHTMXReimagine — HTMX 调用 via54 reimagine 子命令
// 流程: HTMX 上传文件 → 拿路径 → exec via54 reimagine → 输出 HTML
func handleHTMXReimagine(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		htmxError(w, err.Error())
		return
	}

	dst := r.FormValue("_path")
	if dst == "" {
		file, header, err := r.FormFile("screenshot")
		if err == nil {
			defer file.Close()
			ext := strings.ToLower(filepath.Ext(header.Filename))
			filename := fmt.Sprintf("shot_%d%s", time.Now().UnixNano(), ext)
			dst = filepath.Join(uploadDir(), filename)
			out, err := os.Create(dst)
			if err == nil {
				io.Copy(out, file)
				out.Close()
			} else {
				htmxError(w, "保存截图失败: "+err.Error())
				return
			}
		}
	}
	if dst == "" {
		htmxError(w, "截图未上传，或上传数据损坏")
		return
	}

	provider := r.FormValue("provider")
	if provider == "" {
		provider = "openai"
	}
	layoutHint := r.FormValue("layout_hint")
	colorHint := r.FormValue("color_hint")
	fontHint := r.FormValue("font_hint")
	
	// Custom config parameters sent from UI
	customColorVars := r.FormValue("custom_color_vars")
	customFontImport := r.FormValue("custom_font_import")
	customFontFamily := r.FormValue("custom_font_family")

	exePath := selfPath()
	if _, statErr := os.Stat(exePath); statErr != nil {
		htmxError(w, fmt.Sprintf("via54.exe not found: %s", exePath))
		return
	}

	// exec via54 reimagine
	outFilename := fmt.Sprintf("reimagined_%d.html", time.Now().UnixNano())
	outPath := filepath.Join(uploadDir(), outFilename)
	args := []string{"reimagine", "--screenshot", dst, "--provider", provider, "--output", outPath}
	if layoutHint != "" {
		args = append(args, "--layout-hint", layoutHint)
	}
	if colorHint != "" {
		args = append(args, "--color-hint", colorHint)
	}
	if fontHint != "" {
		args = append(args, "--font-hint", fontHint)
	}
	cmd := exec.Command(exePath, args...)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		htmxError(w, fmt.Sprintf("reimagine 失败: %v<br><pre style='font-size:11px;color:var(--text-dim)'>%s</pre>", err, string(cmdOutput)))
		return
	}

	// 读生成的 HTML 并注入自定义样式参数
	htmlBytes, err := os.ReadFile(outPath)
	if err != nil {
		htmxError(w, "读取生成结果失败: "+err.Error())
		return
	}
	
	// Override styles if custom configuration variables are present
	if customColorVars != "" || customFontImport != "" || customFontFamily != "" {
		htmlBytes = []byte(injectCustomStyles(string(htmlBytes), customColorVars, customFontImport, customFontFamily))
		_ = os.WriteFile(outPath, htmlBytes, 0644)
	}

	htmlPreview := string(htmlBytes)
	if len(htmlPreview) > 8000 {
		htmlPreview = htmlPreview[:8000] + "\n\n... (truncated)"
	}

	htmxWrite(w, fmt.Sprintf(`
<div class="output-area" style="margin-bottom:8px;color:var(--accent3)">✅ 复刻完成 (%d bytes HTML) — provider: %s</div>
<details style="margin-top:6px">
  <summary style="cursor:pointer;font-size:12px;color:var(--text-secondary)">📄 查看生成 HTML (前 8000 字符)</summary>
  <pre style="background:var(--bg-inset);padding:10px;border-radius:var(--radius-sm);font-size:11px;overflow:auto;max-height:400px;margin-top:6px;color:var(--text-secondary)"><code>%s</code></pre>
</details>
<div style="margin-top:6px">
  <a href="/uploads/%s" target="_blank" class="btn-primary" style="display:inline-block;text-decoration:none;font-size:12px">🔗 在新窗口打开</a>
  <span style="font-size:11px;color:var(--text-dim);margin-left:8px">/uploads/%s</span>
</div>
`, len(htmlBytes), provider, htmlPreview, outFilename, outFilename))
}


// ═══════════════════════════════════════════════════════
// Local Font Scanning & Custom Styling Injections
// ═══════════════════════════════════════════════════════

var (
	localFontsCache []string
	localFontsMutex sync.Mutex
)

func handleLocalFonts(w http.ResponseWriter, r *http.Request) {
	fonts := getLocalFontsCached()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fonts)
}

func getLocalFontsCached() []string {
	localFontsMutex.Lock()
	defer localFontsMutex.Unlock()

	if len(localFontsCache) > 0 {
		return localFontsCache
	}

	dirs := getSystemFontDirs()
	fontMap := make(map[string]bool)

	for _, dir := range dirs {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".ttf" || ext == ".otf" {
				family, err := readFontFamily(path)
				if err == nil && family != "" {
					family = strings.TrimSpace(family)
					if !strings.HasPrefix(family, ".") && len(family) > 1 {
						fontMap[family] = true
					}
				} else {
					base := filepath.Base(path)
					name := strings.TrimSuffix(base, filepath.Ext(base))
					name = strings.ReplaceAll(name, "-", " ")
					name = strings.ReplaceAll(name, "_", " ")
					name = strings.Title(name)
					fontMap[name] = true
				}
			}
			if len(fontMap) >= 150 {
				return filepath.SkipDir
			}
			return nil
		})
		if len(fontMap) >= 150 {
			break
		}
	}

	var list []string
	for k := range fontMap {
		list = append(list, k)
	}
	sort.Strings(list)
	localFontsCache = list
	return list
}

func getSystemFontDirs() []string {
	var dirs []string
	switch runtime.GOOS {
	case "windows":
		windir := os.Getenv("WINDIR")
		if windir == "" {
			windir = "C:\\Windows"
		}
		dirs = append(dirs, filepath.Join(windir, "Fonts"))
		localappdata := os.Getenv("LOCALAPPDATA")
		if localappdata != "" {
			dirs = append(dirs, filepath.Join(localappdata, "Microsoft\\Windows\\Fonts"))
		}
	case "darwin":
		dirs = append(dirs, "/Library/Fonts", "/System/Library/Fonts", filepath.Join(os.Getenv("HOME"), "Library/Fonts"))
	default:
		dirs = append(dirs, "/usr/share/fonts", "/usr/local/share/fonts", filepath.Join(os.Getenv("HOME"), ".fonts"))
	}
	return dirs
}

func readFontFamily(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	header := make([]byte, 12)
	if _, err := io.ReadFull(f, header); err != nil {
		return "", err
	}

	numTables := binary.BigEndian.Uint16(header[4:6])

	dirs := make([]byte, numTables*16)
	if _, err := io.ReadFull(f, dirs); err != nil {
		return "", err
	}

	var nameOffset uint32
	var nameLen uint32
	for i := 0; i < int(numTables); i++ {
		offset := i * 16
		tag := string(dirs[offset : offset+4])
		if tag == "name" {
			nameOffset = binary.BigEndian.Uint32(dirs[offset+8 : offset+12])
			nameLen = binary.BigEndian.Uint32(dirs[offset+12 : offset+16])
			break
		}
	}

	if nameOffset == 0 {
		return "", io.EOF
	}

	if _, err := f.Seek(int64(nameOffset), io.SeekStart); err != nil {
		return "", err
	}
	nameTable := make([]byte, nameLen)
	if _, err := io.ReadFull(f, nameTable); err != nil {
		return "", err
	}

	if len(nameTable) < 6 {
		return "", io.EOF
	}

	count := binary.BigEndian.Uint16(nameTable[2:4])
	stringOffset := binary.BigEndian.Uint16(nameTable[4:6])

	for i := 0; i < int(count); i++ {
		offset := 6 + i*12
		if offset+12 > len(nameTable) {
			break
		}
		nameID := binary.BigEndian.Uint16(nameTable[offset+6 : offset+8])
		if nameID == 1 {
			length := binary.BigEndian.Uint16(nameTable[offset+8 : offset+10])
			strOff := binary.BigEndian.Uint16(nameTable[offset+10 : offset+12])
			
			start := int(stringOffset) + int(strOff)
			end := start + int(length)
			if end <= len(nameTable) {
				raw := nameTable[start:end]
				platformID := binary.BigEndian.Uint16(nameTable[offset : offset+2])
				if platformID == 0 || platformID == 3 {
					runes := make([]rune, len(raw)/2)
					for j := 0; j < len(runes); j++ {
						runes[j] = rune(binary.BigEndian.Uint16(raw[j*2 : j*2+2]))
					}
					name := string(runes)
					if name != "" {
						return name, nil
					}
				} else {
					name := string(raw)
					if name != "" {
						return name, nil
					}
				}
			}
		}
	}

	return "", io.EOF
}

func injectCustomStyles(html string, colorVars string, fontImport string, fontFamily string) string {
	if fontImport != "" {
		linkTag := fmt.Sprintf(`<link href="%s" rel="stylesheet">`, fontImport)
		if strings.Contains(html, "</head>") {
			html = strings.Replace(html, "</head>", linkTag+"\n</head>", 1)
		} else {
			html = linkTag + "\n" + html
		}
	}

	var cssInject strings.Builder
	cssInject.WriteString("\n  :root {\n")
	if colorVars != "" {
		var vars map[string]string
		if err := json.Unmarshal([]byte(colorVars), &vars); err == nil {
			for k, v := range vars {
				cssInject.WriteString(fmt.Sprintf("    %s: %s !important;\n", k, v))
			}
		}
	}
	if fontFamily != "" {
		cssInject.WriteString(fmt.Sprintf("    --font-sans: %s !important;\n", fontFamily))
		cssInject.WriteString(fmt.Sprintf("    font-family: %s !important;\n", fontFamily))
	}
	cssInject.WriteString("  }\n")

	if strings.Contains(html, "</style>") {
		html = strings.Replace(html, "</style>", cssInject.String()+"</style>", 1)
	} else if strings.Contains(html, "</head>") {
		html = strings.Replace(html, "</head>", fmt.Sprintf("<style>%s</style>\n</head>", cssInject.String()), 1)
	}
	return html
}


// ═══════════════════════════════════════════════════════
// Subcommands Integrated into Web UI (Other Tools)
// ═══════════════════════════════════════════════════════

func handleHTMXPattern(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	style := r.FormValue("pattern")
	var css, previewStyle string
	
	switch style {
	case "gold-grid":
		css = `background-color: #0d0d0d;
background-image: linear-gradient(rgba(212, 175, 55, 0.1) 1px, transparent 1px),
                  linear-gradient(90deg, rgba(212, 175, 55, 0.1) 1px, transparent 1px);
background-size: 20px 20px;`
		previewStyle = "background-color:#0d0d0d; background-image:linear-gradient(rgba(212,175,55,0.1) 1px,transparent 1px),linear-gradient(90deg,rgba(212,175,55,0.1) 1px,transparent 1px); background-size:20px 20px; height:100px; border-radius:var(--radius-sm); border:1px solid var(--border)"
	case "bauhaus-dot":
		css = `background-color: #f5f5f7;
background-image: radial-gradient(rgba(0, 0, 0, 0.15) 1.5px, transparent 1.5px);
background-size: 24px 24px;`
		previewStyle = "background-color:#f5f5f7; background-image:radial-gradient(rgba(0,0,0,0.15) 1.5px,transparent 1.5px); background-size:24px 24px; height:100px; border-radius:var(--radius-sm); border:1px solid rgba(0,0,0,0.1)"
	case "mesh-gradient":
		css = `background: radial-gradient(circle at 0% 0%, #1a1a2e, transparent 50%),
            radial-gradient(circle at 100% 0%, #0d0d0d, transparent 50%),
            radial-gradient(circle at 100% 100%, #1c1505, transparent 50%),
            radial-gradient(circle at 0% 100%, #111111, transparent 50%);`
		previewStyle = "background:radial-gradient(circle at 0% 0%,#1a1a2e,transparent 50%),radial-gradient(circle at 100% 0%,#0d0d0d,transparent 50%),radial-gradient(circle at 100% 100%,#1c1505,transparent 50%),radial-gradient(circle at 0% 100%,#111111,transparent 50%); height:100px; border-radius:var(--radius-sm); border:1px solid var(--border)"
	case "luxury-stripe":
		css = `background: repeating-linear-gradient(45deg, #121212, #121212 10px, #1a1505 10px, #1a1505 20px);`
		previewStyle = "background:repeating-linear-gradient(45deg,#121212,#121212 10px,#1a1505 10px,#1a1505 20px); height:100px; border-radius:var(--radius-sm); border:1px solid var(--border)"
	}

	html := fmt.Sprintf(`
<div class="output-area">
  <strong>🎨 生成 CSS 背景成功!</strong>
  <div style="margin: 10px 0; %s"></div>
  <pre style="background:var(--bg-inset);padding:10px;border-radius:var(--radius-sm);font-size:11px;overflow:auto;margin-top:6px;color:var(--text-secondary)"><code>%s</code></pre>
</div>`, previewStyle, css)
	htmxWrite(w, html)
}

func handleHTMXQuality(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	code := r.FormValue("html_code")
	if code == "" {
		// Read generated output.html as default fallback
		outPath := filepath.Join(baseDir, "output.html")
		if data, err := os.ReadFile(outPath); err == nil {
			code = string(data)
		}
	}
	if code == "" {
		htmxError(w, "未检测到已生成的 output.html, 且输入框为空。请先生成设计，或粘贴代码后再审计")
		return
	}

	// Basic quality audit logic
	var audits []string
	addAudit := func(ok bool, title string) {
		symbol := `<span style="color:var(--accent3);margin-right:6px">✓</span>`
		if !ok {
			symbol = `<span style="color:var(--accent2);margin-right:6px">✗</span>`
		}
		audits = append(audits, fmt.Sprintf("<div style='margin-bottom:4px'>%s %s</div>", symbol, title))
	}

	hasViewport := strings.Contains(code, "width=device-width")
	addAudit(hasViewport, "响应式元标签 Viewport Meta 存在")

	hasStyle := strings.Contains(code, "<style>")
	addAudit(hasStyle, "样式内聚层 Style 声明存在")

	hasVariables := strings.Contains(code, "--bg") || strings.Contains(code, "--accent")
	addAudit(hasVariables, "配色语义变量 CSS Variables 符合 via54Design 规范")

	hasAlt := !strings.Contains(code, "src=") || strings.Contains(code, "alt=")
	addAudit(hasAlt, "页面图片 Alt 无障碍描述覆盖")

	hasCJKFont := strings.Contains(code, "PingFang") || strings.Contains(code, "Microsoft YaHei") || strings.Contains(code, "Meiryo")
	addAudit(hasCJKFont, "中日韩系统后备字体 CJK Fonts 补全")

	html := fmt.Sprintf(`
<div class="output-area">
  <strong>🔍 设计规范审计报告 (via54 quality)</strong>
  <div style="margin-top:8px;font-size:12px;color:var(--text-secondary)">
    %s
  </div>
</div>`, strings.Join(audits, ""))
	htmxWrite(w, html)
}

func handleHTMXList(w http.ResponseWriter, r *http.Request) {
	reg, err := vt.NewRegistry(baseDir)
	if err != nil {
		htmxError(w, "加载预置名录失败: "+err.Error())
		return
	}

	var rows []string
	rows = append(rows, "<tr><th>类型</th><th>标识 ID</th><th>描述</th></tr>")
	for _, l := range reg.Data.Layouts {
		rows = append(rows, fmt.Sprintf("<tr><td>布局 (Layout)</td><td><code>%s</code></td><td>%s</td></tr>", l.ID, l.File))
	}
	for _, c := range reg.Data.ColorSchemes {
		rows = append(rows, fmt.Sprintf("<tr><td>配色 (Color)</td><td><code>%s</code></td><td>%s</td></tr>", c.ID, c.File))
	}
	for _, t := range reg.Data.Typography {
		rows = append(rows, fmt.Sprintf("<tr><td>字体 (Font)</td><td><code>%s</code></td><td>%s</td></tr>", t.ID, t.File))
	}

	html := fmt.Sprintf(`
<style>
.tbl-list { width: 100%%; border-collapse: collapse; font-size:11px; margin-top:6px; color:var(--text-secondary) }
.tbl-list th, .tbl-list td { border: 1px solid var(--border); padding: 6px; text-align: left; }
.tbl-list th { background: var(--bg-inset); color: var(--text-primary) }
</style>
<table class="tbl-list">%s</table>`, strings.Join(rows, ""))
	htmxWrite(w, html)
}

func handleHTMXAI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}
	prompt := r.FormValue("prompt")
	if prompt == "" {
		htmxError(w, "问题不能为空")
		return
	}

	// General AI Assistant response fallback
	response := fmt.Sprintf("🤖 <strong>AI 智能设计顾问:</strong><br><br>关于您的提问 <em>\"%s\"</em>，为您推荐以下 via54Design 最佳设计路径：<br><br>• <strong>包豪斯纯三原色配比</strong>: 可选用 <code>bauhaus-primary</code> 配色，配合 <code>display-sans-bold</code> 展示字体，形成极强对比。<br>• <strong>日式水墨枯山水意境</strong>: 选用 <code>ink-wash</code> 静谧浅灰底色与留白。<br>• <strong>自适应响应式卡片</strong>: 使用 <code>bento-grid-2x2</code> 拼图块自适应缩放。<br><br>您可以在 <strong>Design Sandbox (设计沙盒)</strong> 面板中实时勾选保存并预览这些选项！", prompt)
	
	htmxWrite(w, fmt.Sprintf(`<div class="output-area" style="font-size:12px;line-height:1.6;color:var(--text-secondary)">%s</div>`, response))
}


func handleHTMXTrace(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmxError(w, "use POST")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		htmxError(w, "获取上传图片失败: "+err.Error())
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		htmxError(w, "只支持 PNG 或 JPG 格式图片")
		return
	}

	// Save upload image
	uploadFilename := fmt.Sprintf("trace_in_%d%s", time.Now().UnixNano(), ext)
	uploadPath := filepath.Join(uploadDir(), uploadFilename)
	out, err := os.Create(uploadPath)
	if err != nil {
		htmxError(w, "创建本地图片失败: "+err.Error())
		return
	}
	_, _ = io.Copy(out, file)
	out.Close()

	// Vectorize using TraceImage
	svgPath, err := media.TraceImage(uploadPath, nil)
	if err != nil {
		htmxError(w, "矢量化处理失败: "+err.Error())
		return
	}

	// Read generated SVG content
	svgBytes, err := os.ReadFile(svgPath)
	if err != nil {
		htmxError(w, "读取矢量化结果失败: "+err.Error())
		return
	}
	svgContent := string(svgBytes)
	svgFilename := filepath.Base(svgPath)

	html := fmt.Sprintf(`
<div class="output-area">
  <strong>✅ 矢量化转换完成!</strong>
  <div style="margin: 12px 0; background: var(--bg-inset); padding: 10px; border-radius: var(--radius-sm); border: 1px solid var(--border); display: flex; justify-content: center; align-items: center; max-height: 250px; overflow: hidden;">
    %s
  </div>
  <div style="margin-bottom:8px">
    <a href="/uploads/%s" target="_blank" class="btn-primary" style="display:inline-block;text-decoration:none;font-size:12px">📥 下载 SVG 文件</a>
  </div>
  <details>
    <summary style="cursor:pointer;font-size:12px;color:var(--text-secondary)">📄 查看 SVG 源码</summary>
    <pre style="background:var(--bg-inset);padding:10px;border-radius:var(--radius-sm);font-size:11px;overflow:auto;max-height:200px;margin-top:6px;color:var(--text-secondary)"><code>%s</code></pre>
  </details>
</div>`, svgContent, svgFilename, strings.ReplaceAll(strings.ReplaceAll(svgContent, "<", "&lt;"), ">", "&gt;"))

	htmxWrite(w, html)
}
