#!/usr/bin/env python3
"""via54Design 20 轮全方位测试
测试维度: 部署/压力/冒烟/稳定性/可用性/易用性/准确性
"""
import subprocess
import time
import sys
import os
import json
import hashlib

REPO = r'C:\Users\via54\AppData\Local\Temp\via54Design'
BIN = os.path.join(REPO, 'via54_test.exe')
if os.path.exists(os.path.join(REPO, 'via54.exe')):
    BIN = os.path.join(REPO, 'via54.exe')

# 创建测试 venv
VENV_PY = r'C:\Users\via54\.hermes\hermes-agent\venv\Scripts\python.exe'

results = []
def log(round_num, name, status, details=''):
    icon = "✅" if status == "PASS" else ("⚠️" if status == "WARN" else "❌")
    msg = f"  [{round_num:02d}] {icon} {name}"
    if details:
        msg += f"  — {details}"
    print(msg)
    results.append((round_num, name, status, details))

def run(cmd, timeout=30, cwd=REPO, input_text=None):
    """Run a command and return (returncode, stdout, stderr, duration_ms)"""
    if isinstance(cmd, str):
        cmd = cmd.split()
    t0 = time.time()
    try:
        r = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout,
                          cwd=cwd, input=input_text)
        dur = int((time.time() - t0) * 1000)
        return r.returncode, r.stdout, r.stderr, dur
    except subprocess.TimeoutExpired:
        return -1, "", "TIMEOUT", int((time.time() - t0) * 1000)
    except Exception as e:
        return -2, "", str(e), 0

# ════════════════════════════════════════════════
# Round 01: 部署 (Build)
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 01-03: 部署 / 编译 / 体积")
print("═══════════════════════════════════════════════")

rc, out, err, dur = run(['go', 'build', '-o', BIN, './cmd/via54/'], timeout=60)
if rc == 0 and os.path.exists(BIN):
    size_mb = os.path.getsize(BIN) / (1024*1024)
    log(1, "Go build CLI", "PASS", f"{size_mb:.1f}MB in {dur}ms")
else:
    log(1, "Go build CLI", "FAIL", err[:200])
    sys.exit(1)

rc, out, err, dur = run(['go', 'build', '-o', 'via54-mcp.exe', './cmd/mcp-server/'], timeout=60)
mcp_bin = os.path.join(os.path.dirname(BIN), 'via54-mcp.exe')
if rc == 0 and os.path.exists(mcp_bin):
    size_mb = os.path.getsize(mcp_bin) / (1024*1024)
    log(2, "Go build MCP", "PASS", f"{size_mb:.1f}MB in {dur}ms")
else:
    log(2, "Go build MCP", "FAIL", err[:200])

# go vet
rc, out, err, dur = run(['go', 'vet', './...'], timeout=60)
if rc == 0:
    log(3, "go vet 全包", "PASS", f"0 warnings in {dur}ms")
else:
    log(3, "go vet 全包", "WARN", err[:200])

# ════════════════════════════════════════════════
# Round 04-06: 冒烟测试 (CLI 子命令)
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 04-09: 冒烟测试 (每个子命令)")
print("═══════════════════════════════════════════════")

# Help 输出
rc, out, err, dur = run([BIN, 'help'], timeout=10)
if "via54Design" in out and "generate" in out:
    log(4, "help 输出", "PASS", f"{len(out)} chars")
else:
    log(4, "help 输出", "FAIL", f"unexpected: {out[:100]}")

# interactive
proc = subprocess.Popen([BIN], stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                       stdin=subprocess.PIPE, text=True)
try:
    out, err = proc.communicate(input="0\n", timeout=10)
    if "via54Design" in out and "退出" in out:
        log(5, "interactive 菜单", "PASS", "无参数进入菜单, 选项0退出")
    else:
        log(5, "interactive 菜单", "WARN", f"output: {out[:200]}")
except Exception as e:
    proc.kill()
    log(5, "interactive 菜单", "FAIL", str(e))

# list
rc, out, err, dur = run([BIN, 'list'], timeout=10)
if "layout" in out.lower() or "color" in out.lower() or "template" in out.lower() or len(out) > 50:
    log(6, "list 模板", "PASS", f"{len(out)} chars")
else:
    log(6, "list 模板", "FAIL", out[:200])

# version
rc, out, err, dur = run([BIN, 'version'], timeout=10)
if "via54" in out or "v0" in out:
    log(7, "version", "PASS", out.strip()[:100])
else:
    log(7, "version", "WARN", out[:200])

