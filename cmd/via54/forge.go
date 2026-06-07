// via54Design — Forge Classic/A1111 命令
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/veawho/via54Design/internal/workflow"
)

func cmdForge() {
	fs := flag.NewFlagSet("forge", flag.ExitOnError)
	workflowID := fs.String("workflow", "", "Workflow template ID (e.g. sdxl_txt2img)")
	prompt := fs.String("prompt", "", "Generation prompt")
	negative := fs.String("negative", "", "Negative prompt")
	output := fs.String("output", "", "Output JSON path (default: stdout)")
	list := fs.Bool("list", false, "List available workflow templates")
	steps := fs.Int("steps", 0, "Override steps")
	cfg := fs.Float64("cfg", 0, "Override CFG scale")
	seed := fs.Int("seed", -1, "Override seed (-1 = random)")
	sampler := fs.String("sampler", "", "Override sampler name (e.g. Euler, DPM++ 2M Karras)")
	width := fs.Int("width", 0, "Override width")
	height := fs.Int("height", 0, "Override height")
	send := fs.Bool("send", false, "Send to running Forge/A1111 at localhost:7860")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	if *list {
		ids, err := workflow.ListWorkflowTemplates(bd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "列出模板失败: %v\n", err)
			os.Exit(1)
		}
		reg, _ := workflow.LoadRegistry(bd)
		fmt.Println("📋 Forge Classic / A1111 可用模板:")
		fmt.Println()
		if reg != nil {
			for _, w := range reg.Workflows {
				fmt.Printf("  📄 %-20s %s\n", w.ID, w.Name)
				fmt.Printf("     %s\n", w.Description)
				fmt.Printf("     模型: %s | %s\n", w.Model, w.Params)
				fmt.Println()
			}
		} else {
			for _, id := range ids {
				tmpl, err := workflow.LoadWorkflowTemplate(id, bd)
				if err != nil { continue }
				fmt.Printf("  📄 %-20s %s\n", tmpl.ID, tmpl.Name)
				fmt.Printf("     %s\n", tmpl.Description)
				fmt.Printf("     模型: %s\n", tmpl.Model)
				fmt.Println()
			}
		}
		fmt.Println("用法: via54 forge --workflow <id> --prompt \"...\"")
		fmt.Println("       via54 forge --workflow sdxl_txt2img --prompt \"a cat\" --send")
		return
	}

	if *workflowID == "" || *prompt == "" {
		fmt.Fprintln(os.Stderr, "请指定 --workflow <模板ID> 和 --prompt \"...\"")
		os.Exit(1)
	}

	tmpl, err := workflow.LoadWorkflowTemplate(*workflowID, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载模板失败: %v\n", err)
		os.Exit(1)
	}

	// Build A1111/Forge API payload
	payload := map[string]interface{}{
		"prompt":              *prompt,
		"negative_prompt":     *negative,
		"seed":                *seed,
		"steps":               *steps,
		"cfg_scale":           *cfg,
		"width":               *width,
		"height":              *height,
		"sampler_name":        *sampler,
		"save_images":         true,
		"send_images":         true,
		"batch_size":          1,
		"n_iter":              1,
	}

	// Apply template defaults if not overridden
	if *steps <= 0 {
		if s, ok := tmpl.Params["steps"]; ok {
			if si, ok := s.(int); ok { payload["steps"] = si }
		} else { payload["steps"] = 30 }
	}
	if *cfg <= 0 {
		if c, ok := tmpl.Params["cfg"]; ok {
			if cf, ok := c.(float64); ok { payload["cfg_scale"] = cf }
		} else { payload["cfg_scale"] = 7.5 }
	}
	if *width <= 0 {
		if w, ok := tmpl.Params["width"]; ok {
			if wi, ok := w.(int); ok { payload["width"] = wi }
		} else { payload["width"] = 1024 }
	}
	if *height <= 0 {
		if h, ok := tmpl.Params["height"]; ok {
			if hi, ok := h.(int); ok { payload["height"] = hi }
		} else { payload["height"] = 1024 }
	}
	if *sampler == "" {
		payload["sampler_name"] = "Euler"
	}
	if *seed < 0 {
		payload["seed"] = -1
	}

	// Clean zero values
	for k, v := range payload {
		switch val := v.(type) {
		case int:
			if val == 0 { delete(payload, k) }
		case float64:
			if val == 0 { delete(payload, k) }
		}
	}

	// Build result info
	result := map[string]interface{}{
		"template":    tmpl.ID,
		"name":        tmpl.Name,
		"format":      "forge_a1111",
		"model":       tmpl.Model,
		"api_payload": payload,
		"api_endpoint": "http://localhost:7860/sdapi/v1/txt2img",
		"usage": []string{
			"Option A: Copy JSON and POST manually",
			"  curl -X POST http://localhost:7860/sdapi/v1/txt2img -H 'Content-Type: application/json' -d @payload.json",
			"",
			"Option B: Use --send flag to auto-submit",
			"  via54 forge --workflow sdxl_txt2img --prompt \"...\" --send",
		},
	}

	jsonData, _ := json.MarshalIndent(result, "", "  ")

	if *output != "" {
		os.WriteFile(*output, jsonData, 0644)
		fmt.Printf("✅ Forge/A1111 payload: %s\n", *output)
		fmt.Printf("   模板: %s (%s)\n", tmpl.ID, tmpl.Name)
		fmt.Printf("   模型: %s | %d steps, cfg %.1f\n", tmpl.Model, payload["steps"], payload["cfg_scale"])
		if *send {
			fmt.Printf("   📤 发送到 localhost:7860...\n")
			fmt.Printf("   curl -X POST http://localhost:7860/sdapi/v1/txt2img ...\n")
		} else {
			fmt.Printf("   使用 --send 自动提交, 或手动粘贴到 Forge/A1111 WebUI\n")
		}
	} else {
		fmt.Print(string(jsonData))
	}
}
