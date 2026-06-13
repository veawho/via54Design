// via54Design — 模板引擎测试
// SPDX-License-Identifier: AGPL-3.0-only

package template

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestNewEngine 验证 NewEngine 成功创建
func TestNewEngine(t *testing.T) {
	// 用仓库根目录作为 baseDir (因为仓库根有 templates/ 目录)
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatalf("NewEngine failed: %v", err)
	}
	if eng.Registry == nil {
		t.Error("Engine.Registry should not be nil")
	}
}

// TestNewEngineNotFound 验证 NewEngine 错误处理
func TestNewEngineNotFound(t *testing.T) {
	_, err := NewEngine("/nonexistent/path/should/not/exist")
	if err == nil {
		t.Error("NewEngine should fail for non-existent path")
	}
}

// TestNewRegistry 验证 NewRegistry 成功创建
func TestNewRegistry(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	if reg.Data == nil {
		t.Error("Registry.Data should not be nil")
	}
}

// TestRegistryListAll 验证 ListAll 返回所有分类
func TestRegistryListAll(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatal(err)
	}
	all := reg.ListAll()
	if len(all) == 0 {
		t.Error("Expected at least one category")
	}
	if entries, ok := all["layouts"]; !ok || len(entries) == 0 {
		t.Error("Expected at least one layout")
	}
	if entries, ok := all["color_schemes"]; !ok || len(entries) == 0 {
		t.Error("Expected at least one color scheme")
	}
	if entries, ok := all["typography"]; !ok || len(entries) == 0 {
		t.Error("Expected at least one typography")
	}
	t.Logf("Categories: %d", len(all))
	for cat, entries := range all {
		t.Logf("  %s: %d entries", cat, len(entries))
	}
}

// TestRegistryListNarratology 验证列出所有 narrative models
func TestRegistryListNarratology(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatal(err)
	}
	narrs := reg.ListNarratology()
	if len(narrs) == 0 {
		t.Error("Expected at least one narrative model")
	}
}

// TestResolveLayout 验证根据 ID 查找 layout
func TestResolveLayout(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatal(err)
	}
	// 测试存在的 layout
	path, err := reg.ResolveLayout("hero-split-16-9")
	if err != nil {
		t.Errorf("ResolveLayout for existing layout should succeed: %v", err)
	}
	if path == "" {
		t.Error("ResolveLayout should return a path")
	}
	t.Logf("DEBUG: existing layout path = %q", path)
	// 测试不存在的 layout
	badPath, err2 := reg.ResolveLayout("NONEXISTENT-LAYOUT-12345")
	t.Logf("DEBUG: bad layout path = %q, err = %v", badPath, err2)
	if err2 == nil {
		t.Error("ResolveLayout for non-existent layout should fail")
	}
}

// TestResolveColorScheme 验证根据 ID 查找 color
func TestResolveColorScheme(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatal(err)
	}
	_, err = reg.ResolveColorScheme("ink-wash")
	if err != nil {
		t.Errorf("ResolveColorScheme for existing color should succeed: %v", err)
	}
	_, err = reg.ResolveColorScheme("NONEXISTENT-COLOR-12345")
	if err == nil {
		t.Error("ResolveColorScheme for non-existent color should fail")
	}
}

// TestResolveTypography 验证根据 ID 查找 font
func TestResolveTypography(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		t.Fatal(err)
	}
	_, err = reg.ResolveTypography("ming-hei-editorial")
	if err != nil {
		t.Errorf("ResolveTypography for existing font should succeed: %v", err)
	}
	_, err = reg.ResolveTypography("NONEXISTENT-FONT-12345")
	if err == nil {
		t.Error("ResolveTypography for non-existent font should fail")
	}
}

