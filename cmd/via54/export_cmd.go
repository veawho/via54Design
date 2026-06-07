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

package main

import (
	"github.com/veawho/via54Design/internal/export"
	"flag"
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func cmdExport() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: via54 export [render|pdf|tts|pptx|svg|json|markdown]")
		os.Exit(1)
	}
	switch os.Args[2] {
	case "render":
		fs := flag.NewFlagSet("render", flag.ExitOnError)
		duration := fs.Int("duration", 10, "时长(秒)")
		width := fs.Int("width", 1920, "宽")
		height := fs.Int("height", 1080, "高")
		format := fs.String("format", "mp4", "视频格式: mp4/webm/hevc/frames/apng")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.html"); os.Exit(1) }
		r, err := export.RenderVideoExt(fs.Arg(0), *duration, *width, *height, *format)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ %s: %s\n", *format, r.VideoPath)

	case "pdf":
		fs := flag.NewFlagSet("pdf", flag.ExitOnError)
		output := fs.String("output", "", "输出路径")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.html"); os.Exit(1) }
		p, err := export.ExportPDF(fs.Arg(0), *output)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ PDF: %s\n", p)

	case "tts":
		fs := flag.NewFlagSet("tts", flag.ExitOnError)
		text := fs.String("text", "", "文本")
		output := fs.String("output", "output.mp3", "输出路径")
		voice := fs.String("voice", "", "音色")
		fs.Parse(os.Args[3:])
		if *text == "" { fmt.Fprintln(os.Stderr, "请指定 --text"); os.Exit(1) }
		r, err := export.Synthesize(*text, *output, "", *voice)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ TTS: %s (%d chars)\n", r.AudioPath, r.CharCount)

	case "pptx":
		fs := flag.NewFlagSet("pptx", flag.ExitOnError)
		title := fs.String("title", "via54 演示文稿", "标题")
		output := fs.String("output", "output.pptx", "输出路径")
		widescreen := fs.Bool("16-9", true, "16:9 宽屏")
		styleID := fs.String("style", "accent-bar", "PPTX风格: accent-bar, minimal, editorial, bold")
		themeFile := fs.String("theme", "", "配色主题YAML路径 (可选)")
		fs.Parse(os.Args[3:])

		// 从叙事 JSON 或手动输入
		narrativeJSON := fs.Arg(0)
		var slides []export.PPTXSlide

		if narrativeJSON != "" {
			// 从叙事文件加载
			data, err := os.ReadFile(narrativeJSON)
			if err != nil {
				fmt.Fprintf(os.Stderr, "读取叙事文件失败: %v\n", err)
				os.Exit(1)
			}
			var sc map[string]interface{}
			if err := yaml.Unmarshal(data, &sc); err == nil {
				if beats, ok := sc["beats"].([]interface{}); ok {
					for i, b := range beats {
						bm := b.(map[string]interface{})
						act, _ := bm["act"].(string)
						vo, _ := bm["voiceover"].(string)
						mood, _ := bm["mood"].(string)
						slides = append(slides, export.PPTXSlideFromBeat(act, vo, mood, i+1, len(beats)))
					}
				}
			}
		}

		if len(slides) == 0 {
			// 无叙事文件：生成示例 slide
			slides = []export.PPTXSlide{
				export.PPTXSlideFromBeat(*title, "via54Design 生成", "inspiring", 1, 1),
			}
		}

		if err := export.ExportPPTX(slides, *output, *widescreen, *styleID, *themeFile, baseDir()); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ PPTX: %s (%d slides, 16:9=%v)\n", *output, len(slides), *widescreen)

	case "svg":
		fs := flag.NewFlagSet("svg", flag.ExitOnError)
		output := fs.String("output", "./svg-scenes", "输出目录")
		width := fs.Int("width", 1920, "宽")
		height := fs.Int("height", 1080, "高")
		fs.Parse(os.Args[3:])

		// 从叙事文件或示例
		scenes := buildScenesFromNarrative(fs.Arg(0))
		if len(scenes) == 0 {
			scenes = []export.SVGScene{{
				Title: "示例场景", Voiceover: "via54Design SVG 导出", Mood: "inspiring",
				SceneNo: 1, TotalScenes: 1, BeatName: "sample",
			}}
		}

		paths, err := export.ExportSVG(scenes, *output, *width, *height)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ SVG: %d 个文件 → %s\n", len(paths), *output)
		for _, p := range paths[:min(3, len(paths))] {
			fmt.Printf("     %s\n", p)
		}
		if len(paths) > 3 { fmt.Printf("     ... 共 %d 个\n", len(paths)) }

	case "json":
		fs := flag.NewFlagSet("json", flag.ExitOnError)
		output := fs.String("output", "scenes.json", "输出路径")
		fs.Parse(os.Args[3:])

		scenes := buildSceneDataFromNarrative(fs.Arg(0))
		if len(scenes) == 0 {
			fmt.Fprintln(os.Stderr, "请指定叙事 JSON 文件路径")
			os.Exit(1)
		}
		if err := export.ExportJSON(scenes, *output); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1)
		}
		fmt.Printf("✅ JSON: %s (%d scenes)\n", *output, len(scenes))

	case "markdown":
		fs := flag.NewFlagSet("markdown", flag.ExitOnError)
		title := fs.String("title", "via54 演示文稿", "标题")
		author := fs.String("author", "via54Design", "作者")
		output := fs.String("output", "story.md", "输出路径")
		fs.Parse(os.Args[3:])

		scenes := buildSceneDataFromNarrative(fs.Arg(0))
		if len(scenes) == 0 {
			scenes = []export.SceneData{{
				Title: "示例", Voiceover: "via54Design Markdown 导出", Mood: "inspiring",
				SceneNo: 1, TotalScenes: 1, BeatName: "sample",
			}}
		}
		if err := export.ExportMarkdown(scenes, *title, *author, *output); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1)
		}
		fmt.Printf("✅ Markdown: %s (%d slides, Marp 兼容)\n", *output, len(scenes))
		fmt.Println("   下一步: npx @marp-team/marp-cli story.md --pptx")
		fmt.Println("   或:     npx @marp-team/marp-cli story.md --pdf")

	default:
		fmt.Fprintf(os.Stderr, "未知 export 命令: %s (支持: render, pdf, tts, pptx, svg, json, markdown)\n", os.Args[2])
		os.Exit(1)
	}
}

