// SPDX-License-Identifier: MIT OR AGPL-3.0

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

