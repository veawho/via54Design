# 锂电 30s v6 → v0.7.0 完整开发 INDEX (★ 一站式导航 ★)

> 时间: 2026-06-09 18:16 → 2026-06-12 15:30 (昨晚 → 今天)
> 路径: `G:\agent\developments\via54Design`
> 工作区: `G:\agent\hermes\via54Design-v6` (隔离, **已汇总资产到主项目 2026-06-12 15:25**)

## 🎯 9 版本完整时间线

```
v0.6.2 (commit c046714) → v0.6.3 (b955b92) → v0.6.4 (44e596f) → v0.6.5 (61d1029) →
v0.6.6 (680f581) → v0.6.7 (a5d086f) → v0.6.8 (aac9b6c) → v0.6.9 (3f4937d) → v0.7.0 (572b344)
```

## 📂 主项目资产 (★ 2026-06-12 15:25 汇总后 ★)

### 视频成品 (`minimax-output/lithium_v7/output/`)
- `lithium_30s_v7_zh.mp4` (9.3MB, 28.00s, count-up + 原子能量)
- `lithium_30s_v7_en.mp4` (9.0MB, 28.93s, count-up + Atomic Power)
- `lithium_30s_v7_ja.mp4` (9.2MB, 28.00s, count-up + 原子の力)
- `lithium_30s_v7_base.mp4` (8.9MB, 28.00s, 无音轨)

### 脚本 (`_scripts/`)
- `gen_video.py` — 主流程 (editly + count-up + voice + mux)
- `gen_video_countup.py` — ffmpeg count-up 后处理 (★ v0.7.0 ★)
- `gen_subtitle_png_v2.py` — PNG 字幕生成器 (96pt 升级)
- `gen_subtitle_countup.py` — count-up 5 帧 PNG 生成
- `fetch_market_v2.py` / `fetch_tech_v2.py` / `fetch_outlook_v2.py` — 视频素材

### Spec (`_scripts/spec/`)
- `lithium_30s_v6_template.json5` — ★ v0.7.0 模板 ★ (5 段 3+5+10+7+5=30s, 1366x768@24fps)
- `lithium_30s_v6_base.json5` — 基础 spec
- `lithium_30s_v6_zh.json5` / `_en.json5` / `_ja.json5` — 渲染时 3 语
- `lithium_30s_v6_translated.json5` — 翻译版

### 字幕 PNG (`_scripts/subtitles/`)
- `sub_01-05_{zh,en,ja}.png` — 15 个静态字幕 (96pt 段 5)
- `countup/sub_04_{zh,en,ja}_c{1-5}_*.png` — 15 个 count-up 数字帧

### 视频素材 (`_scripts/stock/`)
- 69 个 mp4 (Pexels 拉的真实视频: 工业/制造/能源/科技)

### 报告 (`docs/`)
| 文件 | 字节 | 时间 | 内容 |
|---|---|---|---|
| `lithium-narrative-script.md` | 1917 | 02:39 | 叙事脚本 v3 |
| `lithium-storyboard-v4.md` | 9052 | 03:09 | 分镜 v4 |
| `lithium-video-experiment.md` | 5619 | 03:01 | 实验 v3 |
| `lithium-video-v4-report.md` | 5551 | 03:17 | v4 报告 |
| `lithium-video-v5-report.md` | 5717 | 11:31 | v5 报告 |
| `lithium-video-v6-plan.md` | 5681 | 11:36 | v6 计划 |
| `lithium-video-v6-report.md` | 270 | 12:20 | v6 报告 (短) |
| `lithium-video-v6-translated.md` | 6637 | 12:54 | v6 翻译 v1 |
| `lithium-video-v6-translated-v2.md` | 6840 | 13:20 | v6 翻译 v2 |
| `lithium-video-v6-translated-v3.md` | 6211 | 13:42 | v6 翻译 v3 |
| `lithium-video-v6-translated-v4.md` | 4358 | 14:10 | v6 翻译 v4 |
| `lithium-video-v6-translated-v5.md` | 4221 | 14:20 | v6 翻译 v5 |
| `lithium-video-v6-translated-v6.md` | 4106 | 14:28 | v6 翻译 v6 (v0.6.8) |
| `lithium-video-v6-translated-v7.md` | 4522 | 14:42 | v6 翻译 v7 (v0.6.9) |
| `lithium-video-v6-translated-v8.md` | 5860 | 15:05 | v6 翻译 v8 (v0.7.0) |
| `prompt-mastery-v3.md` | 9188 | 02:20 | 提示词 v3 |
| `prompt-v3-experiment-report.md` | 3992 | 02:30 | 提示词实验 |
| `prompt-debug-zh.md` | 8011 | 02:04 | 提示词 debug |
| `editly-windows-install.md` | 4980 | 12:00 | editly 安装 |

## 📊 9 版本评分 (5 段平均)

```
v0.6.2 → v0.6.3 → v0.6.4 → v0.6.5 → v0.6.6 → v0.6.7 → v0.6.8 → v0.6.9 → v0.7.0
        (mmx)   (editly) (字幕)  (真电芯) (屋顶)   (PNG)    (全程)  (96pt)  (countup)
6.90  →  ?     →  6.90  →  7.76  →  7.92  →  7.02  →  8.64  →  8.18  →  8.50 ⭐
                                                          ↑ v0.6.8 历史最高
                                                                  ↑ v0.7.0 新高
```

## 🔧 隔离工作区 → 主项目汇总 (2026-06-12 15:25)

**前**: `G:\agent\hermes\via54Design-v6\` 是 v0.6.4-v0.7.0 开发隔离区
**后**: 全部资产已复制到主项目 `G:\agent\developments\via54Design\`
- ✅ 视频: `minimax-output/lithium_v7/output/` (4 个 mp4)
- ✅ 脚本: `_scripts/` (8 个 .py)
- ✅ 字幕: `_scripts/subtitles/` (15 个静态 + 15 个 countup)
- ✅ Spec: `_scripts/spec/` (6 个 json5)
- ✅ 素材: `_scripts/stock/` (69 个 mp4)

## 🐛 v0.6.4-v0.7.0 关键 bug 与修复

1. **editly 智能合并**: 连续同 path clip 合并 → ffmpeg 后处理解决
2. **count-up 时间错位**: 16-23s → 18-25s (匹配段 4 实际区间)
3. **mux 用统一视频源覆盖 en/ja**: 按语言用对应 count-up mp4
4. **ffmpeg 原地覆盖失败**: 输出 .tmp.mp4 → shutil.move
5. **路径 `c5_95%+.png` 不存在**: 实际是 `c5_95%p.png` (PIL 自动转义)
6. **filter_complex 链式叠加错误**: 改成 `[v_n][n+1:v]overlay` 正确语法

## 🏆 v0.7.0 最终交付 (★ 4 创新 ★)

1. **count-up 数字动画**: 段 4 0%→25%→50%→75%→95%+ 流畅切换 (5 PNG + ffmpeg chain overlay)
2. **段 1 改文案**: "锂电新纪元" → "原子能量" (消解 vision 隐喻脱节扣分)
3. **段 5 96pt 升级**: 80pt → 96pt 收尾压轴
4. **3 语 count-up 验证**: zh/en/ja 全部跑通
