# via54Design — 结构化设计引擎

> 人类的灵感是弥散而跳动的，AI 是结构化的、可控的。via54Design 在这两者之间搭建桥梁。

**via54Design** 是一个 Go 语言的结构化设计引擎，将人类的"一句话灵感"转化为 AI 可执行的提示词、叙事脚手架、HTML 设计、演示文稿和视频脚本。不依赖任何图像生成后端即可独立运行。

## 核心能力

```
┌─────────────────────────────────────────────────────┐
│  Prompt 工程   叙事引擎   设计模板   媒体管线         │
│  17 平台统一    4 叙事模型   31 配色     输出全格式    │
│  26 维度控制    YAML 定义    12 字体     Go 原生实现   │
│  图片→提示词    Fountain 剧本 3 布局     零外部依赖    │
├─────────────────────────────────────────────────────┤
│  Web UI (意图驱动)        CLI (14 子命令)             │
│  API (16 端点)            ComfyUI 桥 (30 模板)       │
│  Forge 集成               MCP Server                │
└─────────────────────────────────────────────────────┘
```

---

## 🚀 快速开始

```bash
# 下载二进制
curl -L https://github.com/veawho/via54Design/releases/latest/download/via54.exe -o via54.exe

# 查看所有命令
via54

# 生成提示词
via54 prompt --scene "一只猫在月光下的屋顶上" --platform midjourney

# 生成 HTML 设计
via54 generate --layout hero-split --color ink-wash --font ming-hei-editorial --title "标题"

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

以图生提示词（需 Python + Pillow）：

```bash
python scripts/img2prompt.py photo.jpg                     # 分析→提示词
python scripts/img2prompt.py photo.jpg --desc "补充描述"   # 带用户输入
```

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
via54 generate --layout hero-split --color ink-wash --font ming-hei-editorial \
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
├── cmd/via54/          CLI 入口 (14 子命令)
├── internal/
│   ├── export/         PPTX/PDF/SVG/MD/JSON 导出 (纯 Go)
│   ├── narrate/        叙事引擎 (4 模型, YAML 驱动)
│   ├── prompt/         Prompt 工程 (17 平台, 26 维)
│   ├── pipeline/       LLM 管线编排
│   ├── workflow/       ComfyUI 工作流引擎 (30 模板)
│   ├── quality/        质量门禁
│   ├── template/       模板注册中心
│   ├── media/          媒体管线
│   └── mcp/            MCP Server
├── web/
│   ├── handler.go      API 端点 (16 个)
│   └── templates/      HTML 模板
├── scripts/
│   ├── img2prompt.py   图片分析→提示词
│   ├── doc2ppt.py      文档→PPT 框架
│   └── storyboard2video.py  多图→叙事→视频脚本
├── templates/
│   ├── color-schemes/  31 配色 YAML
│   ├── typography/     12 字体 YAML
│   ├── layouts/        3 布局 YAML
│   ├── narratology/    4 叙事模型 YAML
│   ├── pptx-styles/    4 PPTX 风格 YAML
│   └── workflows/      30 ComfyUI 工作流
└── go.mod              纯 Go，零外部运行时
```

### 设计哲学

1. **YAML 模板 > 硬编码** — 配色/字体/布局/叙事模型/PPTX 风格全是 YAML 定义
2. **纯 Go 单二进制** — 核心引擎 15MB，零运行时依赖
3. **独立运行优先** — 所有核心功能不依赖 Forge/ComfyUI
4. **确定性输出** — 所有 map 遍历前排序 key，同输入同输出
5. **API first** — 所有功能通过 REST API 可调用，Web UI 只是前端

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
