#!/usr/bin/env bash
# ==============================================================
# via54Design — 一句话部署入口
# AI Agent / AI IDE 一句话部署指令:
#   bash <(curl -s https://raw.githubusercontent.com/veawho/via54Design/main/scripts/install.sh)
# 或:
#   curl -fsSL https://raw.githubusercontent.com/veawho/via54Design/main/scripts/install.sh | bash
# ==============================================================
set -euo pipefail

REPO="https://github.com/veawho/via54Design.git"
TMP_DIR="/tmp/via54Design-$(date +%s)"

echo ""
echo "╔══════════════════════════════════════╗"
echo "║      via54Design 一键部署            ║"
echo "╚══════════════════════════════════════╝"
echo ""

# 1. 检查 git
if ! command -v git &>/dev/null; then
    echo "❌ 需要 git，请先安装: https://git-scm.com/"
    exit 1
fi

# 2. 克隆仓库
echo "📦 克隆仓库..."
git clone --depth=1 "$REPO" "$TMP_DIR" 2>/dev/null

# 3. 运行 setup
echo "🔧 运行部署脚本..."
bash "$TMP_DIR/scripts/setup.sh"

# 4. 清理
rm -rf "$TMP_DIR"
echo "✅ 临时文件已清理"
echo ""
echo "🎉 部署完成! 输入 via54 开始使用"
