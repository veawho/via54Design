#!/bin/bash
# gen_voice.sh - mmx TTS 批量生成 5 段中英日旁白
set -e
cd "G:/agent/hermes/via54Design-v6"
mkdir -p voice

# 5 段中文 (男声 jingying = 精英)
declare -A ZH=(
  ["01_hook"]="千度熔炼, 钢铁般的决心, 锂电新纪元由此开启。"
  ["02_trend"]="2026年产能突破八百吉瓦时, 中国锂电, 领跑全球。"
  ["03_tech"]="固态电池技术领跑, 安全性能全面提升, 新能源汽车的未来, 就在此刻。"
  ["04_market"]="板块资金流入同比增长两倍, 锂电投资窗口正在打开。"
  ["05_outlook"]="从中国制造, 到中国智造, 锂电产业, 未来已来。"
)

# 5 段英文 (女声 narrator)
declare -A EN=(
  ["01_hook"]="A thousand degrees of forging, with the will of steel, a new era of lithium batteries begins."
  ["02_trend"]="In 2026, production capacity surpasses 800 GWh. China's lithium battery industry leads the world."
  ["03_tech"]="Solid-state battery technology leads the way, safety fully upgraded, the future of new energy is now."
  ["04_market"]="Sector capital inflow surged 200 percent year over year, the lithium battery investment window is opening."
  ["05_outlook"]="From Made in China to Smart in China, the lithium battery industry, the future is here."
)

# 5 段日文 (女声 yujie 优雅)
declare -A JA=(
  ["01_hook"]="千度の精錬、鋼鉄の決意、リチウム電池新時代の幕開け。"
  ["02_trend"]="2026年、生産能力が800ギガワット時を突破、中国のリチウム電池が世界をリード。"
  ["03_tech"]="全固体電池技術がリード、安全性能が全面的に向上、新エネルギー自動車の未来は今ここにある。"
  ["04_market"]="セクター資金流入は前年比2倍に急増、リチウム電池投資の窓が開いている。"
  ["05_outlook"]="中国製造から中国スマート製造へ、リチウム電池産業、未来はもう来た。"
)

echo "=== 生成中文 (5 段) ==="
for k in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  echo "▶ $k: ${ZH[$k]}"
  mmx speech synthesize --text "${ZH[$k]}" --voice male-qn-jingying --language zh --out "voice/zh_${k}.mp3" --quiet 2>&1 | tail -1
done

echo ""
echo "=== 生成英文 (5 段) ==="
for k in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  echo "▶ $k: ${EN[$k]}"
  mmx speech synthesize --text "${EN[$k]}" --voice female-yujie --language en --out "voice/en_${k}.mp3" --quiet 2>&1 | tail -1
done

echo ""
echo "=== 生成日文 (5 段) ==="
for k in 01_hook 02_trend 03_tech 04_market 05_outlook; do
  echo "▶ $k: ${JA[$k]}"
  mmx speech synthesize --text "${JA[$k]}" --voice female-yujie --language ja --out "voice/ja_${k}.mp3" --quiet 2>&1 | tail -1
done

echo ""
echo "=== 最终 (15 段) ==="
ls -la voice/*.mp3 2>&1
