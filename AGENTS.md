# via54Design — AI 工作上下文

> 版本: v3.0 | 最后更新: 2026-06-08
> 由 CLAUDE.md 迁移至此。多AI工具兼容 — Claude Code / Cursor / Windsurf / Copilot / Cline 均自动发现。
> 同步配置文件: `.cursorrules` (Cursor) · `.windsurfrules` (Windsurf) · `.github/copilot-instructions.md` (Copilot)

---

## 项目定位

**via54Design** — 把人类的审美判断力固化为结构化 YAML 模板，用 Go 核心引擎确定性执行。

不是设计工具，是 **设计模板引擎 + 叙事引擎 + 媒体管线**。CLI 优先，MCP Server 第二。

---

## 谁在开发

| 角色 | 身份 | 职责 |
|------|------|------|
| **巫师叔叔 (via54)** | 4A产品策略总监 | 设计决策、审美判断、YAML模板编写 |
| **AI助手** | 开发助手 | Go代码实现、Shell脚本、测试、文档 |

**关键关系**: 人类提供灵感种子，AI结构化扩展。循环: 人写一句 → AI脚手架 → 人确认 → AI生成。

---

## AI 工具兼容表

| 工具 | 自动读取的文件 | 状态 |
|------|---------------|------|
| **Claude Code** | `AGENTS.md` | ✅ 主上下文 |
| **Cursor** | `.cursorrules` → `AGENTS.md` | ✅ 专用规则 |
| **Windsurf** | `.windsurfrules` → `AGENTS.md` | ✅ 专用规则 |
| **GitHub Copilot** | `.github/copilot-instructions.md` | ✅ 专用指令 |
| **Cline / CLI** | `AGENTS.md` | ✅ 兼容 |
| **Continue.dev** | `AGENTS.md` / `.continuerc.json` | ✅ 兼容 |
| **Aider** | `CONVENTIONS.md` (软链接) | ✅ 可手动链接 |
| **Codex CLI** | `AGENTS.md` | ✅ 兼容 |
| **Hermes Agent** | `SOUL.md` | ✅ 专用灵魂定义 |
| **OpenClaw** | `agent_routing.yaml` + `kanban` | ✅ 多Agent编排层 |
| **CrewAI / AutoGen / LangGraph** | 系统提示词 / Agent 定义 | 🔧 需复制到 Agent 配置 |

---

## 与 Agent 框架集成

### Hermes Agent (主Agent → OpenClaw 编排层)

Hermes Agent 是 via54Design 的默认开发环境。集成方式：

```bash
# 1. SOUL.md 自动注入 system prompt（项目根部自动发现）
# 2. 工具调用: 通过 terminal 执行 via54 CLI
cd /c/Users/via54/AppData/Local/Temp/via54Design
./via54.exe prompt --scene "..." --platform midjourney

# 3. 子Agent派单: 通过 lab_dispatch.py 路由到对应lab
python ~/.hermes/bin/lab_dispatch.py prdlab "Run via54 stress test"

# 4. cron 任务自动化
# 在 Hermes cron jobs.json 中注册定时生成任务
```

**SOUL.md** 位于项目根部，Hermes 自动发现并注入。内容涵盖:
- 核心身份 + 行为准则
- 技术边界 + 关键约束
- 常用工作流 + 项目状态

### OpenClaw (多Agent编排层)

OpenClaw 是 Hermes 的多 Agent 编排层，管理 7 lab + 6 template agents。

```yaml
# agent_routing.yaml 示例 (OpenClaw 路由表)
via54-design:
  description: "via54Design 设计模板引擎"
  workdir: "/c/Users/via54/AppData/Local/Temp/via54Design"
  tools: [terminal, file, web]
  triggers:
    - "via54"
    - "prompt生成"
    - "设计模板"
```

通过 `kanban` 跨 Agent 协作:
```bash
# 跨Agent任务分配
python ~/.hermes/bin/kanban create \
  --task "为via54Design添加Flux平台" \
  --assignee prdlab \
  --context "templates/prompts/flux.yaml already exists, need Go command update"
```

### CrewAI (多Agent编排)
```python
# 在 CrewAI Agent 定义中引用此上下文
agent = Agent(
    role="via54Design Developer",
    goal="Execute design template generation with deterministic Go output",
    backstory="You operate the via54Design engine. See AGENTS.md for full context.",
    allow_delegation=True,
    tools=[Via54CLI()]
)
```

### AutoGen / AG2
```python
# 在 AutoGen AssistantAgent system_prompt 中注入
with open("AGENTS.md") as f:
    context = f.read()
agent = AssistantAgent(
    name="via54Assistant",
    system_message=f"You are the via54Design AI assistant.\n\n{context}"
)
```

