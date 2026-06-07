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

package template

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"gopkg.in/yaml.v3"
)

type Engine struct {
	Registry *Registry
	BaseDir  string
}

func NewEngine(baseDir string) (*Engine, error) {
	reg, err := NewRegistry(baseDir)
	if err != nil {
		return nil, err
	}
	return &Engine{Registry: reg, BaseDir: baseDir}, nil
}

func (e *Engine) Compose(layoutID, colorID, fontID, title string) (*GenerationResult, error) {
	return e.ComposeWithSVG(layoutID, colorID, fontID, title, "", false)
}

func (e *Engine) ComposeWithSVG(layoutID, colorID, fontID, title, letteringSVG string, presentationMode bool) (*GenerationResult, error) {
	lp, _ := e.Registry.ResolveLayout(layoutID)
	cp, _ := e.Registry.ResolveColorScheme(colorID)
	fp, _ := e.Registry.ResolveTypography(fontID)

	layout, err := loadYAML[LayoutTemplate](lp)
	if err != nil {
		return nil, fmt.Errorf("layout: %w", err)
	}
	color, err := loadYAML[ColorSchemeTemplate](cp)
	if err != nil {
		return nil, fmt.Errorf("color: %w", err)
	}
	font, err := loadYAML[TypographyTemplate](fp)
	if err != nil {
		return nil, fmt.Errorf("font: %w", err)
	}

	result := &GenerationResult{
		LayoutID: layoutID,
		ColorID:  colorID,
		FontID:   fontID,
		Title:    title,
		LetteringSVG: letteringSVG,
		PresentationMode: presentationMode,
	}
	result.CSSVars = buildCSSVariables(color, font)
	result.FontImports = buildFontImports(font)
	result.BaseCSS = buildBaseCSS(font, layout, presentationMode)
	result.LayoutCSS = buildLayoutCSS(layout, presentationMode)
	result.HTML = assembleHTML(result, layout)
	return result, nil
}

func loadYAML[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v T
	return &v, yaml.Unmarshal(data, &v)
}

func buildCSSVariables(color *ColorSchemeTemplate, font *TypographyTemplate) string {
	var b strings.Builder
	b.WriteString(":root {\n")
	if color.CSSVariables != "" {
		b.WriteString(color.CSSVariables)
	} else if len(color.Palette) > 0 {
		for _, item := range color.Palette {
			b.WriteString(fmt.Sprintf("  --%s: %s;\n", item.Role, item.Hex))
		}
	} else if color.Colors != nil {
		keys := sortedKeys(color.Colors)
		for _, role := range keys {
			b.WriteString(fmt.Sprintf("  --%s: %s;\n", role, color.Colors[role]))
		}
	}
	keys := sortedKeys(font.Sizes)
	for _, name := range keys {
		b.WriteString(fmt.Sprintf("  --size-%s: %s;\n", name, font.Sizes[name]))
	}
	b.WriteString("}")
	return b.String()
}

// buildLayoutCSS 合并手写 CSS + 自动生成响应式 + 间距变量
func buildLayoutCSS(layout *LayoutTemplate, presentationMode bool) string {
	var parts []string

	// 1. 手写 CSS (核心布局样式)
	if layout.CSS != "" {
		parts = append(parts, layout.CSS)
	}

	// 2. 间距变量注入 (黄金比例)
	parts = append(parts, buildSpacingCSS(layout.Spacing))

	// 3. 断点自动编译
	parts = append(parts, buildResponsiveCSS(layout))

	// 4. 元素级响应式
	parts = append(parts, buildElementResponsiveCSS(layout))

	return strings.Join(parts, "\n\n")
}

// buildSpacingCSS 从 YAML spacing 注入黄金比例 CSS 变量
// 参考: Extra-Strength Responsive Grids 流体间距体系
func buildSpacingCSS(spacing SpacingScale) string {
	if spacing.Base <= 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("/* 间距系统 (黄金比例 φ=1.618) */\n:root {\n")
	step := float64(spacing.Base)
	for i := 1; i <= 12; i++ {
		px := int(step + 0.5)
		b.WriteString(fmt.Sprintf("  --space-step-%d: %dpx;\n", i, px))
		step *= spacing.Ratio
	}
	keys := sortedKeys(spacing.Semantic)
	for _, name := range keys {
		b.WriteString(fmt.Sprintf("  --space-%s: var(--%s);\n", name, spacing.Semantic[name]))
	}
	b.WriteString("}")
	return b.String()
}

