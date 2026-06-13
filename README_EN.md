<div align="center">

# via54Design

**A design template engine + narrative engine + media pipeline — where human inspiration meets AI structure.**

*Hard work, well done, leaves no trace. Right tools, well chosen, leave no friction.*

<br/>

[![License: AGPL-3.0 / MIT](https://img.shields.io/badge/License-AGPL--3.0%20%2F%20MIT-blue.svg)](./LICENSE)
[![Go 1.21+](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Platforms](https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](#-platform-support)
[![Tests](https://img.shields.io/badge/tests-20%2F20%20PASS-brightgreen)](./test_20_rounds.py)
[![MCP](https://img.shields.io/badge/MCP-compatible-purple)](https://modelcontextprotocol.io)

[**English**](#english-readme) · [**中文**](./README.md)

</div>

---

## 🌟 The Philosophy

> *"A horse must learn to run before it knows it can gallop; an idea must spark before it can blaze."*
> — adapted from a Chinese proverb

Human imagination is **scattered, flickering, untamed**. AI is **structured, deterministic, controllable**. We do not pretend to fuse them — we build a bridge.

`via54Design` is not a design tool. It is a **design template engine** that captures the instincts of a 4A creative director in YAML files, then executes them deterministically via a Go core. One line of human inspiration becomes a 90-second brand story. A scattered mood board becomes a 30-slide presentation. A photograph becomes a structured prompt for seventeen AI image platforms.

> *"All the world's a stage, and all the men and women merely players."*
> — William Shakespeare, *As You Like It*

Every story we help you tell is a stage. Every scaffold we generate is a rehearsal. The final performance is yours.

> *"Imitation is the sincerest form of flattery that mediocrity can pay to greatness."*
> — Oscar Wilde

We imitate the great narrative arcs — three-act, hero's journey, cognitive-arc, problem-solution — not because we are mediocre, but because these patterns have survived the test of centuries. We honor them by encoding them as deterministic, executable templates.

---

## 🧠 Turn Your AI Companion Into Your Best Storyteller

**You don't want another tool that aggregates skills. You don't want another workflow optimizer.**

`via54Design` doesn't replace your creativity — **it gives you a brilliant hand that catches inspiration the moment it lands**.

> *"The best ideas are not pulled from the mind — they fall into it like rain on still water."*

Your half-formed idea, in one sentence, becomes a complete narrative scaffold, a slide deck, a storyboard, a set of camera-ready AI prompts. The AI does the structural work. You keep the spark.

> *"文章本天成，妙手偶得之"* (Literary works are made by heaven; the gifted hand but occasionally finds them.)
> — Lu You, Song dynasty poet

A work of literature is pre-ordained; the masterful hand only occasionally discovers it. We are the masterful hand. You supply the literature.

---

## Part I: Story → Video

#### From a single sentence to a 90-second brand story

**Step 1 — The human writes the opening line (an inspiration only you can have)**

> "In 1920s Paris, a Chinese tailor opened a small shop.
> His qipao dresses wove Art Deco's geometric lines into silk.
> No one imagined this garment would change the fashion of two civilizations."

This sentence contains: a character (the tailor), an era (the 1920s), a place (Paris), a conflict (East vs. West), a hook (changing fashion). The AI cannot conjure this seed — it must come from you.

**Step 2 — The AI expands it into a narrative scaffold**

```bash
via54 narrate --seed "In 1920s Paris, a Chinese tailor opened a small shop..." \
  --model heros-journey --duration 90 --format json --output scaffold.json
```

The AI analyzes your seed, matches the best narrative model, and produces a structured scaffold:

```
📋 Hero's Journey  90s
├── Act I · The Ordinary  (0-22s)  mood: calm      narration: Every day, we...
├── Act II · The Meeting  (22-44s) mood: curious   narration: Until one day...
├── Act III · The Ordeal (44-66s) mood: excited   narration: Something changed...
└── Act IV · The Return  (66-90s) mood: inspiring  narration: And every day since...

📋 Shot list: 12 shots (WIDE / MEDIUM / CLOSE-UP / DETAIL cycle)
📋 Fountain script skeleton (4 acts, 8 scenes)
🎞️ LLM-complete prompts (drop-in to Claude / GPT for full script)
```

**Step 3 — You choose the narrative model, confirm the direction**

| Model | Beats | Best for |
|-------|-------|----------|
| `three-act` | Question → Answer → Call | Product launches, brand ads |
| `heros-journey` | Ordinary → Meeting → Transformation → Return | Brand stories, documentaries |
| `cognitive-arc` | Hook → Foundation → Core → Case → Extension → Summary | Explainers, tutorials |
| `problem-solution` | Pain → Solution → Proof → Action | Sales videos, demos |

**Step 4 — Generate the output**

```bash
# One scaffold, many output formats
via54 export pptx scaffold.json --output story.pptx
via54 export markdown scaffold.json --output slides.md
via54 export svg scaffold.json --output ./scenes
via54 export json scaffold.json --output scenes.json
```

---

## Part II: Story → Presentation

The same scaffold exports to multiple presentation formats. No rework.

**Narrative seed → PPTX deck**

```bash
via54 narrate --seed "In 1920s Paris, a Chinese tailor..." \
  --model heros-journey --duration 90 --format json --output scaffold.json

via54 export pptx scaffold.json --output story.pptx \
  --style editorial --theme templates/color-schemes/ink-wash.yaml
```

- 100% pure Go — no Node.js, no unioffice
- 4 layout styles (minimal / editorial / bold / accent-bar)
- 31 color themes, mood-mapped to accent colors
- 16:9 widescreen, text directly editable in PowerPoint

**Narrative seed → HTML design**

```bash
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial \
  --title "The Tailor's Story" --output story.html --presentation
```

---

## Part III: Story → Creative Imagery

**Narrative scaffold → AI image prompts**

```bash
via54 prompt --scene "1920s, a small tailor shop on Paris's Left Bank, a Chinese tailor crafting a qipao" \
  --platform midjourney --output prompt.md
```

Structured prompts across **17 platforms**, **26 control dimensions** (subject / environment / lighting / style / mood / composition / camera / color / texture, etc.).

**Narrative scaffold → SVG vector scenes**

```bash
via54 export svg scaffold.json --output ./scenes
```

One independent SVG per act, 16:9 viewBox, infinitely scalable.

**Narrative scaffold → Prompt pipeline**

```bash
via54 narrate --seed "..." --format json | via54 prompt --from-scaffold /dev/stdin
```

---

### Capability Comparison

| Capability | Input | Output | Engine | Runtime deps |
|------------|-------|--------|--------|--------------|
| 🎬 Story → Video | sentence → narrative JSON | Video script / ComfyUI workflow | narrate + storyboard2video | Optional Forge/ComfyUI |
| 📊 Story → Presentation | sentence → narrative JSON | PPTX / Markdown / PDF / SVG | **pure Go** | **none** |
| 🎨 Story → Imagery | sentence → narrative JSON | Structured prompts (17 platforms) / SVG | **pure Go** | **none** |

---

## 🖥️ Platform Support

| Platform | Arch | Binary | Web UI | PDF/Video export |
|----------|------|--------|:------:|:----------------:|
| **macOS** Intel | amd64 | `via54-darwin-amd64` | ✅ native | needs Node.js + ffmpeg |
| **macOS** Apple Silicon | arm64 | `via54-darwin-arm64` | ✅ native | needs Node.js + ffmpeg |
| **Linux** | amd64 | `via54-linux-amd64` | ✅ native | needs Node.js + ffmpeg |
| **Linux** ARM (RPi/Graviton) | arm64 | `via54-linux-arm64` | ✅ native | needs Node.js + ffmpeg |
| **Windows** | amd64 | `via54-windows-amd64.exe` | ✅ native | needs Node.js + ffmpeg |

> **Core features (prompt / narrative / design / export PPTX/JSON/SVG) are 100% pure Go, zero external dependencies.**
> PDF export and video rendering need Node.js + Playwright + ffmpeg — these are optional; the Web UI does not.

**Quick install:**

```bash
# macOS (Intel or Apple Silicon)
bash hack/setup_macos.sh

# Linux (auto-detects apt / dnf / pacman)
bash hack/setup_linux.sh

# Windows (download binary + install Node.js)
# Download: https://github.com/veawho/via54Design/releases
# Node.js:  https://nodejs.org
```

**Build from source (Go 1.21+ required):**

```bash
make build          # CLI for current platform
make build-mcp      # MCP Server
make cross          # cross-compile 5 platforms
make release        # bundle into zip per platform
```

**Run the test suite:**

```bash
python test_20_rounds.py   # 20-round stress / smoke / stability / usability / accuracy
```

---

## 🧪 Test Status

| Dimension | Status | Evidence |
|-----------|--------|----------|
| **Code simplicity** | ✅ solid | 12,254 lines Go / 58 files, 0 panics, 179 explicit error checks |
| **Architecture** | ✅ solid | Two binaries + 12 internal packages, no cyclic deps |
| **Feature coverage** | ✅ solid | 16 CLI subcommands + 26 API endpoints + 18 MCP tools + 17 platforms |
| **Code quality** | ✅ solid | SPDX 100% (58/58), 0 FIXME, gofmt 100%, `go vet` 0 warnings |
| **Build verification** | ✅ 20/20 | All 20 rounds pass (see `test_20_rounds.py`) |
| **Output determinism** | ✅ pass | Same input → same md5 across 3 runs (prompt / narrate / generate) |
| **Edge cases** | ✅ pass | Empty input / invalid platform / 600-char string all handled gracefully |
| **Stress** | ✅ pass | 20× prompt @ 23ms each, 10× narrate @ 20ms, 5× generate @ 20ms |
| **Integration** | ✅ pass | Web UI HTTP 200, HTMX endpoints return correct fragments, 0 inline JS |

---

## 📜 License

- Go source: `AGPL-3.0-only`
- Templates / scripts / docs: `MIT`
- See [`LICENSE`](./LICENSE) and [`ACKNOWLEDGMENTS`](./ACKNOWLEDGMENTS)

---

## 🚀 Quick Start

```bash
# Download the right binary for your platform
# macOS / Linux:
curl -L https://github.com/veawho/via54Design/releases/latest/download/via54-$(uname -s | tr A-Z a-z)-$(uname -m | sed s/x86_64/amd64/ | sed s/aarch64/arm64/).zip -o via54.zip

# List all commands (or just type `via54` for the interactive menu)
via54

# Generate a prompt
via54 prompt --scene "A cat on a moonlit rooftop" --platform midjourney

# Generate an HTML design
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial --title "Title"

# Generate a narrative scaffold
via54 narrate --seed "A tailor in Paris changed fashion" --model three-act --duration 60

# Launch the Web UI
via54 web --port 8080
```

---

## 📋 CLI Reference

### 🎨 Prompt Engineering — sentence → structured prompt

```bash
via54 prompt --scene "scene description" --platform flux
via54 prompt --list                                       # list all 17 platforms
```

Image-to-prompt (planned, not yet implemented):
- TODO: local VLM image analysis
- Current workaround: use the `--ref` flag to attach a reference image path to the prompt output

| Platform | Format | Distinguishing features |
|----------|--------|-------------------------|
| `flux` / `midjourney` / `dalle3` / `sd3` / `stable_diffusion` | English params | 26 control dimensions |
| `ideogram` / `recraft` / `seedance` | English params | Style / text control |
| `gemini` / `veo` / `sora` / `kling` / `pika` | Natural language | Video generation |
| `comfyui` / `forge` | JSON | API payload |

### 📖 Narrative Engine — sentence → full story

```bash
via54 narrate --seed "one-sentence seed" --model three-act --duration 60
via54 narrate --seed "..." --format json --output scaffold.json
via54 narrate --list                                       # list narrative models
```

| Model | Beats | Best for |
|-------|-------|----------|
| `three-act` | Question → Answer → Call | Product launch, brand ad |
| `heros-journey` | Ordinary → Meeting → Transformation → Return | Brand story, documentary |
| `cognitive-arc` | Hook → Foundation → Core → Case → Extension → Summary | Explainer, tutorial |
| `problem-solution` | Pain → Solution → Proof → Action | Sales video, demo |

Output includes: story outline, beat timeline, Fountain script skeleton, shot list, and LLM-complete prompts.

### 🎨 Design Generation — layout × color × font → HTML

```bash
via54 generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial \
  --title "Title" --output demo.html
via54 generate --presentation ...            # presentation mode (16:9 lock)
```

> **Note**: the real layout IDs are `hero-split-16-9`, `bento-grid-2x2`, `gallery-waterfall` — not `hero`. The CLI output is written to `output.html` (a file, not stdout).

### 📄 Export Pipeline — pure Go, zero external deps

```bash
via54 export pptx scaffold.json              # PPTX (editable)
via54 export pdf story.html                  # PDF
via54 export svg scaffold.json               # SVG scenes
via54 export markdown scaffold.json          # Marp-compatible slides
via54 export json scaffold.json              # structured data
via54 export render story.html --duration 30 # video (needs ffmpeg)
via54 export tts --text "hello"              # text-to-speech
```

PPTX supports style + theme templates:

```bash
via54 export pptx --style minimal            # minimal
via54 export pptx --style editorial          # editorial
via54 export pptx --style bold               # bold
via54 export pptx --style accent-bar         # accent bar (default)

via54 export pptx --theme templates/color-schemes/ink-wash.yaml
```

### 🔌 Forge / ComfyUI Integration

```bash
via54 forge --workflow sdxl_txt2img --prompt "a cat" --send
via54 forge --list
via54 comfyui --workflow sdxl_txt2img --prompt "a cat"
```

### 🖼️ Media Pipeline

```bash
via54 media add-music input.mp4 --mood tech
via54 media convert input.mp4
via54 media trace --input sketch.jpg
```

---

## 🌐 Web UI (Intent-driven)

```bash
via54 web --port 8080
```

Open `http://localhost:8080` — **HTMX-powered, zero JavaScript**:

| Button | What it does |
|--------|--------------|
| 🎨 **Design** | Describe scene → pick style → HTML + scaffold → PPT framework |
| 📝 **Write prompts** | Text / image → structured prompt → submit to Forge |
| 🎬 **Make video** | Upload storyboard → narrative model → video script + prompts |
| ⚡ **Submit to Forge** | Send a prompt directly to Forge for re-generation |

All core features work **standalone** — no Forge/ComfyUI required. The backend only enhances the "re-generate" step.

The Web UI uses **HTMX 1.9.12** for all interactions. No JavaScript, no build step, no npm. 10 HTMX fragment endpoints, all returning server-rendered HTML.

**Endpoints (16 JSON + 10 HTMX):**

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/health` | GET | Health check |
| `/api/templates` | GET | Workflow template list |
| `/api/prompt` | POST | Generate prompt |
| `/api/generate` | POST | Generate HTML design |
| `/api/narrate` | POST | Generate narrative scaffold |
| `/api/build` | POST | Build ComfyUI workflow |
| `/api/export` | POST | Export (PPTX/PDF/SVG/MD/JSON) |
| `/api/upload` | POST | Upload image / document |
| `/api/analyze` | POST | Analyze image features |
| `/api/img2prompt` | POST | Image → prompt |
| `/api/story2ppt` | POST | Document → PPT framework |
| `/api/storyboard` | POST | Multi-image → narrative → video script |
| `/api/video-prompt` | POST | Single image → opening shot prompt |
| `/api/regen` | POST | Submit prompt to Forge |
| `/api/htmx/*` | various | HTMX fragment endpoints (10) |

---

## 📊 Template System

### Color schemes (31 sets)

| Category | Count | Examples |
|----------|-------|----------|
| Chinese traditional | 8 | `ink-wash` (sumi-e), `crimson-elegance` (vermilion), `moon-white` (sky after rain) |
| Japanese | 6 | `tsubaki-camellia` (red camellia), `wabi-sabi`, `rinpa-gold` |
| Classic | 6 | `dark-terminal-blue` (Linear-style), `cosmic-retro` (Perplexity-style) |
| Adobe / Behance | 6 | `spectrum-indigo`, `glassmorphism-pastel`, `neon-dark` |
| Extras | 5 | `bento-dark-glass`, `mono-brand-bold` |

### Typography (12 sets)

Chinese 6 sets (Ming / Kai / Fangsong / Hei / Calligraphy / Sans) + International 6 sets (Serif / Geometric / Display / Didone / Monospace / Rounded).

### Layouts (3 sets)

| Template | Adapts to | Features |
|----------|-----------|----------|
| `hero-split-16-9` | TV → Desktop → Tablet → Phone | Split-pane, responsive stack |
| `bento-grid-2x2` | TV 3×2 → Desktop 2×2 → Phone 1 col | Bento-box grid |
| `gallery-waterfall` | 5 col → 4 → 3 → 2 | Auto-fill waterfall |

### ComfyUI workflows (30 sets)

| Category | Templates |
|----------|-----------|
| Text-to-image | sdxl_txt2img, flux_dev_txt2img, sd15_txt2img, sd3_txt2img, pixart, playground_v2, stable_cascade |
| Image-to-image | sdxl_img2img, flux_img2img, sd3_img2img, sdxl_turbo, sdxl_refiner |
| Inpainting | sdxl_inpaint, flux_fill |
| Upscale | sdxl_upscale |
| ControlNet | controlnet_canny, controlnet_openpose |
| Video | animatediff_txt2vid, hunyuan_txt2vid, wan_txt2vid, ltxv_txt2vid, mochi_txt2vid, cosmos_txt2vid, svd_img2vid, svd_img2vid_ext |
| Advanced | sdxl_advanced (LoRA+IPAdapter+FaceRestore), flux_pro (Tiled), sdxl_tiled, sdxl_img2img_face, lcm_lora |

### PPTX styles (4 sets)

`minimal` / `editorial` / `bold` / `accent-bar` — YAML-defined positions, fonts, colors, accepts `--theme` to reference any of the 31 color schemes.

---

## 🏗️ Architecture

```
via54Design/
├── cmd/
│   ├── via54/          16 files — CLI entrypoint (main + 15 subcommands)
│   └── mcp-server/     MCP Server standalone binary
├── internal/
│   ├── export/         Export engine (pure Go: pptx/svg/json/markdown/pdf/tts)
│   ├── mcp/            MCP Server implementation
│   ├── media/          Media pipeline (download/vectorize/score)
│   ├── narrate/        Narrative engine (4 models, script/shot-list)
│   ├── pattern/        Design pattern extractor
│   ├── prompt/         Prompt engine v2.2 (6 modules)
│   │   ├── types.go         type definitions
│   │   ├── generator.go     generator (16 dimensions)
│   │   ├── templates.go     template loader + negative bank
│   │   ├── quality.go       quality assessment
│   │   ├── version.go       version management
│   │   └── render.go        Markdown/JSON rendering
│   ├── quality/        Quality gate
│   ├── template/       Template engine (layout/color/font)
│   ├── util/           Shared utilities (FindBaseDir)
│   ├── vision/         Computer-vision helpers
│   ├── wasm/           WASM bridge (Rust acceleration)
│   └── workflow/       ComfyUI workflow engine (30 templates)
├── templates/
│   ├── prompts/        17-platform prompt templates (YAML)
│   ├── layouts/        3 layout templates (16:9, 4-breakpoint)
│   ├── color-schemes/  31+ color schemes
│   ├── typography/     12 font definitions
│   ├── narratology/    4 narrative models
│   └── registry.yaml   template registry
├── web/                Web interface (Go handler + HTML/HTMX, **0 JavaScript**)
│   ├── handler.go      28KB HTTP handler
│   └── templates/      HTML fragments + HTMX
├── hack/               Build / deploy scripts
│   ├── build.sh        Go cross-compile (CLI + MCP)
│   ├── install.sh      one-shot installer
│   ├── setup.sh        full setup
│   ├── setup_macos.sh  macOS one-shot (Homebrew + npm + binary)
│   ├── setup_linux.sh  Linux one-shot (apt/dnf/pacman + binary)
│   └── wasm/           Rust WASM source
├── docs/
│   ├── prompts/        camera/lighting/color/composition references
│   ├── template-format.md
│   ├── failure-recovery.md
│   └── deployment-guide.md
├── test_samples/       test fixtures
├── test_20_rounds.py   20-round automated test suite
├── AGENTS.md           AI working context
├── SOUL.md             Hermes agent soul definition
├── Makefile            standard build automation
├── LICENSE             dual license (MIT OR AGPL-3.0)
└── README.md           this file
```

### Design philosophy

1. **YAML templates > hard-coded logic.** Color, font, layout, narrative model, PPTX style — all defined in YAML, swappable at runtime.
2. **Pure-Go single binary.** The core engine is 15 MB, zero runtime dependencies.
3. **Standalone-first.** All core features run without Forge/ComfyUI; those are optional accelerators.
4. **Deterministic output.** Every map iteration is sorted via `sortedKeys()` first — same input always produces the same MD5.
5. **API-first.** Every feature is callable via REST API. The Web UI is just a frontend.
6. **HTMX > JavaScript.** Server-rendered HTML fragments. Zero JS, zero build step, zero npm.

---

## 🔗 Reference projects (peer comparison)

| Category | Reference | ⭐ | via54 differentiation |
|----------|-----------|---|------------------------|
| Prompt engineering | easy-sd, sd-webui-forge | 10k–12k | Unified engine for 17 platforms, YAML-driven, 26-dim control |
| PPT generation | banana-slides | 14.8k | Narrative-driven PPT, 4 narrative models → slides |
| Narrative engine | similar projects | <100 | 4 formal narrative models, Fountain script + shot list |
| ComfyUI management | ComfyUI | 116k | Go execution bridge, 30 templates, deterministic seeds, testable |
| Design templates | huashu-design | 16.7k | Go re-write of the core layer, structured YAML, quality gate |

---

## 📜 License

- Go source: `AGPL-3.0-only`
- Templates / scripts / docs: `MIT`
- See [`LICENSE`](./LICENSE) and [`ACKNOWLEDGMENTS`](./ACKNOWLEDGMENTS)

---

## 🙏 Acknowledgments

- **Aaron Swartz** (1986–2013) and **Tim Berners-Lee** (1955–) — for the open web that made this work possible.
- **Lu You** (1125–1210) — Song dynasty poet, whose words guide the philosophy section.
- **William Shakespeare** (1564–1616) and **Oscar Wilde** (1854–1900) — for the borrowed lines.
- **huashu-design** — the upstream inspiration for the design template engine.
- The entire **open-source community** that built Go, HTMX, the MCP protocol, and every dependency we stand on.

> *"We are all in the gutter, but some of us are looking at the stars."*
> — Oscar Wilde, *Lady Windermere's Fan*

---

## 🇨🇳 Chinese Version

Read this README in Chinese: [**README.md**](./README.md)
