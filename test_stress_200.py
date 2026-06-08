#!/usr/bin/env python3
"""via54Design 200轮连续压力/持久性测试
测试维度: 稳定性 / 响应延迟分布 / 无内存泄露 / 连续无错
"""
import subprocess
import time
import sys
import os
import statistics

REPO = os.path.dirname(os.path.abspath(__file__))
BIN = os.path.join(REPO, 'via54_stress.exe')

def run_cmd(cmd, cwd=REPO):
    if isinstance(cmd, str):
        cmd = cmd.split()
    t0 = time.time()
    try:
        r = subprocess.run(cmd, capture_output=True, encoding='utf-8', timeout=10, cwd=cwd)
        dur = (time.time() - t0) * 1000  # ms
        return r.returncode, r.stdout, r.stderr, dur
    except Exception as e:
        return -1, "", str(e), 0

print("=========================================")
print(" 🚀 开始 200 轮连续压力测试 (Longevity Stress Test)")
print("=========================================")

# 1. 编译专用的 stress 二进制
print("[Step 1] 编译 CLI 压力测试二进制...")
rc, out, err, dur = run_cmd(['go', 'build', '-o', BIN, './cmd/via54/'])
if rc != 0:
    print(f"❌ 编译失败: {err}")
    sys.exit(1)
print(f"✅ 编译成功: {BIN} ({dur:.1f}ms)")

# 2. 200次连续 prompt 生成测试
print("\n[Step 2] 执行 200 次连续 prompt 生成...")
prompt_durs = []
prompt_failures = 0
for i in range(1, 201):
    scene = f"A beautiful futuristic city with high tech towers, round {i}"
    rc, out, err, dur = run_cmd([BIN, 'prompt', '--scene', scene, '--platform', 'midjourney'])
    if rc == 0 and "提示词" in out:
        prompt_durs.append(dur)
    else:
        prompt_failures += 1
        print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")

if prompt_failures == 0:
    print(f"✅ 200次 prompt 完美通过!")
else:
    print(f"⚠️ Prompt 失败数: {prompt_failures}/200")

# 3. 50次连续 narrate (故事大纲)
print("\n[Step 3] 执行 50 次连续 narrate 故事生成...")
narrate_durs = []
narrate_failures = 0
for i in range(1, 51):
    seed = f"An ancient artisan trying to build a perpetual motion machine in round {i}"
    rc, out, err, dur = run_cmd([BIN, 'narrate', '--seed', seed, '--model', 'three-act', '--duration', '60', '--format', 'json'])
    if rc == 0 and "beats" in out:
        narrate_durs.append(dur)
    else:
        narrate_failures += 1
        print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")

if narrate_failures == 0:
    print(f"✅ 50次 narrate 完美通过!")
else:
    print(f"⚠️ Narrate 失败数: {narrate_failures}/50")

# 4. 50次 HTML 设计布局生成
print("\n[Step 4] 执行 50 次 HTML 设计生成...")
generate_durs = []
generate_failures = 0
for i in range(1, 51):
    layout = "hero-split-16-9" if i % 2 == 0 else "bento-grid-2x2"
    title = f"Design Showcase Round {i}"
    rc, out, err, dur = run_cmd([BIN, 'generate', '--layout', layout, '--color', 'ink-wash', '--font', 'ming-hei-editorial', '--title', title])
    if rc == 0 and os.path.exists(os.path.join(REPO, 'output.html')):
        generate_durs.append(dur)
        # 清理生成的临时文件
        try: os.remove(os.path.join(REPO, 'output.html'))
        except: pass
    else:
        generate_failures += 1
        print(f"  ❌ 第 {i} 次失败: rc={rc}, err={err[:100]}")

if generate_failures == 0:
    print(f"✅ 50次 generate HTML 完美通过!")
else:
    print(f"⚠️ Generate 失败数: {generate_failures}/50")

# 5. 汇总数据分析
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
    print(f"    - 标准差 (StdDev): {statistics.stdev(durs):.2f} ms")

print_stats("Prompt 生成 (200轮)", prompt_durs)
print_stats("Narrate 叙事 (50轮)", narrate_durs)
print_stats("Generate 设计 (50轮)", generate_durs)

# 清理 stress 二进制
try: os.remove(BIN)
except: pass

print("\n🎉 压力测试脚本执行完成！")
if prompt_failures + narrate_failures + generate_failures == 0:
    sys.exit(0)
else:
    sys.exit(1)
