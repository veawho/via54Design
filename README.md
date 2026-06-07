# via54Design

> 花叔Design 二次开发 — 结构化模板引擎 + 设计方向顾问 + 质量门禁 + 多模态风格学习
>
> 把人类的审美判断力，转化为机器可执行的确定性模板。

## 定位

**不是「AI 设计工具」**。AI 做不好设计决策。

**而是「把你的审美知识模板化」的工具**。

你（人类）判断什么好看、为什么好看、什么情况下用什么风格。工具把这些判断固化为 **结构化 YAML 模板**——布局、配色、字体、动画、视频节奏——然后确定性执行，每次输出一致。

## 架构

```
┌──────────────────────┐
│   你的审美知识        │ ← 不可替代的人类优势
│   → YAML 模板        │ ← 结构化、可复用、可版本化
├──────────────────────┤
│   Go 核心引擎        │ ← 单二进制 MCP Server + CLI
│   (template-applier) │
├──────────────────────┤
│   Shell 媒体管线      │ ← ffmpeg + Playwright（不改）
├──────────────────────┤
│   MCP Server          │ ← 兼容 Claude Desktop / Cursor / Copilot / Hermes
└──────────────────────┘
```

## 许可

双许可：**MIT OR AGPL-3.0**

- `templates/` YAML 模板：MIT
- `scripts/` Shell 脚本：MIT
- `internal/` Go 源代码：AGPL-3.0
- `references/` 文档：CC BY 4.0

## 快速开始

```bash
# MCP Server 模式（推荐）
go run ./cmd/huashu serve

# CLI 模式
go run ./cmd/huashu generate --layout hero-split --color warm-editorial --font serif-sans
```
