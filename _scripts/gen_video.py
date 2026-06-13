#!/usr/bin/env python
# gen_video.py - 锂电 30s v6 一键生成器 (★ v0.6.4 ★)
# 用法: python gen_video.py [--lang zh|en|ja|all] [--step fetch|video|mux|all]

import json
import subprocess
import sys
import os
from pathlib import Path

ROOT = Path(r"G:\agent\hermes\via54Design-v6")
TEMPLATE = ROOT / "spec" / "lithium_30s_v6_template.json5"

# ★ 字幕黄金模板 (5 段 × 3 语言) ★
I18N = {
    "title_01": {
        "zh": "锂电新纪元",
        "en": "New Era of Lithium",
        "ja": "リチウム新時代",
    },
    "sub_02": {
        "zh": "2026 产能突破 800GWh",
        "en": "2026: 800GWh Milestone",
        "ja": "2026年 800GWh 突破",
    },
    "sub_03": {
        "zh": "固态电池技术领跑全球",
        "en": "Solid-State Tech Leads the World",
        "ja": "全固体電池 世界リード",
    },
    "sub_04": {
        "zh": "板块资金流入 同比增长 200%",
        "en": "Sector Inflow +200% YoY",
        "ja": "資金流入 前年比 200%",
    },
    "title_05": {
        "zh": "投资窗口正在打开",
        "en": "Investment Window Opens",
        "ja": "投資の窓が開く",
    },
}


def render_spec(lang: str) -> Path:
    """模板渲染: 替换 {title_xx} {sub_xx} {lang} 占位符 → 写具体 spec 文件"""
    tpl = TEMPLATE.read_text(encoding="utf-8")
    # 删除头部注释行 (JSON5 支持 // 但 # 不行, 保险删所有 // 注释)
    lines = []
    for line in tpl.splitlines():
        stripped = line.strip()
        if stripped.startswith("//"):
            continue
        lines.append(line)
    tpl = "\n".join(lines)
    # 替换 outPath
    rendered = tpl.replace("{lang}", lang)
    # 替换字幕
    for key, translations in I18N.items():
        rendered = rendered.replace("{" + key + "}", translations[lang])
    # 写文件
    out_spec = ROOT / "spec" / f"lithium_30s_v6_{lang}.json5"
    out_spec.write_text(rendered, encoding="utf-8")
    return out_spec


def run_editly(lang: str) -> Path:
    """跑 editly 拼接 (只视频, 无音轨)"""
    spec = render_spec(lang)
    node = r"C:\Users\via54\tools\node18\node-v18.20.4-win-x64\editly.exe"
    print(f"▶ editly --json {spec.name}")
    r = subprocess.run([node, "--json", str(spec)], capture_output=True, text=True, timeout=600)
    if r.returncode != 0:
        print(f"  ✗ editly failed: {r.stderr[-500:]}")
        return None
    print(f"  ✓ editly done")
    return ROOT / "output" / f"lithium_30s_v6_{lang}.mp4"  # editly 默认命名


