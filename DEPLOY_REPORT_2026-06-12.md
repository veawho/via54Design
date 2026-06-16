# via54Design 部署报告

**部署日期**: 2026-06-12
**部署模式**: 最小化（本地 build · 仅源码 + 二进制，不注册 launchd / 不放 `/usr/local/bin`）
**部署位置**: `/Users/david/Desktop/developments/via54Design/`
**上游仓库**: https://github.com/veawho/via54Design
**部署人**: Hermes Agent（受用户委托）

---

## 1. 环境快照

| 项 | 值 |
|---|---|
| OS | macOS 26.5.1 · Darwin arm64 |
| Go | 1.26.3 darwin/arm64 (`/opt/homebrew/bin/go`) |
| Git | 2.50.1 |
| Python | `/Library/Developer/CommandLineTools/usr/bin/python3` (system 3.x) |
| 磁盘可用 | 31 GB |
| 文件系统 | APFS (`/` → GOFLAGS 空, 不需要 `-buildvcs=false`) |

---

## 2. 部署步骤

### 2.1 克隆
```bash
cd /Users/david/Desktop/developments
git clone https://github.com/veawho/via54Design.git
```
- HEAD: `f1496aa v0.7.1: 时长归一化...`
- 含 436 文件, .git 干净

### 2.2 资源解析修复（已合并到 main）
PR #4 (`fix(util): FindBaseDir support VIA54_BASE_DIR env override`) 状态: **CLOSED（已合并）**。
Git log 显示 commit `5a15cb0` + `4ca79dd` 已在 main 上，`internal/util/paths.go:14-18` 已含 env override：

```go
//  1. VIA54_BASE_DIR 环境变量 (Mac/Linux 安装到 /usr/local/bin 时必须)
if env := os.Getenv("VIA54_BASE_DIR"); env != "" {
    return env
}
```

无需重新打 patch。**部署时使用 `VIA54_BASE_DIR` env var 指向本仓库路径即可。**

### 2.3 本地 build
```bash
cd /Users/david/Desktop/developments/via54Design
mkdir -p bin
go build -o bin/via54 ./cmd/via54/        # 18.1 MB
go build -o bin/via54-mcp ./cmd/mcp-server/  # 16.5 MB
```
- 两二进制均 `chmod +x` 自动
- Go stat-cache 噪音（`permission denied` on `~/go/pkg/mod/cache/download`）属已知现象，不影响产物

### 2.4 烟测 10 项（全 PASS）
| # | 测试 | RC | 详情 |
|---|------|----|----|
| 1 | version | 0 | `via54Design dev` |
| 2 | list | 0 | 71 行（12 字体 + 4 叙事模型 + 3 布局 + 30+ 配色） |
| 3 | prompt --platform midjourney | 0 | 1235 chars, md5=43cd019c |
| 4 | prompt 确定性 | 0 | 同输入 → 同 md5 |
| 5 | generate --layout hero-split-16-9 --color muji-minimal --font hei-modern | 0 | output.html 7660 bytes |
| 6 | export pptx | 0 | output.pptx 4828 bytes |
| 7 | export svg | 0 | svg-scenes/scene-001-sample.svg |
| 8 | export json | 0 | (impl PASS) |
| 9 | quality --html output.html | 0 | PASS · 0 errors / 1 warning / 1 info |
| 10 | narrate --seed ... --model three-act | 0 | 完整 markdown 剧本 |

### 2.5 20 轮 e2e（20/20 PASS）
跑 `test_20_rounds.py`（v2 跨平台版，自动用 `via54` 不带 `.exe`）：

```
  [01-03] 部署/编译/体积     3/3 PASS (build 3.7s + 1.9s + vet 5.4s)
  [04-09] 冒烟 6 命令        6/6 PASS
  [10-12] 准确性确定性       3/3 PASS (narrate md5=84fad168, prompt=06424778, generate=94011a63)
  [13-15] 边界               3/3 PASS (空 scene rc=1 优雅拒绝, 超长 1200→2424 OK)
  [16-18] 压力               3/3 PASS (prompt 20/20 @33ms, narrate 10/10, generate 5/5)
  [19-20] Web + HTMX         2/2 PASS (Web UI HTTP 200, 7770B)
  ────────────
  TOTAL: 20/20 PASS · 0 WARN · 0 FAIL
```

### 2.6 MCP server handshake（额外验证）
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"0.1"}}}' | bin/via54-mcp
```
响应（合法 JSON-RPC 2.0）：
```json
{
  "jsonrpc":"2.0","id":1,"result":{
    "protocolVersion":"2024-11-05",
    "capabilities":{"tools":{"listChanged":true}},
    "serverInfo":{"name":"via54Design","version":"0.2.0"}
  }
}
```
✅ MCP stdio 模式可用。Claude Desktop / Cursor / VS Code 可直接配 stdio mode:
```json
{
  "mcpServers": {
    "via54Design": {
      "command": "/Users/david/Desktop/developments/via54Design/bin/via54-mcp",
      "env": {"VIA54_BASE_DIR": "/Users/david/Desktop/developments/via54Design"}
    }
  }
}
```

---

## 3. 项目目录快照（部署后）

```
/Users/david/Desktop/developments/via54Design/
├── bin/                         ← 本次 build 产物
│   ├── via54        (18 MB)     CLI (13 子命令: serve/generate/narrate/quality/...)
│   └── via54-mcp    (16 MB)     MCP Server (stdio / SSE HTTP)
├── cmd/                         上游源码
│   ├── via54/
│   └── mcp-server/
├── internal/                    核心引擎
├── templates/                   YAML 模板库 (layouts/color-schemes/typography/narratology/prompts)
├── web/                         HTTP UI (handler.go + templates/)
├── hack/                        build.sh / install.sh / wasm/
├── docs/                        template-format.md / deployment-guide.md / prompts/
├── test_samples/                示例输入
├── self_test_*.py               100/200 轮压力测试 (可选)
├── test_20_rounds.py            ✅ 20/20 PASS
├── AGENTS.md / SOUL.md          AI 工作上下文
├── Makefile                     fs-check + build + build-mcp + cross
├── README.md / README_EN.md
└── DEPLOY_REPORT_2026-06-12.md  ← 本文件
```

---

## 4. 使用方式（最小化部署版）

### 4.1 CLI
```bash
export VIA54_BASE_DIR=/Users/david/Desktop/developments/via54Design
alias via54=/Users/david/Desktop/developments/via54Design/bin/via54

