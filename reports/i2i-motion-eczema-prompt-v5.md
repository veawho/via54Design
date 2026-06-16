# i2i + 人物动感 + 左右衣物一致 - Prompt V5

**新参考图**: `composer_2026-06-12_09-26-02-266_5e0807.png`（之前 i2i 生成的写实版皮肤科海报）
**V4 → V5 修正**: 
- ❌ 去掉所有运动场景（健身房/瑜伽/舞蹈/跑步/sports photography）
- ✅ 把"动感"集中在**人物本身**（头发飘、衣物微风、转身姿态），**场景不变**
- ✅ **左右衣物统一**：两边都穿同款白色吊带背心

---

## 🎯 3 个核心变化点

| # | 变化点 | 实现策略 |
|---|--------|---------|
| 1 | **i2i (图生图)** | `mmx image generate --subject-ref type=character,image=<参考图>` |
| 2 | **人物动感**（不是场景动感） | 头发被微风吹起 + 衣物飘动 + 身体微转/侧头 + 局部 motion blur（发梢、衣角）|
| 3 | **左右衣物一致** | 两人都穿**白色吊带背心**（同原图左人物的款式），不再一边白一边蓝 |

---

## 🎨 V5 Prompt（修正版）

```
Photorealistic 16:9 dynamic skincare banner. Sony A7IV with portrait lens, soft natural studio lighting with subtle breeze. Editorial fashion photography aesthetic with gentle motion energy on the subjects themselves.

LEFT Asian young woman with subtle artistic eczema redness on cheeks in medical aesthetic style — clusters of small soft pink-red dots in artistic arrangement, medically accurate but visually graceful, NOT raw or graphic. Her skin still has natural Asian warmth. Hair in soft natural movement as if a gentle breeze is blowing through the studio — a few strands lifting off the shoulder, hair ends slightly lifted with subtle motion energy. Head turned very slightly to her right (3/4 view) with a hint of worried expression — eyebrows drawn mildly together, eyes looking slightly to the side, lips gently closed with subtle tension. She wears white sleeveless tank top with thin spaghetti straps, scoop neckline, fabric lightly catching breeze (subtle movement in the strap area), silver stud earrings, thin silver chain necklace with tiny pendant.

RIGHT Asian young woman same person type, same exact clothing — white sleeveless tank top with thin spaghetti straps identical to left. Skin is radiant glowing clear flawless no blemishes. Hair flowing with movement — more strands lifted, ends flowing with motion suggesting she just turned her head. Bright confident joyful smile, eyes slightly squinted in real joy, head tilted very slightly to her left (3/4 view opposite to left woman, creating mirror composition). Silver stud earrings identical to left, silver chain necklace identical.

Both women wearing IDENTICAL white sleeveless tank tops with thin spaghetti straps scoop neckline — NOT different colors. The two women are mirror twins in clothing, only skin condition and expression differ.

Subtle motion energy throughout: hair ends with very slight motion blur (capturing mid-breeze moment), fabric of tank tops showing micro-movement from air flow, but faces remain sharp and focused. Body language of both women conveys gentle dynamism without being athletic or sporty — soft fashion-editorial motion, not gym workout motion.

Background: same clean soft light blue gradient as the original scene, lighter in center, slightly darker at edges. NO water droplets anywhere. NO splash elements. NO wet surfaces. NO water texture. Clean dry minimalist blue gradient with very subtle soft light particles or subtle atmospheric depth only. No motion blur on background — background stays calm and still while subjects have gentle motion.

HAIR identical base shape on both (with motion overlay): Short pixie-bob hybrid layered tousled, ends just below ear at jawline 5-8cm, NOT touching shoulders, dark natural black, wispy uneven bangs. Motion adds: a few extra strands lifting, ends flowing slightly. Same base hairstyle — same person with motion energy.

Identical accessories: round soft feminine face, silver stud earrings catching light, silver chain pendant necklace, white sleeveless tank top with thin spaghetti straps scoop neckline (SAME on both), natural Asian skin tone.

Lighting: soft natural studio lighting with subtle rim light on hair edges highlighting motion. Cool blue tones overall but with warm skin lighting on subjects.

STRICT PROHIBITIONS: NO water droplets on background, NO splash elements, NO wet surfaces. NO Chinese Japanese English Korean characters, NO numbers, NO typography. NO watermarks logos signatures. NO buttons CTA shapes. NO blue decorative borders. NO fitness studio, NO gym equipment, NO yoga poses, NO dance poses, NO athletic movement, NO sports photography, NO workout scenes, NO exercise clothes. CENTER column empty for post-production text.
```