// 辅助函数: 从叙事 JSON 构建 SVG 场景

func buildScenesFromNarrative(narrativePath string) []export.SVGScene {
	if narrativePath == "" { return nil }
	data, err := os.ReadFile(narrativePath)
	if err != nil { return nil }
	var sc map[string]interface{}
	if err := yaml.Unmarshal(data, &sc); err != nil { return nil }
	beats, ok := sc["beats"].([]interface{})
	if !ok { return nil }

	var scenes []export.SVGScene
	total := len(beats)
	for i, b := range beats {
		bm := b.(map[string]interface{})
		act, _ := bm["act"].(string)
		vo, _ := bm["voiceover"].(string)
		mood, _ := bm["mood"].(string)
		dur, _ := bm["duration"].(int)
		scenes = append(scenes, export.SVGScene{
			Title: act, Voiceover: vo, Mood: mood, BeatName: act,
			SceneNo: i + 1, TotalScenes: total, Duration: dur,
		})
	}
	return scenes
}

// 辅助函数: 从叙事 JSON 构建 SceneData

func buildSceneDataFromNarrative(narrativePath string) []export.SceneData {
	if narrativePath == "" { return nil }
	data, err := os.ReadFile(narrativePath)
	if err != nil { return nil }
	var sc map[string]interface{}
	if err := yaml.Unmarshal(data, &sc); err != nil { return nil }
	beats, ok := sc["beats"].([]interface{})
	if !ok { return nil }

	var scenes []export.SceneData
	total := len(beats)
	for i, b := range beats {
		bm := b.(map[string]interface{})
		act, _ := bm["act"].(string)
		vo, _ := bm["voiceover"].(string)
		mood, _ := bm["mood"].(string)
		st, _ := bm["start_time"].(int)
		dur, _ := bm["duration"].(int)
		scenes = append(scenes, export.SceneData{
			Title: act, Voiceover: vo, Mood: mood, BeatName: act,
			SceneNo: i + 1, TotalScenes: total, Duration: dur,
			Timing: export.SceneTiming{StartSec: st, EndSec: st + dur},
		})
	}
	return scenes
}


