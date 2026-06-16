# 湿疹皮肤科海报 V5 — 修订 prompt 草案 (2026-06-13)

## 来源
- **原图**: `~/Desktop/ref_for_user.png` (皮肤科湿疹海报, manga 风, 中文"动不动就反复湿疹"+ 按钮"皮肤科专家来修救!"+ 蓝边框)
- **修订指令**: "缺少人物角度和人物动作，背景色调改为蓝色，风格为摄影写实风。"

## 修订要点 (用户反馈 → 修复)
| 上一轮问题 | 修复策略 | 修复位置 |
|-----------|---------|---------|
| ❌ 人物角度只写"head tilted slightly"（太弱） | ✅ 明确 **3/4 侧视角** + **eye-level 平视**，左右互为镜像 | COMPOSITION + 每个 LEFT/RIGHT 块 |
| ❌ 人物动作只写"shoulders raised"（不是真动作） | ✅ 明确具体动作：左手抚摸脸颊看镜 + 右手撩发回头笑 | 每个 LEFT/RIGHT 块 |
| ⚠️ 背景蓝色只写"light blue gradient"（不够） | ✅ 蓝色提升为主色基调：深蓝→天蓝渐变 + 蓝雾景深 | STYLE & COLOR 块 |
| ✅ 摄影写实风（已 OK，保留 V2 magic combo） | ✅ 强化 Sony A7IV + 85mm + skin texture | STYLE & COLOR 块 |

## 默认平台
**minimax (jimeng/image-01)** — Phase 2 决策：中文场景 + 写实 + 蓝底清爽 = minimax 直通中文最强，垂类细节保留最好。
如要换平台：MJ/Flux/Gemini/DALL·E3 告诉我一句即可。

---

## Prompt (V5 完整文本)

```
Horizontal banner skincare advertisement poster, 16:9 aspect ratio,
PHOTOREALISTIC style (NOT manga, NOT illustration, NOT anime, NOT line art),
shot on Sony A7IV with 85mm portrait lens, shallow depth of field,
high-end commercial skincare campaign photography aesthetic,
editorial fashion photography composition.

[STYLE & COLOR — BLUE AS DOMINANT TONE]
- Color palette: BLUE AS DOMINANT TONE — deep sapphire blue at frame edges
  fading to bright sky blue in center, soft volumetric blue mist atmosphere,
  subtle blue cast on skin highlights
- Background: smooth blue atmospheric gradient with soft depth-of-field
  bokeh, NO water droplets, NO splash, NO wet surfaces, NO texture patterns
  — just clean dreamy blue gradient with faint blue mist depth
- Skin tones realistic Asian, cool blue undertones preserved, NOT plastic,
  NOT pale, with natural pores and fine vellus hair visible
- Lighting: soft diffused studio rim light with subtle blue bounce,
  gentle warm key light on faces for skin readability
- NO blue decorative borders on top/bottom edges — borderless clean frame

[SYMMETRICAL DUAL-PORTRAIT LAYOUT — MIRROR COMPOSITION]
- Two young East Asian women, both shown from chest up
- Mirror composition: LEFT subject faces 3/4 view to camera right,
  RIGHT subject faces 3/4 view to camera left — they are looking
  at each other across the central empty column
- The CENTER vertical third is left as clean empty space for
  post-production typography overlay
- Both women are EXPLICITLY THE SAME PERSON — same face shape,
  same hairstyle, same outfit, same accessories — only their
  skin condition, expression, and gesture differ

[IDENTICAL DETAILS — both women must have EXACTLY these traits]
- Face shape: round soft feminine face, delicate jawline, East Asian
- Hair: short layered bob-cut just below ear, dark natural black,
  slightly tousled with visible volume at crown, no clear parting,
  wispy uneven bangs — both have the SAME hairstyle
- Earrings: small simple silver stud earrings on both
- Clothing: white sleeveless tank top with thin spaghetti straps,
  simple scoop neckline — both wearing identical white tank top
- Necklace: thin delicate silver chain with tiny pendant at collarbone
- Pose baseline: chest-up framing, eye-level camera angle,
  shoulders visible

[LEFT-SIDE WOMAN — "before" / distressed state]
- VIEWPOINT: 3/4 view (face turned ~30 degrees toward camera right,
  eyes looking past camera into middle distance)
- GESTURE (人物动作): her RIGHT hand raised to her left cheek,
  fingertips gently touching the blemish area as if inspecting
  it — natural self-conscious gesture, NOT dramatic
- HEAD ANGLE: head tilted slightly DOWN and slightly forward,
  chin tucked — conveying worry and introspection
- Expression: distressed and anxious, eyebrows drawn together
  in worried frown, lips pressed together, eyes showing concern
- Skin condition: visible eczema and blemishes on BOTH cheeks —
  clusters of small pink-red inflamed dots and patches, realistic
  skin redness and slight swelling, NOT exaggerated but clearly visible
- Hair detail: a few strands falling forward onto forehead,
  slightly disheveled reinforcing the stressed mood

[RIGHT-SIDE WOMAN — "after" / confident state]
- VIEWPOINT: 3/4 view (face turned ~30 degrees toward camera left,
  eyes looking toward the LEFT subject — toward her past self)
- GESTURE (人物动作): her LEFT hand reaching up behind her head,
  fingers lightly twirling a strand of lifted hair — natural
  confident gesture suggesting she just turned to look at the
  other version of herself
- HEAD ANGLE: head tilted slightly UP and slightly back,
  chin slightly raised — conveying quiet confidence
- Expression: genuinely happy and confident, bright cheerful smile
  with teeth visible, eyes squinted slightly with joy (NOT closed),
  radiating self-assurance
- Skin condition: completely smooth, clear, flawless, radiant,
  glowing complexion — zero blemishes, zero redness, natural
  healthy glow
- Hair detail: a few strands lifted mid-motion as if caught in
  a soft breeze, reinforcing the dynamic confident mood

[STRICT PROHIBITIONS — ZERO TEXT AND NO UI ELEMENTS]
NO Chinese characters, NO Japanese characters, NO English letters,
NO Korean characters, NO numbers, NO typography of any kind,
NO watermarks, NO logos, NO signatures, NO symbols resembling
characters, NO emojis, NO buttons or call-to-action shapes,
NO hand-pointing icons, NO blue decorative borders on top or
bottom edges, NO water droplets, NO splash textures,
NO fitness gym yoga dance sports workout scenes,
NO athletic clothing, NO scene motion background.

[FINAL POLISH]
High-end commercial skincare advertisement photography,
editorial fashion magazine cover quality,
blue and white color harmony with healthy skin tones,
healing and renewal emotional arc,
aspirational yet relatable beauty standard.
```

