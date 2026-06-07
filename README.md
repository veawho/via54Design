#### 致敬Aaron Swartz & Tim Berners-Lee
## 人类的灵感与创造力，因互联网与AI而永生
开发思路：让弥散而跳动的人类灵感与AI碰撞出惊艳的故事；将历久弥新的人类经验，用于规范LLM生产力的可控性。
>
> 如切如磋，如琢如磨。——《诗经·卫风·淇奥》

---

## 🧠 让时时陪伴你的AI，成为你最会讲故事的朋友（核心创意工作流）
**via54Design 不是替代你创作，而是赋予你一双抓住灵感的妙手。**

> 文章本天成，妙手偶得之。
> 粹然无疵瑕，岂复须人为？
> 宋·陆游《文章》

人类的灵感是弥散而跳动的，AI 是结构化的、可控的。via54Design 在这两者之间搭建桥梁——把人类的"一句话灵感"转化为 AI 可执行的叙事脚手架，再通过模板引擎输出视频、演示文稿、创意图片。

---

### 第一部分：故事 → 视频能力

#### 从一句话到 90 秒品牌故事

**Step 1 — 人类写一句开头（人类独有的灵感）**

> "1920年代，一个中国裁缝在巴黎开了一家小店，
> 他做的旗袍融入了 Art Deco 的几何线条。
> 没有人想到，这件衣服会改变两个文明的时尚。"

这一句里有人物（裁缝）、时代（1920s）、地点（巴黎）、冲突（东方 vs 西方）、悬念（改变时尚）。AI 无法凭空创造这个种子——它来自你。

**Step 2 — AI 扩展叙事脚手架**

```bash
via54 narrate --seed "1920年代，一个中国裁缝在巴黎开了一家小店..." \
  --model heros-journey --duration 90 --format json --output scaffold.json
```

AI 分析你的种子，匹配最合适的叙事模型，输出结构化脚手架：

```
📋 英雄之旅 (Hero's Journey)  90秒
├── 第一幕·日常  (0-22s)  mood: calm      旁白: 每天，我们都...
├── 第二幕·相遇  (22-44s) mood: curious   旁白: 直到有一天...
├── 第三幕·蜕变  (44-66s) mood: excited   旁白: 不一样了...
└── 第四幕·回归  (66-90s) mood: inspiring  旁白: 每一天...

📋 分镜表: 12 个 shot（WIDE / MEDIUM / CLOSE-UP / DETAIL 循环）
📋 Fountain 剧本骨架（4 幕 8 场景）
🎞️ LLM 完整提示词（可直接喂给 Claude / GPT 生成完整剧本）
```

**Step 3 — 人类选择叙事模型，确认方向**

你可以在 4 种叙事模型中选择，控制故事的节奏和情绪走向：

| 模型 | 节拍 | 适合讲什么故事 |
|------|------|---------------|
| `three-act` | 设问 → 解答 → 号召 | 产品发布、品牌广告 |
| `heros-journey` | 日常 → 相遇 → 蜕变 → 回归 | 品牌故事、纪录片 |
| `cognitive-arc` | 钩子 → 基础 → 核心 → 案例 → 延展 → 总结 | 科普、教程 |
| `problem-solution` | 痛点 → 方案 → 证明 → 行动 | 销售视频、Demo |

每个模型定义在 `templates/narratology/models/*.yaml`，你可以自由扩展。

**Step 4 — AI 生成多场景 HTML 动画**

```bash
via54 generate --from-narrative scaffold.json --output story.html
```

每个叙事节拍自动映射为独立场景，情绪驱动配色切换：

| 场景 | 时长 | 情绪 | 配色 | 旁白 |
|------|------|------|------|------|
| 日常（上海裁缝铺） | 22s | calm | moon-white | "每天，我们都..." |
| 相遇（巴黎开店） | 22s | curious | ink-wash | "直到有一天..." |
| 蜕变（Art Deco 旗袍） | 22s | excited | candy-duolingo | "不一样了..." |
| 回归（影响时尚） | 24s | inspiring | warm-editorial | "每一天..." |

**Step 5 — 导出视频 + 配乐**

```bash
via54 export render story.html --duration 90 --format mp4   # HTML → MP4
via54 media add-music story.mp4 --mood=ad                    # 背景音乐
```

