# 锂电 30s v6 多语翻译版报告 (v0.6.3 — 自动翻译 + TTS 旁白 + 配乐)

> 时间: 2026-06-12
> 任务: 加自动翻译 + 旁白 + 配乐, 写入 v6 默认规则
> 输出: 3 版本 mp4 (zh/en/ja) + 4 大模块流水线

## 🎯 4 大核心模块集成

```
┌─ 翻译模块 (translate.py) ──────── 关键词+旁白多语
│  ├─ 离线词典 (锂电黄金 15 词)
│  ├─ 5 段关键词 (zh/en/ja) 硬编
│  └─ 5 段旁白文案 (zh/en/ja) 硬编
│
├─ 旁白生成 (gen_voice.sh) ────── mmx speech-2.8-hd
│  ├─ 5 段 × 3 语言 = 15 段 mp3
│  ├─ zh: male-qn-jingying (精英男声)
│  ├─ en: female-yujie + 1.6x 加速
│  └─ ja: female-yujie + 2.0x 加速
│
├─ 配乐生成 (mmx music) ────── music-2.6 instrumental
│  ├─ 3 候选 (epic/tech/business)
│  ├─ 选 epic (Peak -3.9dB, 最响)
│  └─ 切 30s, -18dB 衬底
│
└─ 三路混音 (mux_v6.sh) ────── ffmpeg filter_complex
   ├─ 视频 (无音轨) + 旁白 (1.4) + 配乐 (0.13)
   ├─ duration=longest (不截视频)
   └─ AAC 192k
```

## 🎯 5 段文案 (★ 30s 黄金文案 ★)

| 段 | 时长 | 中 | 英 | 日 |
|---|---|---|---|---|
| hook | 3s | 千度熔炼, 钢铁般的决心, 锂电新纪元由此开启。 | A thousand degrees of forging, with the will of steel, a new era of lithium begins. | 千度の精錬、鋼鉄の決意、リチウム電池新時代の幕開け。 |
| trend | 5s | 2026年产能突破八百吉瓦时, 中国锂电, 领跑全球。 | In 2026, production capacity surpasses 800 GWh. China leads the world. | 2026年、生産能力が800ギガワット時を突破、中国が世界をリード。 |
| tech | 10s | 固态电池技术领跑, 安全性能全面提升, 新能源汽车的未来, 就在此刻。 | Solid-state battery leads the way, safety fully upgraded, future of EV is now. | 全固体電池リード、安全性能向上、新エネルギー自動車の未来は今。 |
| market | 7s | 板块资金流入同比增长两倍, 锂电投资窗口正在打开。 | Sector capital inflow surged 200% YoY, lithium investment window opens. | セクター資金流入2倍急増、リチウム投資の窓が開く。 |
| outlook | 5s | 从中国制造, 到中国智造, 锂电产业, 未来已来。 | From Made in China to Smart in China, lithium industry, future is here. | 中国製造から中国スマート製造へ、未来は来た。 |

## 🎯 3 版本最终输出

| 版本 | 时长 | 大小 | 码率 | 音轨 | 字幕 |
|---|---|---|---|---|---|
| **zh** | 28.00s | 16.5MB | 4.7Mbps | zh_25.31s + bgm_epic_28s | 中文 5 段 |
| **en** | 28.93s | 16.5MB | 4.6Mbps | en_28.69s + bgm_epic_28s | 中文 5 段 (英文字幕待加) |
| **ja** | 28.00s | 16.4MB | 4.7Mbps | ja_27.57s + bgm_epic_28s | 中文 5 段 (日文字幕待加) |

## 🎯 vision 关键帧评分 (中文版)

| 段 | 关键帧 | 评分 | 关键发现 |
|---|---|---|---|
| hook (1s) | 熔融金属 | 7.8 | "锂电新..." 字幕被强光吃 1 字 |
| market (16s) | CAKE K线 | 6.0 | 主题错位 (加密币), 字幕清晰 |
| outlook (21s) | 沙漠公路 | 8.5 | 投资窗口正在打开 完美收尾 |

## 🎯 流水线 1 键跑

