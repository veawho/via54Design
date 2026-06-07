# scripts/ — 部署 & 编译

本目录仅保留部署和编译脚本，所有导出功能已全部迁移到 Go 引擎。

| 文件 | 用途 | 语言 |
|------|------|------|
| `install.sh` | 一句话部署入口 | Bash (16行) |
| `setup.sh` | 自动安装依赖 + 编译 + 配置 | Bash (~180行) |
| `build.sh` | Go 跨平台编译分发 | Bash (~50行) |

## 已迁移到 Go 引擎

所有导出功能已由 `via54 export` 命令覆盖，零外部运行时依赖：

| 功能 | 新命令 | 旧脚本（已删除） |
|------|--------|-----------------|
| PPTX 导出 | `via54 export pptx scaffold.json` | export_deck_pptx.mjs |
| PDF 导出 | `via54 export pdf deck.html` | export_deck_pdf.mjs |
| 单文件 PDF | `via54 export pdf deck.html` | export_deck_stage_pdf.mjs |
| 视频渲染 | `via54 export render --format mp4/webm/hevc/frames/apng` | render-video.js |
| TTS 语音 | `via54 export tts --text "你好"` | tts-doubao.mjs |
| 视频转 GIF | `via54 media convert video.mp4` | convert-formats.sh |
| 配乐 | `via54 media add-music video.mp4 --mood=tech` | add-music.sh |
| 图片搜索 | `via54 media fetch --query "关键词"` | fetch_images.py |
| SVG 导出 | `via54 export svg scaffold.json` | — (新增) |
| JSON 导出 | `via54 export json scaffold.json` | — (新增) |
| Markdown | `via54 export markdown scaffold.json` | — (新增) |
