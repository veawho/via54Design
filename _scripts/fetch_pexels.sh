#!/bin/bash
# fetch_pexels.sh - 拉 Pexels 真实视频素材 (5 段 dtc-testimonial 结构)
# 输出到 G:/agent/hermes/via54Design-v6/stock/

set -e
KEY="aHyfRPK9DM8s7nV4Rv9xVK7aIDmkKoLwOx0tZzcGtxDaI4zhftgRvDPO"
DIR="G:/agent/hermes/via54Design-v6/stock"
mkdir -p "$DIR"

# 5 段 × 关键词 (中文注释: 锂电 30s 行业趋势)
declare -A SEGMENTS=(
  [01_hook]="battery+factory+night+light"
  [02_trend]="data+graph+screen+stock"
  [03_tech]="robotic+arm+manufacturing+precision"
  [04_market]="trading+floor+screen+financial"
  [05_outlook]="electric+vehicle+sunrise+city"
)

cd "$DIR"
echo "Pexels 拉取 (5 段 × 3 候选 = 15 视频)"
echo "=========================================="

for seg in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  query="${SEGMENTS[$seg]}"
  echo ""
  echo "▶ 段 $seg: query='$query'"

  resp=$(curl -sL -H "Authorization: $KEY" \
    "https://api.pexels.com/v1/videos/search?query=${query}&per_page=3&orientation=landscape&min_duration=5&max_duration=20")

  # 列出候选
  echo "$resp" | python -c "
import json, sys
d = json.loads(sys.stdin.read())
print(f'  total: {d.get(\"total_results\", 0)}')
for i, v in enumerate(d.get('videos', [])):
    print(f'  [{i}] id={v[\"id\"]} dur={v[\"duration\"]}s w={v.get(\"width\")} h={v.get(\"height\")}')
" 2>&1

  # 提取下载链接 (HD 优先, SD 降级)
  echo "$resp" | python -c "
import json, sys
d = json.loads(sys.stdin.read())
for v in d.get('videos', []):
    for f in v.get('video_files', []):
        if f.get('quality') == 'hd' and 'mp4' in f.get('file_type', ''):
            print(f['link']); break
    else:
        for f in v.get('video_files', []):
            if f.get('quality') == 'sd' and 'mp4' in f.get('file_type', ''):
                print(f['link']); break
" > "$DIR/${seg}_urls.txt"

  # 串行下载
  i=1
  while IFS= read -r url; do
    if [ -n "$url" ]; then
      out="$DIR/${seg}_${i}.mp4"
      if [ ! -f "$out" ]; then
        echo "  ↓ 下载: $(basename $out)"
        curl -sL -A "Mozilla/5.0" --max-time 60 -o "$out" "$url" 2>&1 | tail -1
      fi
      i=$((i+1))
    fi
  done < "$DIR/${seg}_urls.txt"
done

echo ""
echo "=========================================="
echo "下载结果:"
ls -la "$DIR"/*.mp4 2>&1