```bash
# Step 1: 拉视频
python fetch_pexels_v3.py

# Step 2: 生成旁白
bash gen_voice.sh

# Step 3: 生成配乐 (手工 mmx music 跑 1 次, 缓存复用)
mmx music generate --prompt "Epic cinematic corporate music..." --instrumental --out music/bgm_epic.mp3
ffmpeg -i music/bgm_epic.mp3 -ss 0 -t 30 -y music/bgm_epic_30s.mp3

# Step 4: 拼接 28s 视频 (一次)
editly --json spec/lithium_30s_v6.json5  # 产出 lithium_30s_v6.mp4 28s

# Step 5: 三语混音
for L in zh en ja; do bash mux_v6.sh $L; done
```

## ⚠️ 关键约束 (★ 用户须知 ★)

| # | 约束 | 根因 | 缓解 |
|---|---|---|---|
| 1 | **段 4 主题错位** | Pexels 锂电股 K 线素材稀少, 搜"stock+candle"命中加密币 | 字幕"板块资金流入 200%"补强, 接受 5-6 分 |
| 2 | **段 1 字幕过曝** | "锂电新纪元" 4 字叠在熔融金属高光区, 1 字被吃掉 | 加描边/shadow 或改位置 |
| 3 | **英文/日文需加速** | 1.6x/2.0x 后机器感略重, 但仍在可接受范围 | 用户可调 gen_voice.sh 加速比 |
| 4 | **gl shader 转场不可用** | editly 缺 gl native binding | 用 ffmpeg xfade (fade/dissolve) 14 种 |
| 5 | **字幕没随语种改** | 字幕硬编中文 | spec 改成 `{zh\|en\|ja}` 模板, editly 渲染时插 |

## 🎯 文件清单 (v0.6.3 沉淀)

```
G:/agent/hermes/via54Design-v6/
├── translate.py              翻译模块 (142 行) — 离线词典+在线 API
├── gen_voice.sh              旁白批量生成 (15 段 mp3)
├── fetch_pexels_v3.py        视频拉取 (15 mp4, 73MB)
├── spec/
│   ├── lithium_30s_v6.json5          v6 原始版 (中文, 无音轨)
│   └── lithium_30s_v6_translated.json5  v6 多语版模板
├── voice/                    15 段 mp3 (zh/en/ja)
├── music/                    3 候选配乐 (epic/tech/business) + 3 个 30s 切片
├── output/                   4 个 mp4 (1 原始 + 3 多语混音)
└── frames/                   验证帧

G:/agent/developments/via54Design/
└── docs/
    ├── editly-windows-install.md       v0.6.2
    ├── lithium-video-v6-plan.md        v0.6.2 方案
    ├── lithium-video-v6-report.md      v0.6.2 报告
    └── lithium-video-v6-translated.md  ★ v0.6.3 多语版 (本文件)
```

## 🎯 关键决策 (★ 必看 ★)

| 决策 | 选 | 不选 | 原因 |
|---|---|---|---|
| 翻译源 | **离线黄金词典** | 在线 LibreTranslate | 0 延迟, 0 网络依赖, 锂电 15 词全覆盖 |
| 旁白加速 | **mmx --speed** | 后期 ffmpeg atempo | 一次生成对, 音质保真 |
| 配乐源 | **mmx music-2.6** | Pixabay Music | 已有授权, 30s/180s 灵活切 |
| 配乐音量 | **-18dB (0.13x)** | -10dB (0.3x) | 衬底不抢旁白, 视频完成后可调 |
| 三路混音 | **duration=longest** | shortest | 视频不被截, 配乐/旁白自然落位 |
| 视频拼接 | **editly + JSON5** | ffmpeg filter xfade | editly layer 字幕 + 多音轨支持更好 |

## 🚀 下一步 (★ 路线图 ★)

1. **字幕模板化** — `{zh|en|ja}` 占位符 + gen_video.py 渲染时替换
2. **Pexels 用更精准关键词** — 段 4 找锂电股 K 线 (试 "finance chart" + "industry")
3. **多镜头 Pexels 候选** — 每段 5+ 候选, vision 选 top 2
4. **mmx TTS 多音色** — zh 用 male-qn-badao (霸道总裁)
5. **推 GitHub** — 9 commit 等 PAT
