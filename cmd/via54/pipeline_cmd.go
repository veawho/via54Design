// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/veawho/via54Design/internal/pipeline"
)

func cmdPipeline() {
	if len(os.Args) < 3 {
		pipelineHelp()
		return
	}

	switch os.Args[2] {
	case "expand":
		cmdPipelineExpand()
	case "translate":
		cmdPipelineTranslate()
	case "variant":
		cmdPipelineVariant()
	case "archive":
		cmdPipelineArchive()
	case "export":
		cmdPipelineExport()
	case "reverse":
		cmdPipelineReverse()
	default:
		pipelineHelp()
	}
}

func pipelineHelp() {
	fmt.Println("用法: via54 prompt <子命令> [选项]")
	fmt.Println()
	fmt.Println("子命令:")
	fmt.Println("  expand          扩展现有场景 (LLM + i18n)")
	fmt.Println("    --scene \"...\"  --platform midjourney [--provider openai] [--key ...]")
	fmt.Println("  translate       翻译中文场景为英文")
	fmt.Println("    --text \"中文场景\" [--provider openai]")
	fmt.Println("  variant         展开 {opt1|opt2} 变体语法")
	fmt.Println("    --scene \"{cat|dog}\" --count 4")
	fmt.Println("  archive         提示词存档管理")
	fmt.Println("    --save --scene \"...\" --tags \"tag1,tag2\"")
	fmt.Println("    --search \"keyword\"")
	fmt.Println("    --list")
	fmt.Println("    --delete <id>")
	fmt.Println("  export          导出为 A1111 / ComfyUI 格式")
	fmt.Println("    --scene \"...\" --platform midjourney --format a1111|comfyui")
	fmt.Println("  reverse         反向图片→提示词 (Vision LLM)")
	fmt.Println("    --image path.jpg [--provider openai]")
	fmt.Println()
	fmt.Println("提供者选项: --provider openai|deepseek|ollama|hermes|local")
	fmt.Println("环境变量: VIA54_LLM_ENDPOINT, VIA54_LLM_KEY, VIA54_LLM_MODEL")
}

// ── expand ──

func cmdPipelineExpand() {
	fs := flag.NewFlagSet("expand", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述")
	platform := fs.String("platform", "midjourney", "目标平台")
	provider := fs.String("provider", "openai", "LLM 提供者")
	endpoint := fs.String("endpoint", "", "API 端点 (覆盖 VIA54_LLM_ENDPOINT)")
	key := fs.String("key", "", "API key (覆盖 VIA54_LLM_KEY)")
	model := fs.String("model", "", "模型名 (覆盖 VIA54_LLM_MODEL)")
	output := fs.String("output", "", "输出 JSON 文件")
	fs.Parse(os.Args[3:])

	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"场景描述\"")
		os.Exit(1)
	}

	cfg := pipeline.ConfigFromEnv(*provider)
	if *endpoint != "" {
		cfg.LLMEndpoint = *endpoint
	}
	if *key != "" {
		cfg.LLMKey = *key
	}
	if *model != "" {
		cfg.LLMModel = *model
	}

	if cfg.LLMKey == "" && cfg.ProviderRequiresKey() {
		fmt.Fprintf(os.Stderr, "错误: 提供者 '%s' 需要 API key。设置 VIA54_LLM_KEY 或 --key。\n", *provider)
		os.Exit(1)
	}

	result, err := pipeline.Pipeline(*scene, *platform, *provider, cfg.LLMEndpoint, cfg.LLMKey, cfg.LLMModel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "管道错误: %v\n", err)
		os.Exit(1)
	}

	data, _ := pipeline.ScaffoldToJSON(result)
	if *output != "" {
		os.WriteFile(*output, data, 0644)
		fmt.Printf("✅ 已保存: %s\n", *output)
	} else {
		fmt.Print(string(data))
	}
}

// ── translate ──

