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

// via54Design — PPTX 导出器
// 纯 Go 实现，零外部依赖。PPTX = ZIP + XML。
// 替代 scripts/export_deck_pptx.mjs (Node.js + html2pptx.js)
// v2: 支持风格 + 配色模板系统
package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// PPTXSlide 单张幻灯片内容
type PPTXSlide struct {
	Title    string   // 主标题
	Subtitle string   // 副标题/眉标
	Body     []string // 正文段落
	Color    string   // 强调色 hex
	Image    string   // 图片路径(可选) // TODO: image embedding not yet implemented
}

// PPTXStyleElement 布局元素坐标
type PPTXStyleElement struct {
	Enabled     bool   `yaml:"enabled"`
	X           int    `yaml:"x"`
	Y           int    `yaml:"y"`
	W           int    `yaml:"w"`
	H           int    `yaml:"h"`
	FontSize    int    `yaml:"font_size"`
	Bold        bool   `yaml:"bold"`
	Color       string `yaml:"color"`
	Align       string `yaml:"align"`
	LineSpacing int    `yaml:"line_spacing"`
}

// PPTXAccentBar 装饰条配置
type PPTXAccentBar struct {
	Enabled bool   `yaml:"enabled"`
	Width   int    `yaml:"width"`
	Color   string `yaml:"color"`
}

// PPTXStyle PPTX 风格模板
type PPTXStyle struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Layout      struct {
		Background struct {
			Fill  string `yaml:"fill"`
			Color string `yaml:"color"`
		} `yaml:"background"`
		AccentBar  PPTXAccentBar    `yaml:"accent_bar"`
		Title      PPTXStyleElement `yaml:"title"`
		Subtitle   PPTXStyleElement `yaml:"subtitle"`
		Body       PPTXStyleElement `yaml:"body"`
		PageNumber PPTXStyleElement `yaml:"page_number"`
	} `yaml:"layout"`
}

// PPTXTheme 配色方案（从color-scheme YAML加载）
type PPTXTheme struct {
	ID   string `yaml:"id"`
	Name struct {
		Zh string `yaml:"zh"`
		En string `yaml:"en"`
	} `yaml:"name"`
	Palette []struct {
		Role string `yaml:"role"`
		Hex  string `yaml:"hex"`
	} `yaml:"palette"`
	CSSVariables string `yaml:"css_variables"`
}

// resolveColor 解析颜色值，支持 --var 引用和直接 hex
func resolveColor(colorSpec string, theme *PPTXTheme) string {
	if theme == nil || !strings.HasPrefix(colorSpec, "--") {
		return colorSpec
	}
	varName := strings.TrimPrefix(colorSpec, "--")
	for _, p := range theme.Palette {
		if p.Role == varName {
			return strings.TrimPrefix(p.Hex, "#")
		}
	}
	return strings.TrimPrefix(colorSpec, "#")
}

// resolveColorHex 解析颜色并去 #
func resolveColorHex(colorSpec string, theme *PPTXTheme) string {
	c := resolveColor(colorSpec, theme)
	return strings.TrimPrefix(c, "#")
}

