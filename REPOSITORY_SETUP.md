# Repository settings — manual configuration

> These files document the GitHub repository settings that cannot be set via git.
> Apply them in the GitHub web UI: https://github.com/veawho/via54Design/settings

## 1. About (top-right of repo page)

**Description** (160 chars max):
```
Design template engine + narrative engine + media pipeline. 17 AI image platforms, 4 narrative models, HTMX Web UI, MCP Server. Pure Go single binary, zero runtime deps.
```

**Website**:
```
https://github.com/veawho/via54Design
```

**Topics** (max 20, comma-separated):
```
ai, art, cli, comfyui, design-system, design-templates, golang, html,
image-generation, mcp, model-context-protocol, narratology, pptx,
prompt-engineering, sdwebui, stable-diffusion, storyboard, svg,
template-engine, video-generation, yaml
```

---

## 2. Social Preview

Recommended: 1280×640 PNG showing the via54Design logo + a 3-panel preview
(CLI menu, generated HTML, exported PPTX).

A draft SVG is available at `docs/social-preview.svg` — export to PNG using
any modern browser or `rsvg-convert`.

---

## 3. Features (checkbox in About sidebar)

Enable these in the repository settings:

- [x] ✅ Issues
- [x] ✅ Pull requests
- [x] ✅ Discussions
- [x] ✅ Projects
- [x] ✅ Wiki (optional)
- [x] ✅ Sponsorship (if applicable)
- [x] ✅ Preserve this repository
- [x] ✅ Automatically delete head branches
- [x] ✅ Allow squash merging
- [x] ✅ Allow merge commits
- [ ] Allow rebase merging (off, prefer linear history)
- [ ] Allow auto-merge

---

## 4. Issue & PR templates

Located in `.github/ISSUE_TEMPLATE/`:

- `bug_report.md` — bug reports
- `feature_request.md` — new CLI/API/MCP features
- `template_request.md` — new YAML template (color/font/layout/narrative/prompt)
- `question.md` — usage questions

PR template: `.github/ISSUE_TEMPLATE/pull_request_template.md`

---

## 5. Discussions categories

Recommended categories (enable in Settings → Features → Discussions → Set up discussions):

| Category | Purpose |
|----------|---------|
| 📣 **Announcements** | Releases, breaking changes (admin-only) |
| 💡 **Ideas** | Feature brainstorming before opening an Issue |
| 🙏 **Q&A** | Usage questions, troubleshooting |
| 🎨 **Show and tell** | Templates, generated outputs, workflows |
| 🤝 **Contributors** | How to contribute, design discussions |
| 📚 **Resources** | Tutorials, blog posts, integrations |

---

## 6. Branch protection (`main`)

Recommended rules:

- [x] Require pull request reviews before merging (1 approval)
- [x] Require status checks to pass before merging
  - Select: `Build & Test` (configured in CI)
- [x] Require branches to be up to date before merging
- [x] Require linear history (squash or rebase)
- [x] Do not allow force pushes
- [x] Do not allow deletions

---

## 7. Labels

Default labels to ensure exist:

| Label | Color | Description |
|-------|-------|-------------|
| `bug` | `#d73a4a` | Something isn't working |
| `enhancement` | `#a2eeef` | New feature or request |
| `documentation` | `#0075ca` | Improvements or additions to docs |
| `good first issue` | `#7057ff` | Good for newcomers |
| `help wanted` | `#008672` | Extra attention is needed |
| `question` | `#d876e3` | Further information is requested |
| `template` | `#fbca04` | YAML template addition |
| `wontfix` | `#ffffff` | This will not be worked on |
| `duplicate` | `#cfd3d7` | This issue or PR already exists |
| `needs-triage` | `#ededed` | Awaiting maintainer review |
| `mcp` | `#7b68ee` | Related to MCP Server |
| `web-ui` | `#ffa500` | Related to HTMX Web UI |
| `go` | `#00add8` | Go language change |
| `security` | `#ee0701` | Security issue |

---

## 8. Releases

Recommended release notes format (see `RELEASING.md`):

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

---

## 9. SECURITY.md

Create `.github/SECURITY.md` (or `SECURITY.md` at repo root) for vulnerability
disclosure policy. Use the GitHub-recommended structure:

```markdown
# Security Policy

## Supported Versions
| Version | Supported |
|---------|-----------|
| v0.4.x  | ✅ |
| v0.3.x  | ✅ |
| < v0.3  | ❌ |

## Reporting a Vulnerability
Please report security issues to **via54@users.noreply.github.com** (or open a
private security advisory on GitHub). Do not file a public Issue.
```

---

## 10. FUNDING.yml

If accepting sponsorships, create `.github/FUNDING.yml`:

```yaml
github: [veawho]
```

---
