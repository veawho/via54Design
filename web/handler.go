// via54Design — Web UI HTTP Handlers
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
	"path/filepath"
	"strings"

	"github.com/veawho/via54Design/internal/workflow"
)

//go:embed templates/index.html
var embeddedFiles embed.FS

// baseDir caches the project base directory
var baseDir string

func init() {
	// Try to find the project root
	candidates := []string{".", "..", "../.."}
	for _, c := range candidates {
		path := filepath.Join(c, "templates", "workflows")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)
			baseDir = abs
			return
		}
	}
	// Fallback: use CWD
	baseDir, _ = os.Getwd()
}

// Handler returns an HTTP handler for the web UI.
func Handler(bd string) http.Handler {
	if bd != "" {
		baseDir = bd
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/templates", handleTemplates)
	mux.HandleFunc("/api/build", handleBuild)
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
		"status":  "ok",
		"version": "v0.5.0",
	})
}

func handleTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reg, err := workflow.LoadRegistry(baseDir)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	workflowID, _ := req["workflow_id"].(string)
	prompt, _ := req["prompt"].(string)
	negative, _ := req["negative"].(string)
	if workflowID == "" || prompt == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "workflow_id and prompt required"})
		return
	}

	// Build overrides
	overrides := make(map[string]interface{})
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

	tmpl, err := workflow.LoadWorkflowTemplate(workflowID, baseDir)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Parse optional keyframes
	kfInput, _ := req["keyframes"].(string)
	var kfs []workflow.Keyframe
	if kfInput != "" {
		for _, kf := range strings.Split(kfInput, ",") {
			kf = strings.TrimSpace(kf)
			if kf == "" { continue }
			parts := strings.SplitN(kf, ":", 2)
			if len(parts) == 2 {
				frame := 0
				fmt.Sscanf(parts[0], "%d", &frame)
				kfs = append(kfs, workflow.Keyframe{Frame: frame, Prompt: strings.TrimSpace(parts[1])})
			}
		}
	}

	// Determine output format (comfyui or forge)
	outputFormat, _ := req["format"].(string)
	if outputFormat == "" {
		outputFormat = "comfyui"
	}

	var resultJSON []byte
	nodeCount := 0
	injected := 0
	kfCount := 0

	if outputFormat == "forge" {
		// Build Forge/A1111 API payload
		forgePayload := map[string]interface{}{
			"prompt":          prompt,
			"negative_prompt": negative,
			"steps":           30,
			"cfg_scale":       7.5,
			"width":           1024,
			"height":          1024,
			"sampler_name":    "Euler",
			"save_images":     true,
		}
		if v, ok := overrides["steps"].(int); ok { forgePayload["steps"] = v }
		if v, ok := overrides["cfg"].(float64); ok { forgePayload["cfg_scale"] = v }
		if v, ok := overrides["width"].(int); ok { forgePayload["width"] = v }
		if v, ok := overrides["height"].(int); ok { forgePayload["height"] = v }
		if v, ok := overrides["seed"].(int); ok { forgePayload["seed"] = v }

		resultMap := map[string]interface{}{
			"format":       "forge_a1111",
			"template":     workflowID,
			"api_payload":  forgePayload,
			"api_endpoint": "http://localhost:7860/sdapi/v1/txt2img",
		}
		resultJSON, _ = json.MarshalIndent(resultMap, "", "  ")
		nodeCount = 1
		injected = 1
		kfCount = 0
	} else {
		// Build ComfyUI workflow JSON (default)
		result, err := workflow.BuildWorkflow(tmpl, prompt, negative, overrides, kfs)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		resultJSON = result.JSON
		nodeCount = len(strings.Split(string(result.JSON), "\n")) / 3
		injected = result.Injected
		kfCount = result.Keyframes
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"json":      string(resultJSON),
		"nodes":     nodeCount,
		"injected":  injected,
		"keyframes": kfCount,
		"template":  workflowID,
		"format":    outputFormat,
	})
}
