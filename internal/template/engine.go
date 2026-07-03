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
	"gopkg.in/yaml.v3"
	"html"
	"os"
	"sort"
	"strings"
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
	lp, err := e.Registry.ResolveLayout(layoutID)
	if err != nil {
		return nil, fmt.Errorf("布局 %q: %w", layoutID, err)
	}
	cp, err := e.Registry.ResolveColorScheme(colorID)
	if err != nil {
		return nil, fmt.Errorf("配色 %q: %w", colorID, err)
	}
	fp, err := e.Registry.ResolveTypography(fontID)
	if err != nil {
		return nil, fmt.Errorf("字体 %q: %w", fontID, err)
	}

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
		LayoutID:         layoutID,
		ColorID:          colorID,
		FontID:           fontID,
		Title:            title,
		LetteringSVG:     letteringSVG,
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
	case "image-container":
		return "image"
	case "text-container":
		return "text"
	case "gallery-item":
		return "item"
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
		"Inter": true, "Geist": true, "JetBrains Mono": true, "Fraunces": true,
		"Playfair Display": true, "Noto Serif SC": true, "Noto Sans SC": true,
		"EB Garamond": true, "Nunito": true, "Baloo 2": true,
		"Archivo Black": true, "Archivo": true, "LXGW WenKai": true, "ZCOOL XiaoWei": true,
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

	// XSS 防御: HTML 转义用户输入的 title
	safeTitle := html.EscapeString(r.Title)
	safeTitleHTML := html.EscapeString(titleHTML)
	// bodyContent 内部 (e.g. heroBodyHTML) 也需要安全转义
	_ = safeTitleHTML // 未直接用, 由 heroBodyHTML 内部处理

	// 布局特有 HTML body 内容 (根据 layout ID 选择)
	var bodyContent string
	switch layout.ID {
	case "bento-grid-2x2":
		bodyContent = bentoBodyHTML()
	case "gallery-waterfall":
		bodyContent = galleryBodyHTML()
	case "dashboard-3pane":
		bodyContent = dashboardBodyHTML()
	case "landing-pricing":
		bodyContent = landingPricingBodyHTML(titleHTML)
	case "docs-sidebar":
		bodyContent = docsSidebarBodyHTML()
	case "blog-magazine":
		bodyContent = blogMagazineBodyHTML()
	case "pricing-comparison":
		bodyContent = pricingComparisonBodyHTML()
	case "product-feature-grid":
		bodyContent = productFeatureGridBodyHTML()
	case "settings-form":
		bodyContent = settingsFormBodyHTML()
	default:
		bodyContent = heroBodyHTML(titleHTML)
	}

	// 演示模式：外层包裹 .presentation-mode 容器
	bodyWrapper := ""
	if r.PresentationMode {
		bodyWrapper = `<div class="presentation-mode">` + bodyContent + `</div>`
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
</html>`, vpMeta, safeTitle, r.FontImports, r.CSSVars, r.BaseCSS, r.LayoutCSS, bodyContent)
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
</section>`, html.EscapeString(titleHTML))
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

func dashboardBodyHTML() string {
	return `
<div class="layout-dashboard">
  <aside class="layout-dashboard__sidebar">
    <div class="logo">via54Design</div>
    <nav>
      <a href="#" class="active">Overview</a>
      <a href="#">Analytics</a>
      <a href="#">Projects</a>
      <a href="#">Settings</a>
    </nav>
  </aside>
  <main class="layout-dashboard__main">
    <header>
      <h2>Console Overview</h2>
      <div class="user-profile">User</div>
    </header>
    <div class="metrics-grid">
      <div class="metric-card"><h3>Active Users</h3><p>24.8K</p></div>
      <div class="metric-card"><h3>Conversion Rate</h3><p>3.42%</p></div>
      <div class="metric-card"><h3>Bounce Rate</h3><p>42.1%</p></div>
    </div>
    <div class="content-block">
      <h3>System Performance</h3>
      <div class="chart-placeholder">Chart</div>
    </div>
  </main>
  <aside class="layout-dashboard__aside">
    <h3>Activity Feed</h3>
    <ul>
      <li>User registered 2m ago</li>
      <li>Deployment success 15m ago</li>
      <li>Database backup ok 1h ago</li>
    </ul>
  </aside>
</div>`
}

func landingPricingBodyHTML(titleHTML string) string {
	return fmt.Sprintf(`
<div class="layout-landing">
  <nav class="layout-landing__nav">
    <div class="logo">via54Design</div>
    <div class="nav-links">
      <a href="#">Features</a>
      <a href="#">Pricing</a>
      <a href="#">Docs</a>
    </div>
    <button class="nav-cta">Get Started</button>
  </nav>
  <header class="layout-landing__hero">
    <h1>%s</h1>
    <p>Aesthetic, modular, responsive, and performance-tuned layouts powered by Golang & modern design agents.</p>
    <div class="cta-buttons">
      <a href="#" class="btn-primary">Start Free Trial</a>
      <a href="#" class="btn-secondary">View Docs</a>
    </div>
  </header>
  <section class="layout-landing__pricing">
    <article class="layout-landing__pricing-card">
      <h3>Hobby</h3>
      <div class="price">$0<span>/mo</span></div>
      <p>Perfect for exploring design capabilities.</p>
      <ul>
        <li>10 Projects</li>
        <li>Community Support</li>
      </ul>
      <button class="btn-card">Start Free</button>
    </article>
    <article class="layout-landing__pricing-card layout-landing__pricing-card--featured">
      <h3>Pro</h3>
      <div class="price">$19<span>/mo</span></div>
      <p>Everything you need for production apps.</p>
      <ul>
        <li>Unlimited Projects</li>
        <li>Priority SLA Support</li>
        <li>Custom Domain Mapping</li>
      </ul>
      <button class="btn-card btn-card--featured">Upgrade Pro</button>
    </article>
    <article class="layout-landing__pricing-card">
      <h3>Enterprise</h3>
      <div class="price">Custom</div>
      <p>Dedicated scale and custom integration.</p>
      <ul>
        <li>Single Sign-On (SSO)</li>
        <li>Dedicated Infrastructure</li>
      </ul>
      <button class="btn-card">Contact Sales</button>
    </article>
  </section>
</div>`, html.EscapeString(titleHTML))
}

