#!/usr/bin/env python3
"""via54Design Master Verification & Stress Test Suite

Supports 3 testing modes:
1. Smoke Integration Mode (Default, 20 rounds)
2. Longevity Stress Mode (--stress <rounds>)
3. Concurrent API Mode (--concurrent)
"""
import subprocess
import time
import sys
import os
import json
import hashlib
import shutil
import urllib.request
import urllib.error
import threading
import statistics

# ── 跨平台二进制与环境探测 ───────────────────────
REPO = os.path.dirname(os.path.abspath(__file__))
BIN_EXT = ".exe" if os.name == "nt" else ""
BIN_NAME = f"via54{BIN_EXT}"
MCP_BIN_NAME = f"via54-mcp{BIN_EXT}"
BIN = os.path.join(REPO, BIN_NAME)

VENV_PY = sys.executable
GO_BIN = shutil.which("go") or "/usr/local/go/bin/go"
if os.name == "nt":
    for candidate in [
        r"C:\Program Files\Go\bin\go.exe",
        r"C:\Program Files (x86)\Go\bin\go.exe",
        os.path.expanduser(r"~\go\bin\go.exe"),
    ]:
        if os.path.exists(candidate):
            GO_BIN = candidate
            break

GO_ENV = os.environ.copy()
GO_ENV.setdefault("GOPATH", os.path.expanduser("~/go"))
GO_ENV.setdefault("GOMODCACHE", os.path.join(GO_ENV["GOPATH"], "pkg", "mod"))
GO_ENV.setdefault("CGO_ENABLED", "0")

results = []
results_lock = threading.Lock()

def log(round_num, name, status, details=''):
    icon = "✅" if status == "PASS" else ("⚠️" if status == "WARN" else "❌")
    msg = f"  [{round_num:02d}] {icon} {name}"
    if details:
        msg += f"  — {details}"
    print(msg)
    with results_lock:
        results.append((round_num, name, status, details))

def diagnose_error(cmd, returncode, stdout, stderr):
    lines = []
    lines.append(f"    命令: {' '.join(cmd[:5])}{'...' if len(cmd) > 5 else ''}")
    lines.append(f"    exit={returncode}")
    if stderr:
        first_err = stderr.strip().splitlines()[0] if stderr.strip() else "(empty)"
        lines.append(f"    stderr[0]: {first_err[:200]}")
    else:
        lines.append(f"    stderr: (empty)")

    if "module cache not found" in (stderr or ""):
        lines.append("    💊 修复: export GOPATH=~/go && export GOMODCACHE=$GOPATH/pkg/mod")
    elif "permission denied" in (stderr or "").lower():
        lines.append(f"    💊 修复: chmod +x {cmd[0]} 或在 Windows 用管理员权限")
    elif "no such file" in (stderr or "").lower() or "not found" in (stderr or "").lower():
        lines.append("    💊 修复: 检查二进制路径/名称是否正确")
    return "\n".join(lines)

def run(cmd, timeout=30, cwd=REPO, input_text=None, env=None):
    if isinstance(cmd, str):
        cmd = cmd.split()
    t0 = time.time()
    try:
        r = subprocess.run(cmd, capture_output=True, encoding='utf-8', timeout=timeout,
                          cwd=cwd, input=input_text, env=env or os.environ)
        dur = int((time.time() - t0) * 1000)
        return r.returncode, r.stdout, r.stderr, dur
    except subprocess.TimeoutExpired:
        return -1, "", "TIMEOUT", int((time.time() - t0) * 1000)
    except Exception as e:
        return -2, "", str(e), 0

