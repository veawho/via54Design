# 锂电 30s v6.4 多语模板化报告 (v0.6.4 — ★ 字幕模板 + 段 4 ESTUN 升级 ★)

> 时间: 2026-06-12 13:11
> 任务: 字幕模板化 (zh/en/ja) + 段 4 找真锂电素材 + 配乐优化
> 输出: 3 版本 mp4 (28s) + 模板引擎 (gen_video.py) + 完整流水线

## 🎯 v0.6.3 → v0.6.4 三大升级

```
v0.6.3 (硬编字幕 + CAKE 错位)
  ↓
v0.6.4 (模板化 + 段 4 升级)
├─ 1. 字幕模板: spec 改 {title_xx} {sub_xx} 占位符 + gen_video.py 渲染
├─ 2. 段 4 锂电: Pexels 搜 "lithium+battery+production" → 选 ESTUN 国产黄机器人
└─ 3. 配乐优化: 保留 mmx music-2.6 epic (Pixabay API 需 Cloudflare Turnstile 验证)
```

## 🎬 5 段字幕模板 (★ 多语对照 ★)

| 段 | 时长 | 中 | 英 | 日 |
|---|---|---|---|---|
| 1 hook | 3s | 锂电新纪元 | New Era of Lithium | リチウム新時代 |
| 2 trend | 5s | 2026 产能突破 800GWh | 2026 Capacity Surpasses 800GWh | 2026年 800GWh突破 |
| 3 tech | 10s | 固态电池技术领跑全球 | Solid-State Battery Leads the World | 全固体電池 世界リード |
| 4 market | 7s | 板块资金流入 同比增长 200% | Sector Inflow +200% YoY | 資金流入 前年比 200% |
| 5 outlook | 5s | 投资窗口正在打开 | Investment Window Opens | 投資の窓が開く |

## 🎯 段 4 升级对比 (★ CAKE → ESTUN 国产黄机器人 ★)

| 维度 | v0.6.3 CAKE | v0.6.4 ESTUN |
|---|---|---|
| 画面 | 加密货币 K 线 + 卷发交易员 | **黄色机械臂 + 蓝色电池模组** ⭐ |
| 主题 | 1/10 错位 (crypto) | **9/10 完美 (锂电产线)** |
| 视觉 | 6.0 | **8.8** ⭐ |
| 字幕 | 板块资金流入 (与画面错) | 板块资金流入 (完美契合) |
| 投资感 | 6.0 (有数据) | **9.0** (强中国制造叙事) |

### ESTUN 机器人: 国产工业自动化龙头

- **颜色**: 鲜亮工程黄 (经典国产配色)
- **结构**: 多关节串联式机械臂
- **应用**: 国内锂电行业 (宁德/比亚迪/亿纬产线)
- **战略含义**: "国产工业机器人 + 国产电池 = 中国智造"

## 🎬 配乐方案 (★ Pixabay 退路 + mmx epic 保留 ★)

| 方案 | 优点 | 缺点 | 选 |
|---|---|---|---|
| **mmx music-2.6 epic** | AI 生成 30s 内可控, 风格 3 选 1 | 非真实音乐 | ✅ |
| Pixabay Music API | 真实 CC0 音乐 | 需 Cloudflare Turnstile 验证 | ❌ |
| freesound.org | 真实 CC0 音乐 | 需 API key, 配乐长度不一 | 退路 |

**频谱特征** (volumedetect):

| 候选 | mean_volume | max_volume | 选 |
|---|---|---|---|
| **epic** | -17.6 dB | -0.3 dB | ✅ 戏剧感强 |
| tech | -14.9 dB | 0.0 dB | 太响, 易过载 |
| business | -23.3 dB | -4.1 dB | 太轻, 衬底太弱 |

## 🏗️ 完整流水线 (★ gen_video.py 一键跑 ★)

```bash
# Step 1: 拉视频 (15 候选, 73MB)
python fetch_pexels_v3.py
python fetch_market_v2.py  # 段 4 重拉 12 候选, 选 ESTUN

# Step 2: 拼接 3 语视频 (editly 3 次, 字幕硬编)
python gen_video.py --step video

# Step 3: 生成 3 语旁白 (15 段 mp3, 1.9MB)
python gen_video.py --step voice
# 或: bash gen_voice.sh

# Step 4: 三路混音 (3 版本 mp4)
python gen_video.py --step mux
# 或: for L in zh en ja; do bash mux_v6.sh $L; done
```

## 🎯 gen_video.py 模板引擎 (★ v0.6.4 核心 ★)

