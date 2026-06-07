// SPDX-License-Identifier: MIT OR AGPL-3.0

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
