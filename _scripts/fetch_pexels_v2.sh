#!/bin/bash
# fetch_pexels_v2.sh - 直接 curl 下载 (避开子 shell -o 坑)
set -e
KEY="aHyfRPK9DM8s7nV4Rv9xVK7aIDmkKoLwOx0tZzcGtxDaI4zhftgRvDPO"
DIR="G:/agent/hermes/via54Design-v6/stock"
mkdir -p "$DIR"
cd "$DIR"

declare -A QUERIES=(
  ["01_hook"]="battery+factory+night+light"
  ["02_trend"]="data+graph+screen+stock"
  ["03_tech"]="robotic+arm+manufacturing+precision"
  ["04_market"]="trading+floor+screen+financial"
  ["05_outlook"]="electric+vehicle+sunrise+city"
)

for seg in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  query="${QUERIES[$seg]}"
  echo "▶ $seg: $query"

  # 拉 3 个候选, 写文件
  curl -sL -H "Authorization: $KEY" \
    "https://api.pexels.com/v1/videos/search?query=${query}&per_page=3&orientation=landscape&min_duration=5&max_duration=20" \
    > "$seg.json"

  # 提 URL 列表
  python -c "
import json
d = json.load(open(r'$seg.json'))
with open(r'${seg}_urls.txt', 'w') as f:
    for v in d.get('videos', []):
        for f1 in v.get('video_files', []):
            if f1.get('quality') == 'hd' and 'mp4' in f1.get('file_type', ''):
                f.write(f1['link'] + '\n'); break
        else:
            for f1 in v.get('video_files', []):
                if f1.get('quality') == 'sd' and 'mp4' in f1.get('file_type', ''):
                    f.write(f1['link'] + '\n'); break
"

  # 串行下载
  i=1
  while IFS= read -r url; do
    if [ -n "$url" ]; then
      out="${seg}_${i}.mp4"
      [ -f "$out" ] && [ -s "$out" ] && echo "  ✓ $(basename $out) (cached)" && i=$((i+1)) && continue
      echo "  ↓ $(basename $out) ← $(echo $url | grep -oE 'pexels.com/[^/]+/[^/]+/')"
      curl -sL -A "Mozilla/5.0" --max-time 90 -o "$out" "$url"
      sz=$(stat -c%s "$out" 2>/dev/null || echo 0)
      echo "    $(($sz / 1024)) KB"
      i=$((i+1))
    fi
  done < "${seg}_urls.txt"
done

echo ""
echo "=========================================="
echo "最终素材 ($(ls *.mp4 2>/dev/null | wc -l) 个 mp4):"
ls -la *.mp4 2>&1
