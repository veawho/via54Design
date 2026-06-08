#!/usr/bin/env bash
# Example 3: Gallery Showcase
# 8-tile portfolio grid
set -e
cd "$(dirname "$0")"
go build -o ../../via54.exe ../../cmd/via54/
../../via54.exe generate \
  --layout gallery-waterfall \
  --color rinpa-gold \
  --font elegant-didone \
  --title "2026 品牌组合" \
  --output gallery.html
echo "✓ Generated: gallery.html"
