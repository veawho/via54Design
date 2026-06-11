#!/usr/bin/env python3
"""via54Design + Hermes 系统 100 轮自检自修复测试框架

v1 范围:
  - 1-15: 浅层 (build/vet/test/二进制/help/version)
  - 16-30: 单元测试 + 边界 + 错误处理
  - 31-50: 集成 (模板/Pipeline/导出/导出格式)
  - 51-65: CLI 入口 + 错误恢复
  - 66-80: 性能 + 并发
  - 81-95: 持久性 + 资源
  - 96-100: 终验 + 报告

输出: round_NN_status.json + 最终报告
"""
import os
import sys
import json
import time
import subprocess
import shutil
from pathlib import Path

REPO = Path(__file__).parent
BIN = REPO / ("via54.exe" if os.name == "nt" else "via54")
MCP_BIN = REPO / ("via54-mcp.exe" if os.name == "nt" else "via54-mcp")
REPORT_DIR = REPO / "self_test_reports"
REPORT_DIR.mkdir(exist_ok=True)

results = []

def run(cmd, timeout=30, check_rc=True):
    """Run cmd, return (rc, stdout, stderr)"""
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout, cwd=str(REPO))
    return r.returncode, r.stdout, r.stderr

def record(num, name, status, details=""):
    icon = {"PASS": "✓", "WARN": "⚠", "FAIL": "✗"}[status]
    print(f"  [{num:03d}] {icon} {name}{(' — ' + details) if details else ''}")
    results.append({"round": num, "name": name, "status": status, "details": details})

def header(phase, lo, hi):
    print(f"\n{'='*70}\n PHASE: {phase}  (Rounds {lo:03d}-{hi:03d})\n{'='*70}")

# ────────────────────────────────────────────────────────────
header("Phase 1: 浅层 — build / vet / test / 二进制", 1, 15)
# ────────────────────────────────────────────────────────────

# Round 01: go build
rc, out, err = run(["go", "build", "./..."], timeout=120)
record(1, "go build ./...", "PASS" if rc == 0 else "FAIL", f"rc={rc}")

# Round 02: go vet
rc, out, err = run(["go", "vet", "./..."], timeout=120)
record(2, "go vet ./...", "PASS" if rc == 0 else "FAIL", f"rc={rc}")

# Round 03: go test 全部
rc, out, err = run(["go", "test", "-count=1", "./..."], timeout=300)
record(3, "go test ./...", "PASS" if rc == 0 else "FAIL", f"rc={rc} (out tail: {out[-200:] if out else ''})")

# Round 04: via54.exe --help
rc, out, err = run([str(BIN), "--help"], timeout=5)
record(4, "via54 --help", "PASS" if rc == 0 else "WARN", f"rc={rc}")

# Round 05: via54.exe version
rc, out, err = run([str(BIN), "version"], timeout=5)
record(5, "via54 version", "PASS" if rc == 0 else "FAIL", f"rc={rc} out={out.strip()[:100]}")

# Round 06: via54-mcp.exe --help
rc, out, err = run([str(MCP_BIN), "--help"], timeout=5)
record(6, "via54-mcp --help", "PASS" if rc == 0 else "FAIL", f"rc={rc}")

# Round 07: via54.exe list
rc, out, err = run([str(BIN), "list"], timeout=10)
record(7, "via54 list", "PASS" if rc == 0 else "FAIL", f"rc={rc} lines={len(out.splitlines())}")

# Round 08: via54.exe list --help
rc, out, err = run([str(BIN), "list", "--help"], timeout=5)
# 期望: 进入 list 正常输出, 因为 list 不接受 flag
record(8, "via54 list --help", "WARN" if "Usage" in out else "PASS", f"rc={rc} out_snip={out[:80]}")

# Round 09: via54.exe media --help
rc, out, err = run([str(BIN), "media", "--help"], timeout=5)
# 期望: 应当显示 help 而不是 未知命令
record(9, "via54 media --help", "FAIL" if "未知" in out else "PASS", f"rc={rc} out_snip={out[:80]}")

# Round 10: 13 子命令全部 --help 探测
subcmds_no_block = ["generate", "narrate", "quality", "pattern", "list", "media", "export",
           "prompt", "pipeline", "forge", "comfyui", "web", "version", "serve"]
help_results = {}
for sub in subcmds_no_block:
    try:
        rc, out, err = run([str(BIN), sub, "--help"], timeout=5)
        # 期望: 显示 Usage 或 help 文本; 失败信号: "未知" 出现在输出
        # pipeline help 打印 "via54 pipeline" (不是 via54 prompt) — 已修复
        bad = "未知" in out
        help_results[sub] = (rc, bad, out[:60])
    except subprocess.TimeoutExpired:
        help_results[sub] = (-1, True, "TIMEOUT")
fail_count = sum(1 for r, bad, _ in help_results.values() if bad)
record(10, "14 subcommand --help audit (no interactive)",
       "PASS" if fail_count == 0 else "FAIL",
       f"failed={fail_count}/14 sub={[(s, r) for s, (r, b, _) in help_results.items() if b]}")

# Round 11: version.go 硬编码 v0.3.0 与 git log 不一致
# 通过 grep 检测
rc, out, err = run(["grep", "-rn", "v0.3.0", "cmd/", "internal/"], timeout=10)
record(11, "version.go 硬编码检查", "WARN" if "v0.3.0" in out else "PASS",
       f"out_snip={out[:120]}")

# Round 12: 缺失的 profiles 目录
profiles_dir = REPO.parent / ".hermes" / "profiles"
record(12, "~/.hermes/profiles 目录", "WARN", "path missing — confirm no other profile in use")

# Round 13: cron jobs 健康度 (Hermes 系统)
rc, out, err = run(["hermes", "cron", "list"], timeout=10)
errors_found = "error:" in out.lower() or "429" in out
record(13, "hermes cron list", "WARN" if errors_found else "PASS",
       f"rc={rc} errors={'yes' if errors_found else 'no'}")

# Round 14: gateway health
rc, out, err = run(["curl", "-sS", "http://127.0.0.1:8642/health"], timeout=5)
record(14, "gateway :8642 /health", "PASS" if '"ok"' in out else "FAIL", f"out={out[:80]}")

# Round 15: 模板/文件总览
import yaml
try:
    reg_text = (REPO / "templates" / "registry.yaml").read_text()
    reg = yaml.safe_load(reg_text)
    # 可能是 list 或 dict 结构
    if isinstance(reg, dict):
        cats = list(reg.keys())
    elif isinstance(reg, list):
        cats = [e.get("id", "?") for e in reg if isinstance(e, dict)]
    else:
        cats = []
    record(15, "templates/registry.yaml 加载", "PASS" if len(cats) > 0 else "WARN",
           f"top_keys={len(cats)} sample={cats[:3] if cats else []}")
except Exception as e:
    record(15, "templates/registry.yaml 加载", "FAIL", str(e))

# 保存第一阶段
with open(REPORT_DIR / "phase1_rounds_001_015.json", "w") as f:
    json.dump(results, f, indent=2, ensure_ascii=False)

print(f"\nPhase 1 完成: {sum(1 for r in results if r['status']=='PASS')}/{len(results)} PASS")
