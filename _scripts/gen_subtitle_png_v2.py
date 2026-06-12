#!/usr/bin/env python
# gen_subtitle_png_v2.py - v0.6.8: 5 段 × 3 语 = 15 PNG, 字体 80pt 精致版
# 全程黑底圆角+白字, 段 4 改 "智能产线稼动率 95%+" 配 ESTUN 机器人
import os
from PIL import Image, ImageDraw, ImageFont

W, H = 1920, 1080  # 16:9 视频分辨率

FONT_CANDIDATES = [
    r"C:\Windows\Fonts\msyhbd.ttc",  # 微软雅黑 Bold (中文)
    r"C:\Windows\Fonts\msyh.ttc",    # 微软雅黑
    r"C:\Windows\Fonts\simhei.ttf",  # 黑体
    r"C:\Windows\Fonts\arialbd.ttf", # Arial Bold (英文)
    r"C:\Windows\Fonts\YuGothR.ttc", # 游ゴシック
]

def find_font():
    for f in FONT_CANDIDATES:
        if os.path.exists(f):
            return f
    raise FileNotFoundError("No font found")

FONT = find_font()

# v0.6.8 全 5 段字幕 (★ 段 4 改文案 ★)
# 段 1 hook - hook
# 段 2 trend - 趋势
# 段 3 tech - 技术 (在画面下方, 中等技术, font 70pt 略小)
# 段 4 market - 市场 (改文案 配机器人)
# 段 5 outlook - 展望
SUBTITLES = {
    "01": {  # 段 1 hook
        "zh": "锂电新纪元",
        "en": "New Era of Lithium",
        "ja": "リチウム新時代",
        "size": 96,  # hook 标题最大
        "pos": "bottom-right",  # 段 1 bottom-right 避火
    },
    "02": {  # 段 2 trend
        "zh": "2026 产能突破 800GWh",
        "en": "2026 Output Exceeds 800GWh",
        "ja": "2026年 生産能力 800GWh突破",
        "size": 70,  # trend 略小
        "pos": "bottom-left",  # 段 2 数据型 放左下 (不挡数据图)
    },
    "03": {  # 段 3 tech
        "zh": "固态电池技术领跑全球",
        "en": "Solid-State Tech Leads Global",
        "ja": "全固体電池技術 世界リード",
        "size": 70,
        "pos": "bottom-left",  # 段 3 极片特写 放左下 (不挡极片)
    },
    "04": {  # 段 4 market - ★ v0.6.8 改文案 ★
        "zh": "智能产线稼动率 95%+",  # ★ 改: 配机器人
        "en": "Smart Line Uptime 95%+",  # ★ 改
        "ja": "スマートライン 稼働率 95%+",  # ★ 改
        "size": 70,
        "pos": "bottom-left",  # 段 4 机器人 放左下
    },
    "05": {  # 段 5 outlook
        "zh": "投资窗口正在打开",
        "en": "Investment Window Opening",
        "ja": "投資の窓が開く",
        "size": 96,  # v0.6.9 收尾压轴: 80pt → 96pt 配 hook 段 1
        "pos": "bottom-center",  # 段 5 outlook 收尾居中
    },
}

def gen_subtitle_png(text: str, out_path: str, font_size: int = 80, padding: int = 32, radius: int = 20, bg_color=(0, 0, 0, 220), text_color=(255, 255, 255, 255), position: str = "bottom-left"):
    """
    生成 PNG 透明背景 + 黑底圆角 + 白字 (1920x1080)
    position: bottom-left / bottom-right / bottom-center (段内位置)
    """
    img = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    font = ImageFont.truetype(FONT, font_size)

    bbox = draw.textbbox((0, 0), text, font=font)
    text_w = bbox[2] - bbox[0]
    text_h = bbox[3] - bbox[1]

    box_w = text_w + padding * 2
    box_h = text_h + padding * 2

    # 按 position 计算 box 位置
    edge_margin = 80  # 离边距
    if position == "bottom-left":
        box_x = edge_margin
        box_y = H - box_h - edge_margin
    elif position == "bottom-right":
        box_x = W - box_w - edge_margin
        box_y = H - box_h - edge_margin
    elif position == "bottom-center":
        box_x = (W - box_w) // 2
        box_y = H - box_h - edge_margin
    else:
        box_x = (W - box_w) // 2
        box_y = (H - box_h) // 2

    # 画圆角矩形
    draw.rounded_rectangle(
        [(box_x, box_y), (box_x + box_w, box_y + box_h)],
        radius=radius,
        fill=bg_color,
    )

    # 画白字 (居中)
    text_x = box_x + box_w // 2
    text_y = box_y + box_h // 2
    draw.text(
        (text_x, text_y),
        text,
        font=font,
        fill=text_color,
        anchor="mm",
    )

    img.save(out_path, "PNG")
    print(f"  ✓ {out_path} ({box_w}x{box_h} @ {box_x},{box_y}, pos={position})")
    return out_path

if __name__ == "__main__":
    OUT_DIR = "G:/agent/hermes/via54Design-v6/subtitles"
    os.makedirs(OUT_DIR, exist_ok=True)

    print("=== v0.6.8 生成 5 段 × 3 语 = 15 PNG (80pt 字体, 段 4 改文案) ===\n")
    total = 0
    for seg_id, cfg in SUBTITLES.items():
        for lang in ["zh", "en", "ja"]:
            text = cfg[lang]
            out = os.path.join(OUT_DIR, f"sub_{seg_id}_{lang}.png")
            gen_subtitle_png(text, out, font_size=cfg["size"], position=cfg["pos"])
            total += 1
    print(f"\n=== 完成 {total} 个 PNG ===")
