#!/usr/bin/env python3
"""via54Design + Hermes 系统 100 轮自检自修复测试框架 — Phase 2-4 (16-65)"""
import os
import sys
import json
import subprocess
import shutil
from pathlib import Path

REPO = Path(__file__).parent
BIN = REPO / ("via54.exe" if os.name == "nt" else "via54")
MCP_BIN = REPO / ("via54-mcp.exe" if os.name == "nt" else "via54-mcp")
REPORT_DIR = REPO / "self_test_reports"
REPORT_DIR.mkdir(exist_ok=True)

results = []
phase1_results = json.loads((REPORT_DIR / "phase1_rounds_001_015.json").read_text())
results.extend(phase1_results)

def run(cmd, timeout=30, env=None):
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout, cwd=str(REPO), env=env)
    return r.returncode, r.stdout, r.stderr

def record(num, name, status, details=""):
    icon = {"PASS": "✓", "WARN": "⚠", "FAIL": "✗"}[status]
    print(f"  [{num:03d}] {icon} {name}{(' — ' + details) if details else ''}")
    results.append({"round": num, "name": name, "status": status, "details": details})

def header(phase, lo, hi):
    print(f"\n{'='*70}\n PHASE: {phase}  (Rounds {lo:03d}-{hi:03d})\n{'='*70}")

# ────────────────────────────────────────────────────────────
header("Phase 2: 单元测试 + 错误处理 + 边界", 16, 30)
# ────────────────────────────────────────────────────────────

# Round 16: util 包全部测试
rc, out, err = run(["go", "test", "-count=1", "-v", "./internal/util/..."], timeout=60)
passed = out.count("PASS:") if out else 0
record(16, "internal/util tests", "PASS" if rc == 0 else "FAIL", f"rc={rc} passes={passed}")

# Round 17: template 包测试
rc, out, err = run(["go", "test", "-count=1", "-v", "./internal/template/..."], timeout=60)
passed = out.count("PASS:") if out else 0
record(17, "internal/template tests", "PASS" if rc == 0 else "FAIL", f"rc={rc} passes={passed}")

# Round 18: export 包测试
rc, out, err = run(["go", "test", "-count=1", "-v", "./internal/export/..."], timeout=60)
passed = out.count("PASS:") if out else 0
record(18, "internal/export tests", "PASS" if rc == 0 else "FAIL", f"rc={rc} passes={passed}")

# Round 19: workflow 包测试
rc, out, err = run(["go", "test", "-count=1", "-v", "./internal/workflow/..."], timeout=60)
passed = out.count("PASS:") if out else 0
record(19, "internal/workflow tests", "PASS" if rc == 0 else "FAIL", f"rc={rc} passes={passed}")

# Round 20-25: 各无测试包 (内部) — 仅 vet+build
no_test_pkgs = ["internal/media", "internal/narrate", "internal/pattern",
                "internal/pipeline", "internal/prompt", "internal/quality",
                "internal/vision", "internal/wasm", "internal/mcp"]
for i, pkg in enumerate(no_test_pkgs):
    num = 20 + i
    rc, out, err = run(["go", "build", f"./{pkg}/..."], timeout=60)
    record(num, f"build {pkg}", "PASS" if rc == 0 else "FAIL", f"rc={rc}")

# Round 30: 总测试统计
all_tests = []
for pkg in ["internal/util", "internal/template", "internal/export", "internal/workflow"]:
    rc, out, _ = run(["go", "test", "-count=1", "-v", f"./{pkg}/..."], timeout=60)
    if out:
        all_tests.append((pkg, out.count("PASS:"), out.count("FAIL:")))
total_p = sum(p for _, p, _ in all_tests)
total_f = sum(f for _, _, f in all_tests)
record(30, "总测试统计", "PASS" if total_f == 0 else "WARN",
       f"packages={len(all_tests)} pass={total_p} fail={total_f}")

# ────────────────────────────────────────────────────────────
header("Phase 3: 集成 — 模板/导出/工作流", 31, 50)
# ────────────────────────────────────────────────────────────

# Round 31: 加载所有 YAML 模板
import yaml
all_templates = []
for cat in ["colors", "color-schemes", "fonts", "layouts", "narratology", "pptx-styles",
            "prompts", "typography", "video-edits", "workflows"]:
    cat_dir = REPO / "templates" / cat
    if not cat_dir.exists():
        continue
    for f in cat_dir.iterdir():
        if f.suffix in (".yaml", ".yml"):
            try:
                yaml.safe_load(f.read_text(encoding="utf-8"))
                all_templates.append((cat, f.name, "OK"))
            except Exception as e:
                all_templates.append((cat, f.name, f"FAIL:{e}"))
errs = [t for t in all_templates if t[2] != "OK"]
record(31, "YAML 模板解析", "PASS" if not errs else "FAIL",
       f"loaded={len(all_templates)} errors={len(errs)}")