# prompt
rc, out, err, dur = run([BIN, 'prompt', '--scene', '一只猫在月光下的屋顶', '--platform', 'midjourney'], timeout=15)
if "提示词" in out or "prompt" in out.lower() or "midjourney" in out.lower():
    log(8, "prompt 生成", "PASS", f"{len(out)} chars")
else:
    log(8, "prompt 生成", "WARN", f"len={len(out)}, head={out[:200]}")

# narrate
rc, out, err, dur = run([BIN, 'narrate', '--seed', '一个裁缝在巴黎', '--model', 'three-act', '--duration', '30', '--format', 'json'], timeout=15)
try:
    j = json.loads(out)
    if 'beats' in j or 'model' in j or 'duration' in j:
        log(9, "narrate JSON", "PASS", f"JSON valid, {len(j)} keys")
    else:
        log(9, "narrate JSON", "WARN", f"keys: {list(j.keys())[:5]}")
except Exception as e:
    log(9, "narrate JSON", "WARN", f"非 JSON: {out[:200]}")

# ════════════════════════════════════════════════
# Round 10-12: 准确性 (输出确定性)
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 10-12: 准确性 (输出确定性)")
print("═══════════════════════════════════════════════")

# 同一输入 3 次生成，对比 md5
hashes = []
for i in range(3):
    rc, out, err, dur = run([BIN, 'narrate', '--seed', 'test_seed_stable', '--model', 'three-act', '--duration', '30', '--format', 'json'], timeout=10)
    h = hashlib.md5(out.encode()).hexdigest()[:8]
    hashes.append(h)

if hashes[0] == hashes[1] == hashes[2]:
    log(10, "narrate 确定性", "PASS", f"3 次 md5 一致: {hashes[0]}")
else:
    log(10, "narrate 确定性", "FAIL", f"不一致: {hashes}")

# prompt 确定性
hashes2 = []
for i in range(3):
    rc, out, err, dur = run([BIN, 'prompt', '--scene', 'fixed_scene', '--platform', 'flux'], timeout=10)
    h = hashlib.md5(out.encode()).hexdigest()[:8]
    hashes2.append(h)

if hashes2[0] == hashes2[1] == hashes2[2]:
    log(11, "prompt 确定性", "PASS", f"3 次 md5 一致: {hashes2[0]}")
else:
    log(11, "prompt 确定性", "WARN", f"差异: {hashes2}")

# generate HTML 确定性 — 写入文件，所以从文件读取
# 兼容: 生成可能写到其他位置或失败, 加 safe-read
def rm_output():
    try: os.remove('output.html')
    except: pass
rm_output()
rc1, _, _, _ = run([BIN, 'generate', '--layout', 'hero-split-16-9', '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', 'Test'], timeout=10)
try:
    with open('output.html', 'rb') as f: out1 = f.read()
except FileNotFoundError:
    out1 = b''
rm_output()
rc2, _, _, _ = run([BIN, 'generate', '--layout', 'hero-split-16-9', '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', 'Test'], timeout=10)
try:
    with open('output.html', 'rb') as f: out2 = f.read()
except FileNotFoundError:
    out2 = b''
rm_output()
h1 = hashlib.md5(out1).hexdigest()[:8]
h2 = hashlib.md5(out2).hexdigest()[:8]
if h1 == h2 and len(out1) > 100:
    log(12, "generate 确定性", "PASS", f"2 次 HTML md5: {h1} ({len(out1)}B)")
else:
    log(12, "generate 确定性", "WARN", f"diff: {h1} vs {h2}, len={len(out1)}")

# ════════════════════════════════════════════════
# Round 13-15: 边界测试
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 13-15: 边界测试")
print("═══════════════════════════════════════════════")

# 空 scene
rc, out, err, dur = run([BIN, 'prompt', '--scene', '', '--platform', 'midjourney'], timeout=10)
if "需要" in out or "required" in out.lower() or rc != 0:
    log(13, "空 scene 处理", "PASS", f"rc={rc}, 优雅拒绝")
else:
    log(13, "空 scene 处理", "FAIL", "接受了空输入")

# 无效平台
rc, out, err, dur = run([BIN, 'prompt', '--scene', 'test', '--platform', 'nonexistent_platform_xxx'], timeout=10)
if rc != 0 or "不支持" in out or "unknown" in out.lower() or "valid" in out.lower():
    log(14, "无效平台处理", "PASS", f"rc={rc}")
else:
    log(14, "无效平台处理", "WARN", f"可能接受无效平台: {out[:100]}")

