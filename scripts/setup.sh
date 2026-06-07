#!/usr/bin/env bash
# ==============================================================
# via54Design — 一键部署脚本
# 自动安装依赖、编译二进制、配置环境
# 支持: macOS / Linux / Windows (Git Bash / WSL)
# ==============================================================
set -euo pipefail

VERSION="0.3.0"
REPO="https://github.com/veawho/via54Design.git"

# ── 颜色 ──
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; NC='\033[0m'
ok()  { echo -e "  ${GREEN}✅${NC} $1"; }
warn(){ echo -e "  ${YELLOW}⚠️${NC} $1"; }
fail(){ echo -e "  ${RED}❌${NC} $1"; }
info(){ echo -e "  ${CYAN}ℹ️${NC} $1"; }

# ── 横幅 ──
echo ""
echo -e "${CYAN}╔══════════════════════════════════════╗${NC}"
echo -e "${CYAN}║      via54Design v${VERSION} 一键部署     ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════╝${NC}"
echo ""

# ── 检测 OS ──
OS="$(uname -s)"
ARCH="$(uname -m)"
case "$OS" in
  Darwin)  OS="macos"; PKG="brew"  ;;
  Linux)   OS="linux"; PKG="apt"   ;;
  MINGW*|MSYS*) OS="windows"; PKG="winget" ;;
  *)       fail "不支持的系统: $OS"; exit 1 ;;
esac
ok "系统: $OS ($ARCH)"

# ── 检测包管理器 ──
if command -v brew &>/dev/null; then PKG="brew"
elif command -v apt &>/dev/null; then PKG="apt"
elif command -v winget &>/dev/null; then PKG="winget"
elif command -v choco &>/dev/null; then PKG="choco"
elif command -v scoop &>/dev/null; then PKG="scoop"
fi
ok "包管理器: $PKG"

# ── 1. 检查/安装 Go ──
echo ""
echo -e "${YELLOW}━━━ 1/5 检查 Go ━━━${NC}"
if command -v go &>/dev/null; then
    GO_VER="$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')"
    ok "Go $GO_VER (已安装)"
else
    info "正在安装 Go..."
    case "$PKG" in
        brew)  brew install go ;;
        apt)   sudo apt install -y golang-go ;;
        winget) winget install GoLang.Go ;;
        choco) choco install golang ;;
        scoop) scoop install go ;;
        *)     fail "请手动安装 Go: https://go.dev/dl/"; exit 1 ;;
    esac
    ok "Go 安装完成"
fi

# ── 2. 检查/安装 ffmpeg ──
echo ""
echo -e "${YELLOW}━━━ 2/5 检查 ffmpeg ━━━${NC}"
if command -v ffmpeg &>/dev/null; then
    FF_VER="$(ffmpeg -version 2>&1 | head -1 | grep -oP 'ffmpeg version \K[0-9.]+')"
    ok "ffmpeg $FF_VER (已安装)"
else
    info "正在安装 ffmpeg..."
    case "$PKG" in
        brew)  brew install ffmpeg ;;
        apt)   sudo apt install -y ffmpeg ;;
        winget) winget install FFmpeg ;;
        choco) choco install ffmpeg ;;
        scoop) scoop install ffmpeg ;;
        *)     warn "请手动安装 ffmpeg: https://ffmpeg.org/download.html" ;;
    esac
    if command -v ffmpeg &>/dev/null; then ok "ffmpeg 安装完成"; fi
fi

# ── 3. 检查/安装 Node.js ──
echo ""
echo -e "${YELLOW}━━━ 3/5 检查 Node.js ━━━${NC}"
if command -v node &>/dev/null; then
    N_VER="$(node --version | grep -oP 'v\K[0-9]+')"
    ok "Node.js v$N_VER (已安装)"
else
    info "正在安装 Node.js..."
    case "$PKG" in
        brew)  brew install node ;;
        apt)   curl -fsSL https://deb.nodesource.com/setup_22.x | sudo bash - && sudo apt install -y nodejs ;;
        winget) winget install OpenJS.NodeJS.LTS ;;
        choco) choco install nodejs ;;
        scoop) scoop install nodejs ;;
        *)     warn "请手动安装 Node.js: https://nodejs.org/" ;;
    esac
    if command -v node &>/dev/null; then ok "Node.js 安装完成"; fi
fi

# ── 4. 编译 via54 ──
echo ""
echo -e "${YELLOW}━━━ 4/5 编译 via54 ━━━${NC}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"

if [ ! -f "$REPO_DIR/go.mod" ]; then
    info "克隆仓库..."
    git clone --depth=1 "$REPO" /tmp/via54Design 2>/dev/null || true
    REPO_DIR="/tmp/via54Design"
fi

cd "$REPO_DIR"
go build -o via54 ./cmd/huashu/ 2>&1
ok "编译完成: $(ls -lh via54 | awk '{print $5}')"

# ── 5. 安装 Playwright 浏览器 ──
echo ""
echo -e "${YELLOW}━━━ 5/5 配置 Playwright ━━━${NC}"
if command -v npx &>/dev/null; then
    if [ ! -d "$HOME/.cache/ms-playwright" ]; then
        info "安装 Playwright 浏览器 (首次约 2-3 分钟)..."
        npm install -g playwright 2>/dev/null || true
        npx playwright install chromium 2>&1 | tail -1
        ok "Playwright 浏览器就绪"
    else
        ok "Playwright 浏览器已安装"
    fi
else
    warn "npx 不可用，跳过 Playwright (PDF/视频导出将不可用)"
fi

# ── 添加到 PATH ──
echo ""
echo -e "${YELLOW}━━━ 配置 PATH ━━━${NC}"
INSTALL_DIR="/usr/local/bin"
if [ "$OS" = "windows" ]; then
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi
cp via54 "$INSTALL_DIR/via54" 2>/dev/null || warn "请手动将 via54 添加到 PATH"

# ── 完成 ──
echo ""
echo -e "${GREEN}╔══════════════════════════════════════╗${NC}"
echo -e "${GREEN}║       🎉 via54Design 部署完成        ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════╝${NC}"
echo ""
echo -e "  版本: ${CYAN}v${VERSION}${NC}"
echo -e "  路径: ${CYAN}$(which via54 2>/dev/null || echo "$INSTALL_DIR/via54")${NC}"
echo -e "  大小: ${CYAN}$(ls -lh via54 | awk '{print $5}')${NC}"
echo ""
echo -e "  ${YELLOW}快速开始:${NC}"
echo '    via54                                          # 帮助'
echo '    via54 list                                     # 查看模板'
echo '    via54 generate --layout hero-split --color warm-editorial --font serif-sans --title "我的设计" --output demo.html'
echo '    via54 quality --html demo.html                 # 检查质量'
echo '    via54 serve                                    # MCP Server'
echo ""
echo -e "  ${YELLOW}给 AI 助手:${NC}"
echo '    "帮我用 via54Design 生成一个暖色编辑风格的HTML页面"'
echo '    "给这个HTML做质量检查"'
echo '    "启动 MCP Server 模式"'
echo ""
