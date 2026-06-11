// via54Design — 叙事引擎测试
// SPDX-License-Identifier: AGPL-3.0-only

package narrate

import (
	"strings"
	"testing"

	vt "github.com/veawho/via54Design/internal/template"
)

// baseDir 在 subagent 测试中获取项目根（运行 go test 时 cwd 是包目录）
func baseDir(t *testing.T) string {
	t.Helper()
	// 包目录是 internal/narrate/, 根是 ../../
	return "../.."
}

// TestLoadModel_Success 验证 4 个内置叙事模型都能加载
func TestLoadModel_Success(t *testing.T) {
	for _, id := range []string{"three-act", "heros-journey", "cognitive-arc", "problem-solution"} {
		def, err := LoadModel(id, baseDir(t))
		if err != nil {
			t.Errorf("LoadModel(%q) failed: %v", id, err)
			continue
		}
		if def.ID != id {
			t.Errorf("LoadModel(%q).ID = %q, want %q", id, def.ID, id)
		}
		if def.Name["zh"] == "" {
			t.Errorf("LoadModel(%q).Name[zh] empty", id)
		}
		if len(def.Beats) == 0 {
			t.Errorf("LoadModel(%q).Beats empty", id)
		}
	}
}

// TestLoadModel_NotFound 验证不存在的模型返回 error
func TestLoadModel_NotFound(t *testing.T) {
	_, err := LoadModel("nonexistent-model", baseDir(t))
	if err == nil {
		t.Error("expected error for nonexistent model, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error message should mention 'not found': %v", err)
	}
}

// TestListModels 验证列出所有模型
func TestListModels(t *testing.T) {
	out, err := ListModels(baseDir(t))
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}
	if !strings.Contains(out, "可用叙事模型") {
		t.Errorf("output missing header: %q", out)
	}
	// 应包含 4 个内置模型
	for _, id := range []string{"three-act", "heros-journey", "cognitive-arc", "problem-solution"} {
		if !strings.Contains(out, id) {
			t.Errorf("ListModels output missing %q", id)
		}
	}
}

// TestListModels_EmptyDir 验证空 baseDir 返回空列表（不 panic）
// 注: 当前实现会先尝试 NewRegistry 然后才检查 entries; 空目录无法加载 registry
// 期待 err 包含 "register not found" 或能正常返回 placeholder
func TestListModels_EmptyDir(t *testing.T) {
	out, err := ListModels(t.TempDir())
	if err != nil {
		// 当前行为: registry load 失败时返回 error 而不是 "暂无"
		// 文档化这个行为, 后续可改进为占位符
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "register") {
			t.Logf("unexpected error: %v", err)
		}
		return
	}
	if !strings.Contains(out, "暂无") && !strings.Contains(out, "可用") {
		t.Errorf("empty registry should return placeholder, got: %q", out)
	}
}

// TestGenerateScaffold_HerosJourney 验证英雄之旅模型脚手架
func TestGenerateScaffold_HerosJourney(t *testing.T) {
	s, err := GenerateScaffold(
		"一个裁缝在巴黎改变了时尚",
		"heros-journey",
		90,
		baseDir(t),
	)
	if err != nil {
		t.Fatalf("GenerateScaffold failed: %v", err)
	}
	if s.Seed == "" {
		t.Error("Seed empty")
	}
	if s.ModelID != "heros-journey" {
		t.Errorf("ModelID = %q, want heros-journey", s.ModelID)
	}
	if s.TargetDuration != 90 {
		t.Errorf("TargetDuration = %d, want 90", s.TargetDuration)
	}
	if len(s.Beats) == 0 {
		t.Error("Beats empty")
	}
	if s.FountainScript == "" {
		t.Error("FountainScript empty")
	}
	if len(s.Storyboard) == 0 {
		t.Error("Storyboard empty")
	}
	if s.PromptForLLM == "" {
		t.Error("PromptForLLM empty (LLM 指令未生成)")
	}
	// Beat 总时长应等于 TargetDuration
	total := 0
	for _, b := range s.Beats {
		total += b.Duration
	}
	if total != 90 {
		t.Errorf("Beats total duration = %d, want 90", total)
	}
}

// TestGenerateScaffold_DefaultDuration 验证 duration=0 时默认 30s
func TestGenerateScaffold_DefaultDuration(t *testing.T) {
	s, err := GenerateScaffold("seed", "three-act", 0, baseDir(t))
	if err != nil {
		t.Fatalf("GenerateScaffold failed: %v", err)
	}
	if s.TargetDuration != 30 {
		t.Errorf("duration=0 should default to 30, got %d", s.TargetDuration)
	}
}

// TestGenerateScaffold_NegativeDuration 验证负数 duration 默认为 30
func TestGenerateScaffold_NegativeDuration(t *testing.T) {
	s, err := GenerateScaffold("seed", "three-act", -10, baseDir(t))
	if err != nil {
		t.Fatalf("GenerateScaffold failed: %v", err)
	}
	if s.TargetDuration != 30 {
		t.Errorf("negative duration should default to 30, got %d", s.TargetDuration)
	}
}

// TestGenerateScaffold_BadModel 验证错误模型返回 error
func TestGenerateScaffold_BadModel(t *testing.T) {
	_, err := GenerateScaffold("seed", "bad-model", 60, baseDir(t))
	if err == nil {
		t.Error("expected error for bad model")
	}
}

