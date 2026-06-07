#!/usr/bin/env bash
# via54Engine Rust → WASM 编译脚本
# 依赖: rustup target add wasm32-unknown-unknown

set -euo pipefail
cd "$(dirname "$0")"

echo "=== via54Engine WASM Build ==="

# 检查 wasm32 target
if ! rustup target list --installed | grep -q wasm32; then
    echo "正在添加 wasm32-unknown-unknown target..."
    rustup target add wasm32-unknown-unknown
fi

# 编译 WASM
echo "编译中 (release)..."
cargo build --target wasm32-unknown-unknown --release 2>&1

WASM="target/wasm32-unknown-unknown/release/via54_engine.wasm"
if [ -f "$WASM" ]; then
    # 优化
    if command -v wasm-opt &> /dev/null; then
        wasm-opt -Oz -o "via54-engine.wasm" "$WASM"
        echo "✅ wasm-opt 优化完成"
    else
        cp "$WASM" "via54-engine.wasm"
        echo "✅ (未安装 wasm-opt，使用未优化版本)"
    fi
    
    # 生成 JS 胶水
    wasm-bindgen --target web --out-dir ./pkg "via54-engine.wasm" 2>/dev/null || true
    
    SIZE=$(du -h "via54-engine.wasm" | cut -f1)
    echo "✅ via54-engine.wasm (${SIZE})"
    echo ""
    echo "用法:"
    echo "  浏览器: import init from './pkg/via54_engine.js'"
    echo "  Go:     wazero runtime (internal/wasm/bridge.go)"
    echo "  Node:   const { compose } = require('./pkg/via54_engine.js')"
else
    echo "❌ 编译失败"
    exit 1
fi
