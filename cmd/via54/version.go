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
	"runtime"
	"github.com/veawho/via54Design/internal/wasm"
	"strings"
)

func cmdVersion() {
	fmt.Println("via54Design v0.3.0")
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

