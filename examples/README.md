# via54Design Examples

5 working examples demonstrating the via54Design workflow.
Each example is self-contained and runnable with `via54`.

## Quick Start

```bash
# 1. Build via54 first
go build -o via54.exe ./cmd/via54/

# 2. Try any example
./via54.exe generate --layout hero-split-16-9 --color ink-wash --font ming-hei-editorial --title "我的标题"
```

## Examples

| # | Name | Layout | Color | Font | Use Case |
|---|------|--------|-------|------|----------|
| 1 | [Basic Hero](./01-basic-hero) | hero-split-16-9 | ink-wash | ming-hei-editorial | Product launch |
| 2 | [Bento Dashboard](./02-bento-dashboard) | bento-grid-2x2 | candy-duolingo | hei-modern | Data dashboard |
| 3 | [Gallery Showcase](./03-gallery-showcase) | gallery-waterfall | rinpa-gold | elegant-didone | Brand portfolio |
| 4 | [Narrative Story](./04-narrative-story) | (from scaffold) | (mood-driven) | (narrative-driven) | 30-second video |
| 5 | [Prompt Engineering](./05-prompt-engineering) | n/a | n/a | n/a | 17 AI platforms |

## Running Examples

Each example directory has:
- `run.sh` (or `run.bat` on Windows): one-click run
- `expected.html`: what the output should look like
- `notes.md`: design rationale and customization tips

## Adding Your Example

See [CONTRIBUTING.md](../CONTRIBUTING.md) for the template.
