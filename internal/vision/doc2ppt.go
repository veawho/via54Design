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

// via54Design — Document / Image → PPT Pipeline
// Pure Go implementation using archive/zip + encoding/xml for office formats.
// Replaces scripts/doc2ppt.py — zero external dependencies.
package vision

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ─── Public API ───────────────────────────────────────────────────────────

// ExtractContent auto-detects file type and extracts content.
// Supported: .docx, .md/.markdown, .txt, .pptx, .png/.jpg/.jpeg/.gif
func ExtractContent(path string) map[string]interface{} {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".docx":
		return extractFromDocx(path)
	case ".md", ".markdown":
		return extractFromMD(path)
	case ".txt":
		return extractFromTxt(path)
	case ".pptx":
		return extractFromPPTX(path)
	case ".png", ".jpg", ".jpeg", ".gif":
		return extractFromImage(path)
	default:
		return map[string]interface{}{
			"error": fmt.Sprintf("unsupported format: %s", ext),
		}
	}
}

// GeneratePPTFramework generates a PPT framework from extracted content.
func GeneratePPTFramework(content map[string]interface{}, userPrompt string) map[string]interface{} {
	ctype, _ := content["type"].(string)
	title, _ := content["title"].(string)
	if title == "" {
		title = "演示文稿"
	}

	slides := make([]map[string]interface{}, 0)

	// Cover slide
	subtitle := "via54Design 生成"
	if userPrompt != "" {
		subtitle = userPrompt
	}
	slides = append(slides, map[string]interface{}{
		"type":     "cover",
		"title":    title,
		"subtitle": subtitle,
		"mood":     "inspiring",
	})

	switch ctype {
	case "docx":
		headings, _ := content["headings"].([]interface{})
		for i, h := range headings {
			if i >= 15 {
				break
			}
			if hm, ok := h.(map[string]interface{}); ok {
				text, _ := hm["text"].(string)
				level := 2
				if l, ok := hm["level"].(float64); ok {
					level = int(l)
				}
				slides = append(slides, map[string]interface{}{
					"type":     "content",
					"title":    text,
					"subtitle": fmt.Sprintf("H%d · 章节页", level),
					"mood":     "informative",
					"layout":   "hero-split-16-9",
				})
			}
		}
		slides = append(slides, map[string]interface{}{
			"type": "summary", "title": "总结", "mood": "confident",
		})

	case "markdown":
		headings, _ := content["headings"].([]interface{})
		for i, h := range headings {
			if i >= 12 {
				break
			}
			if hm, ok := h.(map[string]interface{}); ok {
				text, _ := hm["text"].(string)
				level := 1
				if l, ok := hm["level"].(float64); ok {
					level = int(l)
				}
				slides = append(slides, map[string]interface{}{
					"type":     "content",
					"title":    text,
					"subtitle": fmt.Sprintf("Level %d", level),
					"mood":     "informative",
				})
			}
		}
		slides = append(slides, map[string]interface{}{
			"type": "summary", "title": "总结与展望", "mood": "inspiring",
		})

	case "pptx":
		existing, _ := content["slides"].([]interface{})
		for i, s := range existing {
			if i >= 15 {
				break
			}
			if sm, ok := s.(map[string]interface{}); ok {
				slideTitle := fmt.Sprintf("第%.0f页", sm["index"])
				if texts, ok := sm["texts"].([]interface{}); ok && len(texts) > 0 {
					slideTitle, _ = texts[0].(string)
				}
				imgCount := 0
				if ic, ok := sm["image_count"].(float64); ok {
					imgCount = int(ic)
				}
				slides = append(slides, map[string]interface{}{
					"type":       "content",
					"title":      slideTitle,
					"subtitle":   fmt.Sprintf("Slide %.0f", sm["index"]),
					"mood":       "informative",
					"has_images": imgCount > 0,
				})
			}
		}
		slides = append(slides, map[string]interface{}{
			"type": "summary", "title": "总结", "mood": "confident",
		})

	case "image":
		vc, _ := content["visual_context"].(map[string]interface{})
		moods := []string{"professional"}
		if m, ok := vc["moods"].([]interface{}); ok {
			moodsStr := make([]string, 0, len(m))
			for _, mm := range m {
				if str, ok := mm.(string); ok {
					moodsStr = append(moodsStr, str)
				}
			}
			moods = moodsStr
		}
		topics, _ := content["suggested_topics"].([]interface{})
		if len(topics) == 0 {
			topics = []interface{}{"方案介绍"}
		}

		styles := []string{"professional"}
		if st, ok := vc["styles"].([]interface{}); ok {
			stylesStr := make([]string, 0, len(st))
			for _, s := range st {
				if str, ok := s.(string); ok {
					stylesStr = append(stylesStr, str)
				}
			}
			styles = stylesStr
		}

		mainMood := "inspiring"
		if len(moods) > 0 {
			mainMood = moods[0]
		}

		slides = append(slides, map[string]interface{}{
			"type":     "section",
			"title":    "视觉灵感",
			"subtitle": fmt.Sprintf("风格: %s | 情绪: %s", strings.Join(styles, " · "), strings.Join(moods, " · ")),
			"mood":     mainMood,
		})
		for i, topic := range topics {
			if i >= 4 {
				break
			}
			topicStr, _ := topic.(string)
			imageHint := ""
			if gp, ok := content["generated_prompt"].(string); ok {
				if len(gp) > 100 {
					imageHint = gp[:100]
				} else {
					imageHint = gp
				}
			}
			slides = append(slides, map[string]interface{}{
				"type":       "content",
				"title":      topicStr,
				"mood":       "informative",
				"image_hint": imageHint,
			})
		}
		slides = append(slides, map[string]interface{}{
			"type": "summary", "title": "下一步", "mood": "hopeful",
		})

	default: // text / unknown
		preview, _ := content["content_preview"].(string)
		sentences := splitSentences(preview)
		for i, s := range sentences {
			if i >= 8 {
				break
			}
			trunc := s
			if len(trunc) > 60 {
				trunc = trunc[:60]
			}
			slides = append(slides, map[string]interface{}{
				"type":  "content",
				"title": trunc,
				"mood":  "informative",
			})
		}
		slides = append(slides, map[string]interface{}{
			"type": "summary", "title": "总结", "mood": "confident",
		})
	}

	// Build content_info
	contentInfo := map[string]interface{}{
		"paragraphs": 0,
		"headings":   0,
		"images":     0,
	}
	if tp, ok := content["total_paragraphs"].(float64); ok {
		contentInfo["paragraphs"] = int(tp)
	}
	if th, ok := content["total_headings"].(float64); ok {
		contentInfo["headings"] = int(th)
	}
	if ti, ok := content["total_images"].(float64); ok {
		contentInfo["images"] = int(ti)
	} else if imgs, ok := content["images"]; ok {
		contentInfo["images"] = imgs
	}

	titleTrunc := title
	if len(titleTrunc) > 40 {
		titleTrunc = titleTrunc[:40]
	}
	recommendedCmd := fmt.Sprintf(
		"via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial --title \"%s\" --presentation",
		titleTrunc,
	)

	return map[string]interface{}{
		"title":               title,
		"type":                ctype,
		"total_slides":        len(slides),
		"slides":              slides,
		"content_info":        contentInfo,
		"recommended_command": recommendedCmd,
		"user_guidance":       generateGuidance(ctype, content),
	}
}

