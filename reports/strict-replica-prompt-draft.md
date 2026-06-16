# 严格复刻原图（除文字外）- Prompt 草案

**原图**: `composer_2026-06-12_08-53-10-138_952f7c.jpg`（皮肤科湿疹海报，黑白漫画风）
**目标**: 严格复刻所有视觉元素，**完全去除中文文字**（中央留空白给后期）

---

## 🎨 严格复刻 Prompt（待你确认）

```
Horizontal banner advertisement poster in Japanese manga / shoujo comic style, 16:9 aspect ratio, black and white ink illustration with light grayscale shading.

SYMMETRICAL DUAL-PORTRAIT LAYOUT:
- Two young East Asian women, each shown from chest up (bust shot), positioned on the LEFT and RIGHT thirds of the frame
- The CENTER vertical third is intentionally left as a clean empty column for post-production typography overlay (NO text, NO characters, NO symbols in this area)
- Visual focus: side-by-side before/after comparison of the SAME person

LEFT-SIDE WOMAN (the "before" — distressed):
- Hair: short bob-cut, slightly messy / windblown strands falling forward, dark black, hair has visible individual strands drawn with ink lines
- Face: round soft face, narrow eyes looking upward, eyebrows drawn together in a worried frown, lips slightly parted with subtle downward curve
- Skin problem: visible acne and blemishes on both cheeks — clusters of small dots and slight redness concentrated on the left cheek and jawline area
- Expression: distressed, anxious, helpless — conveying skin frustration
- Accessories: a small simple stud earring in the visible ear
- Clothing: white sleeveless tank top with thin spaghetti straps, simple scoop neckline
- Pose: head tilted very slightly forward, shoulders raised in tension, one shoulder slightly higher than the other
- A delicate thin necklace chain visible at the collarbone
- Position: occupying the left ~35% of the frame, face roughly at vertical center

RIGHT-SIDE WOMAN (the "after" — happy):
- Hair: same short bob-cut style as left woman, but hair appears smoother and more groomed
- Face: identical facial structure to left woman (same person), but with bright cheerful smile, eyes slightly squinted in genuine happiness, eyebrows relaxed and slightly arched
- Skin problem: NONE — smooth, clear, flawless complexion
- Expression: confident, joyful, radiant
- Accessories: same small stud earring
- Clothing: same white sleeveless tank top
- Pose: head tilted very slightly back in confidence, shoulders relaxed and down
- Same delicate necklace
- Position: occupying the right ~35% of the frame, mirroring left woman's vertical position

BACKGROUND:
- Smooth grayscale gradient, darker at the top and edges, slightly lighter in the center
- Covered with many water droplets of varying sizes scattered across the entire background
- The water droplets have white highlights making them look three-dimensional
- Droplets are denser and more visible in the upper portion and around the figures

BORDER DECORATION:
- A horizontal BLUE decorative stripe running across the very top edge of the frame (thin, about 3-5% of frame height)
- A matching horizontal BLUE decorative stripe running across the very bottom edge of the frame (same thickness)

LINEWORK STYLE:
- Clean manga ink linework, BLACK outline lines on all subjects (medium line weight, varying thickness)
- Light gray shading using subtle cross-hatching and dot patterns (screentone) to add dimension to hair and skin shadows
- No color anywhere on the figures — strictly black/white/gray
- High contrast between the white skin/clothing and black hair

BOTTOM CENTER (below the women's chest area):
- A small rounded-rectangle white button shape with a thin blue border, centered horizontally
- Inside the button: a small emoji-style hand pointing finger icon 👆 on the right side
- The button is small (occupies maybe 10% of frame width), positioned at the bottom center

OVERALL MOOD:
- Before/after transformation narrative
- The contrast between the two women creates emotional impact
- Clean, commercial-quality Japanese manga illustration aesthetic
- Suitable for skincare / dermatology advertising

STRICT PROHIBITIONS:
- NO Chinese characters, NO Japanese characters, NO English text, NO letters, NO numbers, NO typography anywhere in the image
- NO watermarks, NO logos, NO signatures
- The CENTER column of the image must remain completely empty / clean for later text overlay
```

---

## 📋 Prompt 元素对照表（vs 原图）

| 元素 | Prompt 是否覆盖 | 备注 |
|------|---------------|------|
| 16:9 横版 | ✅ | "16:9 aspect ratio, banner" |
| 左右双人物对称 | ✅ | "SYMMETRICAL DUAL-PORTRAIT LAYOUT" + 详细左/右描述 |
| 中央留白 | ✅ | "CENTER vertical third... clean empty column for post-production typography" |
| 短发亚洲女性 | ✅ | "short bob-cut, East Asian" |
| 左：焦虑表情 | ✅ | "worried frown, lips slightly parted... distressed" |
| 左：脸颊痘痘 | ✅ | "clusters of small dots and slight redness... left cheek and jawline" |
| 右：灿烂微笑 | ✅ | "bright cheerful smile, eyes slightly squinted in genuine happiness" |
| 右：皮肤光滑 | ✅ | "smooth, clear, flawless complexion" |
| 同款白色吊带 | ✅ | "white sleeveless tank top with thin spaghetti straps" 两次 |
| 同款小耳钉 | ✅ | "small simple stud earring" 两次 |
| 同款项链 | ✅ | "delicate thin necklace chain" |
| 水珠背景 | ✅ | "water droplets of varying sizes... white highlights" |
| 黑白灰漫画 | ✅ | "black and white ink illustration with light grayscale shading" |
| 蓝色上下边框 | ✅ | "horizontal BLUE decorative stripe... top/bottom edge" |
| 底部按钮 | ✅ | "small rounded-rectangle white button... blue border... hand pointing finger" |
| 中文字 | ❌ **完全无** | "NO Chinese characters, NO Japanese characters, NO English text" |
| 中央文字 | ❌ **完全无** | "CENTER column... remain completely empty" |
| 网点头发阴影 | ✅ | "screentone dot patterns" |
| 写实 vs 漫画 | 🎨 **漫画** | "Japanese manga / shoujo comic style" |

---

## 🎛️ 风格选项（等你确认）

**目前 prompt 是「漫画风」**（和原图风格 1:1 对齐）。

如果你想换成：
- **写实版**（真人摄影感）→ 我替换风格关键词
- **蓝底清爽版**（蓝色背景）→ 我把灰度背景改成蓝色

**但严格复刻 = 默认保留原图的「黑白漫画」风格**。

---

## ⚙️ 生成参数（计划用）

| 参数 | 值 | 说明 |
|------|----|----|
| `--platform` | minimax | 中文 prompt 直通最好 |
| `--ar` | 16:9 | 匹配原图宽高比 |
| `--n` | 4 | 跑 4 张选最佳 |
| `--seed` | -1 (随机) | 先看多样性 |
| `--prefix` | strict | 文件名前缀 |

预计耗时：~30s，输出 4 张 1280×720 JPEG。

---

## ❓ 等你确认

请确认以下 3 件事：

1. **风格**：保留**漫画风**（和原图一致）✅，还是改**写实风**？
2. **色调**：保留**黑白灰**（原图）✅，还是改**蓝底清爽**？
3. **prompt 草案**：上面那版 OK ✅ 还是要调整？

确认后我立即生图。生成完会再做 4 张相似度排序，把最像原图的那张保存到 reports/。