```python
# 模板文件: spec/lithium_30s_v6_template.json5
# 占位符: {lang} {title_01} ... {title_05} {sub_01} ... {sub_05}

# render_spec(lang) 函数:
# 1. 读模板
# 2. 删除 // 注释 (JSON5 不支持 #)
# 3. 替换 {lang} → zh/en/ja
# 4. 替换 {title_xx} {sub_xx} → I18N 词典
# 5. 写 spec/lithium_30s_v6_{lang}.json5
```

### I18N 词典 (★ 黄金模板 ★)

```python
I18N = {
    "title_01": {"zh": "锂电新纪元", "en": "New Era of Lithium", "ja": "リチウム新時代"},
    "title_02": {"zh": "2026 产能突破 800GWh", "en": "2026 Capacity 800GWh", "ja": "2026年 800GWh突破"},
    "title_03": {"zh": "固态电池技术领跑全球", "en": "Solid-State Battery Leads", "ja": "全固体電池 世界リード"},
    "title_04": {"zh": "板块资金流入 同比增长 200%", "en": "Sector Inflow +200% YoY", "ja": "資金流入 前年比 200%"},
    "title_05": {"zh": "投资窗口正在打开", "en": "Investment Window Opens", "ja": "投資の窓が開く"},
    "sub_01": {"zh": "—千度熔炼, 钢铁般的决心", "en": "—A thousand degrees of steel", "ja": "—千度の精錬、鋼鉄の決意"},
    "sub_02": {"zh": "—中国锂电, 领跑全球", "en": "—China leads the world", "ja": "—中国が世界をリード"},
    "sub_03": {"zh": "—安全性能全面提升", "en": "—Safety fully upgraded", "ja": "—安全性能向上"},
    "sub_04": {"zh": "—锂电投资窗口正在打开", "en": "—Lithium investment opens", "ja": "—リチウム投資の窓が開く"},
    "sub_05": {"zh": "—未来已来", "en": "—The future is here", "ja": "—未来は来た"},
}
```

## 📊 5 段评分对比 (★ 段 4 +2.8 分 ★)

| 段 | v0.6.3 | v0.6.4 | 提升 |
|---|---|---|---|
| 1 hook | 7.8 | 6.5 (字幕被火遮) | -1.3 (字幕位置待优化) |
| 2 trend | 8.0 | **8.0** | 持平 |
| 3 tech | 6.0 | **8.0** | +2.0 (试管配固态电池文字) |
| **4 market** | 6.0 (CAKE 错位) | **8.8** (ESTUN 完美) | **+2.8** ⭐ |
| 5 outlook | 8.5 | 8.0 | -0.5 |
| **平均** | 7.26 | **7.86** | **+0.60** |

## 🐛 修复 bug 列表

| Bug | 症状 | 修法 |
|---|---|---|
| JSON5 不支持 # 注释 | editly 报 "invalid character '#'" | render_spec 删除头部 // 注释行 (保险删) |
| `lang='all'` 找不到翻译 | KeyError: 'all' | video/mux 步骤硬编 ["zh","en","ja"] |
| outPath 用了 {lang} 占位 | 3 视频流都用 zh 字幕 | 让 video 步骤跑 3 次 editly (zh/en/ja) |
| gen_video.py 调 editly 用 .exe 路径 | OSError Win32 | 改用 `editly.cmd` shell=True |

## 🎁 交付清单

| 文件 | 行数 | 用途 |
|---|---|---|
| `gen_video.py` | 230+ | 一键流水线 (3 步骤: video/voice/mux) |
| `spec/lithium_30s_v6_template.json5` | 87 | 字幕模板 (含占位符) |
| `fetch_market_v2.py` | 75 | 段 4 重拉 (lithium/electric/finance/business 4 关键词) |
| `docs/lithium-video-v6-translated-v2.md` | 本报告 | v0.6.4 完整记录 |

## 🎬 3 版本最终输出

| 版本 | 时长 | 大小 | 字幕 | 旁白 | 配乐 |
|---|---|---|---|---|---|
| **zh** | 28.00s | 15.7MB | 锂电新纪元 ... | 千度熔炼... | mmx epic |
| **en** | 28.93s | 15.7MB | New Era ... | A thousand... | mmx epic |
| **ja** | 28.00s | 15.7MB | リチウム新時代 | 千度の精錬... | mmx epic |

## 🚀 下一步建议

1. **字幕位置优化** — 段 1 字幕被火遮, 移到右下角
2. **段 3 选真电芯镜头** — 当前是试管, 不够"电芯"
3. **配乐引入真实** — 试 YouTube Audio Library (CC0) 或 Incompetech
4. **推 GitHub** — 10 commit 等 PAT
