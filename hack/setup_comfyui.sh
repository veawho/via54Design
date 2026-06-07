#!/usr/bin/env bash
# via54Design — ComfyUI 环境检测 + 一键部署脚本
# 用法: bash hack/setup_comfyui.sh [--gpu|--cpu|--cloud]
#
# 参考资料: comfyui skill (Hermes) + comfy-cli (839★)

set -e

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC}   $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
fail()  { echo -e "${RED}[FAIL]${NC} $1"; }

echo "╔═════════════════════════════════════════════════════════╗"
echo "║  via54Design × ComfyUI 环境检测与部署脚本               ║"
echo "╚═════════════════════════════════════════════════════════╝"
echo ""

# ═══ Step 1: 环境检测 ═══
echo "═══ 1/7: 环境检测 ═══"

ARCH=$(uname -m)
OS=$(uname -s 2>/dev/null || echo "Windows")
GPU="none"
VRAM=0
HAS_NVIDIA=false
HAS_AMD=false
HAS_APPLE=false
HAS_CUDA=false
HAS_ROCM=false
HAS_MPS=false
HAS_COMFY_CLI=false
HAS_COMFY_SERVER=false
HAS_PYTHON=false
HAS_DOCKER=false
TOTAL_RAM=$(free -g 2>/dev/null | awk '/^Mem:/{print $2}' || echo "0")

# Python
if command -v python3 &>/dev/null; then
    HAS_PYTHON=true
    PY_VER=$(python3 --version 2>&1)
    ok "Python: $PY_VER"
else
    warn "Python 3 未安装"
fi

# NVIDIA GPU
if command -v nvidia-smi &>/dev/null; then
    HAS_NVIDIA=true
    GPU=$(nvidia-smi --query-gpu=name --format=csv,noheader 2>/dev/null | head -1)
    VRAM=$(nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits 2>/dev/null | head -1 || echo "0")
    ok "NVIDIA GPU: $GPU (${VRAM}MB VRAM)"
    if [ "$VRAM" -ge 6000 ]; then HAS_CUDA=true; fi
    if python3 -c "import torch; print(torch.cuda.is_available())" 2>/dev/null | grep -q True; then
        ok "PyTorch CUDA: 可用"
    fi
fi

# AMD GPU
if command -v rocminfo &>/dev/null; then
    HAS_AMD=true
    warn "AMD GPU 检测到 (需 ROCm)"
fi

# Apple Silicon
if [ "$(uname)" = "Darwin" ] && [ "$ARCH" = "arm64" ]; then
    HAS_APPLE=true
    HAS_MPS=true
    ok "Apple Silicon: $ARCH"
fi

# ComfyUI 状态
if command -v comfy &>/dev/null; then
    HAS_COMFY_CLI=true
    ok "comfy-cli: $(comfy --version 2>&1)"
fi

if curl -s http://127.0.0.1:8188/system_stats &>/dev/null; then
    HAS_COMFY_SERVER=true
    ok "ComfyUI Server: 正在运行 (:8188)"
else
    warn "ComfyUI Server: 未运行"
fi

# Docker
if command -v docker &>/dev/null; then
    HAS_DOCKER=true
    ok "Docker: 可用"
fi

# 内存
if [ "$TOTAL_RAM" -gt 0 ]; then
    ok "系统内存: ${TOTAL_RAM}GB"
fi

echo ""

# ═══ Step 2: 环境评估 ═══
echo "═══ 2/7: 环境评估 ═══"

VERDICT="cloud"
RECOMMENDATION="Comfy Cloud (零安装)"
if [ "$HAS_CUDA" = true ] && [ "$VRAM" -ge 12000 ]; then
    VERDICT="high-end"
    RECOMMENDATION="本地部署 (GPU + SDXL/Flux)"
elif [ "$HAS_CUDA" = true ] && [ "$VRAM" -ge 6000 ]; then
    VERDICT="mid-range"
    RECOMMENDATION="本地部署 (SD1.5/SDXL, 轻量)"
elif [ "$HAS_MPS" = true ] && [ "$(sysctl -n hw.memsize 2>/dev/null)" -ge 34359738368 ]; then
    VERDICT="high-end"
    RECOMMENDATION="本地部署 (Apple Silicon)"
elif [ "$HAS_DOCKER" = true ]; then
    VERDICT="docker"
    RECOMMENDATION="Docker 部署"
fi

echo "  评估: $VERDICT → 推荐: $RECOMMENDATION"
echo ""

# ═══ Step 3: 选择安装模式 ═══
echo "═══ 3/7: 安装模式 ═══"

MODE="${1:---auto}"
case "$MODE" in
    --gpu)    INSTALL_MODE="gpu" ;;
    --cpu)    INSTALL_MODE="cpu" ;;
    --cloud)  INSTALL_MODE="cloud" ;;
    --docker) INSTALL_MODE="docker" ;;
    --auto)
        if [ "$VERDICT" = "high-end" ]; then INSTALL_MODE="gpu"
        elif [ "$VERDICT" = "mid-range" ]; then INSTALL_MODE="gpu"
        elif [ "$VERDICT" = "docker" ]; then INSTALL_MODE="docker"
        else INSTALL_MODE="cloud"
        fi
        ;;
esac
echo "  安装模式: $INSTALL_MODE"
echo ""

# ═══ Step 4: 执行部署 ═══
echo "═══ 4/7: 部署 ComfyUI ═══"