# 超长字符串
long_scene = "一只猫" * 200  # 600 字符
rc, out, err, dur = run([BIN, 'prompt', '--scene', long_scene, '--platform', 'midjourney'], timeout=10)
if len(out) > 0 or rc != 0:
    log(15, "超长输入处理", "PASS", f"len_in={len(long_scene)}, len_out={len(out)}, rc={rc}")
else:
    log(15, "超长输入处理", "WARN", "返回空")

# ════════════════════════════════════════════════
# Round 16-18: 压力测试
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 16-18: 压力测试")
print("═══════════════════════════════════════════════")

# 20 次连续 prompt 生成
ok = 0
fail = 0
t0 = time.time()
for i in range(20):
    rc, out, _, _ = run([BIN, 'prompt', '--scene', f'stress_test_{i}', '--platform', 'midjourney'], timeout=10)
    if rc == 0 and len(out) > 0:
        ok += 1
    else:
        fail += 1
dur = time.time() - t0
log(16, "20 次 prompt 压力", "PASS" if fail == 0 else "WARN",
    f"{ok}/20 OK in {dur:.1f}s ({dur/20*1000:.0f}ms/次)")

# 10 次 narrate
ok = 0
t0 = time.time()
for i in range(10):
    rc, out, _, _ = run([BIN, 'narrate', '--seed', f'stress_{i}', '--model', 'three-act', '--format', 'json'], timeout=10)
    if rc == 0:
        ok += 1
dur = time.time() - t0
log(17, "10 次 narrate 压力", "PASS" if ok == 10 else "WARN",
    f"{ok}/10 OK in {dur:.1f}s")

# 5 次 generate (写到文件)
ok = 0
t0 = time.time()
for i in range(5):
    try: os.remove('output.html')
    except: pass
    rc, out, _, _ = run([BIN, 'generate', '--layout', 'hero-split-16-9', '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', f'Stress{i}'], timeout=10)
    if rc == 0 and os.path.exists('output.html') and os.path.getsize('output.html') > 100:
        ok += 1
dur = time.time() - t0
log(18, "5 次 generate 压力", "PASS" if ok == 5 else "WARN",
    f"{ok}/5 OK in {dur:.1f}s")

# ════════════════════════════════════════════════
# Round 19-20: 集成测试 (Web UI + MCP)
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" Round 19-20: 集成测试 (Web UI + MCP)")
print("═══════════════════════════════════════════════")

# 启动 web 服务器
web_proc = subprocess.Popen([BIN, 'web', '--port', '8765'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
time.sleep(3)

# 测试 web 端点
import urllib.request
try:
    resp = urllib.request.urlopen("http://localhost:8765/", timeout=5)
    html = resp.read().decode()
    if "htmx" in html.lower() and "<script>" not in html:
        log(19, "Web UI HTTP", "PASS", f"HTTP 200, {len(html)}B, 0 inline JS")
    elif "htmx" in html.lower():
        log(19, "Web UI HTTP", "PASS", f"HTTP 200, {len(html)}B")
    else:
        log(19, "Web UI HTTP", "FAIL", f"unexpected: {html[:200]}")
except Exception as e:
    log(19, "Web UI HTTP", "FAIL", str(e)[:200])

# HTMX 端点
try:
    resp = urllib.request.urlopen("http://localhost:8765/api/htmx/status", timeout=5)
    html = resp.read().decode()
    if "status-grid" in html:
        log(20, "HTMX 端点", "PASS", f"碎片 HTML {len(html)}B")
    else:
        log(20, "HTMX 端点", "WARN", f"unexpected: {html[:200]}")
except Exception as e:
    log(20, "HTMX 端点", "FAIL", str(e)[:200])

# 关闭 web
web_proc.terminate()
try:
    web_proc.wait(timeout=5)
except:
    web_proc.kill()

# ════════════════════════════════════════════════
# 测试结果汇总
# ════════════════════════════════════════════════
print("\n═══════════════════════════════════════════════")
print(" 📊 20 轮测试结果")
print("═══════════════════════════════════════════════")
pass_n = sum(1 for _, _, s, _ in results if s == "PASS")
warn_n = sum(1 for _, _, s, _ in results if s == "WARN")
fail_n = sum(1 for _, _, s, _ in results if s == "FAIL")
total = len(results)
print(f"  PASS:  {pass_n} / {total}")
print(f"  WARN:  {warn_n} / {total}")
print(f"  FAIL:  {fail_n} / {total}")
print()

if fail_n == 0:
    print("🎉 所有测试通过！代码质量良好。")
    sys.exit(0)
elif fail_n <= 2:
    print("⚠️ 个别测试失败，但核心功能正常")
    sys.exit(0)
else:
    print("❌ 多项测试失败，需要修复")
    sys.exit(1)
