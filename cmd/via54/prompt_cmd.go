// via54Design — 图片提示词 (Prompt) 命令 v4
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"os"

	"path/filepath"
	"github.com/veawho/via54Design/internal/prompt"
	"github.com/veawho/via54Design/internal/workflow"
)

func cmdPrompt() {
	if len(os.Args) < 3 {
		promptHelp()
		return
	}

	switch os.Args[2] {
	case "edit":
		cmdPromptEdit()
	case "ref":
		cmdPromptRef()
	case "list":
		listPromptPlatforms()
	case "gallery":
		cmdPromptGallery()
	case "comfyui":
		cmdComfyUI()
	case "assess":
		cmdPromptAssess()
	case "version":
		cmdPromptVersion()
	case "send":
		cmdPromptSend()
	default:
		cmdPromptGenerate()
	}
}

func promptHelp() {
	fmt.Println("用法: via54 prompt [命令] [选项]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  <scene>              生成提示词 (默认)")
	fmt.Println("    --scene \"...\"  --platform midjourney/kling/jimeng/gemini")
	fmt.Println("    [--ref ref.jpg] [--output file.md]")
	fmt.Println("  edit                 交互式编辑字段")
	fmt.Println("    --field subject --value \"新主体\" --prompt prompt.json")
	fmt.Println("  list                 列出所有平台")
	fmt.Println("  ref                  参考图分析")
	fmt.Println("    --image ref.jpg --output prompt.json")
	fmt.Println("  gallery              提示词模板市场")
	fmt.Println("  comfyui              ComfyUI 工作流执行")
	fmt.Println("    --workflow workflow.json --prompt \"...\"")
	fmt.Println("  assess              生图质量评估")
	fmt.Println("    --image output.png [--prompt \"used prompt\"]")
	fmt.Println("  version             版本管理 (save/list/diff)")
	fmt.Println("    --save --prompt prompt.json          保存版本")
	fmt.Println("    --list .                             列出版本")
	fmt.Println("    --diff v1 v2                         比较版本")
	fmt.Println("  send                BrowserWing 生图指令")
	fmt.Println("    --prompt prompt.json --platform midjourney")
}

func cmdPromptGenerate() {
	fs := flag.NewFlagSet("prompt", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述")
	platform := fs.String("platform", "midjourney", "平台")
	output := fs.String("output", "", "输出文件")
	ref := fs.String("ref", "", "参考图片路径")
	format := fs.String("format", "markdown", "输出格式: markdown/json")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"场景描述\"")
		os.Exit(1)
	}

	// 验证平台
	if !isValidPlatform(*platform) {
		fmt.Fprintf(os.Stderr, "❌ 未知平台: %q\n", *platform)
		fmt.Fprintf(os.Stderr, "可用平台: %s\n", strings.Join(listPlatforms(), ", "))
		fmt.Fprintln(os.Stderr, "运行 `via54 prompt list` 查看完整列表")
		os.Exit(1)
	}

	s, err := prompt.GeneratePrompt(*scene, *platform, *ref, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}

	if *format == "json" {
		js, _ := s.ToJSON()
		if *output != "" {
			os.WriteFile(*output, []byte(js), 0644)
			fmt.Printf("✅ JSON: %s\n", *output)
		} else {
			fmt.Print(js)
		}
		return
	}

	md, _ := s.RenderMarkdown()
	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 提示词: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

func cmdPromptEdit() {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	field := fs.String("field", "", "字段名")
	value := fs.String("value", "", "新值")
	weight := fs.Float64("weight", 0, "权重 (0.5-2.0)")
	promptFile := fs.String("prompt", "", "prompt JSON 文件")
	platform := fs.String("platform", "midjourney", "平台")
	output := fs.String("output", "", "输出文件")
	fs.Parse(os.Args[3:])

	if *promptFile == "" || *field == "" {
		fmt.Fprintln(os.Stderr, "请指定 --prompt file.json --field name --value \"新值\"")
		os.Exit(1)
	}

	data, err := os.ReadFile(*promptFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取失败: %v\n", err)
		os.Exit(1)
	}

	var s prompt.PromptScaffold
	json.Unmarshal(data, &s)

	if *value != "" {
		s.UpdateField(*field, *value)
	}
	if *weight > 0 {
		s.UpdateWeight(*field, *weight)
	}

	s.Regenerate(*platform, baseDir())

	md, _ := s.RenderMarkdown()
	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 已更新: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

func cmdPromptRef() {
	fs := flag.NewFlagSet("ref", flag.ExitOnError)
	image := fs.String("image", "", "参考图片路径")
	output := fs.String("output", "", "输出文件")
	fs.Parse(os.Args[3:])

	if *image == "" {
		fmt.Fprintln(os.Stderr, "请指定 --image ref.jpg")
		os.Exit(1)
	}

	// 生成带参考图的 prompt
	s, err := prompt.GeneratePrompt("参考图分析", "midjourney", *image, baseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "分析失败: %v\n", err)
		os.Exit(1)
	}

	md, _ := s.RenderMarkdown()
	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 参考图分析: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

func listPromptPlatforms() {
	platforms := listPlatforms()
	fmt.Println("可用平台 (36维度控制 v3.1 — 含视频参数):")
	for _, p := range platforms {
		fmt.Printf("  %s\n", p)
	}
}

// listPlatforms 返回所有有效平台列表
func listPlatforms() []string {
	return []string{
		"midjourney", "flux", "dalle3", "sd3", "stable_diffusion",
		"ideogram", "recraft", "seedance", "gemini", "jimeng",
		"veo", "sora", "kling", "pika",
		"video_generic", "video_camera", "video_keyframe",
	}
}

// isValidPlatform 检查平台名是否在支持列表中
func isValidPlatform(p string) bool {
	for _, v := range listPlatforms() {
		if v == p {
			return true
		}
	}
	return false
}

func cmdPromptGallery() {
	// 列出所有可用 prompt 模板
	bd := baseDir()
	galleryDir := bd + "/templates/prompts/"
	files, err := os.ReadDir(galleryDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取模板目录失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("📋 提示词模板市场:")
	fmt.Println()
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
			name := strings.TrimSuffix(strings.TrimSuffix(f.Name(), ".yaml"), ".yml")
			fmt.Printf("  📄 %-15s  via54 prompt --scene \"...\" --platform %s\n", name, name)
		}
	}
	fmt.Println()
	fmt.Println("自定义模板: 在 templates/prompts/ 下创建 YAML 文件即可注册")
}

