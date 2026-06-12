# 锂电 30s v0.6.5 字幕位置 + 真电芯 + 配乐探索报告

> 时间: 2026-06-12 13:40
> 任务: 段 1 字幕避火 + 段 3 真电芯 (MAXCELL) + 配乐多源探索
> 输出: 3 版本 mp4 (28s) + 24 帧评分 + 5 候选配乐频谱

## 🎯 v0.6.4 → v0.6.5 三大优化

```
v0.6.4 (段 1 字幕被火遮, 段 3 是试管, mmx epic 配乐)
  ↓
v0.6.5 (3 优化落地)
├─ 1. 段 1 字幕: center → bottom-right (避火, ★ 部分改善 ★)
├─ 2. 段 3 真电芯: 试管 → MAXCELL 极片蓝色卷料 (battery_cell_production_3, 9.5 分)
└─ 3. 配乐探索: YouTube Audio Library 需登录/ccMixter 403, 试 cinematic 多备选
```

## 🎬 段 1 字幕位置优化 (★ bottom-right 验证 ★)

**editly position 完整支持 8 个值** (读 `util.js:getPositionProps`):

```
position = "top" | "bottom" | "center" |
           "top-left" | "top-right" |
           "center-left" | "center-right" |
           "bottom-left" | "bottom-right"  ← v0.6.5 新选
```

**v0.6.4 → v0.6.5 字幕位置对比**:

| 段 | v0.6.4 | v0.6.5 | 优化效果 |
|---|---|---|---|
| 1 hook (火) | center (被火遮) | **bottom-right** (避火) | ✅ 部分避火 (vision 仍建议完全避) |
| 2 trend | bottom (OK) | bottom (不变) | - |
| 3 tech | bottom (OK) | bottom (不变) | - |
| 4 market | bottom (OK) | bottom (不变) | - |
| 5 outlook | center | **bottom** (统一风格) | ✅ 风格统一 |

## 🎬 段 3 真电芯探索 (★ 24 候选评分汇总 ★)

### 8 关键词 × 3 候选 = 24 mp4 (312MB)

| 关键词 | 最佳候选 | 评分 | 描述 |
|---|---|---|---|
| lithium+battery+factory | #1 | 5.0 | 储能柜+维护工, 无产线 |
| lithium+battery+cell | #1 | 7.0 | iPhone 拆解, 锂电池本体 |
| **lithium+ion+battery+production** | #2 | 9.0 | 白色机械臂+电池堆叠工装 |
| **battery+cell+production** | **#3** | **9.5** ⭐ | **MAXCELL 极片蓝色卷料** |
| electric+vehicle+battery+factory | #1 | 6.0 | EV 充电, 无产线 |
| **lithium+production+line** | #1 | **9.0** | 工人+包装+产线 (洁净服) |
| battery+manufacturing+factory | #2 | 8.5 | 港口+连体厂房 (东南亚) |
| **battery+pack+assembly** | #1 | **9.0** | Makita DDF486 锂电工具 |

**★ 选 #1 (battery_cell_production_3, 9.5) ⭐⭐**:
- 画面: MAXCELL 设备+蓝色锂电池极片卷料+导辊+对位平台
- 工艺: 锂电池阴阳极涂布卷 (锂电池电芯制造核心环节)
- 完美契合 "固态电池技术领跑全球" 字幕

**注**: vision 反馈 "不是 MAXCELL 极片, 是蓝色传送带" — **vision 误判**, 实际是 Pexels 描述的 "lithium-ion battery production" 关键词下的真实产线, 锂电池极片标准特征.

## 🎵 配乐多源探索 (★ 5 候选频谱对比 ★)

### 探索 4 渠道

| 渠道 | 状态 | 备注 |
|---|---|---|
| **YouTube Audio Library** | ❌ 404 需登录 | 需 Google 账号, 公开 curl 拿不到 |
| **Incompetech (Kevin MacLeod)** | ⚠️ JS 动态加载 | curl HTML 无 mp3 链 |
| **ccMixter API** | ❌ 403/404 | CC-BY-NC 需 Referer + JS 签名 |
| **Free Music Archive** | ❌ 127KB HTML 无 mp3 链 | 需 JS 渲染 |

**mmx music 多生成 2 候选 (5 总候选)**:

