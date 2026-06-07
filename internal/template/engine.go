package template

import (
	"fmt"
	"os"
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
	result.LayoutCSS = layout.CSS
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
		for role, hex := range color.Colors {
			b.WriteString(fmt.Sprintf("  --%s: %s;\n", role, hex))
		}
	}
	for name, size := range font.Sizes {
		b.WriteString(fmt.Sprintf("  --size-%s: %s;\n", name, size))
	}
	b.WriteString("}")
	return b.String()
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
