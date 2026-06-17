# NOTES — via54Design 跟 via54Hermes 整合说明

> **整合日期**: 2026-06-15
> **来源**: per 主人原话 LarkDesign 定义 (per `via54x4-PLAN.md` 跟 CAPABILITY_MATRIX.md section 3)
> **跟本文关系**: via54Design 是 Go 设计引擎 (v0.5.0 stable / v0.6.0 video-pipeline active), **不整合 Telegram / 飞书 / Hermes gateway** (per definition)

---

## 跟 via54Hermes 关系: 0 整合 (per design)

**via54Design 定义** (主人原话): "**有 AI 创意协助/AI 编辑/快速协做加上设计系统"** — **专注设计**。

**不整合**:
- ❌ 不整合 Telegram (per via54Larkbotgo/LarkSkills 1:1 文档 "Larkbotgo LarkSkills 整合 Telegram 11 段, LarkDesign 不需要")
- ❌ 不整合 飞书 (per Larkfix 1:1 文档 "Larkfix 整合飞书 daemon, LarkDesign 不需要")
- ❌ 不整合 Hermes gateway (LarkDesign 是 standalone Go binary, 跟 Hermes 生态无关)

**5 仓库生态 (per design)**:
| 仓库 | 类型 | 整合 Hermes? |
|---|---|---|
| via54Hermes | 知识库 (Hermes Agent 运维 15 事故案例) | (真相源) |
| via54Larkbotgo | Go 飞书 bot skeleton | ✅ 11 段整合 (Telegram/PTB + Larkfix + state.db + memory) |
| via54Larkfix | Python 飞书 daemon (private) | ✅ 11 段整合 (Telegram + Lark approval + state.db + write_file) |
| via54Skills | skills 集 (4 via54* skill + 1 NEW: via54hermes-pitfalls) | ✅ 11 段 skill 速查 (11 个 Hermes 整合坑) |
| **via54Design** | **Go 设计引擎 (v0.5.0 / v0.6.0)** | **❌ 0 整合 (per design, 跟 Hermes 无关)** |

---

## LarkDesign 跟其他 4 仓库的差异 (为什么 0 整合)

via54Design 跟 Larkbotgo Larkfix LarkSkills 不同:

| 维度 | Larkbotgo Larkfix LarkSkills | LarkDesign |
|---|---|---|
| 语言 | Go (Larkbotgo) / Python (Larkfix/LarkSkills) | Go (v0.5.0+) |
| 用途 | 飞书 bot + skills 集 (跟 IM 整合) | 设计引擎 (跟 IM 无关) |
| Hermes 整合 | ✅ 11 段 pitfalls 镜像 | ❌ 0 段 (per design) |
| 运行时 | daemon + IM adapter | standalone Go binary |
| 凭证 | `~/.hermes/.env` + `~/.config/feishu/credentials.json` | `template-registry.yaml` (跟 Hermes 无关) |
| Telegram 整合 | ✅ 11 段 pitfalls 文档化 | ❌ (LarkDesign 不需要 Telegram) |
| 飞书整合 | ✅ Larkfix 跟 Hermes 共享 .env | ❌ (LarkDesign 不需要 飞书) |

**LarkDesign 跟 Hermes 唯一的"交集"**:
- 跟 Larkbotgo 共享 `via54Larkbotgo` 命名空间 (但 Larkbotgo 是 LarkDesign 不直接相关)
- 都用 `make` / `go` build 系统 (但 LarkDesign 不用 Hermes gateway)
- 都有 `.gitignore` + GitHub Actions (但 LarkDesign 的 workflow 是 `video-pipeline.yml`, 不是 VitePress deploy)

---

## 未来可能性 (per Larkfix 仓库 LarkDesign 主题)

per Larkfix `.project.yml` "飞书/Lark 分支" 段:
- **lark/macos-local 分支** 已经存在 (per LarkDesign, push 2026-06-14 12:23, initial commit)
- 这跟 **via54Larkfix 仓库** 无关 — 主人 active dev **Larkbotgo/Larkfix** 不在 LarkDesign 仓库做

**LarkDesign 跟 Larkbotgo/Larkfix 的真实关系**:
- LarkDesign 仓库的 `lark/macos-local` 分支是**主人 active 整合 macOS 的开发分支** (跟 Larkbotgo Larkfix 完全独立的本地 Lark 整合)
- 不影响 Larkbotgo Larkfix LarkSkills LarkHermes 4 仓库的 1:1 对齐

