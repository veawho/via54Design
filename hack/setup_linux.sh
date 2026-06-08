#!/usr/bin/env bash
# via54Design — Linux 依赖安装
# 用法: bash hack/setup_linux.sh
# 支持: apt (Debian/Ubuntu), dnf (Fedora), pacman (Arch)
set -euo pipefail

echo "=== via54Design Linux 依赖安装 ==="

detect_pkg_manager() {
  if command -v apt &>/dev/null; then echo "apt"
  elif command -v dnf &>/dev/null; then echo "dnf"
  elif command -v pacman &>/dev/null; then echo "pacman"
  else echo "unknown"
  fi
}

PKG=$(detect_pkg_manager)
echo "检测到包管理器: $PKG"

install_pkg() {
  local name="$1"
  case "$PKG" in
    apt)    sudo apt install -y "$name" ;;
    dnf)    sudo dnf install -y "$name" ;;
    pacman) sudo pacman -S --noconfirm "$name" ;;
  esac
}

echo "[1/3] Node.js + npm + Playwright (PDF/视频导出)..."
if ! command -v node &>/dev/null; then
  install_pkg nodejs npm
fi
echo "  node $(node --version) ✓"
if ! npx playwright --version &>/dev/null; then
  sudo npm install -g playwright
  npx playwright install chromium
fi
echo "  playwright $(npx playwright --version 2>/dev/null) ✓"

echo "[2/3] ffmpeg (视频渲染)..."
if ! command -v ffmpeg &>/dev/null; then
  install_pkg ffmpeg
fi
echo "  ffmpeg $(ffmpeg -version 2>&1 | head -1) ✓"

echo "[3/3] 下载 via54Design 二进制..."
ARCH="linux-amd64"
[ "$(uname -m)" = "aarch64" ] && ARCH="linux-arm64"
LATEST=$(curl -s https://api.github.com/repos/veawho/via54Design/releases/latest \
  | grep "browser_download_url.*${ARCH}.zip" | cut -d'"' -f4)
if [ -n "$LATEST" ]; then
  curl -L "$LATEST" -o /tmp/via54-${ARCH}.zip
  unzip -o /tmp/via54-${ARCH}.zip -d /tmp/via54/
  sudo cp /tmp/via54/via54-${ARCH}/via54 /usr/local/bin/
  sudo cp /tmp/via54/via54-${ARCH}/via54-mcp /usr/local/bin/
  chmod +x /usr/local/bin/via54 /usr/local/bin/via54-mcp
  rm -rf /tmp/via54 /tmp/via54-${ARCH}.zip
  echo "  ✅ via54 + via54-mcp 已安装到 /usr/local/bin"
else
  echo "  ⚠️ 无法获取最新版本，手动编译: make build build-mcp"
fi

echo ""
echo "=== 完成 ==="
echo "运行 via54 web 启动 Web UI"
echo "运行 via54-mcp 启动 MCP Server"
