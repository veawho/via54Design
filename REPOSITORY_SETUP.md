# Repository settings

> 本项目 (veawho/via54Design) 的设置已通过 `gh CLI` **自动完成**。
> 本文件作为配置参考保留。

## ✅ 1. About (top-right of repo page) — **已完成**

**Description** (158 chars):
```
Design template engine + narrative engine + media pipeline. 17 AI image platforms, 4 narrative models, HTMX Web UI, MCP Server. Pure Go single binary, zero runtime deps.
```

**Website**: `https://github.com/veawho/via54Design`

**Topics** (20 个 — GitHub 上限):
```
ai, cli, comfyui, golang, mcp, prompt-engineering, stable-diffusion, yaml,
design-system, design-templates, forge, image-generation, model-context-protocol,
narratology, pptx, sdwebui, storyboard, svg, template-engine, video-generation
```

---

## ⚠️ 2. Social Preview

推荐尺寸 1280×640 PNG。上传到 Settings → Social preview。

## ✅ 3. Features — **已完成**

| Feature | 状态 |
|---------|------|
| Issues | ✅ 启用 |
| Discussions | ✅ 启用 (6 个默认分类) |
| Squash merging | ✅ 启用 |
| Auto-delete head branches | ✅ 启用 |
| Required conversation resolution | ✅ 启用 |

## ✅ 4. Issue & PR templates — **已完成**

位于 `.github/ISSUE_TEMPLATE/`:
- `bug_report.md`
- `feature_request.md`
- `template_request.md`
- `question.md`
- `pull_request_template.md`

## ⚠️ 5. Discussions categories — **✅ 全部完成 (8 个)**

**6 个默认 + 2 个新建**:

| Category | Slug | 描述 | 状态 |
|----------|------|------|------|
| 📢 Announcements | `/announcements` | Updates from maintainers | 默认 |
| 💛 Contributors | `/contributors` | How to contribute, code review, design discussions | ✨ **新建** |
| 💬 General | `/general` | Chat about anything and everything here | 默认 |
| 💡 Ideas | `/ideas` | Share ideas for new features | 默认 |
| 📊 Polls | `/polls` | Take a vote from the community | 默认 |
| 🙏 Q&A | `/q-a` | Ask the community for help | 默认 |
| 📚 Resources | `/resources` | Tutorials, blog posts, third-party integrations | ✨ **新建** |
| 🙌 Show and tell | `/show-and-tell` | Show off something you've made | 默认 |

> **GitHub API 限制**: Discussion 分类**只能**通过 UI 添加，GraphQL/REST 都不支持创建端点。
> 实际路径: `https://github.com/<user>/<repo>/discussions` → Categories 旁的 🖊️ 铅笔图标 → New category
> 每个分类有 4 种 Discussion Format: Open-ended / Q&A / Announcement / Poll
> (旧版路径 `/settings/categories` 已废弃，返回 404)

## ✅ 6. Branch protection (`main`) — **已完成**

```json
{
  "required_status_checks": { "strict": true, "contexts": ["build"] },
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true,
    "require_last_push_approval": true
  },
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
```

## ✅ 7. Labels — **已完成 (18 个)**

9 个 GitHub 默认 + 9 个项目专用:
- `bug` `enhancement` `documentation` `good first issue` `help wanted`
- `question` `wontfix` `duplicate` `invalid`
- **项目专用**: `template` `mcp` `web-ui` `go` `security`
- **优先级**: `priority/high` `priority/medium` `priority/low`
- **流程**: `needs-triage`

## 📦 8. Releases — 模板

```markdown
## v0.x.0 — YYYY-MM-DD

### ✨ New
- (feature)

### 🐛 Fixed
- (bug)

### 📚 Docs
- (change)

### 🔧 Internal
- (change)

**Full changelog**: https://github.com/veawho/via54Design/compare/v0.x-1...v0.x.0
```

## ✅ 9. SECURITY.md — **已完成**

`SECURITY.md` 在仓库根目录。90 天协调披露策略。

---

## 📊 自动完成 vs 手动完成 (总结)

| 项 | 状态 | 方法 |
|----|------|------|
| About 描述 | ✅ 自动 | `gh repo edit --description` |
| 20 Topics | ✅ 自动 | `gh api PUT /repos/.../topics` |
| Discussions 启用 | ✅ 自动 | `gh repo edit --enable-discussions` |
| 8 Discussion 分类 (6+2) | ✅ **手动** (UI) | /discussions → 🖊️ → New category |
| 18 Labels | ✅ 自动 | `gh label create` |
| 分支保护 | ✅ 自动 | `gh api PUT /branches/main/protection` |
| Squash merge | ✅ 自动 | `gh repo edit --enable-squash-merge` |
| Auto-delete branches | ✅ 自动 | `gh repo edit --delete-branch-on-merge` |
| **首个 Release** | ⚠️ 手动 | `gh release create` |
| **Social preview** | ⚠️ 手动 | 上传 1280×640 PNG |

## 🔧 复现命令

```bash
# 仓库基础
gh repo edit veawho/via54Design \
  --description "..." \
  --homepage "https://github.com/veawho/via54Design" \
  --enable-discussions --enable-issues \
  --enable-squash-merge --delete-branch-on-merge

# Topics
gh api -X PUT repos/veawho/via54Design/topics \
  -f 'names[]=ai' -f 'names[]=cli' ...  (20 个)

# Labels
for label in "template|fbca04|YAML template addition" "mcp|7b68ee|MCP Server" "web-ui|ffa500|HTMX Web UI" ...; do
  IFS='|' read -r name color desc <<< "$label"
  gh label create "$name" --color "$color" --description "$desc"
done

# 分支保护
gh api -X PUT repos/veawho/via54Design/branches/main/protection \
  --input protection.json
```