// Story2PPT is the main entry point: file → content → PPT framework.
func Story2PPT(path string, userPrompt string) map[string]interface{} {
	content := ExtractContent(path)
	if err, hasErr := content["error"]; hasErr && err != nil {
		return content
	}

	framework := GeneratePPTFramework(content, userPrompt)
	framework["source_file"] = path
	framework["user_prompt"] = userPrompt
	return framework
}

// ─── DOCX Extraction ──────────────────────────────────────────────────────

// docxBody is the body element in word/document.xml (namespace-stripped).
type docxBody struct {
	Paragraphs []docxParagraph `xml:"body>p"`
}

type docxParagraph struct {
	PPr  *docxPPr  `xml:"pPr"`
	Runs []docxRun `xml:"r"`
}

type docxPPr struct {
	PStyle *docxStyle `xml:"pStyle"`
}

type docxStyle struct {
	Val string `xml:"val,attr"`
}

type docxRun struct {
	Text string `xml:"t"`
}

// extractFromDocx extracts text structure from .docx files.
func extractFromDocx(path string) map[string]interface{} {
	content, err := readZipXMLContent(path, "word/document.xml")
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("cannot read docx: %v", err)}
	}

	// Strip XML namespace prefixes for easier parsing
	content = stripXMLNS(content)

	var body docxBody
	if err := xml.Unmarshal([]byte(content), &body); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("xml parse error: %v", err)}
	}

	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	headings := make([]map[string]interface{}, 0)
	paragraphs := make([]map[string]interface{}, 0)
	currentSection := "intro"

	for _, p := range body.Paragraphs {
		// Concatenate all run text
		var text string
		for _, r := range p.Runs {
			text += r.Text
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		isHeading := false
		if p.PPr != nil && p.PPr.PStyle != nil {
			styleName := p.PPr.PStyle.Val
			if strings.HasPrefix(styleName, "Heading") || strings.HasPrefix(styleName, "heading") {
				isHeading = true
				level := 2
				// Try to extract heading level from style name (e.g., "Heading 1" → 1)
				parts := strings.Fields(styleName)
				if len(parts) >= 2 {
					fmt.Sscanf(parts[len(parts)-1], "%d", &level)
				}
				headings = append(headings, map[string]interface{}{
					"level": level,
					"text":  text,
				})
				currentSection = text
			}
		}

		if !isHeading {
			paragraphs = append(paragraphs, map[string]interface{}{
				"section": currentSection,
				"text":    text,
			})
		}
	}

	// First heading is the title
	title := baseName
	if len(headings) > 0 {
		if t, ok := headings[0]["text"].(string); ok {
			title = t
		}
	}

	// Content preview
	previewLines := make([]string, 0, 10)
	for i, p := range paragraphs {
		if i >= 10 {
			break
		}
		t, _ := p["text"].(string)
		if len(t) > 120 {
			t = t[:120]
		}
		previewLines = append(previewLines, t)
	}

	return map[string]interface{}{
		"type":             "docx",
		"title":            title,
		"headings":         headings,
		"paragraphs":       paragraphs,
		"total_paragraphs": len(paragraphs),
		"total_headings":   len(headings),
		"images":           []interface{}{},
		"content_preview":  strings.Join(previewLines, "\n"),
	}
}