func cmdComfyUI() {
	fs := flag.NewFlagSet("comfyui", flag.ExitOnError)
	workflowID := fs.String("workflow", "", "ComfyUI workflow template ID (e.g. sdxl_txt2img)")
	workflowPrompt := fs.String("prompt", "", "Generation prompt")
	negativePrompt := fs.String("negative", "", "Negative prompt")
	output := fs.String("output", "", "Output JSON path (default: stdout)")
	list := fs.Bool("list", false, "List available workflow templates")
	steps := fs.Int("steps", 0, "Override steps")
	cfg := fs.Float64("cfg", 0, "Override CFG scale")
	seed := fs.Int("seed", -1, "Override seed (-1 = random)")
	sampler := fs.String("sampler", "", "Override sampler name")
	keyframes := fs.String("keyframes", "", "Keyframe schedule for video: frame:prompt,frame:prompt (e.g. 0:a cat,8:a dog)")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	// ── List mode ──
	if *list {
		ids, err := workflow.ListWorkflowTemplates(bd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "列出工作流模板失败: %v\n", err)
			os.Exit(1)
		}
		reg, _ := workflow.LoadRegistry(bd)

		fmt.Println("📋 ComfyUI 工作流模板:")
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
				if err != nil {
					continue
				}
				fmt.Printf("  📄 %-20s %s\n", tmpl.ID, tmpl.Name)
				fmt.Printf("     %s\n", tmpl.Description)
				fmt.Printf("     模型: %s\n", tmpl.Model)
				fmt.Println()
			}
		}
		fmt.Println("用法: via54 comfyui --workflow <id> --prompt \"...\"")
		fmt.Println("       via54 comfyui --workflow sdxl_txt2img --prompt \"a cat\" --output workflow.json")
		fmt.Println("       via54 comfyui --workflow sdxl_txt2img --prompt \"...\" --steps 40 --cfg 10.0")
		return
	}

	// ── Build mode ──
	if *workflowID == "" {
		fmt.Fprintln(os.Stderr, "请指定 --workflow <模板ID> (或使用 --list 查看可用模板)")
		fmt.Fprintln(os.Stderr, "用法: via54 comfyui --workflow sdxl_txt2img --prompt \"...\"")
		os.Exit(1)
	}

	// Load the template
	tmpl, err := workflow.LoadWorkflowTemplate(*workflowID, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载工作流模板失败: %v\n", err)
		os.Exit(1)
	}

	// Collect overrides
	overrides := make(map[string]interface{})
	if *steps > 0 {
		overrides["steps"] = *steps
	}
	if *cfg > 0 {
		overrides["cfg"] = *cfg
	}
	if *seed >= 0 {
		overrides["seed"] = *seed
	}
	if *sampler != "" {
		overrides["sampler"] = *sampler
	}

	promptText := *workflowPrompt
	if promptText == "" {
		promptText = "empty prompt"
	}
	negText := *negativePrompt

	// Parse keyframes from CLI flag
	var kfs []workflow.Keyframe
	if *keyframes != "" {
		for _, kf := range strings.Split(*keyframes, ",") {
			kf = strings.TrimSpace(kf)
			if kf == "" {
				continue
			}
			parts := strings.SplitN(kf, ":", 2)
			if len(parts) == 2 {
				frame := 0
				fmt.Sscanf(parts[0], "%d", &frame)
				kfs = append(kfs, workflow.Keyframe{Frame: frame, Prompt: strings.TrimSpace(parts[1])})
			}
		}
	}

	// Build the workflow with keyframe support
	result, err := workflow.BuildWorkflow(tmpl, promptText, negText, overrides, kfs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "构建工作流失败: %v\n", err)
		os.Exit(1)
	}

	// Output
	if *output != "" {
		if err := os.WriteFile(*output, result.JSON, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "写入输出文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ ComfyUI 工作流已生成: %s\n", *output)
		fmt.Printf("   模板: %s (%s)\n", tmpl.ID, tmpl.Name)
		fmt.Printf("   类型: %s\n", tmpl.Type)
		fmt.Printf("   模型: %s\n", tmpl.Model)
		fmt.Printf("   节点数: %d\n", countNodes(result.JSON))
		if promptText != "empty prompt" {
			fmt.Printf("   提示词: %q\n", promptText)
		}
		if negText != "" {
			fmt.Printf("   负面提示词: %q\n", negText)
		}
		fmt.Println("   执行: comfy-swap run --workflow " + *output)
	} else {
		fmt.Print(string(result.JSON))
	}
}

