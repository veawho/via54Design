# 结构化模板格式规范 v1.0
> via54Design — 模板系统核心规范
> 目标：将设计模式从「散文描述」转化为「结构化数据」，实现确定性执行 + 版本化演进

## 设计原则

1. **数据驱动 > 散文描述**：每个模板是 YAML 结构，CSS 直接生成，不靠 LLM 重新理解
2. **确定性组合**：布局 × 配色 × 字体 × 动画 可独立版本化、自由组合
3. **可验证**：每个模板自带质量指标（对比度/字号/行宽）
4. **可进化**：模板版本号递增，旧版保留不删

## 模板层级

```
template-registry.yaml  (入口索引)
├── layouts/             (布局模板 — 决定结构)
├── color-schemes/       (配色模板 — 决定色彩)
├── typography/          (字体模板 — 决定排印)
├── animations/          (动画模板 — 决定动效)
└── video-edits/         (视频剪辑模板 — 决定叙事节奏/转场/音效)
```

## YAML Schema

### 布局模板 (layout)

```yaml
# ─── 元数据 ───
id: hero-split-left-image
name: 左右分割式Hero（左图右文）
version: "1.0.0"
created: 2026-06-07
author: via54
category: layout/hero
tags: [hero, split, editorial, image-left, full-bleed]

# ─── 使用条件 ───
when:
  content_has: [hero_image, headline, subtitle, cta]
  suitable_for: [landing, product_page, brand_story]
  min_sections: 2
  max_sections: 4

# ─── 视觉DNA ───
structure:
  type: grid-2col
  ratio: [5, 7]           # 左5/12 : 右7/12
  gap: 0                  # 全出血无间距
  min_height: 100vh       # 全屏

elements:
  - role: image-container
    position: left
    behavior: full-bleed
    child:
      role: image
      aspect_ratio: 3/4
      object_fit: cover
      style: classic-vignette   # 古典暗角效果

  - role: text-container
    position: right
    padding: [120, 80, 120, 80]  # top right bottom left
    align: center-left
    children:
      - role: eyebrow
        tag: p
        style: uppercase-tracking-wider
      - role: headline
        tag: h1
        style: display-large
        size: clamp(48, 6vw, 96)
      - role: body
        tag: p
        style: body-lead
        max_width: 50ch
      - role: cta
        tag: a
        style: button-outline

# ─── CSS实现 ───
css:
  container: |
    display: grid;
    grid-template-columns: 5fr 7fr;
    min-height: 100vh;
  image_col: |
    position: relative;
    overflow: hidden;
    &::after { content: ''; position: absolute; inset: 0;
      background: linear-gradient(90deg, rgba(0,0,0,0.15) 0%, transparent 50%); }
  text_col: |
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 120px 80px;

# ─── 质量门禁 ───
quality:
  min_font_size: 16px
  max_line_length: 75ch
  min_contrast_ratio: 4.5:1
  min_whitespace_percent: 10
  max_elements_per_section: 7

# ─── 案例参考（用于学习/检索） ───
examples:
  - brand: Anthropic
    url: https://anthropic.com
    screenshot: refs/examples/anthropic-hero.png
    notes: "左图右文，奶油底+赤陶橙 accent"
  - brand: Stripe
    url: https://stripe.com
    screenshot: refs/examples/stripe-hero.png
    notes: "斜切渐变hero，结构化栅格"

# ─── 进化日志 ───
changelog:
  "1.0.0": "初始版本，从 Anthropic/Stripe 布局迁移"
```

### 配色模板 (color-scheme)