// ─── MD Extraction ────────────────────────────────────────────────────────

// extractFromMD extracts structure from Markdown files.
func extractFromMD(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("cannot read file: %v", err)}
	}
	content := string(data)

	headings := make([]map[string]interface{}, 0)
	paragraphs := make([]string, 0)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if strings.HasPrefix(stripped, "#") {
			// Count leading # for level
			level := 0
			for _, ch := range stripped {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			text := strings.TrimSpace(stripped[level:])
			headings = append(headings, map[string]interface{}{
				"level": level,
				"text":  text,
			})
		} else if stripped != "" && !strings.HasPrefix(stripped, "```") &&
			!strings.HasPrefix(stripped, "---") && !strings.HasPrefix(stripped, ">") &&
			!strings.HasPrefix(stripped, "-") && !strings.HasPrefix(stripped, "*") &&
			!strings.HasPrefix(stripped, "|") {
			if len(stripped) > 200 {
				paragraphs = append(paragraphs, stripped[:200])
			} else {
				paragraphs = append(paragraphs, stripped)
			}
		}
	}

	// Find title (first H1)
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	title := baseName
	for _, h := range headings {
		if l, ok := h["level"].(int); ok && l == 1 {
			if t, ok := h["text"].(string); ok {
				title = t
			}
			break
		}
	}

	previewLines := paragraphs
	if len(previewLines) > 8 {
		previewLines = previewLines[:8]
	}

	return map[string]interface{}{
		"type":            "markdown",
		"title":           title,
		"headings":        headings,
		"total_headings":  len(headings),
		"content_preview": strings.Join(previewLines, "\n"),
		"raw_length":      len(content),
	}
}

// ─── TXT Extraction ───────────────────────────────────────────────────────

// extractFromTxt extracts plain text content.
func extractFromTxt(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("cannot read file: %v", err)}
	}
	text := string(data)
	preview := text
	if len(preview) > 800 {
		preview = preview[:800]
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return map[string]interface{}{
		"type":            "text",
		"title":           title,
		"content_preview": preview,
		"raw_length":      len(text),
		"total_lines":     strings.Count(text, "\n") + 1,
	}
}

// ─── PPTX Extraction ──────────────────────────────────────────────────────

// pptx types for decoding XML (namespace-stripped)
type pptxPresentation struct {
	SldIDLst *pptxSldIDLst `xml:"sldIdLst"`
}

type pptxSldIDLst struct {
	SldIDs []pptxSldID `xml:"sldId"`
}

type pptxSldID struct {
	ID  int    `xml:"id,attr"`
	RID string `xml:"rid,attr"`
}

type pptxSlideXML struct {
	SpTree *pptxSpTreeXML `xml:"cSld>spTree"`
}

type pptxSpTreeXML struct {
	Sps []pptxSpXML `xml:"sp"`
}

type pptxSpXML struct {
	TxBody *pptxTxBodyXML `xml:"txBody"`
}

type pptxTxBodyXML struct {
	Paragraphs []pptxParagraphXML `xml:"p"`
}

type pptxParagraphXML struct {
	Runs []pptxRunXML `xml:"r"`
}