## 🎨 V5 Prompt 最终版（1478 chars，已验证通过 mmx 1500 chars 限制）

```
Photorealistic 16:9 skincare banner. Sony A7IV soft studio subtle breeze. Editorial fashion photography gentle motion energy on subjects.

Two Asian women mirror twin composition. Both wear IDENTICAL white sleeveless tank top thin straps scoop neckline (NOT different colors), silver earrings silver chain necklace, Asian skin.

LEFT: subtle artistic eczema redness on cheeks — small soft pink-red dots in artistic arrangement, visually graceful NOT raw. Hair soft movement — strands lifting ends lifted by breeze. Head 3/4 view right with mild worried expression, eyebrows drawn, eyes side, lips closed.

RIGHT: skin radiant flawless. Hair flowing movement, strands lifted suggesting she just turned head. Bright confident joyful smile, eyes squinted joy, head tilted 3/4 view left (mirror).

Motion: hair ends slight motion blur mid-breeze, fabric micro-movement, faces sharp. Soft fashion-editorial motion — NOT athletic gym workout.

Background: clean soft light blue gradient lighter center darker edges. NO water droplets splash wet surfaces water texture. Clean dry minimalist blue gradient subtle atmospheric depth. Background calm still while subjects gentle motion.

Lighting: soft studio subtle rim light on hair edges. Cool blue tones warm skin.

STRICT: NO water droplets splash wet surfaces. NO Chinese/Japanese/English/Korean numbers typography watermarks logos. NO buttons CTA shapes blue borders. NO fitness gym yoga dance sports workout. CENTER empty for text.
```

---

## 📋 V4 → V5 关键变化对照

| # | V4（错误）| V5（修正）|
|---|---------|----------|
| 1 | "Sony A7IV with **fast shutter speed** capturing motion" | ✅ "Sony A7IV soft studio"（去掉运动摄影元素）|
| 2 | "**sports photography** aesthetic" | ✅ "**editorial fashion photography**"（去掉运动）|
| 3 | "Left Asian young woman **mid-stride in a fitness studio**" | ✅ "Head 3/4 view right with mild worried expression"（无运动）|
| 4 | "Right Asian young woman **mid-yoga flow or mid-dance turn**" | ✅ "head tilted 3/4 view left"（无运动）|
| 5 | "**Motion background: blurred modern fitness studio**" | ✅ "Background: clean soft light blue gradient... calm still"（静态背景）|
| 6 | "right wears soft **blue athletic tank top**" | ✅ "Both wear IDENTICAL **white sleeveless tank top**"（衣物统一）|
| 7 | 11 重 "**NO fitness gym yoga dance sports workout**" | ✅ 1 句 "**NO fitness gym yoga dance sports workout**" |

---

## ✅ V5 三大修正点

### 1. ❌ → ✅ 删除所有运动场景元素
- 删除：`fitness studio`、`yoga flow`、`dance turn`、`sports photography`、`workout scenes`、`gym equipment`、`athletic movement`
- 改为：editorial fashion photography、soft motion energy on subjects（仅人物本身）

### 2. ✅ 动感集中在人物本身（不是场景）
- 头发被微风吹起（strands lifting）
- 衣物微风感（fabric micro-movement）
- 局部 motion blur（hair ends only）
- 面部仍然 sharp（不模糊）
- 身体姿态：3/4 view 侧头（不是运动姿态）

### 3. ✅ 左右衣物统一（你的明确要求）
- V4 错误：左边白背心 + 右边蓝色背心
- V5 修正：**两边都穿 IDENTICAL 白色 sleeveless tank top**（保留原图款式）
- 描述方式：明示 "Both wear IDENTICAL white sleeveless tank top thin straps scoop neckline (NOT different colors)"