via54 version
via54 list                                    # 列出全部模板
via54 prompt --scene "清晨咖啡" --platform midjourney
via54 generate --layout hero-split-16-9 --color muji-minimal --font hei-modern
via54 export pptx                              # → ./output.pptx
via54 export svg                               # → ./svg-scenes/
via54 quality --html output.html
via54 narrate --seed "城市夜景" --model three-act
```

### 4.2 MCP Server (stdio)
直接调二进制即可（见 2.6 配置片段）。

### 4.3 Web UI (可选启动)
```bash
export VIA54_BASE_DIR=/Users/david/Desktop/developments/via54Design
cd $VIA54_BASE_DIR
./bin/via54 web                                # 默认 :8080 (或 via54 web --port XXXX)
# 浏览器访问 http://localhost:8080
```
> 当前未注册 launchd。需要长期后台运行时再用 `software-deployment-cli-binary-macos` skill 的 Phase 5。

---

## 5. 与 v15 全局部署（`/usr/local/bin/`）对比

| 维度 | v15 全局部署 (2026-06-09) | 本次最小化 (2026-06-12) |
|------|---------------------------|------------------------|
| 二进制位置 | `/usr/local/bin/{via54,via54-mcp}` | `./bin/` (本仓库) |
| launchd 注册 | ✅ com.david.via54-web / via54-mcp | ❌ 不注册 |
| Web 端口 | 8080 (via54-web), 8090 (via54-mcp SSE) | 不占用端口（按需启动） |
| PATH 污染 | 是（全局可见 `via54`） | 否（需 alias 或绝对路径） |
| 卸载成本 | `launchctl bootout` + `sudo rm` | 仅 `rm -rf` 目录 |
| 适用场景 | 日常依赖、开机自启 | 临时验证、PR review、调试、PR 部署复现 |

两份部署互不冲突，可同时存在（PATH 优先取 `/usr/local/bin`）。

---

## 6. 已知问题 / Pitfalls（macOS 部署）

1. **`go build` stat-cache 噪音**：`permission denied` on `~/go/pkg/mod/cache/download` 是 Go 1.20+ 默认行为，不影响产物，可忽略。
2. **PATH 不含 `/opt/homebrew/bin`**：launchd 启动时若 PATH 是 sanitized 默认（`/usr/bin:/bin:/usr/sbin:/sbin`），找不到 `go`/`node`。本次未注册 launchd，**不踩**。
3. **VIA54_BASE_DIR 必须设**：因为本部署没装到 `/usr/local/bin/`，但 `FindBaseDir` 仍会向上找 `templates/` —— 在 `/Users/david/Desktop/developments/via54Design/bin/via54` 向上找能找到，但保险起见**还是设 env**。
4. **布局名带 `-16-9` 后缀**：`hero-split-16-9` / `bento-grid-2x2` / `gallery-waterfall`，README 旧示例写的 `hero` 是错的。
5. **`export pptx` 的 `--output` 被忽略**：实际写到 `CWD/output.pptx`。本次未踩（直接 cd /tmp）。
6. **`export svg` 写到 `./svg-scenes/` 不是 `--output`**：本次踩了（svg_dirs 在 `/tmp/svg-scenes/`）。
7. **WASM 未构建**：`WASM: ❌ (cd hack/wasm && bash build.sh)` —— 需要 Rust toolchain，本次跳过。

---

## 7. 卸载

```bash
# 仅最小化部署版
rm -rf /Users/david/Desktop/developments/via54Design

# 完整清理（含 v15 全局部署）
launchctl bootout gui/501/com.david.via54-web 2>/dev/null
launchctl bootout gui/501/com.david.via54-mcp 2>/dev/null
rm ~/Library/LaunchAgents/com.david.via54-*.plist
sudo rm /usr/local/bin/via54 /usr/local/bin/via54-mcp
rm -rf /Users/david/Desktop/developments/via54Design
```

---

## 8. 验证 checklist

- [x] 二进制 build 成功 (`via54` 18MB + `via54-mcp` 16MB)
- [x] CLI version / help 正常
- [x] 10 项烟测全 PASS
- [x] 20/20 轮 e2e 全 PASS
- [x] MCP stdio initialize handshake 合法响应
- [x] Web UI HTTP 200（Round 19）
- [x] PR #4 fix 已在 main 上, 无需重打 patch
- [x] `VIA54_BASE_DIR` env override 路径验证

**结论**: 项目可立即使用。