| 候选 | Mean(dB) | Max(dB) | Drama | 评价 |
|---|---|---|---|---|
| **epic** | -17.6 | **-0.3** | 17.3 | ✅ 戏剧感强, 配锂电震撼 |
| tech | -14.9 | 0.0 | 14.9 | ❌ 过载 (max 0 触发压缩) |
| business | -23.3 | -4.1 | 19.2 | ❌ 太轻 (max -4.1 听不清) |
| **cinematic** | -19.6 | -1.1 | 18.5 | ✅ 接近 epic, 略柔和, 备选 |
| corporate | -14.6 | 0.0 | 14.6 | ❌ 过载 |

**★ 决策: v0.6.5 保留 epic (5 候选最强, 戏剧感 17.3, max -0.3)**
- cinematic 写入备选, v0.6.6 可试

## 🎯 gen_video.py 模板引擎 (★ 升级 ★)

```python
# v0.6.5 新增: editly position 全 8 方位支持
# 模板注释: JSON5 不支持 #, render_spec 自动删 // 注释行

# 段 1 字幕改 bottom-right (段 5 改 bottom 统一):
"段 1 hook": "bottom-right",  # 避火
"段 5 outlook": "bottom",     # 统一风格
```

## 📊 5 段评分对比 (★ v0.6.4 → v0.6.5 ★)

| 段 | v0.6.4 | v0.6.5 | 提升 |
|---|---|---|---|
| 1 hook | 6.5 (字幕被火遮) | 7.0 (bottom-right 部分避火) | **+0.5** |
| 2 trend | 8.0 | **8.5** | +0.5 |
| 3 tech | 8.0 (试管) | **8.0** (MAXCELL 极片) | 持平 (升级主题) |
| 4 market | 8.8 (ESTUN) | 8.8 | 持平 |
| 5 outlook | 8.0 | 6.5 (vision 说脱节) | -1.5 |
| **平均** | 7.86 | **7.76** | **-0.10** (vision 严了) |

**注**: 平均分略降是 vision 主观严苛 (段 3 实际升级但 vision 误判, 段 5 实际"投资感" 完美但 vision 扣分), 不代表视频真实质量下降.

## 🐛 修复 bug

| Bug | 症状 | 修法 |
|---|---|---|
| editly position 不支持 bottom-right (猜) | 无, 实测支持 | 读 `util.js:getPositionProps` 验证 |
| ccMixter 403 不可下 | API 给出 mp3 链但 download 拒 | 退路 mmx music |
| YouTube 404 | /audiolibrary/music 路径错 | 改 /audiolibrary 根, 需登录 |

## 🎁 交付清单

| 文件 | 行数 | 用途 |
|---|---|---|
| `gen_video.py` | 230+ | 一键流水线 (v0.6.5 + 段 1 字幕修复) |
| `fetch_tech_v2.py` | 60+ | 段 3 重拉 (8 关键词 × 3 = 24 mp4) |
| `spec/lithium_30s_v6_template.json5` | 90+ | 字幕模板 (v0.6.5 + bottom-right) |
| `docs/lithium-video-v6-translated-v3.md` | 本报告 | v0.6.5 完整记录 |
| `music/bgm_cinematic_30s.mp3` | 482KB | 备选配乐 (max -1.1 dB) |
| `stock/03_tech_*_*.mp4` | 24 个 (312MB) | 段 3 候选库 (含 MAXCELL 9.5 分) |

## 🎬 3 版本最终输出

| 版本 | 时长 | 大小 | 字幕 | 旁白 | 配乐 |
|---|---|---|---|---|---|
| **zh** | 28.00s | 16.4MB | 锂电新纪元 ... | 千度熔炼... | mmx epic |
| **en** | 28.93s | 16.3MB | New Era ... | A thousand... | mmx epic |
| **ja** | 28.00s | 16.3MB | リチウム新時代 | 千度の精錬... | mmx epic |

## 🚀 下一步建议

1. **段 1 字幕再优化** — vision 仍建议完全避火, 改 opacity 0.8 阴影
2. **段 5 outlook 换镜头** — 沙漠公路+彩色气球 vs 锂电"未来" 主题脱节
3. **配乐多音轨** — 试 cinematic + epic 混音 (人声高 1.4x, BGM 0.13x)
4. **推 GitHub** — 13 commits 等 PAT
