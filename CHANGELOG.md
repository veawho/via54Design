# Changelog

All notable changes to via54Design are documented in this file.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [0.8.0] - 2026-06-29

- Add integrations/ for OpenMontage (27K), diffusers (33K), Pixelle-Video (23K), KrillinAI (10K). Plan v6 pipeline release.
## [0.6.1 + v6 pipeline] - 2026-06-29

- Plan: v6 video pipeline release v1.0 (move _scripts/v6_*.py to src/v6/)
- Plan: Integrate diffusers API for AI video generation
- Add architecture-comparison docs (OpenMontage 12 pipelines, 52 tools, 500+ skills)

## [Unreleased]

### Added
- **Phase 2** (planned): Dependabot, golangci-lint, release workflow, examples/
- **Phase 3** (planned): Go test coverage 6%→70%, benchmarks, README screenshots, pprof
- **Phase 4** (planned): WebSocket/SSE real-time, Prometheus /metrics, OpenTelemetry

### Changed
- Error handling: `errors.Is/As` from 0 to 30+ usages (better error chain)

---

## [0.6.1] — 2026-06-12

### Added
- **`docs/prompt-mastery-v3.md`** (9.2KB): 中文 AI 生图提示词精准控制 v3 — 8 维度黄金技巧
  - mmx vs SD/MJ 根本差异 (LLM 驱动 ≠ 权重语法)
  - mmx 官方 example 拆解 + 5 微洞察
  - 8 层结构: Meta/Subject/Pose/Scene/Lighting/Camera/Style/Negative
  - 摄影词典 10 顶级光照 + 5 胶片 + 5 焦段 (来自 MJ Reference 12291⭐)
  - 铁律 8 条
- **`templates/prompts/minimax.yaml` v3** (10.8KB): 8 层 + 摄影词典 + 数字标注 `:25` + mmx quirks
  - 新 `mmx_quirks` 字段 (prompt_optimizer + LLM 特性)
  - 新 `photography_lexicon` 摄影词典速查
  - 新 `pose` 独立层
  - `lighting` 加 `color_temperature` 显式色温
  - `camera` 加 `bokeh_shape` 显式光斑形状
  - `negative` 升级为 mmx 友好 "Avoid:" 语法
- **`docs/prompt-v3-experiment-report.md`**: 3 对照实验报告 (v1 中文 vs v2 7 层 vs v3 8 层)

### Verified (本地同 seed=42 · 16:9 · 1n)
- `go vet ./...` — 0 issues
- `go test ./...` — 8/8 packages PASS (export, media, narrate, prompt, quality, template, util, workflow)
- `go build ./...` — 0 error
- 3 张实验图: `minimax-output/v3/cat_a_001.jpg` (207KB) / `cat_b_001.jpg` (280KB) / `cat_c_001.jpg` (264KB)
- vision 评分: A=8.4 / B=7.75 / C=8.0 (mmx image-01 上限 8.0-8.4)

### 核心发现 (★ 沉淀到 SKILL)
1. **mmx image-01 上限 8.0-8.4** (普通生图任务同类主体)
2. **要 9.0+ 必须换模型** (FLUX / MidJourney / GPT-Image-2)
3. **蓝眼虎斑 + 金属书架** mmx 必失败 (训练集偏差, 必给绿眼/木质)
4. **摄影词典 = 画质 9.0 拉满** (Rembrandt + Kodak Portra + 85mm f/1.4)
5. **prompt 长度 sweet spot 200-400 字符** (中文 ≤100字, 英文 ≤80词)
6. **同 seed 字节级复现**: `--seed 42` 验证 v3 = C 图

### 引用源
- HF Diffusers 官方权重语法
- mmx OpenAPI (image-01, prompt_optimizer: true)
- MJ Reference 12291⭐ (willwulfken)
- SD 负向词 91⭐ (mikhail-bot)
- via54Design v2 实测 (commit 09e8dcc, 8.6/10)

---

## [0.6.0] — 2026-06-11

### Added
- **§12 SVG v2 规范** (`docs/12-svg-spec.md`): viewBox=680/382, class t/ts/th, 12/14/24px, stroke-width 1.5
  - 11/11 SVG 模板自检全 PASS
- **Phase D 可观测性** (`cmd/via54/observability.go`): Prometheus /metrics + net/http/pprof
  - 9 metric (请求/错误/延迟/Go runtime, stdlib, 无 client_golang 依赖)
  - 6 pprof 端点 (/debug/pprof/{profile,heap,goroutine,threadcreate,block,cmdline,trace,symbol})
  - 包外层 mux: mcp-go SSEServer 绑死, http.ServeMux 接管 /metrics + /debug/pprof
- **`docs/09-observability.md`** (8.3KB): §9 可观测性规范 (7 端点 + 9 metric)
- **§2 viewport** (`docs/02-viewport-spec.md`): `vpMeta` 类

---

## [0.5.1] — 2026-06-09

### Fixed
- **release.yml (P0)**: `via64-*` typo in 5 file globs → `via54-*` (would have silently dropped all assets from v0.5.0+ releases)
- **FindBaseDir (v0.5.1 core)**: `internal/util/paths.go` now reads `VIA54_BASE_DIR` env var as highest-priority override, so `brew install /usr/local/bin` scenarios can locate `templates/` (commit `4ca79dd`)
- **PR #4 closed** as duplicate: fix was cherry-picked to main before PR review; closing with explanation comment

### Verified (local Windows · Go 1.26.2 · NTFS)
- `go vet ./...` — 0 issues
- `go test ./...` — 4/4 packages PASS (export, template, util, workflow)
- `python test_20_rounds.py` — 20/20 PASS, 0 WARN, 0 FAIL
- `go build ./cmd/via54/` — 18 MB binary
- `go build ./cmd/mcp-server/` — 16 MB binary
- `./via54.exe version` + `list` + `prompt` + `narrate` smoke test — all green

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
