#!/usr/bin/env bash
# via54Design 跨平台编译脚本
# Usage: bash scripts/build.sh [version]
# Output: dist/via54Design-{os}-{arch}.zip

set -euo pipefail

VERSION="${1:-0.2.0}"
OUTDIR="$(dirname "$0")/../dist"
mkdir -p "$OUTDIR"

echo "=== via54Design v$VERSION 跨平台编译 ==="

build() {
    local GOOS="$1" GOARCH="$2" ext=""
    [ "$GOOS" = "windows" ] && ext=".exe"
    local NAME="via54-${GOOS}-${GOARCH}"
    local BINARY="${OUTDIR}/${NAME}/via54${ext}"

    mkdir -p "$(dirname "$BINARY")"
    echo "  → ${NAME}..."

    GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w -X main.version=$VERSION" \
        -o "$BINARY" ./cmd/huashu/

    # Copy templates + scripts
    cp -r templates "$OUTDIR/$NAME/"
    cp -r scripts "$OUTDIR/$NAME/"
    cp template-registry.yaml "$OUTDIR/$NAME/" 2>/dev/null || true
    cp template-format.md "$OUTDIR/$NAME/" 2>/dev/null || true

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
