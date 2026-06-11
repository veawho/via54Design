// via54Design — 验证 §12 SVG 规范 (v0.5.3)
//
// 检查: viewBox=680, class t/ts/th, fill=none, text-anchor, 12/14/24px
package export

import (
	"os"
	"strings"
	"testing"
)

func TestSVG_Section12Compliance(t *testing.T) {
	dir, err := os.MkdirTemp("", "svg12-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	scenes := []SVGScene{
		{Title: "测试场景", Voiceover: "测试旁白", Body: "第一行\n第二行", Mood: "calm", BeatName: "Intro", SceneNo: 1, TotalScenes: 3},
	}
	paths, err := ExportSVG(scenes, dir, 1920, 1080)
	if err != nil {
		t.Fatalf("ExportSVG: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 svg, got %d", len(paths))
	}
	data, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)

	checks := []struct {
		name    string
		mustHave string
	}{
		{"L0 viewBox=680x382 (16:9 §12)", `viewBox="0 0 680 382"`},
		{"L1 class t  正文 12px", `class="t"`},
		{"L1 class ts 副标 14px", `class="ts"`},
		{"L1 class th 标题 24px", `class="th"`},
		{"L2 fill=none on stroke lines", `fill="none"`},
		{"L3 text-anchor=end 场景计数", `text-anchor="end"`},
		{"L4 font-size 12px (§12 small)", `font-size: 12px`},
		{"L4 font-size 14px (§12 sub)", `font-size: 14px`},
		{"L4 font-size 24px (§12 title)", `font-size: 24px`},
		{"L5 stroke-width 1.5 (§12 std)", `stroke-width="1.5"`},
		{"L6 兼容性 width=1920 height=1080", `width="1920" height="1080"`},
	}
	for _, c := range checks {
		if !strings.Contains(s, c.mustHave) {
			t.Errorf("§12 规范缺失: %s (期望含 %q)", c.name, c.mustHave)
		}
	}

	// 反向: 不应再用旧版 font-size=24/48/18/28/16
	banned := []string{
		`font-size="24"`, // 旧版标题
		`font-size="48"`, // 旧版主标题
		`font-size="18"`, // 旧版小字
		`font-size="28"`, // 旧版 voiceover
		`font-size="16"`, // 旧版 scene no
	}
	for _, b := range banned {
		if strings.Contains(s, b) {
			t.Errorf("§12 反向检查: 仍含旧版 %q", b)
		}
	}
}