// countNodes returns the number of top-level keys in a workflow JSON.
func countNodes(data []byte) int {
	var m map[string]interface{}
	if json.Unmarshal(data, &m) != nil {
		return 0
	}
	return len(m)
}

// ─── Assess ───
func cmdPromptAssess() {
	fs := flag.NewFlagSet("assess", flag.ExitOnError)
	image := fs.String("image", "", "生图输出路径")
	promptText := fs.String("prompt", "", "使用的 prompt")
	fs.Parse(os.Args[3:])

	if *image == "" {
		fmt.Fprintln(os.Stderr, "请指定 --image output.png")
		os.Exit(1)
	}

	report := prompt.AssessImage(*image, *promptText)
	fmt.Printf("📊 生图质量评估: %s\n", *image)
	fmt.Printf("  综合评分: %.2f / 1.0\n", report.OverallScore)
	fmt.Printf("  清晰度:   %.2f\n", report.ClarityScore)
	fmt.Printf("  构图:     %.2f\n", report.CompositionScore)
	fmt.Printf("  色彩:     %.2f\n", report.ColorScore)
	fmt.Printf("  Prompt匹配: %.2f\n", report.PromptMatch)
	if len(report.Issues) > 0 {
		fmt.Println("  问题:")
		for _, iss := range report.Issues {
			fmt.Printf("    ⚠️ %s\n", iss)
		}
	}
}

// ─── Version ───
func cmdPromptVersion() {
	fs := flag.NewFlagSet("version", flag.ExitOnError)
	save := fs.Bool("save", false, "保存版本")
	promptFile := fs.String("prompt", "", "prompt JSON 文件")
	list := fs.Bool("list", false, "列出版本")
	versionDir := fs.String("dir", ".via54/prompts", "版本目录")
	diffV1 := fs.String("diff-v1", "", "比较版本 v1")
	diffV2 := fs.String("diff-v2", "", "比较版本 v2")
	fs.Parse(os.Args[3:])

	bd := baseDir()
	dir := *versionDir
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(bd, dir)
	}

	if *save {
		if *promptFile == "" {
			fmt.Fprintln(os.Stderr, "请指定 --prompt prompt.json")
			os.Exit(1)
		}
		data, _ := os.ReadFile(*promptFile)
		var s prompt.PromptScaffold
		json.Unmarshal(data, &s)
		fp, err := prompt.SaveVersion(dir, &s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "保存失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 已保存: %s\n", fp)
	} else if *list {
		versions := prompt.ListVersions(dir)
		if len(versions) == 0 {
			fmt.Println("无版本记录")
			return
		}
		fmt.Println("📋 版本列表:")
		for _, v := range versions {
			fmt.Printf("  %s\n", v)
		}
	} else if *diffV1 != "" && *diffV2 != "" {
		diff, err := prompt.DiffVersions(dir, *diffV1, *diffV2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "比较失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(diff)
	} else {
		fmt.Fprintln(os.Stderr, "请指定操作: --save, --list, 或 --diff-v1 --diff-v2")
		os.Exit(1)
	}
}

// ─── Send (BrowserWing 兼容) ───
func cmdPromptSend() {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	promptFile := fs.String("prompt", "", "prompt JSON 文件")
	platform := fs.String("platform", "midjourney", "目标平台")
	fs.Parse(os.Args[3:])

	if *promptFile == "" {
		fmt.Fprintln(os.Stderr, "请指定 --prompt prompt.json")
		os.Exit(1)
	}

	data, _ := os.ReadFile(*promptFile)
	var s prompt.PromptScaffold
	json.Unmarshal(data, &s)

	// 输出 BrowserWing 兼容的指令
	fmt.Printf(`🌐 使用 BrowserWing (⭐1,292) 自动提交到 %s:

1. 确保 BrowserWing 已安装: npx skills add browserwing/browserwing
2. 执行以下指令:

   browserwing navigate "https://www.midjourney.com"
   browserwing type "#prompt-input" "%s"
   browserwing click "#generate-button"
   browserwing wait 30

或直接复制 prompt 到 %s:
   %s
`, *platform, s.FinalPrompt, *platform, s.FinalPrompt)
}
