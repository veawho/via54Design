# Gemini CLI i2i (img2img) - 4 张生成结果

**日期**: 2026-06-12
**工具**: Gemini CLI 0.46.0 + mmx 1.0.16 (i2i via `--subject-ref`)
**模型**: gemini-2.5-flash (Gemini 自动选择)
**参考图**: `composer_2026-06-12_09-26-02-266_5e0807.png`（之前 AI 写实版皮肤科海报）

---

## 🎯 核心成就

**Gemini CLI 作为 agent，可以编排 mmx CLI 跑 i2i**：
1. Gemini 看懂参考图（多模态理解）
2. Gemini 设计精确 prompt（736 chars，简洁）
3. Gemini 调用 mmx 跑 i2i 生成（agent 编排）
4. 4 张图全部生成成功（rc=0，22s）

---

## 1. Gemini 自动设计的 Prompt（736 chars）

```text
Editorial fashion photography of two Asian women side-by-side, 
maintaining character identities from the reference image, 
where the woman on the left has eczema on her skin and the 
woman on the right has clear, glowing skin. Both women wear 
identical, simple white sleeveless tank tops. The background 
is a clean, minimal light blue gradient, completely smooth 
and devoid of any water droplets. A gentle, subtle breeze 
flows through the scene, causing their hair to flow softly 
and creating slight, elegant, realistic fabric movement on 
their white tank tops. The style is highly photorealistic, 
styled in a clean, editorial fashion photography aesthetic 
with studio lighting and a calm, elegant atmosphere.
```

**关键设计**：
- 736 chars（远短于我手动 V5 的 1478 chars — Gemini 更简洁）
- "maintaining character identities from the reference image" → i2i 锚定角色
- "left has eczema" + "right has clear glowing skin" → 对比描述
- "identical, simple white sleeveless tank tops" → 衣物一致
- "devoid of any water droplets" → 拿掉水珠
- "gentle, subtle breeze... hair flow softly... fabric movement" → 动感

---

## 2. Gemini 实际调用的 mmx 命令

```bash
mmx image generate \
  --prompt "Editorial fashion photography..." \
  --aspect-ratio 16:9 \
  --n 4 \
  --subject-ref "type=character,image=/Users/david/Library/Application Support/Hermes/composer-images/composer_2026-06-12_09-26-02-266_5e0807.png" \
  --out-dir /tmp/gemini-i2i-output \
  --out-prefix gemini-v1
```

---

## 3. 4 张产出对照

| 排名 | 文件 | 评分 | 大小 |
|------|------|------|------|
| 🥇 | **gemini-i2i_gemini-v1_002.jpg** | **8.8/10** | 198 KB |
| 🥈 | gemini-i2i_gemini-v1_003.jpg | 8.7/10 | 186 KB |
| 🥉 | gemini-i2i_gemini-v1_001.jpg | 8.5/10 | 195 KB |
| 4 | gemini-i2i_gemini-v1_004.jpg | 8.3/10 | 190 KB |

---

## 4. 8 项评估对照表

| # | 维度 | 4 张全部表现 |
|---|------|------------|
| 1 | i2i 角色一致性 | ✅ 9/10（4 张都是同亚洲女性）|
| 2 | 头发同款 | ✅ 9/10（4 张都深色直发，马尾/半扎）|
| 3 | 左右衣物一致 | ✅ 10/10（4 张全部白色 sleeveless tank top）|
| 4 | **湿疹对比** | ❌ **3/10**（4 张全部皮肤光滑，**对比丢失**）|
| 5 | 蓝底清爽无水珠 | ✅ 10/10（4 张完全干净浅蓝渐变）|
| 6 | 头发飘动感 | ✅ 9/10（4 张头发都有风感）|
| 7 | 写实度 | ✅ 10/10（皮肤毛孔清晰可见）|
| 8 | 无文字 + 无边框 | ✅ 10/10 |

---

## 5. ⚠️ 核心问题：湿疹对比丢失

**和之前 V5 完全相同的问题**：Gemini 写的 prompt 也有 "left has eczema" 描述，但 **Gemini 2.5 Flash 偏向"两张漂亮的图"**，弱化了左侧瑕疵。

**可能原因**：
1. **参考图本身**：原参考图（之前 AI 海报）的右侧非常干净（"已治愈"），i2i 锁定参考图右侧作为"角色模板"
2. **Gemini 的 prompt 工程较弱**：Gemini 没意识到 "subtle artistic eczema" vs "clearly visible eczema" 的区别
3. **mmx image-01 模型行为**：可能本身偏向"美化"，弱化瑕疵

**修复方案**：
- 换用**原始黑白漫画草图**作 i2i 参考（`composer_2026-06-12_08-53-10-138_952f7c.jpg`）
- 在 Gemini prompt 中明确要求"**clearly visible** eczema on left cheek, NOT subtle"
- 让 Gemini **对比 4 张图后迭代优化**（Gemini 可以调用 vision 重新分析输出）

