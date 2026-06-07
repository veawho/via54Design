package quality

import (
	"fmt"
	"regexp"
	"strings"
)

type Issue struct {
	Severity string
	Category string
	Message  string
}

type Report struct {
	Verdict    string
	TotalScore int
	Issues     []Issue
	Summary    map[string]int
	HTMLSize   int
	CSSBlocks  int
	TotalLines int
}

type Checker struct {
	html string
}

func New(htmlContent string) *Checker {
	return &Checker{html: htmlContent}
}

func (c *Checker) RunAll() *Report {
	allIssues := []Issue{}

	allIssues = append(allIssues, c.checkHTMLStructure()...)
	allIssues = append(allIssues, c.checkCSS()...)
	allIssues = append(allIssues, c.checkContent()...)
	allIssues = append(allIssues, c.checkAntiSlop()...)

	summary := map[string]int{"error": 0, "warning": 0, "info": 0}
	for _, iss := range allIssues {
		summary[iss.Severity]++
	}

	verdict := "PASS"
	if summary["error"] > 0 {
		verdict = "FAIL"
	} else if summary["warning"] > 0 {
		verdict = "WARNING"
	}

	cssRe := regexp.MustCompile("(?i)<style[^>]*>(.*?)</style>")

	return &Report{
		Verdict:    verdict,
		Issues:     allIssues,
		Summary:    summary,
		HTMLSize:   len(c.html),
		CSSBlocks:  len(cssRe.FindAllString(c.html, -1)),
		TotalLines: strings.Count(c.html, "\n") + 1,
	}
}

func (c *Checker) checkHTMLStructure() []Issue {
	var issues []Issue
	u := strings.ToUpper(c.html)

	if !strings.Contains(u, "<!DOCTYPE HTML") {
		issues = append(issues, Issue{"error", "html", "Missing DOCTYPE"})
	}
	if !strings.Contains(u, "<HTML") {
		issues = append(issues, Issue{"error", "html", "Missing <html>"})
	}
	if !strings.Contains(c.html, "<head>") {
		issues = append(issues, Issue{"error", "html", "Missing <head>"})
	}
	if !strings.Contains(c.html, "<body") {
		issues = append(issues, Issue{"error", "html", "Missing <body>"})
	}
	if !strings.Contains(u, "</HTML>") {
		issues = append(issues, Issue{"warning", "html", "Missing </html>"})
	}
	if !strings.Contains(c.html[:min(len(c.html), 2000)], "charset") {
		issues = append(issues, Issue{"warning", "html", "Missing charset"})
	}
	if !strings.Contains(c.html[:min(len(c.html), 2000)], "viewport") {
		issues = append(issues, Issue{"warning", "html", "Missing viewport"})
	}
	return issues
}

func (c *Checker) checkCSS() []Issue {
	var issues []Issue
	re := regexp.MustCompile("(?i)<style[^>]*>(.*?)</style>")
	matches := re.FindAllStringSubmatch(c.html, -1)

	for i, match := range matches {
		block := match[1]
		opens := strings.Count(block, "{")
		closes := strings.Count(block, "}")
		if opens != closes {
			issues = append(issues, Issue{"error", "css",
				fmt.Sprintf("CSS block #%d: brace mismatch (%d opens, %d closes)", i+1, opens, closes)})
		}
		if strings.Contains(block, "!important") {
			issues = append(issues, Issue{"warning", "css",
				fmt.Sprintf("CSS block #%d: !important found", i+1)})
		}
	}
	return issues
}

func (c *Checker) checkContent() []Issue {
	var issues []Issue

	checks := []struct{
		pattern string
		desc string
		severity string
	}{
		{"<section[^>]*>\\s*</section>", "empty <section>", "warning"},
		{"<div[^>]*>\\s*</div>", "empty <div>", "warning"},
		{"<p[^>]*>\\s*</p>", "empty <p>", "warning"},
		{"<h\\d[^>]*>\\s*</h\\d>", "empty heading", "warning"},
	}
	for _, ch := range checks {
		re := regexp.MustCompile(ch.pattern)
		if n := len(re.FindAllString(c.html, -1)); n > 0 {
			issues = append(issues, Issue{ch.severity, "content",
				fmt.Sprintf("%d x %s", n, ch.desc)})
		}
	}

	imgRe := regexp.MustCompile("(?i)<img[^>]+src=")
	if len(imgRe.FindAllString(c.html, -1)) == 0 &&
	   !strings.Contains(c.html, "background-image") {
		issues = append(issues, Issue{"info", "content", "No <img> tags"})
	}

	return issues
}

func (c *Checker) checkAntiSlop() []Issue {
	var issues []Issue

	// Simple emoji icon check - look for common patterns
	emojiCount := 0
	for _, em := range []string{"⭐", "🔴", "✅", "❌", "🔥", "💡", "🎯", "🚀", "✨", "💪"} {
		emojiCount += strings.Count(c.html, em)
	}
	if emojiCount > 0 {
		issues = append(issues, Issue{"warning", "anti-slop",
			fmt.Sprintf("%d emoji icons found", emojiCount)})
	}

	faceRe := regexp.MustCompile("(?i)<svg[^>]*>.*?(?:face|eye|nose|mouth|person).*?</svg>")
	if n := len(faceRe.FindAllString(c.html, -1)); n > 0 {
		issues = append(issues, Issue{"warning", "anti-slop",
			fmt.Sprintf("%d SVG face drawings found", n)})
	}

	if strings.Contains(c.html, "#0D1117") {
		issues = append(issues, Issue{"info", "anti-slop", "#0D1117 GitHub-dark detected"})
	}

	return issues
}

func CheckHTML(htmlContent string) *Report {
	return New(htmlContent).RunAll()
}
