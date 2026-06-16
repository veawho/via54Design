# 严格复刻原图（除文字外）- Prompt V3（头发严格一致）

**原图**: `composer_2026-06-12_08-53-10-138_952f7c.jpg`（皮肤科湿疹海报）
**V3 改进**: 头发描述严格对照原图，**两侧头发必须完全一致**（除"凌乱 vs 整洁"的小差异）

---

## 🎨 V3 Prompt（仅显示 V2 → V3 改动部分，完整版在生成时使用）

```
[HEADER - 保留 V2]
Photorealistic horizontal banner advertisement poster, 16:9 aspect ratio.
Shot on Sony A7IV with 85mm portrait lens, shallow depth of field,
high-end commercial skincare campaign photography aesthetic.
Photorealistic style (NOT manga, NOT illustration, NOT anime, NOT line art).
Color palette: clean blue fresh tone — light blue to pale blue gradient background, cool soft studio lighting.
Background covered with realistic water droplets of varying sizes with white specular highlights.
NO blue decorative borders on top/bottom edges — borderless clean frame.
Skin tones realistic and natural (NOT pale, NOT plastic), preserve realistic Asian skin tone with subtle warm undertone.

[NEW V3 — HAIR DESCRIPTION — STRICT AND DETAILED]

Hair is the SAME on both women (same hairstyle, same length, same parting, same color, same volume) — only the GROOMED STATE differs slightly.

HAIR — IDENTICAL DETAILS ON BOTH WOMEN:
1. STYLE: Short pixie-bob hybrid cut (NOT a typical bob, NOT pixie) — a layered tousled cut that ends JUST BELOW the ear lobe, sitting at jawline-to-neck-top level. Length approximately 5-8 cm below the ear.
2. HAIR LENGTH: 
   - Back and sides: ends JUST below the jawline, reaching the upper neck area (NOT touching shoulders)
   - Top: slightly longer strands that fall forward onto the forehead and around the temples
3. PARTING: NO clear parting line — hair falls naturally forward and to the sides, somewhat disheveled and windblown. Hair sweeps across and partially covers the forehead with wispy uneven strands.
4. FRONT/BANGS: 
   - There ARE wispy uneven bangs, but they are NOT a uniform blunt fringe
   - Bangs consist of multiple individual strands of varying lengths falling onto the forehead and between the eyes
   - Some strands are LONGER reaching down to the eyebrow level, some are SHORTER just above the eyebrows
   - Bangs are parted loosely in the middle (slight middle part visible) with strands falling on both sides of the forehead
5. SIDE PROFILE:
   - Hair tucks behind the ear on the OUTER side (visible ear side) — the ear is EXPOSED and visible
   - Hair flares out slightly at the bottom edges, with strands flicking outward
6. VOLUME & TEXTURE:
   - Slightly tousled, windblown, layered look
   - NOT flat against the head — has visible volume and lift at the crown and around the temples
   - Individual strands visible — NOT a solid block of hair
   - Strands have natural movement and asymmetry
7. COLOR: 
   - Dark natural black (NOT jet black plastic-looking, NOT brown)
   - With subtle realistic dark brown undertones in highlights
   - Natural shine from the soft cool lighting, NO colored highlights
8. HAIRLINE: 
   - Natural hairline, slightly visible at the temples
   - Some baby hairs along the hairline for realism

LEFT-SIDE WOMAN HAIR (the "before"):
- Hair appears slightly MORE tousled and messy, with more visible flyaway strands
- A few strands stick out slightly from the main silhouette
- The overall hair silhouette has slight asymmetry — left side flares out more
- Hair looks slightly damp at the tips (subtle wet-look)
- Some strands across the forehead are clearly separated and visible against the skin

RIGHT-SIDE WOMAN HAIR (the "after"):
- Same exact cut, length, parting, color as left woman
- Hair appears slightly MORE groomed and smooth — fewer flyaways
- Slightly more symmetrical silhouette than left
- Same slight dampness at tips (subtle wet-look, consistent with the fresh-washed theme)
- Same wispy bangs, same middle part, same forehead coverage

[REST - 保留 V2: 面部 / 表情 / 服饰 / 配饰 / 姿势 / 背景 / 禁止项]

[STRICT PROHIBITIONS - 保留 V2 的 11 重否定]
NO Chinese characters, NO Japanese characters, NO English letters, NO Korean characters, NO numbers, NO typography, NO watermarks, NO logos, NO signatures, NO emojis, NO buttons or call-to-action shapes, NO blue decorative borders on top/bottom edges.
```

