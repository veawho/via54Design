# 严格复刻原图（除文字外）- Prompt V2（修改版）

**原图**: `composer_2026-06-12_08-53-10-138_952f7c.jpg`（皮肤科湿疹海报）
**目标**: 严格复刻**人物状态** + **写实风** + **蓝底清爽** + 无文字 + 无蓝边框

---

## 🎨 V2 Prompt（修改后）

```
Photorealistic horizontal banner advertisement poster, 16:9 aspect ratio.
Shot on Sony A7IV with 85mm portrait lens, shallow depth of field,
high-end commercial skincare campaign photography aesthetic.

[STYLE & COLOR OVERRIDE - vs original sketch]
- Photorealistic style (NOT manga, NOT illustration, NOT anime, NOT line art)
- Color palette: clean blue fresh tone — light blue to pale blue gradient background, cool soft studio lighting
- Background covered with realistic water droplets of varying sizes with white specular highlights
- NO blue decorative borders on top/bottom edges — borderless clean frame
- Skin tones realistic and natural (NOT pale, NOT plastic), preserve realistic Asian skin tone with subtle warm undertone

[SYMMETRICAL DUAL-PORTRAIT LAYOUT]
- Two young East Asian women, both shown from chest up (bust shot), positioned on the LEFT and RIGHT thirds of the frame
- The CENTER vertical third is left as clean empty space for post-production typography overlay
- Both women are EXPLICITLY THE SAME PERSON — same face, same hairstyle, same outfit, same accessories, only skin condition and expression differ

[IDENTICAL DETAILS — both women must have EXACTLY these traits]
- Hair: short bob-cut, dark black with realistic individual strands, slightly wet-looking tips with water droplets clinging
- Face shape: round soft feminine face, delicate jawline
- Earrings: small simple silver stud earrings
- Clothing: white sleeveless tank top with thin spaghetti straps, simple scoop neckline
- Necklace: thin delicate silver chain with tiny pendant visible at collarbone
- Pose baseline: head at the same vertical height, shoulders visible, facing forward
- Both with fresh, dewy, slightly wet skin (as if just after washing face) — water droplets glistening on cheeks, neck, shoulders

[LEFT-SIDE WOMAN — the "before" / distressed]
- Expression: distressed and anxious, eyebrows drawn together in worried frown, eyes looking slightly upward with helplessness, lips slightly parted with subtle downward curve, mouth corners pulled down very slightly
- Skin condition: visible acne and blemishes on both cheeks — clusters of small red dots and slight inflammation concentrated on the LEFT cheek (her left, viewer's right side of her face) extending down to the jawline area, with some blemishes also on the right cheek
- Pose modifier: head tilted very slightly forward, shoulders raised and tensed, one shoulder slightly higher than the other
- Mood: helpless, frustrated by skin condition

[RIGHT-SIDE WOMAN — the "after" / happy]
- Expression: genuinely happy and confident, bright cheerful smile with teeth visible, eyes slightly squinted in real joy (Duchenne smile), eyebrows relaxed and naturally arched
- Skin condition: completely smooth, clear, flawless, radiant, glowing complexion — NO blemishes, NO redness, NO visible pores issues
- Pose modifier: head tilted very slightly back in quiet confidence, shoulders relaxed and down
- Mood: confident, joyful, radiant

[BACKGROUND]
- Smooth light blue to pale blue gradient, slightly darker at edges, lighter in center
- Covered with many realistic water droplets of varying sizes across the entire background
- Droplets have white specular highlights making them look three-dimensional and wet
- Water droplets denser in the upper portion and around the figures
- Soft cool ambient lighting, no harsh shadows

[STRICT PROHIBITIONS — ZERO TEXT IN IMAGE]
- NO Chinese characters
- NO Japanese characters
- NO English letters or text
- NO Korean characters
- NO numbers
- NO typography of any kind
- NO watermarks
- NO logos
- NO signatures
- NO symbols that resemble characters
- The CENTER vertical column of the image MUST remain completely empty and clean for later text overlay
- NO buttons or call-to-action shapes in the image
- NO emojis of any kind
- NO hand-pointing icons
- NO blue decorative borders on top or bottom edges
```

