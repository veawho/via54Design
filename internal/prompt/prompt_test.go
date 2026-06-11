// via54Design — 提示词引擎测试
// SPDX-License-Identifier: AGPL-3.0-only

package prompt

import (
	"strings"
	"testing"
)

func baseDir(t *testing.T) string {
	t.Helper()
	return "../.."
}

// TestIsVideoPlatform 验证视频平台判定
func TestIsVideoPlatform(t *testing.T) {
	video := []string{"kling", "veo", "sora", "pika", "seedance"}
	image := []string{"midjourney", "flux", "dalle3", "stable_diffusion", "sd3"}
	for _, p := range video {
		if !isVideoPlatform(p) {
			t.Errorf("%s should be video platform", p)
		}
	}
	for _, p := range image {
		if isVideoPlatform(p) {
			t.Errorf("%s should NOT be video platform", p)
		}
	}
}

// TestApplyWeights 验证权重格式
func TestApplyWeights(t *testing.T) {
	cases := []struct {
		in    string
		w     float64
		want  string
	}{
		{"cat", 1.0, "cat"},
		{"cat", 0, "cat"},
		{"cat", 1.5, "(cat:1.5)"},
		{"cat", 2.0, "(cat:2.0)"},
		{"cat", 0.5, "(cat:0.5)"},
	}
	for _, c := range cases {
		got := applyWeights(c.in, c.w)
		if got != c.want {
			t.Errorf("applyWeights(%q, %v) = %q, want %q", c.in, c.w, got, c.want)
		}
	}
}

// TestMin 验证 min helper
func TestMin(t *testing.T) {
	if min(1, 2) != 1 {
		t.Error("min(1,2) should be 1")
	}
	if min(2, 1) != 1 {
		t.Error("min(2,1) should be 1")
	}
	if min(5, 5) != 5 {
		t.Error("min(5,5) should be 5")
	}
}

// TestGeneratePrompt_Midjourney 验证 midjourney 平台生成
func TestGeneratePrompt_Midjourney(t *testing.T) {
	s, err := GeneratePrompt("a cat on the roof", "midjourney", "", baseDir(t))
	if err != nil {
		t.Fatalf("GeneratePrompt failed: %v", err)
	}
	if s.Platform != "midjourney" {
		t.Errorf("Platform = %q, want midjourney", s.Platform)
	}
	if s.Seed != "a cat on the roof" {
		t.Errorf("Seed = %q", s.Seed)
	}
	if s.FinalPrompt == "" {
		t.Error("FinalPrompt empty")
	}
	// FinalPrompt 应至少包含 midjourney 默认 params (--ar 16:9, --v 6.1 等)
	if !strings.Contains(s.FinalPrompt, "--ar") {
		t.Errorf("FinalPrompt should contain midjourney params, got %q", s.FinalPrompt)
	}
	// midjourney 应该有 negative
	if len(s.Negative) == 0 {
		t.Error("midjourney should have negative prompts")
	}
	if !strings.Contains(s.FinalPrompt, "--no") {
		t.Error("FinalPrompt should include --no for negatives")
	}
	// FinalPrompt 走模板默认 fields (subject 留空等待 LLM 填), seed 存 Seed 字段
	if s.Seed != "a cat on the roof" {
		t.Errorf("Seed = %q, want %q", s.Seed, "a cat on the roof")
	}
}

// TestGeneratePrompt_UnknownPlatform 验证未知平台走 generic 分支
func TestGeneratePrompt_UnknownPlatform(t *testing.T) {
	s, err := GeneratePrompt("test scene", "unknown_platform_xyz", "", baseDir(t))
	if err != nil {
		t.Fatalf("GeneratePrompt with unknown platform should not error, got: %v", err)
	}
	if s.Model != "通用" {
		t.Errorf("unknown platform should fallback to '通用', got %q", s.Model)
	}
	if s.FinalPrompt != "test scene" {
		t.Errorf("generic FinalPrompt should equal seed, got %q", s.FinalPrompt)
	}
}

// TestGeneratePrompt_EmptyScene 验证空场景
func TestGeneratePrompt_EmptyScene(t *testing.T) {
	s, err := GeneratePrompt("", "midjourney", "", baseDir(t))
	if err != nil {
		t.Fatalf("empty scene failed: %v", err)
	}
	if s.FinalPrompt == "" {
		t.Error("FinalPrompt should not be empty even with empty scene")
	}
}

// TestGeneratePrompt_VideoPlatform 验证视频平台（kling）
func TestGeneratePrompt_VideoPlatform(t *testing.T) {
	s, err := GeneratePrompt("a bird flies", "kling", "", baseDir(t))
	if err != nil {
		t.Fatalf("kling prompt failed: %v", err)
	}
	// FinalPrompt 应至少包含 kling 默认 params
	if !strings.Contains(s.FinalPrompt, "--ar") {
		t.Errorf("kling FinalPrompt missing params, got %q", s.FinalPrompt)
	}
	// FinalPrompt 走模板默认 fields, seed 存 Seed 字段
	if s.Seed != "a bird flies" {
		t.Errorf("kling Seed = %q, want %q", s.Seed, "a bird flies")
	}
	// kling 应该有 negative
	if len(s.Negative) == 0 {
		t.Error("kling should have negatives")
	}
}

// TestGeneratePrompt_WithRefImage 验证带参考图
func TestGeneratePrompt_WithRefImage(t *testing.T) {
	s, err := GeneratePrompt("scene", "midjourney", "/tmp/ref.jpg", baseDir(t))
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if s.RefImage != "/tmp/ref.jpg" {
		t.Errorf("RefImage = %q", s.RefImage)
	}
}