type pptxRunXML struct {
	Text string `xml:"t"`
}

// extractFromPPTX extracts content from existing .pptx files.
func extractFromPPTX(path string) map[string]interface{} {
	r, err := zip.OpenReader(path)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("cannot open pptx: %v", err)}
	}
	defer r.Close()

	// Read the rels file to map rIds to slide file paths
	relsContent, err := readZipFileOpen(r, "ppt/_rels/presentation.xml.rels")
	sldMap := make(map[string]string)
	if err == nil {
		relsContent = stripXMLNS(relsContent)
		sldMap = parseRels(relsContent)
	}

	// Read presentation.xml to get slide order
	presContent, err := readZipFileOpen(r, "ppt/presentation.xml")
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("cannot read presentation.xml: %v", err)}
	}
	presContent = stripXMLNS(presContent)

	var pres pptxPresentation
	if err := xml.Unmarshal([]byte(presContent), &pres); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("xml parse presentation: %v", err)}
	}

	// Process slides
	slides := make([]map[string]interface{}, 0)
	allText := make([]string, 0)

	if pres.SldIDLst == nil {
		return map[string]interface{}{
			"type":   "pptx",
			"title":  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
			"total_slides": 0,
			"slides": slides,
			"content_preview": "",
		}
	}

	for i, sld := range pres.SldIDLst.SldIDs {
		if i >= 20 { // limit to first 20 slides
			break
		}

		// Resolve slide file path from rels map
		sldFile := sldMap[sld.RID]
		if sldFile == "" {
			sldFile = fmt.Sprintf("ppt/slides/slide%d.xml", i+1)
		}

		sldContent, err := readZipFileOpen(r, sldFile)
		if err != nil {
			continue
		}
		sldContent = stripXMLNS(sldContent)

		var slide pptxSlideXML
		if err := xml.Unmarshal([]byte(sldContent), &slide); err != nil {
			continue
		}

		texts := make([]string, 0)
		if slide.SpTree != nil {
			for _, sp := range slide.SpTree.Sps {
				if sp.TxBody != nil {
					for _, p := range sp.TxBody.Paragraphs {
						var paraText string
						for _, r := range p.Runs {
							paraText += r.Text
						}
						paraText = strings.TrimSpace(paraText)
						if paraText != "" {
							texts = append(texts, paraText)
							allText = append(allText, paraText)
						}
					}
				}
			}
		}

		slides = append(slides, map[string]interface{}{
			"index":       i + 1,
			"texts":       texts,
			"image_count": 0, // PPTX image detection not implemented in pure Go
		})
	}

	// Count total text blocks
	totalTextBlocks := 0
	for _, s := range slides {
		if texts, ok := s["texts"].([]string); ok {
			totalTextBlocks += len(texts)
		}
	}

	previewLines := allText
	if len(previewLines) > 15 {
		previewLines = previewLines[:15]
	}

	return map[string]interface{}{
		"type":              "pptx",
		"title":             strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		"total_slides":      len(slides),
		"total_images":      0, // image detection not implemented
		"total_text_blocks": totalTextBlocks,
		"slides":            slides,
		"content_preview":   strings.Join(previewLines, "\n"),
	}
}

// ─── Image Extraction ─────────────────────────────────────────────────────

