# 结构化模板格式规范 v2.0
> via54Design — 模板系统核心规范
> 将设计模式从「散文描述」转化为「结构化数据」，实现确定性执行 + 版本化演进

## 设计原则

1. **数据驱动 > 散文描述**：每个模板是 YAML 结构，CSS 自动生成，不靠 LLM 重新理解
2. **确定性组合**：布局 × 配色 × 字体 × 叙事 可独立版本化、自由组合
3. **可验证**：每个模板自带质量指标（对比度/字号/行宽）
4. **可进化**：模板版本号递增，旧版保留不删

## 模板层级

```
templates/registry.yaml  (入口索引)
├── layouts/              (布局模板 — 决定结构)
├── color-schemes/        (配色模板 — 决定色彩)
├── typography/           (字体模板 — 决定排印)
├── animations/           (动画模板 — 决定动效)
├── narratology/          (叙事模型 — 决定故事结构)
│   ├── models/           (叙事模型定义)
│   └── guides/           (分镜/叙事指南)
└── video-edits/          (视频剪辑模板 — 参考文档)
```

## YAML Schema

### 布局模板 (layout) — v2

```yaml
id: hero-split-16-9                 # 唯一 ID
name:                               # 多语言名称
  zh: "左右分割式 Hero (16:9)"
  en: "Split Hero (16:9)"
version: "2.0.0"                     # 语义化版本
category: layout/hero                # 分类（用于 registry 分组）
tags: [hero, split, 16-9]           # 筛选标签

# ─── 使用条件（供 LLM 自动选择）───
when:
  content_has: [hero_image, headline]
  suitable_for: [landing, brand_story]

# ─── 视口配置 ───
viewport:
  baseline: "16:9"                  # 设计基准比例
  min_height: 100dvh
  max_width: 1920px                 # TV 端最大宽度
  presentation_mode: false           # 是否默认开启演示锁定

# ─── 结构 ───
structure:
  type: grid-2col                   # grid-2col / grid-3col / flex / bento
  ratio: "5fr 7fr"                  # 列比例
  gap: "0"

# ─── 间距系统（黄金比例 φ=1.618）───
spacing:
  base: 4                           # 基准 4px
  ratio: 1.618
  semantic:
    section: "step-8"              # → --space-section: var(--space-step-8)
    card: "step-5"                 # → --space-card: var(--space-step-5)

# ─── 响应式断点 ───
responsive:
  - name: tv                        # 断点名称
    min_width: 1920                 # 最小宽度
    max_width: 0                    # 0 = 无上限
    columns: "5fr 7fr"              # 覆盖栅格
    font_scale: 1.3                 # 字体缩放
    safe_area: [80, 120, 80, 120]   # TV overscan 安全区 [上右下左]
    spacing_scale: 1.2
    hide_roles: []                  # 按 role 隐藏元素
    stack: false                    # 是否堆叠
    stack_order: []                 # 堆叠顺序 [text, image]

  - name: phone
    min_width: 0
    max_width: 767
    columns: "1fr"
    font_scale: 0.72
    safe_area: [0, 20, 0, 20]
    stack: true
    stack_order: [text, image]
    hide_roles: [eyebrow]

# ─── 元素树 ───
elements:
  - role: image-container
    position: left
    behavior: full-bleed
    z_index: 1
    children:
      - role: image
        tag: div
        behavior: cover

  - role: text-container
    position: right
    padding: [120, 80, 120, 80]
    max_width: 560px
    children:
      - role: eyebrow
        tag: p
        font_size: "clamp(11px, 0.7vw, 14px)"
        letter_spacing: "0.15em"
        text_transform: uppercase
        responsive:                 # ← 元素级响应式（v2 新增）
          phone:
            hide: true
          tablet:
            font_size: "12px"
      - role: headline
        tag: h1
        font_size: "clamp(36px, 5vw, 88px)"
        font_weight: 700

# ─── CSS（手写核心样式，媒体查询由 engine 自动生成）───
css: |                              # 可空，引擎会从 responsive[] 自动编译
  .layout-hero-split { ... }

# ─── 质量 ───
quality:
  min_font_size: 15px
  max_line_length: 50ch
  min_contrast_ratio: 4.5:1

# ─── 案例 ───
examples:
  - brand: Anthropic
    url: https://anthropic.com

changelog:
  "2.0.0": "16:9 基准 + 四端适配 + CSS 自动编译 + 黄金比例间距"
  "1.0.0": "初始版"
```

### 配色模板 (color-scheme)

```yaml
id: ink-wash
name:
  zh: "水墨"
  en: "Ink Wash"

mood: [calm, minimal, zen]
season: all

palette:
  - role: background
    hex: "#F5F0E6"
    name_zh: "澄心堂纸"
    cultural_note: "南唐李后主御用宣纸色"
  - role: text_primary      # 6 语义角色: background / text_primary
    hex: "#1A1A1A"          #         text_secondary / accent
  - role: text_secondary    #         accent_alt / border
    hex: "#6B6B6B"
  - role: accent
    hex: "#C43C3A"
    name_zh: "朱砂印"
```

### 字体模板 (typography)

```yaml
id: ming-hei-editorial
name:
  zh: "明黑配"

fonts:
  display: "'Source Han Serif', serif"    # 标题字体
  body: "'Source Han Sans', sans-serif"   # 正文字体
  mono: "'JetBrains Mono', monospace"     # 等宽字体

sizes:
  display: "clamp(42px, 5vw, 88px)"      # 响应式字号
  body: "clamp(16px, 1.2vw, 21px)"

google_fonts: ["Source+Han+Serif:wght@400;700", "Source+Han+Sans"]
```

### 叙事模型 (narratology/model)

```yaml
id: three-act
name:
  zh: "三幕剧"
  en: "Three-Act"

beats:
  - id: setup
    name:
      zh: "设问"
    mood: mysterious
    duration_weight: 0.2
    voiceover_template: "是否曾经..."
    sub_beats:                    # 子节拍（可选）
      - id: reveal
        weight: 0.5

shot_types: [WIDE, MEDIUM, CLOSE-UP, DETAIL]
camera_moves: [Static, Slow zoom, Dolly in]
```

## 引擎自动生成

以下 CSS **不需要手写**，引擎根据 YAML 自动编译：

| YAML 字段 | 引擎自动生成 |
|-----------|-------------|
| `spacing` | `--space-step-1..12` CSS 变量 (φ=1.618) |
| `responsive[].columns` | `@media { grid-template-columns }` |
| `responsive[].font_scale` | `calc(1em * scale)` |
| `responsive[].safe_area` | `padding` 覆盖 |
| `responsive[].hide_roles` | `display: none` |
| `responsive[].stack` | `grid-template-columns: 1fr` + order |
| `element.responsive` | 独立 `@media` 块 |
| `presentation_mode` | `.presentation-mode` 16:9 锁定容器 |

## 引擎不自动生成（需要手写）

| 内容 | 原因 |
|------|------|
| 布局核心 `display: grid` / `flex` | 布局类型太多，无法穷举 |
| 元素具体样式（border-radius、transition） | 视觉细节因设计而异 |
| 背景图片/渐变 | 无通用模式 |
| hover 交互 | 交互设计不可预测 |