支持多格式：`mp4` / `webm` (VP9) / `hevc` (H.265) / `frames` (PNG序列) / `apng`。

**完整管线（一行命令）**

```bash
via54 narrate --seed "你的故事开头" --model heros-journey --format json \
  | via54 generate --from-narrative /dev/stdin --output story.html --presentation
```

---

### 第二部分：故事 → 演示能力

同一个叙事脚手架，可以导出多种演示格式，无需重新创作。

**Step 1 — 相同的叙事种子**

```bash
via54 narrate --seed "1920年代，一个中国裁缝在巴黎..." \
  --model heros-journey --duration 90 --format json --output scaffold.json
```

**Step 2 — 导出为 PPTX 演示文稿**

```bash
via54 export pptx scaffold.json --output story.pptx
```

- 纯 Go 实现，零外部依赖（不依赖 Node.js / unioffice）
- 每幕一张幻灯片，情绪映射强调色（左侧装饰条）
- 16:9 宽屏，标题 + 旁白 + 页码
- 文字直接在 PPT 中可编辑

**Step 3 — 导出为 Markdown 幻灯片（Marp 兼容）**

```bash
via54 export markdown scaffold.json --output slides.md
npx @marp-team/marp-cli slides.md --pptx    # 转 PPTX
npx @marp-team/marp-cli slides.md --pdf     # 转 PDF
npx @marp-team/marp-cli slides.md --html    # 转 HTML
```

- 兼容 Marp (⭐11,917) 幻灯片生态
- 情绪标注为 CSS class，可用自定义主题
- 支持 YAML frontmatter（title / author / theme）

**Step 4 — 导出为结构化 JSON**

```bash
via54 export json scaffold.json --output scenes.json
```

- 包含 timing（start_sec / end_sec）
- 可供外部工具或自定义管线消费

**Step 5 — 导出为 PDF**

```bash
via54 export pdf story.html --output story.pdf
```

---

### 第三部分：故事 → 创意图片能力

同一个叙事脚手架，还可以导出为矢量图片，用于海报、社交媒体、印刷等场景。

**Step 1 — 叙事脚手架 → SVG 矢量场景**

```bash
via54 narrate --seed "1920年代，一个中国裁缝在巴黎..." \
  --model heros-journey --duration 90 --format json --output scaffold.json

via54 export svg scaffold.json --output ./scenes
```

输出 `./scenes/` 目录，每幕生成一个独立 SVG 文件：

```
scenes/
├── scene-001-日常-(ordinary_world).svg   # 22s
├── scene-002-相遇-(call_to_adventure).svg  # 22s
├── scene-003-蜕变-(transformation).svg     # 22s
└── scene-004-回归-(return).svg             # 24s
```

**SVG 特性**：

- 16:9 viewBox，无限缩放不失真
- 情绪映射配色：
  - `calm` → 绿色调 `#f0f4e8`
  - `curious` → 暖白调 `#f5f0e6`
  - `excited` → 深绿调 `#1a3a1a`
  - `inspiring` → 暖黄调 `#fdf6e3`
- 12px 强调色装饰条（左侧）
- 旁白以斜体呈现（底部半透明背景）
- 页码标注（1/4, 2/4, 3/4, 4/4）

**Step 2 — 照片 → SVG 矢量化**

```bash
via54 media trace --input logo-sketch.jpg --output logo.svg     # 手绘Logo → 矢量
via54 media trace --input handwriting.jpg --output title.svg     # 书法 → 矢量
via54 media trace --input photo.jpg --output portrait.svg        # 照片 → 矢量
```

基于 VTracer (⭐6,150) 引擎，保留原始笔触质感。

**Step 3 — 场景 → AI 生图提示词（新增）**

人类写一句基础场景，AI 生成结构化提示词，人类修改确认后喂给 Midjourney / 可灵 / 即梦 / Gemini。

```bash
via54 prompt --scene "1920年代，巴黎左岸的小裁缝店里，一位中国裁缝在制作旗袍" \
  --platform midjourney --output prompt.md
```

输出结构化提示词，包含可编辑字段：

