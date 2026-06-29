#!/usr/bin/env python
# translate.py - 自动翻译模块 (中文 → 英文 + 日文, 备选 LibreTranslate)
# 用于 Pexels 关键词 / 旁白文案 多语生成

import json
import urllib.request
import urllib.parse
import sys

# 锂电 5 段关键词中英映射表 (★ 行业黄金词典 ★, 离线备)
KEYWORD_DICT = {
    # 中文 : (English, Japanese)
    "工业钢铁": ("industrial+steel", "工業鉄鋼"),
    "工厂": ("factory+manufacturing", "工場"),
    "数据可视化": ("data+visualization+hologram", "データ可視化"),
    "实验室": ("laboratory+technician", "実験室"),
    "电池": ("battery+cell", "電池"),
    "股票": ("stock+market+candle", "株式"),
    "电动车": ("electric+vehicle+car", "電気自動車"),
    "日出": ("sunrise+golden+hour", "日の出"),
    "城市": ("city+skyline+modern", "都市"),
    "金融仪表": ("finance+dashboard+chart", "金融ダッシュボード"),
    "全息": ("hologram+technology", "ホログラム"),
    "技术专业": ("technician+precision", "技術専門"),
    "投资": ("investment+growth", "投資"),
    "未来": ("future+innovation", "未来"),
    "新能源": ("renewable+energy+green", "新エネルギー"),
}

# 5 段默认中英 + 日文关键词
SEGMENTS_I18N = {
    "01_hook":   {"zh": "工业钢铁+工厂+火焰",          "en": "factory+industrial+steel+fire",      "ja": "工場+鉄鋼+火炎"},
    "02_trend":  {"zh": "数据可视化+全息+仪表盘",      "en": "data+visualization+dashboard+hologram", "ja": "データ可視化+ホログラム+ダッシュボード"},
    "03_tech":   {"zh": "实验室+技术员+电池",          "en": "technician+laboratory+battery",        "ja": "実験室+技術者+電池"},
    "04_market": {"zh": "股票+市场+蜡烛+金融",        "en": "stock+market+candle+finance",         "ja": "株式+市場+ローソク+金融"},
    "05_outlook":{"zh": "电动车+公路+日出",            "en": "electric+car+road+sunrise",            "ja": "電気自動車+道路+日の出"},
}

# 5 段中文旁白 + 英文/日文 (★ 30s 文案 5 段 ★)
NARRATION_I18N = {
    "01_hook": {
        "zh": "千度熔炼, 钢铁般的决心, 锂电新纪元由此开启。",
        "en": "A thousand degrees of forging, with the will of steel, a new era of lithium batteries begins.",
        "ja": "千度の精錬、鋼鉄の決意、リチウム電池新時代の幕開け。",
    },
    "02_trend": {
        "zh": "2026年产能突破八百吉瓦时, 中国锂电, 领跑全球。",
        "en": "In 2026, production capacity surpasses 800 GWh. China's lithium battery industry leads the world.",
        "ja": "2026年、生産能力が800ギガワット時を突破、中国のリチウム電池が世界をリード。",
    },
    "03_tech": {
        "zh": "固态电池技术领跑, 安全性能全面提升, 新能源汽车的未来, 就在此刻。",
        "en": "Solid-state battery technology leads the way, safety performance fully upgraded, the future of new energy vehicles is now.",
        "ja": "全固体電池技術がリード、安全性能が全面的に向上、新エネルギー自動車の未来は今ここにある。",
    },
    "04_market": {
        "zh": "板块资金流入同比增长两倍, 锂电投资窗口正在打开。",
        "en": "Sector capital inflow surged 200% year over year, the lithium battery investment window is opening.",
        "ja": "セクター資金流入は前年比2倍に急増、リチウム電池投資の窓が開いている。",
    },
    "05_outlook": {
        "zh": "从中国制造, 到中国智造, 锂电产业, 未来已来。",
        "en": "From Made in China to Smart in China, the lithium battery industry, the future is here.",
        "ja": "中国製造から中国スマート製造へ、リチウム電池産業、未来はもう来た。",
    },
}


def translate_offline(zh_text: str, target: str = "en") -> str:
    """离线查表翻译 (★ 黄金词典 ★, 0 网络延迟)"""
    # 简单分词查表 + fallback
    if target == "en":
        result = []
        # 按空格 / 中文标点分割
        import re
        for word in re.split(r'[+\s,，、]', zh_text):
            word = word.strip()
            if not word:
                continue
            if word in KEYWORD_DICT:
                result.append(KEYWORD_DICT[word][0].replace("+", " "))
            else:
                # 留原文
                result.append(word)
        return "+".join(result)
    elif target == "ja":
        result = []
        import re
        for word in re.split(r'[+\s,，、]', zh_text):
            word = word.strip()
            if not word:
                continue
            if word in KEYWORD_DICT:
                result.append(KEYWORD_DICT[word][1])
            else:
                result.append(word)
        return "+".join(result)
    return zh_text


def translate_online(zh_text: str, target: str = "en") -> str:
    """在线翻译 (LibreTranslate 公开 API, 兜底)"""
    try:
        url = "https://libretranslate.de/translate"
        data = urllib.parse.urlencode({
            "q": zh_text,
            "source": "zh",
            "target": target,
            "format": "text",
        }).encode()
        req = urllib.request.Request(url, data=data, headers={"Content-Type": "application/x-www-form-urlencoded"}, method="POST")
        with urllib.request.urlopen(req, timeout=10) as resp:
            return json.loads(resp.read())["translatedText"]
    except Exception as e:
        # fallback 离线
        return translate_offline(zh_text, target)


if __name__ == "__main__":
    # CLI: python translate.py <zh_text> [target=en|ja]
    if len(sys.argv) < 2:
        # 默认导出 5 段 i18n
        print(json.dumps({
            "segments": SEGMENTS_I18N,
            "narrations": NARRATION_I18N,
        }, ensure_ascii=False, indent=2))
    else:
        text = sys.argv[1]
        target = sys.argv[2] if len(sys.argv) > 2 else "en"
        result = translate_offline(text, target)
        print(f"[{target}] {text} → {result}")
