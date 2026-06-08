# Acknowledgments

> *"No man is an island, entire of itself; every man is a piece of the continent."*
> — John Donne, Meditation XVII

`via54Design` is a small project standing on the shoulders of many giants.
This file lists the people, projects, and ideas that made it possible.

## Inspirations

### The philosophers

- **Aaron Swartz** (1986–2013) — co-author of RSS, founder of Reddit,
  open-access champion. His life reminds us that the Internet's gifts
  belong to everyone.
- **Tim Berners-Lee** (b. 1955) — inventor of the World Wide Web. He gave
  us a tool for the universal sharing of knowledge, and then gave it
  away.
- **Lu You** (陆游, 1125–1210) — Song-dynasty poet whose line
  *"文章本天成，妙手偶得之"* (a work of literature is pre-ordained;
  the masterful hand only occasionally discovers it) guides our
  philosophy of human-AI collaboration.
- **William Shakespeare** (1564–1616) — for the borrowed lines and the
  reminder that every human story is a stage.
- **Oscar Wilde** (1854–1900) — for the reminder that imitating the great
  patterns is the sincerest form of flattery.

### The upstream project

- **[huashu-design](https://github.com/alchaincyf/huashu-design)** (MIT, 16k+ ★) — the
  original Node.js design template engine. `via54Design`'s core engine is a
  from-scratch Go re-implementation of its template-composition idea, with
  expanded scope (narrative engine, prompt engine, MCP server, 5 platforms).

## Open-source dependencies

### Direct Go dependencies

| Package | License | Use |
|---------|---------|-----|
| [`mcp-go`](https://github.com/mark3labs/mcp-go) | MIT | MCP Server protocol |
| [`yaml.v3`](https://github.com/go-yaml/yaml) | Apache-2.0 + MIT | YAML template parser |
| [`wazero`](https://github.com/tetratelabs/wazero) | Apache-2.0 | WASM runtime (transitive) |
| [`uuid`](https://github.com/google/uuid) | BSD-3-Clause | ID generation (transitive) |
| [`golang.org/x/sys`](https://pkg.go.dev/golang.org/x/sys) | BSD-3-Clause | OS syscall helpers |
| [`golang.org/x/text`](https://pkg.go.dev/golang.org/x/text) | BSD-3-Clause | Text encoding |

### Runtime / optional

| Tool | License | Use |
|------|---------|-----|
| [Go](https://go.dev) | BSD-3-Clause | The toolchain |
| [HTMX](https://htmx.org) | BSD-2-Clause | Web UI without JavaScript |
| [Node.js](https://nodejs.org) | MIT | PDF / video export (Playwright) |
| [Playwright](https://playwright.dev) | Apache-2.0 | Browser automation |
| [FFmpeg](https://ffmpeg.org) | LGPL-2.1+ / GPL-2+ | Video encoding |
| [Rust](https://www.rust-lang.org) | MIT / Apache-2.0 | WASM module (optional) |

## Reference projects

We learned from (and stole ideas from) the following open-source projects:

| Project | License | What we took |
|---------|---------|--------------|
| [ComfyUI](https://github.com/comfyanonymous/ComfyUI) | GPL-3.0 | Workflow node graph pattern |
| [banana-slides](https://github.com) | MIT | PPTX structure patterns |
| [sd-webui-forge](https://github.com/lllyasviel/stable-diffusion-webui-forge) | Apache-2.0 | Prompt engineering insights |
| [easy-sd](https://github.com/cheeaun/easy-sd) | MIT | API ergonomics |

## Standards we follow

- **Semantic Versioning** 2.0.0
- **Conventional Commits** 1.0.0
- **SPDX License List** (every file has `SPDX-License-Identifier`)
- **GitHub Community Health** (issue templates, PR template, CODEOWNERS)
- **OpenSSF Best Practices** (in progress)

## Security researchers

We thank the following individuals for responsibly disclosing security issues
to us. *(none yet — be the first!)*

## Community contributors

Every PR, every Issue, every Discussion post makes this project better.
You are the people this project exists for.

— veawho (via54), 2026