```
📋 平台: midjourney | 格式: midjourney
⏳ [subject]      （LLM填充：主体描述）
⏳ [environment]  （LLM填充：环境描述）
✅ [lighting]     自然光
✅ [style]        摄影写实
✅ [mood]         平静
✅ [composition]  中景

最终 prompt:
（LLM填充：主体描述）, （LLM填充：环境描述）, 自然光, 摄影写实, 平静, 中景 --ar 16:9 --v 6.1 --style raw --s 250
```

人类可以直接修改 YAML 字段值，然后重新生成最终 prompt。支持平台：

| 平台 | 格式 | 特点 |
|------|------|------|
| `midjourney` | `subject, env, style --params` | --ar / --v / --style / --s |
| `kling` | 结构化 + 运镜参数 | duration / motion / cfg_scale |
| `jimeng` | 中文结构化 | 国风 / 写实 / 3D |
| `gemini` | 自然语言段落 | 适合 Gemini Imagen |

参考项目: ai-media-generator (⭐70), Ultimate-AI-Media-Generator-Skill (⭐57)

**矢量标题嵌入设计**

```bash
via54 generate --lettering-svg ./title.svg --color ink-wash \
  --font calligraphy-accent --title "山水之间" --output poster.html
```

书法/手写标题直接嵌入 HTML 设计，不依赖任何字体文件。

---

### 三种能力对比

| 能力 | 输入 | 输出 | 引擎 | 依赖 |
|------|------|------|------|------|
| 🎬 **故事→视频** | 一句话 → 叙事JSON | MP4/WebM/HEVC + 配乐 | narrate + generate + Playwright | 需 ffmpeg |
| 📊 **故事→演示** | 一句话 → 叙事JSON | PPTX / Markdown / JSON / PDF | **纯 Go** | **零** |
| 🎨 **故事→创意图片** | 一句话 → 叙事JSON | SVG 矢量文件 | **纯 Go** + VTracer | 需 VTracer (trace) |
| 🖼️ **场景→生图提示词** | 一句话场景 | 结构化 Prompt (MJ/Kling/即梦/Gemini) | **纯 Go** (YAML模板) | **零** |

---

## 🗣️ 自然语言操作示例（给 AI 助手）

把下面的任意一句话发给 AI（Claude / Cursor / Copilot），它能自动完成：

| 你想做什么 | 对 AI 说 |
|-----------|----------|
| **构思故事** | "用 via54Design 帮我构思一个品牌故事，种子是：'一个中国裁缝在巴黎改变了时尚'" |
| **全链路** | "从这句话开始，帮我做一个 60 秒的品牌故事动画：'一个农民用无人机种出了最好吃的大米'" |
| **做 HTML** | "用 via54Design 生成一个暖色编辑风格的页面，标题叫'羽图鉴'" |
| **检查质量** | "帮我检查这个 HTML 文件的质量: via54 quality --html demo.html" |
| **提取风格** | "从我的 HTML 里提取配色和字体方案: via54 pattern --html demo.html" |
| **加背景音乐** | "给这个视频配一首科技感 BGM: via54 media add-music demo.mp4 --mood=tech" |
| **转 GIF** | "把视频转成 GIF: via54 media convert demo.mp4" |
| **取素材图** | "帮我找一些鹦鹉的公共领域插画: via54 media fetch --query parrot" |
| **录视频** | "把这个 HTML 录成 30 秒 1080p 视频: via54 export render demo.html --duration 30" |
| **转 PDF** | "把这个页面导出为 PDF: via54 export pdf demo.html" |
| **语音合成** | "把这段文字转成语音: via54 export tts --text '你好世界' --out hello.mp3" |
| **矢量化** | "把我手写的 logo 草稿转成 SVG: via54 media trace --input sketch.jpg" |
| **列出模板** | "看看有哪些设计模板可以用: via54 list" |

---

## ⚡ 快速上手

```bash
# 查看所有命令
via54

via54 list                         # 列出所有模板（配色40+ / 字体12 / 布局3 / 叙事模型4）

# 叙事驱动：一句话 → 多场景 HTML 动画
via54 narrate --seed "你的故事种子" --model three-act --format json | via54 generate --from-narrative /dev/stdin

# 标准设计模板生成
via54 generate --layout hero-split --color ink-wash --font ming-hei-editorial --title "我的设计" --output demo.html

# 查看生成的HTML
open demo.html     # Mac
start demo.html   # Windows
xdg-open demo.html # Linux
```

