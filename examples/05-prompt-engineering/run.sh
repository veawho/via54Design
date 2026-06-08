#!/usr/bin/env bash
# Example 5: Prompt Engineering
# Generate prompts for 17 AI platforms
set -e
cd "$(dirname "$0")"
go build -o ../../via54.exe ../../cmd/via54/

# 1. Generate prompt for Midjourney
../../via54.exe prompt \
  --scene "中国平安 logo 焕新, 牡丹花 + 祥云, 高级感" \
  --platform midjourney \
  --output prompt-mj.md

# 2. Same scene, different platform
../../via54.exe prompt \
  --scene "中国平安 logo 焕新, 牡丹花 + 祥云, 高级感" \
  --platform flux \
  --output prompt-flux.md

# 3. List all 17 platforms
../../via54.exe prompt list

echo "✓ Generated: prompt-mj.md, prompt-flux.md"
