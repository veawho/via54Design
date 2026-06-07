// via54Design — 设计质量门禁 v2
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package quality

import (
	"fmt"
	"regexp"
	"strings"
)

func (c *Checker) checkLayout() []Issue {
	var issues []Issue
	layouts := []string{"cover", "statement", "kpi-tower", "loop-diagram", "duo-compare",
		"image-hero", "closing", "section-title", "stats-grid", "timeline",
		"comparison", "process-flow", "team-grid", "quote-hero", "data-table",
		"map-location", "gallery-grid", "pricing-card", "faq-accordion", "cta-banner",
		"contact-form", "footer-closing"}
	for _, layout := range layouts {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)class\s*=\s*["']layout-%s["']|data-layout\s*=\s*["']%s["']`, layout, layout))
		if re.MatchString(c.html) {
			issues = append(issues, Issue{"info", "layout", fmt.Sprintf("Layout '%s' detected", layout)})
		}
	}
	return issues
}

func (c *Checker) checkColorCompliance() []Issue {
	var issues []Issue
	if !strings.Contains(c.html, "--") {
		issues = append(issues, Issue{"warning", "color", "No CSS custom properties (--var) found"})
	}
	hexRe := regexp.MustCompile(`#[0-9a-fA-F]{6}`)
	hexMatches := hexRe.FindAllString(c.html, -1)
	if len(hexMatches) > 10 {
		issues = append(issues, Issue{"warning", "color", fmt.Sprintf("%d hardcoded hex colors", len(hexMatches))})
	}
	return issues
}

func (c *Checker) checkTypography() []Issue {
	var issues []Issue
	if !strings.Contains(c.html, "font-family") {
		issues = append(issues, Issue{"warning", "typography", "No font-family declaration"})
	}
	if !strings.Contains(c.html, "line-height") {
		issues = append(issues, Issue{"warning", "typography", "No line-height set"})
	}
	return issues
}

func (c *Checker) checkResponsive() []Issue {
	var issues []Issue
	mediaRe := regexp.MustCompile(`(?i)@media`)
	if n := len(mediaRe.FindAllString(c.html, -1)); n == 0 {
		issues = append(issues, Issue{"warning", "responsive", "No @media queries (not responsive)"})
	}
	if !strings.Contains(c.html, "viewport") {
		issues = append(issues, Issue{"error", "responsive", "Missing viewport meta tag"})
	}
	return issues
}

func (c *Checker) checkAccessibility() []Issue {
	var issues []Issue
	imgRe := regexp.MustCompile(`(?i)<img[^>]+>`)
	imgs := imgRe.FindAllString(c.html, -1)
	noAlt := 0
	for _, img := range imgs {
		if !strings.Contains(img, "alt=") { noAlt++ }
	}
	if noAlt > 0 {
		issues = append(issues, Issue{"warning", "a11y", fmt.Sprintf("%d <img> missing alt text", noAlt)})
	}
	return issues
}

func (c *Checker) RunAllV2() *Report {
	r := c.RunAll()
	r.Issues = append(r.Issues, c.checkLayout()...)
	r.Issues = append(r.Issues, c.checkColorCompliance()...)
	r.Issues = append(r.Issues, c.checkTypography()...)
	r.Issues = append(r.Issues, c.checkResponsive()...)
	r.Issues = append(r.Issues, c.checkAccessibility()...)
	summary := map[string]int{"error": 0, "warning": 0, "info": 0}
	for _, iss := range r.Issues { summary[iss.Severity]++ }
	r.Summary = summary
	verdict := "PASS"
	if summary["error"] > 0 { verdict = "FAIL"
	} else if summary["warning"] > 2 { verdict = "WARNING" }
	r.Verdict = verdict
	return r
}
