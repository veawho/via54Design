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
	"github.com/veawho/via54Design/internal/template"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func cmdGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	layout := fs.String("layout", "", "布局模板ID")
	color := fs.String("color", "", "配色模板ID")
	font := fs.String("font", "", "字体模板ID")
	title := fs.String("title", "via54Design", "页面标题")
	output := fs.String("output", "output.html", "输出路径")
	letteringSVG := fs.String("lettering-svg", "", "SVG lettering 文件路径 (手写/书法文字)")
	fromNarrative := fs.String("from-narrative", "", "叙事脚手架 JSON 路径 (via54 narrate --format json 的输出)")
	presentation := fs.Bool("presentation", false, "演示模式: 锁定 16:9 (PPT/视频输出)")
	fs.Parse(os.Args[2:])

	bd := baseDir()
	eng, err := template.NewEngine(bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "失败: %v\n", err)
		os.Exit(1)
	}

	// ── 叙事驱动生成 ──
	if *fromNarrative != "" {
		generateFromNarrative(eng, *fromNarrative, *layout, *color, *font, *output, *presentation)
		return
	}

	if (*layout == "" || *color == "" || *font == "") && *letteringSVG == "" {
		fmt.Fprintln(os.Stderr, "请指定: --layout, --color, --font")
		fmt.Fprintln(os.Stderr, "或者: --from-narrative scaffold.json (叙事驱动模式)")
		os.Exit(1)
	}
	var result *template.GenerationResult
	if *letteringSVG != "" {
		data, _ := os.ReadFile(*letteringSVG)
		if *layout == "" {
			*layout = "hero-split-left-image"
		}
		if *color == "" {
			*color = "ink-wash"
		}
		if *font == "" {
			*font = "serif-sans-editorial"
		}
		result, err = eng.ComposeWithSVG(*layout, *color, *font, *title, string(data), *presentation)
	} else {
		result, err = eng.ComposeWithSVG(*layout, *color, *font, *title, "", *presentation)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
	result.SaveToFile(*output)
	fmt.Printf("✅ %s (%d bytes)\n   layout=%s color=%s font=%s\n", *output, len(result.HTML), result.LayoutID, result.ColorID, result.FontID)
}