// ExportPPTX 从 slide 列表生成 PPTX 文件
// styleID: "accent-bar" (默认), "minimal", "editorial", "bold"
// themePath: 配色 YAML 路径 (可选)
// baseDir:   templates 基础路径
// PPTX 是 ZIP 包，内含 OOXML (Office Open XML)
// 参考: ECMA-376 Office Open XML 标准
func ExportPPTX(slides []PPTXSlide, outputPath string, widescreen bool, styleID, themePath, baseDir string) error {
	// 加载风格模板
	style := loadPPTXStyle(styleID, baseDir)
	if style == nil {
		style = defaultPPTXStyle()
	}
	// 加载配色主题
	var theme *PPTXTheme
	if themePath != "" {
		theme = loadPPTXTheme(themePath)
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// ── 必需: [Content_Types].xml ──
	if err := writeZip(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Default Extension="png" ContentType="image/png"/>
  <Default Extension="jpg" ContentType="image/jpeg"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
  <Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
  <Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
  <Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
  <Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
`+genSlideTypes(len(slides))+`</Types>`); err != nil {
		return err
	}

	// ── 关系: _rels/.rels ──
	if err := writeZip(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`); err != nil {
		return err
	}

	// ── docProps ──
	if err := writeZip(w, "docProps/core.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties">
  <dc:title>via54Design</dc:title>
  <cp:lastModifiedBy>via54Design</cp:lastModifiedBy>
</cp:coreProperties>`); err != nil {
		return err
	}

	if err := writeZip(w, "docProps/app.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties">
  <Application>via54Design</Application>
  <Slides>`+strconv.Itoa(len(slides))+`</Slides>
</Properties>`); err != nil {
		return err
	}

	// ── Theme ──
	if err := writeZip(w, "ppt/theme/theme1.xml", pptxTheme()); err != nil {
		return err
	}

	// ── Slide Master ──
	if err := writeZip(w, "ppt/slideMasters/slideMaster1.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:cSld><p:spTree><p:nvGrpSpPr><p:nvPr/><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
</p:sldMaster>`); err != nil {
		return err
	}
	if err := writeZip(w, "ppt/slideMasters/_rels/slideMaster1.xml.rels", `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/></Relationships>`); err != nil {
		return err
	}

	// ── Slide Layout ──
	if err := writeZip(w, "ppt/slideLayouts/slideLayout1.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="blank">
  <p:cSld><p:spTree><p:nvGrpSpPr><p:nvPr/><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
</p:sldLayout>`); err != nil {
		return err
	}

	// ── Presentation ──
	slideRels := ""
	slideList := ""
	for i := range slides {
		num := i + 1
		slideList += fmt.Sprintf(`<p:sldId id="%d" r:id="rId%d"/>`, 256+i, i+2)
		slideRels += fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`, i+2, num)
	}

	sz := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>
  <p:sldIdLst>` + slideList + `</p:sldIdLst>
  <p:sldSz cx="12192000" cy="6858000"/>` // 16:9 widescreen
	if !widescreen {
		sz = strings.Replace(sz, `cx="12192000" cy="6858000"`, `cx="9144000" cy="6858000"`, 1) // 4:3
	}
	sz += `<p:notesSz cx="6858000" cy="9144000"/></p:presentation>`
	if err := writeZip(w, "ppt/presentation.xml", sz); err != nil {
		return err
	}

	if err := writeZip(w, "ppt/_rels/presentation.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
`+slideRels+`</Relationships>`); err != nil {
		return err
	}

	// ── 每张 Slide ──
	for i, s := range slides {
		num := i + 1
		slideXML := buildSlideXML(s, num, len(slides), style, theme)
		if err := writeZip(w, fmt.Sprintf("ppt/slides/slide%d.xml", num), slideXML); err != nil {
			return err
		}
		if err := writeZip(w, fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num), `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`); err != nil {
			return err
		}
	}

	w.Close()

	// 写入文件
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("mkdir output: %w", err)
	}
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func writeZip(w *zip.Writer, name, content string) error {
	f, err := w.Create(name)
	if err != nil {
		return fmt.Errorf("zip create %s: %w", name, err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("zip write %s: %w", name, err)
	}
	return nil
}

func genSlideTypes(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		b.WriteString(fmt.Sprintf(`  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i))
	}
	return b.String()
}

