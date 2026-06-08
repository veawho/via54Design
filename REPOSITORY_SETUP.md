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

## ⚠️ 5. Discussions categories — **API 限制，需手动**

**已存在 6 个默认分类** (无法通过 API 创建更多):

| Category | Slug | 用途 |
|----------|------|------|
| Announcements | `/announcements` | Releases, breaking changes |
| General | `/general` | 通用讨论 |
| Ideas | `/ideas` | Feature brainstorming |
| Polls | `/polls` | 投票 |
| Q&A | `/q-a` | Usage questions |
| Show and tell | `/show-and-tell` | Templates, generated outputs |

> **GitHub GraphQL 与 REST API 都不支持** 创建 Discussion 分类。
> 访问 https://github.com/veawho/via54Design/settings/categories 手动添加。

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
| 18 Labels | ✅ 自动 | `gh label create` |
| 分支保护 | ✅ 自动 | `gh api PUT /branches/main/protection` |
| Squash merge | ✅ 自动 | `gh repo edit --enable-squash-merge` |
| Auto-delete branches | ✅ 自动 | `gh repo edit --delete-branch-on-merge` |
| **Discussion 分类** | ⚠️ 手动 | GitHub UI (API 限制) |
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
