#!/usr/bin/env python
# fetch_market_v2.py - 段 4 重拉 (找锂电股/锂电产线, 替换 CAKE)
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

# 段 4 新关键词 (锂电股 / 锂电产线 / 金融行业图表)
SEGMENT_04 = {
    "lithium+battery+production":  3,  # 锂电真实产线 (5101 hits)
    "electric+vehicle+factory":   3,  # 电动车厂 (7003 hits)
    "finance+chart+industry":      3,  # 金融行业图表 (6318 hits)
    "business+growth+chart":       3,  # 商业增长图
}


def fetch_video(url: str, out_path: str) -> bool:
    """下载单个视频"""
    if os.path.exists(out_path) and os.path.getsize(out_path) > 100000:
        return True
    r = subprocess.run([
        "curl", "-sL", "-A", "Mozilla/5.0", "--max-time", "120", "-o", out_path, url
    ], capture_output=True, timeout=130)
    return os.path.exists(out_path) and os.path.getsize(out_path) > 100000


def main():
    new_files = []
    for query, n_candidates in SEGMENT_04.items():
        print(f"\n▶ Query: {query} (n={n_candidates})")
        api_url = f"https://api.pexels.com/v1/videos/search?query={query}&per_page={n_candidates}&orientation=landscape&min_duration=5&max_duration=15"
        r = subprocess.run([
            "curl", "-sL", "-H", f"Authorization: {KEY}", api_url
        ], capture_output=True, text=True, timeout=30)
        try:
            data = json.loads(r.stdout)
        except json.JSONDecodeError:
            print(f"  ✗ JSON parse failed"); continue

        for i, v in enumerate(data.get("videos", [])[:n_candidates], 1):
            vid_id = v["id"]
            files = v.get("video_files", [])
            cands = sorted(
                [f for f in files if "mp4" in f.get("file_type", "")],
                key=lambda f: -(f.get("width", 0) or 0)
            )
            # 选 1280x720-1920x1080
            sel = next((f for f in cands if 1280 <= (f.get("width", 0) or 0) <= 1920), cands[0] if cands else None)
            if not sel:
                continue
            url = sel["link"]
            # 命名: 04_market_new_1.mp4, 04_market_new_2.mp4 ...
            safe_query = query.replace("+", "_")[:20]
            out_name = f"04_market_{safe_query}_{i}.mp4"
            out_path = os.path.join(DIR, out_name)
            print(f"  [{i}] id={vid_id} dur={v['duration']}s {sel.get('width')}x{sel.get('height')} → {out_name}")
            if fetch_video(url, out_path):
                size_kb = os.path.getsize(out_path) / 1024
                print(f"      ✓ {size_kb:.0f}KB")
                new_files.append(out_name)
            else:
                print(f"      ✗ 下载失败")

    print(f"\n=== 段 4 新候选: {len(new_files)} 个 ===")
    for f in new_files:
        size_mb = os.path.getsize(os.path.join(DIR, f)) / 1024 / 1024
        print(f"  {f}: {size_mb:.1f}MB")


if __name__ == "__main__":
    main()
