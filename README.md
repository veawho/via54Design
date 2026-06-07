# via54Design

> via54Design — 结构化模板引擎 + 设计方向顾问 + 质量门禁 + 多模态风格学习
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
│   (12 MB, 13 命令)   │
├──────────────────────┤
│   MCP Server          │ ← 兼容 Claude Desktop / Cursor / Copilot / Hermes
└──────────────────────┘
```

## 命令

```bash
# 设计模板
via54 generate --layout hero-split --color warm-editorial --font serif-sans --title "提案" --output output.html
via54 quality --html output.html
via54 pattern --html output.html --name "项目名"
via54 list

# 媒体管线
via54 media fetch --query "Edward Lear parrot" --out ./img
via54 media add-music input.mp4 --mood=tech
via54 media convert input.mp4

# 导出
via54 export render input.html --duration 30
via54 export pdf input.html
via54 export tts --text "你好" --out voice.mp3

# MCP Server
via54 serve
```

## 语言

| 语言 | 用途 | 文件数 |
|------|------|--------|
| Go | 核心引擎、CLI、MCP、媒体管线、导出 | 13 |
| Rust | WASM 高速模板引擎（可选编译） | 7 |
| JavaScript | PPTX 导出（唯一残留） | 3 |

## 许可

双许可：**MIT OR AGPL-3.0**

- `templates/` YAML 模板：MIT
- `scripts/` JS/Shell 脚本：MIT
- `internal/` Go 源代码：AGPL-3.0
- `internal/wasm/` Rust 源代码：AGPL-3.0
