#!/usr/bin/env python3
"""via54Design WebUI 变更检测与自测脚本 (用于 Hermes cron)

由 cron job 定时触发 (每 10 分钟)。
检测 web/ cmd/ internal/ templates/workflows/ 的文件变更。
有变更时自动编译 + 启动server + 测试 + 报告。

输出格式: 仅在有变更或错误时输出 (silent when clean)
"""

import hashlib, json, os, subprocess, sys, time, http.client, glob

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
EXE = os.path.join(REPO, "via54.exe")
CHECKSUM_FILE = os.path.join(REPO, ".webui_checksums.json")
PORT = 19995  # internal test port, won't conflict with user's server

WATCH_PATTERNS = [
    "web/**",
    "cmd/**/*.go",
    "internal/**/*.go",
    "templates/workflows/**",
    "hack/*.py",
]

def get_checksums():
    checksums = {}
    for pattern in WATCH_PATTERNS:
        for f in glob.glob(os.path.join(REPO, pattern), recursive=True):
            if os.path.isfile(f) and not f.endswith('.pyc'):
                try:
                    checksums[f] = hashlib.md5(open(f, "rb").read()).hexdigest()
                except: pass
    return checksums

def load_saved():
    try:
        with open(CHECKSUM_FILE) as f:
            return json.load(f)
    except: return {}

def save_checksums(cs):
    with open(CHECKSUM_FILE, 'w') as f:
        json.dump(cs, f, indent=2)

def detect_changes(old, new):
    changed = []
    for f, h in new.items():
        if old.get(f) != h:
            changed.append(f)
    # Also detect deleted files
    for f in old:
        if f not in new:
            changed.append(f)
    return changed

def test_api(port):
    """Quick API health check"""
    try:
        conn = http.client.HTTPConnection("localhost", port, timeout=5)
        conn.request("GET", "/api/health")
        r = conn.getresponse()
        data = r.read().decode()
        conn.close()
        d = json.loads(data)
        return d.get("status") == "ok"
    except: return False

def full_test():
    """Compile + start server + quick test + stop"""
    # Compile
    r = subprocess.run(["go", "build", "-o", "via54.exe", "./cmd/via54/"],
                       cwd=REPO, capture_output=True, text=True, timeout=60)
    if r.returncode != 0:
        print(f"❌ BUILD FAILED:\n{r.stderr[:500]}")
        return False
    
    # Start server
    server = subprocess.Popen([EXE, "web", "--port", str(PORT)], cwd=REPO,
                               stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    time.sleep(3)
    
    try:
        # Test key endpoints
        endpoints = [
            ("GET", "/api/health"),
            ("GET", "/api/templates"),
        ]
        for method, path in endpoints:
            try:
                conn = http.client.HTTPConnection("localhost", PORT, timeout=5)
                conn.request(method, path)
                r = conn.getresponse()
                r.read()
                conn.close()
                if r.status != 200:
                    print(f"❌ {method} {path} -> {r.status}")
                    return False
            except Exception as e:
                print(f"❌ {method} {path} -> {e}")
                return False
        
        # Quick stress sample
        errors = 0
        for _ in range(20):
            try:
                conn = http.client.HTTPConnection("localhost", PORT, timeout=3)
                conn.request("GET", "/api/health")
                r = conn.getresponse()
                r.read()
                conn.close()
                if r.status != 200: errors += 1
            except: errors += 1
        
        if errors > 0:
            print(f"❌ Stress test: {errors} errors")
            return False
        
        print(f"✅ All tests passed (20/20 health checks)")
        return True
    finally:
        server.terminate()
        server.wait()

# ─── Main ───
if __name__ == "__main__":
    old = load_saved()
    new = get_checksums()
    changed = detect_changes(old, new)
    
    if not changed:
        sys.exit(0)  # Silent exit — nothing changed
    
    print(f"🔄 检测到 {len(changed)} 个文件变更:")
    for f in changed[:10]:
        rel = os.path.relpath(f, REPO)
        print(f"   • {rel}")
    if len(changed) > 10:
        print(f"   ... 及 {len(changed)-10} 个其他文件")
    print()
    
    ok = full_test()
    
    if ok:
        save_checksums(new)
        print(f"\n✅ 全部通过 — checksums 已更新")
        sys.exit(0)
    else:
        print(f"\n❌ 测试失败 — 请检查输出")
        sys.exit(1)
