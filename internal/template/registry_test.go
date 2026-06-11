// via54Design — 模板注册表 + 引擎测试 (补充)
// SPDX-License-Identifier: AGPL-3.0-only

package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseDir(t *testing.T) string {
	t.Helper()
	return "../.."
}

// TestNewRegistry_Success 验证 registry 加载
func TestNewRegistry_Success(t *testing.T) {
	r, err := NewRegistry(baseDir(t))
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	if r.Data == nil {
		t.Error("registry data nil")
	}
	if len(r.Data.Layouts) == 0 {
		t.Error("no layouts in registry")
	}
	if len(r.Data.ColorSchemes) == 0 {
		t.Error("no color schemes in registry")
	}
	if len(r.Data.Typography) == 0 {
		t.Error("no typography in registry")
	}
}

// TestNewRegistry_MissingDir 验证错误 baseDir
func TestNewRegistry_MissingDir(t *testing.T) {
	_, err := NewRegistry("/nonexistent/path/abc")
	if err == nil {
		t.Error("expected error for missing dir")
	}
}

// TestResolveLayout_Hit 验证 layout 解析
func TestResolveLayout_Hit(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	// 找一个真实 layout
	if len(r.Data.Layouts) == 0 {
		t.Skip("no layouts in registry")
	}
	id := r.Data.Layouts[0].ID
	path, err := r.ResolveLayout(id)
	if err != nil {
		t.Errorf("ResolveLayout(%q) failed: %v", id, err)
	}
	if path == "" {
		t.Error("resolved path empty")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("resolved file not exists: %s", path)
	}
}

// TestResolveLayout_Miss 验证 layout not found
func TestResolveLayout_Miss(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	_, err := r.ResolveLayout("nonexistent-layout-xyz")
	if err == nil {
		t.Error("expected error for nonexistent layout")
	}
}

// TestResolveColorScheme_Hit 验证配色解析
func TestResolveColorScheme_Hit(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	if len(r.Data.ColorSchemes) == 0 {
		t.Skip("no color schemes")
	}
	id := r.Data.ColorSchemes[0].ID
	_, err := r.ResolveColorScheme(id)
	if err != nil {
		t.Errorf("ResolveColorScheme(%q) failed: %v", id, err)
	}
}

// TestResolveTypography_Hit 验证字体解析
func TestResolveTypography_Hit(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	if len(r.Data.Typography) == 0 {
		t.Skip("no typography")
	}
	id := r.Data.Typography[0].ID
	_, err := r.ResolveTypography(id)
	if err != nil {
		t.Errorf("ResolveTypography(%q) failed: %v", id, err)
	}
}

// TestResolveNarratology_Hit 验证叙事模型解析
func TestResolveNarratology_Hit(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	_, err := r.ResolveNarratology("three-act")
	if err != nil {
		t.Errorf("ResolveNarratology(three-act) failed: %v", err)
	}
}

// TestListAll 验证 ListAll
func TestListAll(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	all := r.ListAll()
	for _, key := range []string{"layouts", "color_schemes", "typography", "narratology"} {
		if _, ok := all[key]; !ok {
			t.Errorf("ListAll missing %q", key)
		}
	}
}

// TestCompose_RealLayout 验证完整 compose 流程
func TestCompose_RealLayout(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	eng := &Engine{Registry: r, BaseDir: baseDir(t)}
	if len(r.Data.Layouts) == 0 || len(r.Data.ColorSchemes) == 0 || len(r.Data.Typography) == 0 {
		t.Skip("registry incomplete")
	}
	res, err := eng.Compose(
		r.Data.Layouts[0].ID,
		r.Data.ColorSchemes[0].ID,
		r.Data.Typography[0].ID,
		"测试标题",
	)
	if err != nil {
		t.Fatalf("Compose failed: %v", err)
	}
	if res.HTML == "" {
		t.Error("HTML empty")
	}
	if !strings.Contains(res.HTML, "<!DOCTYPE") {
		t.Error("HTML missing DOCTYPE")
	}
	if !strings.Contains(res.HTML, "测试标题") {
		t.Error("HTML missing title")
	}
}

// TestCompose_BadLayout 验证错误 layout 返回 error
func TestCompose_BadLayout(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	eng := &Engine{Registry: r, BaseDir: baseDir(t)}
	_, err := eng.Compose("nonexistent", r.Data.ColorSchemes[0].ID, r.Data.Typography[0].ID, "title")
	if err == nil {
		t.Error("expected error for bad layout")
	}
}

// TestSortedKeys_Deterministic 验证 sortedKeys (map 遍历确定性)
func TestSortedKeys_Deterministic(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	a := sortedKeys(m)
	b := sortedKeys(m)
	if len(a) != 3 {
		t.Errorf("sortedKeys length = %d, want 3", len(a))
	}
	// 多次调用顺序一致
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("sortedKeys not deterministic at %d", i)
		}
	}
	// 应该是有序的
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			t.Errorf("not sorted: %v", a)
		}
	}
}

// TestSortedKeys_Empty 验证空 map
func TestSortedKeys_Empty(t *testing.T) {
	got := sortedKeys(map[string]int{})
	if len(got) != 0 {
		t.Errorf("empty map should return empty slice, got %v", got)
	}
}

// TestResolve_ResolvesRelativePath 验证相对路径解析
func TestResolve_ResolvesRelativePath(t *testing.T) {
	r, _ := NewRegistry(baseDir(t))
	// 注册表中的 file 字段应该是相对路径
	for _, e := range r.Data.Layouts {
		if filepath.IsAbs(e.File) {
			t.Errorf("layout %q file should be relative, got %s", e.ID, e.File)
		}
	}
}