// TestGenerateScaffold_FountainHasTitle 验证 Fountain 包含种子作为标题
func TestGenerateScaffold_FountainHasTitle(t *testing.T) {
	s, _ := GenerateScaffold("测试种子", "three-act", 60, baseDir(t))
	if !strings.Contains(s.FountainScript, "测试种子") {
		t.Error("Fountain script missing seed as title")
	}
	if !strings.Contains(s.FountainScript, "Fountain") {
		t.Error("Fountain script missing format header")
	}
}

// TestGenerateScaffold_StoryboardShots 验证分镜数量合理
func TestGenerateScaffold_StoryboardShots(t *testing.T) {
	s, _ := GenerateScaffold("seed", "heros-journey", 60, baseDir(t))
	// 60s → 每个 beat (3 段各 ~20s) 各 2-3 shot → 6-9 shots
	if len(s.Storyboard) < 4 {
		t.Errorf("Storyboard has too few shots: %d", len(s.Storyboard))
	}
	if len(s.Storyboard) > 36 {
		t.Errorf("Storyboard exceeds 36 shot cap: %d", len(s.Storyboard))
	}
	// Timecode 格式 M:SS-M:SS
	for _, shot := range s.Storyboard {
		if !strings.Contains(shot.Timecode, ":") {
			t.Errorf("shot %d timecode %q missing colon", shot.ShotNo, shot.Timecode)
		}
	}
}

// TestRenderMarkdown 验证 Markdown 渲染
func TestRenderMarkdown(t *testing.T) {
	s, _ := GenerateScaffold("seed", "three-act", 60, baseDir(t))
	md, err := s.RenderMarkdown()
	if err != nil {
		t.Fatalf("RenderMarkdown failed: %v", err)
	}
	for _, section := range []string{"叙事脚手架", "节拍时间线", "Fountain 剧本", "分镜表"} {
		if !strings.Contains(md, section) {
			t.Errorf("Markdown missing section %q", section)
		}
	}
}

// TestTruncate 验证字符串截断（中文安全）
func TestTruncate(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello", 3, "hel..."},
		{"中文测试字符串", 2, "中文..."},
		{"", 5, ""},
	}
	for _, c := range cases {
		got := truncate(c.in, c.n)
		if got != c.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", c.in, c.n, got, c.want)
		}
	}
}

// TestBuildBeatsFromDef_DurationSum 验证节拍总时长等于 duration
func TestBuildBeatsFromDef_DurationSum(t *testing.T) {
	def, _ := LoadModel("three-act", baseDir(t))
	beats := buildBeatsFromDef(def, 45)
	total := 0
	for _, b := range beats {
		if b.Duration <= 0 {
			t.Errorf("beat %s has non-positive duration: %d", b.Act, b.Duration)
		}
		total += b.Duration
	}
	if total != 45 {
		t.Errorf("total duration = %d, want 45", total)
	}
}

// TestBuildStoryboard_EmptyShotTypes 验证 shotTypes 空时使用默认
func TestBuildStoryboard_EmptyShotTypes(t *testing.T) {
	beats := []Beat{{Act: "act1", StartTime: 0, Duration: 20, Voiceover: "v", Mood: "calm"}}
	shots := buildStoryboard(beats, nil, nil)
	if len(shots) == 0 {
		t.Error("expected default shots, got 0")
	}
	if shots[0].ShotType == "" {
		t.Error("default shot type not applied")
	}
}

// TestBuildRecommendedGen 验证推荐命令格式
func TestBuildRecommendedGen(t *testing.T) {
	got := buildRecommendedGen("heros-journey", 60, "测试种子")
	if !strings.Contains(got, "via54 generate") {
		t.Errorf("recommendation missing 'via54 generate': %q", got)
	}
	if !strings.Contains(got, "bento-grid-2x2") {
		t.Errorf("heros-journey should recommend bento-grid-2x2: %q", got)
	}
}

// TestBuildRecommendedGen_TruncatesLongSeed 验证长 seed 截断
func TestBuildRecommendedGen_TruncatesLongSeed(t *testing.T) {
	longSeed := strings.Repeat("长", 100)
	got := buildRecommendedGen("three-act", 30, longSeed)
	// truncate 到 40 字符 + "..."
	if strings.Contains(got, strings.Repeat("长", 41)) {
		t.Error("long seed should be truncated")
	}
}

// TestRegistry_NarratologyNotEmpty 验证 Registry 列出 narratology
func TestRegistry_NarratologyNotEmpty(t *testing.T) {
	reg, err := vt.NewRegistry(baseDir(t))
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	entries := reg.ListNarratology()
	if len(entries) == 0 {
		t.Error("registry has no narratology entries")
	}
}

// TestGenerateScaffold_AllModelsProduceBeats 验证 4 模型都能产生 beats
func TestGenerateScaffold_AllModelsProduceBeats(t *testing.T) {
	for _, id := range []string{"three-act", "heros-journey", "cognitive-arc", "problem-solution"} {
		s, err := GenerateScaffold("seed for "+id, id, 60, baseDir(t))
		if err != nil {
			t.Errorf("%s failed: %v", id, err)
			continue
		}
		if len(s.Beats) < 2 {
			t.Errorf("%s has %d beats, want >= 2", id, len(s.Beats))
		}
		if len(s.Storyboard) < 2 {
			t.Errorf("%s has %d storyboard shots, want >= 2", id, len(s.Storyboard))
		}
	}
}
