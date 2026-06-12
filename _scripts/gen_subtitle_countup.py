#!/usr/bin/env python
# gen_subtitle_countup.py - v0.6.9: 段 4 count-up 数字动画
# 生成 5 张 PNG (0%→25%→50%→75%→95%+), 段 4 持续 7s, 5 帧 × 1.4s 切换
# 用 editly 多 image 层 + 切换时间错开
import os
from PIL import Image, ImageDraw, ImageFont

W, H = 1920, 1080

FONT_CANDIDATES = [
    r"C:\Windows\Fonts\msyhbd.ttc",
    r"C:\Windows\Fonts\msyh.ttc",
    r"C:\Windows\Fonts\arialbd.ttf",
    r"C:\Windows\Fonts\YuGothR.ttc",
]

FONT = next(f for f in FONT_CANDIDATES if os.path.exists(f))

# 段 4 5 帧数字进度 (0%→25%→50%→75%→95%+)
# 每帧 1.4s, 总 7s
COUNTS = [
    ("0%",    "00.0s-01.4s"),
    ("25%",   "01.4s-02.8s"),
    ("50%",   "02.8s-04.2s"),
    ("75%",   "04.2s-05.6s"),
    ("95%+",  "05.6s-07.0s"),
]

# 段 4 三语
PREFIX = {
    "zh": "智能产线稼动率",
    "en": "Smart Line Uptime",
    "ja": "スマートライン 稼働率",
}

def gen_countup_png(count: str, prefix: str, out_path: str, font_size: int = 70, padding: int = 32, radius: int = 20):
    """生成 count-up PNG: 黑底圆角 + "智能产线稼动率 95%+" 整段字幕"""
    img = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 字号: 前缀 56pt, 数字 96pt (大字更醒目)
    font_prefix = ImageFont.truetype(FONT, 56)
    font_count = ImageFont.truetype(FONT, 96)

    text = f"{prefix} {count}"

    bbox = draw.textbbox((0, 0), text, font=font_count)
    text_w = bbox[2] - bbox[0]
    text_h = bbox[3] - bbox[1]

    box_w = text_w + padding * 2
    box_h = text_h + padding * 2

    edge_margin = 80
    box_x = edge_margin
    box_y = H - box_h - edge_margin

    # 画圆角矩形
    draw.rounded_rectangle(
        [(box_x, box_y), (box_x + box_w, box_y + box_h)],
        radius=radius,
        fill=(0, 0, 0, 220),
    )

    # 画白字 (居中)
    text_x = box_x + box_w // 2
    text_y = box_y + box_h // 2
    draw.text(
        (text_x, text_y),
        text,
        font=font_count,
        fill=(255, 255, 255, 255),
        anchor="mm",
    )

    img.save(out_path, "PNG")
    print(f"  ✓ {out_path} ({box_w}x{box_h} @ {box_x},{box_y}, count={count})")
    return out_path

if __name__ == "__main__":
    OUT_DIR = "G:/agent/hermes/via54Design-v6/subtitles/countup"
    os.makedirs(OUT_DIR, exist_ok=True)

    print("=== v0.6.9 段 4 count-up PNG: 3 语 × 5 帧 = 15 张 ===\n")
    total = 0
    for lang in ["zh", "en", "ja"]:
        prefix = PREFIX[lang]
        for i, (count, _time) in enumerate(COUNTS):
            out = os.path.join(OUT_DIR, f"sub_04_{lang}_c{i+1}_{count.replace('+', 'p')}.png")
            gen_countup_png(count, prefix, out)
            total += 1
    print(f"\n=== 完成 {total} 张 count-up PNG ===")
    print("\n时间分配 (段 4 持续 7s):")
    for count, time in COUNTS:
        print(f"  {count}: {time}")
