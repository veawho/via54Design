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

	"github.com/veawho/via54Design/internal/util"
)

func main() {
	if len(os.Args) < 2 {
		cmdInteractive()
		return
	}
	// --help / -h 应成功退出 (CLI 惯例)
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		help()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "serve":
		cmdServe()
	case "interactive", "menu", "i":
		cmdInteractive()
	case "generate":
		cmdGenerate()
	case "narrate":
		cmdNarrate()
	case "quality":
		cmdQuality()
	case "pattern":
		cmdPattern()
	case "list":
		cmdList()
	case "media":
		cmdMedia()
	case "export":
		cmdExport()
	case "prompt":
		cmdPrompt()
	case "pipeline":
		cmdPipeline()
	case "present":
		cmdPresent()
	case "comfyui":
		cmdComfyUI()
	case "forge":
		cmdForge()
	case "img":
		cmdImg()
	case "web":
		cmdWeb()
	case "version":
		cmdVersion()
	default:
		help()
		os.Exit(1)
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
	fmt.Println("  forge ...           Forge Classic / A1111 生图")
	fmt.Println("  img ...             生图 (通过 mmx MiniMax CLI, 调 mmx image generate)")
	fmt.Println("  comfyui ...         ComfyUI 工作流执行 (21模板)")
	fmt.Println("  pipeline ...        LLM提示词管道 (扩展/翻译/存档)")
	fmt.Println("  web ...            Web UI (全功能控制面板)")
	fmt.Println("  interactive       交互式菜单（推荐新手）")
	fmt.Println("  version           版本信息")
	fmt.Println()
	fmt.Println("详见: https://github.com/veawho/via54Design")
}

func baseDir() string {
	return util.FindBaseDir()
}
