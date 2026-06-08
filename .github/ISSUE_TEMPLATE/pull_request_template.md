---
name: Pull Request
about: Submit a change to via54Design
title: '[type] description'
---

## Type of change

- [ ] 🐛 Bug fix (non-breaking change that fixes an issue)
- [ ] ✨ New feature (non-breaking change that adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to change)
- [ ] 📚 Documentation only
- [ ] 🧪 Test only
- [ ] 🔧 Refactor / chore (no functional change)
- [ ] 🎨 New template (YAML only, no Go change)
- [ ] 🌐 Translation / i18n

## What does this change?

A concise 1–3 sentence summary.

## Why?

The motivation. Reference any related issue (`Closes #123`) or design doc.

## How is it tested?

- [ ] Manual test: I ran `...` and got `...`
- [ ] New unit test added (`go test ./...`)
- [ ] Updated `test_20_rounds.py`
- [ ] Cross-platform check: tested on Windows / macOS / Linux

## Checklist

- [ ] My code follows the project's Go style guide (`gofmt`, `go vet ./...` clean)
- [ ] I have added an SPDX-License-Identifier header to any new `.go` / `.rs` files
- [ ] I have updated the relevant documentation (README, AGENTS.md, or `docs/`)
- [ ] I have added or updated tests where applicable
- [ ] For YAML template changes: I have validated the template via `via54 list` / `via54 generate` / `via54 narrate`
- [ ] For dependency changes: I have run `go mod tidy` and verified the lockfile is committed
- [ ] My commit messages follow `type: description` format

## Screenshots / output

If the change is visible (Web UI, exported PPTX, etc.), attach before/after.