---

## 📋 头发对照表（vs 原图）

| 元素 | 原图 | V3 Prompt 描述 |
|------|------|---------------|
| 长度 | 耳下到脖子上 | ✅ "JUST below jawline... upper neck area... NOT touching shoulders" |
| 风格 | 短碎发 + 层次 | ✅ "Short pixie-bob hybrid... layered tousled" |
| 刘海 | 不规则碎刘海 | ✅ "wispy uneven bangs... NOT uniform blunt fringe... multiple individual strands of varying lengths" |
| 分缝 | 中分偏松 | ✅ "loosely in the middle... slight middle part visible" |
| 蓬松度 | 有体积感 | ✅ "NOT flat... visible volume and lift at crown and around temples" |
| 耳朵 | 露出可见 | ✅ "tucks behind ear on OUTER side... ear EXPOSED and visible" |
| 发色 | 黑色 | ✅ "Dark natural black... subtle realistic dark brown undertones" |
| 单根发丝 | 可见 | ✅ "Individual strands visible — NOT a solid block" |
| 风吹感 | 有 | ✅ "tousled, windblown... some strands stick out slightly" |
| 左 vs 右差异 | 左更凌乱、右更整齐 | ✅ "MORE tousled and messy" vs "MORE groomed and smooth" |
| 湿润感 | 略带湿 | ✅ "slightly damp at the tips (subtle wet-look)" |

---

## 🔄 V2 → V3 主要改动

| 部分 | V2 描述 | V3 描述 |
|------|---------|---------|
| 头发长度 | "short bob-cut, dark black" | ✅ 精确到 "pixie-bob hybrid... 5-8 cm below ear... NOT touching shoulders" |
| 刘海 | 未明确 | ✅ 新增整段 "FRONT/BANGS"（碎刘海细节） |
| 分缝 | 未明确 | ✅ 新增 "NO clear parting... slight middle part visible" |
| 耳朵处理 | 未明确 | ✅ 新增 "tucks behind ear on OUTER side... ear EXPOSED" |
| 蓬松度 | "slightly wet-looking tips" | ✅ 新增 "NOT flat against head... volume at crown" |
| 单根发丝 | "individual strands" | ✅ 强化 "NOT a solid block of hair" |
| 左 vs 右差异 | 只说 "more groomed" | ✅ 6 项精确对比（messy/groomed/symmetry/damp/forehead strands） |
| 整体一致性 | "same hairstyle" | ✅ 8 项 "HAIR — IDENTICAL DETAILS ON BOTH WOMEN" 明示 |

---

## ⚙️ 生成参数

```bash
via54 img --scene "<V3 prompt 全文>" \
  --platform minimax \
  --ar 16:9 \
  --n 4 \
  --prefix strict-v3
```

预计 ~30s 出 4 张 1280×720 JPEG。

---

## ❓ 再次确认

V3 重点强化了头发一致性。请检查：

1. ✅ **短发类型**: "pixie-bob hybrid... JUST below ear... 5-8 cm"
2. ✅ **刘海细节**: "wispy uneven bangs... NOT uniform blunt fringe"
3. ✅ **中分偏松**: "loosely in the middle... slight middle part visible"
4. ✅ **耳朵露出**: "tucks behind ear... ear EXPOSED and visible"
5. ✅ **蓬松度**: "NOT flat... volume at crown and around temples"
6. ✅ **发色**: "Dark natural black... subtle dark brown undertones"
7. ✅ **左 vs 右差异**: messy vs groomed（仅整理状态不同）
8. ✅ **湿润感**: "slightly damp at tips (subtle wet-look)"

**确认就生成。** 还有调整告诉我。
