// SPDX-License-Identifier: MIT OR AGPL-3.0

// via54Design — 叙事引擎 (Narrate Engine)
// 从一句话到完整故事大纲 + Fountain 剧本 + 分镜表
// 参考: huobao-drama (⭐12.6k) + Fountain screenplay format
//
// 架构: 叙事模型定义在 YAML 文件中 (templates/narratology/models/*.yaml)
//       引擎从 Registry 统一加载，与 layout/color/font 一致
package narrate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	vt "github.com/veawho/via54Design/internal/template"
	"gopkg.in/yaml.v3"
)

// ─── 叙事模型定义 (从 YAML 加载) ───

type NarrativeModelDef struct {
	ID            string                `yaml:"id"`
	Name          map[string]string     `yaml:"name"`
	Description   map[string]string     `yaml:"description"`
	Source        string                `yaml:"source"`
	SuitableFor   []string              `yaml:"suitable_for"`
	Beats         []BeatDef             `yaml:"beats"`
	ShotTypes     []string              `yaml:"shot_types"`
	CameraMoves   []string              `yaml:"camera_moves"`
}

type BeatDef struct {
	ID             string            `yaml:"id"`
	Name           map[string]string `yaml:"name"`
	Mood           string            `yaml:"mood"`
	DurationWeight float64           `yaml:"duration_weight"`
	DurationMin    int               `yaml:"duration_min,omitempty"`
	DurationMax    int               `yaml:"duration_max,omitempty"`
	VoiceoverTmpl  string            `yaml:"voiceover_template"`
	Transition     string            `yaml:"transition,omitempty"`
	SFX            string            `yaml:"sfx,omitempty"`
	Weight         float64           `yaml:"weight,omitempty"`
	SubBeats       []BeatDef         `yaml:"sub_beats,omitempty"`
}

// ─── 叙事脚手架输出 ───

type NarrativeScaffold struct {
	Seed            string                `yaml:"seed" json:"seed"`
	ModelID         string                `yaml:"model_id" json:"model_id"`
	ModelName       string                `yaml:"model_name" json:"model_name"`
	Description     string                `yaml:"description" json:"description"`
	TargetDuration  int                   `yaml:"target_duration" json:"target_duration"`
	ExpandedOutline string                `yaml:"expanded_outline,omitempty" json:"expanded_outline,omitempty"`
	Beats           []Beat                `yaml:"beats,omitempty" json:"beats,omitempty"`
	FountainScript  string                `yaml:"fountain_script,omitempty" json:"fountain_script,omitempty"`
	Storyboard      []Shot                `yaml:"storyboard,omitempty" json:"storyboard,omitempty"`
	PromptForLLM    string                `yaml:"prompt_for_llm,omitempty" json:"prompt_for_llm,omitempty"`
	RecommendedGen  string                `yaml:"recommended_generate,omitempty" json:"recommended_generate,omitempty"`
}

type Beat struct {
	Act       string `yaml:"act" json:"act"`
	StartTime int    `yaml:"start_time" json:"start_time"`
	Duration  int    `yaml:"duration" json:"duration"`
	Event     string `yaml:"event" json:"event"`
	Voiceover string `yaml:"voiceover" json:"voiceover"`
	Mood      string `yaml:"mood" json:"mood"`
	Transition string `yaml:"transition,omitempty" json:"transition,omitempty"`
	SFX       string `yaml:"sfx,omitempty" json:"sfx,omitempty"`
}

type Shot struct {
	ShotNo    int    `yaml:"shot_no" json:"shot_no"`
	Timecode  string `yaml:"timecode" json:"timecode"`
	Duration  int    `yaml:"duration" json:"duration"`
	ShotType  string `yaml:"shot_type" json:"shot_type"`
	Camera    string `yaml:"camera" json:"camera"`
	Visual    string `yaml:"visual" json:"visual"`
	Audio     string `yaml:"audio" json:"audio"`
	Voiceover string `yaml:"voiceover" json:"voiceover"`
	Mood      string `yaml:"mood" json:"mood"`
}

// ─── 模型加载 ───

// LoadModel 从 Registry 加载叙事模型定义
func LoadModel(id, baseDir string) (*NarrativeModelDef, error) {
	reg, err := vt.NewRegistry(baseDir)
	if err != nil {
		return nil, fmt.Errorf("registry load failed: %w", err)
	}

	filePath, err := reg.ResolveNarratology(id)
	if err != nil {
		return nil, fmt.Errorf("model '%s' not found in registry: %w", id, err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read model file %s: %w", filePath, err)
	}

	var def NarrativeModelDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse model %s: %w", filePath, err)
	}

	if def.ID == "" {
		return nil, fmt.Errorf("model file %s missing 'id' field", filePath)
	}

	return &def, nil
}

