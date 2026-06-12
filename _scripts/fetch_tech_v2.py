#!/usr/bin/env python
# fetch_tech_v2.py - 段 3 重拉: 找真锂电电芯/电池工厂, 替换试管
import json
import subprocess
import os

# Pexels API key: 从环境变量 PEXELS_API_KEY 读
# 旧 key 已在 2026-06-12 泄露, 已在 GitHub 公开 commit 出现
# 用户应 revoke 旧 key 并到 https://www.pexels.com/api/ 申请新 key
KEY = os.environ.get("PEXELS_API_KEY", "")
if not KEY:
    raise SystemExit(
        "错误: 未设置 PEXELS_API_KEY 环境变量\n"
        "获取方式: https://www.pexels.com/api/ 注册后即可获得\n"
        "设置方式 (Windows): set PEXELS_API_KEY=你的新key\n"
        "设置方式 (bash): export PEXELS_API_KEY=你的新key"
    )
DIR = r"G:\agent\hermes\via54Design-v6\stock"
os.makedirs(DIR, exist_ok=True)

# 8 个新关键词, 优先电芯/电池工厂
queries = [
    "lithium+battery+factory",
    "lithium+battery+cell",
    "lithium+ion+battery+production",
    "battery+cell+production",
    "electric+vehicle+battery+factory",
    "lithium+production+line",
    "battery+manufacturing+factory",
    "battery+pack+assembly",
]

print("=== 8 关键词 × 3 候选 = 24 mp4 ===")
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
    print(f"  {q:42s}: {len(videos)} returned, total {total}")
    if not videos:
        continue
    for i, v in enumerate(videos[:3], 1):
        # 选 SD 优先 (1300x720 或 sd 都可)
        video_files = v.get("video_files", [])
        # 优先 1366x768 (匹配 v6 视频) > 1920x1080 (HD) > 1280x720 (SD) > 其他
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
        out = os.path.join(DIR, f"03_tech_{safe}_{i}.mp4")
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
