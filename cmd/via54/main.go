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
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

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
	case "serve":     cmdServe()
	case "generate":  cmdGenerate()
	case "narrate":   cmdNarrate()
	case "quality":   cmdQuality()
	case "pattern":   cmdPattern()
	case "list":      cmdList()
	case "media":     cmdMedia()
	case "export":    cmdExport()
	case "prompt":    cmdPrompt()
	case "comfyui":   cmdComfyUI()
	case "forge":     cmdForge()
	case "web":       cmdWeb()
	case "version":   cmdVersion()
	default:          help()
	}
}

func help() {
	fmt.Println("via54 — 设计模板引擎 + 叙事引擎 (via54Design)")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  serve             启动 MCP Server (stdio)")
	fmt.Println("  generate ...       模板组合生成HTML")
	fmt.Println("  narrate ...        叙事脚手架")
	fmt.Println("  quality ...        质量门禁检查")
	fmt.Println("  pattern ...        从HTML提取设计模式")
	fmt.Println("  list               列出所有可用模板")
	fmt.Println("  media ...          媒体管线")
	fmt.Println("  export ...         导出 (pptx/svg/json/markdown/video/pdf/tts)")
	fmt.Println("  prompt ...         图片提示词 (scene→MJ/Kling/即梦/Gemini)")
	fmt.Println("  web ...            Web UI (ComfyUI workflow control panel)")
	fmt.Println("  version           版本信息")
	fmt.Println()
	fmt.Println("详见: https://github.com/veawho/via54Design")
	os.Exit(1)
}

func baseDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil { return dir }
	if _, err := os.Stat(filepath.Join(dir, "templates", "registry.yaml")); err == nil { return dir }
	wd, _ := os.Getwd()
	return wd
}