// ListModels 从 Registry 列出所有可用叙事模型
func ListModels(baseDir string) (string, error) {
	reg, err := vt.NewRegistry(baseDir)
	if err != nil {
		return "", err
	}

	entries := reg.ListNarratology()
	if len(entries) == 0 {
		return "暂无叙事模型注册\n", nil
	}

	var b strings.Builder
	b.WriteString("可用叙事模型:\n\n")
	for _, e := range entries {
		// 只显示 model 类目
		if e.Category == "narratology/model" {
			b.WriteString(fmt.Sprintf("  %-20s  %s\n", e.ID, e.Name))
			if len(e.Tags) > 0 {
				b.WriteString(fmt.Sprintf("  %-20s  适用: %s\n\n", "", strings.Join(e.Tags, ", ")))
			}
		}
	}
	b.WriteString("运行 via54 narrate --seed \"一句话\" --model <id> --duration <秒>\n")
	return b.String(), nil
}

// ─── 生成器 ───

// GenerateScaffold 从 seed sentence 生成完整叙事脚手架
func GenerateScaffold(seed string, modelID string, duration int, baseDir string) (*NarrativeScaffold, error) {
	def, err := LoadModel(modelID, baseDir)
	if err != nil {
		return nil, err
	}

	if duration <= 0 {
		duration = 30
	}

	zhName := def.Name["zh"]
	if zhName == "" {
		zhName = def.ID
	}
	zhDesc := def.Description["zh"]
	if zhDesc == "" {
		zhDesc = def.Description["en"]
	}

	scaffold := &NarrativeScaffold{
		Seed:           seed,
		ModelID:        modelID,
		ModelName:      zhName,
		Description:    zhDesc,
		TargetDuration: duration,
		ExpandedOutline: buildOutlinePrompt(seed, def, duration),
		PromptForLLM:   buildLLMPrompt(seed, def, duration),
	}

	// 生成 beat 骨架
	scaffold.Beats = buildBeatsFromDef(def, duration)

	// Fountain 剧本
	scaffold.FountainScript = buildFountainTemplate(scaffold)

	// 分镜表
	scaffold.Storyboard = buildStoryboard(scaffold.Beats, def.ShotTypes, def.CameraMoves)

	// 推荐 generate 命令
	scaffold.RecommendedGen = buildRecommendedGen(modelID, duration, seed)

	return scaffold, nil
}

func buildOutlinePrompt(seed string, def *NarrativeModelDef, duration int) string {
	var b strings.Builder
	zhName := def.Name["zh"]
	if zhName == "" { zhName = def.ID }
	zhDesc := def.Description["zh"]
	if zhDesc == "" { zhDesc = def.Description["en"] }

	b.WriteString(fmt.Sprintf("# 叙事大纲扩展指令\n\n"))
	b.WriteString(fmt.Sprintf("## 种子句子\n> %s\n\n", seed))
	b.WriteString(fmt.Sprintf("## 选用模型\n**%s** — %s\n", zhName, zhDesc))
	b.WriteString(fmt.Sprintf("## 来源\n%s\n\n", def.Source))
	b.WriteString(fmt.Sprintf("## 节拍结构\n"))
	for _, beat := range def.Beats {
		bn := beat.Name["zh"]
		if bn == "" { bn = beat.ID }
		b.WriteString(fmt.Sprintf("- %s (%s, mood: %s, weight: %.0f%%)\n", bn, beat.ID, beat.Mood, beat.DurationWeight*100))
		if len(beat.SubBeats) > 0 {
			for _, sb := range beat.SubBeats {
				sbn := sb.Name["zh"]
				if sbn == "" { sbn = sb.ID }
				b.WriteString(fmt.Sprintf("  └ %s (%s, mood: %s, weight: %.0f%%)\n", sbn, sb.ID, sb.Mood, sb.Weight*100))
			}
		}
	}
	b.WriteString(fmt.Sprintf("\n## 输出要求\n"))
	b.WriteString(fmt.Sprintf("1. 将种子句子扩展为%d秒的完整故事大纲\n", duration))
	b.WriteString(fmt.Sprintf("2. 以 %s 结构组织\n", zhName))
	b.WriteString(fmt.Sprintf("3. 每个节拍给出: 时间码、画面描述、对白/旁白、情绪\n"))
	b.WriteString(fmt.Sprintf("4. 最后输出一个Fountain格式的完整剧本\n"))
	return b.String()
}

