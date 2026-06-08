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
	"github.com/veawho/via54Design/internal/pattern"
	"os"
)

func cmdPattern() {
	fs := flag.NewFlagSet("pattern", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	name := fs.String("name", "unnamed", "作品名称")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" {
		fmt.Fprintln(os.Stderr, "请指定 --html")
		os.Exit(1)
	}
	data, _ := os.ReadFile(*htmlFile)
	p, yaml := pattern.ExtractFromHTML(string(data), *name)
	fmt.Printf("\n=== 模式提取: %s ===\n", *name)
	fmt.Printf("  🎨 %d colors | 📐 %v | 🎞️ %s\n", p.Colors.TotalUnique, p.Layout.Types, p.Animations.Complexity)
	fmt.Printf("  🔤 Display: %s  Body: %s\n", p.Fonts.Display, p.Fonts.Body)
	fmt.Printf("\n─── YAML 候选 ───\n%s\n", yaml)
}
