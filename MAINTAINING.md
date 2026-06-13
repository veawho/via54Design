# Maintainer Guide

> So you want to be a `via54Design` maintainer? Welcome — read this first.

## Roles

| Role | Responsibilities | Rights |
|------|------------------|--------|
| **Owner** (`veawho`) | Direction, releases, security | All |
| **Maintainer** | Triage issues, review PRs, merge | Triage label, merge |
| **Contributor** | File issues, send PRs | Comment, file issues |

## Triage workflow

1. **New issue** arrives → bot applies `needs-triage`
2. **Maintainer** reviews within 7 days, applies one of:
   - `bug` + `priority/high|medium|low` → actionable
   - `enhancement` + `priority/medium|low` → roadmap candidate
   - `question` → close once answered
   - `duplicate` → link to original, close
   - `wontfix` → explain why, close
3. **High-priority bugs** go in the next patch release
4. **Enhancements** discussed in `Discussions → Ideas` first

## Pull request flow

1. PR opened → CI runs `build.yml` on 3 OSes × 2 Go versions = **6 jobs**
2. Maintainer reviews (code style, tests, docs, license headers)
3. Approval → squash-merge with conventional commit message
4. Author gets a thank-you comment

### Conventional commits

We use the [Conventional Commits](https://www.conventionalcommits.org/) spec:

| Type | When | Example |
|------|------|---------|
| `feat:` | New user-facing feature | `feat: add ComfyUI workflow engine` |
| `fix:` | Bug fix | `fix: prompt platform validation rejects typos` |
| `refactor:` | No-behavior-change cleanup | `refactor: extract baseDir() into internal/util` |
| `docs:` | README / docs only | `docs: add English README_EN.md` |
| `test:` | Tests only | `test: add 20-round integration suite` |
| `chore:` | Tooling, CI, deps | `chore: bump mcp-go to v0.54.1` |
| `license:` | License header changes | `license: add SPDX headers to all .go files` |

Squash-merge preserves the PR title → becomes the commit message.

## Release process

1. Bump version in `Makefile` (`VERSION` variable)
2. Run `make release` locally — produces 5 platform zips in `dist/`
3. Push a `v0.x.0` tag → CI runs `release.yml` workflow
4. GitHub Actions attaches the zips to the release
5. Manually write the release notes (see `REPOSITORY_SETUP.md` for template)
6. Announce in `Discussions → Announcements`

## When to break backward compatibility

`via54Design` is pre-1.0, so breaking changes are allowed in any minor release.
After v1.0, the following require a major version bump:

- Removing a CLI subcommand
- Changing a YAML template field name
- Changing the MCP tool schema
- Changing the HTTP API contract (path, params, response shape)

Adding new optional fields / subcommands / endpoints is always backward-compatible.

## Code review checklist

- [ ] SPDX header on every new `.go` / `.rs` file
- [ ] `gofmt -l .` produces no output
- [ ] `go vet ./...` produces no warnings
- [ ] No new external dependencies without prior discussion
- [ ] Tests added or updated (`test_20_rounds.py` if user-facing)
- [ ] Docs updated (README, AGENTS.md, `docs/`)
- [ ] Commit message follows conventional format
- [ ] PR description explains the *why*, not just the *what*
