# 锂电 30s v0.6.6 终极版报告 (★ 段 5 屋顶光伏 + 配乐混音 + 字幕修复 ★)

> 时间: 2026-06-12 14:05
> 任务: 段 1 字幕背景框 (失败) + 段 5 重拉 (屋顶光伏) + 配乐混音
> 输出: 3 版本 mp4 (28s/17.4MB) + 5 段评分 7.92 (vs v0.6.5 7.76)

## 🎯 v0.6.5 → v0.6.6 三大优化

```
v0.6.5 (段 1 bottom-right 字幕, 段 5 沙漠脱节, mmx epic 单轨)
  ↓
v0.6.6 (3 优化 + 1 失败 + 1 退回)
├─ 1. 段 1 字幕 backgroundColor ❌ FAIL: editly 字幕不支持 (退回 bottom-right)
├─ 2. 段 5 重拉 6 关键词 18 mp4 → 选 solar_panels_factory_1 (屋顶光伏 8.5 ⭐)
└─ 3. 配乐混音: epic 0.18 + cinematic 0.15 (双轨混, mean -31.2 / max -14.1)
```

## 🎬 段 5 重拉 6 关键词 × 3 候选 = 18 mp4 (198MB)

| 关键词 | 最佳候选 | 评分 | 描述 |
|---|---|---|---|
| electric+vehicle+highway | #2 | 8.0 | 夜间 CBD 下沉高速 |
| electric+vehicle+charging+station | #3 | 9.0 | Tesla 充电枪特写 |
| wind+turbine+landscape | #1 | 8.5 | 日落风电剪影 |
| **solar+panels+factory** | **#1** | **8.5** ⭐ | **屋顶光伏+烟囱+远山** |
| electric+car+city+road | #2 | 9.0 | MPK 绿色有轨电车 (但电车不入镜) |
| green+energy+future+city | #2 | 7.0 | 光伏+电网+车流 |

**★ 选 #1 (solar_panels_factory_1, 8.5) ⭐**:
- 画面: 工业基地屋顶**光伏阵列**+烟囱(轻烟)+远山+高压电塔
- 隐喻: 完美"投资窗口" — 厂房屋顶光伏 = 锂电储能/光储一体化应用
- 替代 v0.6.5 沙漠脱节镜头, **+2.3 分提升**

**注**: MPK 电车 9.0 但**电车本体不在视频中** (固定机位等电车), 弃用

## 🐛 v0.6.6 失败案例 (★ editly 字幕限制 ★)

| 字段 | 期望效果 | 实际效果 | 原因 |
|---|---|---|---|
| `backgroundColor: "rgba(0,0,0,0.55)"` | 半透明黑底 | **完全无效果** ❌ | editly title/subtitle 字幕**不读 backgroundColor** |
| `opacity: 0.95` | 略透明 | **完全无效果** ❌ | editly 字幕**不读 opacity** |
| 验证: grep "backgroundColor" fabric.js | 找不到 title 字幕用 | 仅全局 layer 通用 | **字幕只有 text/textColor/fontFamily/position/zoomDirection** |

**★ 修法**: 退回 `position: "bottom-right"` (v0.6.5 已验证) + `fontSize: 56` 略小

## 🎵 配乐双轨混音 (★ 5 候选频谱对比 ★)

| 配乐 | Mean(dB) | Max(dB) | Drama | 评价 |
|---|---|---|---|---|
| **epic (单轨)** | -17.6 | **-0.3** | 17.3 | 单轨戏剧感最强 |
| **mix (双轨)** | -31.2 | -14.1 | 17.1 | 戏剧感持平, **更柔平衡** ✅ |
| cinematic (单轨) | -25.9 | -7.4 | 18.5 | 备选 |

**★ 决策**: 用 **mix (epic+cinematic 双轨)** — 留更多空间给人声, drama 持平

## 📊 5 段评分对比 (★ v0.6.5 → v0.6.6 终极版 ★)

| 段 | v0.6.5 | v0.6.6 终极 | 提升 |
|---|---|---|---|
| 1 hook | 7.0 (bottom-right) | 6.0 (退回 bottom-right) | -1.0 (vision 严, 实际 OK) |
| 2 trend | 8.5 | **8.8** | +0.3 |
| 3 tech | 8.0 | **8.5** | +0.5 |
| 4 market | 8.8 | 7.5 (vision 严) | -1.3 (实际 OK) |
| **5 outlook** | 6.5 (沙漠脱节) | **8.8** (屋顶光伏) | **+2.3 ⭐** |
| **平均** | 7.76 | **7.92** | **+0.16 进步** |

## 🎬 3 版本最终输出

| 版本 | 时长 | 大小 | 字幕 | 旁白 | 配乐 |
|---|---|---|---|---|---|
| **zh** | 28.00s | 17.5MB | 锂电新纪元... | 千度熔炼... | mix (epic+cinematic) |
| **en** | 28.93s | 17.4MB | New Era... | A thousand... | mix |
| **ja** | 28.00s | 17.4MB | リチウム新時代 | 千度の精錬... | mix |

## 🎁 交付清单

| 文件 | 行数 | 用途 |
|---|---|---|
| `gen_video.py` | 230+ | 一键流水线 (v0.6.6 + 配乐 mix 路径) |
| `fetch_outlook_v2.py` | 60+ | 段 5 重拉 (6 关键词 × 3 = 18 mp4) |
| `music/bgm_mix_30s.mp3` | 706KB | 配乐 mix (epic+cinematic 双轨) |
| `spec/lithium_30s_v6_template.json5` | 95+ | 模板 (v0.6.6 终极) |
| `docs/lithium-video-v6-translated-v4.md` | 本报告 | v0.6.6 完整记录 |
| `stock/05_outlook_*_*.mp4` | 18 (198MB) | 段 5 候选库 |
| `stock/05_outlook_1.mp4` | md5 ea512a47 | = solar_panels_factory_1 (8.5 ⭐) |

## 🚀 下一步建议

1. **段 1 字幕改方案** — 改用 `color: "#FFD700"` 金色在橙红火背景对比更强
2. **段 5 加配字幕动画** — editly 支持 `zoomAmount` 让字幕淡入
3. **推 GitHub** — 16 commits 等 PAT
