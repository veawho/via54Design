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
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/veawho/via54Design/internal/wasm"
)

// Version 通过 ldflags 注入 (Makefile: -ldflags "-X main.Version=...")
// 默认 "dev" (未注入时) 避免暴露 v0.3.0 之类的旧硬编码
var Version = "dev"

func cmdVersion() {
	// --help / -h 在父命令层 (CLI 惯例)
	if len(os.Args) >= 3 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
		fmt.Println("用法: via54 version")
		fmt.Println()
		fmt.Println("显示当前 via54Design 版本号、Go runtime 和 WASM 引擎状态。")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --help, -h    显示本帮助")
		fmt.Println()
		fmt.Println("Version 通过 ldflags 注入 (main.Version), 未注入时为 \"dev\"。")
		os.Exit(0)
	}
	fmt.Printf("via54Design %s\n", Version)
	fmt.Printf("Go: %s %s/%s\n", strings.TrimPrefix(runtime.Version(), "go"), runtime.GOOS, runtime.GOARCH)
	we := wasm.NewEngine(baseDir())
	if we.Available() {
		fmt.Println("WASM: ✅ via54-engine loaded")
	} else {
		fmt.Println("WASM: ❌ (cd hack/wasm && bash build.sh)")
	}
	fmt.Println("Lang: Go + Rust (WASM optional)")
}

// ─── Template ───
