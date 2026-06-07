// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

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
