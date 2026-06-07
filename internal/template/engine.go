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
	}
	result.CSSVars = buildCSSVariables(color, font)
	result.FontImports = buildFontImports(font)
	result.BaseCSS = buildBaseCSS(font)
	result.LayoutCSS = layout.CSS
	result.HTML = assembleHTML(result)
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
	// 优先使用显式定义的 google_fonts
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

	// 回退：从 fonts map 自动推断
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

func buildBaseCSS(font *TypographyTemplate) string {
	body := getOrDefault(font.Fonts, "body", "'Inter', sans-serif")
	disp := getOrDefault(font.Fonts, "display", body)
	mono := getOrDefault(font.Fonts, "mono", "'JetBrains Mono', monospace")
	return fmt.Sprintf(`* { box-sizing:border-box; margin:0; padding:0; }
html { scroll-behavior:smooth; }
body {
  font-family: %s; line-height:1.7;
  color: var(--text-primary, #1A1A1A);
  background: var(--background, #FFFFFF);
}
h1,h2,h3,h4 { font-family: %s; line-height:1.1; }
code,pre { font-family: %s; }
a { color: var(--accent, inherit); text-decoration:none; }
img { max-width: 100%%%%; height:auto; display:block; }
.container { max-width:1200px; margin:0 auto; padding:0 40px; }
@media (max-width:768px) { .container { padding:0 24px; } }`, body, disp, mono)
}


// displayName 从 name 字段提取中文名或英文名
// DisplayName 获取模板中文名或英文名
func DisplayName(name interface{}) string {
	switch v := name.(type) {
	case string:
		return v
	case map[string]interface{}:
		if zh, ok := v["zh"]; ok { return zh.(string) }
		if en, ok := v["en"]; ok { return en.(string) }
	}
	return ""
}
func getOrDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return def
}

func assembleHTML(r *GenerationResult) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>%s</title>
%s
<style>
%s
%s
%s
</style>
</head>
<body>
<main>
<section class="hero-split">
  <div class="hero-split__image"><!-- img --></div>
  <div class="hero-split__text">
    <p class="hero-split__eyebrow">EYEBROW</p>
    <h1 class="hero-split__headline">标题</h1>
    <p class="hero-split__body">副标题</p>
    <a class="hero-split__cta" href="#">CTA</a>
  </div>
</section>
</main>
</body>
</html>`, r.Title, r.FontImports, r.CSSVars, r.BaseCSS, r.LayoutCSS)
}

func (r *GenerationResult) SaveToFile(path string) error {
	return os.WriteFile(path, []byte(r.HTML), 0644)
}
