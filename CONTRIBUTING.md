# Contributing to via54Design

> "I would have written a shorter letter, but I did not have the time."
> — Blaise Pascal (often cited in software engineering)

Thank you for your interest in contributing to via54Design! This project
thrives on community contributions — whether it's a bug fix, a new template,
a translation, or a documentation improvement.

---

## 🚀 Quick Start

### 1. Fork & Clone
```bash
git clone https://github.com/<your-fork>/via54Design
cd via54Design
```

### 2. Set Up Development Environment
```bash
# Go 1.26.2 or later
go version

# Optional but recommended
make fs-check         # Check your filesystem (exFAT/FAT32 needs GOFLAGS=-buildvcs=false)
go build -o via54.exe ./cmd/via54/
go build -o via54-mcp.exe ./cmd/mcp-server/
```

### 3. Verify Everything Works
```bash
make test             # 11 Go unit tests
make test-e2e         # 20 end-to-end rounds
make test-stress      # 200 pressure rounds
make test-concurrent  # 200 concurrent requests
```

### 4. Make Your Change
```bash
git checkout -b feat/your-feature-name
# ... edit code ...
```

### 5. Verify Before PR
```bash
go vet ./...
gofmt -l .            # Should be empty
go test ./... -count=1
```

---

## 📝 Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/) with
emoji for visual scanning:

```
<emoji> <type>(<scope>): <subject>

<body>

<footer>
```

### Types
| Type       | Description                                      | Example                                 |
|-----------|--------------------------------------------------|----------------------------------------|
| `feat`    | New feature                                       | `feat(prompt): add Midjourney v7 support` |
| `fix`     | Bug fix                                           | `fix(export): prevent XSS in template` |
| `refactor`| Code change that neither fixes bug nor adds feature | `refactor(template): extract color resolver` |
| `docs`    | Documentation only                                | `docs(readme): add Go 1.26 install steps` |
| `test`    | Adding or correcting tests                        | `test(template): add benchmark` |
| `chore`   | Build, CI, deps, tooling                          | `chore(ci): pin Go to 1.26.2` |
| `perf`    | Performance improvement                           | `perf(export): cache YAML parses` |

### Emojis (optional but encouraged)
🎉 `:tada:` initial commit • ✨ `:sparkles:` feat • 🐛 `:bug:` fix • 📝 `:memo:` docs
🔥 `:fire:` remove code/file • 💚 `:green_heart:` CI fix • 🚀 `:rocket:` deploy

### Example
```bash
git commit -m "✨ feat(prompt): add Ideogram v2 platform support

- Add templates/prompts/ideogram.yaml with 6-dim control
- Update cmd/via54/prompt_cmd.go listPlatforms()
- Add tests in internal/prompt/ideogram_test.go
- Update docs/prompts/ideogram.md

Verified:
  ✓ go test ./internal/prompt/...  (3/3 pass)
  ✓ via54 prompt --scene '...' --platform ideogram  (correct output)
  ✓ gofmt clean
"
```

---

## 🧪 Testing Requirements

Before opening a PR, ensure:

| Check | Command | Pass Criteria |
|-------|---------|---------------|
| `go vet` | `go vet ./...` | 0 errors |
| `gofmt` | `gofmt -l .` | empty output |
| `go test` | `go test ./... -count=1` | all pass |
| `go build` | `go build ./...` | 0 errors |
| E2E | `make test-e2e` | 20/20 pass |
| Stress | `make test-stress` | 200/200 pass |

If your change touches a public API, add or update `internal/<pkg>/*_test.go`.
We aim for **70%+ test coverage** in new code.

---

## 🎨 Adding a New Template

Templates are YAML files in `templates/`. To add a layout:

### 1. Create the YAML file
`templates/layouts/my-new-layout.yaml`:
```yaml
id: my-new-layout
name: 我的新布局
description: 16:9 比例, 适合产品发布
category: hero

structure:
  type: split
  columns: 2
  ratio: 0.6   # image : text
  
text:
  position: right
  max_chars: 50
  alignment: left

styles:
  spacing: 1.618   # Golden ratio
  typography: editorial
```

### 2. Test it
```bash
go build -o via54.exe ./cmd/via54/
./via54.exe list | grep my-new-layout
./via54.exe generate --layout my-new-layout --color ink-wash --font ming-hei-editorial --title "测试"
```

### 3. Update registry (if needed)
`templates/registry.yaml`:
```yaml
layouts:
  - id: my-new-layout
    category: hero
    tags: [产品, 发布]
```

---

## 🌐 Adding a New AI Platform

1. Create `templates/prompts/<platform>.yaml` (copy existing as template)
2. Update `cmd/via54/prompt_cmd.go` `listPlatforms()` slice
3. Update `internal/prompt/generator.go` if format differs
4. Add test in `internal/prompt/<platform>_test.go`
5. Update `docs/prompts/<platform>.md`

See `templates/prompts/midjourney.yaml` as the canonical reference.

---

## 🔒 Security

**DO NOT** open a public issue for security vulnerabilities.
Email `security@via54design.local` (or contact via GitHub Security tab).
We follow a 90-day coordinated disclosure policy. See `SECURITY.md`.

---

## 📜 License

By contributing, you agree that your contributions will be licensed under:
- **AGPL-3.0** for source code (`.go`, `.yaml`, `.ts`, etc.)
- **MIT** for documentation and templates (`docs/`, `templates/`)

This is the project's dual license (see `LICENSE` and `NOTICE`).

---

## 🌏 Communication

- **GitHub Issues**: bugs, features, questions
- **GitHub Discussions**: ideas, show-and-tell, Q&A, contributors, resources
- **Pull Requests**: code changes (1 approval required)
- **Security**: private email (see SECURITY.md)

---

## 🏆 Recognition

Contributors are recognized in `ACKNOWLEDGMENTS.md`. Significant contributors
(>10 merged PRs or 1 major feature) are listed in `README.md` Hall of Fame.

---

## ❓ Questions?

- Read `AGENTS.md` (comprehensive AI work context)
- Read `ARCHITECTURE.md` (technical architecture)
- Read `docs/` directory
- Ask in GitHub Discussions (Q&A category)

We appreciate your time and effort! 🙏