---

## Source element vs prompt coverage table

| Element from ref | Prompt covers? | Notes |
|------------------|---------------|-------|
| 16:9 横幅 | ✅ | "16:9 aspect ratio, banner" |
| 左右双人物对称 | ✅ | "SYMMETRICAL DUAL-PORTRAIT LAYOUT — MIRROR COMPOSITION" |
| 中央留白放文字 | ✅ | "CENTER vertical third... clean empty space for post-production typography" |
| 双人物是同一人 | ✅ | "Both women are EXPLICITLY THE SAME PERSON" + 7 identical attributes |
| 短发亚洲女性 | ✅ | "short layered bob-cut, East Asian" |
| 皮肤瑕疵对比 | ✅ | LEFT "visible eczema on both cheeks, clusters of pink-red inflamed dots" |
| 修复后光滑 | ✅ | RIGHT "completely smooth, clear, flawless, radiant" |
| 蓝色主色调 | ✅ V5 NEW | "BLUE AS DOMINANT TONE — deep sapphire to bright sky blue gradient" |
| 摄影写实风 | ✅ | "PHOTOREALISTIC style + Sony A7IV + 85mm + skin texture + natural pores" |
| **人物角度 3/4 侧** | ✅ V5 NEW | "3/4 view (~30 degrees), mirror composition" |
| **人物具体动作** | ✅ V5 NEW | LEFT "指尖轻抚脸颊瑕疵处" + RIGHT "手撩发于脑后" |
| 无文字 (post-PIL overlay) | ✅ | 11 strict prohibitions |
| 无蓝色边框 | ✅ | "NO blue decorative borders on top/bottom edges" |
| 无按钮/UI | ✅ | "NO buttons, NO call-to-action shapes, NO hand-pointing icons" |
| 无水珠（V2 有，V5 删） | ✅ V5 NEW | "NO water droplets, NO splash" — 改纯净蓝渐变 |

## 修复路径对比

| 项 | V2 (2026-06-12 第一轮) | V5 (本轮) |
|----|---------------------|-----------|
| 人物角度 | "head tilted very slightly forward/back"（弱） | "3/4 view, ~30度 turn, mirror composition"（明确） |
| 人物动作 | "shoulders raised and tensed / shoulders relaxed"（弱） | LEFT "右手抚脸颊瑕疵处" + RIGHT "左手撩发于脑后"（具体） |
| 蓝色主调 | "light blue to pale blue gradient"（浅蓝） | "deep sapphire → bright sky blue, 蓝雾景深"（深蓝主调） |
| 背景 | "realistic water droplets"（水珠） | "clean dreamy blue gradient with faint blue mist"（纯净蓝渐变） |
| 风格 | photorealistic magic combo（保留） | photorealistic magic combo（保留）+ 编辑级 fashion photography |

## Generation plan (待 Devin 确认后再跑)

```bash
cd ~/Desktop/developments/via54Design
./via54 img --scene "<full prompt above>" \
  --platform minimax \
  --ar 16:9 \
  --n 4 \
  --prefix eczema-v5-blue-motion
```

预期: 4 张图 → 视觉评估 → 跨组排名 → 一句话总结（如 2026-06-12 偏好节奏）。

## 等 Devin 确认以下 N 点

1. ✅ 默认平台 **minimax** 可以吗？如换 MJ/Flux/Gemini 我重写。
2. ✅ 人物角度用 **3/4 侧互为镜像** ok？还是想要正脸 / low angle / 仰视 / 俯视？
3. ✅ 人物动作用 **"抚脸颊瑕疵" + "撩发回头笑"** ok？备选：托腮 / 手抵额头 / 转身回头看 / 抬手挡光等。
4. ✅ 蓝色主调用 **"深蓝→天蓝渐变 + 蓝雾景深"** ok？还是想要更纯净的浅蓝 / 冷蓝 / 午夜深蓝？
5. ✅ 删掉 V2 的 **水珠背景** 改纯净蓝渐变 ok？
6. ✅ 同款白色背心、银项链、银耳钉 — 同款细节 ok？

确认后我跑生成 + 给你 4 张图 + 视觉评估 + 跨组排名。