---

## LarkDesign 仓库的"11 段"等价物 (per Larkbotgo Larkfix LarkSkills 11 段)

虽然 LarkDesign 不需要 11 段 pitfalls (per design), **但** LarkDesign 有自己的 11 段 "Go 整合 pitfalls" 暗含在 `AGENTS.md` 跟 `CHANGELOG.md`:

| Larkbotgo Larkfix LarkSkills 11 段 (Hermes) | LarkDesign 等价物 (Go) |
|---|---|
| 坑 1: TELEGRAM_PROXY socks5 | (N/A LarkDesign 不需要) |
| 坑 2: chat_id vs bot_id | (N/A LarkDesign 不需要) |
| 坑 3: write_file silent fail | (N/A LarkDesign 是 Go 不用 write_file) |
| 坑 4: memory old_text vs old_string | (N/A LarkDesign 不用 memory tool) |
| 坑 5: hermes config set model flatten | (N/A LarkDesign 不用 hermes config) |
| 坑 6: PyYAML \U Unicode escape | ✅ LarkDesign AGENTS.md "v0.5.0 / v0.3.0+ macOS 修复" 段 |
| 坑 7: hermes gateway status stale | (N/A) |
| 坑 8: X-Hermes-Api-Key header | (N/A) |
| 坑 9: Lark approval skill 损坏 | (N/A LarkDesign 不用 Lark) |
| 坑 10: state.db 跟 memory_log.db | (N/A LarkDesign 不用 state.db) |
| 坑 11: via54-mcp template-registry.yaml | ✅ LarkDesign 仓库 `templates/registry.yaml` 是 LarkDesign 设计 registry 真相源 (跟 Larkbotgo 仓库 ~/.hermes 共享) |

**LarkDesign 跟 Larkbotgo Larkfix LarkSkills 唯一交集**: 坑 11 (template-registry.yaml 跟 templates/registry.yaml)。

---

## 真相源

- **本地 LarkDesign 仓库**: `~/Desktop/developments/via54Design/`
- **GitHub 仓库**: `veawho/via54Design` (default branch: `main`, active branch: `feature/video-pipeline`)
- **AGENTS.md** (22KB): 主人 LarkDesign active dev guide
- **CHANGELOG.md** (8KB): LarkDesign v0.5.0 / v0.6.0 changelog
- **README.md** (22KB) + **README_EN.md** (22KB): LarkDesign 双语 (per master 原话 "1 个仓库可以多 README")
- **lark/macos-local 分支** (push 2026-06-14 12:23): LarkDesign 主人 active dev macOS 整合
- **feature/video-pipeline** (push 2026-06-15 active): LarkDesign v0.6.0 视频产线 active dev


---

## 本轮 18 件修复 (per 2026-06-15 IM 平台统一 session) — LarkDesign 视角

> **整合日期**: 2026-06-15
> **整合来源**: 本 session 全部验证的修复 + Hermes 官方 PR + GitHub issue tracker
> **4 仓库同步**: Larkbotgo Larkfix LarkSkills LarkDesign + CAPABILITY_MATRIX

### LarkDesign 跟其他 4 仓库的 18 件修复对应关系

| 18 件修复类别 | Larkbotgo | Larkfix | LarkSkills | **LarkDesign (per design 0 整合)** |
|---|---|---|---|---|
| A. Hermes GitHub Issue + PR 修复 (3) | ✅ | ✅ | ✅ | **❌ 不需要** (Go 设计引擎不整合 Hermes gateway) |
| B. IM 平台统一 (4) | ✅ Telegram + Feishu 整合 | ✅ Feishu daemon | ✅ skill 速查 | **❌ 不需要** (per design) |
| C. 4 仓库 + 1 doc 整合 (5) | ✅ Larkbotgo 主整合 | ✅ Larkfix 整合 | ✅ LarkSkills 整合 | **⚠ 部分** (NOTES_INTEGRATION.md 文档同步, 0 代码整合) |
| D. LarkDesign 完美 sync (3) | ✅ 远端 workflow | (N/A) | ✅ 远端 workflow | **✅ LarkDesign 主** (main=feature/video-pipeline 1:1 sync) |
| E. B16 stress test (1) | ✅ 92% HTTP 200 | ✅ Larkfix daemon 验证 | ✅ LarkSkills skill 集成 | **❌ 不需要** (LarkDesign 不跑 daemon) |
| F. Cross-tool 模型路由 (1) | ✅ Hermes config | (N/A) | (N/A) | **❌ 不需要** (LarkDesign 独立 Go CLI) |

