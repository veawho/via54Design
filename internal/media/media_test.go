// via54Design — 媒体管线测试
// SPDX-License-Identifier: AGPL-3.0-only

package media

import (
	"strings"
	"testing"
)

// TestSanitize 验证文件名清理
// 注: Go 的 \w 默认 ASCII only, 所以中文字符被替换为 _
func TestSanitize(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"hello world", "hello_world"},
		{"abc-def.jpg", "abc-def.jpg"},
		{"file<>:\"|?*", "file_"},
		{"normal_name", "normal_name"},
		{"path/to/file", "path_to_file"},
		{"a b c d", "a_b_c_d"},
	}
	for _, c := range cases {
		got := sanitize(c.in)
		if got != c.want {
			t.Errorf("sanitize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestTruncate 验证字符串截断
func TestTruncate(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{"hello world", 20, "hello world"},
		{"hello world", 5, "hello"},
		{"", 10, ""},
		{"abcdef", 0, ""},
	}
	for _, c := range cases {
		got := truncate(c.in, c.n)
		if got != c.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", c.in, c.n, got, c.want)
		}
	}
}

// TestAddMusic_BadMood 验证错误 mood 返回 error
func TestAddMusic_BadMood(t *testing.T) {
	err := AddMusic("/nonexistent.mp4", "totally-bad-mood", "/tmp/out.mp4")
	if err == nil {
		t.Error("expected error for bad mood")
	}
	if !strings.Contains(err.Error(), "未知 mood") {
		t.Errorf("error should mention 未知 mood: %v", err)
	}
}

// TestFindAsset_Missing 验证资产未找到
func TestFindAsset(t *testing.T) {
	got := findAsset("nonexistent-bgmx.mp3")
	if got != "" {
		t.Errorf("findAsset for missing file should return '', got %q", got)
	}
}

// TestAddSuffix 验证文件后缀
func TestAddSuffix(t *testing.T) {
	cases := []struct {
		in     string
		suffix string
		want   string
	}{
		{"video.mp4", "-with-music", "video-with-music.mp4"},
		{"file.txt", "-v2", "file-v2.txt"},
		{"noext", "-x", "noext-x"},
	}
	for _, c := range cases {
		got := addSuffix(c.in, c.suffix)
		if got != c.want {
			t.Errorf("addSuffix(%q, %q) = %q, want %q", c.in, c.suffix, got, c.want)
		}
	}
}

// TestTrimExt 验证去扩展名
func TestTrimExt(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"video.mp4", "video"},
		{"file.tar.gz", "file.tar"},
		{"noext", "noext"},
	}
	for _, c := range cases {
		got := trimExt(c.in)
		if got != c.want {
			t.Errorf("trimExt(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestBgMMoods_ValidMoods 验证所有 mood 都有对应 BGM 文件名
func TestBgMMoods_ValidMoods(t *testing.T) {
	required := []string{"tech", "ad", "educational", "tutorial"}
	for _, m := range required {
		if _, ok := bgmMoods[m]; !ok {
			t.Errorf("bgmMoods missing %q", m)
		}
		if !strings.HasSuffix(bgmMoods[m], ".mp3") {
			t.Errorf("bgm for %q should be .mp3, got %q", m, bgmMoods[m])
		}
	}
}

// TestFetchImages_EmptyQueries 验证空查询不 panic
func TestFetchImages_EmptyQueries(t *testing.T) {
	dir := t.TempDir()
	results, err := FetchImages(nil, dir, 1)
	if err != nil {
		t.Errorf("empty queries should not error, got: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("empty queries should return 0 results, got %d", len(results))
	}
}

// TestConvertFormats_MissingInput 验证缺输入文件时优雅失败
func TestConvertFormats_MissingInput(t *testing.T) {
	// Skip: 需要 ffmpeg
	t.Skip("requires ffmpeg binary")
}

// TestMixVoiceover_MissingFiles 验证缺文件时优雅失败
func TestMixVoiceover_MissingFiles(t *testing.T) {
	t.Skip("requires ffmpeg binary")
}