func cmdPipelineTranslate() {
	fs := flag.NewFlagSet("translate", flag.ExitOnError)
	text := fs.String("text", "", "要翻译的文本")
	provider := fs.String("provider", "openai", "LLM 提供者")
	endpoint := fs.String("endpoint", "", "API 端点")
	key := fs.String("key", "", "API key")
	model := fs.String("model", "", "模型名")
	fs.Parse(os.Args[3:])

	if *text == "" {
		fmt.Fprintln(os.Stderr, "请指定 --text \"中文场景\"")
		os.Exit(1)
	}

	cfg := pipeline.ConfigFromEnv(*provider)
	if *endpoint != "" {
		cfg.LLMEndpoint = *endpoint
	}
	if *key != "" {
		cfg.LLMKey = *key
	}
	if *model != "" {
		cfg.LLMModel = *model
	}

	// Only require API key if the text actually contains Chinese
	needsKey := cfg.ProviderRequiresKey() && pipeline.ContainsChinese(*text)
	if cfg.LLMKey == "" && needsKey {
		fmt.Fprintf(os.Stderr, "错误: 提供者 '%s' 需要 API key 来翻译中文。\n", *provider)
		os.Exit(1)
	}

	result, wasTranslated, err := pipeline.TranslateToEnglish(*text, func(messages []pipeline.ChatMessage) (string, error) {
		return pipeline.CallLLM(messages, cfg.LLMEndpoint, cfg.LLMKey, cfg.LLMModel)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "翻译错误: %v\n", err)
		os.Exit(1)
	}

	if wasTranslated {
		fmt.Println(result)
	} else {
		fmt.Println("(无需翻译 — 输入不包含中文)")
		fmt.Println(result)
	}
}

// ── variant ──

func cmdPipelineVariant() {
	fs := flag.NewFlagSet("variant", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述 (含 {opt1|opt2} 语法)")
	count := fs.Int("count", 4, "生成变体数量")
	fs.Parse(os.Args[3:])

	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"{cat|dog} in a park\"")
		os.Exit(1)
	}

	variants := pipeline.ExpandVariants(*scene, *count)
	fmt.Printf("变体 (%d):\n", len(variants))
	for i, v := range variants {
		fmt.Printf("  %d. %s\n", i+1, v)
	}
}

// ── archive ──

func cmdPipelineArchive() {
	if len(os.Args) < 4 {
		archiveHelp()
		return
	}

	switch os.Args[3] {
	case "save":
		cmdArchiveSave()
	case "list":
		cmdArchiveList()
	case "search":
		cmdArchiveSearch()
	case "delete":
		cmdArchiveDelete()
	default:
		archiveHelp()
	}
}

func archiveHelp() {
	fmt.Println("用法: via54 prompt archive <子命令> [选项]")
	fmt.Println()
	fmt.Println("子命令:")
	fmt.Println("  save             保存提示词到存档")
	fmt.Println("    --scene \"...\" --platform midjourney [--tags \"tag1,tag2\"]")
	fmt.Println("  list             列出最近条目")
	fmt.Println("    [--limit 20]")
	fmt.Println("  search           搜索存档")
	fmt.Println("    --query \"keyword\" [--limit 10]")
	fmt.Println("  delete           删除条目")
	fmt.Println("    --id <record_id>")
}

