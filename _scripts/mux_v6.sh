#!/bin/bash
# mux_v6.sh - 旁白 + 配乐 + 视频 三合一混音
# 用法: bash mux_v6.sh <lang>   # lang=zh|en|ja

set -e
LANG="${1:-zh}"
VID="G:/agent/hermes/via54Design-v6/output/lithium_30s_v6.mp4"
OUT="G:/agent/hermes/via54Design-v6/output/lithium_30s_v6_${LANG}.mp4"
BGM="G:/agent/hermes/via54Design-v6/music/bgm_epic_30s.mp3"

# 拼接 5 段旁白成 1 个完整音轨
VOICE_LIST=""
for k in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  VOICE_LIST="$VOICE_LIST G:/agent/hermes/via54Design-v6/voice/${LANG}_${k}.mp3"
done

VOICE_TMP="G:/agent/hermes/via54Design-v6/output/_voice_${LANG}.mp3"

echo "▶ 拼接 ${LANG} 5 段旁白 → $VOICE_TMP"
# 5 段 mp3 用 concat demuxer 拼接 (mp3 编码参数一致)
cat > /tmp/concat_voice.txt <<EOF
file 'G:/agent/hermes/via54Design-v6/voice/${LANG}_01_hook.mp3'
file 'G:/agent/hermes/via54Design-v6/voice/${LANG}_02_trend.mp3'
file 'G:/agent/hermes/via54Design-v6/voice/${LANG}_03_tech.mp3'
file 'G:/agent/hermes/via54Design-v6/voice/${LANG}_04_market.mp3'
file 'G:/agent/hermes/via54Design-v6/voice/${LANG}_05_outlook.mp3'
EOF
ffmpeg -y -f concat -safe 0 -i /tmp/concat_voice.txt -c copy "$VOICE_TMP" 2>&1 | tail -2

# 验证旁白时长
VOICE_DUR=$(ffprobe -v error -show_entries format=duration -of csv=p=0 "$VOICE_TMP")
echo "  旁白总时长: ${VOICE_DUR}s"

# ffmpeg 混音: 视频 (静音) + 旁白 (1.0) + 配乐 (0.13) → 输出
echo "▶ 三路混音 → $OUT"
# 视频 28s + 旁白 25-28s + 配乐 atrim 28s, 三路自然走完, ffmpeg 取最长
ffmpeg -y \
  -i "$VID" \
  -i "$VOICE_TMP" \
  -i "$BGM" \
  -filter_complex "
    [1:a]volume=1.4[voice];
    [2:a]volume=0.13,atrim=0:28.0[bgm];
    [voice][bgm]amix=inputs=2:duration=longest:dropout_transition=2[mix]
  " \
  -map 0:v -map "[mix]" \
  -c:v copy -c:a aac -b:a 192k \
  "$OUT" 2>&1 | tail -5

echo ""
echo "=== 最终输出 ==="
ls -la "$OUT" 2>&1
ffprobe -v error -show_entries format=duration,size,bit_rate -show_entries stream=codec_name,width,height -of default=noprint_wrappers=1 "$OUT" 2>&1