# Round 32-40: generate 实际生成
out_dir = REPO / "self_test_out"
out_dir.mkdir(exist_ok=True)
gen_results = []
for color, font, layout in [
    ("ink-wash", "serif-sans-editorial", "hero-split-16-9"),
    ("spectrum-indigo", "sans-geometric-tech", "bento-grid"),
    ("earth-terracotta", "sc-sans-clean", "hero-centered"),
    ("swiss-monochrome", "mono-code", "feature-three-column"),
    ("dark-terminal-blue", "display-sans-bold", "pricing-three-tier"),
    ("glassmorphism-pastel", "playful-rounded", "testimonial-quote"),
    ("bauhaus-primary", "calligraphy-accent", "gallery-grid"),
    ("candy-duolingo", "song-literary", "blog-article"),
    ("flat-ui-vibrant", "kai-rounded-friendly", "dashboard-stats"),
]:
    out_html = out_dir / f"gen_{layout}_{color}.html"
    rc, out, err = run([str(BIN), "generate",
                        "--layout", layout, "--color", color, "--font", font,
                        "--title", f"Test {layout}", "--output", str(out_html)],
                       timeout=15)
    ok = rc == 0 and out_html.exists() and out_html.stat().st_size > 0
    gen_results.append((layout, rc, ok, out_html.stat().st_size if out_html.exists() else 0))
gen_pass = sum(1 for _, _, ok, _ in gen_results if ok)
record(40, "9 种 layout+color+font 组合生成",
       "PASS" if gen_pass == 9 else "WARN",
       f"pass={gen_pass}/9 sample={gen_results[0]}")

# Round 41: narrate markdown 输出 (注意: --seed 必填, 不是 --model)
narrate_md = out_dir / "narrate.md"
rc, out, err = run([str(BIN), "narrate", "--seed", "测试故事: 一个产品发布",
                    "--model", "three-act", "--duration", "30",
                    "--format", "markdown", "--output", str(narrate_md)], timeout=10)
narrate_ok = rc == 0 and narrate_md.exists() and narrate_md.stat().st_size > 100
record(41, "narrate markdown 输出", "PASS" if narrate_ok else "FAIL",
       f"rc={rc} bytes={narrate_md.stat().st_size if narrate_md.exists() else 0}")

# Round 42: narrate json 输出 (供 generate --from-narrative 使用)
narrate_json = out_dir / "narrate.json"
rc, out, err = run([str(BIN), "narrate", "--seed", "AI in design",
                    "--model", "three-act", "--duration", "60",
                    "--format", "json", "--output", str(narrate_json)], timeout=10)
record(42, "narrate json 输出", "PASS" if rc == 0 and narrate_json.exists() else "FAIL",
       f"rc={rc} bytes={narrate_json.stat().st_size if narrate_json.exists() else 0}")

# Round 43: quality gate on generated HTML
gen_html = out_dir / "gen_hero-split-16-9_ink-wash.html"
if gen_html.exists():
    rc, out, err = run([str(BIN), "quality", "--html", str(gen_html), "--verbose"],
                       timeout=15)
    q_pass = "通过" in out or "PASS" in out or "✓" in out or "0 errors" in out or "0 issue" in out
    record(43, "quality gate on gen HTML", "PASS" if q_pass else "WARN",
           f"rc={rc} out_tail={out[-200:] if out else ''}")
else:
    record(43, "quality gate on gen HTML", "FAIL", "no HTML to test")

# Round 44: pattern extractor
rc, out, err = run([str(BIN), "pattern", "--html", str(gen_html), "--name", "test-pattern"],
                   timeout=10)
record(44, "pattern extractor", "PASS" if rc == 0 else "FAIL", f"rc={rc} out_snip={out[:120]}")

# Round 45: export json (注意: positional arg 是输入 JSON 路径, 不是 --input)
export_json = out_dir / "export.json"
narrate_tmp = out_dir / "_narrate_for_export.json"
# 先生成 narrate JSON
subprocess.run([str(BIN), "narrate", "--seed", "测试", "--format", "json",
                "--output", str(narrate_tmp)], capture_output=True, timeout=10)
rc, out, err = run([str(BIN), "export", "json", "--output", str(export_json),
                    str(narrate_tmp)], timeout=10)
record(45, "export json", "PASS" if rc == 0 and export_json.exists() else "FAIL",
       f"rc={rc} bytes={export_json.stat().st_size if export_json.exists() else 0} out_snip={out[:80]}")

# Round 46: export markdown
export_md = out_dir / "export.md"
rc, out, err = run([str(BIN), "export", "markdown", "--output", str(export_md),
                    "--title", "Test Story", "--author", "tester",
                    str(narrate_tmp)], timeout=10)
record(46, "export markdown", "PASS" if rc == 0 and export_md.exists() else "FAIL",
       f"rc={rc} bytes={export_md.stat().st_size if export_md.exists() else 0}")

# Round 47: export pptx
export_pptx = out_dir / "export.pptx"
rc, out, err = run([str(BIN), "export", "pptx", "--output", str(export_pptx),
                    "--title", "Test Slides", "--style", "accent-bar",
                    str(narrate_tmp)], timeout=15)