---

## 💻 命令参考

```bash
via54                              # 帮助
via54 version                      # 版本信息
via54 list                         # 列出所有模板

# 叙事引擎 — 人类+AI协作讲故事
via54 narrate --list                                          # 查看4种叙事模型
via54 narrate --seed "一句话" --model three-act                # 输出叙事脚手架 (markdown)
via54 narrate --seed "..." --model heros-journey --duration 60 # 指定时长和模型
via54 narrate --seed "..." --format json --output scaffold.json # 输出JSON (供generate消费)
via54 narrate --seed "..." --model cognitive-arc              # 四选一叙事模型

# 叙事驱动生成 (全管线)
via54 narrate --seed "你的故事" --format json | via54 generate --from-narrative /dev/stdin

# 设计模板
via54 generate --layout <id> --color <id> --font <id> --title "标题" --output out.html
via54 generate --lettering-svg ./vector.svg --title "标题" --output out.html
via54 generate --layout <id> --color <id> --font <id> --lettering-svg ./art.svg --title "标题" --output out.html
via54 quality --html out.html
via54 pattern --html out.html --name "项目名"

# 媒体管线
via54 media add-music input.mp4 --mood tech|ad|educational
via54 media convert input.mp4
via54 media fetch --query "关键词" --out ./img --count 3
via54 media trace --input photo.jpg --output logo.svg       # 照片→SVG矢量化
via54 media trace --input handwriting.jpg --output title.svg # 书法/签名→SVG

# 导出 (纯 Go，零外部依赖)
via54 export render input.html --duration 30 --width 1920 --height 1080 --format mp4   # MP4 视频
via54 export render input.html --format webm                                            # WebM (VP9)
via54 export render input.html --format frames                                          # PNG 序列帧
via54 export pdf input.html                                                              # PDF
via54 export pptx --output deck.pptx scaffold.json                                       # PPTX 演示文稿 (从叙事)
via54 export svg --output ./scenes scaffold.json                                         # SVG 矢量稿 (每场景独立文件)
via54 export json --output scenes.json scaffold.json                                     # 结构化场景数据
via54 export markdown --output slides.md scaffold.json                                   # Marp 兼容幻灯片
via54 export tts --text "你好" --out voice.mp3                                           # 语音合成

# 图片提示词 (Prompt)
via54 prompt --list                                                      # 查看可用平台
via54 prompt --scene "场景描述" --platform midjourney                     # 生成 Midjourney 提示词
via54 prompt --scene "场景" --platform kling --output prompt.md           # 输出到文件

# MCP Server (独立二进制: via54-mcp, 兼容: via54 serve)
via54 serve
# 推荐使用独立二进制: via54-mcp
```

### 命令标志速查