func buildSlideXML(s PPTXSlide, num, total int, style *PPTXStyle, theme *PPTXTheme) string {
	// 安全 fallback
	if style == nil {
		style = defaultPPTXStyle()
	}

	// 颜色解析
	bgColor := resolveColorHex(style.Layout.Background.Color, theme)
	accentColor := style.Layout.AccentBar.Color
	accentColor = resolveColorHex(accentColor, theme)
	if s.Color != "" {
		accentColor = strings.TrimPrefix(s.Color, "#")
	}

	// 构建 body 段落
	var bodyXML strings.Builder
	for _, line := range s.Body {
		bodyXML.WriteString(fmt.Sprintf(`<a:p><a:r><a:rPr lang="zh-CN" sz="%d" dirty="0"/><a:t>%s</a:t></a:r></a:p>`,
			style.Layout.Body.FontSize, escapeXML(line)))
	}

	// 标题字号自适应
	titleSize := style.Layout.Title.FontSize
	if len(s.Title) > 20 && titleSize > 3200 {
		titleSize = 3200
	}

	// ── 构建 spTree 内部 ──
	var spTree strings.Builder
	spTree.WriteString(`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>`)
	spTree.WriteString(`<p:grpSpPr/>`)

	// 背景
	spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Bg"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:solidFill><a:srgbClr val="%s"/></a:solidFill></p:spPr></p:sp>`, bgColor))

	// 强调色装饰条 (条件)
	if style.Layout.AccentBar.Enabled {
		barW := style.Layout.AccentBar.Width
		if barW == 0 {
			barW = 152400
		}
		spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="3" name="Accent"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="6858000"/></a:xfrm><a:solidFill><a:srgbClr val="%s"/></a:solidFill></p:spPr></p:sp>`, barW, accentColor))
	}

	// 眉标 (Subtitle)
	if s.Subtitle != "" {
		sub := style.Layout.Subtitle
		subColor := resolveColorHex(sub.Color, theme)
		spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="4" name="Eyebrow"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm></p:spPr><p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="%d" dirty="0" b="0"><a:solidFill><a:srgbClr val="%s"/></a:solidFill></a:rPr><a:t>%s</a:t></a:r></a:p></p:txBody></p:sp>`,
			sub.X, sub.Y, sub.W, sub.H, sub.FontSize, subColor, escapeXML(s.Subtitle)))
	}

	// 标题
	title := style.Layout.Title
	titleColor := resolveColorHex(title.Color, theme)
	titleBold := "0"
	if title.Bold {
		titleBold = "1"
	}
	spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="5" name="Title"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm></p:spPr><p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="%d" dirty="0" b="%s"><a:solidFill><a:srgbClr val="%s"/></a:solidFill></a:rPr><a:t>%s</a:t></a:r></a:p></p:txBody></p:sp>`,
		title.X, title.Y, title.W, title.H, titleSize, titleBold, titleColor, escapeXML(s.Title)))

	// 正文
	body := style.Layout.Body
	if s.Body != nil && len(s.Body) > 0 {
		spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="6" name="Body"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm></p:spPr><p:txBody><a:bodyPr/>%s</p:txBody></p:sp>`,
			body.X, body.Y, body.W, body.H, bodyXML.String()))
	}

	// 页码 (条件)
	if style.Layout.PageNumber.Enabled {
		pn := style.Layout.PageNumber
		pnColor := resolveColorHex(pn.Color, theme)
		spTree.WriteString(fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="7" name="PageNum"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm></p:spPr><p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="%d" dirty="0"><a:solidFill><a:srgbClr val="%s"/></a:solidFill></a:rPr><a:t>%d / %d</a:t></a:r></a:p></p:txBody></p:sp>`,
			pn.X, pn.Y, pn.W, pn.H, pn.FontSize, pnColor, num, total))
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:cSld>
    <p:spTree>
      %s
    </p:spTree>
  </p:cSld>
</p:sld>`, spTree.String())
}