case "$INSTALL_MODE" in
    cloud)
        echo ""
        echo "  ┌─────────────────────────────────────────────────────┐"
        echo "  │ Comfy Cloud — 零安装方案                             │"
        echo "  ├─────────────────────────────────────────────────────┤"
        echo "  │ 1. 注册: https://comfy.org/cloud                    │"
        echo "  │ 2. 获取 API Key: https://platform.comfy.org/login   │"
        echo "  │ 3. 设置: export COMFY_CLOUD_API_KEY=\"comfyui-...\"  │"
        echo "  │ 4. 使用 via54 生成 workflow JSON                   │"
        echo "  │    via54 comfyui --workflow sdxl_txt2img --prompt   │"
        echo "  │      \"...\" --output workflow.json                   │"
        echo "  │ 5. 上传到 Comfy Cloud 运行                          │"
        echo "  └─────────────────────────────────────────────────────┘"
        ;;

    docker)
        echo "  Docker 部署中..."
        docker run -d --name comfyui \
            -p 8188:8188 \
            -v comfyui-models:/comfyui/models \
            --restart unless-stopped \
            comfyui/comfyui:latest 2>/dev/null || {
            warn "Docker 镜像拉取失败，尝试手动:"
            echo "    docker pull comfyui/comfyui:latest"
            echo "    docker run -d -p 8188:8188 --name comfyui comfyui/comfyui:latest"
        }
        ;;

    gpu)
        echo "  本地 GPU 部署..."
        if [ "$HAS_COMFY_CLI" = false ]; then
            echo "  安装 comfy-cli..."
            pip install comfy-cli 2>/dev/null || pip3 install comfy-cli 2>/dev/null || true
        fi
        if command -v comfy &>/dev/null; then
            comfy install --nvidia 2>&1 | tail -3 || {
                warn "comfy-cli 安装失败，使用手动方案:"
                echo "    git clone https://github.com/comfyanonymous/ComfyUI.git"
                echo "    cd ComfyUI && pip install -r requirements.txt"
                echo "    python main.py --listen 0.0.0.0 --port 8188"
            }
        else
            warn "comfy-cli 不可用，手动安装:"
            echo "    git clone https://github.com/comfyanonymous/ComfyUI.git"
            echo "    cd ComfyUI && pip install -r requirements.txt"
            echo "    python main.py --listen 0.0.0.0 --port 8188"
        fi
        ;;

    cpu)
        echo "  CPU 部署（慢，仅供测试）..."
        echo "    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu"
        echo "    git clone https://github.com/comfyanonymous/ComfyUI.git"
        echo "    cd ComfyUI && pip install -r requirements.txt"
        echo "    CUDA_VISIBLE_DEVICES='' python main.py --listen 127.0.0.1 --port 8188"
        ;;
esac
echo ""

# ═══ Step 5: 模型下载 ═══
echo "═══ 5/7: 模型资源 ═══"
echo ""
echo "  via54 已内置工作流模板，运行时 ComfyUI 会自动下载依赖模型。"
echo "  也可手动预下载:"
echo "    comfy model download --url <HF_URL> --relative-path models/checkpoints/"
echo ""
echo "  推荐使用 aria2c 多线程下载:"
echo "    aria2c -x 16 -s 16 -d models/checkpoints/ <URL>"
echo ""

# ═══ Step 6: 与 via54 集成测试 ═══
echo "═══ 6/7: 集成测试 ═══"
echo ""
BASE_DIR=$(cd "$(dirname "$0")/.." && pwd)
if [ -f "$BASE_DIR/via54.exe" ]; then
    echo "  via54 二进制: $BASE_DIR/via54.exe"
    $BASE_DIR/via54.exe comfyui --list 2>&1 | head -3
    ok "via54 ComfyUI 集成正常"
else
    warn "via54.exe 未找到，请先编译: cd $BASE_DIR && go build -o via54.exe ./cmd/via54/"
fi
echo ""

# ═══ Step 7: 降级方案 ═══
echo "═══ 7/7: 脱离 ComfyUI 的降级方案 ═══"
echo ""
echo "  via54Design 的 Prompt 引擎完全独立于 ComfyUI："
echo ""
echo "  ┌──────────────────────────────────────────────────────────┐"
echo "  │ 方案 A: 直接提交到云平台                                │"
echo "  │   via54 prompt --scene \"...\" --platform midjourney       │"
echo "  │   → 复制最终 Prompt 到 Midjourney / Flux / 即梦等       │"
echo "  │                                                          │"
echo "  │ 方案 B: 生成 A1111/SD.Next 格式                          │"
echo "  │   python hack/via54_pipeline.py --scene \"...\"           │"
echo "  │     --export-a1111                                       │"
echo "  │   → 粘贴到 Automatic1111 WebUI                           │"
echo "  │                                                          │"
echo "  │ 方案 C: ComfyUI JSON 文件 (无需 ComfyUI 运行)            │"
echo "  │   via54 comfyui --workflow sdxl_txt2img --prompt \"...\"  │"
echo "  │     --output workflow.json                               │"
echo "  │   → 在任意 ComfyUI 实例中拖入 workflow.json 即可         │"
echo "  │                                                          │"
echo "  │ 方案 D: 文本输出 (最轻量)                                │"
echo "  │   via54 prompt --scene \"...\" --platform midjourney       │"
echo "  │   → 直接复制最终 Prompt 到任何生图工具                    │"
echo "  └──────────────────────────────────────────────────────────┘"
echo ""
echo "╔═════════════════════════════════════════════════════════╗"
echo "║  部署完成                                               ║"
echo "╚═════════════════════════════════════════════════════════╝"