### LangGraph / LangChain
```python
# 在 LangGraph 节点中作为 system_prompt 注入
from langchain_core.prompts import ChatPromptTemplate
with open("AGENTS.md") as f:
    agents_context = f.read()
prompt = ChatPromptTemplate.from_messages([
    ("system", f"You are the via54Design development assistant.\n\n{agents_context}"),
    ("human", "{input}")
])
```

### MCP Server 模式
via54Design 自带 MCP Server (`cmd/mcp-server/`)，支持 MCP 协议的工具均可直接调用:
```bash
# Claude Desktop / Cursor / VS Code + MCP
via54 serve  # 启动 stdio MCP Server
# 提供: template.generate, template.list, prompt.generate, quality.assess, narrate.create
```

---

## 技术栈

| 层 | 语言 | 用途 | 外部依赖 |
|----|------|------|---------|
| 核心引擎 | **Go** | CLI + MCP Server + 模板引擎 + 叙事引擎 + 提示词引擎 + 质量门禁 | 仅 `mcp-go` + `yaml.v3` |
| 设计模板 | **YAML** | 布局/配色/字体/叙事/提示词 | 纯数据文件 |
| 媒体管线 | **Shell** | ffmpeg + Playwright | 系统工具 |
| 测试套件 | **Python 3.6+** | `test_20_rounds.py` 端到端测试 (零依赖) | stdlib only |

---

## 目录结构

```
via54Design/
├── cmd/
│   ├── via54/          CLI 入口 (13个子命令)
│   └── mcp-server/     MCP Server 入口
├── internal/
│   ├── export/         导出引擎 (纯Go: pptx/svg/json/markdown/pdf/tts)
│   ├── mcp/            MCP Server 实现
│   ├── media/          媒体管线 (下载/矢量化/配乐)
│   ├── narrate/        叙事引擎 (4模型, 剧本/分镜)
│   ├── pattern/        设计模式提取
│   ├── prompt/         提示词引擎 v2.2 (6模块)
│   │   ├── types.go
│   │   ├── generator.go
│   │   ├── templates.go
│   │   ├── quality.go
│   │   ├── version.go
│   │   └── render.go
│   ├── quality/        质量门禁
│   ├── template/       模板引擎 (布局/配色/字体)
│   └── wasm/           WASM桥接 (Rust)
├── templates/
│   ├── prompts/        4平台提示词模板 YAML
│   ├── layouts/        3布局模板 (16:9, 四端响应)
│   ├── color-schemes/  30+配色方案
│   ├── typography/     12字体定义
│   ├── narratology/    4叙事模型
│   └── registry.yaml   模板注册表
├── web/                Web界面
│   ├── handler.go      28KB HTTP处理器
│   └── templates/      HTML模板 + JS
├── hack/               构建/部署脚本
│   ├── build.sh        Go跨平台编译 (CLI+MCP)
│   ├── install.sh      一键部署
│   ├── setup.sh        完整安装
│   └── wasm/           Rust WASM源码
├── docs/
│   ├── prompts/        镜头/布光/配色/构图参考
│   ├── template-format.md
│   ├── failure-recovery.md
│   └── deployment-guide.md
├── test_samples/       测试样本
├── AGENTS.md           ← 本文件 (AI工作上下文)
├── SOUL.md             Hermes灵魂定义
├── Makefile            标准构建自动化
├── LICENSE             双许可 (MIT OR AGPL-3.0)
└── README.md           项目文档
```

---

## 架构原则

```
用户请求
  │
  ▼
via54 (Go CLI) ← 结构化执行: 模板组合/叙事/提示词/导出/质量门禁
  │
  ├── templates/YAML  ← 所有设计数据由YAML驱动
  ├── export/         ← PPTX/SVG/JSON/Markdown 纯Go实现
  └── quality/        ← 质量评分+问题清单
```

- **Go 纯二进制**: 核心引擎 15MB, 零运行时依赖
- **YAML 是唯一数据源**: 添加新平台/配色/字体只需新增 YAML 文件
- **确定性优先**: 所有 map 遍历前用 `sortedKeys()` 排序，输出 md5 可复现
- **测试脚本**: `test_20_rounds.py` (Python 3.6+ stdlib only) 端到端验证

---

## 关键约束

### 代码风格
- Go: 标准库优先，少加第三方依赖
- 错误处理: `fmt.Errorf("xxx: %w", err)` 带上下文
- 所有 map 遍历必须用 `sortedKeys()` 泛型函数保证确定性
- `internal/` + `cmd/` Go 源码：**AGPL-3.0-only**
- 每个 .go 文件必须有 `SPDX-License-Identifier` 头部

### 测试要求
- MVP验证: 至少3轮×4维 (有效性/稳定性/准确性/完整性)
- 边界测试: 空输入、无效平台、超长字符、特殊字符
- 输出确定性: 同一输入必须产生同一 md5
- 压力测试: 200次连续生成无错误

