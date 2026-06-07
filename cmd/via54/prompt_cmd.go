// via54Design — 图片提示词 (Prompt) 命令 v2
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/veawho/via54Design/internal/prompt"
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
	fmt.Println("可用平台 (16维度控制):")
	for _, p := range []struct{n, d string}{
		{"midjourney", "Midjourney 图片生成 — 16维度 + 负面词库 + Token权重"},
		{"kling", "可灵 AI 视频/图片"},
		{"jimeng", "即梦 AI 图片生成"},
		{"gemini", "Google Gemini Imagen"},
	} {
		fmt.Printf("  %-15s  %s\n", p.n, p.d)
	}
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
	workflow := fs.String("workflow", "", "ComfyUI workflow JSON")
	workflowPrompt := fs.String("prompt", "", "生图 prompt")
	output := fs.String("output", "comfyui_output.json", "输出路径")
	fs.Parse(os.Args[3:])

	if *workflow == "" || *workflowPrompt == "" {
		fmt.Fprintln(os.Stderr, "请指定 --workflow workflow.json --prompt \"...\"")
		fmt.Fprintln(os.Stderr, "ComfyUI 工作流模板: docs/prompts/comfyui-workflows/")
		os.Exit(1)
	}

	// 读取 workflow 并注入 prompt
	data, err := os.ReadFile(*workflow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取 workflow 失败: %v\n", err)
		os.Exit(1)
	}

	// 注入 prompt 到 workflow 的 CLIPTextEncode 节点
	var wf map[string]interface{}
	json.Unmarshal(data, &wf)

	injected := 0
	for _, node := range wf {
		n, ok := node.(map[string]interface{})
		if !ok { continue }
		cls, _ := n["class_type"].(string)
		if cls == "CLIPTextEncode" {
			if inputs, ok := n["inputs"].(map[string]interface{}); ok {
				inputs["text"] = *workflowPrompt
				injected++
			}
		}
	}

	outData, _ := json.MarshalIndent(wf, "", "  ")
	os.WriteFile(*output, outData, 0644)
	fmt.Printf("✅ ComfyUI workflow 已注入 prompt (%d 个节点): %s\n", injected, *output)
	fmt.Println("  执行: comfy-swap run --workflow " + *output)
}
