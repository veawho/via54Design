// SPDX-License-Identifier: MIT OR AGPL-3.0

package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/veawho/via54Design/internal/pattern"
)

func cmdPattern() {
	fs := flag.NewFlagSet("pattern", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	name := fs.String("name", "unnamed", "作品名称")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, _ := os.ReadFile(*htmlFile)
	p, yaml := pattern.ExtractFromHTML(string(data), *name)
	fmt.Printf("\n=== 模式提取: %s ===\n", *name)
	fmt.Printf("  🎨 %d colors | 📐 %v | 🎞️ %s\n", p.Colors.TotalUnique, p.Layout.Types, p.Animations.Complexity)
	fmt.Printf("  🔤 Display: %s  Body: %s\n", p.Fonts.Display, p.Fonts.Body)
	fmt.Printf("\n─── YAML 候选 ───\n%s\n", yaml)
}


