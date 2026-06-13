#!/usr/bin/env python
# fetch_outlook_v2.py - 段 5 重拉: 找真锂电未来镜头 (EV road + 风电 + 智能), 替换沙漠
import json
import subprocess
import os

KEY = "aHyfRPK9DM8s7nV4Rv9xVK7aIDmkKoLwOx0tZzcGtxDaI4zhftgRvDPO"
DIR = r"G:\agent\hermes\via54Design-v6\stock"
os.makedirs(DIR, exist_ok=True)

# 6 个新关键词, 优先锂电/EV/绿色能源"未来" 视觉
queries = [
    "electric+vehicle+highway",
    "electric+vehicle+charging+station",
    "wind+turbine+landscape",
    "solar+panels+factory",
    "electric+car+city+road",
    "green+energy+future+city",
]

print("=== 6 关键词 × 3 候选 = 18 mp4 ===")
results = []
for q in queries:
    safe = q.replace("+", "_")
    api_url = f"https://api.pexels.com/videos/search?query={q}&per_page=3&orientation=landscape&size=sd"
    headers = ["-H", f"Authorization: {KEY}"]
    r = subprocess.run(
        ["curl", "-s"] + headers + [api_url],
        capture_output=True, text=True, timeout=30
    )
    try:
        data = json.loads(r.stdout)
    except Exception as e:
        print(f"  ✗ {q}: parse err {e}")
        continue
    videos = data.get("videos", [])
    total = data.get("total_results", 0)
    print(f"  {q:38s}: {len(videos)} returned, total {total}")
    if not videos:
        continue
    for i, v in enumerate(videos[:3], 1):
        video_files = v.get("video_files", [])
        chosen = None
        for vf in video_files:
            w = vf.get("width", 0)
            h = vf.get("height", 0)
            if 1200 <= w <= 1400 and 700 <= h <= 800:
                chosen = vf; break
        if not chosen:
            for vf in video_files:
                if vf.get("quality") == "hd":
                    chosen = vf; break
        if not chosen:
            chosen = video_files[0] if video_files else None
        if not chosen: continue
        link = chosen["link"]
        out = os.path.join(DIR, f"05_outlook_{safe}_{i}.mp4")
        if os.path.exists(out) and os.path.getsize(out) > 100000:
            print(f"    ✓ {i}: 缓存 {os.path.basename(out)} ({os.path.getsize(out)//1024}KB)")
            results.append(out); continue
        r = subprocess.run(["curl", "-sL", "-o", out, link], capture_output=True, timeout=120)
        if os.path.exists(out) and os.path.getsize(out) > 100000:
            print(f"    ✓ {i}: 拉 {os.path.basename(out)} ({os.path.getsize(out)//1024}KB)")
            results.append(out)
        else:
            print(f"    ✗ {i}: 失败")
print(f"\n=== 总拉: {len(results)} mp4 ===")
