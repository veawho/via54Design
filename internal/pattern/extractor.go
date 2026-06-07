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

package pattern

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Patterns 从HTML中提取的设计模式
type Patterns struct {
	Colors     ColorInfo   `json:"colors"`
	Fonts      FontInfo    `json:"fonts"`
	Layout     LayoutInfo  `json:"layout"`
	Animations AnimInfo    `json:"animations"`
	Metrics    MetricsInfo `json:"metrics"`
}

type ColorInfo struct {
	Palette      []ColorEntry       `json:"palette"`
	Roles        map[string]string  `json:"roles"`
	TotalUnique  int                `json:"total_unique"`
	ContrastInfo map[string]string  `json:"contrast_info"`
}

type ColorEntry struct {
	Hex  string `json:"hex"`
	Freq int    `json:"freq"`
}

type FontInfo struct {
	Families []string `json:"families"`
	Display  string   `json:"display"`
	Body     string   `json:"body"`
	Google   []string `json:"google_fonts"`
}

type LayoutInfo struct {
	Types       []string `json:"types"`
	Sections    int      `json:"sections"`
	CardCount   int      `json:"card_count"`
	Responsive  bool     `json:"responsive"`
}

type AnimInfo struct {
	HasMotion  bool     `json:"has_motion"`
	Types      []string `json:"types"`
	Complexity string   `json:"complexity"`
}

type MetricsInfo struct {
	TotalLines    int     `json:"total_lines"`
	CSSSizePct    float64 `json:"css_size_pct"`
	ImageCount    int     `json:"image_count"`
	SVGCount      int     `json:"svg_count"`
	ExternalLinks int     `json:"external_links"`
}

// Extractor 模式提取器
type Extractor struct {
	html string
}

func New(htmlContent string) *Extractor {
	return &Extractor{html: htmlContent}
}

func (e *Extractor) ExtractAll() *Patterns {
	return &Patterns{
		Colors:     e.extractColors(),
		Fonts:      e.extractFonts(),
		Layout:     e.extractLayout(),
		Animations: e.extractAnimations(),
		Metrics:    e.extractMetrics(),
	}
}

func (e *Extractor) extractColors() ColorInfo {
	// 提取所有 hex 颜色
	hexRe := regexp.MustCompile(`#[0-9A-Fa-f]{6}`)
	allHex := hexRe.FindAllString(e.html, -1)

	// 去重并统计频率
	freq := map[string]int{}
	for _, h := range allHex {
		freq[strings.ToUpper(h)]++
	}

	// 按频率排序
	type pair struct{ hex string; freq int }
	var sorted []pair
	for h, f := range freq {
		sorted = append(sorted, pair{h, f})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})

	// 截取前10
	top := sorted
	if len(top) > 10 {
		top = top[:10]
	}

	palette := []ColorEntry{}
	for _, p := range top {
		palette = append(palette, ColorEntry{Hex: p.hex, Freq: p.freq})
	}

	// 识别语义角色
	roles := map[string]string{}
	if len(sorted) > 0 {
		roles["background"] = sorted[0].hex
	}
	if len(sorted) > 1 {
		roles["text_primary"] = sorted[1].hex
	}
	// 找 accent (与背景对比度高)
	if len(sorted) > 2 {
		bg := sorted[0].hex
		for _, p := range sorted[2:] {
			if contrastRatio(bg, p.hex) > 3.0 {
				roles["accent"] = p.hex
				break
			}
		}
	}

	return ColorInfo{
		Palette:     palette,
		Roles:       roles,
		TotalUnique: len(sorted),
	}
}