// buildResponsiveCSS 从 YAML responsive[] 自动编译媒体查询
// 覆盖: columns / safe_area / font_scale / spacing_scale
func buildResponsiveCSS(layout *LayoutTemplate) string {
	if len(layout.Responsive) == 0 {
		return ""
	}

	className := layoutClassName(layout.ID)
	var b strings.Builder
	b.WriteString("/* 自动编译响应式断点 */\n")

	for _, bp := range layout.Responsive {
		if bp.MinWidth == 0 && bp.MaxWidth == 0 {
			continue
		}

		// ── 媒体查询条件 ──
		if bp.MaxWidth > 0 {
			b.WriteString(fmt.Sprintf("@media (min-width: %dpx) and (max-width: %dpx) {\n", bp.MinWidth, bp.MaxWidth))
		} else {
			b.WriteString(fmt.Sprintf("@media (min-width: %dpx) {\n", bp.MinWidth))
		}

		// ── 栅格布局 ──
		if bp.Columns != "" {
			b.WriteString(fmt.Sprintf("  .%s { grid-template-columns: %s; }\n", className, bp.Columns))
		}

		// ── Stack (堆叠 + 控制顺序) ──
		if bp.Stack {
			b.WriteString(fmt.Sprintf("  .%s { grid-template-columns: 1fr; }\n", className))
			for i, role := range bp.StackOrder {
				elClass := fmt.Sprintf("%s__%s", className, elementCSSRole(role))
				b.WriteString(fmt.Sprintf("  .%s { order: %d; }\n", elClass, i+1))
			}
		}

		// ── 安全区域 ──
		if len(bp.SafeArea) == 4 {
			// 作用于 text-container
			b.WriteString(fmt.Sprintf("  .%s__text { padding: %dpx %dpx %dpx %dpx; }\n",
				className, bp.SafeArea[0], bp.SafeArea[1], bp.SafeArea[2], bp.SafeArea[3]))
		}

		// ── 字体缩放 ──
		if bp.FontScale > 0 && bp.FontScale != 1.0 {
			b.WriteString(fmt.Sprintf("  .%s { --bp-font-scale: %.2f; }\n", className, bp.FontScale))
			b.WriteString(fmt.Sprintf("  .%s h1, .%s h2, .%s p { font-size: calc(1em * %.2f); }\n",
				className, className, className, bp.FontScale))
		}

		// ── 隐藏元素 ──
		for _, role := range bp.HideRoles {
			elClass := fmt.Sprintf("%s__%s", className, elementCSSRole(role))
			b.WriteString(fmt.Sprintf("  .%s { display: none; }\n", elClass))
		}

		b.WriteString("}\n")
	}
	return b.String()
}

// buildElementResponsiveCSS 从 Element.Responsive 编译元素级响应式
func buildElementResponsiveCSS(layout *LayoutTemplate) string {
	if len(layout.Responsive) == 0 {
		return ""
	}
	className := layoutClassName(layout.ID)
	var b strings.Builder
	b.WriteString("/* 元素级响应式 */\n")

	// 收集所有元素到拍平列表
	var elements []Element
	var walk func(elems []Element)
	walk = func(elems []Element) {
		for _, e := range elems {
			elements = append(elements, e)
			if len(e.Children) > 0 {
				walk(e.Children)
			}
		}
	}
	walk(layout.Elements)

	for _, el := range elements {
		if len(el.Responsive) == 0 {
			continue
		}
		elClass := fmt.Sprintf("%s__%s", className, elementCSSRole(el.Role))

		for bpName, res := range el.Responsive {
			// 找对应断点的尺寸
			bp := findBreakpoint(layout.Responsive, bpName)
			if bp == nil {
				continue
			}

			// 媒体查询
			if bp.MaxWidth > 0 {
				b.WriteString(fmt.Sprintf("@media (min-width: %dpx) and (max-width: %dpx) {\n", bp.MinWidth, bp.MaxWidth))
			} else {
				b.WriteString(fmt.Sprintf("@media (min-width: %dpx) {\n", bp.MinWidth))
			}

			if res.Hide {
				b.WriteString(fmt.Sprintf("  .%s { display: none; }\n", elClass))
			}
			if res.Order != 0 {
				b.WriteString(fmt.Sprintf("  .%s { order: %d; }\n", elClass, res.Order))
			}
			if res.FontSize != "" {
				b.WriteString(fmt.Sprintf("  .%s { font-size: %s; }\n", elClass, res.FontSize))
			}
			if len(res.Padding) == 4 {
				b.WriteString(fmt.Sprintf("  .%s { padding: %dpx %dpx %dpx %dpx; }\n",
					elClass, res.Padding[0], res.Padding[1], res.Padding[2], res.Padding[3]))
			}

			b.WriteString("}\n")
		}
	}
	return b.String()
}

