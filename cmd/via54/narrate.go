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
	"github.com/veawho/via54Design/internal/narrate"
	"os"
)

func cmdNarrate() {
	fs := flag.NewFlagSet("narrate", flag.ExitOnError)
	seed := fs.String("seed", "", "一句话种子 (必填)")
	model := fs.String("model", "three-act", "叙事模型ID")
	duration := fs.Int("duration", 30, "目标视频时长(秒)")
	output := fs.String("output", "", "输出文件路径 (默认: stdout)")
	format := fs.String("format", "markdown", "输出格式: markdown / json")
	listModels := fs.Bool("list", false, "列出所有叙事模型")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	if *listModels {
		out, err := narrate.ListModels(bd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1)
		}
		fmt.Print(out)
		return
	}
	if *seed == "" {
		fmt.Fprintln(os.Stderr, "请指定 --seed \"一句话描述\"")
		fmt.Fprintln(os.Stderr, "或使用 --list 查看可用叙事模型")
		os.Exit(1)
	}

	scaffold, err := narrate.GenerateScaffold(*seed, *model, *duration, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "叙事生成失败: %v\n", err)
		os.Exit(1)
	}

	if *format == "json" {
		js, err := scaffold.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON 序列化失败: %v\n", err)
			os.Exit(1)
		}
		if *output != "" {
			os.WriteFile(*output, []byte(js), 0644)
			fmt.Printf("✅ 叙事脚手架 JSON 已保存: %s\n", *output)
		} else {
			fmt.Print(js)
		}
		return
	}

	md, err := scaffold.RenderMarkdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "渲染失败: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 叙事脚手架已保存: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

// ─── Media (Shell+Python 迁移) ───

