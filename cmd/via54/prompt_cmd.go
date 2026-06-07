// via54Design — 图片提示词 (Prompt) 命令
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/veawho/via54Design/internal/prompt"
)

func cmdPrompt() {
	fs := flag.NewFlagSet("prompt", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述 (必填)")
	platform := fs.String("platform", "midjourney", "平台: midjourney/kling/jimeng/gemini")
	output := fs.String("output", "", "输出文件 (默认 stdout)")
	listPlatforms := fs.Bool("list", false, "列出所有平台")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	if *listPlatforms {
		fmt.Println("可用平台:")
		for _, p := range []string{"midjourney", "kling", "jimeng", "gemini"} {
			fmt.Printf("  %-15s  %s\n", p, platformDesc(p))
		}
		return
	}
	if *scene == "" {
		fmt.Fprintln(os.Stderr, "请指定 --scene \"场景描述\"")
		fmt.Fprintln(os.Stderr, "或使用 --list 查看可用平台")
		os.Exit(1)
	}

	scaffold, err := prompt.GeneratePrompt(*scene, *platform, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "提示词生成失败: %v\n", err)
		os.Exit(1)
	}

	md, err := scaffold.RenderMarkdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "渲染失败: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 提示词已保存: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

func platformDesc(p string) string {
	switch p {
	case "midjourney": return "Midjourney 图片生成"
	case "kling":      return "可灵 AI 视频/图片"
	case "jimeng":     return "即梦 AI 图片生成"
	case "gemini":     return "Google Gemini Imagen"
	default:           return ""
	}
}
