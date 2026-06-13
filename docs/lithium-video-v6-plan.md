# 锂电 30s 视频 v6 真实视频素材方案 (★ 终极方案 ★)

> 时间: 2026-06-12
> 任务: **真实视频素材组合** (不用 AI 生成图 + zoompan 假冒视频)
> 工具栈: **Pexels API (免费 CC0) + editly (声明式拼接) + ffmpeg + mmx TTS/配乐**

## 🎯 调研找到 4 大核心金矿

| 项目 | ⭐ | 价值 |
|---|---|---|
| **`poseljacob/agentic-video-editor`** | 437 | **AI Director 选镜 + EditPlan + 5 维评分 + 自动重试** (Gemini + FFmpeg) |
| **`mifi/editly`** | 5435 | **声明式 JSON5 spec 多素材拼接** (Node.js + FFmpeg) |
| **`neutraltone/awesome-stock-resources`** | 14261 | **14K⭐ 免费 stock 视频列表** (Pexels/Pixabay/Mixkit/Coverr) |
| **`mifi/lossless-cut`** | 41218 | **无损视频瑞士军刀** (辅助用) |

## 📐 终极方案架构 (5 步流水线)

```
[1] 关键词清单 (锂电相关)
   ↓
[2] Pexels API 拉真实视频素材 (锂电+工厂+夜景+全息+数据) 
   ↓
[3] 选最优素材 (按 [锂电主题/画质/时长] 排序)
   ↓
[4] editly 声明式拼接 (5 段 30s, dtc-testimonial 结构)
   + mmx TTS 中文旁白 + mmx music 配乐
   ↓
[5] 输出 lithium_30s_v6.mp4 (★ 真实视频素材 ★)
```

## 🎬 5 段镜头 (★ 复用 dtc-testimonial.yaml 黄金结构 ★)

| 段 | 名称 | 时长 | 任务 | 关键词 |
|---|---|---|---|---|
| 1 | **hook 震撼开场** | 3s | **震撼画面** | cyberpunk city night, factory, electric vehicle, technology |
| 2 | **trend 行业趋势** | 5s | 数据/图表 | data visualization, stock chart, financial dashboard |
| 3 | **technology 技术细节** | 10s | 工厂/电池微观 | battery factory, assembly line, robotic arm, battery cell |
| 4 | **market 投资市场** | 7s | 投资/资本 | trading floor, financial market, battery stock, investment |
| 5 | **outlook 投资召唤** | 5s | 收尾/EV | electric vehicle, sunrise, futuristic city, energy |

**总和: 3+5+10+7+5 = 30s 精确**

## 📜 Pexels API 集成方案

### 注册 + 拿 key
1. 访问 https://www.pexels.com/api/ 注册
2. 免费拿 API key (200 req/hour, 20000 req/month)
3. 写入 .env: `PEXELS_API_KEY=xxx`

### 9 关键词视频拉取 (5 段 × 平均 1.8 备用)
```bash
# 段 1 hook: 震撼开场
curl -H "Authorization: $PEXELS_API_KEY" \
  "https://api.pexels.com/videos/search?query=cyberpunk+night+city+neon&per_page=5&orientation=landscape"

# 段 2 trend: 数据/图表
curl -H "Authorization: $PEXELS_API_KEY" \
  "https://api.pexels.com/videos/search?query=data+visualization+stock+chart&per_page=5"

# 段 3 technology: 工厂/电池
curl -H "Authorization: $PEXELS_API_KEY" \
  "https://api.pexels.com/videos/search?query=battery+factory+assembly+robotic&per_page=5"

# 段 4 market: 投资/交易
curl -H "Authorization: $PEXELS_API_KEY" \
  "https://api.pexels.com/videos/search?query=trading+floor+financial+market&per_page=5"

# 段 5 outlook: EV + 黎明
curl -H "Authorization: $PEXELS_API_KEY" \
  "https://api.pexels.com/videos/search?query=electric+vehicle+sunrise+city&per_page=5"
```

## 🎛️ editly 声明式 spec 模板

```json5
{
  outPath: 'lithium_30s_v6.mp4',
  width: 1920,
  height: 1080,
  fps: 30,
  defaults: {
    duration: 6,
    transition: { duration: 0.5, name: 'fade' },
    layer: { fontPath: '...' }
  },
  clips: [
    { duration: 3, layers: [
      { type: 'video', path: 'hook1.mp4' },
      { type: 'title', text: '锂电新纪元', position: 'top' }
    ]},
    { duration: 5, layers: [
      { type: 'video', path: 'trend1.mp4' },
      { type: 'subtitle', text: '2026 产能突破 800GWh' }
    ]},
    { duration: 10, layers: [
      { type: 'video', path: 'tech1.mp4' }
    ]},
    { duration: 7, layers: [
      { type: 'video', path: 'market1.mp4' }
    ]},
    { duration: 5, layers: [
      { type: 'video', path: 'outlook1.mp4' },
      { type: 'title', text: '投资窗口正在打开' }
    ]}
  ],
  audioFilePath: 'narration.mp3',
  loopAudio: false,
  outputVolume: 1.0,
  clipsAudioVolume: 0.3  // 背景 -10dB
}
```

## 🎯 与 v3-v5 对比

| 版本 | 方案 | 视频素材 | 平均分 | 末帧 |
|---|---|---|---|---|
| v3 | image-01 + zoompan | ❌ 假视频 | 7.4 | 7.25 |
| v4 | image-01 + zoompan (5层公式) | ❌ 假视频 | 7.36 | 8.75 |
| v5 | image-01 + zoompan (修复) | ❌ 假视频 | TBD | 8.98 |
| **v6** | **Pexels 真实视频 + editly** | ✅ **真视频** | 预期 **8.0+** | 预期 **8.5+** |

## ⚠️ 风险 + 缓解

| 风险 | 缓解 |
|---|---|
| Pexels 没合适锂电素材 | Mixkit + Pixabay 备选 + 多关键词 |
| Pexels 视频风格不统一 (真实拍摄 vs 3D) | editly 调色统一 + 字幕过渡 |
| 真实视频有时长不匹配 (Pexels 5-30s) | ffmpeg cutTo/cutFrom 截取精确时长 |
| 网络限速 (Pexels 海外 CDN) | 串行下载 + retry |
| mmx 旁白/配乐同步 | editly 音频轨道 + apad |

## 🔧 实施步骤

1. **注册 Pexels API key** (用户操作, 1 分钟)
2. **安装 editly**: `npm i -g editly` (用户有 Node.js)
3. **编写 fetch_pexels.sh** 拉 25 个视频素材
4. **编写 editly spec 模板** lithium_spec.json5
5. **跑 editly 拼接** 5 段 30s
6. **混音**: TTS 旁白 + 配乐 + 视频音频 (amix)
7. **vision 评分 + 提交 v6 report**

## 🎁 沉淀到 SKILL v0.6.2 (新章节)

- **真实视频素材优先** > AI 视频生成 (当 AI 视频 quality 不达标时)
- **Pexels API** = 黄金免费源 (CC0, 高质量, 200 req/h)
- **editly** = 声明式多素材拼接 + transitions + text overlays
- **dtc-testimonial.yaml 结构** = 30s 短视频黄金 5 段 (3+5+10+7+5)
- **真实视频 vs AI 视频 trade-off**: 真实视频时序真实但风格不一致, AI 视频风格统一但时序不自然