// TestUpdateField 验证 UpdateField
func TestUpdateField(t *testing.T) {
	s, _ := GeneratePrompt("seed", "midjourney", "", baseDir(t))
	s.UpdateField("subject", "a dog")
	if s.Fields["subject"] != "a dog" {
		t.Errorf("UpdateField failed: %q", s.Fields["subject"])
	}
}

// TestUpdateWeight 验证 UpdateWeight
func TestUpdateWeight(t *testing.T) {
	s, _ := GeneratePrompt("seed", "midjourney", "", baseDir(t))
	s.UpdateWeight("subject", 1.8)
	if s.Weights["subject"] != 1.8 {
		t.Errorf("UpdateWeight failed: %v", s.Weights["subject"])
	}
}

// TestRegenerate 验证 Regenerate
func TestRegenerate(t *testing.T) {
	s, _ := GeneratePrompt("seed", "midjourney", "", baseDir(t))
	original := s.FinalPrompt
	s.UpdateField("subject", "changed value")
	s.Regenerate("midjourney", baseDir(t))
	if s.FinalPrompt == original {
		t.Error("Regenerate should update FinalPrompt after field change")
	}
}

// TestRegenerate_UnknownPlatform 验证 Regenerate 走 unknown 分支
func TestRegenerate_UnknownPlatform(t *testing.T) {
	s, _ := GeneratePrompt("seed text", "midjourney", "", baseDir(t))
	s.Regenerate("unknown_xyz", baseDir(t))
	if s.FinalPrompt != "seed text" {
		t.Errorf("unknown Regenerate should set FinalPrompt to seed, got %q", s.FinalPrompt)
	}
}

// TestLoadTemplate_Midjourney 验证 midjourney yaml 模板
func TestLoadTemplate_Midjourney(t *testing.T) {
	tmpl := loadTemplate("midjourney", baseDir(t))
	if tmpl == nil {
		t.Skip("midjourney template not found in test env")
	}
	if tmpl.ID == "" {
		t.Error("template ID empty")
	}
	if len(tmpl.Sections) < 5 {
		t.Errorf("expected >= 5 sections, got %d", len(tmpl.Sections))
	}
}

// TestBuildWeightedPrompt_Deterministic 验证同输入产生同输出
func TestBuildWeightedPrompt_Deterministic(t *testing.T) {
	tmpl := loadTemplate("midjourney", baseDir(t))
	if tmpl == nil {
		t.Skip("midjourney template not found")
	}
	fields := map[string]string{"subject": "a cat"}
	weights := map[string]float64{}
	a := buildWeightedPrompt(tmpl, fields, weights, nil)
	b := buildWeightedPrompt(tmpl, fields, weights, nil)
	if a != b {
		t.Error("buildWeightedPrompt is not deterministic")
	}
}

// TestBuildWeightedPrompt_WithNegative 验证 negative 加入
// 注: 当前实现用 ", " (带空格) join negatives
func TestBuildWeightedPrompt_WithNegative(t *testing.T) {
	tmpl := loadTemplate("midjourney", baseDir(t))
	if tmpl == nil {
		t.Skip("midjourney template not found")
	}
	neg := []string{"blurry", "ugly"}
	got := buildWeightedPrompt(tmpl, map[string]string{"subject": "x"}, nil, neg)
	// 实际格式: "--no blurry, ugly" (带空格, 这是当前行为)
	if !strings.Contains(got, "--no") {
		t.Errorf("--no marker missing: %q", got)
	}
	if !strings.Contains(got, "blurry") || !strings.Contains(got, "ugly") {
		t.Errorf("negatives not in output: %q", got)
	}
}

// TestBuildWeightedPrompt_EmptyFields 验证空 fields
func TestBuildWeightedPrompt_EmptyFields(t *testing.T) {
	tmpl := loadTemplate("midjourney", baseDir(t))
	if tmpl == nil {
		t.Skip("midjourney template not found")
	}
	got := buildWeightedPrompt(tmpl, map[string]string{}, nil, nil)
	if got == "" {
		t.Error("empty fields should still produce non-empty output (params)")
	}
}

// TestNegativeBank_HasEntries 验证 NegativeBank 包含主流平台
func TestNegativeBank_HasEntries(t *testing.T) {
	required := []string{"midjourney", "kling", "jimeng", "gemini"}
	for _, p := range required {
		if _, ok := NegativeBank[p]; !ok {
			t.Errorf("NegativeBank missing %q", p)
		}
	}
}

// TestGeneratePrompt_Deterministic 验证同输入产生同输出 (md5 思想)
func TestGeneratePrompt_Deterministic(t *testing.T) {
	a, _ := GeneratePrompt("a cat in moonlight", "midjourney", "", baseDir(t))
	b, _ := GeneratePrompt("a cat in moonlight", "midjourney", "", baseDir(t))
	if a.FinalPrompt != b.FinalPrompt {
		t.Error("GeneratePrompt not deterministic")
	}
}

// TestInjectReferenceImage 验证参考图注入到 LLM 占位
func TestInjectReferenceImage(t *testing.T) {
	s := &PromptScaffold{
		Fields: map[string]string{"subject": "（LLM填充：主体）"},
	}
	injectReferenceImage(s, "/tmp/myref.jpg")
	if !strings.Contains(s.Fields["subject"], "参考图") {
		t.Errorf("reference image not injected: %q", s.Fields["subject"])
	}
	if !strings.Contains(s.Fields["subject"], "myref.jpg") {
		t.Errorf("ref basename not in field: %q", s.Fields["subject"])
	}
}
