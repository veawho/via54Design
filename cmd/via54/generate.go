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
	"errors"
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
	stdout := fs.Bool("stdout", false, "输出到 stdout (覆盖 --output)")
	letteringSVG := fs.String("lettering-svg", "", "SVG lettering 文件路径 (手写/书法文字)")
	fromNarrative := fs.String("from-narrative", "", "叙事脚手架 JSON 路径 (via54 narrate --format json 的输出)")
	presentation := fs.Bool("presentation", false, "演示模式: 锁定 16:9 (PPT/视频输出)")
	strict := fs.Bool("strict", false, "严格模式: layout/color/font ID 必须存在, 否则报错 (默认宽松)")
	fs.Parse(os.Args[2:])

	bd := baseDir()
	eng, err := template.NewEngine(bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "失败: %v\n", err)
		os.Exit(1)
	}

	// ── 严格模式 ID 校验 ──
	if *strict {
		if err := validateIDs(eng, *layout, *color, *font); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			fmt.Fprintln(os.Stderr, "提示: 用 'via54 list' 查看可用 ID; 用 --strict=false (默认) 跳过校验")
			os.Exit(1)
		}
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

	// ── --stdout 模式: 输出到 stdout 而非文件 ──
	if *stdout {
		fmt.Print(result.HTML)
		return
	}

	result.SaveToFile(*output)
	fmt.Printf("✅ %s (%d bytes)\n   layout=%s color=%s font=%s\n", *output, len(result.HTML), result.LayoutID, result.ColorID, result.FontID)
}

// validateIDs 严格模式下校验 layout/color/font ID 是否存在
func validateIDs(eng *template.Engine, layoutID, colorID, fontID string) error {
	if layoutID == "" || colorID == "" || fontID == "" {
		return fmt.Errorf("严格模式: layout/color/font 都不能为空")
	}
	all := eng.Registry.ListAll()
	ids := map[string]map[string]bool{
		"layouts":       make(map[string]bool),
		"color_schemes": make(map[string]bool),
		"typography":    make(map[string]bool),
	}
	for cat, entries := range all {
		if _, ok := ids[cat]; !ok {
			ids[cat] = make(map[string]bool)
		}
		for _, e := range entries {
			ids[cat][e.ID] = true
		}
	}
	checks := []struct {
		id    string
		category string
		label string
	}{
		{layoutID, "layouts", "layout"},
		{colorID, "color_schemes", "color"},
		{fontID, "typography", "font"},
	}
	for _, c := range checks {
		if !ids[c.category][c.id] {
			// 提示相似 ID
			suggestions := findSimilar(c.id, ids[c.category])
			hint := ""
			if len(suggestions) > 0 {
				hint = fmt.Sprintf(" (相近: %s)", strings.Join(suggestions, ", "))
			}
			return fmt.Errorf("%s ID %q 不存在%s", c.label, c.id, hint)
		}
	}
	return nil
}

// findSimilar 综合相似度: Levenshtein 距离 + 前缀匹配
// 用户 typo 模式: "heroo-split" → "hero-split-16-9" (1 字符替换 + 后缀差异)
//                "bento-2x2"   → "bento-grid-2x2"  (中段不同)
func findSimilar(target string, pool map[string]bool) []string {
	targetLower := strings.ToLower(target)
	threshold := len(targetLower) / 3
	if threshold < 2 {
		threshold = 2
	}
	type scored struct {
		id    string
		score int
	}
	var all []scored
	for id := range pool {
		if id == target {
			continue
		}
		idLower := strings.ToLower(id)
		d := levenshtein(targetLower, idLower)
		// 前缀匹配加分(距离相同但前缀更长优先)
		prefixBonus := 0
		for i := 0; i < minBytes(len(targetLower), len(idLower)); i++ {
			if targetLower[i] == idLower[i] {
				prefixBonus++
			} else {
				break
			}
		}
		// 综合得分: 距离 - 前缀奖励
		score := d - prefixBonus/2
		if d <= threshold || (prefixBonus >= 4 && d <= threshold+3) {
			all = append(all, scored{id, score})
		}
	}
	// 按得分排序(越小越相似)
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].score < all[i].score {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	var suggestions []string
	for _, s := range all {
		suggestions = append(suggestions, s.id)
		if len(suggestions) >= 5 {
			break
		}
	}
	return suggestions
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
		dp[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		dp[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[len(a)][len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
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
		// errors.As: 提取 yaml 详细错误位置 (line/column)
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			fmt.Fprintf(os.Stderr, "解析叙事文件失败 (yaml 类型错误): %v\n", yamlErr)
		} else {
			fmt.Fprintf(os.Stderr, "解析叙事文件失败: %v\n", err)
		}
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
		} else {
			// Contextualize static placeholders with beat attributes
			bodyContent = strings.ReplaceAll(bodyContent, "EYEBROW", beat.Act)
			bodyContent = strings.ReplaceAll(bodyContent, "副标题内容", beat.Voiceover)
			bodyContent = strings.ReplaceAll(bodyContent, "CTA 按钮", "查看详情")
			bodyContent = strings.ReplaceAll(bodyContent, "Aesthetic, modular, responsive, and performance-tuned layouts powered by Golang & modern design agents.", beat.Voiceover)
			bodyContent = strings.ReplaceAll(bodyContent, "月活用户", "分镜场景")
			bodyContent = strings.ReplaceAll(bodyContent, "ARR", "叙事节奏")
			bodyContent = strings.ReplaceAll(bodyContent, "客户留存率", "情感指数")
			bodyContent = strings.ReplaceAll(bodyContent, "活跃项目", "镜头时长")
			bodyContent = strings.ReplaceAll(bodyContent, "Console Overview", "场景视觉概述")
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