func cmdArchiveSave() {
	fs := flag.NewFlagSet("save", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述")
	platform := fs.String("platform", "midjourney", "目标平台")
	tags := fs.String("tags", "", "逗号分隔的标签")
	provider := fs.String("provider", "openai", "LLM 提供者")
	endpoint := fs.String("endpoint", "", "API 端点")
	key := fs.String("key", "", "API key")
	model := fs.String("model", "", "模型名")
	output := fs.String("output", "", "输出 JSON 文件")
	fs.Parse(os.Args[4:])

	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"...\"")
		os.Exit(1)
	}

	cfg := pipeline.ConfigFromEnv(*provider)
	if *endpoint != "" {
		cfg.LLMEndpoint = *endpoint
	}
	if *key != "" {
		cfg.LLMKey = *key
	}
	if *model != "" {
		cfg.LLMModel = *model
	}

	if cfg.LLMKey == "" && cfg.ProviderRequiresKey() {
		// No API key — save raw scene without LLM expansion
		result := &pipeline.PromptScaffold{
			Scene:     *scene,
			Platform:  *platform,
			Fields:    make(map[string]string),
			RawPrompt: *scene,
		}
		saveToArchive(result, *output, *tags)
		return
	}

	// Run pipeline with LLM
	result, err := pipeline.Pipeline(*scene, *platform, *provider, cfg.LLMEndpoint, cfg.LLMKey, cfg.LLMModel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "管道错误: %v\n", err)
		os.Exit(1)
	}

	saveToArchive(result, *output, *tags)
}

// saveToArchive saves a scaffold to file and archive.
func saveToArchive(result *pipeline.PromptScaffold, output string, tags string) {
	if output != "" {
		data, _ := pipeline.ScaffoldToJSON(result)
		os.WriteFile(output, data, 0644)
		fmt.Printf("✅ 已保存: %s\n", output)
	}

	// Save to archive
	var tagList []string
	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tagList = append(tagList, t)
			}
		}
	}

	arch := pipeline.NewArchive("")
	id, err := arch.Save(result, tagList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "存档错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 已存档: id=%s\n", id)
	if result.RawPrompt != "" {
		preview := result.RawPrompt
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("   提示词: %s\n", preview)
	}
}

func cmdArchiveList() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	limit := fs.Int("limit", 20, "显示条目数")
	fs.Parse(os.Args[4:])

	arch := pipeline.NewArchive("")
	records, err := arch.List(*limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "列出存档错误: %v\n", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		fmt.Println("存档为空。")
		return
	}

	fmt.Printf("最近 %d 条存档条目:\n\n", len(records))
	for _, r := range records {
		tagsStr := "(无标签)"
		if len(r.Tags) > 0 {
			tagsStr = strings.Join(r.Tags, ", ")
		}
		scenePreview := r.Scene
		if len(scenePreview) > 60 {
			scenePreview = scenePreview[:60] + "..."
		}
		created := r.CreatedAt
		if len(created) > 19 {
			created = created[:19]
		}
		fmt.Printf("  [%s] %s | %s\n", r.ID, created, scenePreview)
		fmt.Printf("        标签: %s\n", tagsStr)
		fmt.Println()
	}
}

func cmdArchiveSearch() {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	query := fs.String("query", "", "搜索关键词")
	limit := fs.Int("limit", 10, "最大结果数")
	fs.Parse(os.Args[4:])

	if *query == "" {
		fmt.Fprintln(os.Stderr, "请指定 --query \"关键词\"")
		os.Exit(1)
	}

	arch := pipeline.NewArchive("")
	results, err := arch.Search(*query, *limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "搜索错误: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Printf("未找到匹配 '%s' 的结果。\n", *query)
		return
	}

	fmt.Printf("找到 %d 个结果 (关键词: '%s'):\n\n", len(results), *query)
	for _, r := range results {
		tagsStr := "(无标签)"
		if len(r.Tags) > 0 {
			tagsStr = strings.Join(r.Tags, ", ")
		}
		scenePreview := r.Scene
		if len(scenePreview) > 80 {
			scenePreview = scenePreview[:80] + "..."
		}
		created := r.CreatedAt
		if len(created) > 19 {
			created = created[:19]
		}
		fmt.Printf("  [%s] %s\n", r.ID, created)
		fmt.Printf("        场景: %s\n", scenePreview)
		fmt.Printf("        标签: %s\n", tagsStr)
		fmt.Println()
	}
}

