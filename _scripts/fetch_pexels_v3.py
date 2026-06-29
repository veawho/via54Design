#!/usr/bin/env python
# fetch_pexels_v3.py - 用 Python 直接处理 Pexels API + 下载
import json
import subprocess
import os

KEY = "aHyfRPK9DM8s7nV4Rv9xVK7aIDmkKoLwOx0tZzcGtxDaI4zhftgRvDPO"
DIR = r"G:\agent\hermes\via54Design-v6\stock"
os.makedirs(DIR, exist_ok=True)

# 5 段关键词
SEGMENTS = {
    "01_hook":   "factory+industrial+steel+fire",
    "02_trend":  "data+visualization+dashboard+hologram",
    "03_tech":   "technician+laboratory+battery",
    "04_market": "stock+market+candle+finance",
    "05_outlook":"electric+car+road+sunrise",
}

for seg, query in SEGMENTS.items():
    print(f"\n▶ {seg}: {query}")
    # 1) 拉 Pexels API
    api_url = f"https://api.pexels.com/v1/videos/search?query={query}&per_page=3&orientation=landscape&min_duration=5&max_duration=20"
    result = subprocess.run(
        ["curl", "-sL", "-H", f"Authorization: {KEY}", api_url],
        capture_output=True, text=True, timeout=30
    )
    try:
        data = json.loads(result.stdout)
    except json.JSONDecodeError as e:
        print(f"  ✗ JSON 解析失败: {e}")
        print(f"  响应: {result.stdout[:200]}")
        continue

    total = data.get("total_results", 0)
    print(f"  total: {total}")
    videos = data.get("videos", [])
    if not videos:
        print(f"  ✗ 无结果")
        continue

    # 2) 提取 URL 列表 (SD 优先, 锂电 30s 用 1280x720 够)
    urls = []
    for v in videos:
        files = v.get("video_files", [])
        # 找 1280x720 或类似 HD, 但 SD 也可
        candidates = sorted(
            [f for f in files if "mp4" in f.get("file_type", "")],
            key=lambda f: -(f.get("width", 0) or 0)  # 高分辨率优先
        )
        if candidates:
            # 优先选 1280-1920 宽的 (避免过大 4K)
            for c in candidates:
                w = c.get("width", 0) or 0
                if 1280 <= w <= 1920:
                    urls.append(c["link"]); break
            else:
                urls.append(candidates[0]["link"])
    print(f"  URLs: {len(urls)}")

    # 3) 串行下载
    for i, url in enumerate(urls, 1):
        out_path = os.path.join(DIR, f"{seg}_{i}.mp4")
        if os.path.exists(out_path) and os.path.getsize(out_path) > 0:
            sz = os.path.getsize(out_path) // 1024
            print(f"  ✓ {os.path.basename(out_path)} cached ({sz}KB)")
            continue
        print(f"  ↓ {os.path.basename(out_path)}")
        dl = subprocess.run(
            ["curl", "-sL", "-A", "Mozilla/5.0", "--max-time", "90", "-o", out_path, url],
            capture_output=True, text=True, timeout=120
        )
        if os.path.exists(out_path) and os.path.getsize(out_path) > 0:
            sz = os.path.getsize(out_path) // 1024
            print(f"    {sz}KB ✓")
        else:
            print(f"    ✗ 下载失败 (0 bytes)")

print("\n" + "=" * 50)
print("最终素材:")
mps = sorted([f for f in os.listdir(DIR) if f.endswith(".mp4")])
for f in mps:
    p = os.path.join(DIR, f)
    print(f"  {f}: {os.path.getsize(p)//1024}KB")
print(f"  总计: {len(mps)} 个")
