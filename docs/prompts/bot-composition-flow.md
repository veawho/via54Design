---
name: via54-prompt-generation
description: |
  通过飞书私聊/群聊接收用户的中文需求或参考图片+文字需求,调 via54Design
  工具生成 AI 绘图提示词。流程: 文字 → 3 段构图方案 → 用户选 → 完整英文 prompt。
  流程 B: 图片+文字 → vision 描述 → via54 → 中文拆解版 → 用户修改 → 完整英文 prompt。
keywords:
  - via54
  - prompt generation
  - 提示词
  - 设计模板
  - i2i
  - image to prompt
  - 构图方案
  - 英文 prompt
version: "1.0.0"
---

# via54-prompt-generation SKILL

> 通过 via54Design 工具,为飞书用户生成 AI 绘图/视频的提示词。
> 14 平台 + 26 维度 + 4 叙事模型 = via54Design 全能力。

---

## 何时使用

**Trigger**:
- 飞书私聊/群聊消息内容包含 "我想做一张..." / "根据参考图片..." / "生成提示词" / "1:1 提示词"
- 用户发参考图 (post / image 消息类型) + 文字需求
- 用户明确说"基于参考图 + 文字"生成

**不适用**:
- 纯闲聊
- 非设计类请求 (查天气 / 写代码 / 翻译)

---

## 工具

| 工具 | 路径 | 用途 |
|------|------|------|
| `~/.local/bin/via54` | via54Design CLI v3.0 | 生成 3 段构图 + 完整英文 prompt |
| `via54 prompt list` | 列出所有 14 平台 | midjourney / flux / dalle3 / sd3 / stable_diffusion / ideogram / recraft / seedance / gemini / veo / sora / kling / pika / jimeng |
| `via54 prompt generate --scene "..." --platform <plat>` | 渲染完整英文 prompt | 平台特定优化 (字数 / 关键词 / 风格) |
| `tools/vision_tools.vision_analyze_tool` | 图片 → 中文描述 | minimax-cn MiniMax-M3 |

---

## 流程 A: 文字指令生成提示词

### Step 1: 用户发起需求 (私聊/群聊)

**用户示例**:
> "我想做一张湿疹主题的面向年轻女性的广告海报。"

### Step 2: 思考 3 种构图方案

**bot 行为**:
1. 收到用户中文需求
2. 调 `via54 prompt` 思考 3 种构图方案
   ```bash
   ~/.local/bin/via54 prompt \
     --scene "湿疹主题广告海报,面向年轻亚洲女性,护肤产品营销,需展现治疗前后皮肤状态对比" \
     --platform jimeng \
     --mode compositions \
     --count 3
   ```
3. bot **不直接返 via54 CLI 输出**(技术细节),而是把 3 段**简化成用户友好描述**返给用户:

```
请选择您喜欢的构图方案:

【方案 1: 左右对称对比构图】
左侧女性(焦虑+湿疹),右侧女性(自信+健康肌肤),中间"湿疹主题"标题,
深蓝色边框,水珠凝结玻璃背景,情感反差强烈。

【方案 2: 前后景层叠构图】
前景女性(右侧/微笑/健康)半身右侧位,后景女性(左侧/焦虑/湿疹)半身左侧位模糊化,
中央"皮肤病专科医生"按钮,引导互动。

【方案 3: 经典电影分屏构图】
左侧怼脸特写(湿疹皮肤纹理),右侧远景玻璃窗前女性背影剪影,
中央标题横排,治愈系医疗广告风格。

回复 "1" / "2" / "3" 选择方案,或回复"重新生成"换方案。
```

### Step 3: 用户回复 + 生成完整英文 prompt

**用户回复** (任一):
- "1" → 选方案 1
- "重新生成" → 跳回 Step 2
- "2,然后改成蓝色调" → 选方案 2 + 修改意见

**bot 行为**:
1. 调 `via54 prompt generate` 渲染完整英文 prompt:
   ```bash
   ~/.local/bin/via54 prompt \
     --scene "<3 段构图 + 用户修改意见>" \
     --platform <user-specified or "通用">
   ```
2. **如果用户指定平台** (jimeng/midjourney/flux/gemini/sora 等),按平台优化:
   - jimeng: ≤1500 字 + 中文友好
   - midjourney: ≤4500 字 + `--style raw` + `--v 6`
   - flux: ≤4000 字 + 自然语言
   - gemini: ≤8000 字 + 详细描述
   - sora: ≤4000 字 + 视频特定
3. **如果用户没指定**,用"通用"模式生成细节拉满的英文 prompt:
   ```
   A cinematic split-screen photograph of two young Asian women...
   [详细到每个细节: 构图/光线/色调/姿态/表情/服饰/背景/镜头]
   ```
