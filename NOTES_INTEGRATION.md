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