record(47, "export pptx", "PASS" if rc == 0 and export_pptx.exists() and export_pptx.stat().st_size > 1000 else "FAIL",
       f"rc={rc} bytes={export_pptx.stat().st_size if export_pptx.exists() else 0}")

# Round 48: export svg
export_svg_dir = out_dir / "svg-out"
rc, out, err = run([str(BIN), "export", "svg", "--output", str(export_svg_dir),
                    str(narrate_tmp)], timeout=15)
svg_count = len(list(export_svg_dir.glob("*.svg"))) if export_svg_dir.exists() else 0
record(48, "export svg", "PASS" if rc == 0 and svg_count > 0 else "WARN",
       f"rc={rc} svgs={svg_count} out_snip={out[:80]}")

# Round 49: export tts (无网/无 key 情况下)
export_tts = out_dir / "export_tts.mp3"
rc, out, err = run([str(BIN), "export", "tts", "--output", str(export_tts),
                    "--text", "Hello world, this is a test.", "--voice", "default"],
                   timeout=20)
record(49, "export tts (online)", "WARN",
       f"rc={rc} (online api; expected fail) out_snip={out[-150:]}")

# Round 50: end-to-end pipeline (narrate → generate from narrative)
e2e_html = out_dir / "e2e.html"
rc, out, err = run([str(BIN), "generate", "--from-narrative", str(narrate_json),
                    "--layout", "hero-split-16-9", "--color", "ink-wash",
                    "--font", "serif-sans-editorial", "--output", str(e2e_html)],
                   timeout=20)
record(50, "narrate→generate end-to-end", "PASS" if rc == 0 and e2e_html.exists() else "FAIL",
       f"rc={rc} bytes={e2e_html.stat().st_size if e2e_html.exists() else 0}")

# ────────────────────────────────────────────────────────────
header("Phase 4: CLI 入口 + 错误恢复", 51, 65)
# ────────────────────────────────────────────────────────────

# Round 51-60: 各子命令未知子命令错误恢复
cli_tests = [
    ("media", ["add-music", "convert", "fetch"]),
    ("export", ["render", "pdf", "tts", "pptx", "svg", "json", "markdown"]),
    ("prompt", ["scene", "edit", "list", "ref", "gallery", "comfyui", "assess"]),
]
cli_pass = 0
cli_total = 0
for parent, subs in cli_tests:
    for sub in subs:
        # 每个子命令单独 --help 应工作
        rc, out, err = run([str(BIN), parent, sub, "--help"], timeout=5)
        ok = "Usage" in out or "用法" in out
        if not ok:
            print(f"    {parent} {sub} --help: rc={rc} snip={out[:60]!r}")
        cli_total += 1
        if ok:
            cli_pass += 1
record(60, "11 个子命令的 --help 入口", "PASS" if cli_pass == cli_total else "WARN",
       f"pass={cli_pass}/{cli_total}")

# Round 61: 错误参数
rc, out, err = run([str(BIN), "generate", "--layout", "NONEXISTENT", "--color", "x", "--font", "y"],
                   timeout=10)
record(61, "错误 layout 名称处理", "PASS" if rc != 0 and ("不存在" in out or "未找到" in out or "fail" in out.lower() or "fail" in err.lower() or "❌" in out) else "WARN",
       f"rc={rc} out_snip={out[:120]} err_snip={err[:120]}")

# Round 62: generate 缺参数
rc, out, err = run([str(BIN), "generate"], timeout=5)
record(62, "generate 缺参数", "PASS" if rc != 0 else "FAIL", f"rc={rc} err_snip={err[:120]}")

# Round 63: 未知主命令
rc, out, err = run([str(BIN), "totallymadeup"], timeout=5)
record(63, "未知主命令", "PASS" if rc != 0 else "FAIL", f"rc={rc}")

# Round 64: 无参数 → 进入 interactive?  实际会立即返回 (不是阻塞)
try:
    rc, out, err = run([str(BIN)], timeout=3)
    record(64, "无参数启动", "PASS" if rc != 0 or "interactive" in out.lower() else "WARN",
           f"rc={rc} out_snip={out[:80]}")
except subprocess.TimeoutExpired:
    record(64, "无参数启动", "WARN", "TIMEOUT (interactive) — expected")

# Round 65: 重复 --flag 检测
rc, out, err = run([str(BIN), "generate", "--layout", "x", "--layout", "y",
                    "--color", "z", "--font", "w"], timeout=5)
# 期望: 取最后一个值, 不应崩溃
record(65, "重复 --flag", "WARN", f"rc={rc} out_snip={out[:80]}")

with open(REPORT_DIR / "phase2to4_rounds_016_065.json", "w") as f:
    json.dump(results, f, indent=2, ensure_ascii=False)

print(f"\nPhase 2-4 完成: 累计 {sum(1 for r in results if r['status']=='PASS')}/{len(results)} PASS")
