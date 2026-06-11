// via54Design — 质量门禁测试
// SPDX-License-Identifier: AGPL-3.0-only

package quality

import (
	"strings"
	"testing"
)

const validHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Test</title>
  <style>body { font-family: sans-serif; line-height: 1.5; color: #333; }</style>
</head>
<body>
  <h1>标题</h1>
  <p>内容段落</p>
  <img src="test.jpg" alt="测试图">
  <button class="cta" :focus-visible>按钮</button>
</body>
</html>`

// TestCheckHTMLStructure_Valid 验证合法 HTML 不报错
func TestCheckHTMLStructure_Valid(t *testing.T) {
	rep := CheckHTML(validHTML)
	for _, iss := range rep.Issues {
		if iss.Severity == "error" {
			t.Errorf("valid HTML produced error: %s — %s", iss.Category, iss.Message)
		}
	}
}

// TestCheckHTMLStructure_MissingDoctype 验证缺 DOCTYPE
func TestCheckHTMLStructure_MissingDoctype(t *testing.T) {
	rep := CheckHTML("<html><head></head><body></body></html>")
	found := false
	for _, iss := range rep.Issues {
		if iss.Category == "html" && strings.Contains(iss.Message, "DOCTYPE") {
			found = true
		}
	}
	if !found {
		t.Error("expected DOCTYPE error, not found")
	}
}

// TestCheckCSS_BraceMismatch 验证 CSS 大括号不匹配
func TestCheckCSS_BraceMismatch(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>body { color: red; </style></head><body></body></html>`
	rep := CheckHTML(html)
	found := false
	for _, iss := range rep.Issues {
		if iss.Category == "css" && strings.Contains(iss.Message, "brace") {
			found = true
		}
	}
	if !found {
		t.Error("expected CSS brace mismatch error, not found")
	}
}

// TestCheckCSS_ImportantDetected 验证 !important 警告
func TestCheckCSS_ImportantDetected(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>a { color: red !important; }</style></head><body></body></html>`
	rep := CheckHTML(html)
	found := false
	for _, iss := range rep.Issues {
		if iss.Category == "css" && strings.Contains(iss.Message, "!important") {
			found = true
		}
	}
	if !found {
		t.Error("expected !important warning, not found")
	}
}

// TestCheckContent_EmptySection 验证空 section
func TestCheckContent_EmptySection(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>body{}</style></head><body><section></section><p></p></body></html>`
	rep := CheckHTML(html)
	emptySection := false
	emptyP := false
	for _, iss := range rep.Issues {
		if iss.Category == "content" {
			if strings.Contains(iss.Message, "empty <section>") {
				emptySection = true
			}
			if strings.Contains(iss.Message, "empty <p>") {
				emptyP = true
			}
		}
	}
	if !emptySection {
		t.Error("empty <section> not detected")
	}
	if !emptyP {
		t.Error("empty <p> not detected")
	}
}

// TestCheckAntiSlop_Emoji 验证 emoji 检测
func TestCheckAntiSlop_Emoji(t *testing.T) {
	html := `<!DOCTYPE html><html><head></head><body>⭐⭐⭐ 标题</body></html>`
	rep := CheckHTML(html)
	found := false
	for _, iss := range rep.Issues {
		if iss.Category == "anti-slop" && strings.Contains(iss.Message, "emoji") {
			found = true
		}
	}
	if !found {
		t.Error("emoji not detected by anti-slop")
	}
}

// TestCheckAntiSlop_GitHubDark 验证 GitHub-dark 色检测
func TestCheckAntiSlop_GitHubDark(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>body{background:#0D1117;}</style></head><body></body></html>`
	rep := CheckHTML(html)
	found := false
	for _, iss := range rep.Issues {
		if iss.Category == "anti-slop" && strings.Contains(iss.Message, "0D1117") {
			found = true
		}
	}
	if !found {
		t.Error("GitHub-dark color not flagged")
	}
}

// TestCheckV2_Responsive 验证 V2 响应式检测
func TestCheckV2_Responsive(t *testing.T) {
	html := `<!DOCTYPE html><html><head><meta name="viewport"></head><body></body></html>`
	c := New(html)
	issues := c.checkResponsive()
	if len(issues) == 0 {
		t.Error("expected responsive warning, got 0")
	}
	hasMedia := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "@media") {
			hasMedia = true
		}
	}
	if !hasMedia {
		t.Error("expected @media warning")
	}
}