// layoutClassName 从布局 ID 生成 CSS 类名
func layoutClassName(id string) string {
	// hero-split-16-9 → layout-hero-split
	// bento-grid-2x2 → layout-bento
	// gallery-waterfall → layout-gallery
	switch {
	case len(id) >= 11 && id[:11] == "hero-split-":
		return "layout-hero-split"
	case len(id) >= 5 && id[:5] == "bento":
		return "layout-bento"
	case len(id) >= 7 && id[:7] == "gallery":
		return "layout-gallery"
	}
	return "layout-" + id
}

// elementCSSRole 从 role 生成 CSS 类名片段
func elementCSSRole(role string) string {
	// image-container → image
	// text-container → text
	// card-icon → card-icon
	// gallery-item → item
	switch role {
	case "image-container": return "image"
	case "text-container":  return "text"
	case "gallery-item":    return "item"
	}
	return role
}

// findBreakpoint 按名称查找断点
func findBreakpoint(bps []BreakpointDef, name string) *BreakpointDef {
	for _, bp := range bps {
		if bp.Name == name {
			return &bp
		}
	}
	return nil
}

func buildFontImports(font *TypographyTemplate) string {
	if len(font.GoogleFonts) > 0 {
		var b strings.Builder
		b.WriteString(`<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
`)
		for _, gf := range font.GoogleFonts {
			b.WriteString(fmt.Sprintf(
				`<link href="https://fonts.googleapis.com/css2?family=%s&display=swap" rel="stylesheet">`+"\n", gf))
		}
		return b.String()
	}

	gf := map[string]bool{
		"Inter":true,"Geist":true,"JetBrains Mono":true,"Fraunces":true,
		"Playfair Display":true,"Noto Serif SC":true,"Noto Sans SC":true,
		"EB Garamond":true,"Nunito":true,"Baloo 2":true,
		"Archivo Black":true,"Archivo":true,"LXGW WenKai":true,"ZCOOL XiaoWei":true,
	}
	for _, family := range font.Fonts {
		primary := strings.Trim(strings.Split(family, ",")[0], " '\"")
		if gf[primary] {
			return fmt.Sprintf(
				`<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=%s:wght@400;500;700&display=swap" rel="stylesheet">`,
				strings.ReplaceAll(primary, " ", "+"))
		}
	}
	return ""
}

func buildBaseCSS(font *TypographyTemplate, layout *LayoutTemplate, presentationMode bool) string {
	body := getOrDefault(font.Fonts, "body", "'Inter', sans-serif")
	disp := getOrDefault(font.Fonts, "display", body)
	mono := getOrDefault(font.Fonts, "mono", "'JetBrains Mono', monospace")

	// Viewport lock — 仅在演示模式激活 16:9
	vpLock := ""
	if presentationMode && layout.Viewport.Baseline != "" {
		vpLock = fmt.Sprintf(`/* 演示模式: 16:9 锁定 */
.presentation-mode {
  aspect-ratio: %s;
  margin: 0 auto;
  max-height: 100dvh;
  overflow: hidden;
  background: var(--presentation-bg, #000);
  display: flex;
  align-items: center;
  justify-content: center;
}
.presentation-mode > .layout-hero-split,
.presentation-mode > .layout-bento,
.presentation-mode > .layout-gallery {
  aspect-ratio: %s;
  max-height: 100dvh;
}
@media (max-aspect-ratio: %s) {
  .presentation-mode > .layout-hero-split,
  .presentation-mode > .layout-bento,
  .presentation-mode > .layout-gallery {
    width: auto;
    height: 100dvh;
  }
}
`, layout.Viewport.Baseline, layout.Viewport.Baseline, layout.Viewport.Baseline)
	}

	return fmt.Sprintf(`* { box-sizing:border-box; margin:0; padding:0; }
html { scroll-behavior:smooth; }
%s
body {
  font-family: %s; line-height:1.7;
  color: var(--text-primary, #1A1A1A);
  background: var(--background, #FFFFFF);
}
h1,h2,h3,h4 { font-family: %s; line-height:1.1; }
code,pre { font-family: %s; }
a { color: var(--accent, inherit); text-decoration:none; }
img { max-width: 100%%; height:auto; display:block; }
:root { --layout-max-width: %s; }`,
		vpLock, body, disp, mono,
		getLayoutMaxWidth(layout))
}

