#!/usr/bin/env python3
"""via54Design Web API并发/吞吐量压力测试
测试维度: 并发稳定性 / 线程安全 / API响应延迟统计
"""
import subprocess
import time
import sys
import os
import json
import urllib.request
import urllib.error
import threading

REPO = os.path.dirname(os.path.abspath(__file__))
BIN = os.path.join(REPO, 'via54_test.exe')
if os.path.exists(os.path.join(REPO, 'via54.exe')):
    BIN = os.path.join(REPO, 'via54.exe')

PORT = "8765"
BASE_URL = f"http://localhost:{PORT}"

# 用于保存线程测试结果
results = []
results_lock = threading.Lock()

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
        dur = (time.time() - t0) * 1000
        return e.code, str(e), dur
    except Exception as e:
        dur = (time.time() - t0) * 1000
        return -1, str(e), dur

def worker(worker_id, num_requests):
    for i in range(num_requests):
        # 轮流测试三个主要端点
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

        with results_lock:
            results.append({
                "worker_id": worker_id,
                "request_num": i,
                "endpoint": endpoint,
                "status": status,
                "success": success,
                "duration": dur,
                "error": None if success else str(res)[:100]
            })

print("=========================================")
print(" 🚀 开始 API 并发压力测试 (Concurrency Test)")
print("=========================================")

# 1. 确保 CLI/Web 二进制已编译
if not os.path.exists(BIN):
    print("编译测试二进制...")
    r = subprocess.run(['go', 'build', '-o', BIN, './cmd/via54/'], cwd=REPO)
    if r.returncode != 0:
        print("❌ 编译失败")
        sys.exit(1)

# 2. 启动 Web 服务
print(f"启动 Web 服务器，端口: {PORT}...")
web_proc = subprocess.Popen([BIN, 'serve'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=REPO)
# serve/web 都可以启动服务，在serve命令下或者web命令下启动
# 检查 cmd/serve.go 或者 cmd/web.go
# 让我们安全启动 via54 web --port 8765 或者是 via54 serve 或者是 via54 web
# 刚才的 test_20_rounds.py 用的是 BIN, 'web', '--port', '8765'，我们沿用该格式。
web_proc.terminate() # 先关闭上面启动的 serve
web_proc = subprocess.Popen([BIN, 'web', '--port', PORT], stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=REPO)

time.sleep(3) # 等待服务器完全就绪

# 3. 运行并发线程
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

# 4. 关闭 Web 服务器
print("正在关闭 Web 服务器...")
web_proc.terminate()
try:
    web_proc.wait(timeout=5)
except:
    web_proc.kill()

# 5. 分析结果
total_requests = len(results)
success_requests = sum(1 for r in results if r["success"])
failures = total_requests - success_requests
durations = [r["duration"] for r in results if r["success"]]

print("\n=========================================")
print(" 📊 并发测试数据汇总 (Concurrency Summary)")
print("=========================================")
print(f"  总并发请求数: {total_requests}")
print(f"  成功请求数:   {success_requests}")
print(f"  失败请求数:   {failures}")
print(f"  总测试时间:   {total_time:.2f} s")
print(f"  吞吐率 (RPS):  {total_requests / total_time:.2f} req/s")

if durations:
    import statistics
    print(f"  延迟统计:")
    print(f"    - 平均值 (Avg): {statistics.mean(durations):.2f} ms")
    print(f"    - 最小值 (Min): {min(durations):.2f} ms")
    print(f"    - 最大值 (Max): {max(durations):.2f} ms")
    print(f"    - 中位数 (Med): {statistics.median(durations):.2f} ms")
else:
    print("  ❌ 没有成功的请求以统计延迟。")

if failures > 0:
    print("\n  ⚠️ 失败请求日志示例:")
    failed_samples = [r for r in results if not r["success"]][:5]
    for idx, fs in enumerate(failed_samples):
        print(f"    [{idx+1}] Endp: {fs['endpoint']}, Status: {fs['status']}, Err: {fs['error']}")
    sys.exit(1)
else:
    print("\n🎉 并发压力测试完全通过！服务器运行非常稳定且线程安全。")
    sys.exit(0)
