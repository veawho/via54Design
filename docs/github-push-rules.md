# GitHub Push Rules — via54Design

> **核心原则**: 仓库只存**源码 + 文档 + 模板 + 配置**. 测试产物/中间产物**永远不**进仓库.

---

## 🚫 绝对不推

| 类型 | 例子 | 原因 |
|---|---|---|
| **测试视频** | Pexels stock mp4 (69 个, 710MB), ffmpeg 中间渲染产物 | 体积大 + 可再生 (脚本拉) |
| **测试图片** | 渲染的 PNG 字幕 (30+), vision 验证抽帧 (15 张 jpg) | 可再生 (gen_subtitle_*.py) |
| **生成产物** | 编译二进制 (`via54.exe`), ffmpeg 渲染 mp4 | 编译产物 |
| **API 缓存** | `minimax-output/` (mmx 调用的所有输出) | mmx 缓存, 重生成即可 |
| **本地配置** | `.env`, `.git-credentials`, IDE 缓存 | 含密钥 |
| **测试报告** | `self_test_reports/*.json`, `frames_*.jpg` | 一次性产物 |
| **下载的素材** | `web/uploads/`, `_scripts/stock/` | 体积大, 拉取即可 |

## ✅ 应该推

| 类型 | 例子 |
|---|---|
| **Go 源码** | `cmd/`, `internal/`, `*.go` |
| **Python 脚本** | `_scripts/*.py` (生成脚本, 不是产物) |
| **模板/配置** | `_scripts/spec/*.json5`, `templates/` |
| **文档** | `docs/*.md`, `README.md`, `CHANGELOG.md` |
| **CI/规则** | `.github/`, `Makefile`, `*.yml` |
| **本地脚本** | `scripts/` (跑构建的脚本) |
| **AGENTS/SOUL** | `AGENTS.md`, `SOUL.md`, `.cursorrules` |

---

## 🔄 标准推送流程 (按 GitHub 规则)

### 1. 本地先 .gitignore 验证

```bash
# 看 untracked + staged (应只含源码/文档/模板)
git status -s

# 看跟踪中的大文件 (> 1MB)
git rev-list --objects --all | \
  git cat-file --batch-check='%(objectname) %(objecttype) %(objectsize) %(rest)' | \
  awk '/^blob/ {print $3, $4}' | sort -n | tail -20
```

### 2. Conventional Commits 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

| type | 用途 |
|---|---|
| `feat` | 新功能 |
| `fix` | bug 修复 |
| `docs` | 文档 |
| `refactor` | 重构 |
| `chore` | 杂项 (build/deps) |
| `ci` | CI 配置 |
| `style` | 格式 |
| `test` | 测试 |

**本项目例子**:
- `feat(v0.7.0): 锂电 30s 视频 (count-up + 3 语 + 8.50⭐)`
- `fix(v0.7.1): en 时长归一化 28.925s→30.000s (apad + atrim + -t 30)`
- `chore(cleanup): 移除 _scripts/stock (710MB) + _scripts/subtitles (15MB)`

### 3. 推送走 PR 流程 (保护 main)

```bash
# 1) 创建 feature 分支
git checkout -b feat/<name>

# 2) 提交
git add -p  # interactive 二次确认
git commit -m "feat(<scope>): <subject>"

# 3) 推自己 fork
git push origin feat/<name>

# 4) GitHub 上开 PR
gh pr create --title "feat: <subject>" --body "..."

# 5) CI 跑过 + review 后 merge
gh pr merge --squash
```

**绝对不要直接 push main** (就算 admin 也别 — 触发保护, 也坏规矩).

### 4. GitHub Release 规则

| 推/不推 | 类型 |
|---|---|
| ✅ 可推 | 编译二进制 (`via54.exe`, `via54-mcp.exe`) |
| ✅ 可推 | 文档 PDF/HTML |
| ❌ **绝对不** | 演示视频/测试产物 (体积大, 另存 Releases/对象存储) |
| ❌ **绝对不** | 抽帧截图/PNG |

**视频/PNG 应该**: 留在本地 `minimax-output/` + 用户自己保存, 不进 Git.

### 5. Commit 前自检清单

- [ ] `git status -s` 全是源码/文档/模板 (无 `*.mp4`/`*.png`/`*.exe`)
- [ ] `git diff --stat` 大小合理 (< 几百 KB)
- [ ] 改动有对应 commit message (Conventional Commits)
- [ ] 没有混 .gitignore 中列的文件
- [ ] 之前 commit 没遗留 .DS_Store / Thumbs.db

---

## 🚨 历史教训 (2026-06-12)

### ❌ 已发生的违规

1. **691MB stock mp4 推送** — `_scripts/stock/` 69 个文件 (4.7MB+39.8MB+...)
2. **15MB PNG 字幕推送** — `_scripts/subtitles/` 30 个 PNG
3. **v0.7.0-beta Release 4 mp4** — `lithium_30s_v7_*.mp4` 4 个 8.5-8.9MB

### ✅ 修正动作

1. 删 v0.7.0-beta Release + 4 assets (HTTP 204 ×5)
2. 本地 `git rm -rf --cached _scripts/stock _scripts/subtitles`
3. 加 `.gitignore` 规则: `_scripts/stock/`, `_scripts/subtitles/`, `minimax-output/`, `frames_*/`
4. 本文档 `docs/github-push-rules.md` 作永久规范

### 📊 仓库净化前/后

| | 修复前 | 修复后 |
|---|---|---|
| 跟踪文件数 | 24,210+ | ~80 |
| 仓库体积 | ~750MB | ~2MB |
| Release 数 | 3 | 2 (删 v0.7.0-beta) |

---

## 🛠 工具命令 (快速用)

```bash
# 1. 找仓库内 > 1MB 的大文件
git rev-list --objects --all | \
  git cat-file --batch-check='%(objectname) %(objecttype) %(objectsize) %(rest)' | \
  awk '/^blob/ {print $3/1024/1024 "MB", $4}' | sort -rn | head -20

# 2. 找仓库内 mp4/png/jpg
git ls-files | grep -E '\.(mp4|png|jpg|jpeg|gif|webp)$'

# 3. 看 .gitignore 哪些没匹配
git check-ignore -v $(git ls-files | head -20)

# 4. dry-run add (看会进什么)
git add --dry-run -A
```

---

**Last updated**: 2026-06-12
**Enforcer**: 任何 .go/.py 之外的二进制/媒体文件 → 先问 "能再生吗?" → 能就 .gitignore, 不能就单独存储
