// via54Design — Markdown 幻灯片导出 (slidev/marp 兼容)
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/veawho/via54Design/internal/narrate"
)

func cmdPresent() {
	fs := flag.NewFlagSet("present", flag.ExitOnError)
	title := fs.String("title", "Presentation", "Slide deck title")
	seed := fs.String("seed", "", "Story seed for narrative (optional)")
	layout := fs.String("layout", "hero-split-16-9", "Template layout")
	color := fs.String("color", "ink-wash", "Color scheme")
	font := fs.String("font", "ming-hei-editorial", "Font")
	format := fs.String("format", "marp", "Output format: marp, slidev, revealjs")
	slides := fs.Int("slides", 5, "Number of slides")
	output := fs.String("output", "", "Output file (default: stdout)")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	// Build slide content
	var slideContent []string
	if *seed != "" {
		s, err := narrate.GenerateScaffold(*seed, "three-act", 30, bd)
		if err == nil && s != nil {
			// Use narrative structure for slides
			acts := []string{"Act 1: Setup", "Act 2: Confrontation", "Act 3: Resolution"}
			for idx, act := range acts {
				if idx < *slides/3+1 {
					slideContent = append(slideContent, fmt.Sprintf("## %s\n\n%s", act, "via54Design"))
				}
			}
		}
	}
	if len(slideContent) == 0 {
		for slideNum := 1; slideNum <= *slides; slideNum++ {
			slideContent = append(slideContent, fmt.Sprintf("## Slide %d\n\nContent for slide %d.", slideNum, slideNum))
		}
	}

	var outputText string
	switch *format {
	case "slidev":
		// slidev format (Vite-based, markdown frontmatter)
		outputText = fmt.Sprintf(`---
theme: default
title: %s
layout: default
---
`, *title)
		for _, content := range slideContent {
			outputText += fmt.Sprintf("\n---\n%s\n", content)
		}
		outputText += "\n---\n# Thank You\n\nBuilt with via54Design + Slidev\n"
		// slidev output

	case "revealjs":
		// reveal.js format (HTML sections)
		outputText = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <title>%s</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/reveal.css">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/theme/night.css">
</head>
<body>
<div class="reveal">
<div class="slides">
<section><h1>%s</h1></section>
`, *title, *title)
		for _, content := range slideContent {
			outputText += fmt.Sprintf("<section>%s</section>\n", content)
		}
		outputText += `<section><h1>Thank You</h1></section>
</div></div>
<script src="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/reveal.js"></script>
<script>Reveal.initialize();</script>
</body>
</html>`

	default:
		// marp format (default) - CommonMark + --- separators
		var b strings.Builder
		b.WriteString(fmt.Sprintf("---\nmarp: true\ntitle: %s\ntheme: uncover\n---\n\n", *title))
		b.WriteString(fmt.Sprintf("# %s\n\n", *title))
		b.WriteString(fmt.Sprintf("<!-- \n  via54Design + Marp\n  layout: %s\n  color: %s\n  font: %s\n-->\n\n", *layout, *color, *font))
		for _, content := range slideContent {
			b.WriteString("---\n\n")
			b.WriteString(content)
			b.WriteString("\n\n")
		}
		b.WriteString("---\n\n# Thank You\n\nBuilt with via54Design + Marp\n")
		outputText = b.String()
	}

	if *output != "" {
		os.WriteFile(*output, []byte(outputText), 0644)
		fmt.Printf("✅ Presentation: %s (%s format, %d slides)\n", *output, *format, *slides)
	} else {
		fmt.Print(outputText)
	}
}
