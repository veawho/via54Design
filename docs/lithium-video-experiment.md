# 锂电行业 30s 视频生成实验报告

> 时间: 2026-06-12 02:46
> 任务: 一天一个行业之锂电 — 30s 行业趋势评述
> 风格: 科技感画面 + 中文旁白 + 科技配乐, 无人物出镜
> 工具栈: mmx video + mmx speech + mmx music + ffmpeg

## 最终交付

**`minimax-output/lithium/lithium_30s_final.mp4`**
- 时长: **29.875s** (≈30s 完美)
- 分辨率: 1366×768 (16:9)
- 编码: H.264 video + AAC audio
- 码率: 1.89 Mbps
- 文件大小: 7.0 MB
- 帧数: 717 帧 @ 24fps

## 旁白稿 (24.7s, 5 句 / ~120 字)

> 锂电池产业, 正站在能源革命的临界点。
> 固态电池能量密度突破 500 瓦时每公斤,
> 钠离子电池开始进入储能市场。
> 2026 年, 全球动力电池产能预计突破 2.5 太瓦时,
> 中国企业市占率超六成。
> 技术迭代加速, 投资窗口正在打开。

## 5 段素材 (3 视频 + 1 旁白 + 1 配乐)

| # | 资源 | 命令 | 状态 | 时长/大小 |
|---|---|---|---|---|
| 1 | 视频段 1 科技发展 | `mmx video generate` | ✅ Success | 5.875s / 1.6MB |
| 2 | 视频段 2 技术突破 | `mmx video generate` | ✅ Success | 5.875s / 1.0MB |
| 3 | 视频段 3 投资趋势 | `mmx video generate` | ✅ Success | 5.875s / 2.1MB |
| 4 | 旁白 TTS | `mmx speech synthesize --voice male-qn-qingse` | ✅ Success | 24.7s / 388KB |
| 5 | 配乐 (instrumental) | `mmx music generate --instrumental --duration 30` | ✅ Success | 30s / 941KB |

## 3 段视频 prompt (英文 8 层 + 摄影词典)

### 段 1: 科技发展 (微观)
```
Microscopic cinematic view inside a lithium-ion battery.
Glowing blue lithium ions flowing through layered cathode and anode structures.
Floating electrolyte particles, electric current visualization,
blue and teal neon light trails.
Abstract scientific visualization, particle physics animation,
dark background with luminous energy flows.
Camera slowly pushes through battery layers, smooth cinematic motion,
photorealistic VFX, 8K,
no people, no text, no watermark, no logo.
```

### 段 2: 技术突破 (生产线)
```
Futuristic solid-state battery assembly line in a clean white factory.
Robotic arms precisely placing battery modules,
electric vehicle charging station with blue energy flowing,
clean energy infrastructure.
White and blue color palette, industrial cinematic lighting,
slow pan across production line, photorealistic,
no people visible, no text, no watermark.
```

### 段 3: 投资趋势 (夜景)
```
Futuristic aerial view of a massive global battery manufacturing
facility at dusk, transitioning to neon-lit cityscape at night.
Holographic data dashboard showing market statistics,
electric vehicle traffic, glowing power grid.
Cinematic drone shot pulling back,
teal and orange color palette, photorealistic,
no people, no text, no watermark, smooth motion.
```

## ffmpeg 后期流水线 (4 步)

```bash
# Step 1: 拉长每段视频到 10s (setpts × 1.702)
ffmpeg -i seg1_science.mp4 -filter:v "setpts=PTS*1.702" -r 24 -an seg1_10s.mp4
# 5.875s × 1.702 = 9.958s ≈ 10s

# Step 2: 拼接 3 段 → 29.875s
ffmpeg -i seg1_10s.mp4 -i seg2_10s.mp4 -i seg3_10s.mp4 \
  -filter_complex "[0:v][1:v][2:v]concat=n=3:v=1:a=0[outv]" \
  -map "[outv]" -c:v libx264 -preset fast -crf 22 lithium_video_raw.mp4

# Step 3: 旁白 100% + 配乐 18% 混音 (amix)
ffmpeg -i narration.mp3 -i bgm_30s.mp3 \
  -filter_complex "[0:a]volume=1.0[nar];[1:a]volume=0.18[bgm];[nar][bgm]amix=inputs=2:duration=shortest[mix]" \
  -map "[mix]" -c:a libmp3lame -b:a 192k audio_mix.mp3

# Step 4: 视频 + 音频合成 (apad 补静音到 29.875s)
ffmpeg -i lithium_video_raw.mp4 -i audio_mix.mp3 \
  -filter_complex "[1:a]apad,atrim=0:29.875[aud]" \
  -map 0:v -map "[aud]" -c:v copy -c:a aac -b:a 192k \
  lithium_30s_final.mp4
```

## 5 帧 vision 评分 (4 维度)

| 帧 | 科技感 | 画面 | 配色 | 商用 | 总分 |
|---|---|---|---|---|---|
| 0s (科技发展起) | 8.5 | 8.0 | 9.0 | 8.0 | **8.4** |
| 6s (科技发展) | — | — | — | — | (跳过) |
| 15s (技术突破末) | 8.5 | 7.5 | 9.0 | 7.0 | **8.0** |
| 22s (投资趋势中) | — | — | — | — | (跳过) |
| 30s (投资趋势末) | 8.5 | 7.5 | 9.0 | 7.0 | **8.0** |
| **平均** | 8.5 | 7.67 | 9.0 | 7.33 | **8.13** |

## ★ 沉淀到 SKILL (v0.6.1 patch 2)

### 视频 30s 标准流水线
1. **3 段 × 5.875s mmx video** = 17.6s 原始素材
2. **setpts 拉长 × 1.702** = 9.96s/段 (TPS 慢动作感)
3. **ffmpeg concat 拼接** = 29.875s (用视频为准, audio apad 补静音)
4. **mmx speech synthesize** 旁白 (24.7s, 1 句 5 短句最稳)
5. **mmx music --instrumental** 配乐 (music-2.6, 30s+)
6. **amix 混音**: 旁白 100% + 配乐 18% (-15dB 背景)
7. **apad + atrim** 音轨对齐视频 (29.875s)

### mmx 限制
- ❌ **video ≤10s/段** (Hailuo-2.3 默认 5.875s, 无 30s 单段选项)
- ❌ **video 无音频合成** (需 ffmpeg 后期)
- ❌ **video 无 --prompt-optimizer / --seed**
- ✅ **speech ≤10k chars 同步 TTS** (speech-2.8-hd)
- ⚠️ **music 输出默认 ~3 分钟** (即使 --duration 30)
- ⚠️ **mmx 文件默认名** (`speech_*.mp3` / `music_*.mp3`), 需 mv 重命名

### 30s 视频 prompt 黄金模板
- 句数: 5 短句
- 字数: 100-130 中文字符
- 节奏: 前 2 句 5s+5s (段 1), 中 2 句 5s+5s (段 2), 末 1 句 10s (段 3)
- 风格: "Investor-friendly, no jargon, data-backed"

## 引用
- 旁白稿: `docs/lithium-narrative-script.md`
- 实验报告: `docs/lithium-video-experiment.md` (本文件)
- mmx 工具: video/speech/music
- 后期工具: ffmpeg 8.0 (concat + amix + apad)
