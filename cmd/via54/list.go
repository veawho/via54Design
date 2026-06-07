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
	"fmt"
	"os"
	"github.com/veawho/via54Design/internal/template"
	"strings"
)

func cmdList() {
	eng, err := template.NewEngine(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1) }
	for cat, entries := range eng.Registry.ListAll() {
		fmt.Printf("\n=== %s ===\n", cat)
		// 叙事学分组显示
		if cat == "narratology" {
			fmt.Println("  ── guides ──")
			for _, e := range entries {
				if e.Category == "narratology/guide" || e.Category == "" {
					fmt.Printf("    %-32s %s\n", e.ID, e.Name)
				}
			}
			fmt.Println("  ── models ──")
			for _, e := range entries {
				if e.Category == "narratology/model" {
					extra := ""
					if len(e.Tags) > 0 {
						extra = fmt.Sprintf("  [%s]", strings.Join(e.Tags, ", "))
					}
					fmt.Printf("    %-32s %s%s\n", e.ID, e.Name, extra)
				}
			}
			continue
		}
		for _, e := range entries {
			fmt.Printf("  %-32s %s\n", e.ID, e.Name)
		}
	}
}

// ─── Narrate (叙事引擎) ───

