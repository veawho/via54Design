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