func docsSidebarBodyHTML() string {
	return `
<div class="layout-docs">
  <aside class="layout-docs__sidebar">
    <h3>Getting Started</h3>
    <a href="#" class="active">Introduction</a>
    <a href="#">Installation</a>
    <a href="#">Quick Start</a>
    <h3>Customization</h3>
    <a href="#">Color Schemes</a>
    <a href="#">Layout Templates</a>
  </aside>
  <article class="layout-docs__content">
    <h1>Documentation Guide</h1>
    <p>Welcome to the official developer guide. Discover how to create stunning UI components, define responsive layouts, and configure advanced color schemes using the via54Design system.</p>
    <h2>Core Concepts</h2>
    <p>Layouts are defined in structured YAML configurations featuring multiple screen adaptions, baseline aspect ratios, and golden ratio spacing systems.</p>
  </article>
  <aside class="layout-docs__toc">
    <h3>On This Page</h3>
    <a href="#">Overview</a>
    <a href="#">Core Concepts</a>
    <a href="#">Next Steps</a>
  </aside>
</div>`
}

func blogMagazineBodyHTML() string {
	return `
<article class="layout-blog">
  <header class="layout-blog__header">
    <div class="meta">DESIGN ARCHIVE • 2026</div>
    <h1>The Art of Aesthetics: Vibe Coding in the Age of Agents</h1>
    <p class="subtitle">Exploring Lovable, v0, and the craft of high-fidelity visual design systems.</p>
  </header>
  <div class="layout-blog__image" style="background:var(--accent,#ccc); min-height:400px; border-radius:12px;"></div>
  <div class="layout-blog__content">
    <p>Design is not just what it looks like and feels like. Design is how it works. In the modern era of agentic workflows, software development is moving towards high-speed visual iterations.</p>
    <blockquote>"Aesthetics is the language of quality. If a product looks elegant and responds instantly, it breeds trust."</blockquote>
    <p>By leveraging standard color variables, golden ratio spacing scales, and fluid typography, we build websites that adapt gracefully to any viewport.</p>
  </div>
</article>`
}

func pricingComparisonBodyHTML() string {
	return `
<div class="layout-pricing-comparison">
  <header>
    <h2>Compare Plan Features</h2>
    <p>Select the plan that fits your operational scale.</p>
  </header>
  <table class="layout-pricing-comparison__table">
    <thead>
      <tr>
        <th>Feature</th>
        <th>Hobby</th>
        <th>Pro</th>
        <th>Enterprise</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>API Access</td>
        <td>100 req/day</td>
        <td>Unlimited</td>
        <td>Custom SLA</td>
      </tr>
      <tr>
        <td>Custom Branding</td>
        <td>❌</td>
        <td>✅</td>
        <td>✅</td>
      </tr>
      <tr>
        <td>SSO Login</td>
        <td>❌</td>
        <td>❌</td>
        <td>✅</td>
      </tr>
    </tbody>
  </table>
</div>`
}

func productFeatureGridBodyHTML() string {
	return `
<div class="layout-product-features">
  <header>
    <h2>Built for High-Speed Visual Delivery</h2>
    <p>Everything you need to craft pixel-perfect web interfaces.</p>
  </header>
  <div class="layout-product-features__grid">
    <div class="layout-product-features__card">
      <div class="icon">⚡</div>
      <h3>Ultra-fast Compiles</h3>
      <p>Go compiled binaries output fully integrated layouts in milliseconds.</p>
    </div>
    <div class="layout-product-features__card">
      <div class="icon">🎨</div>
      <h3>Golden Spacing</h3>
      <p>Mathematical design harmony utilizing standard golden ratio scaling.</p>
    </div>
    <div class="layout-product-features__card">
      <div class="icon">📱</div>
      <h3>4-Device Responsive</h3>
      <p>Optimized for TV, Desktop, Tablet, and Phone layout structures.</p>
    </div>
  </div>
</div>`
}

func settingsFormBodyHTML() string {
	return `
<div class="layout-settings">
  <aside class="layout-settings__nav">
    <a href="#" class="active">General Profile</a>
    <a href="#">API Keys</a>
    <a href="#">Team Access</a>
  </aside>
  <form class="layout-settings__form" onsubmit="event.preventDefault()">
    <h2>Account Settings</h2>
    <div class="form-group">
      <label>Workspace Name</label>
      <input type="text" value="My Awesome Workspace">
    </div>
    <div class="form-group">
      <label>Contact Email</label>
      <input type="email" value="admin@workspace.com">
    </div>
    <div class="form-actions">
      <button class="btn-save">Save Changes</button>
    </div>
  </form>
</div>`
}


func (r *GenerationResult) SaveToFile(path string) error {
	return os.WriteFile(path, []byte(r.HTML), 0644)
}

// sortedKeys 返回 map 的排序后 keys，保证遍历确定性
func sortedKeys[K string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
