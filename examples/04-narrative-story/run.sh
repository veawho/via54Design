#!/usr/bin/env bash
# Example 4: Narrative Story
# Story-driven 3-act video script → storyboard
set -e
cd "$(dirname "$0")"
go build -o ../../via54.exe ../../cmd/via54/

# 1. Narrate from seed
../../via54.exe narrate \
  --seed "中国平安从 1988 创立到 2026 的品牌升级之旅" \
  --model three-act \
  --output scaffold.json

# 2. Generate storyboard HTML from scaffold
../../via54.exe generate \
  --from-narrative scaffold.json \
  --output storyboard.html

# 3. Export to markdown (for Marp/Keynote)
../../via54.exe export markdown storyboard.html

echo "✓ Generated: storyboard.html + scenes.md"
