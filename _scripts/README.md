# _scripts 目录 (★ 锂电 30s v6 视频生产管线 ★)

> v0.6.2 → v0.7.0 阶段开发, 9 版本累计
> 主项目: `G:\agent\developments\via54Design`
> 工作区 (隔离, 已汇总): `G:\agent\hermes\via54Design-v6`

## 📂 目录结构

```
_scripts/
├── gen_video.py               # ★ 主流程 ★ editly + count-up + voice + mux
├── gen_video_countup.py       # v0.7.0 ffmpeg count-up 后处理
├── gen_subtitle_png.py        # v0.6.4 字幕 PNG 生成器
├── gen_subtitle_png_v2.py     # v0.6.8 PNG 升级 (96pt 段 5)
├── gen_subtitle_countup.py    # v0.6.9 count-up 5 帧 PNG 生成
├── fetch_market_v2.py         # 段 4 视频素材 (Pexels)
├── fetch_tech_v2.py           # 段 3 视频素材
├── fetch_outlook_v2.py        # 段 5 视频素材
├── spec/                      # editly spec (5 段 3+5+10+7+5=30s)
│   ├── lithium_30s_v6_template.json5   # ★ v0.7.0 模板 ★
│   ├── lithium_30s_v6_base.json5
│   ├── lithium_30s_v6_zh.json5
│   ├── lithium_30s_v6_en.json5
│   ├── lithium_30s_v6_ja.json5
│   └── lithium_30s_v6_translated.json5
├── subtitles/                 # 5 段 × 3 语 = 15 个静态字幕 PNG
│   ├── sub_01-05_{zh,en,ja}.png
│   └── countup/               # 段 4 count-up 5 帧 × 3 语 = 15 PNG
│       └── sub_04_{zh,en,ja}_c{1-5}_*.png
└── stock/                     # 69 个 Pexels 视频素材
    ├── 01_hook_{1,2,3}.mp4
    ├── 02_trend_{1,2,3}.mp4
    ├── 03_tech_{1,2,3,battery_cell_production_1}.mp4
    ├── 04_market_{1,2,3,4,5,6}.mp4
    └── 05_outlook_{1,2,3,4,5,6}.mp4
```

## 🎬 视频输出 (★ v0.7.0 成品 ★)

主项目: `minimax-output/lithium_v7/output/`
- `lithium_30s_v7_zh.mp4` (9.3MB, 28.00s) — 锂电新纪元中文
- `lithium_30s_v7_en.mp4` (9.0MB, 28.93s) — 锂电新纪元英文 (语速差异)
- `lithium_30s_v7_ja.mp4` (9.2MB, 28.00s) — 锂电新纪元日文

## 🔄 完整流程 (★ 4 步 ★)

```bash
# 1. 拉 Pexels 视频
python fetch_market_v2.py    # 段 4 ESTUN 锂电产线
python fetch_tech_v2.py       # 段 3 MAXCELL 极片
python fetch_outlook_v2.py    # 段 5 屋顶光伏

# 2. 生成字幕 PNG (15 + 15 = 30 个)
python gen_subtitle_png_v2.py     # 静态 5 段 × 3 语
python gen_subtitle_countup.py    # count-up 5 帧 × 3 语

# 3. 主流程 (editly + count-up + voice + mux)
python gen_video.py --step all
```

## 🏆 9 版本评分

```
v0.6.4 (6.90) → v0.6.5 (7.76) → v0.6.6 (7.92) → v0.6.7 (7.02)
             → v0.6.8 (8.64 ⭐) → v0.6.9 (8.18) → v0.7.0 (8.50 ⭐ 历史新高)
```

## 📚 报告

主项目 `docs/`:
- `lithium-v6-to-v7-index.md` — 9 版本完整 INDEX
- `lithium-video-v6-translated-v8.md` — v0.7.0 详细报告
- `lithium-video-v6-translated-v7.md` — v0.6.9 报告
- `lithium-video-v6-translated-v6.md` — v0.6.8 报告
- `editly-windows-install.md` — editly 安装指南

## 🐛 已知问题

- **en.mp4 时长 28.925s** (vs zh/ja 28.000s) — 英文字幕/配音语速自然延伸
- **editly 智能合并** — 连续同 path clip 被合并, 用 ffmpeg 后处理绕过
- **mmx video quota** — 5/5 用完, ~12.6h 重置

## ⚙️ 工具链

- **Node.js 18.20.4**: `C:/Users/via54/tools/node18/node-v18.20.4-win-x64/editly.cmd`
- **Python 3.11.15**: `python` / `pip`
- **ffmpeg** + **ffprobe**
- **PIL** (Pillow): PNG 字幕生成
- **minimax CLI (mmx)**: TTS + music + image (video quota 用完)
