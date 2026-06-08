// via54Design — 导出模块测试
// SPDX-License-Identifier: AGPL-3.0-only

package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportJSON 验证 JSON 导出
func TestExportJSON(t *testing.T) {
	scenes := []SceneData{
		{Title: "测试1", Voiceover: "旁白1", Body: "正文1", Mood: "calm"},
		{Title: "测试2", Voiceover: "旁白2", Body: "正文2", Mood: "excited"},
	}
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.json")

	if err := ExportJSON(scenes, outPath); err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}

	// 验证内容 (ExportJSON 包了 wrapper { version, total_scenes, scenes })
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var wrapped struct {
		Version    string      `json:"version"`
		TotalScenes int        `json:"total_scenes"`
		Scenes      []SceneData `json:"scenes"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		t.Fatal(err)
	}
	if len(wrapped.Scenes) != 2 {
		t.Errorf("expected 2 scenes, got %d", len(wrapped.Scenes))
	}
	if wrapped.Scenes[0].Title != "测试1" {
		t.Errorf("scene[0].Title = %q, want %q", wrapped.Scenes[0].Title, "测试1")
	}
}

// TestExportJSONEmptyScenes 验证空场景
func TestExportJSONEmptyScenes(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "empty.json")
	if err := ExportJSON(nil, outPath); err != nil {
		t.Errorf("ExportJSON with nil scenes should succeed, got: %v", err)
	}
}

// TestExportJSONInvalidPath 验证错误路径
func TestExportJSONInvalidPath(t *testing.T) {
	scenes := []SceneData{{Title: "test"}}
	// 不存在的目录
	if err := ExportJSON(scenes, "/nonexistent/dir/file.json"); err == nil {
		t.Error("ExportJSON to invalid path should fail")
	}
}

// TestExportMarkdown 验证 Markdown 导出 (Marp 格式)
func TestExportMarkdown(t *testing.T) {
	scenes := []SceneData{
		{Title: "标题A", Voiceover: "旁白A", Body: "正文A", Mood: "calm"},
		{Title: "标题B", Voiceover: "旁白B", Body: "正文B", Mood: "excited"},
	}
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "story.md")

	if err := ExportMarkdown(scenes, "测试标题", "test-author", outPath); err != nil {
		t.Fatalf("ExportMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	// 验证 Marp frontmatter
	if !strings.Contains(content, "marp: true") {
		t.Error("markdown should contain Marp frontmatter")
	}
	if !strings.Contains(content, "title: 测试标题") {
		t.Error("markdown should contain title")
	}
	if !strings.Contains(content, "author: test-author") {
		t.Error("markdown should contain author")
	}
	if !strings.Contains(content, "标题A") {
		t.Error("markdown should contain scene 1 title")
	}
}

// TestPPTXSlideFromBeat 验证 PPTX slide 转换
func TestPPTXSlideFromBeat(t *testing.T) {
	slide := PPTXSlideFromBeat("act1", "测试旁白", "calm", 1, 5)
	if slide.Title == "" {
		t.Error("Title should not be empty")
	}
	if slide.Subtitle == "" {
		t.Error("Subtitle should not be empty")
	}
	if len(slide.Body) == 0 {
		t.Error("Body should have at least 1 element")
	}
}

// TestPPTXSlideFromBeatConsistent 验证 PPTX 转换稳定
func TestPPTXSlideFromBeatConsistent(t *testing.T) {
	// 同样输入应产生同样输出
	s1 := PPTXSlideFromBeat("act1", "测试旁白", "calm", 1, 5)
	s2 := PPTXSlideFromBeat("act1", "测试旁白", "calm", 1, 5)
	if s1.Title != s2.Title {
		t.Errorf("Same input should produce same title, got %q vs %q", s1.Title, s2.Title)
	}
}

// TestCheckDependencies 验证依赖检查 (短超时)
func TestCheckDependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CheckDependencies in short mode")
	}
	missing := CheckDependencies()
	t.Logf("Missing dependencies: %v", missing)
	// 不做严格断言
}

// BenchmarkExportJSON 性能基准
func BenchmarkExportJSON(b *testing.B) {
	scenes := make([]SceneData, 10)
	for i := range scenes {
		scenes[i] = SceneData{
			Title:     "场景",
			Voiceover: "旁白内容",
			Body:      "正文内容",
			Mood:      "calm",
		}
	}
	tmpDir := b.TempDir()
	outPath := filepath.Join(tmpDir, "bench.json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := ExportJSON(scenes, outPath); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPPTXSlideFromBeat 性能基准
func BenchmarkPPTXSlideFromBeat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PPTXSlideFromBeat("act1", "旁白内容", "calm", 1, 5)
	}
}