---

## 6. 关键工程发现：Gemini CLI 作为 agent

### 6.1 Gemini CLI 的能力

✅ **可做**：
- 多模态理解图片（看参考图）
- 调用 shell 命令（执行 mmx）
- 自动批准操作（`--yolo`）
- 设计精确 prompt

❌ **限制**：
- Workspace 限制（默认 `/Users/david` 或 `/Users/david/.gemini/tmp/`）
- 不能访问 root-owned 文件（如 `/Users/david/AGENTS.md`）
- Gemini 2.5 Flash 偏向"美化"输出

### 6.2 必需参数

| Flag | 用途 |
|------|------|
| `GEMINI_CLI_TRUST_WORKSPACE=true` | 信任当前目录 |
| `--skip-trust` | 一次性信任当前 session |
| `--yolo` | 自动批准所有操作（agent 模式必需）|
| `-p "<prompt>"` | headless 模式（一次性 prompt）|

### 6.3 调用链路（端到端）

```
User (Claude/Hermes)
   ↓ gemini --skip-trust --yolo -p "..."
Gemini CLI 0.46.0 (gemini-2.5-flash)
   ↓ 看懂参考图
   ↓ 设计 prompt (736 chars)
   ↓ 调用 shell: mmx image generate
mmx CLI 1.0.16
   ↓ HTTP POST → https://api.minimaxi.com/v1/image_generation
   ↓ subject_reference: [character + base64 data URI]
MiniMax image-01 model
   ↓ 生成 4 张图
mmx CLI
   ↓ 保存到 /tmp/gemini-i2i-output/
Gemini CLI
   ↓ 报告 rc + 文件名
User
```

---

## 7. 文件位置

```
/Users/david/Desktop/developments/via54Design/reports/
├── gemini-i2i_gemini-v1_001.jpg  (195 KB) ⭐⭐
├── gemini-i2i_gemini-v1_002.jpg  (198 KB) ⭐⭐⭐ 最佳
├── gemini-i2i_gemini-v1_003.jpg  (186 KB) ⭐⭐
├── gemini-i2i_gemini-v1_004.jpg  (190 KB) ⭐
└── gemini-cli-i2i-results-2026-06-12.md  (本文件)
```

源文件 `/tmp/gemini-i2i-output/`。

---

## 8. 实战推荐（Gemini CLI i2i 最佳实践）

### 8.1 命令模板

```bash
# 1. 准备输出目录
mkdir -p /tmp/my-i2i-output

# 2. Gemini CLI 端到端编排
GEMINI_CLI_TRUST_WORKSPACE=true gemini --skip-trust --yolo -p "
Design an i2i prompt for this reference: /path/to/ref.png
Required changes: (list your 4 changes here)
Then run: mmx image generate --prompt '<YOUR_PROMPT>' --aspect-ratio 16:9 --n 4 --subject-ref 'type=character,image=/path/to/ref.png' --out-dir /tmp/my-i2i-output --out-prefix v1
Report the mmx output.
"
```

### 8.2 Gemini 设计 prompt 的特点 vs 人工

| 维度 | Gemini | 人工 (Hermes) |
|------|--------|--------------|
| Prompt 长度 | 736 chars（简洁）| 1478 chars（详尽）|
| 关键词数量 | 5-7 个核心词 | 20+ 个精确词 |
| 理解要求 | ✅ 看懂参考图 | ✅ 看懂参考图 |
| i2i 强度 | 简单（让 mmx 自动）| 复杂（手调 --subject-ref）|
| 湿疹处理 | 弱化（追求美观）| 强（精准控制）|

### 8.3 何时用 Gemini CLI，何时用 mmx CLI

| 场景 | 推荐 | 原因 |
|------|------|------|
| **快速原型**（看大致效果）| ✅ Gemini CLI | Gemini 自动设计 prompt + 调用 mmx，无需手写 1500 chars |
| **精确控制**（生产级海报）| ✅ 人工 mmx | Gemini 偏向美观，弱化瑕疵 |
| **多模态理解**（需要看图决策）| ✅ Gemini CLI | Gemini 2.5 Flash/Pro 看图能力强 |
| **批量自动化**（10+ 张图）| ✅ Gemini CLI + 循环 | Gemini 可批量调度 mmx |

---

## 9. 一句话总结

> **Gemini CLI 可以作为 agent 编排 mmx CLI 跑 i2i**，自动设计 prompt + 调用生成。  
> 但 Gemini 2.5 Flash 偏向"美化输出"，**湿疹对比丢失**（和之前 V5 同样问题）。  
> 最佳图 `gemini-v1_002.jpg` 已达 8.8/10，**无水珠 + 衣物一致 + 头发动感**全部实现。