func (e *Extractor) extractFonts() FontInfo {
	// font-family
	fontRe := regexp.MustCompile(`(?i)font-family\s*:\s*([^;}]+)`)
	fontMatches := fontRe.FindAllStringSubmatch(e.html, -1)

	// Google Fonts imports
	gfRe := regexp.MustCompile(`(?i)family=([^:&"]+)`)
	gfMatches := gfRe.FindAllStringSubmatch(e.html, -1)

	familySet := map[string]bool{}
	var families []string
	for _, m := range fontMatches {
		parts := strings.Split(m[1], ",")
		for _, p := range parts {
			f := strings.Trim(strings.TrimSpace(p), "'\"")
			if f != "" && !isGenericFont(f) && !familySet[f] {
				familySet[f] = true
				families = append(families, f)
			}
		}
	}

	var gf []string
	for _, m := range gfMatches {
		gf = append(gf, strings.ReplaceAll(m[1], "+", " "))
	}

	// 识别角色
	var display, body string
	for _, f := range families {
		if strings.Contains(strings.ToLower(f), "serif") {
			if display == "" { display = f }
		} else if strings.Contains(strings.ToLower(f), "sans") || strings.Contains(strings.ToLower(f), "inter") {
			if body == "" { body = f }
		}
	}
	// 回退
	if display == "" && len(families) > 0 { display = families[0] }
	if body == "" {
		if len(families) > 1 { body = families[1] }
	}

	return FontInfo{
		Families: families,
		Display:  display,
		Body:     body,
		Google:   gf,
	}
}

func (e *Extractor) extractLayout() LayoutInfo {
	var types []string

	if strings.Contains(e.html, "display: grid") || strings.Contains(e.html, "display:grid") {
		colsRe := regexp.MustCompile(`grid-template-columns\s*:\s*([^;}]+)`)
		if cols := colsRe.FindString(e.html); cols != "" {
			types = append(types, "grid")
		}
	}
	if strings.Contains(e.html, "display: flex") || strings.Contains(e.html, "display:flex") {
		types = append(types, "flexbox")
	}
	if strings.Contains(e.html, "column-count") || strings.Contains(e.html, "columns:") {
		types = append(types, "masonry/columns")
	}

	sections := strings.Count(e.html, "<section")
	cards := 0
	for _, kw := range []string{"card", "gallery-item", "bento"} {
		re := regexp.MustCompile(`class="[^"]*` + kw + `[^"]*"`)
		cards += len(re.FindAllString(e.html, -1))
	}

	return LayoutInfo{
		Types:      types,
		Sections:   sections,
		CardCount:  cards,
		Responsive: strings.Contains(e.html, "@media"),
	}
}

func (e *Extractor) extractAnimations() AnimInfo {
	var types []string

	reKf := regexp.MustCompile(`@keyframes\s+(\S+)`)
	if kfs := reKf.FindAllString(e.html, -1); len(kfs) > 0 {
		types = append(types, fmt.Sprintf("keyframes(%d)", len(kfs)))
	}

	reTr := regexp.MustCompile(`transition\s*:\s*([^;}]+)`)
	if trs := reTr.FindAllString(e.html, -1); len(trs) > 0 {
		types = append(types, fmt.Sprintf("transitions(%d)", len(trs)))
	}

	if strings.Contains(e.html, "@keyframes") {
		types = append(types, "css-animations")
	}
	if strings.Contains(e.html, "transform") {
		types = append(types, "transforms")
	}

	complexity := "low"
	if len(types) > 5 { complexity = "high" } else if len(types) > 2 { complexity = "medium" }

	return AnimInfo{
		HasMotion:  len(types) > 0,
		Types:      types,
		Complexity: complexity,
	}
}

func (e *Extractor) extractMetrics() MetricsInfo {
	lines := strings.Count(e.html, "\n") + 1

	styleRe := regexp.MustCompile(`(?i)<style[^>]*>(.*?)</style>`)
	cssLen := 0
	for _, m := range styleRe.FindAllStringSubmatch(e.html, -1) {
		cssLen += len(m[1])
	}
	cssPct := 0.0
	if len(e.html) > 0 {
		cssPct = float64(cssLen) / float64(len(e.html)) * 100
	}

	imgs := len(regexp.MustCompile(`(?i)<img`).FindAllString(e.html, -1))
	svgs := len(regexp.MustCompile(`<svg`).FindAllString(e.html, -1))
	exts := len(regexp.MustCompile(`src="https?://`).FindAllString(e.html, -1)) +
		len(regexp.MustCompile(`href="https?://`).FindAllString(e.html, -1))

	return MetricsInfo{
		TotalLines:    lines,
		CSSSizePct:    cssPct,
		ImageCount:    imgs,
		SVGCount:      svgs,
		ExternalLinks: exts,
	}
}

