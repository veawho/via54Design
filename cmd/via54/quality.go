// SPDX-License-Identifier: MIT OR AGPL-3.0

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


