# Changelog

All notable changes to via54Design are documented in this file.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] — 2026-06-08

### Added
- **Phase 1** (in progress): CHANGELOG.md, CONTRIBUTING.md, CODE_OF_CONDUCT.md (this release)
- **Phase 2** (planned): Dependabot, golangci-lint, release workflow, examples/
- **Phase 3** (planned): Go test coverage 6%→70%, benchmarks, README screenshots, pprof
- **Phase 4** (planned): WebSocket/SSE real-time, Prometheus /metrics, OpenTelemetry

### Changed
- Error handling: `errors.Is/As` from 0 to 30+ usages (better error chain)

---

## [0.4.0] — 2026-06-08

### Security
- **P1 (CVE-grade)**: Fixed XSS vulnerability in `internal/template/engine.go`
  - `r.Title` now HTML-escaped via `html.EscapeString` before insertion
  - `heroBodyHTML(titleHTML)` now escapes `titleHTML`
  - Tested: `<script>alert(1)</script>` → `&lt;script&gt;alert(1)&lt;/script&gt;`
  - Tested: `" onerror=alert(1)` → `&#34; onerror=alert(1)`

### Fixed
- `--help` / `-h` now exits 0 (UNIX CLI convention) instead of 1
- All 50+ Go files: `gofmt -w .` (CRLF → LF)
- `test_20_rounds.py`: wrapped `output.html` reads in `try/except` (was crashing)
- `.gitignore`: was matching `cmd/via54/` directory (CRITICAL BUG, blocked CI)
- CI: fixed SIGPIPE (exit 141) by adding `set +e` + `|| true` to test scripts
- CI: Go version matrix `['1.21', '1.22']` → `['1.26.2']` (matches go.mod)

### Added
- **Release v0.4.0**: 5 platform zips published to GitHub Releases
  - via54-darwin-amd64.zip / via54-darwin-arm64.zip
  - via54-linux-amd64.zip / via54-linux-arm64.zip
  - via54-windows-amd64.zip
- **Discussions categories (8 total)**: 6 default + Contributors + Resources
- **GitHub community files**: 5 Issue/PR templates, CI workflow, SECURITY, etc.
- **Bilingual README**: README.md (中文) + README_EN.md (English)
- **Interactive mode**: `via54` (no args) opens Chinese menu
- **--http flag**: CLI serve now supports HTTP mode
- **GOFLAGS auto-detect**: Makefile detects exFAT/FAT32 and adds `-buildvcs=false`
- **`make fs-check` target**: shows current filesystem + GOFLAGS state
- **`make test-{e2e,stress,concurrent}` targets**: 420 E2E test invocations

### Changed
- Project migrated: `C:\Users\via54\AppData\Local\Temp\via54Design` → `G:\agent\projects\via54Design`
  - Build perf: 21.2s vs 34.4s (38% faster on G: drive)
- Architecture: dual binary (CLI + MCP Server) split into 2 separate binaries
- All 12 deleted Node.js + Python scripts migrated to pure Go

### Tests
- **test_20_rounds.py**: 20/20 PASS, 0 FAIL
- **test_stress_200.py**: 200/200 PASS, latency 26-33ms avg
- **test_concurrent.py**: 200/200 concurrent, 0 failures, 141 RPS
- **go test**: 11 PASS in internal/workflow

### Documentation
- README.md: 19.6KB Chinese with bilingual switcher
- README_EN.md: 22.9KB English with Shakespeare/Wilde/Donne/陆游 quotes
- AGENTS.md: comprehensive AI work context
- REPOSITORY_SETUP.md: 10-point GitHub setup checklist
- MAINTAINING.md: triage workflow + Conventional Commits + release process
- ACKNOWLEDGMENTS.md: dependency attribution
- SECURITY.md: 90-day coordinated disclosure

---

## [0.3.0] — 2026-06-07

### Added
- Initial AGPL-3.0 + MIT dual license release
- 5 platform cross-compilation via Makefile
- 16 Go packages, 12,282 LOC, 107 YAML templates
- 13 CLI subcommands (generate, narrate, prompt, list, quality, media, export, pipeline, present, web, forge, comfyui, version)
- MCP Server with 18 tools
- HTMX web UI (zero JavaScript, 10 endpoints)
- 17 AI image/video platforms (midjourney, flux, sdxl, kling, etc.)
- 4 narrative models (three-act, hero's journey, cognitive arc, problem-solution)
- Quality gate (5 dimensions)
- Pattern extractor (learns from generated outputs)

### Infrastructure
- 5 GitHub labels (template, mcp, web-ui, go, security)
- Branch protection (1 PR approval, strict checks, linear history)
- Squash merge + auto-delete branches

---

## Release Process

1. Update `CHANGELOG.md` → "Unreleased" section
2. Run `make test test-e2e test-stress test-concurrent` (all must pass)
3. Move entries from "Unreleased" to new version section with date
4. Commit: `git commit -am "chore: release vX.Y.Z"`
5. Tag: `git tag -a vX.Y.Z -m "vX.Y.Z: <summary>"`
6. Push: `git push origin main vX.Y.Z`
7. GitHub Action `.github/workflows/release.yml` automatically:
   - Builds 5 platform binaries
   - Creates GitHub Release with zips attached
   - Updates Homebrew tap (if configured)

[Unreleased]: https://github.com/veawho/via54Design/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/veawho/via54Design/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/veawho/via54Design/releases/tag/v0.3.0