---

## 📋 V1 → V2 修改对照

| # | 修改点 | V1（草案） | V2（修改后） |
|---|--------|-----------|------------|
| 1 | **人物状态完全一致** | 部分覆盖（"same person"） | ✅ 完全明细：同发型 + 同脸型 + 同耳钉 + 同吊带 + 同项链 + 同湿发感 + 同湿润皮肤 |
| 2 | **风格改为写实** | "Japanese manga / shoujo comic style" | ✅ "Photorealistic... NOT manga, NOT illustration, NOT anime" + Sony A7IV |
| 3 | **色调改为蓝底清爽** | 黑白灰 + 蓝边框 | ✅ "clean blue fresh tone — light blue to pale blue gradient" |
| 4 | **去掉文字** | 部分覆盖 | ✅ 11 重否定（中/日/英/韩/数字/符号/水印/logo/签名/emoji/按钮） |
| 5 | **去掉蓝边框** | "horizontal BLUE decorative stripe... top/bottom" | ✅ "NO blue decorative borders on top/bottom edges" + "borderless clean frame" |

---

## 🎯 V2 Prompt 元素 vs 原图对照

### 保留原图（人物状态）

| 元素 | V2 是否覆盖 | 备注 |
|------|------------|------|
| 短发亚洲女性 × 2 | ✅ | "short bob-cut, dark black... young East Asian women" |
| 同款白色吊带 | ✅ | "white sleeveless tank top with thin spaghetti straps" 两次 |
| 同款小耳钉 | ✅ | "small simple silver stud earrings" |
| 同款项链 | ✅ | "thin delicate silver chain with tiny pendant" |
| 左：焦虑表情 | ✅ | "worried frown, eyes looking slightly upward... lips slightly parted with downward curve" |
| 左：脸颊痘痘 | ✅ | "clusters of small red dots... LEFT cheek extending to jawline... some on right cheek" |
| 右：灿烂微笑 | ✅ | "bright cheerful smile with teeth visible... Duchenne smile" |
| 右：皮肤光滑 | ✅ | "completely smooth, clear, flawless, radiant" |
| 水珠背景 | ✅ | "realistic water droplets... white specular highlights" |

### 替换原图（风格/色调）

| 元素 | V2 替换 |
|------|--------|
| 漫画风 | → 写实风（photorealistic + Sony A7IV） |
| 黑白灰 | → 蓝底清爽（light blue gradient + 写实肤色） |
| 蓝色上下边框 | → **无边框**（borderless clean frame） |
| 中央中文文字 | → **空白**（11 重否定 + 显式留白说明） |
| 底部按钮 + 👆 | → **无**（NO buttons, NO emojis, NO call-to-action） |

---

## ⚙️ 生成参数

```bash
via54 img --scene "<V2 prompt 全文>" \
  --platform minimax \
  --ar 16:9 \
  --n 4 \
  --prefix strict-v2
```

预计耗时 ~30s，输出 4 张 1280×720 JPEG 到 `/tmp/via54-strict-v2-<ts>/`，然后复制到 `reports/strict-v2_*.jpg`。

---

## ❓ 再次确认

请检查以下 4 点是否符合你的要求：

1. ✅ **人物状态完全一致**：prompt 里第 18-26 行（同发型/同脸型/同耳钉/同吊带/同项链/同湿发）
2. ✅ **写实风格**：第 1-3 行（Photorealistic + Sony A7IV + 写实摄影）
3. ✅ **蓝底清爽**：第 4-9 行（light blue gradient + 写实肤色）
4. ✅ **去文字**：第 39-49 行（11 重否定）
5. ✅ **去蓝边框**：第 9-10 行（NO blue decorative borders + borderless clean frame）

**确认就生成**（4 张 n=4，~30s）。  
**还需调整**就告诉我哪里改。
