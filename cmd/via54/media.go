// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"github.com/veawho/via54Design/internal/media"
	"os"
	"strings"
)

func cmdMedia() {
	// --help / -h 应在父命令层处理 (CLI 惯例)
	if len(os.Args) >= 3 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
		fmt.Println("用法: via54 media <子命令> [选项]")
		fmt.Println()
		fmt.Println("子命令:")
		fmt.Println("  add-music    给视频添加 BGM")
		fmt.Println("    --mood     配乐 mood (tech/calm/epic) (default \"tech\")")
		fmt.Println("    --output   输出路径")
		fmt.Println("  convert      转 60fps/GIF")
		fmt.Println("  fetch        批量下载图片")
		fmt.Println("    --query    关键词 (逗号分隔)")
		fmt.Println("    --out      输出目录 (default \"./img\")")
		fmt.Println("    --count    每关键词张数 (default 2)")
		fmt.Println("  trace        图片矢量化 (vtracer)")
		fmt.Println("    --input    输入图片 (JPG/PNG)")
		fmt.Println("    --output   输出 SVG 路径")
		os.Exit(0)
	}
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: via54 media [add-music|convert|fetch|trace] (用 --help 查看详情)")
		os.Exit(1)
	}
	switch os.Args[2] {
	case "add-music":
		fs := flag.NewFlagSet("add-music", flag.ExitOnError)
		mood := fs.String("mood", "tech", "配乐 mood")
		output := fs.String("output", "", "输出路径")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "请指定 input.mp4")
			os.Exit(1)
		}
		if err := media.AddMusic(fs.Arg(0), *mood, *output); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ BGM 添加完成")

	case "convert":
		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "请指定 input.mp4")
			os.Exit(1)
		}
		mp4, gif, err := media.ConvertFormats(fs.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 60fps: %s\n", mp4)
		if gif != "" {
			fmt.Printf("✅ GIF:   %s\n", gif)
		}

	case "fetch":
		fs := flag.NewFlagSet("fetch", flag.ExitOnError)
		query := fs.String("query", "", "关键词 (逗号分隔)")
		out := fs.String("out", "./img", "输出目录")
		count := fs.Int("count", 2, "每关键词张数")
		fs.Parse(os.Args[3:])
		if *query == "" {
			fmt.Fprintln(os.Stderr, "请指定 --query")
			os.Exit(1)
		}
		queries := strings.Split(*query, ",")
		results, err := media.FetchImages(queries, *out, *count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 下载 %d 张到 %s\n", len(results), *out)

	case "trace":
		fs := flag.NewFlagSet("trace", flag.ExitOnError)
		input := fs.String("input", "", "输入图片路径 (JPG/PNG)")
		output := fs.String("output", "", "输出 SVG 路径 (可选)")
		fs.Parse(os.Args[3:])
		if *input == "" {
			fmt.Fprintln(os.Stderr, "请指定 --input")
			os.Exit(1)
		}
		svgPath, err := media.TraceImage(*input, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
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
