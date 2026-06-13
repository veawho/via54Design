#!/bin/bash
# via54Design 模式检测脚本
echo "========================================="
echo "  via54Design 模式检测"
echo "========================================="

# 1. 检测内存
RAM_GB=$(systeminfo 2>/dev/null | grep "Total Physical Memory" | awk '{print int($4/1024/1024)}' || echo "16")
echo "内存: ${RAM_GB}GB"

# 2. 检测 GPU
VRAM_GB=0
GPU_NAME="无独立显卡"
if command -v nvidia-smi &> /dev/null; then
    VRAM_MB=$(nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits 2>/dev/null | head -1 || echo "0")
    VRAM_GB=$((VRAM_MB / 1024))
    if [ $VRAM_GB -gt 0 ]; then
        GPU_NAME=$(nvidia-smi --query-gpu=name --format=csv,noheader 2>/dev/null | head -1)
    fi
fi
echo "GPU: ${GPU_NAME} (${VRAM_GB}GB VRAM)"

# 3. 推荐模式
if [ $RAM_GB -ge 16 ] && [ $VRAM_GB -ge 12 ]; then
    MODE="A"
    REASON="${RAM_GB}GB RAM + ${VRAM_GB}GB VRAM 满足全量模式"
elif [ $RAM_GB -ge 16 ]; then
    MODE="B"
    REASON="${RAM_GB}GB RAM 适合最小化模式 (via54 + 在线 API)"
else
    MODE="B"
    REASON="${RAM_GB}GB RAM 不足以支撑任何重型本地栈"
fi

echo ""
echo "推荐: 模式 ${MODE}"
echo "理由: ${REASON}"
echo ""
echo "[1] 确认  [2] 切换  [3] 跳过"
read -p "选择: " CHOICE

if [ "$CHOICE" = "1" ] || [ "$CHOICE" = "2" ]; then
    cat > .via54-mode << EOF
mode: ${MODE}
hardware:
  ram_gb: ${RAM_GB}
  vram_gb: ${VRAM_GB}
locked_at: $(date -Iseconds)
EOF
    echo "✅ 模式 ${MODE} 已锁定"
fi