func getLayoutMaxWidth(layout *LayoutTemplate) string {
	if layout.Viewport.MaxWidth != "" {
		return layout.Viewport.MaxWidth
	}
	return "1920px"
}

func getOrDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return def
}

func assembleHTML(r *GenerationResult, layout *LayoutTemplate) string {
	titleHTML := r.Title
	if r.LetteringSVG != "" {
		titleHTML = r.LetteringSVG
	}

	// Viewport meta with device adaptation
	vpMeta := `<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">`

	// 布局特有 HTML body 内容 (根据 layout ID 选择)
	bodyContent := heroBodyHTML(titleHTML)
	if layout.ID == "bento-grid-2x2" {
		bodyContent = bentoBodyHTML()
	} else if layout.ID == "gallery-waterfall" {
		bodyContent = galleryBodyHTML()
	}

	// 演示模式：外层包裹 .presentation-mode 容器
	bodyWrapper := ""
	if r.PresentationMode {
		bodyWrapper = `<div class="presentation-mode">\n  ` + bodyContent + `\n</div>`
		bodyContent = bodyWrapper
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
%s
<title>%s</title>
%s
<style>
%s
%s
%s
</style>
</head>
<body>
<main>%s</main>
</body>
</html>`, vpMeta, r.Title, r.FontImports, r.CSSVars, r.BaseCSS, r.LayoutCSS, bodyContent)
}

func heroBodyHTML(titleHTML string) string {
	return fmt.Sprintf(`
<section class="layout-hero-split">
  <div class="layout-hero-split__image"><!-- img placeholder --></div>
  <div class="layout-hero-split__text">
    <p class="layout-hero-split__eyebrow">EYEBROW</p>
    <h1 class="layout-hero-split__headline">%s</h1>
    <p class="layout-hero-split__body">副标题内容</p>
    <a class="layout-hero-split__cta" href="#">CTA 按钮</a>
  </div>
</section>`, titleHTML)
}

func bentoBodyHTML() string {
	return `
<div class="layout-bento">
  <div class="layout-bento__card">
    <div class="layout-bento__card-icon">📊</div>
    <div class="layout-bento__card-value">12,847</div>
    <p class="layout-bento__card-label">月活用户</p>
    <span class="layout-bento__card-trend up">↑ 23.5%</span>
  </div>
  <div class="layout-bento__card">
    <div class="layout-bento__card-icon">💰</div>
    <div class="layout-bento__card-value">$8.2M</div>
    <p class="layout-bento__card-label">ARR</p>
    <span class="layout-bento__card-trend up">↑ 12.1%</span>
  </div>
  <div class="layout-bento__card">
    <div class="layout-bento__card-icon">🎯</div>
    <div class="layout-bento__card-value">94.7%</div>
    <p class="layout-bento__card-label">客户留存率</p>
    <span class="layout-bento__card-trend down">↓ 0.3%</span>
  </div>
  <div class="layout-bento__card">
    <div class="layout-bento__card-icon">🚀</div>
    <div class="layout-bento__card-value">3,201</div>
    <p class="layout-bento__card-label">活跃项目</p>
    <span class="layout-bento__card-trend up">↑ 45.2%</span>
  </div>
</div>`
}

func galleryBodyHTML() string {
	out := `<div class="layout-gallery">`
	for i := 1; i <= 8; i++ {
		out += fmt.Sprintf(`
  <div class="layout-gallery__item">
    <div class="layout-gallery__image" style="background:var(--accent-2,#ccc)"></div>
    <div class="layout-gallery__caption">
      <h3>作品 %d</h3>
      <p>标签描述</p>
    </div>
  </div>`, i)
	}
	out += "</div>"
	return out
}

func (r *GenerationResult) SaveToFile(path string) error {
	return os.WriteFile(path, []byte(r.HTML), 0644)
}

// sortedKeys 返回 map 的排序后 keys，保证遍历确定性
func sortedKeys[K string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
