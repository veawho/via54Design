// SPDX-License-Identifier: MIT OR AGPL-3.0

package main

// via54Design — 衍生自 huashu-design (MIT) by alchaincyf
// CLI 入口: 命令分发

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "serve":      cmdServe()
	case "generate":   cmdGenerate()
	case "narrate":    cmdNarrate()
	case "quality":    cmdQuality()
	case "pattern":    cmdPattern()
	case "list":       cmdList()
	case "media":      cmdMedia()
	case "export":     cmdExport()
	case "version":    cmdVersion()
	default:           help()
	}
}

func help() {
	fmt.Println("via54 — 设计模板引擎 + 叙事引擎 (via54Design)")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  serve             启动 MCP Server (stdio)")
	fmt.Println("  generate ...       模板组合生成HTML")
	fmt.Println("                     --layout <id> --color <id> --font <id> [--presentation]")
	fmt.Println("                     --from-narrative scaffold.json")
	fmt.Println("  narrate ...        叙事脚手架 — 一句话到故事大纲+剧本+分镜")
	fmt.Println("                     --seed \"一句话\" --model three-act [--duration 30]")
	fmt.Println("                     [--format markdown|json] [--output file]")
	fmt.Println("  quality ...        质量门禁检查")
	fmt.Println("  pattern ...        从HTML提取设计模式")
	fmt.Println("  list               列出所有可用模板")
	fmt.Println("  media ...          媒体管线 (ffmpeg + trace)")
	fmt.Println("    add-music     input.mp4 --mood=tech")
	fmt.Println("    convert       input.mp4")
	fmt.Println("    fetch         --query parrot --out ./img")
	fmt.Println("    trace         --input photo.jpg --output logo.svg")
	fmt.Println("  export ...         导出 (全 Go 管线)")
	fmt.Println("    render        input.html --format mp4/webm/hevc/frames/apng")
	fmt.Println("    pdf           input.html [--output out.pdf]")
	fmt.Println("    pptx          scaffold.json [--output deck.pptx]")
	fmt.Println("    svg           scaffold.json [--output ./scenes]")
	fmt.Println("    json          scaffold.json [--output scenes.json]")
	fmt.Println("    markdown      scaffold.json [--output slides.md]")
	fmt.Println("    tts           --text '你好' --out voice.mp3")
	fmt.Println("  version           版本信息")
	fmt.Println()
	fmt.Println("MCP: Claude Desktop / Cursor / Copilot / VS Code / Hermes")
	fmt.Println("Docs: https://github.com/veawho/via54Design")
}

func baseDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil { return dir }
	if _, err := os.Stat(filepath.Join(dir, "templates", "registry.yaml")); err == nil { return dir }
	wd, _ := os.Getwd()
	return wd
}
