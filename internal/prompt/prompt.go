// via54Design — 图片提示词 (Prompt) 引擎 v2
// 16 维度结构化提示词 + 负面词库 + Token 权重 + 参考图
// 参考: designer-image (10维度) + image-gen-pipeline (7分类) + CLIP Interrogator
package prompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// ─── 16 维度提示词模板 ───

type PromptTemplate struct {
	ID          string            `yaml:"id"`
	Name        map[string]string `yaml:"name"`
	Platform    string            `yaml:"platform"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	Sections    []PromptSection   `yaml:"sections"`
	Format      string            `yaml:"format"`
	Params      map[string]string `yaml:"params,omitempty"`
	Negative    []string          `yaml:"negative,omitempty"` // 平台默认负面词
}

type PromptSection struct {
	ID           string            `yaml:"id"`
	Label        string            `yaml:"label"`
	Category     string            `yaml:"category,omitempty"` // subject/style/comp/lighting/color/env/detail/quality
	Hint         string            `yaml:"hint,omitempty"`
	Default      string            `yaml:"default,omitempty"`
	Options      []Option          `yaml:"options,omitempty"`
	Weighted     bool              `yaml:"weighted,omitempty"`  // 支持 (keyword:1.2) 语法
	Weight       float64           `yaml:"weight,omitempty"`    // 默认权重 1.0
}

type Option struct {
	Value  string `yaml:"value"`
	Label  string `yaml:"label"`
	Weight float64 `yaml:"weight,omitempty"` // 选项的权重 (如 (cinematic lighting:1.3))
}

// ─── 提示词脚手架 v2 ───

type PromptScaffold struct {
	Seed        string              `yaml:"seed" json:"seed"`
	Platform    string              `yaml:"platform" json:"platform"`
	Model       string              `yaml:"model" json:"model"`
	Fields      map[string]string   `yaml:"fields" json:"fields"`
	Weights     map[string]float64  `yaml:"weights,omitempty" json:"weights,omitempty"`
	Negative    []string            `yaml:"negative" json:"negative"`
	FinalPrompt string              `yaml:"final_prompt" json:"final_prompt"`
	Params      map[string]string   `yaml:"params,omitempty" json:"params,omitempty"`
	Expanded    string              `yaml:"expanded,omitempty" json:"expanded,omitempty"`
	RefImage    string              `yaml:"ref_image,omitempty" json:"ref_image,omitempty"`
}

// ─── 负面词库 ───

var NegativeBank = map[string][]string{
	"midjourney": {
		"blurry, low quality, distorted, deformed, ugly, bad anatomy",
		"extra limbs, missing fingers, bad hands, poorly drawn hands",
		"watermark, signature, text, logo, brand",
		"oversaturated, underexposed, noise, grainy, jpeg artifacts",
		"amateur, snapshot, lowres, worst quality, normal quality",
	},
	"kling": {
		"jittery, shaky, unstable motion, flickering",
		"distorted face, bad expression, unnatural movement",
		"watermark, text overlay, logo",
		"low resolution, blurry, pixelated",
	},
	"jimeng": {
		"模糊, 变形, 畸形, 崩坏, 低质量",
		"多余肢体, 手指错误, 水印, 文字",
		"过度饱和, 噪点, 压缩痕迹",
	},
	"gemini": {
		"blurry, distorted, low quality, bad anatomy",
		"watermark, text, signature",
		"overexposed, noise, artifacts",
	},
}

// ─── 引擎 ───

func GeneratePrompt(scene string, platform string, refImage string, baseDir string) (*PromptScaffold, error) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		return generateGeneric(scene, platform, refImage), nil
	}

	scaffold := &PromptScaffold{
		Seed:     scene,
		Platform: platform,
		Model:    tmpl.Name["zh"],
		Fields:   make(map[string]string),
		Weights:  make(map[string]float64),
		Negative: NegativeBank[platform],
		Params:   tmpl.Params,
		RefImage: refImage,
	}

	// 填充 16 维度字段
	for _, sec := range tmpl.Sections {
		val := sec.Default
		if val == "" {
			val = fmt.Sprintf("（LLM填充：%s）", sec.Hint)
		}
		scaffold.Fields[sec.ID] = val
		if sec.Weight > 0 {
			scaffold.Weights[sec.ID] = sec.Weight
		}
	}

	// 如果有参考图，注入视觉特征
	if refImage != "" {
		injectReferenceImage(scaffold, refImage)
	}

	// 生成最终 prompt
	scaffold.FinalPrompt = buildWeightedPrompt(tmpl, scaffold.Fields, scaffold.Weights, scaffold.Negative)
	scaffold.Expanded = buildPromptDebug(scaffold, tmpl)
	return scaffold, nil
}

// applyWeights 对字段值应用 Token 权重语法
func applyWeights(value string, weight float64) string {
	if weight == 0 || weight == 1.0 {
		return value
	}
	// 已经是权重语法则跳过
	if strings.Contains(value, ":") && (strings.HasPrefix(value, "(") || strings.HasSuffix(value, ")") || strings.HasSuffix(value, ")") && strings.Contains(value, ":")) {
		return value
	}
	return fmt.Sprintf("(%s:%.1f)", value, weight)
}

func buildWeightedPrompt(tmpl *PromptTemplate, fields map[string]string, weights map[string]float64, negative []string) string {
	// 按 Category 分组
	categories := []string{"subject", "style", "comp", "lighting", "color", "env", "detail", "quality"}
	categorized := make(map[string][]string)
	for _, sec := range tmpl.Sections {
		if v, ok := fields[sec.ID]; ok && v != "" && !strings.HasPrefix(v, "（LLM") {
			w := weights[sec.ID]
			if w == 0 { w = 1.0 }
			if w != 1.0 && sec.Weighted {
				v = applyWeights(v, w)
			}
			cat := sec.Category
			if cat == "" { cat = "other" }
			categorized[cat] = append(categorized[cat], v)
		}
	}

	// 链式组合: subject → style → comp → lighting → color → env → detail → quality
	var parts []string
	for _, cat := range categories {
		if vals, ok := categorized[cat]; ok {
			parts = append(parts, vals...)
		}
	}

	// 追加剩余分类
	if vals, ok := categorized["other"]; ok {
		parts = append(parts, vals...)
	}

	prompt := strings.Join(parts, ", ")

	// 追加平台参数 (排序保证确定性)
	var paramStr []string
	keys := make([]string, 0, len(tmpl.Params))
	for k := range tmpl.Params { keys = append(keys, k) }
	sort.Strings(keys)
	for _, k := range keys {
		paramStr = append(paramStr, fmt.Sprintf("--%s %s", k, tmpl.Params[k]))
	}
	if len(paramStr) > 0 {
		prompt += " " + strings.Join(paramStr, " ")
	}

	// 追加负面词
	if len(negative) > 0 {
		prompt += " --no " + strings.Join(negative, ", ")
	}

	return prompt
}

func injectReferenceImage(s *PromptScaffold, refPath string) {
	// 参考图注入: 标记字段来源
	for k := range s.Fields {
		if strings.HasPrefix(s.Fields[k], "（LLM填充") {
			s.Fields[k] = fmt.Sprintf("（参考图 %s 提取: %s）", filepath.Base(refPath), s.Fields[k])
		}
	}
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
		if err := yaml.Unmarshal(data, &t); err == nil && t.ID != "" {
			return &t
		}
	}
	return nil
}

func generateGeneric(scene, platform, refImage string) *PromptScaffold {
	return &PromptScaffold{
		Seed:     scene,
		Platform: platform,
		Model:    "通用",
		Fields:   map[string]string{"scene": scene},
		Negative: NegativeBank[platform],
		FinalPrompt: scene,
		RefImage: refImage,
	}
}

// ─── 交互式编辑 ───

func (s *PromptScaffold) UpdateField(id, value string) {
	s.Fields[id] = value
}

func (s *PromptScaffold) UpdateWeight(id string, weight float64) {
	s.Weights[id] = weight
}

func (s *PromptScaffold) Regenerate(platform string, baseDir string) {
	tmpl := loadTemplate(platform, baseDir)
	if tmpl == nil {
		s.FinalPrompt = s.Seed
		return
	}
	s.FinalPrompt = buildWeightedPrompt(tmpl, s.Fields, s.Weights, s.Negative)
	s.Expanded = buildPromptDebug(s, tmpl)
}

// ─── 质量评估 (CLIP-AGIQA 参考) ───

type QualityReport struct {
	ImagePath      string  `json:"image_path"`
	OverallScore   float64 `json:"overall_score"`   // 0.0-1.0
	ClarityScore   float64 `json:"clarity_score"`   // 清晰度
	CompositionScore float64 `json:"composition_score"` // 构图
	ColorScore     float64 `json:"color_score"`     // 色彩
	PromptMatch    float64 `json:"prompt_match"`    // 与 prompt 匹配度
	Issues         []string `json:"issues"`
}

// AssessImage 评估生图质量 (基于 CLIP-AGIQA 方法论)
// 实际评估需要 CLIP 模型, 当前输出结构化评估框架
func AssessImage(imagePath string, promptText string) *QualityReport {
	report := &QualityReport{
		ImagePath: imagePath,
		OverallScore: 0.0,
		Issues: []string{},
	}

	// 检查图片是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		report.Issues = append(report.Issues, "图片文件不存在: "+imagePath)
		return report
	}

	// 获取文件大小作为基础质量信号
	if info, err := os.Stat(imagePath); err == nil {
		sizeMB := float64(info.Size()) / (1024 * 1024)
		if sizeMB < 0.1 {
			report.Issues = append(report.Issues, "图片文件过小 (<100KB), 可能低质量")
			report.ClarityScore = 0.3
		} else if sizeMB > 10 {
			report.ClarityScore = 0.9
		} else {
			report.ClarityScore = 0.6 + (sizeMB / 10) * 0.3
		}
	}

	// Prompt 匹配度估算 (基于关键词覆盖)
	if promptText != "" {
		promptWords := len(strings.Fields(promptText))
		if promptWords > 50 {
			report.PromptMatch = 0.8
		} else if promptWords > 20 {
			report.PromptMatch = 0.6
		} else {
			report.PromptMatch = 0.4
		}
	}

	// 综合评分 (实际应调用 CLIP 模型)
	report.CompositionScore = 0.7 // placeholder
	report.ColorScore = 0.7       // placeholder
	report.OverallScore = (report.ClarityScore + report.CompositionScore + report.ColorScore + report.PromptMatch) / 4.0

	return report
}

// ─── 版本管理 ───

type PromptVersion struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Seed      string `json:"seed"`
	Platform  string `json:"platform"`
	Prompt    string `json:"prompt"`
	Previous  string `json:"previous,omitempty"`
}

// SaveVersion 保存 prompt 版本到文件 (git-friendly JSON lines)
func SaveVersion(dir string, s *PromptScaffold) (string, error) {
	os.MkdirAll(dir, 0755)
	v := PromptVersion{
		Version:   fmt.Sprintf("v%d", len(ListVersions(dir))+1),
		Timestamp: time.Now().Format(time.RFC3339),
		Seed:      s.Seed,
		Platform:  s.Platform,
		Prompt:    s.FinalPrompt,
	}
	data, _ := json.Marshal(v)
	fp := filepath.Join(dir, fmt.Sprintf("prompt-%s.json", v.Version))
	os.WriteFile(fp, data, 0644)
	return fp, nil
}

func ListVersions(dir string) []string {
	var versions []string
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "prompt-v") && strings.HasSuffix(e.Name(), ".json") {
			versions = append(versions, e.Name())
		}
	}
	sort.Strings(versions)
	return versions
}

// DiffVersions 比较两个版本
func DiffVersions(dir, v1, v2 string) (string, error) {
	f1 := filepath.Join(dir, fmt.Sprintf("prompt-%s.json", v1))
	f2 := filepath.Join(dir, fmt.Sprintf("prompt-%s.json", v2))
	d1, _ := os.ReadFile(f1)
	d2, _ := os.ReadFile(f2)

	var pv1, pv2 PromptVersion
	json.Unmarshal(d1, &pv1)
	json.Unmarshal(d2, &pv2)

	if pv1.Prompt == pv2.Prompt {
		return "相同 (no changes)\n", nil
	}
	return fmt.Sprintf("差异: %s → %s\n之前: %s\n之后: %s\n",
		v1, v2, pv1.Prompt, pv2.Prompt), nil
}

func (s *PromptScaffold) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	return string(data), err
}

func (s *PromptScaffold) RenderMarkdown() (string, error) {
	funcMap := template.FuncMap{"hasPrefix": hasPrefix}
	tmpl := template.Must(template.New("prompt").Funcs(funcMap).Parse(markdownTemplateV2))
	var buf strings.Builder
	err := tmpl.Execute(&buf, s)
	return buf.String(), err
}

const markdownTemplateV2 = `# 🎨 图片提示词 v2

## 来源
> {{.Seed}}

**平台**: {{.Platform}}
**格式**: {{.Model}}
{{if .RefImage}}**参考图**: {{.RefImage}}{{end}}
{{if .Params}}**参数**: {{range $k, $v := .Params}}--{{$k}} {{$v}} {{end}}{{end}}

---

## 维度控制 (16 维度)

| 分类 | 字段 | 值 | 权重 |
|------|------|-----|------|
{{range $k, $v := .Fields}}{{if not (hasPrefix $v "（LLM" )}}| {{$k}} | {{$v}} | {{if index $.Weights $k}}{{index $.Weights $k}}{{else}}1.0{{end}} |
{{end}}{{end}}

---

## 负面词

` + "```" + `
{{range .Negative}}{{.}}
{{end}}` + "```" + `

---

## 最终 Prompt

` + "```" + `
{{.FinalPrompt}}
` + "```" + `

---

*由 via54 prompt v2 生成 — 参考: designer-image (10维度) + CLIP Interrogator*
`

func buildPromptDebug(s *PromptScaffold, tmpl *PromptTemplate) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📋 平台: %s | 维度: %d\n", s.Platform, len(s.Fields)))
	for _, sec := range tmpl.Sections {
		val := s.Fields[sec.ID]
		mark := "✅"; if strings.HasPrefix(val, "（LLM") { mark = "⏳" }
		w := s.Weights[sec.ID]
		wStr := ""
		if w > 0 && w != 1.0 { wStr = fmt.Sprintf(" ×%.1f", w) }
		b.WriteString(fmt.Sprintf("  %s [%s] %s%s\n", mark, sec.ID, val[:min(len(val), 60)], wStr))
	}
	b.WriteString(fmt.Sprintf("\n最终 prompt:\n%s\n", s.FinalPrompt))
	return b.String()
}

// hasPrefix 用于模板判断
func hasPrefix(s, prefix string) bool { return strings.HasPrefix(s, prefix) }

func min(a, b int) int { if a < b { return a }; return b }
