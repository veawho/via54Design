#!/usr/bin/env bash
# via54Design 跨平台编译脚本
# Usage: bash hack/build.sh [version]
# Output: dist/via54Design-{os}-{arch}.zip (含 CLI + MCP 双二进制)

set -euo pipefail

VERSION="${1:-0.3.0}"
OUTDIR="$(dirname "$0")/../dist"
mkdir -p "$OUTDIR"

echo "=== via54Design v$VERSION 跨平台编译 ===\n"

build() {
    local GOOS="$1" GOARCH="$2" ext=""
    [ "$GOOS" = "windows" ] && ext=".exe"
    local NAME="via54-${GOOS}-${GOARCH}"
    local DIR="${OUTDIR}/${NAME}"
    mkdir -p "$DIR"

    echo "  → ${NAME}..."

    cd "$(dirname "$0")/.."

    # CLI 二进制
    GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w" \
        -o "${DIR}/via54${ext}" ./cmd/via54/

    # MCP Server 二进制
    GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w" \
        -o "${DIR}/via54-mcp${ext}" ./cmd/mcp-server/

    # Copy templates + docs
    cp -r templates "$DIR/"
    cp -r docs "$DIR/" 2>/dev/null || true
    cp README.md "$DIR/"

    # Zip
    cd "$OUTDIR"
    zip -r "${NAME}.zip" "$NAME/" > /dev/null 2>&1
    rm -rf "$NAME/"
    cd - > /dev/null

    local SIZE=$(du -h "${OUTDIR}/${NAME}.zip" | cut -f1)
    echo "    ✅ ${NAME}.zip (${SIZE})"
}

build "darwin"  "amd64"
build "darwin"  "arm64"
build "linux"   "amd64"
build "linux"   "arm64"
build "windows" "amd64"

echo ""
echo "=== 完成: $(ls -lh "$OUTDIR"/*.zip | wc -l) 个平台 ==="
ls -lh "$OUTDIR"/*.zip
