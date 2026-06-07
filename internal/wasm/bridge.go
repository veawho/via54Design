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

package wasm

import (
	"fmt"
	"os"
	"path/filepath"
)

type Engine struct {
	wasmPath  string
	available bool
}

func NewEngine(baseDir string) *Engine {
	candidates := []string{
		filepath.Join(baseDir, "internal", "wasm", "via54-engine.wasm"),
		filepath.Join(baseDir, "via54-engine.wasm"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return &Engine{wasmPath: p, available: true}
		}
	}
	return &Engine{available: false}
}

func (e *Engine) Available() bool { return e.available }
func (e *Engine) Path() string   { return e.wasmPath }

func (e *Engine) Compose(layoutYAML, colorYAML, fontYAML, title string) (string, error) {
	if !e.available {
		return "", fmt.Errorf("WASM not available (cd internal/wasm && bash build.sh)")
	}
	// wazero runtime integration:
	//   r := wazero.NewRuntime(ctx)
	//   compiled, _ := r.CompileModule(ctx, wasmBytes)
	//   mod, _ := r.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	//   compose := mod.ExportedFunction("compose")
	//   results, _ := compose.Call(ctx, ...)
	// See: https://github.com/tetratelabs/wazero
	return "(WASM compose - see bridge.go for wazero integration)", nil
}

func (e *Engine) MustHave() {
	if !e.available {
		fmt.Fprintf(os.Stderr, "⚠️  via54-engine.wasm not found. Build with:\n")
		fmt.Fprintf(os.Stderr, "   cd internal/wasm && bash build.sh\n")
		fmt.Fprintf(os.Stderr, "   (requires Rust: rustup target add wasm32-unknown-unknown)\n")
		fmt.Fprintf(os.Stderr, "   Falling back to native Go engine (same functionality, no WASM).\n")
	}
}