def run_mux(lang: str) -> Path:
    """ffmpeg 三路混音 (旁白 + 配乐 + 视频)"""
    # 拼接 5 段旁白
    voice_list = [ROOT / "voice" / f"{lang}_{k}.mp3" for k in
                  ["01_hook", "02_trend", "03_tech", "04_market", "05_outlook"]]
    voice_tmp = ROOT / "output" / f"_voice_{lang}.mp3"
    # 写 concat 列表
    concat_file = ROOT / "output" / f"_concat_voice_{lang}.txt"
    concat_file.write_text(
        "\n".join([f"file '{v.as_posix()}'" for v in voice_list]),
        encoding="utf-8"
    )
    subprocess.run([
        "ffmpeg", "-y", "-f", "concat", "-safe", "0",
        "-i", str(concat_file), "-c", "copy", str(voice_tmp)
    ], check=True, capture_output=True)

    voice_dur = float(subprocess.run([
        "ffprobe", "-v", "error", "-show_entries", "format=duration",
        "-of", "csv=p=0", str(voice_tmp)
    ], capture_output=True, text=True).stdout.strip())
    print(f"  旁白总时长: {voice_dur:.2f}s")

    # v0.7.1: 时长归一化 (en 28.925s → 30s, 加 1.075s 静音补齐)
    # 用 apad 给旁白补静音到 30s, 让 mux 输出 3 语都 = 30.000s
    pad_dur = max(0, 30.0 - voice_dur)
    print(f"  末尾补静音: {pad_dur:.3f}s")

    # v0.7.0: 按语言用对应 count-up 视频源 (不用统一的 lithium_30s_v6.mp4)
    vid = ROOT / "output" / f"lithium_30s_v6_{lang}.mp4"  # 当前语言的 count-up 版
    if not vid.exists():
        # 退路: 用 lithium_30s_v6.mp4 (zh 版)
        vid = ROOT / "output" / "lithium_30s_v6.mp4"
    bgm = ROOT / "music" / "bgm_epic_30s.mp3"
    out = ROOT / "output" / f"lithium_30s_v6_{lang}.mp4"

    # v0.7.1: 旁白补静音 (apad) + 配乐 atrim 0:30 + 视频 -t 30 强制截到 30s
    filter_complex = (
        f"[1:a]volume=1.4,apad,atrim=0:30.0[voice];"
        f"[2:a]volume=0.13,atrim=0:30.0[bgm];"
        "[voice][bgm]amix=inputs=2:duration=first:dropout_transition=2[mix]"
    )
    r = subprocess.run([
        "ffmpeg", "-y",
        "-i", str(vid),
        "-i", str(voice_tmp),
        "-i", str(bgm),
        "-filter_complex", filter_complex,
        "-map", "0:v", "-map", "[mix]",
        "-t", "30",  # v0.7.1: 强制截到 30s (避免 en 28.925s 偏差)
        "-c:v", "copy", "-c:a", "aac", "-b:a", "192k",
        str(out) + ".tmp.mp4"  # v0.7.0: 临时文件, 避免原地覆盖
    ], capture_output=True, text=True)
    if r.returncode != 0:
        print(f"  ✗ mux failed: {r.stderr[-500:]}")
        return None
    # 移动临时文件到目标
    import shutil
    shutil.move(str(out) + ".tmp.mp4", str(out))
    return out


def main():
    args = sys.argv[1:]
    lang = "all"
    step = "all"
    for i, a in enumerate(args):
        if a == "--lang" and i + 1 < len(args):
            lang = args[i + 1]
        elif a == "--step" and i + 1 < len(args):
            step = args[i + 1]

    langs = ["zh", "en", "ja"] if lang == "all" else [lang]

    if step in ("fetch", "all"):
        print("=== Step 1: 拉 Pexels 视频 ===")
        r = subprocess.run([sys.executable, str(ROOT / "fetch_pexels_v3.py")])
        if r.returncode != 0:
            print("✗ fetch failed"); return

    if step in ("video", "all"):
        print("=== Step 2: 拼接视频 (editly, 无音轨) ===")
        node = "C:/Users/via54/tools/node18/node-v18.20.4-win-x64/editly.cmd"
        # video 步骤默认 3 语 (不论 --lang, 跑全 3 语)
        video_langs = ["zh", "en", "ja"]
        for L in video_langs:
            print(f"  --- {L} ---")
            spec_lang = render_spec(L)
            r = subprocess.run([node, "--json", str(spec_lang)], capture_output=True, text=True, timeout=600, shell=True)
            if r.returncode != 0:
                print(f"  ✗ editly {L} failed: {r.stderr[-500:]}"); continue
            print(f"    ✓ editly {L} → output/lithium_30s_v6_{L}.mp4")
        # 找最近的 zh 视频流作为 mux 源
        zh_video = ROOT / "output" / "lithium_30s_v6_zh.mp4"
        if zh_video.exists():
            target = ROOT / "output" / "lithium_30s_v6.mp4"
            import shutil
            shutil.copy(zh_video, target)
            print(f"  ✓ mux 源 (复制) {zh_video.name} → {target.name}")

    if step in ("voice", "all"):
        print("=== Step 3: 生成旁白 ===")
        subprocess.run(["bash", str(ROOT / "gen_voice.sh")], check=False)

    if step in ("mux", "all"):
        print("=== Step 4: 三语混音 ===")
        # mux 步骤也跑全 3 语 (不论 --lang)
        for L in ["zh", "en", "ja"]:
            print(f"--- {L} ---")
            out = run_mux(L)
            if out:
                # ffprobe 验证
                d = float(subprocess.run([
                    "ffprobe", "-v", "error", "-show_entries", "format=duration",
                    "-of", "csv=p=0", str(out)
                ], capture_output=True, text=True).stdout.strip())
                sz = out.stat().st_size / 1024 / 1024
                print(f"  ✓ {out.name}: {d:.2f}s / {sz:.1f}MB")

    print("\n=== 全部完成 ===")


if __name__ == "__main__":
    main()
