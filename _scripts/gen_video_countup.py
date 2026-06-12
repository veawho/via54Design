#!/usr/bin/env python
# gen_video_countup.py - v0.7.0: ffmpeg 后处理 count-up 数字动画
# 段 4 7s 区间 (16-23s) 叠 5 个 PNG 用 enable between 切换
# editly 生成 base mp4 (段 4 c5 终极 PNG) → ffmpeg 后处理
import os, sys, subprocess

# 段 4 时间: clip 3 结束 16s (3+5+10-2=16) → 23s (16+7)
SEG4_START = 16.0
SEG4_END = 23.0
SEG4_DURATION = SEG4_END - SEG4_START  # 7s
COUNT_DURATION = SEG4_DURATION / 5  # 1.4s

def apply_countup(in_mp4, out_mp4, lang):
    """5 个 PNG 用 enable between 切换, 段 4 真正 count-up 动画"""
    counts = ["0%", "25%", "50%", "75%", "95%+"]
    base_dir = "G:/agent/hermes/via54Design-v6/subtitles/countup"

    # 5 个 PNG 路径 (与 gen_subtitle_countup.py 输出一致: c5_95%p.png 实际文件名)
    pngs = [
        f"{base_dir}/sub_04_{lang}_c1_0%.png",
        f"{base_dir}/sub_04_{lang}_c2_25%.png",
        f"{base_dir}/sub_04_{lang}_c3_50%.png",
        f"{base_dir}/sub_04_{lang}_c4_75%.png",
        f"{base_dir}/sub_04_{lang}_c5_95%p.png",
    ]

    # 5 个 overlay inputs
    cmd = ["ffmpeg", "-y", "-i", in_mp4]
    for png in pngs:
        cmd += ["-i", png]

    # filter_complex: 链式叠加, 每帧基于上一帧
    last_v = "0:v"
    filter_parts = []
    for i, png in enumerate(pngs):
        start = SEG4_START + i * COUNT_DURATION
        end = start + COUNT_DURATION
        out_v = f"v{i+1}"
        filter_parts.append(f"[{last_v}][{i+1}:v]overlay=enable='between(t,{start},{end})':x=(W-w)/2:y=(H-h)-80[{out_v}]")
        last_v = out_v

    cmd += [
        "-filter_complex", ";\n".join(filter_parts),
        "-map", f"[{last_v}]",
        "-map", "0:a?",
        "-c:a", "copy",
        "-c:v", "libx264",
        "-preset", "fast",
        "-crf", "23",
        out_mp4
    ]

    print(f"  ▶ ffmpeg count-up {lang}: {os.path.basename(in_mp4)} → {os.path.basename(out_mp4)}")
    print(f"    5 PNG: {os.path.basename(pngs[0])} ... {os.path.basename(pngs[-1])}")
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=120)
    if result.returncode != 0:
        print(f"    ✗ FAILED: {result.stderr[-500:]}")
        return False
    print(f"    ✓ OK")
    return True

if __name__ == "__main__":
    import sys
    lang = sys.argv[1] if len(sys.argv) > 1 else "zh"

    base_dir = "G:/agent/hermes/via54Design-v6/output"
    in_mp4 = f"{base_dir}/lithium_30s_v6_{lang}_base.mp4"
    out_mp4 = f"{base_dir}/lithium_30s_v6_{lang}_countup.mp4"
    final_mp4 = f"{base_dir}/lithium_30s_v6_{lang}.mp4"

    print(f"=== v0.7.0 ffmpeg count-up {lang} ===")
    if apply_countup(in_mp4, out_mp4, lang):
        # 替换 final
        os.replace(out_mp4, final_mp4)
        # 删 base
        if os.path.exists(in_mp4):
            os.remove(in_mp4)
        print(f"  ✓ {final_mp4}")