func buildLLMPrompt(seed string, def *NarrativeModelDef, duration int) string {
	var b strings.Builder
	zhName := def.Name["zh"]
	if zhName == "" { zhName = def.ID }

	b.WriteString("你是一个专业的叙事设计师和编剧。请根据以下信息完成一个完整的视频剧本。\n\n")
	b.WriteString("### 用户的一句话种子\n")
	b.WriteString(fmt.Sprintf("\"%s\"\n\n", seed))
	b.WriteString(fmt.Sprintf("### 视频时长\n%d秒\n\n", duration))
	b.WriteString(fmt.Sprintf("### 叙事模型\n%s\n\n", zhName))
	b.WriteString(fmt.Sprintf("来源: %s\n\n", def.Source))

	var beatIDs []string
	for _, bd := range def.Beats {
		bn := bd.Name["zh"]
		if bn == "" { bn = bd.ID }
		beatIDs = append(beatIDs, bn)
	}
	b.WriteString(fmt.Sprintf("节拍: %s\n\n", strings.Join(beatIDs, " → ")))

	b.WriteString(`### 输出格式 (三段式)

## 第一段：故事大纲
每个节拍给1-2句话描述，包含情绪标注和旁白初稿。

## 第二段：Fountain 剧本
  **SCENE 1: 场景标题**  (以 **开头)
  场景描述文字（纯文本）

  旁白
  (语气说明)
  对白

  > 转场方式

## 第三段：分镜表
| Shot | Time | Duration | Type | Camera | Visual | Audio | Voiceover | Mood |
`)
	return b.String()
}

func buildBeatsFromDef(def *NarrativeModelDef, duration int) []Beat {
	var beats []Beat
	curTime := 0
	remaining := duration

	for i, bd := range def.Beats {
		// 计算时长
		segDur := int(float64(duration) * bd.DurationWeight)
		if bd.DurationMin > 0 && segDur < bd.DurationMin {
			segDur = bd.DurationMin
		}
		if bd.DurationMax > 0 && segDur > bd.DurationMax {
			segDur = bd.DurationMax
		}
		// 最后一段占满剩余时间
		if i == len(def.Beats)-1 {
			segDur = remaining
		}
		if segDur <= 0 { segDur = 1 }

		zhName := bd.Name["zh"]
		if zhName == "" { zhName = bd.ID }

		vt := bd.VoiceoverTmpl
		if vt == "" { vt = bd.ID + "..." }

		beats = append(beats, Beat{
			Act:       fmt.Sprintf("%s (%s)", zhName, bd.ID),
			StartTime: curTime,
			Duration:  segDur,
			Event:     fmt.Sprintf("（LLM填充：%s场景）", bd.ID),
			Voiceover: vt,
			Mood:      bd.Mood,
			Transition: bd.Transition,
			SFX:       bd.SFX,
		})

		curTime += segDur
		remaining -= segDur
	}
	return beats
}

func buildFountainTemplate(s *NarrativeScaffold) string {
	var b strings.Builder
	b.WriteString("Fountain Screenplay\n")
	b.WriteString(fmt.Sprintf("Title: %s\n", s.Seed))
	b.WriteString(fmt.Sprintf("Model: %s\n", s.ModelID))
	b.WriteString(fmt.Sprintf("Duration: %ds\n\n", s.TargetDuration))

	for i, beat := range s.Beats {
		b.WriteString(fmt.Sprintf("**SCENE %d: %s**\n", i+1, beat.Act))
		b.WriteString("（LLM填充：场景描述）\n\n")
		b.WriteString("旁白\n")
		b.WriteString(fmt.Sprintf("(%s)\n", beat.Mood))
		b.WriteString(fmt.Sprintf("%s\n\n", beat.Voiceover))
		if beat.SFX != "" {
			b.WriteString(fmt.Sprintf(">> 音效: %s\n\n", beat.SFX))
		}
		if i < len(s.Beats)-1 {
			trans := beat.Transition
			if trans == "" { trans = "切至" }
			b.WriteString(fmt.Sprintf("> %s SCENE %d\n\n", trans, i+2))
		}
	}
	return b.String()
}

