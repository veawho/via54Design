// via54Design — v2.3 i2i engine tests
// SPDX-License-Identifier: AGPL-3.0-only

package prompt

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestParseBilingual_BothHalves 验证 "中文 (English)" 拆分
func TestParseBilingual_BothHalves(t *testing.T) {
	b := ParseBilingual("摄影写实 (Photorealistic)")
	if b.ZH != "摄影写实" {
		t.Errorf("ZH = %q, want 摄影写实", b.ZH)
	}
	if b.EN != "Photorealistic" {
		t.Errorf("EN = %q, want Photorealistic", b.EN)
	}
}

// TestParseBilingual_PureEnglish 验证纯英文
func TestParseBilingual_PureEnglish(t *testing.T) {
	b := ParseBilingual("warm amber tones")
	if b.ZH != "" {
		t.Errorf("ZH should be empty for pure English, got %q", b.ZH)
	}
	if b.EN != "warm amber tones" {
		t.Errorf("EN = %q, want warm amber tones", b.EN)
	}
}

// TestParseBilingual_PureChinese 验证纯中文
func TestParseBilingual_PureChinese(t *testing.T) {
	b := ParseBilingual("暖琥珀色")
	if b.ZH != "暖琥珀色" {
		t.Errorf("ZH = %q, want 暖琥珀色", b.ZH)
	}
	if b.EN != "" {
		t.Errorf("EN should be empty for pure Chinese, got %q", b.EN)
	}
}

// TestParseBilingual_ParensNotBilingual 验证 "POV (first person)" 不被错误切分
// 关键: 第一个括号前是英文, 不应被当成中文+EN
func TestParseBilingual_ParensNotBilingual(t *testing.T) {
	b := ParseBilingual("POV (first person)")
	// 因为"POV"没有 CJK 字符, 我们的 hasCJK 检查会决定整个是英文
	if b.ZH != "" {
		t.Errorf("ZH should be empty, got %q", b.ZH)
	}
	if b.EN == "" {
		t.Errorf("EN should not be empty, got empty")
	}
}

// TestMaxCharsForPlatform 验证 16 平台限额
func TestMaxCharsForPlatform(t *testing.T) {
	cases := map[string]int{
		"midjourney": 4500,
		"jimeng":     1500,
		"gemini":     8000,
		"dalle3":     4000,
		"flux":       4000,
		"sd3":        2000,
		"stable_diffusion": 2000,
		"kling":      2500,
	}
	for plat, want := range cases {
		got := MaxCharsForPlatform(plat)
		if got != want {
			t.Errorf("MaxCharsForPlatform(%q) = %d, want %d", plat, got, want)
		}
	}
}

// TestGenerateI2I_NoRef 验证无 ref 时, subject 留 LLM填充
func TestGenerateI2I_NoRef(t *testing.T) {
	r, err := GenerateI2I(I2IRequest{
		Scene:    "温馨家庭客厅",
		Platform: "jimeng",
		BaseDir:  baseDir(t),
	})
	if err != nil {
		t.Fatalf("GenerateI2I: %v", err)
	}
	if r.Platform != "jimeng" {
		t.Errorf("Platform = %q, want jimeng", r.Platform)
	}
	if r.FinalChars < 100 {
		t.Errorf("FinalChars = %d, want > 100 (engine should at least emit params + negative)", r.FinalChars)
	}
	if r.MaxChars != 1500 {
		t.Errorf("MaxChars = %d, want 1500 (jimeng default)", r.MaxChars)
	}
	// subject 应为 LLM填充(没 ref 没 override)
	if !strings.HasPrefix(r.Sections["subject"], "（LLM填充") {
		t.Errorf("subject = %q, want to start with （LLM填充", r.Sections["subject"])
	}
}

// TestGenerateI2I_WithRefLock 验证 ref 锁住 subject/secondary/environment
func TestGenerateI2I_WithRefLock(t *testing.T) {
	refDesc := "modern Scandinavian living room, white walls, beige sofa, plants"
	r, err := GenerateI2I(I2IRequest{
		Scene:          "温馨家庭客厅",
		Platform:       "jimeng",
		RefImage:       "/tmp/fake-ref.jpg",
		RefDescription: refDesc,
		BaseDir:        baseDir(t),
	})
	if err != nil {
		t.Fatalf("GenerateI2I: %v", err)
	}
	// subject 应被 ref 锁
	if !strings.Contains(r.Sections["subject"], "preserve visual identity") {
		t.Errorf("subject should be ref-locked, got %q", r.Sections["subject"])
	}
	if !strings.Contains(r.Sections["subject"], refDesc) {
		t.Errorf("subject should contain ref description, got %q", r.Sections["subject"])
	}
	// environment 也应被锁
	if !strings.Contains(r.Sections["environment"], "preserve visual identity") {
		t.Errorf("environment should be ref-locked, got %q", r.Sections["environment"])
	}
	// subject 和 environment 不应完全相同 (ref-lock 语义要求互补)
	if r.Sections["subject"] == r.Sections["environment"] {
		t.Errorf("subject and environment should differ; both = %q", r.Sections["subject"])
	}
	// 主体中文版应含 ref 描述
	if !strings.Contains(r.SectionsZH["subject"], refDesc) {
		t.Errorf("subject ZH should contain ref description, got %q", r.SectionsZH["subject"])
	}
}

