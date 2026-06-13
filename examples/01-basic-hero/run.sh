#!/usr/bin/env bash
# Example 1: Basic Hero
# Minimal example: 1 layout + 1 color + 1 font + 1 title
set -e
cd "$(dirname "$0")"
go build -o ../../via54.exe ../../cmd/via54/
../../via54.exe generate \
  --layout hero-split-16-9 \
  --color ink-wash \
  --font ming-hei-editorial \
  --title "中国平安 2026" \
  --output hero.html
echo "✓ Generated: hero.html ($(wc -c < hero.html) bytes)"