func cmdArchiveDelete() {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.String("id", "", "记录 ID")
	fs.Parse(os.Args[4:])

	if *id == "" {
		fmt.Fprintln(os.Stderr, "请指定 --id <record_id>")
		os.Exit(1)
	}

	arch := pipeline.NewArchive("")
	found, err := arch.Delete(*id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "删除错误: %v\n", err)
		os.Exit(1)
	}

	if found {
		fmt.Printf("✅ 已删除记录: %s\n", *id)
	} else {
		fmt.Fprintf(os.Stderr, "记录未找到: %s\n", *id)
		os.Exit(1)
	}
}

// ── export ──

func cmdPipelineExport() {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述")
	platform := fs.String("platform", "midjourney", "目标平台")
	format := fs.String("format", "a1111", "输出格式: a1111|comfyui")
	provider := fs.String("provider", "openai", "LLM 提供者")
	endpoint := fs.String("endpoint", "", "API 端点")
	key := fs.String("key", "", "API key")
	model := fs.String("model", "", "模型名")
	fs.Parse(os.Args[3:])

	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"...\"")
		os.Exit(1)
	}

	cfg := pipeline.ConfigFromEnv(*provider)
	if *endpoint != "" {
		cfg.LLMEndpoint = *endpoint
	}
	if *key != "" {
		cfg.LLMKey = *key
	}
	if *model != "" {
		cfg.LLMModel = *model
	}

	if cfg.LLMKey == "" && cfg.ProviderRequiresKey() {
		fmt.Fprintf(os.Stderr, "错误: 提供者 '%s' 需要 API key。\n", *provider)
		os.Exit(1)
	}

	result, err := pipeline.Pipeline(*scene, *platform, *provider, cfg.LLMEndpoint, cfg.LLMKey, cfg.LLMModel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "管道错误: %v\n", err)
		os.Exit(1)
	}

	switch *format {
	case "a1111":
		fmt.Print(pipeline.ExportA1111(result))
	case "comfyui":
		fmt.Print(pipeline.ExportComfyUI(result, "6"))
	default:
		fmt.Fprintf(os.Stderr, "未知格式: %s (支持: a1111, comfyui)\n", *format)
		os.Exit(1)
	}
}

// ── reverse ──

func cmdPipelineReverse() {
	fs := flag.NewFlagSet("reverse", flag.ExitOnError)
	image := fs.String("image", "", "图片路径")
	provider := fs.String("provider", "openai", "LLM 提供者")
	endpoint := fs.String("endpoint", "", "API 端点")
	key := fs.String("key", "", "API key")
	model := fs.String("model", "", "模型名")
	fs.Parse(os.Args[3:])

	if *image == "" {
		fmt.Fprintln(os.Stderr, "请指定 --image path.jpg")
		os.Exit(1)
	}

	cfg := pipeline.ConfigFromEnv(*provider)
	if *endpoint != "" {
		cfg.LLMEndpoint = *endpoint
	}
	if *key != "" {
		cfg.LLMKey = *key
	}
	if *model != "" {
		cfg.LLMModel = *model
	}

	if cfg.LLMKey == "" && cfg.ProviderRequiresKey() {
		fmt.Fprintf(os.Stderr, "错误: 提供者 '%s' 需要 API key。\n", *provider)
		os.Exit(1)
	}

	expansion, err := pipeline.ReverseImage(*image, cfg.LLMKey, cfg.LLMEndpoint, cfg.LLMModel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "反向分析错误: %v\n", err)
		os.Exit(1)
	}

	scaffold := &pipeline.PromptScaffold{
		Scene:    fmt.Sprintf("[Reverse from: %s]", *image),
		Platform: "midjourney",
		Fields:   expansion.Fields,
		Negative: expansion.Negative,
	}
	scaffold.RawPrompt = pipeline.BuildRawPrompt(scaffold)

	data, _ := pipeline.ScaffoldToJSON(scaffold)
	var out map[string]interface{}
	json.Unmarshal(data, &out)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}
