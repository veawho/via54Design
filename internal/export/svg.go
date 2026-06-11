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

// via54Design — SVG 矢量导出 (v2: §12 design-audit 规范)
//
// 规范: viewBox=680 (16:9 画布内部坐标系, 缩放属性保 width/height)
//       class 系统: t (正文 12px) / ts (副标 14px) / th (标题 24px)
//       fill 默认 none, 文字/text 元素显式 fill
//       元素 stroke-width 标准 1.5
//       text-anchor 在需要居中/居右时显式
//       字号 12/14 (小) / 24 (大)
package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SVGScene SVG 场景数据
type SVGScene struct {
	Title       string
	Voiceover   string
	Body        string
	Mood        string
	BeatName    string
	SceneNo     int
	TotalScenes int
	Duration    int
	Width       int
	Height      int
}

// §12 设计规范常量
const (
	SVGViewBoxW = 680  // §12 内部画布宽
	SVGViewBoxH = 382  // §12 内部画布高 (680/16*9 ≈ 382, 16:9)
	SVGFontSizeT  = 12 // 正文
	SVGFontSizeTS = 14 // 副标
	SVGFontSizeTH = 24 // 标题
	SVGStrokeW    = 1.5
)

// ExportSVG 从场景列表生成多个 SVG 文件 (v2 §12 规范)
func ExportSVG(scenes []SVGScene, outputDir string, width, height int) ([]string, error) {
	if width <= 0 {
		width = 1920
	}
	if height <= 0 {
		height = 1080
	}

	os.MkdirAll(outputDir, 0755)
	var paths []string

	moodBg := map[string]string{
		"mysterious": "#1a1a2e", "aspirational": "#2d1b69", "confident": "#8b0000",
		"urgent": "#1a0000", "calm": "#f0f4e8", "curious": "#f5f0e6",
		"excited": "#1a3a1a", "inspiring": "#fdf6e3", "informative": "#e8f0fe",
		"focused": "#fff8f0", "warm": "#fff3e0", "practical": "#f5f0e6",
		"frustrated": "#2a2a2a", "tense": "#000000", "triumphant": "#002200",
		"hopeful": "#f0fff0",
	}
	moodText := map[string]string{
		"mysterious": "#e0e0ff", "aspirational": "#e0d0ff", "confident": "#ffffff",
		"urgent": "#ffcccc", "calm": "#1a1a1a", "curious": "#1a1a1a",
		"excited": "#ffffff", "inspiring": "#1a1a1a", "informative": "#1a1a1a",
		"focused": "#1a1a1a", "warm": "#1a1a1a", "practical": "#1a1a1a",
		"frustrated": "#ffffff", "tense": "#ffffff", "triumphant": "#ffffff",
		"hopeful": "#1a1a1a",
	}
	moodAccent := map[string]string{
		"mysterious": "#4a4aff", "aspirational": "#b088ff", "confident": "#ff4444",
		"urgent": "#ff0044", "calm": "#88aa66", "curious": "#cc9966",
		"excited": "#66ff66", "inspiring": "#dd8844", "informative": "#4488cc",
		"focused": "#dd6633", "warm": "#ff8844", "practical": "#cc7744",
		"frustrated": "#666666", "tense": "#ff0000", "triumphant": "#ffdd00",
		"hopeful": "#88cc88",
	}

	for _, s := range scenes {
		bg := moodBg[s.Mood]
		if bg == "" {
			bg = "#f5f0e6"
		}
		tc := moodText[s.Mood]
		if tc == "" {
			tc = "#1a1a1a"
		}
		ac := moodAccent[s.Mood]
		if ac == "" {
			ac = "#C43C3A"
		}

		bodyLines := strings.Split(s.Body, "\n")
		var bodySVG strings.Builder
		// v2: 12px 正文, 间距 18 (在 viewBox=680 坐标系下)
		for i, line := range bodyLines {
			if line == "" {
				continue
			}
			y := 180 + i*18
			bodySVG.WriteString(fmt.Sprintf(
				`  <text class="t" x="40" y="%d" fill="%s">%s</text>`,
				y, tc, escapeXML(line)))
		}

		// v2: §12 坐标 (40 = 左 padding, 360 = 标题 y, 320 = 副标 y, 100 = 章节 y)
		bt := strings.ToUpper(s.BeatName)
		// 计算 voiceover 区域 (底部)
		voY := SVGViewBoxH - 30
		voBoxY := voY - 50
		// 计算 scene no (右下角, text-anchor=end)
		pageNoY := SVGViewBoxH - 10

		svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <defs>
    <style>
      .t  { font-family: system-ui, -apple-system, sans-serif; font-size: %dpx; font-weight: 400; fill: var(--text, #1a1a1a); }
      .ts { font-family: system-ui, -apple-system, sans-serif; font-size: %dpx; font-weight: 500; fill: var(--text, #1a1a1a); letter-spacing: 0.1em; }
      .th { font-family: system-ui, -apple-system, sans-serif; font-size: %dpx; font-weight: 700; fill: var(--text, #1a1a1a); }
      .accent { fill: var(--accent, #C43C3A); }
      .muted  { fill: var(--muted,  #888888); opacity: 0.7; }
    </style>
    <linearGradient id="bg" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%%" stop-color="%s"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.95"/>
    </linearGradient>
  </defs>

  <rect width="%d" height="%d" fill="url(#bg)" fill-opacity="1"/>

  <line x1="0" y1="0" x2="0" y2="%d" stroke="%s" stroke-width="%g" fill="none"/>
  <text class="ts muted" x="40" y="40">%s</text>
  <text class="th" x="40" y="80" fill="%s">%s</text>

  <line x1="40" y1="100" x2="%d" y2="100" stroke="%s" stroke-width="0.75" fill="none" opacity="0.3"/>

%s

  <rect x="40" y="%d" width="%d" height="40" rx="4" fill="none" stroke="%s" stroke-width="%g" opacity="0.4"/>
  <text class="ts" x="50" y="%d" fill="%s" font-style="italic">%s</text>

  <text class="t muted" x="%d" y="%d" text-anchor="end" fill="%s">%d / %d</text>
</svg>`,
			SVGViewBoxW, SVGViewBoxH, width, height,
			SVGFontSizeT, SVGFontSizeTS, SVGFontSizeTH,
			bg, bg,
			SVGViewBoxW, SVGViewBoxH, SVGViewBoxH, ac, SVGStrokeW,
			bt, tc, escapeXML(s.Title),
			SVGViewBoxW-40, ac,
			bodySVG.String(),
			voBoxY, SVGViewBoxW-80, ac, SVGStrokeW,
			voY-12, tc, escapeXML(s.Voiceover),
			SVGViewBoxW-40, pageNoY, tc, s.SceneNo, s.TotalScenes)

		filename := fmt.Sprintf("scene-%03d-%s.svg", s.SceneNo, slugify(s.BeatName))
		fp := filepath.Join(outputDir, filename)
		os.WriteFile(fp, []byte(svg), 0644)
		paths = append(paths, fp)
	}
	return paths, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "·", "")
	s = strings.ReplaceAll(s, "（", "")
	s = strings.ReplaceAll(s, "）", "")
	return s
}