4. bot 把英文 prompt 包装成**中文提示** + **英文 prompt 块**:
   ```
   ✅ 完整英文 prompt 已生成 (jimeng 平台 / 1500 字):
   
   ```
   [英文 prompt 块,用户可直接复制到 jimeng 网页/APP]
   ```
   
   💡 使用提示:
   - 复制上方英文 prompt → 打开 jimeng 网页/APP → 粘贴 → 生图
   - 如想调整,直接回复"修改: <您的意见>"
   - 如想换平台,回复"平台: midjourney"
   - 如想换构图,回复"重新生成"
   ```

### 过程步骤 (中间对话)

**用户在确认构图方案前**,所有对话都用中文需求沟通,**不用先生成完整英文 prompt 再拆解中文**(浪费时间)。

**bot 处理逻辑**:
- 用户给修改意见 → bot 把意见合并到当前方案 → 再生成 3 段新方案 (走 Step 2)
- 用户选方案 → 才生成完整英文 prompt (走 Step 3)
- **永远避免**:先渲英文 → 再让 LLM 翻译中文 → 再让用户改 → 再渲英文 (浪费 token)

### 最终步骤 (用户确认 1+ 个方案)

**必须返**:
- 1 个或多个方案对应的完整英文 prompt (平台特定优化)
- 中文操作提示 (复制使用 / 修改 / 换平台)

---

## 流程 B: 参考图片+文字指令

### Step 1: 用户发参考图 + 文字需求

**用户示例**:
> (发参考图 + 文字) "根据参考图片,我想做一张湿疹主题的面向年轻女性的广告海报。"

### 中间过程

1. **下载参考图到本地** (用飞书 im/v1/messages/{msg_id}/resources/{image_key}?type=image API)
2. **调 vision_analyze_tool 描述图** (minimax-cn MiniMax-M3, 200-500 字纯中文):
   ```
   一幅采用黑白线描风格绘制的横向宣传图,主色调为灰白黑三色...
   左侧女性(雀斑+湿疹),右侧女性(光洁+自信)...
   ```
3. **调 via54 把 vision 中文描述 + 用户文字需求 → 渲染英文 prompt**:
   ```bash
   ~/.local/bin/via54 prompt \
     --scene "<vision 中文描述> + <用户文字需求>" \
     --platform <user-specified or "通用">
   ```
4. **bot 拿到完整英文 prompt → 调 LLM 简化 + 中文拆解**:
   ```
   ### 主体
   - 主体 (Subject): 左右两位年轻亚洲女性
   - 辅助 (Secondary): 中间"湿疹主题"标题
   
   ### 风格
   - 风格流派: 摄影写实
   - 媒介: digital photography
   - 色调: 蓝色调 (按你之前要求)
   
   ### 构图
   - 镜头: medium shot
   - 构图: rule of thirds
   
   ### 光线
   - 光线: cinematic lighting
   
   ### 色彩
   - 色调: cool blue grading
   ```
5. **bot 把"完整英文 prompt 文件" + "中文拆解版"两个一起返给用户**

### Step 2: bot 首次回复 (关键 - 含 2 点确认逻辑)

**bot 必须保留以下逻辑**(首次返图+文):

> **我已根据参考图 + 文字需求生成中文拆解版提示词**。
>
> 在生成完整英文 prompt 前,**需要确认 2 点**:
>
> 1. **左侧湿疹皮肤** — 要保留"明显红斑/皮屑"的真实细节,还是更克制一点(避免观感不适)?
> 2. **右侧健康皮肤** — 要"发光透亮"的高光感,还是"自然健康"的哑光感?
>
> 你直接打字告诉我,我来调整。

**为什么必须有**:参考图的"明显红斑"细节可能让生成图观感不适(医疗广告伦理),"高光感 vs 哑光感"是营销策略选择(高端感 vs 真实感)。让用户先确认这两点,避免生成图后反复修改。

### Step 3: 用户回复 + 处理

**用户回复**(任一):
- "1: 克制 / 2: 哑光" → bot 把意见合到 vision 中文描述 → 调 via54 重渲英文 prompt
- "1: 真实细节 / 2: 高光感" → bot 直接生成完整英文 prompt + 中文拆解版
- "重新生成 3 段构图" → 走流程 A 的 Step 2

### Step 4: bot 返完整英文 prompt + 中文拆解版 + 操作提示

同流程 A 的 Step 3。

---

## 平台特定优化

| 平台 | 字数限制 | 关键参数 | 风格 |
|------|---------|----------|------|
| jimeng | ≤1500 字 | 中文友好 | 中文营销 |
| midjourney | ≤4500 字 | `--style raw --v 6` | 高端艺术 |
| flux | ≤4000 字 | 自然语言 | 真实感 |
| gemini | ≤8000 字 | 详细描述 | 详细精确 |
| sora | ≤4000 字 | 视频特定 | 动态镜头 |
| stable_diffusion | ≤380 字 | 短关键词 | 通用 |
| 通用 (默认) | 不限 | 详细段落 | 全场景 |

