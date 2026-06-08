#!/usr/bin/env bash
# via54Design — macOS 依赖安装
# 用法: bash hack/setup_macos.sh
set -euo pipefail

echo "=== via54Design macOS 依赖安装 ==="

# Homebrew check
if ! command -v brew &>/dev/null; then
  echo "正在安装 Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

echo "[1/4] Go (编译环境)..."
if ! command -v go &>/dev/null; then
  brew install go
fi
echo "  go $(go version | grep -oP 'go\S+') ✓"

echo "[2/4] Node.js + Playwright (PDF/视频导出)..."
if ! command -v node &>/dev/null; then
  brew install node
fi
echo "  node $(node --version) ✓"
if ! npx playwright --version &>/dev/null; then
  npm install -g playwright
  npx playwright install chromium
fi
echo "  playwright $(npx playwright --version 2>/dev/null) ✓"

echo "[3/4] ffmpeg (视频渲染)..."
if ! command -v ffmpeg &>/dev/null; then
  brew install ffmpeg
fi
echo "  ffmpeg $(ffmpeg -version 2>&1 | head -1) ✓"

echo "[4/4] 下载 via54Design 二进制..."
LATEST=$(curl -s https://api.github.com/repos/veawho/via54Design/releases/latest \
  | grep "browser_download_url.*darwin-arm64.zip" | cut -d'"' -f4)
if [ -n "$LATEST" ]; then
  curl -L "$LATEST" -o /tmp/via54-darwin-arm64.zip
  unzip -o /tmp/via54-darwin-arm64.zip -d /tmp/via54/
  cp /tmp/via54/via54-darwin-arm64/via54 /usr/local/bin/
  cp /tmp/via54/via54-darwin-arm64/via54-mcp /usr/local/bin/
  chmod +x /usr/local/bin/via54 /usr/local/bin/via54-mcp
  rm -rf /tmp/via54 /tmp/via54-darwin-arm64.zip
  echo "  ✅ via54 + via54-mcp 已安装到 /usr/local/bin"
else
  echo "  ⚠️ 无法获取最新版本，手动编译: make build build-mcp"
fi

echo ""
echo "=== 完成 ==="
echo "运行 via54 web 启动 Web UI"
echo "运行 via54-mcp 启动 MCP Server"
