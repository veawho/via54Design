#!/usr/bin/env python
# gen_subtitle_png.py - 用 PIL 生成黑底圆角+白字 PNG, 给段 1 hook (避火) 用
# 输入: 字幕文本, 输出: PNG 透明背景黑底圆角白字
# editly 用 image 层 + 50% 透明度叠到火上
import os
from PIL import Image, ImageDraw, ImageFont

W, H = 1920, 1080  # 16:9 视频分辨率

# 字体路径: Windows 自带
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

# 段 1 字幕 (3 语, 25pt 黑底白字 + 12px 内边距 + 圆角 24)
SUBTITLES = {
    "zh": "锂电新纪元",
    "en": "New Era of Lithium",
    "ja": "リチウム新時代",
}

def gen_subtitle_png(text: str, out_path: str, font_size: int = 96, padding: int = 36, radius: int = 24, bg_color=(0, 0, 0, 220), text_color=(255, 255, 255, 255)):
    """
    生成 PNG 透明背景 + 黑底圆角 + 白字
    bg_color/text_color 第 4 位 = alpha
    """
    # 创建透明背景图层
    img = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 加载字体
    font = ImageFont.truetype(FONT, font_size)

    # 测量文字 bbox (用 getbbox 拿真实大小)
    bbox = draw.textbbox((0, 0), text, font=font)
    text_w = bbox[2] - bbox[0]
    text_h = bbox[3] - bbox[1]

    # 字幕框尺寸
    box_w = text_w + padding * 2
    box_h = text_h + padding * 2

    # 段 1 字幕位置: bottom-right (右下角, 避火)
    # 离右边 80px, 离下边 80px
    box_x = W - box_w - 80
    box_y = H - box_h - 80

    # 画圆角矩形 (PIL 2.0+ 直接支持 rounded_rectangle)
    draw.rounded_rectangle(
        [(box_x, box_y), (box_x + box_w, box_y + box_h)],
        radius=radius,
        fill=bg_color,
    )

    # 画白字 (居中)
    # PIL 2.0+ 用 anchor="mm" 居中
    text_x = box_x + box_w // 2
    text_y = box_y + box_h // 2
    draw.text(
        (text_x, text_y),
        text,
        font=font,
        fill=text_color,
        anchor="mm",
    )

    # 保存 PNG
    img.save(out_path, "PNG")
    print(f"  ✓ {out_path} ({box_w}x{box_h} @ {box_x},{box_y})")
    return out_path

if __name__ == "__main__":
    OUT_DIR = "G:/agent/hermes/via54Design-v6/subtitles"
    os.makedirs(OUT_DIR, exist_ok=True)

    print("=== 生成段 1 字幕 PNG (3 语, bottom-right 黑底圆角白字) ===")
    for lang, text in SUBTITLES.items():
        out = os.path.join(OUT_DIR, f"sub_01_{lang}.png")
        gen_subtitle_png(text, out, font_size=96)

    print("\n=== 段 5 outlook PNG (3 语, bottom 居中黑底圆角白字) ===")
    OUTLOOK = {
        "zh": "投资窗口正在打开",
        "en": "Investment Window Opening",
        "ja": "投資の窓が開く",
    }
    for lang, text in OUTLOOK.items():
        out = os.path.join(OUT_DIR, f"sub_05_{lang}.png")
        gen_subtitle_png(text, out, font_size=80)