### Git 惯例
# 提交信息: `type: 描述`
# type: `feat` / `fix` / `refactor` / `docs` / `license` / `test`
# 每次大改动后必须有 `git push`
# 变更清单附在提交信息尾部

### 常见陷阱 (AI 注意)
- `empty_scene` 返回 exit=1 是正确行为，不是 bug
- PPTX 导出用 `archive/zip` + `encoding/xml`，不用第三方库
- `baseDir()` 在 prompt_cmd.go 中未定义 → 由 `main.go` 的 `baseDir()` 提供
- `lab_dispatch.py` 派单到子 agent 时用 Python，不是 Shell
- Provider 为 `ollama/hermes/local` 时不需 API key
- **平台 ID**: 真实 ID 是 `hero-split-16-9 / bento-grid-2x2 / gallery-waterfall`，不是 `hero`
- **generate 输出**: 写到 `output.html` 文件，不是 stdout
- **项目唯一 Python 文件**: `test_20_rounds.py` (测试套件)，不是 LLM 管道。AGENTS.md 旧版提到的 `hack/via54_pipeline.py` / `scripts/img2prompt.py` 已不存在

---

## 许可证地图

| 目录 | 许可 | 说明 |
|------|------|------|
| `internal/` + `cmd/` | **AGPL-3.0-only** | Go 源码，防云服务商闭源提取 |
| `templates/` + `hack/` + `docs/` | **MIT** | 模板/脚本/文档，宽松 |
| `go.mod` | `AGPL-3.0-only OR MIT` | 模块级双重许可声明 |

---

## 当前状态

- **版本**: v0.4.0
- **Go二进制**: ~15MB, 单文件, 零外部依赖
- **测试套件**: `test_20_rounds.py` (Python 3.6+ stdlib only)
- **14平台**: midjourney/flux/dalle3/sd3/stable_diffusion/ideogram/recraft/seedance/gemini/veo/sora/kling/pika/jimeng + 3 video
- **26维度**: subject~emotion 全字段 + 权重控制
- **稳定性**: 20轮测试 0 错误, 100% 确定性
- **SPDX覆盖率**: 58/58 Go 源文件 (100%)

---

## 快速参考

### 常用命令
```bash
make build build-mcp                  # 编译 (推荐用 Makefile, 自动检测 exFAT)
make fs-check                        # 查看当前文件系统类型 + GOFLAGS 状态
make cross                           # 跨平台编译 (5 平台)
./via54.exe prompt --scene "..." --platform midjourney  # 生成提示词
./via54.exe prompt list               # 列表所有平台
python test_20_rounds.py             # 20 轮端到端测试
```

### GOFLAGS 自动检测
Makefile 会自动检测**当前目录的文件系统**:
- **NTFS / APFS / ext4 / btrfs / xfs**: GOFLAGS 为空 (VCS 嵌入正常, 无问题)
- **exFAT / FAT32 / vfat**: 自动加 `-buildvcs=false` (避免 file lock 缺失导致 git 索引损坏)

**覆盖示例**:
```bash
# 强制使用 -v 显示详细构建
GOFLAGS="-v" make build

# 禁用自动检测, 用 VCS 嵌入 (仅在 NTFS/APFS 上安全)
GOFLAGS="" make build

# 查看当前 GOFLAGS 计算结果
make fs-check
```

### 添加新平台
1. 新建 `templates/prompts/<name>.yaml` (复制现有模板改参数)
2. 更新 `cmd/via54/prompt_cmd.go` 的 `listPlatforms()` 数组
3. `go build` + 验证

---

## Bot 集成 (2026-06-16 新增)

via54Design 通过飞书 Bot 集成,完整流程见 `docs/prompts/bot-composition-flow.md`:

### A. 纯文字生成提示词
1. 飞书私聊/群聊发需求 → bot
2. bot 调 via54Design 生成 3 种构图方案
3. 用户选 1/2/3 → 生成完整英文 prompt
4. 用户可"修改" / "换平台" / "重新生成"

### B. 参考图+文字生成提示词
1. 飞书发图 + 文字需求 → bot
2. vision_analyze_tool 识别图 (minimax-cn MiniMax-M3)
3. bot 渲染中文拆解版 + **首次回复含 2 点确认逻辑**:
   - 左侧湿疹皮肤: 保留"明显红斑/皮屑"真实细节 vs 克制?
   - 右侧健康皮肤: "发光透亮"高光感 vs "自然健康"哑光感?
4. 用户回复 → 生成完整英文 prompt

### 集成点
- m12 bot: `~/.hermes/scripts/m12_full_channel_bot.py`
- 飞书 Channel SDK: 5 大能力 (policy/safety/inbound/outbound/transport)
- 14 平台 + 26 维度 + 4 叙事模型 全部 via54Design CLI 暴露