```yaml
id: warm-editorial-cream
name: 暖色出版物 — 奶油纸底+赤陶橙
version: "1.0.0"
category: color-scheme/warm

colors:
  background: "#F5F0E8"      # 奶油纸底
  text_primary: "#191919"     # 近黑
  text_secondary: "#6B6258"   # 深灰褐
  accent: "#CC785C"           # 赤陶橙
  accent_hover: "#B8654A"
  border: "#E8DFD4"           # 极淡米灰
  success: "#4A7C59"
  error: "#C44A4A"

# 使用场景
when:
  brand_tone: [warm, editorial, natural, premium]
  audience: [readers, professionals, knowledge_workers]

# 可访问性校验
contrast:
  background_text_primary: "6.8:1 ✅"
  background_text_secondary: "3.5:1 ⚠️（仅用于装饰）"
  background_accent: "2.1:1 ❌（accent仅用于大块/按钮，不用于文字）"

# CSS变量
css_variables: |
  --bg: #F5F0E8;
  --text-primary: #191919;
  --text-secondary: #6B6258;
  --accent: #CC785C;
  --accent-hover: #B8654A;
  --border: #E8DFD4;
```

### 动画模板 (animation)

```yaml
id: hero-reveal-stagger
name: Hero渐进式展开 — 逐元素错位入场
version: "1.0.0"
category: animation/hero

# 时间轴
timeline:
  total_duration: 2.0  # 秒
  elements:
    - id: background
      type: fade-in
      start: 0.0
      duration: 0.8
      easing: ease-out
    
    - id: headline
      type: slide-up
      start: 0.3
      duration: 0.6
      easing: cubic-bezier(0.16, 1, 0.3, 1)
      offset: 40px
    
    - id: subtitle
      type: slide-up
      start: 0.7
      duration: 0.5
      easing: cubic-bezier(0.16, 1, 0.3, 1)
      offset: 30px
    
    - id: cta
      type: fade-in
      start: 1.1
      duration: 0.4
      easing: ease-out

# CSS实现
css: |
  @keyframes slide-up {
    from { opacity: 0; transform: translateY(var(--offset)); }
    to { opacity: 1; transform: translateY(0); }
  }
  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }

# 质量
quality:
  max_stagger_gap: 0.4s  # 最大错位间隔，防止观众等待
  min_element_duration: 0.3s  # 最小动画时长
```

### 视频剪辑模板 (video-edit)

```yaml
id: launch-film-30s
name: 30秒品牌宣传片 — 苹果发布会级
version: "1.0.0"
category: video/launch-film

specs:
  total_duration: 30       # 秒
  frame_rate: 25
  resolution: [1920, 1080]

# 叙事弧
arc:
  - phase: hook            # 0-5s 钩子
    duration: 5
    mood: mysterious
    music: bgm-tech.mp3
    sfx_density: high      # 6个/10秒
    
  - phase: reveal          # 5-15s 产品揭示
    duration: 10
    mood: aspirational
    music: bgm-tech.mp3
    sfx_density: medium    # 3个/10秒
    
  - phase: feature         # 15-25s 功能展示
    duration: 10
    mood: confident
    music: bgm-educational.mp3
    sfx_density: low       # 1-2个/10秒
    
  - phase: closing         # 25-30s 品牌收尾
    duration: 5
    mood: inspiring
    music: bgm-educational.mp3
    sfx_density: high

# 转场规则
transitions:
  hook_to_reveal: "cross-fade 0.8s"
  reveal_to_feature: "slide-left 0.5s"
  feature_to_closing: "fade-to-black 1.0s + logo stamp sfx"

# 音效配方（参考 audio-design-rules.md）
sfx_recipe: A  # 发布hero型：密集SFX + 低频BGM

# 模板来源
source: "references/launch-film-director-notes.md"
```

## 使用流程

### 设计时（LLM调用）
```
LLM确定用户需求 → 查询 template-registry.yaml 匹配模板
→ 读取对应YAML → 注入用户内容 → 按CSS字段生成HTML
→ 通过 quality 字段自检
```

### 进化时（反馈驱动）
```
用户反馈“好” → pattern-extractor.py 提取该产品的视觉特征
→ 生成新模板候选 → 人工审校 → 合并到注册表 + 版本号+1
```

### 质量门禁执行
```
preflight.py 检查环境 → quality-gate.py 验证CSS/对比度/字号
→ 输出报告（通过/警告/失败）
```