// TestComposeHTML 验证 Compose 端到端
func TestComposeHTML(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	result, err := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", "测试标题")
	if err != nil {
		t.Fatalf("Compose failed: %v", err)
	}
	if result == nil {
		t.Fatal("Compose returned nil result")
	}
	if !strings.Contains(result.HTML, "<!DOCTYPE html>") {
		t.Error("Generated HTML should contain DOCTYPE")
	}
	if !strings.Contains(result.HTML, "测试标题") {
		t.Error("Generated HTML should contain the title")
	}
	if !strings.Contains(result.HTML, "layout-hero-split") {
		t.Error("Generated HTML should have layout-hero-split class")
	}
}

// TestComposeHTMLInvalidLayout 验证错误布局名
func TestComposeHTMLInvalidLayout(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	_, err = eng.Compose("NONEXISTENT", "ink-wash", "ming-hei-editorial", "测试")
	if err == nil {
		t.Error("Compose with invalid layout should fail")
	}
}

// TestXSSProtection 验证 XSS 转义 (P1 安全修复)
func TestXSSProtection(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	xssPayload := "<script>alert(1)</script>"
	result, err := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", xssPayload)
	if err != nil {
		t.Fatal(err)
	}
	// XSS payload 应被转义
	if strings.Contains(result.HTML, "<script>alert(1)</script>") {
		t.Error("XSS payload should be HTML-escaped, but found unescaped version")
	}
	if !strings.Contains(result.HTML, "&lt;script&gt;") {
		t.Error("XSS payload should be escaped as &lt;script&gt;")
	}
}

// TestXSSOnErrorAttribute 验证 attribute 注入也被转义
func TestXSSOnErrorAttribute(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	xssPayload := `" onerror="alert(1)`
	result, err := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", xssPayload)
	if err != nil {
		t.Fatal(err)
	}
	// 检查未转义
	if strings.Contains(result.HTML, `" onerror=`) {
		t.Error("attribute injection should be escaped")
	}
	// 应有 HTML 实体 (&#34; 或 &quot;)
	if !strings.Contains(result.HTML, "&#34;") && !strings.Contains(result.HTML, "&quot;") {
		t.Error("attribute injection should be escaped to HTML entity")
	}
}

// TestComposeBentoLayout 验证 bento 布局
func TestComposeBentoLayout(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	result, err := eng.Compose("bento-grid-2x2", "candy-duolingo", "hei-modern", "Bento")
	if err != nil {
		t.Fatalf("Compose bento failed: %v", err)
	}
	if !strings.Contains(result.HTML, "layout-bento") {
		t.Error("Bento HTML should have layout-bento class")
	}
}

// TestComposeGalleryLayout 验证 gallery 布局
func TestComposeGalleryLayout(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	result, err := eng.Compose("gallery-waterfall", "rinpa-gold", "elegant-didone", "Gallery")
	if err != nil {
		t.Fatalf("Compose gallery failed: %v", err)
	}
	if !strings.Contains(result.HTML, "layout-gallery") {
		t.Error("Gallery HTML should have layout-gallery class")
	}
}

// TestDeterministicOutput 验证同样输入产生同样输出
func TestDeterministicOutput(t *testing.T) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		t.Fatal(err)
	}
	r1, _ := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", "Determinism")
	r2, _ := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", "Determinism")
	if r1.HTML != r2.HTML {
		t.Error("Same input should produce same output (deterministic)")
	}
}

// BenchmarkComposeHero 性能基准
func BenchmarkComposeHero(b *testing.B) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eng.Compose("hero-split-16-9", "ink-wash", "ming-hei-editorial", "基准测试")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkComposeBento 性能基准 (Bento)
func BenchmarkComposeBento(b *testing.B) {
	wd, _ := filepath.Abs("../..")
	eng, err := NewEngine(wd)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eng.Compose("bento-grid-2x2", "candy-duolingo", "hei-modern", "Bento")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveLayout 解析性能
func BenchmarkResolveLayout(b *testing.B) {
	wd, _ := filepath.Abs("../..")
	reg, err := NewRegistry(wd)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := reg.ResolveLayout("hero-split-16-9")
		if err != nil {
			b.Fatal(err)
		}
	}
}
