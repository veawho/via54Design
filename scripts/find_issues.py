#!/usr/bin/env python
# find_issues.py - 自动检索潜在推广 Issues 脚本
import json
import subprocess
import os
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent
HISTORY_FILE = ROOT / "scripts" / "commented_issues.json"
CANDIDATES_FILE = ROOT / "scripts" / "candidate_issues.json"

# 检索关键词定义
QUERIES = [
    "pptx notes state:open type:issue",
    "speaker notes pptx state:open type:issue",
    "markdown to pptx state:open type:issue",
    "slides to video state:open type:issue",
]

def load_history():
    if HISTORY_FILE.exists():
        try:
            return set(json.loads(HISTORY_FILE.read_text(encoding="utf-8")))
        except Exception:
            return set()
    return set()

def save_history(history):
    HISTORY_FILE.parent.mkdir(parents=True, exist_ok=True)
    HISTORY_FILE.write_text(json.dumps(list(history), indent=2), encoding="utf-8")

def search_issues():
    history = load_history()
    candidates = []

    print("🔎 开始扫描 GitHub Issues...")
    for q in QUERIES:
        print(f"  -> 检索: {q}")
        cmd = [
            "gh", "api", 
            f"search/issues?q={q}",
            "--jq", ".items[] | {html_url: .html_url, number: .number, title: .title, repository_url: .repository_url}"
        ]
        r = subprocess.run(cmd, capture_output=True, text=True, shell=True)
        if r.returncode != 0:
            print(f"    ✗ 检索失败: {r.stderr.strip()}")
            continue
        
        output = r.stdout.strip()
        if not output:
            continue
            
        for line in output.splitlines():
            try:
                item = json.loads(line)
                url = item["html_url"]
                if url not in history:
                    candidates.append(item)
            except Exception:
                continue

    # 保存候选 Issue
    CANDIDATES_FILE.write_text(json.dumps(candidates, indent=2, ensure_ascii=False), encoding="utf-8")
    print(f"✅ 扫描结束。发现 {len(candidates)} 个新的候选 Issues。")
    return candidates

if __name__ == "__main__":
    search_issues()
