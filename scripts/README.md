# scripts/ — 遗留管线

本目录仅保留未被 Go 迁移的脚本。

## 保留

| 文件 | 原因 |
|------|------|
| `export_deck_pptx.mjs` | HTML→可编辑 PPTX，Go 无成熟替代 |
| `export_deck_pdf.mjs` | deck→PDF (配套) |
| `export_deck_stage_pdf.mjs` | deck→PDF (单文件版) |
| `build.sh` | Go 跨平台编译 |

## 已迁移到 Go

以下功能已合并到 Go 引擎：

| 旧文件 | 迁移目标 | 新命令 |
|--------|----------|--------|
| add-music.sh → | `internal/media/media.go` → | `via54 media add-music` |
| convert-formats.sh → | `internal/media/media.go` → | `via54 media convert` |
| fetch_images.py → | `internal/media/fetch.go` → | `via54 media fetch` |
| render-video.js → | `internal/export/render.go` → | `via54 export render` |
| tts-doubao.mjs → | `internal/export/tts.go` → | `via54 export tts` |
| verify.py → | `internal/quality/checker.go` → | `via54 quality` |
