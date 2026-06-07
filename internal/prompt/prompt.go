// via54Design — 图片提示词 (Prompt) 引擎
// 从人类的基础场景描述 → 结构化 prompt → 多平台生图
// 参考: Hao0321/ai-media-generator (⭐70) + Midjourney Prompt Engineering
package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// ─── 提示词模板定义 (从 YAML 加载) ───

type PromptTemplate struct {
	ID          string            `yaml:"id"`
	Name        map[string]string `yaml:"name"`
	Platform    string            `yaml:"platform"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`

	Sections []PromptSection `yaml:"sections"`
	Format   string          `yaml:"format"` // midjourney / kling / jimeng / gemini / comfyui
	Params   map[string]string `yaml:"params,omitempty"` // 平台参数
}

type PromptSection struct {
	ID      string `yaml:"id"`
	Label   string `yaml:"label"`
	Hint    string `yaml:"hint,omitempty"`
	Default string `yaml:"default,omitempty"`
	Options []Option `yaml:"options,omitempty"`
}

type Option struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

// ─── 提示词脚手架 ───

type PromptScaffold struct {
	Seed        string              `yaml:"seed" json:"seed"`
	Platform    string              `yaml:"platform" json:"platform"`
	Model       string              `yaml:"model" json:"model"`
	Fields      map[string]string   `yaml:"fields" json:"fields"`
	FinalPrompt string              `yaml:"final_prompt" json:"final_prompt"`
	Params      map[string]string   `yaml:"params,omitempty" json:"params,omitempty"`
	Expanded    string              `yaml:"expanded,omitempty" json:"expanded,omitempty"`
}

// ─── 引擎 ───

// GeneratePrompt 从场景描述生成结构化提示词
func GeneratePrompt(scene string, platform string, baseDir string) (*PromptScaffold, error) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		// fallback: 通用提示词
		return generateGeneric(scene, platform), nil
	}

	scaffold := &PromptScaffold{
		Seed:     scene,
		Platform: platform,
		Model:    tmpl.Name["zh"],
		Fields:   make(map[string]string),
		Params:   tmpl.Params,
	}

	// 填充字段
	for _, sec := range tmpl.Sections {
		val := sec.Default
		if val == "" {
			val = fmt.Sprintf("（LLM填充：%s）", sec.Hint)
		}
		scaffold.Fields[sec.ID] = val
	}

	// 生成最终 prompt
	scaffold.FinalPrompt = buildFinalPrompt(tmpl, scaffold.Fields)
	scaffold.Expanded = buildPromptDebug(scaffold, tmpl)

	return scaffold, nil
}

// UpdateField 更新某个字段的值（人类介入修改）
func (s *PromptScaffold) UpdateField(id, value string) {
	s.Fields[id] = value
}

// Regenerate 重新生成最终 prompt
func (s *PromptScaffold) Regenerate(platform string, baseDir string) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		s.FinalPrompt = s.Seed
		return
	}
	s.FinalPrompt = buildFinalPrompt(tmpl, s.Fields)
	s.Expanded = buildPromptDebug(s, tmpl)
}

func loadTemplate(platform, baseDir string) *PromptTemplate {
	candidates := []string{
		filepath.Join(baseDir, "templates", "prompts", fmt.Sprintf("%s.yaml", platform)),
		filepath.Join(baseDir, "templates", "prompts", fmt.Sprintf("%s.yml", platform)),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil { continue }
		var t PromptTemplate
		if err := yaml.Unmarshal(data, &t); err == nil {
			return &t
		}
	}
	return nil
}

func generateGeneric(scene, platform string) *PromptScaffold {
	return &PromptScaffold{
		Seed:     scene,
		Platform: platform,
		Model:    "通用",
		Fields:   map[string]string{"scene": scene},
		FinalPrompt: scene,
	}
}

func buildFinalPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	switch tmpl.Format {
	case "midjourney":
		return buildMidjourneyPrompt(tmpl, fields)
	case "kling":
		return buildKlingPrompt(tmpl, fields)
	case "jimeng":
		return buildJimengPrompt(tmpl, fields)
	case "gemini":
		return buildGeminiPrompt(tmpl, fields)
	default:
		return buildGenericPrompt(tmpl, fields)
	}
}

// ─── 各平台提示词构建 ───

func buildMidjourneyPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	var parts []string
	for _, sec := range tmpl.Sections {
		if v, ok := fields[sec.ID]; ok && v != "" {
			parts = append(parts, v)
		}
	}
	prompt := strings.Join(parts, ", ")

	// 追加参数
	var params []string
	for k, v := range tmpl.Params {
		params = append(params, fmt.Sprintf("--%s %s", k, v))
	}
	if len(params) > 0 {
		prompt += " " + strings.Join(params, " ")
	}
	return prompt
}

func buildKlingPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	return buildMidjourneyPrompt(tmpl, fields) // 结构类似
}

func buildJimengPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	return buildMidjourneyPrompt(tmpl, fields)
}

func buildGeminiPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	// Gemini/Imagen 偏好自然语言描述
	var parts []string
	for _, sec := range tmpl.Sections {
		if v, ok := fields[sec.ID]; ok && v != "" {
			label := sec.Label
			if strings.HasPrefix(v, "（LLM") {
				parts = append(parts, fmt.Sprintf("%s: %s", label, v))
			} else {
				parts = append(parts, v)
			}
		}
	}
	return strings.Join(parts, "\n")
}

func buildGenericPrompt(tmpl *PromptTemplate, fields map[string]string) string {
	var parts []string
	for _, sec := range tmpl.Sections {
		if v, ok := fields[sec.ID]; ok && v != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", sec.Label, v))
		}
	}
	return strings.Join(parts, "\n")
}

func buildPromptDebug(s *PromptScaffold, tmpl *PromptTemplate) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📋 平台: %s | 格式: %s\n\n", s.Platform, tmpl.Format))
	b.WriteString("字段列表:\n")
	for _, sec := range tmpl.Sections {
		val := s.Fields[sec.ID]
		mark := "✅" 
		if strings.HasPrefix(val, "（LLM") { mark = "⏳" }
		b.WriteString(fmt.Sprintf("  %s [%s] %s\n", mark, sec.ID, val[:min(len(val), 60)]))
	}
	b.WriteString(fmt.Sprintf("\n最终 prompt:\n%s\n", s.FinalPrompt))
	return b.String()
}

// RenderMarkdown 输出 markdown 格式
func (s *PromptScaffold) RenderMarkdown() (string, error) {
	tmpl := template.Must(template.New("prompt").Parse(markdownTemplate))
	var buf strings.Builder
	err := tmpl.Execute(&buf, s)
	return buf.String(), err
}

const markdownTemplate = `# 🎨 图片提示词 (Prompt)

## 来源场景
> {{.Seed}}

**平台**: {{.Platform}}
**格式**: {{.Model}}

---

## 字段

| 字段 | 值 |
|------|-----|
{{range $k, $v := .Fields}}| {{$k}} | {{$v}} |
{{end}}

---

## 最终 Prompt

` + "```" + `
{{.FinalPrompt}}
` + "```" + `

---

## 平台参数

{{range $k, $v := .Params}}| ` + "`" + `--{{$k}}` + "`" + ` | {{$v}} |
{{else}}无额外参数{{end}}

---

*由 via54 prompt 生成 — 参考: ai-media-generator (⭐70)*
`

func min(a, b int) int { if a < b { return a }; return b }