| 命令 | 标志 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| **generate** | `--layout` | string | — | 布局模板 ID (`via54 list` 查看) |
| | `--color` | string | — | 配色模板 ID |
| | `--font` | string | — | 字体模板 ID |
| | `--title` | string | "via54Design" | 页面标题 |
| | `--output` | string | "output.html" | 输出文件路径 |
| | `--presentation` | bool | false | 16:9 演示锁定 (PPT/视频) |
| | `--lettering-svg` | string | — | 手写/书法 SVG 路径 |
| | `--from-narrative` | string | — | 叙事脚手架 JSON (via54 narrate --format json) |
| **narrate** | `--seed` | string | **必填** | 一句话故事种子 |
| | `--model` | string | "three-act" | 叙事模型 ID |
| | `--duration` | int | 30 | 目标视频时长(秒) |
| | `--format` | string | "markdown" | 输出格式 (markdown/json) |
| | `--output` | string | stdout | 输出文件路径 |
| | `--list` | bool | false | 列出所有叙事模型 |
| **media** | `add-music` | `--mood` | string | "tech" | 配乐情绪 (tech/ad/educational) |
| | `fetch` | `--query` | string | **必填** | 搜索关键词 |
| | | `--out` | string | "./img" | 输出目录 |
| | | `--count` | int | 2 | 每关键词张数 |
| | `trace` | `--input` | string | **必填** | 输入图片路径 |
| | | `--output` | string | — | 输出 SVG 路径 |
| **export** | `render` | `--duration` | int | 10 | 时长(秒) |
| | | `--width` | int | 1920 | 宽 |
| | | `--height` | int | 1080 | 高 |
| | | `--format` | string | "mp4" | 视频格式: mp4/webm/hevc/frames/apng |
| | `pdf` | `--output` | string | — | 输出路径 |
| | `pptx` | `--output` | string | "output.pptx" | 输出路径 |
| | | `--16-9` | bool | true | 16:9 宽屏 |
| | | `--title` | string | "via54 演示文稿" | 标题 |
| | `svg` | `--output` | string | "./svg-scenes" | 输出目录 |
| | | `--width` | int | 1920 | 宽 |
| | | `--height` | int | 1080 | 高 |
| | `json` | `--output` | string | "scenes.json" | 输出路径 |
| | `markdown` | `--output` | string | "story.md" | 输出路径 |
| | | `--title` | string | "via54 演示文稿" | 标题 |
| | | `--author` | string | "via54Design" | 作者 |
| | `tts` | `--text` | string | **必填** | 文本 |
| | | `--output` | string | "output.mp3" | 输出路径 |
| | | `--voice` | string | — | 音色 |
| **quality** | `--html` | string | **必填** | HTML 文件路径 |
| | `--verbose` | bool | false | 显示 info 级问题 |
| **pattern** | `--html` | string | **必填** | HTML 文件路径 |
| | `--name` | string | "unnamed" | 作品名称 |

---

## 🎨 配色模版一览

40+ 套配色方案，每套包含 6 个语义角色（背景/正文/辅助/强调/强调2/边框）。

### 中国传统配色 (8)

| 方案 | 季节 | 气质 | 强调色 | 典故 |
|------|------|------|--------|------|
| `crimson-elegance` | 秋 | 热烈·庄重 | `#C23A2B` 朱砂 | 唐代宫廷正色 |
| `pine-spring` | 春 | 自然·雅致 | `#5B8C5A` 松花绿 | 苏轼松花酿酒 |
| `daylily-warmth` | 夏 | 温暖·忘忧 | `#E8A838` 萱草黄 | 诗经忘忧草 |
| `ultramarine-deep` | 冬 | 深邃·宁静 | `#2E5CB8` 群青 | 敦煌青金石 |
| `rosewood-noble` | 秋 | 高贵·典雅 | `#7A4B5C` 紫檀 | 明清一寸紫檀一寸金 |
| `moon-white` | 冬 | 素雅·空灵 | `#5B7B8C` 雨过天青 | 宋徽宗汝窑梦色 |
| `autumn-fragrance` | 秋 | 古朴·温润 | `#9B8B6E` 秋香 | 明人书房桂花色 |
| `ink-wash` | 四季 | 极简·禅意 | `#C43C3A` 朱砂印 | 王维始创水墨画 |

```bash
via54 generate --layout hero-split --color ink-wash --font ming-hei-editorial --title "寒山寺" --output demo.html
```

### 日系配色 (6)

| 方案 | 来源 | 气质 | 强调色 | 文化背景 |
|------|------|------|--------|----------|
| `tsubaki-camellia` | 资深堂 | 优雅·知性 | `#BF3A2B` 椿色 | 山茶花口红传奇 |
| `wabi-sabi` | 千利休茶道 | 残缺·空寂 | `#6B5B4A` 焦茶 | 侘寂美学 |
| `muji-minimal` | 原研哉 | 极简·功能 | `#B27C5A` 亚麻 | 空无的设计哲学 |
| `sakura-blossom` | 花见 | 温柔·短暂 | `#E8A0B4` 薄紅 | 一期一会 |
| `indigo-craft` | 阿波藍 | 匠人·深沉 | `#264C7B` 藍色 | Japan Blue |
| `rinpa-gold` | 尾形光琳 | 华美·装饰 | `#C89B3C` 金 | 风神雷神屏风 |

### 经典配色 (6)

