// SPDX-License-Identifier: MIT OR AGPL-3.0

package main

import (
	"flag"
	"fmt"
	"github.com/veawho/via54Design/internal/media"
	"os"
	"strings"
)

func cmdMedia() {
	if len(os.Args) < 3 { fmt.Fprintln(os.Stderr, "用法: via54 media [add-music|convert|fetch|trace]"); os.Exit(1) }
	switch os.Args[2] {
	case "add-music":
		fs := flag.NewFlagSet("add-music", flag.ExitOnError)
		mood := fs.String("mood", "tech", "配乐 mood")
		output := fs.String("output", "", "输出路径")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.mp4"); os.Exit(1) }
		if err := media.AddMusic(fs.Arg(0), *mood, *output); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1)
		}
		fmt.Println("✅ BGM 添加完成")

	case "convert":
		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.mp4"); os.Exit(1) }
		mp4, gif, err := media.ConvertFormats(fs.Arg(0))
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ 60fps: %s\n", mp4)
		if gif != "" { fmt.Printf("✅ GIF:   %s\n", gif) }

	case "fetch":
		fs := flag.NewFlagSet("fetch", flag.ExitOnError)
		query := fs.String("query", "", "关键词 (逗号分隔)")
		out := fs.String("out", "./img", "输出目录")
		count := fs.Int("count", 2, "每关键词张数")
		fs.Parse(os.Args[3:])
		if *query == "" { fmt.Fprintln(os.Stderr, "请指定 --query"); os.Exit(1) }
		queries := strings.Split(*query, ",")
		results, err := media.FetchImages(queries, *out, *count)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ 下载 %d 张到 %s\n", len(results), *out)

	case "trace":
		fs := flag.NewFlagSet("trace", flag.ExitOnError)
		input := fs.String("input", "", "输入图片路径 (JPG/PNG)")
		output := fs.String("output", "", "输出 SVG 路径 (可选)")
		fs.Parse(os.Args[3:])
		if *input == "" { fmt.Fprintln(os.Stderr, "请指定 --input"); os.Exit(1) }
		svgPath, err := media.TraceImage(*input, nil)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		if *output != "" {
			os.Rename(svgPath, *output)
			svgPath = *output
		}
		fmt.Printf("✅ SVG: %s\n", svgPath)
		default:
		fmt.Fprintf(os.Stderr, "未知 media 命令: %s (支持: add-music, convert, fetch)\n", os.Args[2])
		os.Exit(1)
	}
}

// ─── Export (全 Go 导出管线) ───

