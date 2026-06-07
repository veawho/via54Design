# via54Design

> 结构化模板引擎 + 设计方向顾问 + 质量门禁
>
> 把人类的审美判断力，转化为机器可执行的确定性模板。

## 🚀 一句话部署

AI Agent / AI IDE 自动识别这条命令：

```bash
bash <(curl -s https://raw.githubusercontent.com/veawho/via54Design/main/scripts/install.sh)
```

或在终端手动执行：

```bash
curl -fsSL https://raw.githubusercontent.com/veawho/via54Design/main/scripts/install.sh | bash
```

部署后输入 `via54` 即可使用。

---

## 🗣️ 自然语言操作示例

### 给 AI 助手（Claude / Cursor / Copilot）的指令

把下面任意一句话发给 AI，它能自动完成：

| 你想做什么 | 对 AI 说 |
|-----------|----------|
| **做 HTML** | "用 via54Design 生成一个暖色编辑风格的页面，标题叫'羽图鉴'" |
| **检查质量** | "帮我检查这个 HTML 文件的质量: via54 quality --html demo.html" |
| **提取风格** | "从我的 HTML 里提取配色和字体方案: via54 pattern --html demo.html" |
| **启动服务** | "启动 via54Design 的 MCP Server，我要在 Claude Desktop 里用" |
| **加背景音乐** | "给这个视频配一首科技感 BGM: via54 media add-music demo.mp4 --mood=tech" |
| **转 GIF** | "把视频转成 GIF: via54 media convert demo.mp4" |
| **取素材图** | "帮我找一些鹦鹉的公共领域插画: via54 media fetch --query parrot" |
| **录视频** | "把这个 HTML 录成 30 秒 1080p 视频: via54 export render demo.html --duration 30" |
| **转 PDF** | "把这个页面导出为 PDF: via54 export pdf demo.html" |
| **语音合成** | "把这段文字转成语音: via54 export tts --text '你好世界' --out hello.mp3" |
| **列出模板** | "看看有哪些设计模板可以用: via54 list" |

### 组合示例

```
"帮我用 via54Design 做一个暖色编辑风格的品牌页面，检查质量，
 然后录成 30 秒视频，再加一段科技感背景音乐。"

→ 实际执行的命令:
  via54 generate --layout hero-split --color warm-editorial --font serif-sans --title "品牌页" --output brand.html
  via54 quality --html brand.html
  via54 export render brand.html --duration 30
  via54 media add-music brand.mp4 --mood=tech
```

---

## 📦 手动安装

### 依赖

| 工具 | 用途 | 安装 |
|------|------|------|
| Go 1.21+ | 核心引擎 | [go.dev/dl](https://go.dev/dl/) |
| ffmpeg | 视频/音频处理 | `brew install ffmpeg` / `apt install ffmpeg` |
| Node.js 18+ | PPTX 导出 | [nodejs.org](https://nodejs.org/) |

```bash
# 编译
go build -o via54 ./cmd/huashu/

# 安装 Playwright 浏览器
npx playwright install chromium

# 运行
./via54
```

---

## 💻 命令参考

```bash
via54                              # 帮助
via54 version                      # 版本信息
via54 list                         # 列出所有模板

# 设计模板
via54 generate --layout <id> --color <id> --font <id> --title "标题" --output out.html
via54 quality --html out.html
via54 pattern --html out.html --name "项目名"

# 媒体管线
via54 media add-music input.mp4 --mood tech|ad|educational|tutorial
via54 media convert input.mp4
via54 media fetch --query "关键词" --out ./img --count 3

# 导出
via54 export render input.html --duration 30 --width 1920 --height 1080
via54 export pdf input.html
via54 export tts --text "你好" --out voice.mp3

# MCP Server (兼容 Claude Desktop / Cursor / Copilot)
via54 serve
```

### MCP 配置

**Claude Desktop** (`claude_desktop_config.json`):
```json
{ "mcpServers": { "via54Design": { "command": "via54", "args": ["serve"] } } }
```

**Cursor** (`.cursor/mcp.json`):
```json
{ "mcpServers": { "via54Design": { "command": "via54", "args": ["serve"] } } }
```

---

## 📐 架构

```
┌──────────────────────┐
│   你的审美知识        │ ← 不可替代的人类优势
│   → YAML 模板        │ ← 结构化、可复用、可版本化
├──────────────────────┤
│   Go 核心引擎        │ ← 12MB 单二进制, 13 命令
├──────────────────────┤
│   MCP Server          │ ← Claude / Cursor / Copilot / Hermes
└──────────────────────┘
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

## 致谢

基于 [huashu-design](https://github.com/alchaincyf/huashu-design) by alchaincyf（花叔·花生），MIT 许可。