| 方案 | 灵感来源 | 气质 | 强调色 |
|------|----------|------|--------|
| `warm-editorial-cream` | Anthropic | 温暖·智性 | `#CC785C` 赤陶橙 |
| `dark-terminal-blue` | Linear | 科技·发光 | `#5E6AD2` 紫蓝 |
| `swiss-monochrome` | Vercel | 极简·权威 | `#000000` 纯黑 |
| `bauhaus-primary` | Khan Academy | 几何·活力 | `#E63946` 红 |
| `candy-duolingo` | Duolingo | 游戏·亲和 | `#58CC02` 绿 |
| `cosmic-retro` | Perplexity | 太空·复古 | `#2B4F91` 钴蓝 |

### Adobe / Behance 流行配色 (6)

| 方案 | 设计运动 | 气质 | 强调色 |
|------|----------|------|--------|
| `spectrum-indigo` | Adobe Spectrum | 专业·创意 | `#5C5CE0` Indigo |
| `flat-ui-vibrant` | Flat Design | 扁平·活力 | `#3498DB` 碧蓝 |
| `millennial-sage` | Instagram美学 | 温柔·趋势 | `#E8A0BF` 千禧粉 |
| `glassmorphism-pastel` | 玻璃拟态 | 未来·梦幻 | `#7C6BF5` 紫罗兰 |
| `neon-dark` | 赛博朋克 | 暗黑·霓虹 | `#FF2D95` 荧光粉 |
| `earth-terracotta` | 返璞归真 | 大地·温暖 | `#C06C4C` 陶土橙 |

---

## 🔤 字体模版一览

12 套字体配对方案，每套包含显示(display)、正文(body)、等宽(mono)三层字体栈。

### 中文排版 (6)

| 方案 | 分类 | 显示字体 | 正文字体 | 灵感来源 |
|------|------|---------|---------|----------|
| `ming-hei-editorial` | 明体+黑体 | Source Han Serif | Source Han Sans | JustFont 经典推荐 |
| `kai-rounded-friendly` | 楷体+圆体 | LXGW WenKai | Nunito | 霞鹜文楷开源项目 |
| `song-literary` | 仿宋+明体 | Source Han Serif | FangSong | 宋刻本文学传统 |
| `hei-modern` | 黑体统一 | Source Han Sans | Source Han Sans | Adobe+Google 联合开发 |
| `calligraphy-accent` | 书法+黑体 | ZCOOL XiaoWei | Inter | 站酷小薇开源书法 |
| `sc-sans-clean` | 中文无衬线 | Noto Sans SC | Noto Sans SC | Google Noto 项目 |

### 国际排版 (6)

| 方案 | 分类 | 显示字体 | 正文字体 | 灵感来源 |
|------|------|---------|---------|----------|
| `serif-sans-editorial` | 过渡衬线+人文无衬线 | Fraunces | Inter | Anthropic/Claude |
| `sans-geometric-tech` | 几何无衬线+等宽 | Inter | Inter | Linear/Vercel |
| `display-sans-bold` | Grotesque 展示 | Archivo Black | Inter | Apple Keynote |
| `elegant-didone` | 迪多体高反差 | Playfair Display | Inter 300 | Vogue / Harper's |
| `mono-code` | 等宽主导 | JetBrains Mono | Inter | JetBrains/Cursor |
| `playful-rounded` | 圆体亲和 | Baloo 2 | Nunito | Duolingo/Khan Academy |

---

## 🧱 布局模板一览

3 套布局模板，全部以 **16:9 视频基准** 设计，支持 **TV/Desktop/Tablet/Phone 四端适配**。CSS 自动编译 + 黄金比例间距 + 元素级响应式。

| 模板 | 中文名 | 结构 | TV ≥1920 | Desktop 1280-1919 | Tablet 768-1279 | Phone <768 |
|------|--------|------|----------|-------------------|-----------------|------------|
| `hero-split-16-9` | 左右分割 Hero | 左图右文 5:7 | 5:7 分屏 ×1.3字 safe-area 120px | 5:7 标准 | 堆叠 text↑image↓ | 堆叠 ×0.72 无眉标 |
| `bento-grid-2x2` | Bento 便当格 | 2×2 卡片 | **3×2 六格** ×1.3字 | 2×2 标准 | 2 列 | **1列瀑布** ×0.75 |
| `gallery-waterfall` | 画廊瀑布流 | 自动填充网格 | **5列** 16:9锁定 | **4列** | **3列** | **2列** 常显标题 |

