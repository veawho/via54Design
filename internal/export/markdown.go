// SPDX-License-Identifier: MIT OR AGPL-3.0

// via54Design — Markdown 导出 (Marp 兼容)
// 输出可直接喂给 Marp (⭐11,917) 的 markdown 幻灯片
package export

import (
	"fmt"
	"os"
	"strings"
)

// ExportMarkdown 从场景数据生成 Marp 兼容的 markdown 幻灯片
// Marp: https://marp.app  |  GitHub: marp-team/marp (⭐11,917)
// 输出文件可直接: npx @marp-team/marp-cli story.md --pptx
// 或: npx @marp-team/marp-cli story.md --pdf
func ExportMarkdown(scenes []SceneData, title, author string, outputPath string) error {
	var b strings.Builder

	// YAML frontmatter
	b.WriteString("---\n")
	b.WriteString("marp: true\n")
	b.WriteString(fmt.Sprintf("title: %s\n", title))
	b.WriteString(fmt.Sprintf("author: %s\n", author))
	if author == "" { b.WriteString("author: via54Design\n") }
	b.WriteString("theme: uncover\n")
	b.WriteString("size: 16:9\n")
	b.WriteString("---\n\n")

	for i, s := range scenes {
		if i > 0 {
			b.WriteString("\n---\n\n")
		}

		// 环境色标记 (Marp 自定义 CSS class)
		moodClass := s.Mood
		if moodClass == "" { moodClass = "default" }

		// Beat 标签
		beatLabel := s.BeatName
		if beatLabel == "" { beatLabel = fmt.Sprintf("Scene %d", s.SceneNo) }

		b.WriteString(fmt.Sprintf("<!-- _class: %s -->\n", moodClass))
		b.WriteString(fmt.Sprintf("## %s\n\n", s.Title))

		if s.Voiceover != "" {
			b.WriteString(fmt.Sprintf("> *%s*\n\n", s.Voiceover))
		}

		if s.Body != "" {
			for _, line := range strings.Split(s.Body, "\n") {
				if strings.TrimSpace(line) != "" {
					b.WriteString(fmt.Sprintf("%s\n\n", line))
				}
			}
		}

		// 底部信息
		b.WriteString(fmt.Sprintf("<!-- \n  节拍: %s | 情绪: %s | 时长: %ds\n  页码: %d / %d\n-->\n",
			beatLabel, s.Mood, s.Duration, s.SceneNo, s.TotalScenes))
	}

	return os.WriteFile(outputPath, []byte(b.String()), 0644)
}