// TestGenerateI2I_OverridesRespected 验证 overrides 优先于 ref-lock
func TestGenerateI2I_OverridesRespected(t *testing.T) {
	r, err := GenerateI2I(I2IRequest{
		Scene:          "温馨家庭客厅",
		Platform:       "jimeng",
		RefImage:       "/tmp/fake-ref.jpg",
		RefDescription: "some ref",
		Overrides: map[string]string{
			"subject": "红发蓝裙美少女",
		},
		BaseDir: baseDir(t),
	})
	if err != nil {
		t.Fatalf("GenerateI2I: %v", err)
	}
	if r.Sections["subject"] != "红发蓝裙美少女" {
		t.Errorf("override should win over ref-lock, got %q", r.Sections["subject"])
	}
	// environment 没 override, 应被 ref 锁
	if !strings.Contains(r.Sections["environment"], "preserve visual identity") {
		t.Errorf("environment should still be ref-locked, got %q", r.Sections["environment"])
	}
}

// TestGenerateI2I_8Platforms 验证 8 个平台都能正常输出
func TestGenerateI2I_8Platforms(t *testing.T) {
	platforms := []string{"midjourney", "jimeng", "gemini", "stable_diffusion", "flux", "dalle3", "kling", "veo"}
	for _, p := range platforms {
		r, err := GenerateI2I(I2IRequest{
			Scene: "温馨家庭客厅", Platform: p, BaseDir: baseDir(t),
		})
		if err != nil {
			t.Errorf("platform %s: %v", p, err)
			continue
		}
		if r.MaxChars <= 0 {
			t.Errorf("platform %s: MaxChars = %d, want > 0", p, r.MaxChars)
		}
		if r.FinalChars <= 0 {
			t.Errorf("platform %s: FinalChars = %d, want > 0", p, r.FinalChars)
		}
	}
}

// TestGenerateI2I_JSONSerializable 验证 I2IResult 可以 JSON 序列化
// (飞书 inbox_watcher 跨进程消费依赖此特性)
func TestGenerateI2I_JSONSerializable(t *testing.T) {
	r, err := GenerateI2I(I2IRequest{
		Scene: "test", Platform: "jimeng", BaseDir: baseDir(t),
	})
	if err != nil {
		t.Fatalf("GenerateI2I: %v", err)
	}
	js, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	// 必须包含飞书端用得上的关键字段
	for _, key := range []string{`"FinalEN"`, `"FormattedZH"`, `"Platform"`, `"Sections"`, `"SectionsZH"`} {
		if !strings.Contains(string(js), key) {
			t.Errorf("JSON missing %s: %s", key, js[:200])
		}
	}
}

// TestPlanFill 验证限额最优填充
func TestPlanFill(t *testing.T) {
	r, err := GenerateI2I(I2IRequest{
		Scene: "x", Platform: "jimeng", BaseDir: baseDir(t),
	})
	if err != nil {
		t.Fatalf("GenerateI2I: %v", err)
	}
	plan := PlanFill(r, "fake ref desc")
	if plan.MaxChars != 1500 {
		t.Errorf("plan MaxChars = %d, want 1500", plan.MaxChars)
	}
	if plan.RemainingBudget <= 0 {
		t.Errorf("RemainingBudget = %d, want > 0", plan.RemainingBudget)
	}
	if len(plan.Targets) == 0 {
		t.Errorf("plan should have fill targets, got 0")
	}
	// target chars 应在 [30, 350] 范围
	for _, tgt := range plan.Targets {
		if tgt.TargetChars < 30 || tgt.TargetChars > 350 {
			t.Errorf("target %s: chars = %d, want 30-350", tgt.SectionID, tgt.TargetChars)
		}
	}
}

// TestNormalizeFieldZH 验证双语归一化
func TestNormalizeFieldZH(t *testing.T) {
	if got := NormalizeFieldZH("摄影写实 (Photorealistic)"); got != "摄影写实" {
		t.Errorf("got %q, want 摄影写实", got)
	}
	if got := NormalizeFieldZH("pure english"); got != "pure english" {
		t.Errorf("got %q, want 'pure english' (pass-through)", got)
	}
	if got := NormalizeFieldZH("（LLM填充：xxx）"); !strings.HasPrefix(got, "（LLM") {
		t.Errorf("LLM placeholder should be preserved, got %q", got)
	}
}
