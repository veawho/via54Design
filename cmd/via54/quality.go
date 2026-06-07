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
	"os"
	"github.com/veawho/via54Design/internal/quality"
)

func cmdQuality() {
	fs := flag.NewFlagSet("quality", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	verbose := fs.Bool("verbose", false, "显示 info")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, _ := os.ReadFile(*htmlFile)
	report := quality.CheckHTML(string(data))
	fmt.Printf("\n=== 质量门禁: %s ===\n", report.Verdict)
	fmt.Printf("文件: %d bytes | CSS块: %d | 行: %d\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	fmt.Printf("问题: %d 错误 / %d 警告 / %d 信息\n\n", report.Summary["error"], report.Summary["warning"], report.Summary["info"])
	for _, iss := range report.Issues {
		if iss.Severity == "info" && !*verbose { continue }
		icon := map[string]string{"error": "❌", "warning": "⚠️", "info": "ℹ️"}[iss.Severity]
		fmt.Printf("  %s [%s] %s\n", icon, iss.Category, iss.Message)
	}
}