// TestCheckV2_ColorCompliance_NoVars 验证 v2 颜色合规（无 CSS 变量）
func TestCheckV2_ColorCompliance_NoVars(t *testing.T) {
	html := `<!DOCTYPE html><html><head></head><body><p>test</p></body></html>`
	c := New(html)
	issues := c.checkColorCompliance()
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "CSS custom properties") {
			found = true
		}
	}
	if !found {
		t.Error("expected CSS variables warning, not found")
	}
}

// TestCheckV2_ColorCompliance_ManyHex 验证硬编码 hex 过多
func TestCheckV2_ColorCompliance_ManyHex(t *testing.T) {
	hex := strings.Repeat("#FFAA00 ", 15)
	html := `<!DOCTYPE html><html><head><style>:root{` + hex + `}</style></head><body></body></html>`
	c := New(html)
	issues := c.checkColorCompliance()
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "hardcoded hex") {
			found = true
		}
	}
	if !found {
		t.Error("expected hardcoded hex warning, not found")
	}
}

// TestCheckV2_Accessibility_MissingAlt 验证缺失 alt 检测
func TestCheckV2_Accessibility_MissingAlt(t *testing.T) {
	html := `<!DOCTYPE html><html><head></head><body><img src="x.jpg"><img src="y.jpg" alt="y"></body></html>`
	c := New(html)
	issues := c.checkAccessibility()
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "missing alt") {
			found = true
		}
	}
	if !found {
		t.Error("expected missing alt warning, not found")
	}
}

// TestCheckV2_AntiCliche 验证反陈词滥调
func TestCheckV2_AntiCliche(t *testing.T) {
	html := `<!DOCTYPE html><html><head></head><body><p>Our cutting-edge solution leverages synergy to disrupt the industry.</p></body></html>`
	c := New(html)
	issues := c.checkAntiCliche()
	if len(issues) < 3 {
		t.Errorf("expected multiple cliche issues, got %d", len(issues))
	}
}

// TestReport_VerdictPass 验证 PASS verdict
func TestReport_VerdictPass(t *testing.T) {
	rep := CheckHTML(validHTML)
	if rep.Verdict == "FAIL" {
		t.Errorf("valid HTML should not FAIL: %+v", rep.Issues)
	}
}

// TestReport_VerdictFail 验证 FAIL verdict
func TestReport_VerdictFail(t *testing.T) {
	rep := CheckHTML("not even html")
	if rep.Verdict != "FAIL" {
		t.Errorf("garbage should FAIL, got %s", rep.Verdict)
	}
	if rep.Summary["error"] == 0 {
		t.Error("expected error summary > 0")
	}
}

// TestReport_HasHTMLSize 验证 HTMLSize 字段
func TestReport_HasHTMLSize(t *testing.T) {
	rep := CheckHTML(validHTML)
	if rep.HTMLSize != len(validHTML) {
		t.Errorf("HTMLSize = %d, want %d", rep.HTMLSize, len(validHTML))
	}
}

// TestReport_CSSBlocksCount 验证 CSS 块数
func TestReport_CSSBlocksCount(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>a{}</style><style>b{}</style></head><body></body></html>`
	rep := CheckHTML(html)
	if rep.CSSBlocks != 2 {
		t.Errorf("CSSBlocks = %d, want 2", rep.CSSBlocks)
	}
}

// TestCheckV2_Typography_NoFamily 验证缺 font-family
func TestCheckV2_Typography_NoFamily(t *testing.T) {
	html := `<!DOCTYPE html><html><head><style>body{color:red}</style></head><body></body></html>`
	c := New(html)
	issues := c.checkTypography()
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "font-family") {
			found = true
		}
	}
	if !found {
		t.Error("expected font-family warning, not found")
	}
}
