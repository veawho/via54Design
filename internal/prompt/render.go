// SPDX-License-Identifier: AGPL-3.0-only
package prompt

import (
	"encoding/json"
	"strings"
	"text/template"
)

func (s *PromptScaffold) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	return string(data), err
}

func (s *PromptScaffold) RenderMarkdown() (string, error) {
	funcMap := template.FuncMap{"hasPrefix": strings.HasPrefix}
	tmpl := template.Must(template.New("prompt").Funcs(funcMap).Parse(markdownTemplate))
	var buf strings.Builder
	err := tmpl.Execute(&buf, s)
	return buf.String(), err
}

const markdownTemplate = `# 🎨 图片提示词 v2

## 来源
> {{.Seed}}

**平台**: {{.Platform}} | **格式**: {{.Model}}
{{if .RefImage}}**参考图**: {{.RefImage}}{{end}}
{{if .Params}}**参数**: {{range $k, $v := .Params}}--{{$k}} {{$v}} {{end}}{{end}}

---

## 维度控制

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
`

func hasPrefix(s, prefix string) bool { return strings.HasPrefix(s, prefix) }