// generateFromNarrative 从叙事脚手架 JSON 生成多场景 HTML 动画
func generateFromNarrative(eng *template.Engine, narrativePath, layoutOverride, colorOverride, fontOverride, output string, presentationMode bool) {
	data, err := os.ReadFile(narrativePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取叙事文件失败: %v\n", err)
		os.Exit(1)
	}

	var scaffold narrate.NarrativeScaffold
	if err := yaml.Unmarshal(data, &scaffold); err != nil {
		fmt.Fprintf(os.Stderr, "解析叙事文件失败: %v\n", err)
		os.Exit(1)
	}

	// 叙事映射到视觉模板
	moodColor := map[string]string{
		"mysterious":   "dark-terminal-blue",
		"aspirational": "cosmic-retro",
		"confident":    "neo-brutalist-vibrant",
		"urgent":       "neon-dark",
		"calm":         "moon-white",
		"curious":      "ink-wash",
		"excited":      "candy-duolingo",
		"inspiring":    "warm-editorial-cream",
		"informative":  "bauhaus-primary",
		"focused":      "crimson-elegance",
		"warm":         "daylily-warmth",
		"practical":    "earth-terracotta",
		"insightful":   "ultramarine-deep",
		"hopeful":      "pine-spring",
		"frustrated":   "dark-terminal-blue",
		"tense":        "neon-dark",
		"triumphant":   "candy-duolingo",
	}

	// 默认视觉模板
	lID := layoutOverride
	if lID == "" {
		lID = "hero-split-left-image"
	}
	cID := colorOverride
	if cID == "" {
		cID = "ink-wash"
	}
	fID := fontOverride
	if fID == "" {
		fID = "ming-hei-editorial"
	}

	// 为每个 beat 生成场景 HTML，拼合成多场景页面
	var scenesHTML strings.Builder
	scenesHTML.WriteString(`<!DOCTYPE html><html lang="zh-CN"><head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: system-ui; overflow-x: hidden; }
.scene { min-height: 100vh; display: flex; flex-direction: column; justify-content: center;
  align-items: center; padding: 4rem 2rem; position: relative; transition: opacity 0.5s; }
.scene .time { position: absolute; top: 1rem; right: 1rem; font-size: 0.8rem; opacity: 0.5; }
.scene .mood-tag { position: absolute; top: 1rem; left: 1rem; font-size: 0.7rem;
  text-transform: uppercase; letter-spacing: 0.1em; padding: 0.2rem 0.6rem;
  border-radius: 4px; background: rgba(255,255,255,0.1); }
.scene h2 { font-size: 1.8rem; margin-bottom: 1rem; text-align: center; max-width: 40rem; }
.scene .voiceover { font-size: 2.4rem; font-weight: 700; text-align: center;
  margin: 1.5rem 0; max-width: 36rem; line-height: 1.4; }
.scene .sub { font-size: 1rem; opacity: 0.7; text-align: center;
  margin-top: 1rem; max-width: 30rem; }
.scene .act-label { font-size: 0.9rem; opacity: 0.6; margin-bottom: 0.5rem;
  text-transform: uppercase; letter-spacing: 0.15em; }
.scene-nav { position: fixed; bottom: 2rem; left: 50%; transform: translateX(-50%);
  display: flex; gap: 0.5rem; z-index: 100; }
.scene-nav button { padding: 0.5rem 1.2rem; border: 1px solid currentColor;
  background: transparent; cursor: pointer; border-radius: 4px; font-size: 0.9rem;
  transition: all 0.2s; }
.scene-nav button:hover { background: rgba(255,255,255,0.15); }
.scene-nav .progress { display: flex; align-items: center; gap: 0.3rem;
  font-size: 0.8rem; opacity: 0.5; }
</style></head><body>
`)

	for i, beat := range scaffold.Beats {
		// 根据情绪选配色
		sceneColor := cID
		if mc, ok := moodColor[beat.Mood]; ok && colorOverride == "" {
			sceneColor = mc
		}

		// 每个场景用独立配色生成
		sceneResult, err := eng.ComposeWithSVG(lID, sceneColor, fID,
			fmt.Sprintf("%s — %s", scaffold.Seed, beat.Act), "", presentationMode)
		if err != nil {
			// fallback: plain scene
			sceneResult = &template.GenerationResult{
				HTML: fmt.Sprintf("<div class=\"scene\"><h2>%s</h2><p class=\"voiceover\">%s</p></div>",
					beat.Act, beat.Voiceover),
			}
		}

		// 提取 body 内容
		html := sceneResult.HTML
		bodyStart := strings.Index(html, "<body")
		bodyEnd := strings.Index(html, "</body>")
		bodyContent := ""
		if bodyStart > 0 && bodyEnd > 0 {
			contentStart := strings.Index(html[bodyStart:], ">")
			if contentStart > 0 {
				bodyContent = html[bodyStart+contentStart+1 : bodyEnd]
			}
		}
		if bodyContent == "" {
			bodyContent = fmt.Sprintf("<h2>%s</h2><p>%s</p>", beat.Act, beat.Voiceover)
		}

		scenesHTML.WriteString(fmt.Sprintf(`<div class="scene" id="scene-%d" data-beat="%s" data-mood="%s" data-duration="%d">
  <div class="act-label">%s</div>
  <div class="time">%ds → %ds</div>
  <div class="mood-tag">%s</div>
  %s
  <div class="voiceover">%s</div>
</div>
`, i+1, beat.Act, beat.Mood, beat.Duration, beat.Act,
			beat.StartTime, beat.StartTime+beat.Duration, beat.Mood,
			bodyContent, beat.Voiceover))
	}

	// 导航 + JS 自动播放
	scenesHTML.WriteString(`<div class="scene-nav">`)
	for i := range scaffold.Beats {
		scenesHTML.WriteString(fmt.Sprintf(`<button onclick="location.href='#scene-%d'" data-idx="%d">%d</button>`, i+1, i, i+1))
	}
	scenesHTML.WriteString(`<div class="progress"><span id="progress-text">1/`)
	scenesHTML.WriteString(fmt.Sprintf("%d", len(scaffold.Beats)))
	scenesHTML.WriteString(`</span></div></div>
<script>
let currentScene = 0;
const scenes = document.querySelectorAll('.scene');
scenes.forEach((s, i) => { s.style.display = i === 0 ? 'flex' : 'none'; });
document.querySelectorAll('.scene-nav button').forEach(btn => {
  btn.addEventListener('click', () => {
    const idx = parseInt(btn.dataset.idx);
    scenes.forEach((s, i) => { s.style.display = i === idx ? 'flex' : 'none'; });
    currentScene = idx;
  });
});
</script></body></html>`)

	finalHTML := scenesHTML.String()
	os.WriteFile(output, []byte(finalHTML), 0644)
	fmt.Printf("✅ 叙事驱动动画: %s (%d 场景, %d 字节)\n   via54 narrate → via54 generate --from-narrative\n",
		output, len(scaffold.Beats), len(finalHTML))
}
