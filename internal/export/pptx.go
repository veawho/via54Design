// SPDX-License-Identifier: MIT OR AGPL-3.0

// via54Design — PPTX 导出器
// 纯 Go 实现，零外部依赖。PPTX = ZIP + XML。
// 替代 scripts/export_deck_pptx.mjs (Node.js + html2pptx.js)
package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PPTXSlide 单张幻灯片内容
type PPTXSlide struct {
	Title    string   // 主标题
	Subtitle string   // 副标题/眉标
	Body     []string // 正文段落
	Color    string   // 强调色 hex
	Image    string   // 图片路径(可选)
}

// ExportPPTX 从 slide 列表生成 PPTX 文件
// PPTX 是 ZIP 包，内含 OOXML (Office Open XML)
// 参考: ECMA-376 Office Open XML 标准
func ExportPPTX(slides []PPTXSlide, outputPath string, widescreen bool) error {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// ── 必需: [Content_Types].xml ──
	writeZip(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
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
`+genSlideTypes(len(slides))+`</Types>`)

	// ── 关系: _rels/.rels ──
	writeZip(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`)

	// ── docProps ──
	writeZip(w, "docProps/core.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties">
  <dc:title>via54Design</dc:title>
  <cp:lastModifiedBy>via54Design</cp:lastModifiedBy>
</cp:coreProperties>`)

	writeZip(w, "docProps/app.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties">
  <Application>via54Design</Application>
  <Slides>`+strconv.Itoa(len(slides))+`</Slides>
</Properties>`)

	// ── Theme ──
	writeZip(w, "ppt/theme/theme1.xml", pptxTheme())

	// ── Slide Master ──
	writeZip(w, "ppt/slideMasters/slideMaster1.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:cSld><p:spTree><p:nvGrpSpPr><p:nvPr/><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
</p:sldMaster>`)
	writeZip(w, "ppt/slideMasters/_rels/slideMaster1.xml.rels", `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/></Relationships>`)

	// ── Slide Layout ──
	writeZip(w, "ppt/slideLayouts/slideLayout1.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="blank">
  <p:cSld><p:spTree><p:nvGrpSpPr><p:nvPr/><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>
</p:sldLayout>`)

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
	writeZip(w, "ppt/presentation.xml", sz)

	writeZip(w, "ppt/_rels/presentation.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
`+slideRels+`</Relationships>`)

	// ── 每张 Slide ──
	for i, s := range slides {
		num := i + 1
		slideXML := buildSlideXML(s, num)
		writeZip(w, fmt.Sprintf("ppt/slides/slide%d.xml", num), slideXML)
		writeZip(w, fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num), `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`)
	}

	w.Close()

	// 写入文件
	os.MkdirAll(filepath.Dir(outputPath), 0755)
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func writeZip(w *zip.Writer, name, content string) {
	f, _ := w.Create(name)
	f.Write([]byte(content))
}

func genSlideTypes(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		b.WriteString(fmt.Sprintf(`  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i))
	}
	return b.String()
}

func buildSlideXML(s PPTXSlide, num int) string {
	// 颜色解析
	accentHex := "C43C3A"
	if s.Color != "" {
		accentHex = strings.TrimPrefix(s.Color, "#")
	}

	// 构建 body 段落
	var bodyXML strings.Builder
	for _, line := range s.Body {
		bodyXML.WriteString(fmt.Sprintf(`<a:p><a:r><a:rPr lang="zh-CN" sz="1800" dirty="0"/><a:t>%s</a:t></a:r></a:p>`, escapeXML(line)))
	}

	titleSize := "4400"
	if len(s.Title) > 20 {
		titleSize = "3200"
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
      <p:grpSpPr/>
      <!-- 背景 -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="2" name="Bg"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:solidFill><a:srgbClr val="F5F0E6"/></a:solidFill></p:spPr>
      </p:sp>
      <!-- 强调色装饰条 -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="3" name="Accent"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr>
          <a:xfrm><a:off x="0" y="0"/><a:ext cx="152400" cy="6858000"/></a:xfrm>
          <a:solidFill><a:srgbClr val="%s"/></a:solidFill>
        </p:spPr>
      </p:sp>
      <!-- 眉标 (Subtitle) -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="4" name="Eyebrow"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="914400" y="685800"/><a:ext cx="8229600" cy="457200"/></a:xfrm></p:spPr>
        <p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="1200" dirty="0" b="0"><a:solidFill><a:srgbClr val="%s"/></a:solidFill></a:rPr><a:t>%s</a:t></a:r></a:p></p:txBody>
      </p:sp>
      <!-- 标题 -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="5" name="Title"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="914400" y="1371600"/><a:ext cx="8229600" cy="1371600"/></a:xfrm></p:spPr>
        <p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="%s" dirty="0" b="1"><a:solidFill><a:srgbClr val="1A1A1A"/></a:solidFill></a:rPr><a:t>%s</a:t></a:r></a:p></p:txBody>
      </p:sp>
      <!-- 正文 -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="6" name="Body"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="914400" y="2971800"/><a:ext cx="8229600" cy="3200400"/></a:xfrm></p:spPr>
        <p:txBody><a:bodyPr/>%s</p:txBody>
      </p:sp>
      <!-- 页码 -->
      <p:sp>
        <p:nvSpPr><p:cNvPr id="7" name="PageNum"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
        <p:spPr><a:xfrm><a:off x="914400" y="6400800"/><a:ext cx="1371600" cy="274320"/></a:xfrm></p:spPr>
        <p:txBody><a:bodyPr/><a:p><a:r><a:rPr lang="zh-CN" sz="900" dirty="0"><a:solidFill><a:srgbClr val="999999"/></a:solidFill></a:rPr><a:t>%d / %d</a:t></a:r></a:p></p:txBody>
      </p:sp>
    </p:spTree>
  </p:cSld>
</p:sld>`, accentHex, accentHex, escapeXML(s.Subtitle), titleSize, escapeXML(s.Title), bodyXML.String(), num, num)
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

var _ = color.White // 保留 color 导入
