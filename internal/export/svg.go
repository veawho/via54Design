// via54Design — SVG 矢量导出
// 每个叙事 scene → 独立 SVG 文件 (16:9 矢量，可无限缩放)
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

// ExportSVG 从场景列表生成多个 SVG 文件
func ExportSVG(scenes []SVGScene, outputDir string, width, height int) ([]string, error) {
	if width <= 0 { width = 1920 }
	if height <= 0 { height = 1080 }

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
		if bg == "" { bg = "#f5f0e6" }
		tc := moodText[s.Mood]
		if tc == "" { tc = "#1a1a1a" }
		ac := moodAccent[s.Mood]
		if ac == "" { ac = "#C43C3A" }

		bodyLines := strings.Split(s.Body, "\n")
		var bodySVG strings.Builder
		for i, line := range bodyLines {
			if line == "" { continue }
			y := 120 + i*40
			bodySVG.WriteString(fmt.Sprintf(
				`    <text x="100" y="%d" font-family="system-ui" font-size="24" fill="%s">%s</text>`,
				y, tc, escapeXML(line)))
		}

		bt := strings.ToUpper(s.BeatName)
		svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%%" stop-color="%s"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.95"/>
    </linearGradient>
  </defs>
  <rect width="%d" height="%d" fill="url(#bg)"/>
  <rect x="0" y="0" width="12" height="%d" fill="%s"/>
  <text x="60" y="80" font-family="system-ui" font-size="18" fill="%s" opacity="0.6" letter-spacing="4">%s</text>
  <text x="60" y="160" font-family="system-ui" font-size="48" font-weight="700" fill="%s">%s</text>
  %s
  <rect x="60" y="%d" width="%d" height="80" rx="8" fill="%s" opacity="0.1"/>
  <text x="80" y="%d" font-family="system-ui" font-size="28" font-style="italic" fill="%s">%s</text>
  <text x="%d" y="%d" font-family="system-ui" font-size="16" fill="%s" opacity="0.4" text-anchor="end">%d / %d</text>
</svg>`,
			width, height, width, height, bg, bg,
			width, height, height, ac,
			tc, bt,
			tc, escapeXML(s.Title),
			bodySVG.String(),
			height-140, width-120, ac,
			height-90, tc, escapeXML(s.Voiceover),
			width-60, height-30, tc, s.SceneNo, s.TotalScenes)

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