func buildStoryboard(beats []Beat, shotTypes, cameraMoves []string) []Shot {
	var shots []Shot
	shotNo := 1

	if len(shotTypes) == 0 {
		shotTypes = []string{"WIDE", "MEDIUM", "CLOSE-UP", "DETAIL", "MEDIUM", "WIDE"}
	}
	if len(cameraMoves) == 0 {
		cameraMoves = []string{"Static", "Slow zoom", "Dolly in", "Static", "Pan", "Crane up"}
	}

	for _, beat := range beats {
		nShots := 2
		if beat.Duration > 15 {
			nShots = 3
		}
		if beat.Duration > 40 {
			nShots = 4
		}

		segDur := beat.Duration / nShots
		for j := 0; j < nShots && shotNo <= 36; j++ {
			st := shotTypes[(shotNo-1)%len(shotTypes)]
			ca := cameraMoves[(shotNo-1)%len(cameraMoves)]
			startT := beat.StartTime + j*segDur
			tc := fmt.Sprintf("%d:%02d-%d:%02d",
				startT/60, startT%60,
				(startT+segDur)/60, (startT+segDur)%60)

			shots = append(shots, Shot{
				ShotNo:    shotNo,
				Timecode:  tc,
				Duration:  segDur,
				ShotType:  st,
				Camera:    ca,
				Visual:    fmt.Sprintf("（LLM填充：画面 %d）", shotNo),
				Audio:     "（LLM填充：音效）",
				Voiceover: beat.Voiceover,
				Mood:      beat.Mood,
			})
			shotNo++
		}
	}
	return shots
}

func buildRecommendedGen(modelID string, duration int, seed string) string {
	// 根据叙事模型推荐合适的视觉模板
	layoutMap := map[string]string{
		"three-act":        "hero-split-left-image",
		"heros-journey":    "bento-grid-2x2",
		"cognitive-arc":    "gallery-waterfall",
		"problem-solution": "hero-split-left-image",
	}
	colorMap := map[string]string{
		"three-act":        "dark-terminal-blue",
		"heros-journey":    "warm-editorial-cream",
		"cognitive-arc":    "crimson-elegance",
		"problem-solution": "neon-dark",
	}

	l := layoutMap[modelID]
	if l == "" { l = "hero-split-left-image" }
	c := colorMap[modelID]
	if c == "" { c = "ink-wash" }

	return fmt.Sprintf("via54 generate --layout %s --color %s --font ming-hei-editorial --title \"%s\"",
		l, c, truncate(seed, 40))
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n { return s }
	return string(runes[:n]) + "..."
}

// ─── 渲染输出 ───

func (s *NarrativeScaffold) RenderMarkdown() (string, error) {
	tmpl := template.Must(template.New("narrative").Parse(markdownTemplate))
	var b strings.Builder
	err := tmpl.Execute(&b, s)
	return b.String(), err
}

const markdownTemplate = `# 🎬 叙事脚手架

**种子**: {{.Seed}}
**模型**: {{.ModelName}} — {{.Description}}
**时长**: {{.TargetDuration}}秒

---

## 一、故事大纲 (给 LLM/人类扩写用)

{{.ExpandedOutline}}

---

## 二、节拍时间线

| 幕/节拍 | 时间 | 时长 | 事件 | 旁白 | 情绪 | 转场 |
|---------|------|------|------|------|------|------|
{{range .Beats}}| {{.Act}} | {{.StartTime}}s | {{.Duration}}s | {{.Event}} | {{.Voiceover}} | _{{.Mood}}_ | {{if .Transition}}{{.Transition}}{{else}}-{{end}} |
{{end}}

---

## 三、Fountain 剧本骨架

` + "```screenplay" + `
{{.FountainScript}}
` + "```" + `

---

## 四、分镜表

| Shot | Timecode | Dur | Type | Camera | Visual | Audio | Voiceover | Mood |
|------|----------|-----|------|--------|--------|-------|-----------|------|
{{range .Storyboard}}| {{.ShotNo}} | {{.Timecode}} | {{.Duration}}s | {{.ShotType}} | {{.Camera}} | {{.Visual}} | {{.Audio}} | {{.Voiceover}} | _{{.Mood}}_ |
{{end}}

---

## 五、LLM 完整提示词

` + "```" + `
{{.PromptForLLM}}
` + "```" + `

---

## 六、推荐生成命令

` + "```bash" + `
{{.RecommendedGen}}
` + "```" + `

---

*由 via54 narrate 生成 — 叙事引擎 v2.0 (基于 YAML 模型模板)*
*参考: huobao-drama (⭐12.6k) + Fountain screenplay format*
`

// ToJSON 输出结构化 JSON（可供 generate --from-narrative 消费）
func (s *NarrativeScaffold) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var _ = os.Getenv // keep os import alive