// ToYAML 生成YAML模板候选
func (p *Patterns) ToYAML(name string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Pattern extracted from: %s\n", name))
	b.WriteString("# Auto-generated by via54Design pattern extractor\n\n")

	b.WriteString("colors:\n")
	b.WriteString("  roles:\n")
	for role, hex := range p.Colors.Roles {
		b.WriteString(fmt.Sprintf("    %s: \"%s\"\n", role, hex))
	}
	b.WriteString("  palette:\n")
	for _, c := range p.Colors.Palette {
		b.WriteString(fmt.Sprintf("    - hex: \"%s\"  # freq=%d\n", c.Hex, c.Freq))
	}

	b.WriteString("\nfonts:\n")
	b.WriteString(fmt.Sprintf("  display: \"%s\"\n", p.Fonts.Display))
	b.WriteString(fmt.Sprintf("  body: \"%s\"\n", p.Fonts.Body))
	if len(p.Fonts.Google) > 0 {
		b.WriteString(fmt.Sprintf("  google: %v\n", p.Fonts.Google))
	}

	b.WriteString("\nlayout:\n")
	b.WriteString(fmt.Sprintf("  types: %v\n", p.Layout.Types))
	b.WriteString(fmt.Sprintf("  sections: %d\n", p.Layout.Sections))
	b.WriteString(fmt.Sprintf("  responsive: %v\n", p.Layout.Responsive))

	b.WriteString("\nmetrics:\n")
	b.WriteString(fmt.Sprintf("  lines: %d\n", p.Metrics.TotalLines))
	b.WriteString(fmt.Sprintf("  images: %d\n", p.Metrics.ImageCount))
	b.WriteString(fmt.Sprintf("  svgs: %d\n", p.Metrics.SVGCount))

	return b.String()
}

// ExtractFromHTML 便捷函数
func ExtractFromHTML(htmlContent, name string) (*Patterns, string) {
	p := New(htmlContent).ExtractAll()
	yaml := p.ToYAML(name)
	return p, yaml
}

// ─── helpers ───

func isGenericFont(f string) bool {
	lower := strings.ToLower(f)
	generics := []string{"serif", "sans-serif", "monospace", "cursive", "fantasy",
		"system-ui", "ui-serif", "ui-sans-serif", "ui-monospace", "ui-rounded",
		"-apple-system", "blinkmacsystemfont", "segoe ui", "helvetica neue",
		"arial", "times new roman", "georgia"}
	for _, g := range generics {
		if lower == g {
			return true
		}
	}
	return false
}

func hexToRGB(hex string) (int, int, int) {
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 3 {
		h = string(h[0]) + string(h[0]) + string(h[1]) + string(h[1]) + string(h[2]) + string(h[2])
	}
	if len(h) != 6 { return 0, 0, 0 }
	r := int(hexByte(h[0:2]))
	g := int(hexByte(h[2:4]))
	b := int(hexByte(h[4:6]))
	return r, g, b
}

func hexByte(s string) uint8 {
	var v uint8
	fmt.Sscanf(s, "%x", &v)
	return v
}

func relativeLuminance(r, g, b int) float64 {
	// sRGB linearization
	linear := func(c float64) float64 {
		c = c / 255.0
		if c <= 0.03928 { return c / 12.92 }
		return ((c + 0.055) / 1.055) * ((c + 0.055) / 1.055) * ((c + 0.055) / 1.055)
	}
	return 0.2126*linear(float64(r)) + 0.7152*linear(float64(g)) + 0.0722*linear(float64(b))
}

func contrastRatio(c1, c2 string) float64 {
	r1, g1, b1 := hexToRGB(c1)
	r2, g2, b2 := hexToRGB(c2)
	l1 := relativeLuminance(r1, g1, b1)
	l2 := relativeLuminance(r2, g2, b2)
	if l1 < l2 { l1, l2 = l2, l1 }
	return (l1 + 0.05) / (l2 + 0.05)
}
