# 管线故障处理与降级策略

> 当 via54Design 的某个环节失败时，本节记录已知故障模式 + 推荐降级路径。
> 这不是配置文件，是运行态问题的诊断参考。

## 图片获取 (`via54 media fetch`)

```
wikimedia_commons (公有领域)
  ↓ 限流或失败
unsplash_pexels (免费图库)
  ↓ 无结果或无网络
honest_placeholder (灰块 + "图待补")
```

`via54 media fetch` 内置 Wikimedia Commons → Unsplash 两级回退。
若均失败，引擎生成 `placeholder` 灰块，不做假 SVG。

## 字体加载

```
Google Fonts CDN
  ↓ CDN 不可达
系统字体栈
  ├── serif:  Noto Serif SC / Source Han Serif / STSong / SimSun
  ├── sans:   PingFang SC / Microsoft YaHei / Hiragino Sans GB
  └── mono:   JetBrains Mono / Consolas / Courier New
```

Go 引擎的 `buildFontImports()` 自动处理此降级。
CDN 可用 → 加载 Google Fonts；不可用 → CSS 使用系统字体栈。

## 视频录制 (`via54 export render`)

```
Playwright (首选，支持 headless)
  ↓ Playwright 未安装或无显示
手动打开 HTML → 浏览器截图
  ↓ 需要自动化
ffmpeg x11grab (仅 Linux 有显示环境)
```

```
via54 export render demo.html         # Playwright 录制
# 如果失败:
# 1. 手动打开 demo.html
# 2. 浏览器 DevTools 截图 / 录屏
# 3. 或: npm install playwright && npx playwright install chromium
```

## 导出管线

| 导出格式 | 引擎 | 故障模式 | 降级 |
|---------|------|---------|------|
| HTML | Go (无依赖) | 几乎不会失败 | — |
| PPTX | Go (ZIP+XML) | 磁盘空间不足 | 换路径 |
| SVG | Go | 同上 | — |
| JSON | Go | 同上 | — |
| Markdown | Go | 同上 | — |
| PDF | Go + Playwright | Playwright 未装 | 浏览器打印 |
| TTS | Go + Edge TTS | 网络不可达 | 换语音引擎或跳过 |

## 设计原则

所有降级的最高优先级是 **诚实**：

- 图片获取不到 → 灰块 + "图待补"，不做假图
- 字体加载失败 → 系统字体，不假装加载成功
- 视频渲染失败 → 告诉用户，不输出空文件
- 任何"优雅降级"不可行时 → 明确报错，不静默失败
