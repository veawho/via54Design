#!/usr/bin/env bash
# Example 2: Bento Dashboard
# 4-quadrant data display
set -e
cd "$(dirname "$0")"
go build -o ../../via54.exe ../../cmd/via54/
../../via54.exe generate \
  --layout bento-grid-2x2 \
  --color candy-duolingo \
  --font hei-modern \
  --title "Q1 2026 数据看板" \
  --output dashboard.html
echo "✓ Generated: dashboard.html"