```bash
# 普通网页 — 自由布局
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial --title "品牌故事"

# 演示模式 — 16:9 锁定 (PPT/视频输出)
via54 generate --layout bento-grid-2x2 --color dark-terminal-blue --font mono-code --presentation

# Bento 数据卡片
via54 generate --layout bento-grid-2x2 --color dark-terminal-blue --font mono-code --title "Dashboard"

# 画廊
via54 generate --layout gallery-waterfall --color crimson-elegance --font calligraphy-accent --title "Gallery"
```

### 间距系统

每个布局内置 **黄金比例间距** (`base=4, φ=1.618`)，以 CSS 变量注入：

```css
--space-step-1: 4px;   --space-step-2: 6px;    --space-step-3: 10px;
--space-step-4: 17px;  --space-step-5: 27px;   --space-step-6: 44px;
--space-step-7: 72px;  --space-step-8: 116px;  --space-step-9: 188px;
--space-card: var(--space-step-5);   /* 27px */
--space-section: var(--space-step-8); /* 116px */
```

### 响应式自动编译

YAML `responsive[]` 自动生成 CSS `@media` 查询，覆盖 **columns/stack/safe_area/font_scale/hide_roles**，无需手写媒体查询。

### 元素级响应式

每个 Element 支持按断点配置 hide/order/fontSize/padding：

```yaml
- role: eyebrow
  responsive:
    phone:
      hide: true
    tablet:
      font_size: "12px"
```

---

## 🎬 叙事模型一览

4 种叙事模型，以 YAML 模板定义，可自由扩展。放在 `templates/narratology/models/` 下：

| 模型 | 中文名 | 节拍 | 适用场景 |
|------|--------|------|----------|
| `three-act` | 三幕剧 | 设问 → 解答 → 号召 | 产品发布、品牌广告 |
| `heros-journey` | 英雄之旅 | 日常 → 相遇 → 蜕变 → 回归 | 品牌故事、纪录片 |
| `cognitive-arc` | 认知弧 | 钩子 → 基础 → 核心 → 案例 → 延展 → 总结 | 科普、教程 |
| `problem-solution` | 问题-解法 | 痛点 → 方案 → 证明 → 行动 | 销售视频、Demo |

添加新模型：只需在 `templates/narratology/models/` 下创建 YAML 文件并在 `registry.yaml` 注册，无需修改 Go 代码。

---

## 📦 安装

```bash
# 编译两个二进制
go build -o via54 ./cmd/via54/
go build -o via54-mcp ./cmd/mcp-server/

# 或使用 Makefile
make all
```

### 依赖