// extractFromImage extracts visual analysis from image for PPT context.
func extractFromImage(path string) map[string]interface{} {
	analysis := AnalyzeImageToMap(path)
	prompt := BuildPromptFromAnalysisMap(analysis, "")

	colorfulness, _ := analysis["colorfulness"].(string)
	suggestedTopics := []string{"综合提案", "方案介绍", "案例分享"}
	if strings.Contains(colorfulness, "vibrant") {
		suggestedTopics = []string{"品牌展示", "产品亮点", "创意概念"}
	} else if strings.Contains(colorfulness, "muted") {
		suggestedTopics = []string{"专业报告", "数据分析", "行业洞察"}
	}

	moods, _ := analysis["suggested_moods"].([]interface{})
	moodList := make([]string, 0, len(moods))
	for _, m := range moods {
		if str, ok := m.(string); ok {
			moodList = append(moodList, str)
		}
	}

	styles, _ := analysis["suggested_styles"].([]interface{})
	styleList := make([]string, 0, len(styles))
	for _, s := range styles {
		if str, ok := s.(string); ok {
			styleList = append(styleList, str)
		}
	}

	hexes := make([]string, 0)
	if dc, ok := analysis["dominant_colors"].([]interface{}); ok {
		for i, ci := range dc {
			if i >= 5 {
				break
			}
			if cm, ok := ci.(map[string]interface{}); ok {
				if h, ok := cm["hex"].(string); ok {
					hexes = append(hexes, h)
				}
			}
		}
	}

	brightness, _ := analysis["brightness_label"].(string)

	return map[string]interface{}{
		"type":             "image",
		"title":            strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		"analysis":         analysis,
		"generated_prompt": prompt,
		"suggested_topics": suggestedTopics,
		"visual_context": map[string]interface{}{
			"colors":     hexes,
			"brightness": brightness,
			"moods":      moodList,
			"styles":     styleList,
		},
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────

// generateGuidance generates guidance text for the user.
func generateGuidance(ctype string, content map[string]interface{}) string {
	guides := map[string]string{
		"docx": "📄 检测到 Word 文档。已提取章节结构，建议逐章确认标题，\n然后选择配色方案和字体，点击「生成完整演示」。",
		"markdown": "📝 检测到 Markdown 文件。已解析标题层级，\n建议补充每节的核心要点（每节 3-5 个 bullet），然后生成。",
		"pptx": "📊 检测到已有 PPTX 文件。已提取每页文本和图片信息，\n可以选择保留原始风格或应用 via54 模板重新设计。",
		"image": "🖼️ 检测到图片。已分析视觉特征（色彩/亮度/情绪），\n建议基于「",
		"text": "📃 检测到纯文本。已自动分段，建议标注各段标题，\n并补充关键数据和图片需求。",
	}
	guide, ok := guides[ctype]
	if !ok {
		return "文件已分析完成，请确认内容并继续。"
	}
	if ctype == "image" {
		topics, _ := content["suggested_topics"].([]interface{})
		topicStr := "方案"
		if len(topics) > 0 {
			if t, ok := topics[0].(string); ok {
				topicStr = t
			}
		}
		guide += topicStr + "」方向补充文案。"
	}
	return guide
}

// readZipXMLContent reads a specific file from a ZIP archive at the given path.
func readZipXMLContent(zipPath, xmlPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	return readZipFileOpen(r, xmlPath)
}

// readZipFileOpen reads a named file from an opened zip.ReadCloser.
func readZipFileOpen(r *zip.ReadCloser, name string) (string, error) {
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			data := make([]byte, f.UncompressedSize64)
			_, err = rc.Read(data)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}
	}
	return "", fmt.Errorf("file not found in zip: %s", name)
}

// stripXMLNS removes XML namespace prefixes from element and attribute names,
// and removes xmlns declarations, so encoding/xml can parse the result.
func stripXMLNS(xmlContent string) string {
	// Remove xmlns declarations
	re := regexp.MustCompile(`\s+xmlns:[a-zA-Z][a-zA-Z0-9]*="[^"]*"`)
	result := re.ReplaceAllString(xmlContent, "")

	// Strip namespace prefixes from opening tags: <ns:tag → <tag
	re = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*):`)
	result = re.ReplaceAllString(result, "<")

	// Strip namespace prefixes from closing tags: </ns:tag → </tag
	re = regexp.MustCompile(`</([a-zA-Z][a-zA-Z0-9]*):`)
	result = re.ReplaceAllString(result, "</")

	// Strip namespace prefixes from self-closing tags
	re = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*):([a-zA-Z][a-zA-Z0-9]*)/>`)
	result = re.ReplaceAllString(result, "<$2/>")

	// Strip namespace prefixes from attribute names: ns:attr= → attr=
	re = regexp.MustCompile(`([\s])([a-zA-Z][a-zA-Z0-9]*):([a-zA-Z][a-zA-Z0-9]*)=`)
	result = re.ReplaceAllString(result, "$1$3=")

	return result
}

// splitSentences splits text into sentences by common Chinese/English delimiters.
func splitSentences(text string) []string {
	if text == "" {
		return nil
	}
	re := regexp.MustCompile(`[。！？\n]`)
	parts := re.Split(text, -1)
	sentences := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 10 {
			sentences = append(sentences, p)
		}
	}
	return sentences
}

// parseRels parses relationship XML and returns a map of rId → target path.
func parseRels(relsXML string) map[string]string {
	result := make(map[string]string)
	re := regexp.MustCompile(`<Relationship\s+Id="([^"]+)"\s+Type="[^"]*"\s+Target="([^"]+)"`)
	matches := re.FindAllStringSubmatch(relsXML, -1)
	for _, m := range matches {
		if len(m) >= 3 {
			target := m[2]
			// Make target relative to the right directory
			if strings.HasPrefix(target, "slides/") {
				target = "ppt/" + target
			}
			result[m[1]] = target
		}
	}
	return result
}