**bot 默认行为**:
1. 第一次: **不指定平台**,生成"通用"细节拉满版
2. 用户说"平台: <X>" → 重新生成平台特定优化版

---

## 多轮对话 session 记忆

bot 必须记住:
- 当前选定方案编号 (1/2/3)
- 用户修改意见列表
- 目标平台 (jimeng/midjourney/flux/...)
- 当前完整英文 prompt (用于"修改" 时基线)

**session 存**:
```
/tmp/m12_bot/sessions/{chat_id}.json
{
  "phase": "compositions" | "modify" | "final",
  "selected_composition": 1,
  "modifications": ["蓝色调", "摄影写实", "去掉按钮"],
  "platform": "jimeng",
  "english_prompt": "...",
  "chinese_breakdown": "..."
}
```

**bot 每次收到消息**:
1. 读 session (chat_id)
2. 决定当前 phase:
   - 无 session → 进入"compositions" phase (Step 2)
   - 有 session + 用户给方案号 → 进入"final" phase (Step 3)
   - 有 session + 用户给修改意见 → 进入"modify" phase (重渲英文 + 应用修改)
   - 有 session + 用户说"重新生成" → 重置 → "compositions"
3. 处理
4. 写 session

---

## 错误处理

| 错误 | 处理 |
|------|------|
| 用户发图但没文字 | bot 反问"请补文字需求" |
| vision 描述失败 | bot 用"通用版"提示词 (描述部分空着) |
| via54 CLI 失败 | bot 重试一次,失败返"via54 工具异常,请稍后再试" |
| 用户指定不存在的平台 | bot 列出 14 平台让用户选 |
| 用户消息超长 | bot 截断到 4000 字 |
| 多轮修改累计修改超 5 次 | bot 提示"建议先确认当前方案,再批量应用修改" |

---

## 示例对话

### 流程 A 示例

```
[用户] 我想做一张湿疹主题的面向年轻女性的广告海报

[bot] 请选择您喜欢的构图方案:

【方案 1: 左右对称对比构图】...
【方案 2: 前后景层叠构图】...
【方案 3: 经典电影分屏构图】...

回复 "1" / "2" / "3" 选择方案,或回复"重新生成"换方案。

[用户] 1

[bot] ✅ 完整英文 prompt 已生成 (通用版 / 1200 字):

```
A cinematic split-screen photograph of two young Asian women in their 20s...
```

💡 使用提示:
- 复制上方英文 prompt → 打开 AI 绘图工具 → 粘贴 → 生图
- 修改: 直接回复"修改: <意见>"
- 换平台: "平台: jimeng"
```

### 流程 B 示例

```
[用户] (发图) 根据参考图片,我想做一张湿疹主题的面向年轻女性的广告海报

[bot] 我已根据参考图 + 文字需求生成中文拆解版提示词:

### 主体
- 主体: 左右两位年轻亚洲女性
- 左侧: 焦虑+湿疹红斑
- 右侧: 灿烂笑容+光洁肌肤
...

在生成完整英文 prompt 前,需要确认 2 点:

1. 左侧湿疹皮肤 — 要保留"明显红斑/皮屑"的真实细节,还是更克制一点?
2. 右侧健康皮肤 — 要"发光透亮"的高光感,还是"自然健康"的哑光感?

你直接打字告诉我,我来调整。

[用户] 1 克制 / 2 哑光

[bot] ✅ 完整英文 prompt 已生成 (通用版 / 1300 字):

```
A cinematic split-screen photograph of two young Asian women...
[左侧皮肤克制细腻 / 右侧哑光健康肤色]
...
```
```

---

## 真测试命令

```bash
# 测 via54 CLI 可用
~/.local/bin/via54 prompt list

# 测 3 段构图
~/.local/bin/via54 prompt \
  --scene "湿疹主题广告海报, 两位年轻亚洲女性对比" \
  --platform jimeng \
  --mode compositions \
  --count 3

# 测完整英文 prompt
~/.local/bin/via54 prompt \
  --scene "<scene>" \
  --platform midjourney
```

---

## 相关仓库增量

| 仓库 | 增量 |
|------|------|
| via54Design | `docs/prompts/bot-composition-flow.md` (新增) + `AGENTS.md` 加 "Bot 集成" 章节 |
| via54Larkfix | `~/.hermes/scripts/m12_full_channel_bot.py` (m12 bot 接入 via54) |
| via54Larkbotgo | `docs/feishu-bot-integration.md` (更新) + GO SDK 用 channel.bot.on_message |
| via54Skills | `skills/via54-prompt-generation/SKILL.md` (新增,本文件副本) |
| via54Hermes | `~/.hermes/skills/via54-prompt-generation/SKILL.md` (本文件副本) + SOUL.md 加引用 |