| 工具 | 用途 | 安装 |
|------|------|------|
| Go 1.21+ | 核心引擎 | [go.dev/dl](https://go.dev/dl/) |
| ffmpeg | 视频/音频处理 | `brew install ffmpeg` / `apt install ffmpeg` |
| Playwright (可选) | HTML→视频/PDF | `npm install playwright && npx playwright install chromium` |

---

# MCP 配置

via54Design 提供独立 MCP Server 二进制 `via54-mcp`，与主 CLI 分开部署。

**Claude Desktop** (`claude_desktop_config.json`):
```json
{ "mcpServers": { "via54Design": { "command": "via54-mcp" } } }
```

**Cursor** (`.cursor/mcp.json`):
```json
{ "mcpServers": { "via54Design": { "command": "via54-mcp" } } }
```

**也可以使用主 CLI 的 serve 命令**（兼容旧配置）:
```json
{ "mcpServers": { "via54Design": { "command": "via54", "args": ["serve"] } } }
```

**安装**:
```bash
# 编译两个二进制
make all
# 或分别编译
go build -o via54-mcp ./cmd/mcp-server/
go build -o via54 ./cmd/via54/
```

---

## 📐 架构

```
你的创意灵感                ← 不可替代的人类优势
    │   通过 narrate --seed 注入
    ▼
┌─ 叙事引擎 ──────────────────────────────┐
│  YAML 4 种叙事模型 (templates/narratology/)│ ← 可扩展
│  narrate → 脚手架 + 剧本 + 分镜        │
│  generate --from-narrative → 多场景      │
└─────────────────────────────────────────┘
    │
    ▼
┌─ 设计模板引擎 (internal/template/) ─────┐
│  40+ 配色 / 12 字体 / 3 布局 / 响应式    │
│  13MB 单二进制 (Go), WASM 加速可选 (Rust)│
│  CSS 自动编译 + 黄金比例间距             │
└─────────────────────────────────────────┘
    │
    ▼
┌─ 导出管线 (internal/export/) ──────────┐
│  HTML → PPTX / SVG / JSON / Markdown    │ ← 纯 Go
│  HTML → MP4/WebM/HEVC (Playwright)      │ ← 可选依赖
│  HTML → PDF / TTS                       │
└─────────────────────────────────────────┘
    │
    ▼
┌─ 媒体管线 (internal/media/) ───────────┐
│  配乐 (ffmpeg) / 格式转换               │
│  照片→SVG 矢量化 (VTracer)              │
│  图片搜索 (Wikimedia/Unsplash)          │
└─────────────────────────────────────────┘

cmd/via54/   ← CLI (10 文件, 主入口 + 9 子命令)
cmd/mcp-server/  ← MCP Server 独立二进制
hack/        ← 构建/部署脚本
docs/        ← 模板格式规范 + 故障恢复指南
```

## 语言

| 语言 | 用途 | 文件数 |
|------|------|--------|
| Go | 核心引擎、CLI、MCP、全部导出 | 22 |
| Rust | WASM 加速引擎（可选编译） | 5 |
| Bash | 部署/构建脚本 | 3 |

## 许可

双许可：**MIT OR AGPL-3.0** (SPDX: `MIT OR AGPL-3.0`)

| 目录 | 许可 | 说明 |
|------|------|------|
| `templates/` | MIT | YAML 模板定义 |
| `hack/` | MIT | 部署/编译脚本 + Rust WASM |
| `docs/` | MIT | 文档 |
| `internal/` + `cmd/` | AGPL-3.0 | Go 源代码 |

## 致谢

via54Design 建立在以下开源项目的基础上，感谢所有贡献者：

| 项目 | 许可 | 用途 |
|------|------|------|
| **[huashu-design](https://github.com/alchaincyf/huashu-design)** by alchaincyf（花叔·花生） | MIT | 基础设计模板引擎，本项目由此衍生 |
| **[VTracer](https://github.com/visioncortex/vtracer)** | MIT | 照片→SVG 矢量化引擎 |
| **[wazero](https://github.com/tetratelabs/wazero)** | Apache-2.0 | WebAssembly 运行时（Rust WASM 加速） |
| **[mcp-go](https://github.com/mark3labs/mcp-go)** | MIT | MCP Server 框架 |
| **[Extra-Strength Responsive Grids](https://github.com/johnpolacek/extra-strength-responsive-grids)** | — | 流体 CSS Grid 系统参考 |
| **[Marp](https://github.com/marp-team/marp)** | MIT | Markdown 幻灯片生态（export markdown 兼容格式） |
| **[golang-standards/project-layout](https://github.com/golang-standards/project-layout)** | — | Go 项目目录结构参考 |

### 设计参考

| 项目/资源 | 说明 |
|-----------|------|
| **Fountain screenplay format** ([fountain.io](https://fountain.io)) | 剧本格式标准（narrate 的 Fountain 输出） |
| **Aristotle's Poetics / Syd Field's Screenplay** | 三幕剧叙事理论（three-act 模型） |
| **Joseph Campbell's The Hero with a Thousand Faces** | 英雄之旅叙事理论（heros-journey 模型） |
| **Cognitive Load Theory (Sweller, 1988)** | 认知弧叙事模型（cognitive-arc 模型） |
| **Direct Response Marketing (Caples/Schwartz)** | 问题-解法叙事模型（problem-solution 模型） |
| **Pinterest Waterfall / Masonry Layout** | 瀑布流布局参考（gallery-waterfall） |
| **Apple Bento Grid Design Language** | Bento 便当格布局参考（bento-grid-2x2） |

### 评价参考

| 项目 | 说明 |
|------|------|
| **huobao-drama** (⭐12,623) | 一句话生成短剧，“叙事脚手架”模式验证 |
| **presenterm** (⭐8,494) | 终端幻灯片工具，演示模式概念验证 |
| **unioffice** (⭐4,875) | 纯 Go Office 文档库，PPTX 实现参考