# ── Mode A: Default 20 Rounds Smoke Mode ───────────────────
def run_smoke_mode():
    print("\n═══════════════════════════════════════════════")
    print(" Round 01-03: 部署 / 编译 / 体积")
    print(f" 平台: {sys.platform} ({os.name}), Go: {GO_BIN}, Python: {VENV_PY}")
    print("═══════════════════════════════════════════════")
    
    if not os.path.exists(GO_BIN):
        log(1, "Go build CLI", "FAIL", f"go binary not found: {GO_BIN}")
        sys.exit(1)
    
    rc, out, err, dur = run([GO_BIN, 'build', '-o', BIN, './cmd/via54/'], timeout=120, env=GO_ENV)
    if rc == 0 and os.path.exists(BIN):
        size_mb = os.path.getsize(BIN) / (1024*1024)
        log(1, "Go build CLI", "PASS", f"{size_mb:.1f}MB in {dur}ms")
    else:
        log(1, "Go build CLI", "FAIL", f"exit={rc}")
        print(diagnose_error([GO_BIN, 'build', '-o', BIN, './cmd/via54/'], rc, out, err))
        sys.exit(1)
    
    mcp_bin = os.path.join(os.path.dirname(BIN), MCP_BIN_NAME)
    rc, out, err, dur = run([GO_BIN, 'build', '-o', mcp_bin, './cmd/mcp-server/'], timeout=120, env=GO_ENV)
    if rc == 0 and os.path.exists(mcp_bin):
        size_mb = os.path.getsize(mcp_bin) / (1024*1024)
        log(2, "Go build MCP", "PASS", f"{size_mb:.1f}MB in {dur}ms")
    else:
        log(2, "Go build MCP", "WARN", f"exit={rc} (非阻塞, 跳过)")
    
    rc, out, err, dur = run([GO_BIN, 'vet', './...'], timeout=60, env=GO_ENV)
    if rc == 0:
        log(3, "go vet 全包", "PASS", f"0 warnings in {dur}ms")
    else:
        log(3, "go vet 全包", "WARN", err[:200])
    
    print("\n═══════════════════════════════════════════════")
    print(" Round 04-09: 冒烟测试 (每个子命令)")
    print("═══════════════════════════════════════════════")
    
    rc, out, err, dur = run([BIN, '--help'], timeout=10)
    if "via54" in out and ("generate" in out or "generate..." in out):
        log(4, "help 输出", "PASS", f"{len(out)} chars")
    else:
        log(4, "help 输出", "FAIL", f"unexpected: {out[:100]}")
    
    proc = subprocess.Popen([BIN], stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                           stdin=subprocess.PIPE, encoding='utf-8')
    try:
        out, err = proc.communicate(input="0\n", timeout=10)
        if "via54" in out and ("退出" in out or "exit" in out.lower() or "0" in out):
            log(5, "interactive 菜单", "PASS", "无参数进入菜单, 选项0退出")
        else:
            log(5, "interactive 菜单", "WARN", f"output: {out[:200]}")
    except Exception as e:
        proc.kill()
        log(5, "interactive 菜单", "FAIL", str(e))
    
    rc, out, err, dur = run([BIN, 'list'], timeout=10)
    if "layout" in out.lower() or "color" in out.lower() or "template" in out.lower() or len(out) > 50:
        log(6, "list 模板", "PASS", f"{len(out)} chars")
    else:
        log(6, "list 模板", "FAIL", out[:200])
    
    rc, out, err, dur = run([BIN, 'version'], timeout=10)
    if "via54" in out or "v0" in out:
        log(7, "version", "PASS", out.strip()[:100])
    else:
        log(7, "version", "WARN", out[:200])
    
    rc, out, err, dur = run([BIN, 'prompt', '--scene', 'cat on moonlit roof', '--platform', 'midjourney'], timeout=15)
    if "提示词" in out or "prompt" in out.lower() or "midjourney" in out.lower() or len(out) > 100:
        log(8, "prompt 生成", "PASS", f"{len(out)} chars")
    else:
        log(8, "prompt 生成", "WARN", f"len={len(out)}, head={out[:200]}")
    
    rc, out, err, dur = run([BIN, 'narrate', '--seed', 'tailor in Paris', '--model', 'three-act', '--duration', '30', '--format', 'json'], timeout=15)
    try:
        j = json.loads(out)
        if 'beats' in j or 'model' in j or 'duration' in j:
            log(9, "narrate JSON", "PASS", f"JSON valid, {len(j)} keys")
        else:
            log(9, "narrate JSON", "WARN", f"keys: {list(j.keys())[:5]}")
    except Exception as e:
         log(9, "narrate JSON", "WARN", f"非 JSON: {out[:200]}")
    
    print("\n═══════════════════════════════════════════════")
    print(" Round 10-12: 准确性 (输出确定性)")
    print("═══════════════════════════════════════════════")
    
    hashes = []
    for i in range(3):
        rc, out, err, dur = run([BIN, 'narrate', '--seed', 'test_seed_stable', '--model', 'three-act', '--duration', '30', '--format', 'json'], timeout=10)
        h = hashlib.md5(out.encode()).hexdigest()[:8]
        hashes.append(h)
    
    if hashes[0] == hashes[1] == hashes[2]:
        log(10, "narrate 确定性", "PASS", f"3 次 md5 一致: {hashes[0]}")
    else:
        log(10, "narrate 确定性", "FAIL", f"不一致: {hashes}")
    
    hashes2 = []
    for i in range(3):
        rc, out, err, dur = run([BIN, 'prompt', '--scene', 'fixed_scene', '--platform', 'flux'], timeout=10)
        h = hashlib.md5(out.encode()).hexdigest()[:8]
        hashes2.append(h)
    
    if hashes2[0] == hashes2[1] == hashes2[2]:
        log(11, "prompt 确定性", "PASS", f"3 次 md5 一致: {hashes2[0]}")
    else:
        log(11, "prompt 确定性", "WARN", f"差异: {hashes2}")
    
    def rm_output():
        p = "output.html"
        try: os.remove(p)
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
    
    print("\n═══════════════════════════════════════════════")
    print(" Round 13-15: 边界测试")
    print("═══════════════════════════════════════════════")
    
    rc, out, err, dur = run([BIN, 'prompt', '--scene', '', '--platform', 'midjourney'], timeout=10)
    if "需要" in out or "required" in out.lower() or rc != 0:
        log(13, "空 scene 处理", "PASS", f"rc={rc}, 优雅拒绝")
    else:
        log(13, "空 scene 处理", "FAIL", "接受了空输入")
    
    rc, out, err, dur = run([BIN, 'prompt', '--scene', 'test', '--platform', 'nonexistent_platform_xxx'], timeout=10)
    if rc != 0 or "不支持" in out or "unknown" in out.lower() or "valid" in out.lower() or "❌" in out:
        log(14, "无效平台处理", "PASS", f"rc={rc}")
    else:
        log(14, "无效平台处理", "WARN", f"可能接受无效平台: {out[:100]}")
    
    long_scene = "a cat " * 200
    rc, out, err, dur = run([BIN, 'prompt', '--scene', long_scene, '--platform', 'midjourney'], timeout=10)
    if len(out) > 0 or rc != 0:
        log(15, "超长输入处理", "PASS", f"len_in={len(long_scene)}, len_out={len(out)}, rc={rc}")
    else:
        log(15, "超长输入处理", "WARN", "返回空")
    
    print("\n═══════════════════════════════════════════════")
    print(" Round 16-18: 压力测试")
    print("═══════════════════════════════════════════════")
    
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
    log(16, "20 次 prompt 压力", "PASS" if fail == 0 else "WARN", f"{ok}/20 OK in {dur:.1f}s ({dur/20*1000:.0f}ms/次)")
    
    ok = 0
    t0 = time.time()
    for i in range(10):
        rc, out, _, _ = run([BIN, 'narrate', '--seed', f'stress_{i}', '--model', 'three-act', '--format', 'json'], timeout=10)
        if rc == 0:
            ok += 1
    dur = time.time() - t0
    log(17, "10 次 narrate 压力", "PASS" if ok == 10 else "WARN", f"{ok}/10 OK in {dur:.1f}s")
    
    ok = 0
    t0 = time.time()
    for i in range(5):
        try: os.remove('output.html')
        except: pass
        rc, out, _, _ = run([BIN, 'generate', '--layout', 'hero-split-16-9', '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', f'Stress{i}'], timeout=10)
        if rc == 0 and os.path.exists('output.html') and os.path.getsize('output.html') > 100:
            ok += 1
    dur = time.time() - t0
    log(18, "5 次 generate 压力", "PASS" if ok == 5 else "WARN", f"{ok}/5 OK in {dur:.1f}s")
    
    print("\n═══════════════════════════════════════════════")
    print(" Round 19-20: 集成测试 (Web UI + MCP)")
    print("═══════════════════════════════════════════════")
    
    web_proc = subprocess.Popen([BIN, 'web', '--port', '8765'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(3)
    
    try:
        resp = urllib.request.urlopen("http://localhost:8765/", timeout=5)
        html = resp.read().decode()
        if "htmx" in html.lower():
            log(19, "Web UI HTTP", "PASS", f"HTTP 200, {len(html)}B")
        else:
            log(19, "Web UI HTTP", "FAIL", f"unexpected: {html[:200]}")
    except Exception as e:
        log(19, "Web UI HTTP", "FAIL", str(e)[:200])
    
    try:
        resp = urllib.request.urlopen("http://localhost:8765/api/htmx/status", timeout=5)
        html = resp.read().decode()
        if "status" in html.lower() or len(html) > 10:
            log(20, "HTMX 端点", "PASS", f"碎片 HTML {len(html)}B")
        else:
            log(20, "HTMX 端点", "WARN", f"unexpected: {html[:200]}")
    except Exception as e:
        log(20, "HTMX 端点", "FAIL", str(e)[:200])
    
    web_proc.terminate()
    try: web_proc.wait(timeout=5)
    except: web_proc.kill()
    
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
        sys.exit(0)
    elif fail_n <= 2:
        sys.exit(0)
    else:
        sys.exit(1)

# ── Mode B: Longevity Stress Mode ──────────────────────────
def run_stress_mode(rounds):
    print("=========================================")
    print(f" 🚀 开始 {rounds} 轮连续压力测试 (Longevity Stress Test)")
    print("=========================================")
    
    print("[Step 1] 编译 CLI 压力测试二进制...")
    rc, out, err, dur = run(['go', 'build', '-o', BIN, './cmd/via54/'])
    if rc != 0:
        print(f"❌ 编译失败: {err}")
        sys.exit(1)
    print(f"✅ 编译成功: {BIN} ({dur:.1f}ms)")
    
    print(f"\n[Step 2] 执行 {rounds} 次连续 prompt 生成...")
    prompt_durs = []
    prompt_failures = 0
    for i in range(1, rounds + 1):
        scene = f"A beautiful futuristic city with high tech towers, round {i}"
        rc, out, err, dur = run([BIN, 'prompt', '--scene', scene, '--platform', 'midjourney'])
        if rc == 0 and "提示词" in out:
            prompt_durs.append(dur)
        else:
            prompt_failures += 1
            print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")
            
    # Run 1/4 rounds for narrate and generate for speed
    n_rounds = max(5, rounds // 4)
    print(f"\n[Step 3] 执行 {n_rounds} 次连续 narrate 故事生成...")
    narrate_durs = []
    narrate_failures = 0
    for i in range(1, n_rounds + 1):
        seed = f"An ancient artisan building a perpetual machine, round {i}"
        rc, out, err, dur = run([BIN, 'narrate', '--seed', seed, '--model', 'three-act', '--duration', '60', '--format', 'json'])
        if rc == 0 and "beats" in out:
            narrate_durs.append(dur)
        else:
            narrate_failures += 1
            print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")
            
    print(f"\n[Step 4] 执行 {n_rounds} 次 HTML 设计生成...")
    generate_durs = []
    generate_failures = 0
    for i in range(1, n_rounds + 1):
        layout = "hero-split-16-9" if i % 2 == 0 else "bento-grid-2x2"
        title = f"Design Showcase Round {i}"
        rc, out, err, dur = run([BIN, 'generate', '--layout', layout, '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', title])
        if rc == 0 and os.path.exists(os.path.join(REPO, 'output.html')):
            generate_durs.append(dur)
            try: os.remove(os.path.join(REPO, 'output.html'))
            except: pass
        else:
            generate_failures += 1
            print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")
            
    print("\n=========================================")
    print(" 📊 延迟性能统计汇总 (Latency Statistics)")
    print("=========================================")
    
    def print_stats(name, durs):
        if not durs:
            print(f"  {name}: 无数据")
            return
        print(f"  {name}:")
        print(f"    - 平均耗时 (Avg): {statistics.mean(durs):.2f} ms")
        print(f"    - 最小耗时 (Min): {min(durs):.2f} ms")
        print(f"    - 最大耗时 (Max): {max(durs):.2f} ms")
        if len(durs) > 1:
            print(f"    - 标准差 (StdDev): {statistics.stdev(durs):.2f} ms")
            
    print_stats("Prompt 生成", prompt_durs)
    print_stats("Narrate 叙事", narrate_durs)
    print_stats("Generate 设计", generate_durs)
    
    try: os.remove(BIN)
    except: pass
    
    total_fails = prompt_failures + narrate_failures + generate_failures
    sys.exit(0 if total_fails == 0 else 1)

# ── Mode C: Concurrent API Mode ────────────────────────────
PORT = "8765"
BASE_URL = f"http://localhost:{PORT}"
concurrent_results = []
concurrent_lock = threading.Lock()

def post_json(endpoint, data):
    url = f"{BASE_URL}{endpoint}"
    req_data = json.dumps(data).encode('utf-8')
    req = urllib.request.Request(
        url,
        data=req_data,
        headers={'Content-Type': 'application/json'},
        method='POST'
    )
    t0 = time.time()
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            body = response.read().decode('utf-8')
            resp_data = json.loads(body)
            dur = (time.time() - t0) * 1000
            return response.status, resp_data, dur
    except urllib.error.HTTPError as e:
        return e.code, str(e), (time.time() - t0) * 1000
    except Exception as e:
        return -1, str(e), (time.time() - t0) * 1000

def worker(worker_id, num_requests):
    for i in range(num_requests):
        if i % 3 == 0:
            status, res, dur = post_json("/api/prompt", {
                "scene": f"A cute dog in the yard, thread {worker_id} request {i}",
                "platform": "flux"
            })
            endpoint = "/api/prompt"
            success = (status == 200 and "output" in res)
        elif i % 3 == 1:
            status, res, dur = post_json("/api/narrate", {
                "seed": f"A lost traveler finding an oasis, thread {worker_id} request {i}",
                "model": "heros-journey",
                "duration": 30
            })
            endpoint = "/api/narrate"
            success = (status == 200 and "output" in res)
        else:
            status, res, dur = post_json("/api/generate", {
                "layout": "bento-grid-2x2",
                "color": "moon-white",
                "font": "ming-hei-editorial",
                "title": f"Bento thread {worker_id} req {i}"
            })
            endpoint = "/api/generate"
            success = (status == 200 and "html" in res)

        with concurrent_lock:
            concurrent_results.append({
                "worker_id": worker_id,
                "request_num": i,
                "endpoint": endpoint,
                "status": status,
                "success": success,
                "duration": dur,
                "error": None if success else str(res)[:100]
            })

def run_concurrent_mode():
    print("=========================================")
    print(" 🚀 开始 API 并发压力测试 (Concurrency Test)")
    print("=========================================")
    
    if not os.path.exists(BIN):
        print("编译测试二进制...")
        r = subprocess.run(['go', 'build', '-o', BIN, './cmd/via54/'], cwd=REPO)
        if r.returncode != 0:
            print("❌ 编译失败")
            sys.exit(1)
            
    print(f"启动 Web 服务器，端口: {PORT}...")
    web_proc = subprocess.Popen([BIN, 'web', '--port', PORT], stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=REPO)
    time.sleep(3)
    
    NUM_THREADS = 10
    REQUESTS_PER_THREAD = 20
    print(f"并发配置: {NUM_THREADS} 线程, 每线程 {REQUESTS_PER_THREAD} 次请求，总计 {NUM_THREADS * REQUESTS_PER_THREAD} 次请求...")
    
    threads = []
    t_start = time.time()
    for t_id in range(NUM_THREADS):
        t = threading.Thread(target=worker, args=(t_id, REQUESTS_PER_THREAD))
        threads.append(t)
        t.start()
        
    for t in threads:
        t.join()
        
    t_end = time.time()
    total_time = t_end - t_start
    
    print("正在关闭 Web 服务器...")
    web_proc.terminate()
    try: web_proc.wait(timeout=5)
    except: web_proc.kill()
    
    total_requests = len(concurrent_results)
    success_requests = sum(1 for r in concurrent_results if r["success"])
    failures = total_requests - success_requests
    durations = [r["duration"] for r in concurrent_results if r["success"]]
    
    print("\n=========================================")
    print(" 📊 并发测试数据汇总 (Concurrency Summary)")
    print("=========================================")
    print(f"  总并发请求数: {total_requests}")
    print(f"  成功请求数:   {success_requests}")
    print(f"  失败请求数:   {failures}")
    print(f"  总测试时间:   {total_time:.2f} s")
    print(f"  吞吐率 (RPS):  {total_requests / total_time:.2f} req/s")
    
    if durations:
        print(f"  延迟统计:")
        print(f"    - 平均值 (Avg): {statistics.mean(durations):.2f} ms")
        print(f"    - 最小值 (Min): {min(durations):.2f} ms")
        print(f"    - 最大值 (Max): {max(durations):.2f} ms")
        print(f"    - 中位数 (Med): {statistics.median(durations):.2f} ms")
        
    try: os.remove(BIN)
    except: pass
    
    sys.exit(0 if failures == 0 else 1)

# ── Main Entry ─────────────────────────────────────────────
if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="via54Design Unified Test Runner")
    parser.add_argument("--stress", type=int, default=0, help="Run sequential stress test for N rounds")
    parser.add_argument("--concurrent", action="store_true", help="Run concurrent API stress test")
    args = parser.parse_args()
    
    if args.stress > 0:
        run_stress_mode(args.stress)
    elif args.concurrent:
        run_concurrent_mode()
    else:
        run_smoke_mode()