**LarkDesign 总占比**: 2/18 件相关 (C 部分文档同步 + D 部分 sync), **16/18 跟 LarkDesign 无关** (per design 0 整合)。

### LarkDesign 18 件修复中的 2 件贡献

1. **D 项 LarkDesign 完美 sync** (cddd264):
   - LarkDesign main = feature/video-pipeline 1:1
   - 8 conflict 解 (重置 + 重建 + cherry-pick)
   - LarkDesign NOTES_INTEGRATION.md push 到 main

2. **C 项 LarkDesign 文档同步** (本文件):
   - NOTES_INTEGRATION.md 5 段 (跟 via54Hermes 0 整合 per design)
   - LarkDesign Larkbotgo Larkfix LarkSkills LarkHermes 5 仓库生态表
   - LarkDesign Larkbotgo Larkfix LarkSkills LarkHermes 4 仓库镜像章节 1:1 token verify

### LarkDesign 跟 Larkbotgo Larkfix LarkSkills 同步状态

| LarkDesign 章节 | Larkbotgo Larkfix LarkSkills 对应 | 1:1 镜像 |
|---|---|---|
| NOTES_INTEGRATION.md 段 1 (0 整合 per design) | Larkbotgo hermes-pitfalls 13 段 + Larkfix references 13 段 + LarkSkills via54hermes-pitfalls 13 段 | ✅ (per design 反向镜像) |
| NOTES_INTEGRATION.md 5 仓库生态表 | Larkbotgo Larkfix LarkSkills LarkDesign LarkHermes 5 仓库 1:1 对齐 | ✅ |
| LarkDesign Larkbotgo Larkfix LarkSkills LarkHermes 4 仓库差异表 | Larkbotgo Larkfix LarkSkills LarkDesign LarkHermes 8 维度对比 | ✅ |

### B16 stress test 报告 (LarkDesign 视角)

> **来源**: `/tmp/B16_test_v2_results.txt` (50 轮 stress test, exit code 0)
> **LarkDesign 关系**: **0** (LarkDesign 是 Go 设计引擎, 不跑 m12 bot daemon)

| 指标 | 值 | LarkDesign 关系 |
|---|---|---|
| 总测试轮次 | 50 | ❌ |
| HTTP 200 OK | 46/50 (92%) | ❌ |
| HTTP 0 EXC | 4/50 (8%) | ❌ |
| m12 bot 进程 | 0 hang | ❌ |
| 修法整合 | Larkfix _send_path_degraded (commit 952892c) + Larkbotgo reference 1:1 镜像 | ❌ (LarkDesign 不跑 daemon) |

**LarkDesign 不参与 B16 stress test** — LarkDesign 是 standalone Go binary (跟 Hermes gateway 无关)。

### Larkbotgo Larkfix LarkSkills LarkDesign LarkHermes 5 仓库 18 件修复最终对齐表

| 仓库 | HEAD (本轮后) | 18 件覆盖 | 1:1 token verify |
|---|---|---|---|
| **via54Larkbotgo** | 69d4519 | 18/18 (主整合) | ✅ zh + en 镜像 |
| **via54Larkfix** | 23d4c13 | 18/18 (主整合) | ✅ Larkbotgo Larkfix 1:1 |
| **via54Skills** | a4619b8 | 18/18 (skill 速查索引) | ✅ Larkbotgo Larkfix LarkSkills 1:1 |
| **via54Design** | (NOTES 加 18 件说明, HEAD 跟远端 1:1) | 2/18 (per design 0 整合) | ⚠ 部分 (NOTES 同步) |
| **CAPABILITY_MATRIX.md** | (section 12 加) | 18/18 (跨仓库总结) | ✅ |

### LarkDesign 关键 takeaway

- ✅ LarkDesign NOTES_INTEGRATION.md 5 段 + 18 件修复 0 整合 per design (跟其他 4 仓库 100% 一致)
- ✅ LarkDesign main = feature/video-pipeline = cddd264 (1:1 完美 sync)
- ⚠ LarkDesign 不参与 IM 平台统一 / B16 stress test / Hermes gateway (per definition)
- ⚠ LarkDesign 唯一参与: LarkDesign 仓库本身的 sync (D 项) + 文档同步 (C 项)