func pptxTheme() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="via54">
  <a:themeElements>
    <a:clrScheme name="via54">
      <a:dk1><a:srgbClr val="1A1A1A"/></a:dk1>
      <a:lt1><a:srgbClr val="F5F0E6"/></a:lt1>
      <a:dk2><a:srgbClr val="6B6B6B"/></a:dk2>
      <a:lt2><a:srgbClr val="FFFFFF"/></a:lt2>
      <a:accent1><a:srgbClr val="C43C3A"/></a:accent1>
      <a:accent2><a:srgbClr val="2A2A2A"/></a:accent2>
      <a:accent3><a:srgbClr val="D9D4CA"/></a:accent3>
    </a:clrScheme>
    <a:fontScheme name="via54">
      <a:majorFont><a:latin typeface="Source Han Sans"/><a:ea typeface="Source Han Sans"/></a:majorFont>
      <a:minorFont><a:latin typeface="Source Han Sans"/><a:ea typeface="Source Han Sans"/></a:minorFont>
    </a:fontScheme>
    <a:fmtScheme name="via54"/>
  </a:themeElements>
</a:theme>`
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// PPTXSlideFromBeat 从叙事节拍构建 PPTXSlide
func PPTXSlideFromBeat(act, voiceover, mood string, idx, total int) PPTXSlide {
	colorMap := map[string]string{
		"mysterious": "2E5CB8", "aspirational": "C89B3C", "confident": "C23A2B",
		"urgent": "FF2D95", "calm": "5B8C5A", "curious": "7A4B5C",
		"excited": "58CC02", "inspiring": "CC785C", "informative": "3498DB",
		"focused": "C23A2B", "warm": "E8A838", "practical": "C06C4C",
		"insightful": "2E5CB8", "hopeful": "5B8C5A", "frustrated": "1A1A1A",
		"tense": "000000", "triumphant": "58CC02",
	}
	c := colorMap[mood]
	if c == "" {
		c = "C43C3A"
	}
	return PPTXSlide{
		Title:    act,
		Subtitle: "via54 叙事引擎",
		Body:     []string{voiceover},
		Color:    c,
	}
}

func loadPPTXStyle(styleID, baseDir string) *PPTXStyle {
	if styleID == "" {
		styleID = "accent-bar"
	}
	path := filepath.Join(baseDir, "templates", "pptx-styles", styleID+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("loadPPTXStyle: read %s: %v\n", path, err)
		return nil
	}
	var s PPTXStyle
	if err := yaml.Unmarshal(data, &s); err != nil {
		fmt.Printf("loadPPTXStyle: unmarshal %s: %v\n", path, err)
		return nil
	}
	return &s
}

func loadPPTXTheme(themePath string) *PPTXTheme {
	data, err := os.ReadFile(themePath)
	if err != nil {
		fmt.Printf("loadPPTXTheme: read %s: %v\n", themePath, err)
		return nil
	}
	var t PPTXTheme
	if err := yaml.Unmarshal(data, &t); err != nil {
		fmt.Printf("loadPPTXTheme: unmarshal %s: %v\n", themePath, err)
		return nil
	}
	return &t
}

func defaultPPTXStyle() *PPTXStyle {
	s := &PPTXStyle{
		ID:   "default",
		Name: "默认",
	}
	s.Layout.Background.Fill = "solid"
	s.Layout.Background.Color = "F5F0E6"
	s.Layout.AccentBar = PPTXAccentBar{Enabled: true, Width: 152400, Color: "C43C3A"}
	s.Layout.Title = PPTXStyleElement{X: 914400, Y: 1371600, W: 8229600, H: 1371600, FontSize: 4400, Bold: true, Color: "1A1A1A"}
	s.Layout.Subtitle = PPTXStyleElement{X: 914400, Y: 685800, W: 8229600, H: 457200, FontSize: 1200, Color: "C43C3A"}
	s.Layout.Body = PPTXStyleElement{X: 914400, Y: 2971800, W: 8229600, H: 3200400, FontSize: 1800, Color: "4A4A4A"}
	s.Layout.PageNumber = PPTXStyleElement{X: 914400, Y: 6400800, W: 1371600, H: 274320, FontSize: 900, Color: "999999", Enabled: true}
	return s
}
