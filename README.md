#### 致敬Aaron Swartz & Tim Berners-Lee
## 人类的灵感与创造力，因互联网与AI而永生
开发思路：让弥散而跳动的人类灵感与AI碰撞出惊艳的故事；将历久弥新的人类经验，用于规范LLM生产力的可控性。

<div align="center">

[**🇨🇳 中文**](#) · [**🇬🇧 English**](./README_EN.md)

</div>

>
> 如切如磋，如琢如磨。——《诗经·卫风·淇奥》

---

## 🧠 让时时陪伴你的AI，成为你最会讲故事的朋友
你想要的一定不是帮你把技能简单集合，更不是帮你优化工作流
**via54Design AI无法替代你的创造力，而是赋予你一双抓住灵感的妙手。**

> 文章本天成，妙手偶得之。
> 粹然无疵瑕，岂复须人为？
> —— 宋·陆游《文章》

人类的灵感是弥散而跳动的，AI 是结构化的、可控的。via54Design 在这两者之间搭建桥梁——把人类的"一句话灵感"转化为 AI 可执行的叙事脚手架，再通过模板引擎输出视频、演示文稿、创意图片。

---

### 第一部分：故事 → 视频能力

#### 从一句话到 90 秒品牌故事

**Step 1 — 人类写一句开头（人类独有的灵感）**

> "1920年代，一个中国裁缝在巴黎开了一家小店，
> 他做的旗袍融入了 Art Deco 的几何线条。
> 没有人想到，这件衣服会改变两个文明的时尚。"

这一句里有人物（裁缝）、时代（1920s）、地点（巴黎）、冲突（东方 vs 西方）、悬念（改变时尚）。AI 无法凭空创造这个种子——它来自你。

**Step 2 — AI 扩展叙事脚手架**

```bash
via54 narrate --seed "1920年代，一个中国裁缝在巴黎开了一家小店..." \
  --model heros-journey --duration 90 --format json --output scaffold.json
```

AI 分析你的种子，匹配最合适的叙事模型，输出结构化脚手架：

```
📋 英雄之旅 (Hero's Journey)  90秒
├── 第一幕·日常  (0-22s)  mood: calm      旁白: 每天，我们都...
├── 第二幕·相遇  (22-44s) mood: curious   旁白: 直到有一天...
├── 第三幕·蜕变  (44-66s) mood: excited   旁白: 不一样了...
└── 第四幕·回归  (66-90s) mood: inspiring  旁白: 每一天...

📋 分镜表: 12 个 shot（WIDE / MEDIUM / CLOSE-UP / DETAIL 循环）
📋 Fountain 剧本骨架（4 幕 8 场景）
🎞️ LLM 完整提示词（可直接喂给 Claude / GPT 生成完整剧本）
```

**Step 3 — 人类选择叙事模型，确认方向**

| 模型 | 节拍 | 适合讲什么故事 |
|------|------|---------------|
| `three-act` | 设问 → 解答 → 号召 | 产品发布、品牌广告 |
| `heros-journey` | 日常 → 相遇 → 蜕变 → 回归 | 品牌故事、纪录片 |
| `cognitive-arc` | 钩子 → 基础 → 核心 → 案例 → 延展 → 总结 | 科普、教程 |
| `problem-solution` | 痛点 → 方案 → 证明 → 行动 | 销售视频、Demo |

**Step 4 — 生成输出**

```bash
# 同一叙事脚手架，可输出多种格式
via54 export pptx scaffold.json --output story.pptx      # PPTX 演示文稿
via54 export markdown scaffold.json --output slides.md   # Marp 幻灯片
via54 export svg scaffold.json --output ./scenes          # SVG 矢量场景
via54 export json scaffold.json --output scenes.json      # 结构化数据
```

---

### 第二部分：故事 → 演示能力

同一个叙事脚手架，可以导出多种演示格式，无需重新创作。

**叙事种子 → PPTX 演示文稿**

```bash
via54 narrate --seed "1920年代，一个中国裁缝在巴黎..." \
  --model heros-journey --duration 90 --format json --output scaffold.json

via54 export pptx scaffold.json --output story.pptx --style editorial --theme templates/color-schemes/ink-wash.yaml
```

- 纯 Go 实现，零外部依赖（不依赖 Node.js / unioffice）
- 4 种风格模板（minimal / editorial / bold / accent-bar）
- 31 种配色主题，情绪映射强调色
- 16:9 宽屏，文字直接在 PPT 中可编辑

**叙事种子 → 设计生成**

```bash
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial \
  --title "裁缝的故事" --output story.html --presentation
```

---

### 第三部分：故事 → 创意图片能力

**叙事脚手架 → AI 生图提示词**

```bash
via54 prompt --scene "1920年代，巴黎左岸的小裁缝店里，一位中国裁缝在制作旗袍" \
  --platform midjourney --output prompt.md
```

输出结构化提示词，支持 17 个平台，26 维控制参数（主体/环境/光照/风格/情绪/构图/镜头/色彩/质感等）。

**叙事脚手架 → SVG 矢量场景**

```bash
via54 export svg scaffold.json --output ./scenes
```

每幕生成独立 SVG 文件，16:9 viewBox，无限缩放不失真。

**叙事脚手架 → Prompt 工作流**

```bash
via54 narrate --seed "..." --format json | via54 prompt --from-scaffold /dev/stdin
```

---

### 三种能力对比

| 能力 | 输入 | 输出 | 引擎 | 依赖 |
|------|------|------|------|------|
| 🎬 故事→视频 | 一句话 → 叙事JSON | 视频脚本 / ComfyUI 工作流 | narrate + storyboard2video | 可选 Forge/ComfyUI |
| 📊 故事→演示 | 一句话 → 叙事JSON | PPTX / Markdown / PDF / SVG | **纯 Go** | **零** |
| 🎨 故事→创意图片 | 一句话 → 叙事JSON | 结构化 Prompt (17平台) / SVG | **纯 Go** | **零** |

---

## 📊 对比参考项目

| 类别 | 参考项目 | ⭐ | via54 差异化 |
|------|---------|---|-------------|
| Prompt 工程 | easy-sd, sd-webui-forge | 10k-12k | 17 平台统一引擎，YAML 驱动，26 维控制 |
| PPT 生成 | banana-slides | 14.8k | 叙事驱动的 PPT 生成，4 种叙事模型→幻灯片 |
| 叙事引擎 | 同类项目 | <100 | 4 种正式叙事模型，Fountain 剧本+分镜表 |
| ComfyUI 管理 | ComfyUI | 116k | Go 执行桥，30 模板，确定性种子，可测试 |
| 设计模板 | huashu-design | 16.7k | Go 重写核心层，结构化 YAML 模板，质量门禁 |

---

## 🖥️ 平台支持

| 平台 | 架构 | 二进制 | Web UI | PDF/视频导出 |
|------|------|--------|:------:|:----------:|
| **macOS** Intel | amd64 | `via54-darwin-amd64` | ✅ 原生 | 需 Node.js + ffmpeg |
| **macOS** Apple Silicon | arm64 | `via54-darwin-arm64` | ✅ 原生 | 需 Node.js + ffmpeg |
| **Linux** | amd64 | `via54-linux-amd64` | ✅ 原生 | 需 Node.js + ffmpeg |
| **Linux** ARM (RPi/Graviton) | arm64 | `via54-linux-arm64` | ✅ 原生 | 需 Node.js + ffmpeg |
| **Windows** | amd64 | `via54-windows-amd64.exe` | ✅ 原生 | 需 Node.js + ffmpeg |

> **核心功能（提示词/叙事/设计/导出PPTX/JSON/SVG）100% 纯 Go，零外部依赖。**
> PDF导出和视频渲染需 Node.js + Playwright + ffmpeg —— 可选功能，Web UI 不需要。

**快速安装:**
```bash
# macOS (Intel 或 Apple Silicon)
bash hack/setup_macos.sh

# Linux (apt/dnf/pacman 自动检测)
bash hack/setup_linux.sh

# Windows (下载二进制 + 安装 Node.js)
# 下载: https://github.com/veawho/via54Design/releases
# Node.js: https://nodejs.org
```

**手动编译 (需要 Go 1.21+):**
```bash
make build          # 当前平台 CLI
make build-mcp      # MCP Server
make cross          # 5 平台交叉编译
make release        # 打包为 zip 发布包
```

**测试 (20 轮):**
```bash
python test_20_rounds.py   # 全方位压力/冒烟/稳定性/可用性/易用性/准确性测试
```

---

## 🧪 测试状态

| 维度 | 状态 | 证据 |
|------|------|------|
| **代码简洁** | ✅ 良好 | 12,254 行 Go / 58 文件, 0 panic, 179 错误检查 |
| **架构清晰** | ✅ 良好 | 双二进制 + 12 internal 包, 无环依赖 |
| **功能完整** | ✅ 良好 | 16 CLI 子命令 + 26 API endpoint + 18 MCP 工具 + 17 平台 |
| **产品质量** | ✅ 良好 | SPDX 100% (58/58), 0 FIXME, gofmt 100%, go vet 0 warnings |
| **构建测试** | ✅ 20/20 | Round 01-20 全通过 (详见 `test_20_rounds.py`) |
| **确定性输出** | ✅ 通过 | 同一输入 3 次 md5 一致 (prompt/narrate/generate) |
| **边界测试** | ✅ 通过 | 空输入/无效平台/超长字符串均优雅处理 |
| **压力测试** | ✅ 通过 | 20 次 prompt @ 23ms/次, 10 次 narrate @ 20ms/次, 5 次 generate @ 20ms/次 |
| **集成测试** | ✅ 通过 | Web UI HTTP 200, HTMX 端点返回正确碎片, 0 inline JS |

---

## 🖥️ 双模式运行配置 (v0.5.0 新增)

via54Design 支持**两种运行模式**，根据硬件自动检测并锁定。

### 模式 A: 全量运行 (≥16GB RAM + ≥12GB VRAM)

适用于配备独立 GPU 的工作站 / AI 设计师。

| 组件 | 内存 | 说明 |
|------|------|------|
| via54.exe | 20MB | 核心引擎（必装） |
| ComfyUI / Forge | 8-12GB VRAM | 本地生图（可选） |
| vtracer | 50MB | 位图→SVG 矢量化 |
| 在线 API | 0 常驻 | LLM/视频/音乐按需调用 |

**适用场景**: 专业图像生产、AI 训练、离线工作流

### 模式 B: 最小化运行 (16GB RAM, 无独显)

适用于笔记本 / 轻量办公 / 演示环境。

| 组件 | 内存 | 说明 |
|------|------|------|
| via54.exe | 20MB | 唯一本地常驻 |
| msedge.exe + opencli | 300MB (按需) | 浏览器自动化 |
| 在线 API | 0 常驻 | **全部走云端** |

**适用场景**: 海报/PPT/网页设计、临时演示、成本敏感

### 核心铁律: 两种模式都不包含本地 LLM

| 原因 | 详情 |
|------|------|
| **质量** | 7B 本地模型远不如 GPT-4/Claude/SenseNova |
| **内存** | 7B 模型常驻 4-10GB RAM 严重浪费 |
| **场景** | 设计任务以确定性模板为主，LLM 仅辅助提示词 |
| **成本** | 在线 API 按需付费优于硬件一次性投入 |

**所有 LLM 推理一律走在线 API**:
- SenseNova (`sensenova-6.7-flash-lite` / `deepseek-v4-flash`)
- OpenAI (`gpt-4o` / `dalle-3`)
- Claude (`claude-sonnet-4`)

### 模式自动检测

```bash
# 首次启动: 自动检测硬件 + 推荐模式
bash scripts/detect_mode.sh

# 输出示例:
# =========================================
#   via54Design 模式检测
# =========================================
# 内存: 16GB
# GPU: 无独立显卡
#
# 推荐: 模式 B (最小化运行)
# 理由: 16GB RAM 适合最小化模式 (via54 + 在线 API)
# =========================================
#
# [1] 确认  [2] 切换  [3] 跳过
```

### 模式锁定文件

检测完成后写入 `.via54-mode`：

```yaml
mode: B                    # A = 全量 / B = 最小化

hardware:
  ram_gb: 16
  vram_gb: 0
  gpu: "Intel UHD Graphics"

# 模式 B 配置 (16GB 内存)
local_services:
  - via54.exe (20MB)
  - msedge.exe + opencli (按需)

online_services:
  - sensenova-u1-fast (生图)
  - sensenova-6.7-flash-lite (LLM)
  - openai/dall-e-3 (生图)
  - kling-ai (生视频)
  - suno (生音乐)

# 显式禁用
disabled_services:
  - Ollama 本地 LLM
  - ComfyUI/Forge
```

### 模式对比

| 维度 | 模式 A (全量) | 模式 B (最小化) |
|------|---------------|-----------------|
| **内存需求** | 16GB RAM + 12GB VRAM | 16GB RAM |
| **本地生图** | ✅ ComfyUI/Forge | ❌ 在线 API |
| **本地矢量化** | ✅ vtracer | ✅ vtracer (可选) |
| **本地 LLM** | ❌ 在线 API | ❌ 在线 API |
| **离线能力** | 生图可离线 | 完全依赖网络 |
| **成本模型** | 硬件一次性 | 按需付费 |
| **启动时间** | 30-60s (加载模型) | <1s (仅 via54) |

**切换建议**:
- 主力做设计/AI 训练 → 模式 A
- 主力做办公/演示/轻设计 → 模式 B
- 不确定 → 先用模式 B，需求增长再升级

---

## 📜 许可

- Go 源码: `AGPL-3.0-only`
- 模板/脚本/文档: `MIT`

## 🚀 快速开始

```bash
# 下载对应平台的二进制
# macOS/Linux:
curl -L https://github.com/veawho/via54Design/releases/latest/download/via54-$(uname -s | tr A-Z a-z)-$(uname -m | sed s/x86_64/amd64/ | sed s/aarch64/arm64/).zip -o via54.zip

# 查看所有命令
via54

# 生成提示词
via54 prompt --scene "一只猫在月光下的屋顶上" --platform midjourney

# 生成 HTML 设计
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial --title "标题"

# 生成叙事脚手架
via54 narrate --seed "一个裁缝在巴黎改变了时尚" --model three-act --duration 60

# 启动 Web UI
via54 web --port 8080
```

---

## 📋 CLI 命令参考

### 🎨 Prompt 工程 — 一句话→结构化提示词

```bash
via54 prompt --scene "场景描述" --platform flux          # 文本→提示词
via54 prompt --list                                       # 查看 17 个平台
```

以图生提示词（暂未实现 - 计划中）：
- 待添加：本地 VLM 参考图片分析
- 现状：可通过 `--ref` 字段附加参考图路径到 prompt 输出中

| 平台 | 格式 | 特点 |
|------|------|------|
| `flux` / `midjourney` / `dalle3` / `sd3` / `stable_diffusion` | 英文参数 | 26 维控制参数 |
| `ideogram` / `recraft` / `seedance` | 英文参数 | 风格/文字控制 |
| `gemini` / `veo` / `sora` / `kling` / `pika` | 自然语言 | 视频生图 |
| `comfyui` / `forge` | JSON | API payload |

### 📖 叙事引擎 — 一句话→完整故事

```bash
via54 narrate --seed "一句话种子" --model three-act --duration 60   # 叙事脚手架
via54 narrate --seed "..." --format json --output scaffold.json      # JSON 输出
via54 narrate --list                                                  # 查看模型
```

| 模型 | 节拍 | 适用场景 |
|------|------|---------|
| `three-act` | 设问→解答→号召 | 产品发布、品牌广告 |
| `heros-journey` | 日常→相遇→蜕变→回归 | 品牌故事、纪录片 |
| `cognitive-arc` | 钩子→基础→核心→案例→延展→总结 | 科普、教程 |
| `problem-solution` | 痛点→方案→证明→行动 | 销售视频、Demo |

输出包含：故事大纲、节拍时间线、Fountain 剧本骨架、分镜表、LLM 完整提示词。

### 🎨 设计生成 — 布局×配色×字体→HTML

```bash
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial \
  --title "标题" --output demo.html         # 生成 HTML
via54 generate --presentation ...            # 演示文稿模式 (16:9 锁定)
```

### 📄 导出管线 — 纯 Go，零外部依赖

```bash
via54 export pptx scaffold.json              # PPTX (可编辑)
via54 export pdf story.html                  # PDF
via54 export svg scaffold.json               # SVG 场景
via54 export markdown scaffold.json          # Marp 兼容幻灯片
via54 export json scaffold.json             # 结构化数据
via54 export render story.html --duration 30 # 视频 (需 ffmpeg)
via54 export tts --text "你好"               # 语音合成
```

PPTX 支持风格模板（`--style` + `--theme`）：

```bash
# 4 种布局风格
via54 export pptx --style minimal            # 极简
via54 export pptx --style editorial          # 编辑
via54 export pptx --style bold               # 撞色
via54 export pptx --style accent-bar         # 装饰条 (默认)

# 31 种配色主题
via54 export pptx --theme templates/color-schemes/ink-wash.yaml
```

### 🔌 Forge/ComfyUI 集成

```bash
via54 forge --workflow sdxl_txt2img --prompt "a cat" --send        # 提交 Forge
via54 forge --list                                                  # 查看工作流
via54 comfyui --workflow sdxl_txt2img --prompt "a cat"              # 构建 ComfyUI JSON
```

### 🖼️ 媒体管线

```bash
via54 media add-music input.mp4 --mood tech      # 配乐
via54 media convert input.mp4                    # 格式转换
via54 media trace --input sketch.jpg             # 照片→SVG 矢量化
```

---

## 🌐 Web UI (意图驱动)

```bash
via54 web --port 8080
```

浏览器打开 `http://localhost:8080`，五个意图按钮：

| 按钮 | 功能 |
|------|------|
| 🎨 **做设计** | 描述场景→选择风格→生成 HTML + 文档/叙事→PPT 框架 |
| 📝 **写提示词** | 文本/图片→结构化 Prompt → 提交 Forge |
| 🎬 **做视频** | 上传故事板图片→叙事模型→视频脚本 + prompts |
| ⚡ **提交 Forge** | 直接提交提示词到 Forge 再生成 |

所有核心功能**独立运行**，无需 Forge/ComfyUI。后端仅增强"再生成"步骤。

API 端点（16 个）：

| 端点 | 方法 | 功能 |
|------|------|------|
| `/api/health` | GET | 健康检查 |
| `/api/templates` | GET | 工作流模板列表 |
| `/api/prompt` | POST | 生成提示词 |
| `/api/generate` | POST | 生成 HTML 设计 |
| `/api/narrate` | POST | 生成叙事脚手架 |
| `/api/build` | POST | 构建 ComfyUI 工作流 |
| `/api/export` | POST | 导出 (PPTX/PDF/SVG/MD/JSON) |
| `/api/upload` | POST | 上传图片/文档 |
| `/api/analyze` | POST | 分析图片特征 |
| `/api/img2prompt` | POST | 图片→提示词 |
| `/api/story2ppt` | POST | 文档→PPT 框架 |
| `/api/storyboard` | POST | 多图→叙事→视频脚本 |
| `/api/video-prompt` | POST | 单图→开场画面提示词 |
| `/api/regen` | POST | 提交 Forge 再生成 |

---

## 📊 模板系统

### 配色 (31 套)

| 类别 | 数量 | 示例 |
|------|------|------|
| 中国传统 | 8 | `ink-wash` 水墨, `crimson-elegance` 朱砂, `moon-white` 雨过天青 |
| 日系 | 6 | `tsubaki-camellia` 椿色, `wabi-sabi` 侘寂, `rinpa-gold` 金 |
| 经典 | 6 | `dark-terminal-blue` Linear风, `cosmic-retro` Perplexity风 |
| Adobe/Behance | 6 | `spectrum-indigo`, `glassmorphism-pastel`, `neon-dark` |
| 额外 | 5 | `bento-dark-glass`, `mono-brand-bold` 等 |

### 字体 (12 套)

中文 6 套（明体/楷体/仿宋/黑体/书法/无衬线）+ 国际 6 套（衬线/几何/展示/迪多/等宽/圆体）。

### 布局 (3 套)

| 模板 | 适配 | 特性 |
|------|------|------|
| `hero-split-16-9` | TV→Desktop→Tablet→Phone | 左右分割，响应式堆叠 |
| `bento-grid-2x2` | TV 3×2 → Desktop 2×2 → Phone 1列 | 模拟便当盒网格 |
| `gallery-waterfall` | 5列→4列→3列→2列 | 自动填充瀑布流 |

### ComfyUI 工作流 (30 套)

| 类别 | 模板 |
|------|------|
| 文生图 | sdxl_txt2img, flux_dev_txt2img, sd15_txt2img, sd3_txt2img, pixart, playground_v2, stable_cascade |
| 图生图 | sdxl_img2img, flux_img2img, sd3_img2img, sdxl_turbo, sdxl_refiner |
| 局部重绘 | sdxl_inpaint, flux_fill |
| 超分 | sdxl_upscale |
| 控制网络 | controlnet_canny, controlnet_openpose |
| 视频 | animatediff_txt2vid, hunyuan_txt2vid, wan_txt2vid, ltxv_txt2vid, mochi_txt2vid, cosmos_txt2vid, svd_img2vid, svd_img2vid_ext |
| 高级 | sdxl_advanced (LoRA+IPAdapter+FaceRestore), flux_pro (Tiled), sdxl_tiled, sdxl_img2img_face, lcm_lora |

### PPTX 风格 (4 套)

`minimal` / `editorial` / `bold` / `accent-bar` — YAML 定义位置/字体/颜色，支持 `--theme` 引用 31 配色方案。

---

## 🏗️ 架构

```
via54Design/
├── cmd/
│   ├── via54/          13 文件 — CLI 入口 (main + 10子命令)
│   └── mcp-server/     MCP Server 独立二进制
├── internal/
│   ├── export/         导出引擎 (纯Go: pptx/svg/json/markdown/pdf/tts)
│   ├── mcp/            MCP Server 实现
│   ├── media/          媒体管线 (下载/矢量化/配乐)
│   ├── narrate/        叙事引擎 (4模型, 剧本/分镜)
│   ├── pattern/        设计模式提取
│   ├── prompt/         提示词引擎 v2.2 (6模块)
│   │   ├── types.go    类型定义
│   │   ├── generator.go 生成器 (16维度)
│   │   ├── templates.go 模板加载 + 负面词库
│   │   ├── quality.go  质量评估
│   │   ├── version.go  版本管理
│   │   └── render.go   Markdown/JSON渲染
│   ├── quality/        质量门禁
│   ├── template/       模板引擎 (布局/配色/字体)
│   └── wasm/           WASM桥接 (Rust加速)
├── templates/
│   ├── prompts/        4平台提示词模板 (YAML)
│   ├── layouts/        3布局模板 (16:9, 四端适配)
│   ├── color-schemes/  30+配色方案
│   ├── typography/     12字体定义
│   ├── narratology/    4叙事模型
│   └── registry.yaml   模板注册表
├── web/                Web界面 (Go handler + HTML/JS)
│   ├── handler.go      28KB HTTP处理器
│   └── templates/      HTML模板 + JS
├── hack/               构建/部署脚本
│   ├── build.sh        Go跨平台编译 (CLI+MCP双二进制)
│   ├── install.sh      一键部署入口
│   ├── setup.sh        完整安装脚本
│   └── wasm/           Rust WASM源码
├── docs/
│   ├── prompts/        镜头/布光/配色/构图/质量参考
│   ├── template-format.md 模板格式规范
│   ├── failure-recovery.md 故障恢复指南
│   └── deployment-guide.md 部署文档
├── test_samples/       测试样本
├── AGENTS.md           AI工作上下文
├── SOUL.md             Hermes灵魂定义
├── .via54-mode         双模式锁定 (A=全量 / B=最小化)
├── scripts/
│   └── detect_mode.sh  模式自动检测脚本
├── Makefile            标准构建自动化
├── LICENSE             双许可(MIT OR AGPL-3.0)
└── README.md           项目文档
```

### 设计哲学

1. **YAML 模板 > 硬编码** — 配色/字体/布局/叙事模型/PPTX 风格全是 YAML 定义
2. **纯 Go 单二进制** — 核心引擎 15MB，零运行时依赖
3. **独立运行优先** — 所有核心功能不依赖 Forge/ComfyUI
4. **确定性输出** — 所有 map 遍历前排序 key，同输入同输出
5. **API first** — 所有功能通过 REST API 可调用，Web UI 只是前端
6. **双模式自适应** — 硬件检测 + 模式锁定 (模式A全量 / 模式B最小化)

---

## 🔗 对比参考项目

| 类别 | 参考项目 | ⭐ | via54 差异化 |
|------|---------|---|-------------|
| Prompt 工程 | easy-sd, sd-webui-forge | 10k-12k | 17 平台统一引擎，YAML 驱动，26 维控制 |
| PPT 生成 | banana-slides | 14.8k | 叙事驱动的 PPT 生成，4 种叙事模型→幻灯片 |
| 叙事引擎 | 同类项目 | <100 | 4 种正式叙事模型，Fountain 剧本+分镜表 |
| ComfyUI 管理 | ComfyUI | 116k | Go 执行桥，30 模板，确定性种子，可测试 |
| 设计模板 | huashu-design | 16.7k | Go 重写核心层，结构化 YAML 模板，质量门禁 |

---

## 📜 许可

- Go 源码: `AGPL-3.0-only`
- 模板/脚本/文档: `MIT`
- 参见: `LICENSE` 和 `ACKNOWLEDGMENTS`

---

## 🇬🇧 English Version

Read this README in English: [**README_EN.md**](./README_EN.md)


## 🔗 集成 (v0.8.0 新增)

via54Design v0.8.0 跟踪 4 个高星视频生成项目, 计划 v0.9.0 集成 diffusers, v1.0 完整 pipeline:

- [huggingface/diffusers](https://github.com/huggingface/diffusers) (33.9K) - State-of-the-art diffusion models
- [calesthio/OpenMontage](https://github.com/calesthio/OpenMontage) (27.9K) - Agentic video production (12 pipelines, 52 tools)
- [AIDC-AI/Pixelle-Video](https://github.com/AIDC-AI/Pixelle-Video) (23.8K) - AI fully automated short video engine
- [krillinai/KrillinAI](https://github.com/krillinai/KrillinAI) (10.4K) - AI video translation & dubbing

详见 [integrations/README.md](integrations/README.md) 和 [REFERENCES.md](REFERENCES.md).
