# 部署 & 编译

本目录说明部署方式和已迁移的历史脚本。**所有导出功能已全部迁移到 Go 引擎**（`via54 export` 子命令），零外部运行时依赖。

## 构建脚本（位于 `hack/`）

| 文件 | 用途 | 语言 |
|------|------|------|
| `hack/install.sh` | 一句话部署入口 | Bash (16行) |
| `hack/setup.sh` | 自动安装依赖 + 编译 + 配置 | Bash (~180行) |
| `hack/setup_macos.sh` | macOS 一键安装 (Homebrew) | Bash (~50行) |
| `hack/setup_linux.sh` | Linux 一键安装 (apt/dnf/pacman) | Bash (~60行) |
| `hack/build.sh` | Go 跨平台编译分发 | Bash (~50行) |

## 已迁移到 Go 引擎

所有 Python/Node.js 导出脚本已删除，由 `via54 export` 纯 Go 命令替代：

| 功能 | 新命令 | 旧脚本（已删除） | 原语言 |
|------|--------|-----------------|--------|
| PPTX 导出 | `via54 export pptx scaffold.json` | `scripts/export_deck_pptx.mjs` | Node.js |
| PDF 导出 | `via54 export pdf deck.html` | `scripts/export_deck_pdf.mjs` | Node.js |
| 单文件 PDF | `via54 export pdf deck.html` | `scripts/export_deck_stage_pdf.mjs` | Node.js |
| 视频渲染 | `via54 export render --format mp4/webm/hevc/frames/apng` | `scripts/render-video.js` | Node.js |
| TTS 语音 | `via54 export tts --text "..."` | `scripts/tts-doubao.mjs` | Node.js |
| 视频转 GIF | `via54 media convert video.mp4` | `scripts/convert-formats.sh` | Bash |
| 配乐 | `via54 media add-music video.mp4 --mood=tech` | `scripts/add-music.sh` | Bash |
| 图片搜索 | `via54 media fetch --query "关键词"` | `scripts/fetch_images.py` | Python |
| 反向图片→提示词 | 待实现 (VLM) | `scripts/img2prompt.py` | Python |
| 文档→PPT 框架 | `via54 export pptx` (`/api/story2ppt`) | `scripts/doc2ppt.py` | Python |
| 故事板→视频 | `via54 export render` | `scripts/storyboard2video.py` | Python |
| LLM 提示词管道 | 不再需要 — Go 端直接处理 | `hack/via54_pipeline.py` | Python |

**当前项目唯一 Python 文件**: `test_20_rounds.py`（端到端测试套件，stdlib only）
