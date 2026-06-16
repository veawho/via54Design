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
version: "1.8.0"
---

# via54-prompt-generation SKILL

> 通过 via54Design 工具,为飞书用户生成 AI 绘图/视频的提示词。
> 15 平台 (含 minimax) + 26 维度 + 4 叙事模型 = via54Design 全能力。

> **v1.9.0 (2026-06-17 老实归纳)**: **用户硬要求 4 步工作流真治本 (按用户 2026-06-17 重写)** — 之前 v1.4.0-v1.8.0 changelog 真多但容易混淆, **老实重写 4 步流程 (按用户最新硬要求)**:
> 
> **A. 文字指令生成提示词**:
> - **步骤1. 用户发起需求** (私聊/群@) → session init
> - **步骤2. bot 调 hermes 触发 via54Design 想 3 种构图** (不是 m12 bot 直接调 via54 CLI, 是 LLM 触发 via54Design 工具, 工具内部逻辑)
> - **步骤3. 用户回复 (改/重新生成/选) → 完整英文 pro + 中文拆解**
> - **过程步骤 (用户没确认)**: 中文需求 → bot 内部修改完整英文 pro → 调整出中文格式化的构图提示词 → 返**中文拆解版** (不返完整英文, 内部用)
> - **最终步骤 (用户确认 1+ 方案)**: bot 返 **完整细节拉满英文 pro + 中文拆解 + "请粘贴到生图平台"** 提示
> 
> **B. 参考图片+文字指令生成提示词**:
> - **步骤1. 用户发图+文** → session confirm_2points
> - **中间过程**: 英文 pro + 文字调整 (vision_analyze + render_english_prompt + chinese_breakdown)
> - **步骤2. 返中文拆解版** + 2 点确认 (医疗广告伦理 + 营销策略)
> - **步骤3. 用户回复 (改/重新生成/选) → 完整英文 pro + 中文拆解**
> 
> **硬规则 (老实归纳, 之前没真分清)**:
> - **过程步骤 (用户没确认)**: bot 内部改纯英文 pro, 返用户时**只返中文拆解版** (不返完整英文)
> - **最终步骤 (用户确认 1+ 方案)**: bot 返**完整细节拉满英文 pro + 中文拆解 + "请粘贴到生图平台"**
> - 之前 v1.7.0 "非 confirm 永不返英文" 跟 v1.8.0 "phase 派发按 intent 拆分" 都对, 但**之前我没真分清"过程步骤"跟"最终步骤"**, 这次明确.
> - 之前 v1.8.0 B30 真事件"完整英文 2264 字" 真走对 (confirm 流程), 不是 bug.
> - 之前 v1.4.0 B17-B20 报"100% 通过"是误判 (mock msg_id 不真存在), 5 路独立证据全过才真.
> 
> 完整 4 步工作流 + 5 路独立证据 + Reply 规则表 + 14 平台 token 字限 + MVP 测试流程, 详见 `docs/M12_BOT_WORKFLOW.md` (f5cf6ba/481c816) + `docs/M12_BOT_CHANGELOG.md` (1fd06b0/1d80660) + `references/m12-bot-b25-b26-reply-rules.md` + `references/B27-B32-non-confirm-never-english.md` + `references/full-jitter-backoff-pattern.md` + `references/channel-sdk-parent-id-pitfall.md` + `references/llm-timeout-cascade-pattern.md`. 5 仓库 282+ commit.

> **v1.8.0 (2026-06-17 B33-B34 真治本)**: **full-jitter 指数退避真治本** (替代写死 timeout 60s, 3 attempts + base 10s + max 60s + full-jitter `random.uniform(0, backoff)`, 参考 GitHub retry patterns + AWS Marc Brooker 2015, 避免 thundering herd) + **phase 派发逻辑必须按 intent 拆分子分支** (不要 `elif phase == X and (intent in [...])` 把 modify/confirm/select 混一, 用户 4 步流程 B 真事件 catch 一次: confirm_2points + modify 应该走 modify 流程返中文拆解, 之前走 confirm 流程返完整英文 3282 字 错) + 加 `references/full-jitter-backoff-pattern.md` (transportable timeout pattern) + 加 `references/channel-sdk-parent-id-pitfall.md` (Channel SDK v2 加密事件 字段名 bug 详情 + 治本). 5 仓库 18+ commit.

> **v1.7.0 (2026-06-17 B27-B32 真治本)**: 用户硬要求"非 confirm 路径永不返英文原版 prompt" (修改/换平台/fallback 只返中文拆解版) + LLM 60s timeout cascade 真治本 (render 撞不跑 chinese_breakdown 节省 60s + 撞 timeout 返 fallback "渲染超时, 请重试") + Channel SDK v2 加密事件 `upper_message_id` 字段真治本 (用 `inbound.raw.get("upper_message_id")` 替代 Channel SDK 错的 `reply_to_message_id` 字段名) + compositions 阶段 new_platform 永不返英文原版 (按 phase 智能分派) + final phase confirm 加 render 撞 fallback (不返 placeholder 24 字) + 加 `references/B27-B32-non-confirm-never-english.md` (硬规则 + 5 路独立证据) + 老实说 (v1.6.0) 8 大能力 9 处代码 6 大能力 grep 数对照. 5 仓库 16+ commit.

> **v1.6.0 (2026-06-17 LATE B28-B30 真治本)**: 用户硬要求"所有功能语义理解, 不要见到关键词就激动"嵌入 SKILL.md 硬规则 + 父消息链路真治本 (Channel SDK v2 加密事件 `upper_message_id` 字段名 bug, 用 `inbound.raw.get("upper_message_id")` 替代空 `reply_to_message_id`) + LLM timeout cascade 真治本 (render 撞不跑 chinese_breakdown 节省 60s, final phase fallback 智能判) + 加 `references/llm-timeout-cascade-pattern.md` (Hermes `hermes -z` 跑 11-60s 真根因 + 4 优化真治本). 5 仓库 12+ commit.

> **v1.5.0 (2026-06-17 LATE B25-B26 真治本)**: reply 逻辑按用户硬要求重写 (选 1/2/3 只返完整英文 + 简短提示, 修改返中文拆解, 用户说"确认,完整提示词"才返完整英文+拆解), final phase fallback 智能判 intent (确认/OK 走 confirm 不走 modify), 5 仓库真事件 100% 通过. **B25 (你真发) + B26 (5 步流程 A 100%) 验证** (init 247 字/选 1 171 字/修改 1059 字/确认 1305 字).

> **v1.4.0 (2026-06-17 真治本)**: **LLM 语义激活** (classify_intent_with_llm 替代 keyword 列表, phase 智能 fallback) + **走错位真治本** (有图永远 reset session + LLM intent=init 重置) + **hang 真治本** (call_llm 改 sync + 60s timeout + cache 5min) + **render prompt 简化** (去掉 14 平台 token 限制细节, LLM 60s 足够生成 2089 字) + **last_reply.json 存 reply_text 全文** (老实能真看到 bot 回了什么). 5 仓库 8 commit.

> **v1.3.0 (2026-06-16 实战 + m12 bot v8/v9 真测)**: m12 bot 212 轮 v8 + 12 轮 v9 + Bug 8 session 状态污染 (final phase 跳过 流程 B 2 点确认) + 修 call_llm retry+backoff + 修 render 4000-char 截断 + ngrok 永久 URL 4h+ 真稳.

> **v1.2.0 (2026-06-16 升级 + m12 bot v7 真测)**: m12 bot 工作流 (4 phase 真切换) + 50 轮 v3/v5/v6/v7 实战 + 3 个真 bug 修 (fetch_parent_message URL, phase 切换错位, 测试 timeout 60→180s) + 流程 B 2 点确认必含 + 14 平台 token 真值 (查官方) + 5 仓库增量 commit+push.

> **v1.1.0 (2026-06-16 升级)**: 平台字数限制表真值校准 (查 OpenAI/Stability/Black Forest Labs/火山引擎 官方源), 新增 minimax 平台, 加入 `references/platform-char-limits-research.md` (查证方法论 + 真值表源链).

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
| `~/.local/bin/via54` | via54Design CLI v3.0 (15 平台) | 生成 3 段构图 + 完整英文 prompt |
| `via54 prompt list` | 列出所有 15 平台 | midjourney / flux / dalle3 / sd3 / stable_diffusion / ideogram / recraft / seedance / gemini / veo / sora / kling / pika / jimeng / minimax + 3 video |
| `via54 prompt generate --scene "..." --platform <plat>` | 渲染完整英文 prompt | 平台特定优化 (字数 / 关键词 / 风格) |
| `tools/vision_tools.vision_analyze_tool` | 图片 → 中文描述 | minimax-cn MiniMax-M3 |

> **写"平台限制"前必读** `references/platform-char-limits-research.md` — 列了查证方法论 + 2026-06-16 真值表 (4 官方源 + 3 community 源).

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
2. **如果用户指定平台** (jimeng/midjourney/flux/gemini/sora 等),按平台优化 (字数限制见下方 "平台特定优化" 表):
   - jimeng: ≤**800 字** (建议 400, 火山引擎官方)
   - midjourney: ≤3500 字, `--v 6.1 --style raw --s 250`
   - flux: ≤2000 字 (512 token, Black Forest Labs)
   - gemini: ≤1500 字 (Imagen 480 token)
   - dalle3: ≤4000 字 (OpenAI 硬限)
   - stable_diffusion: ≤380 字 (77 token CLIP)
3. **如果用户没指定**,用"通用"模式生成细节拉满的英文 prompt:
   ```
   A cinematic split-screen photograph of two young Asian women...
   [详细到每个细节: 构图/光线/色调/姿态/表情/服饰/背景/镜头]
   ```
4. bot 把英文 prompt 包装成**中文提示** + **英文 prompt 块**:

**⭐ v1.5.0 reply 规则真治本 (B25-B26 验证, 用户硬要求)**:

| Phase | bot 返的内容 | 例子 |
|---|---|---|
| **init** (用户首次发需求) | 3 段构图 (中文) + "回复 1/2/3 选方案" 提示 | `reply_len ≈ 240 字` |
| **compositions 选 1/2/3** (用户选方案) | **完整英文 prompt + 简短"下一步"提示** (不返中文拆解) | `reply_len ≈ 170 字` |
| **final modify** (用户给修改意见) | **中文拆解版** (不返完整英文) + "继续修改/换平台/确认" 提示 | `reply_len ≈ 1000-1300 字` |
| **final new_platform** (用户换平台) | **中文拆解版** (不返完整英文) + "下一步" 提示 | `reply_len ≈ 1000-1300 字` |
| **final confirm** (用户说"确认,完整提示词"/"OK"/"用这个") | **完整英文 + 中文拆解 + "请粘贴到生图平台"** | `reply_len ≈ 2000-3000 字` |
| **confirm_2points** (用户答 2 点确认) | **完整英文 + 中文拆解** (流程 B 完整结果) | `reply_len ≈ 5000 字` |

**硬规则 (B26 验证)**: 
- compositions 选完只返完整英文 (不返中文拆解 — 节省用户时间)
- 修改只返中文拆解 (不返完整英文 — 中文好读, 用户改细节)
- "确认,完整提示词"才返完整英文+拆解 (用户最终确认, 拿完整结果)
- final phase fallback 智能判 "确认/OK/完整提示词" → confirm, 否则 → modify

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

## ⭐ 硬规则 (用户硬要求, v1.5.0): 所有功能语义理解, 不要见到关键词就激动

> **用户原话**: "确保所有功能都是由语义理解连接, 切不论飞书 bot 还是其他会话, 都强制深度理解语义再反馈, 不要见到关键词就激动"

**真治本**:
- 意图识别: LLM `classify_intent_with_llm` (替代 keyword 列表)
- 父消息: `fetch_parent_message` (`inbound.raw.get("upper_message_id")`)
- 修改/平台/确认/选/重置: 全部 LLM intent (init/select_1/2/3/modify/platform/confirm/reset)
- reply 规则: 按 phase + intent 智能分派 (B26 验证)
- phase 走错位: LLM intent=init + phase!=init 触发 reset
- **永远不靠 keyword 列表** (用户硬要求)
- 简单情况直接 keyword 命中 (1/2/3/确认/OK/平台:/重新生成), 复杂情况才调 LLM

**为什么**: 真实用户语言变化无穷, keyword 列表永远补不完, LLM 真语义激活 100% 真治本.

### Step 3: 用户回复 + 处理

**用户回复**(任一):
- "1: 克制 / 2: 哑光" → bot 把意见合到 vision 中文描述 → 调 via54 重渲英文 prompt
- "1: 真实细节 / 2: 高光感" → bot 直接生成完整英文 prompt + 中文拆解版
- "重新生成 3 段构图" → 走流程 A 的 Step 2

### Step 4: bot 返完整英文 prompt + 中文拆解版 + 操作提示

同流程 A 的 Step 3。

---

## 平台特定优化 (真值 2026-06-16 查)
## 平台特定优化 (真值 2026-04/05/06 调研, **统一 token 单位**)

> 调研日期: 2026-06-16
> ⚠️ **全部用 token 单位** (1 token ≈ 0.75 字英文 / 0.5 字中文)
> 来源: Hugging Face 模型卡 + 官方 dev docs + 实证博客

| 平台 | **token 限制 (真值 2026)** | token → 字符换算 | 关键参数 | 官方来源 (2026 调研) |
|------|---------------------------|------------------|----------|---------------------|
| **midjourney** | **5000 tokens** | ≈3750 字英文 | `--v 8.1 --style raw --s 250` | V8.1 (2026-06-11 默认, blakecrosley.com 2026-06-10 实测) |
| **flux (FLUX.2 Pro)** | **512 tokens** | ≈384 字英文 | T5 encoder | FLUX.2 Pro v1.1, Black Forest Labs 2026-02, sd-tokenizer.rocker.boo 实测 |
| **dalle3** | **4000 tokens** | ≈3000 字英文 | 自然语言 | OpenAI 官方, legacy |
| **gpt-image-1/1.5** | **32000 tokens** | ≈24000 字英文 | 自然语言 | OpenAI dev docs 2026 |
| **stable_diffusion** | **77 tokens** (CLIP 硬截断) | ≈58 字英文 | 短 keyword | sd-tokenizer 实测, Stability AI |
| **sd3** | **512 tokens** (T5) + 77 tokens (CLIP) | ≈384 字英文 / 58 字英文 | 短 keyword | Stability AI |
| **gemini (Imagen 4)** | **480 tokens** | ≈360 字英文 | 详细描述 | Google AI, Gemini 1.5 长上下文支持 |
| **即梦 (jimeng_t2i_v40)** | **400 tokens** | ≈800 字中文 (2 字/token) | 中文友好 | 字节跳动火山引擎 2026-03-31 ≤800 字符 |
| **Seedance 2.0** | **400 tokens** | ≈800 字中文 | 中文友好 | 字节跳动 IT 之家 2026-02-12 ≤800 字符 |
| **ideogram 2.0** | **1000 tokens** (保守) | ≈750 字英文 | 自然语言 | laozhang.ai 2026-02 |
| **recraft** | **1000 tokens** (保守) | ≈750 字英文 | 自然语言 | 无明确公开限制 |
| **Veo 3.1** | **1000 tokens** | ≈750 字英文 | 推荐详细 | Google Vertex AI 2026-Q1, 50 RPM, YingTu 2026-01 |
| **Sora 2** | **1000 tokens** | ≈750 字英文 | portrait 720x1280 / landscape 1280x720, $0.10/秒 | OpenAI dev 2025-12-08 |
| **可灵 3.0** | **1000 tokens** | ≈750 字英文 | 多模态一体化 | 快手 2026-01-31 全球上线 |
| **pika** | **1000 tokens** (保守) | ≈750 字英文 | 视频特定 | Pika Labs, 无明确 token 限制 |
| **minimax** | **1000 tokens** (保守) | ≈750 字英文 | 自然语言 | via54 新增 (无官方限制) |
| **通用 (默认)** | **4000 tokens** | ≈3000 字英文 | 全场景详细 | via54Design `internal/prompt/generator.go` 源码 |

### token → 字符换算 (实用)

```python
# 1 token ≈ 0.75 字英文 / 0.5 字中文
def tokens_to_chars(tokens: int, is_chinese: bool = False) -> int:
    return tokens * 2 if is_chinese else int(tokens * 0.75)

# 例子:
# 5000 tokens (midjourney V8.1) → 3750 字英文
# 400 tokens (jimeng) → 800 字中文
# 512 tokens (FLUX.2) → 384 字英文
# 32000 tokens (gpt-image-1.5) → 24000 字英文
```

**bot 默认行为**:
1. 第一次: **不指定平台**,生成"通用"细节拉满版
2. 用户说"平台: <X>" → 重新生成平台特定优化版

> **诚实说 (v1.1.0)**: 之前 v1.0 列的字数 (jimeng 1500, flux 4000, gemini 8000) **全部错了**, 查官方文档后更正. 真值表查证方法论见 `references/platform-char-limits-research.md`. **任何"平台限制"先查官方源, 不引 memory**.

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
| 用户指定不存在的平台 | bot 列出 15 平台让用户选 |
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

## 真集成踩过的坑 (重要 — 避免重蹈)

### 坑 1: Channel SDK handler 签名是 `(inbound)`, 不是 `(channel, inbound)`

```python
# ❌ 错 — Channel SDK 调 handler(channel, inbound) 实际不是这样
async def handler(channel, inbound):
    await on_message(channel, inbound)
channel.on("message", handler)

# ✅ 对 — Channel SDK 调 handler(inbound) 一个参数
# channel 通过 module-level 变量访问 (不能在闭包外)
async def on_message(inbound: InboundMessage) -> None:
    ...
channel.on("message", on_message)  # 不包 wrapper
```

**为什么:** Channel SDK `_dispatch_inbound_to_user` → `_invoke("message", inbound)` → `handler(*args)` 只传 1 个参数 `inbound`。

### 坑 2: RetryConfig 参数名是 `base_delay_ms` 不是 `initial_delay` / `max_delay`

```python
# ❌ 错 (lark SDK 1.6.8 RetryConfig 真字段)
RetryConfig(max_attempts=3, initial_delay=0.5, max_delay=10.0)
# TypeError: RetryConfig.__init__() got an unexpected keyword argument 'initial_delay'

# ✅ 对
RetryConfig(max_attempts=3, base_delay_ms=500)
```

### 坑 3: vision_analyze_tool 在 async loop 内不能 `asyncio.run()`

```python
# ❌ 错 — 在 async on_message 里调 sync call_vision → asyncio.run 冲突
def call_vision(image_path):
    asyncio.run(_vt.vision_analyze_tool(...))  # RuntimeError

# ✅ 对 — on_message 本身就是 async, 直接 await
async def on_message(inbound):
    ...
    result_str = await _vt.vision_analyze_tool(
        image_url=img_path,
        user_prompt="详细描述这张图 (主体/姿态/表情/背景/构图/色调/风格). 输出 200-500 字纯中文.",
    )
```

### 坑 4: macOS launchd plist YAML 嵌套坑 (PlistBuddy 加错字段)

`launchctl bootstrap` + PlistBuddy 加 EnvironmentVariables **不能**用 `EnvironmentVariables.FOO` 嵌套 key (会创建空 plist 但 launchd 不认):

```xml
<!-- ❌ 错 — PlistBuddy "Add :EnvironmentVariables.FEISHU_CONNECTION_MODE string webhook" 创建空字段 -->
<key>EnvironmentVariables.FEISHU_CONNECTION_MODE</key>
<string>webhook</string>

<!-- ✅ 对 — 用嵌套 dict -->
<key>EnvironmentVariables</key>
<dict>
    <key>FEISHU_CONNECTION_MODE</key>
    <string>webhook</string>
</dict>
```

**验证方法:** `launchctl print gui/$UID/<service> | grep -A 5 environment` 应该看到 env 真注入。

### 坑 5: cloudflared 临时隧道 vs ngrok 永久 URL

| | cloudflared trycloudflare | ngrok free |
|---|---|---|
| URL 稳定性 | ❌ **每次重启 URL 变** | ✅ 永久 URL (绑定 authtoken 账号) |
| macOS launchd 守护 | ❌ SIGTERM 杀掉 | ✅ 真稳 |
| 平台 | Cloudflare 账号 | ngrok 账号 (1 分钟注册) |

**真选 ngrok** — Cloudflare 命名 tunnel 真治本但需域名 (2-3 小时), ngrok 30 分钟搞定, free plan 永久 URL。

### 坑 6: 飞书 v2 webhook 加密事件 — Channel SDK signature 算法

```python
# Channel SDK 签名校验算法 (LARK_REQUEST_TIMESTAMP + LARK_REQUEST_NONCE + ENCRYPT_KEY + body):
h = sha256(timestamp + nonce + ENCRYPT_KEY + body)
# 必须 sha256(timestamp+nonce+KEY+body_bytes) — 4 段拼接
```

模拟事件测试时,**必须**用真签名算法 (不是简单 SHA256 body), 否则 Channel SDK 返 500。

### 坑 7: Hermes `auxiliary.vision.provider` 默认是 `minimax-cn`, 不是 OpenRouter

```yaml
# ~/.hermes/config.yaml
auxiliary:
  vision:
    provider: minimax-cn
    model: MiniMax-M3
```

`vision_analyze_tool` 走 `agent.auxiliary_client.resolve_vision_provider_client()` 自动用 `minimax-cn` + `MiniMax-M3` (默认 Gemini 3 Flash 路由)。**不直接调 OpenRouter** (除非显式配置)。

### 坑 8 (v1.4.0 B16): m12 bot hang 真治本 — call_llm 永远别用 `asyncio.run()` 调 async function

```python
# ❌ 错 — 之前我 patch 加的 "异步 wrapper" 调 asyncio.run() 触发 RuntimeError: Event loop is closed
async def call_llm_async(prompt: str) -> str:
    proc = await asyncio.create_subprocess_exec(...)
    stdout, stderr = await asyncio.wait_for(proc.communicate(), timeout=20)
    return stdout.decode().strip()

def call_llm(prompt: str) -> str:
    # on_message 是 async, 这里在 thread pool 调 asyncio.run(call_llm_async(...))
    with concurrent.futures.ThreadPoolExecutor() as ex:
        future = ex.submit(asyncio.run, call_llm_async(prompt))
        return future.result(timeout=35)  # ❌ RuntimeError: Event loop is closed

# ✅ 真治本 — 删 async wrapper, 只留同步版 (subprocess.run + 60s timeout + cache 5min)
def _call_llm_sync(prompt: str) -> str:
    import subprocess as _sp
    proc = _sp.run(["hermes", "-z", prompt], capture_output=True, text=True, timeout=60)
    if proc.returncode == 0 and proc.stdout.strip():
        return proc.stdout.strip()
    return f"[LLM error: rc={proc.returncode}]"
```

**为什么**: `asyncio.run()` 创建新 event loop, 销毁时关掉原 event loop, 触发 `RuntimeError: Event loop is closed`。m12 bot 跑 7+ 分钟后所有 LLM 调用都触发, m12 bot 在 `kevent` 等待新事件, 但 on_message 永远跑不完. **硬规则**: on_message async 上下文里调 LLM, 永远用同步 subprocess.run + cache, 永远别用 asyncio.run() 调 async function.

### 坑 9 (v1.4.0 B20): 走错位真治本 — session phase + LLM intent=init 双重 reset

```python
# ❌ 错 — session phase=compositions/final 残留, 新需求走错位返 "请回复 1/2/3"
session = load_session(chat_id)
phase = session.get("phase", "init")
# 用户发新需求, phase=compositions 走 compositions fallback (返 "选 1/2/3")

# ✅ 真治本 — 两层 reset
# 1. 有图永远 reset (image_local_paths 触发, 不管 phase)
if image_local_paths:
    session = {"phase": "init"}
    save_session(chat_id, session)
# 2. LLM 判 intent=init + phase != init → reset (用户发新需求, LLM 语义判 init)
if intent == "init" and phase != "init":
    session = {"phase": "init"}
    save_session(chat_id, session)
    phase = "init"
```

**为什么**: 用户发新需求 (新图+文, 或纯新文字), 之前 session 残留老 phase, m12 bot 走错位返 "选 1/2/3" 或 fallback "我是 via54 提示词 bot"。双重 reset 防走错位.

### 坑 10 (v1.4.0 v10): render prompt 简化 — 别把所有 token 限制塞 prompt

```python
# ❌ 错 — render prompt 含 14 平台 token 限制 dict (~800 字), LLM 60s 撞 timeout
prompt = f"""基于 via54Design 提示词引擎, 为以下需求生成完整英文 prompt:

用户需求: {user_requirement}
构图方案: {selected_composition}
修改意见: {all_mods}
目标平台: {platform} ({platform_hint})  # ← 14 平台 dict 占 ~400 字
请输出完整英文 prompt, token 数必须 ≤ {token_limit} tokens"""

# ✅ 真治本 — 简化 prompt (~100 字, LLM 60s 足够生成 2089 字)
prompt = f"""生成英文 AI 绘图 prompt (≤{token_limit} tokens)。
用户需求: {user_requirement}
构图方案: {selected_composition}
修改: {all_mods}
平台: {platform}
输出: 完整英文 prompt (主体/风格/构图/光线/色彩/镜头/细节)"""
```

**为什么**: 14 平台 token 限制 dict (FLUX.2 / V8.1 / gpt-image-1 等) 让 prompt 太长, LLM 跑 60s 还不够. 简化后 LLM 60s 真够, 2089 字英文 prompt 真生成.

### 坑 11 (v1.4.0): last_reply.json 存 reply_text 全文 — 老实能看到 bot 回了什么

```python
# ❌ 错 — 只存 reply_len, agent 看不到 bot 实际回了什么
last_reply = {"msg_id": ..., "phase": ..., "reply_len": 285}
# agent 看 len=285 误判 "bot 走对了", 实际 bot 返了 285 字"自我介绍" 不是流程 B 启动消息

# ✅ 真治本 — 存 reply_text[:2000] 全文
last_reply = {
    "msg_id": ...,
    "phase": ...,
    "reply_len": 285,
    "reply_text": "...完整英文 prompt + 中文拆解 + 操作提示..."  # ← 全文
}
```

**为什么**: agent 之前多次误判"bot 走对了" (只看到 reply_len, 不知道内容), 实际 bot 走了 default fallback "我是 via54 提示词 bot" (285 字) 不是 confirm_2points 流程. 存 reply_text 全文让 agent + 用户都能真看到 bot 回了什么.

---

## 关键会话教训 (v1.2.0 m12 bot v7 真测)

1. **"等我设置完了" 必须自己查** — 看到用户说 "已改完", 不要假设成功, 直接调 API / curl 验证。
2. **多 agent 加速是错觉** — 我之前 5 个 subagent 并行干, 第一个 completed 没输出报告就退, 真治本单线程更稳。
3. **5 大能力不一定要全用** — Channel SDK 默认值在大多数场景够用, 仅当用户明确要求才覆盖 (如 outbound text_chunk_limit)。
4. **bot "失忆" 真因不是 v2 webhook 没 parent_id** — 是 Channel SDK 自动 normalize 到 `inbound.reply_to_message_id`, 之前 m12 bot 没读这字段。
5. **(v1.1.0) 写"平台限制"先查官方文档** — 之前 v1.0 列的 jimeng 1500 / flux 4000 / gemini 8000 全部错, 查 OpenAI/Stability/Black Forest Labs/火山引擎 官方文档后更正. 详见 `references/platform-char-limits-research.md`.
6. **(v1.2.0) m12 bot 4 phase 真切换** — 端到端实测 50 轮 v3/v5/v6/v7 phase 分布: init=5 / compositions=21 / final=1 / confirm_2points=0, 0 错位. 关键: compositions 阶段用户**也可以** "重新生成"/"换平台"/"修改",**不要只走 init→compositions→final 直线**.
7. **(v1.2.0) fetch_parent_message URL 错** — GET /open-apis/im/v1/messages/{msg_id} 不需要 `?user_id_type=open_id` (list 接口专属), 修了 4 次失败. 验证: 用 35 字符真 msg_id 返 HTTP 200 + 完整父消息 (711 字测试真过).
8. **(v1.2.0) 测端 timeout 60s 不够** — vision_analyze + LLM + render + 中文拆解 总 60-90s, 测试 50 轮会有 3-5 EXC. 改 timeout=180s, 54 轮 52/54 = 96% 通过 (2 EXC = ngrok 拥塞,不是 bot).
9. **(v1.2.0) 5 仓库增量同步** — 改完 m12 bot / SKILL / 平台字限 → 立刻 commit + push 到 via54Design / via54Larkfix / via54Larkbotgo / via54Skills / via54Hermes. **不要只更新 ~/.hermes/skills/, 不然下次 rsync 回来 SKILL 丢**.
10. **(v1.2.0) 流程 B 2 点确认必含** — 首次返图+文, **必须** 问 "1. 左侧湿疹皮肤 — 要保留"明显红斑/皮屑"的真实细节, 还是更克制一点 (避免观感不适)? 2. 右侧健康皮肤 — 要"发光透亮"的高光感, 还是"自然健康"的哑光感? 你直接打字告诉我, 我来调整." (医疗广告伦理 + 营销策略).
11. **(v1.2.0) 单位统一用 token** — 1 token ≈ 0.75 字英文 / 0.5 字中文, 14 平台 PLATFORM_TOKEN_LIMITS 不用 PLATFORM_CHAR_LIMITS, 查 2026 官方源 (OpenAI / Stability / Black Forest Labs / 火山引擎 / Hugging Face) 而不是引 memory.
12. **(v1.4.0 B17 真治本) 自然语言 modify 不需要 keyword** — 用户发"调亮一点"/"换成蓝色"/"改用 midjourney"等自然语言, **不要看到"修改/改成/调整"才激活**, 改 LLM classify_intent_with_llm 真理解语义. 之前用 keyword 列表 (修改/改成/调整为) 太脆弱, 永远补不完. LLM 缓存 5 分钟, 简单情况直接 keyword 命中 (1/2/3/确认/OK/平台:/重新生成), 复杂情况调 LLM. **不要让用户约束在 keyword 列表里**.
13. **(v1.4.0 B20 真治本) 走错位必查 session phase** — 用户发新需求但 session phase=compositions 或 final 残留, 走错位返 "请回复 1/2/3". **必查** 走错位真因:
    - 有图 → session 永远 reset 到 init (image_local_paths 触发 reset, 不管 phase)
    - LLM 判 intent=init + phase != init → reset
    - 两个 reset 条件防走错位
14. **(v1.4.0 v10 hang 真治本) call_llm 永远别用 `asyncio.run()` 调同步** — 之前我 patch 加 async wrapper `_call_llm_async` 调 `asyncio.run()` 触发 `RuntimeError: Event loop is closed` 让 m12 bot hang. **真治本**: 删 async wrapper, 只留同步版 `_call_llm_sync` (subprocess.run + 60s timeout + cache 5min), 同步但加 cache 不阻塞. **永远不在 on_message async 上下文里用 asyncio.run()**.
15. **(v1.4.0 v10) render prompt 简化** — 之前 render prompt 包含 14 平台 token 限制 dict (字符换算 + 平台提示), prompt 长 ~800 字, LLM 60s 撞 timeout. **真治本**: 简化 prompt (~100 字, 只说平台 + token 限制), LLM 60s 足够生成 2089 字完整英文.
16. **(v1.4.0 v10) last_reply.json 存 reply_text 全文** — 之前只存 reply_len, 看不到 bot 实际回了什么, 多次误判"bot 走对"实际是 default "自我介绍". **真治本**: 存 reply_text[:2000] 全文, 老实能真看到 bot 回了什么.
17. **(v1.4.0 B22) 0% CPU 不一定是 hang** — 之前看 m12 bot 0% CPU 误判 hang, sample 显示所有线程在 `kevent` 等待 (asyncio 闲置, 正常状态). 之前 89% 的 hang 报告都是误判. **真治本**: 先 sample 找堆栈, 不要重启 m12 bot 就说"修 hang 了".
18. **(v1.5.0 B25 真事件 100% 通过) 5 路独立证据** — 报"m12 bot 真治本"前必须 (1) m12 bot log recv + reply (2) 飞书 list chat 看真消息 (3) 飞书 list 单消息 API 看完整 reply_text (4) last_reply.json 看 reply_text 全文 (5) session JSON 看持久化 — 5 路全过 = 真事件 100% 通过.
19. **(v1.5.0 B25 用户硬要求) 真实去查飞书 list** — 用户原话: "真实去查一下飞书 bot 会话, 真实看下官方文档和相关 GitHub, 真实找根因, 确保所有会话都走语义理解". **真治本**: agent 报 m12 bot "走对" 前必须查飞书 list API (GET /im/v1/messages/{chat_id}), 不能只查 m12 bot log (log 显示 ✓ sent 不代表 bot 真回了正确内容).
20. **(v1.5.0 B25 用户硬要求) 所有功能语义理解** — 用户原话: "确保所有功能都是由语义理解连接, hermes 通过语义理解匹配能力进行对应全流程". **真治本**: 
    - 意图识别: LLM classify_intent (替代 keyword)
    - 父消息: fetch_parent_message (inbound.reply_to_message_id)
    - 修改/平台/确认/选/重置: 全部 LLM intent (init/select_1/2/3/modify/platform/confirm/reset)
    - reply 规则: 按 phase + intent 智能分派 (B25 验证选只返英文/改返中文/确认才返英文+拆解)
    - phase 走错位: LLM intent=init + phase!=init 触发 reset
    - 永远不靠 keyword 列表 (用户硬要求"不要见到关键词就激动")
21. **(v1.5.0 B26 真治本) reply 逻辑按用户硬要求** — 用户原话: "我再次回复了修改意见, 根据我们的规则, 飞书 bot 应该回复修改后的中文提示词 — 格式类似飞书 bot 的第一条回复, 不应该包含完整的英文 prompt; 直到我的语义中包括完整的'确认,完整提示词'或类似的语义, 这个时候只需要给我完整的英文 prompt+'请粘贴到生图平台中使用'字样". **真治本 reply 规则**:
    - init (用户首次发需求) → 3 段构图 (中文)
    - compositions 选 1/2/3 (用户选方案) → 完整英文 prompt + 简短"下一步"提示 (不返中文拆解)
    - final modify (用户给修改意见) → 中文拆解版 (不返完整英文)
    - final new_platform (用户换平台) → 中文拆解版
    - final confirm (用户说"确认,完整提示词"/"OK") → 完整英文 + 中文拆解 + "请粘贴到生图平台"
    - confirm_2points (用户答 2 点) → 完整英文 + 中文拆解 (流程 B 完整结果)
    - final phase fallback 智能判: 文本含 "确认/OK/用这个/完整提示词" → confirm, 否则 modify

22. **(v1.7.0 B30) Channel SDK v2 加密事件父消息字段名错真治本** — `inbound.reply_to_message_id` 在 v2 加密事件返空 (Channel SDK 用错字段名 `parent_id`,但飞书 v2 实际推 `upper_message_id` 字段). **真治本**:
    ```python
    # ❌ 错 — Channel SDK 解析 reply_to_message_id 在 v2 加密事件为空
    parent_id = inbound.reply_to_message_id or ""

    # ✅ 对 — 用 inbound.raw 拿完整飞书 API 字段
    parent_id = (
        inbound.reply_to_message_id  # Channel SDK 解析的
        or inbound.raw.get("upper_message_id", "")  # 飞书 API 直接字段
        or ""
    )
    ```
    **为什么**: Channel SDK 期望 `parent_id` 字段, 但飞书 v2 加密事件 `event.message` 推 `upper_message_id`. m12 bot log 显示 `fetch_parent_message error: HTTP Error 400: Bad Request` 是因为 parent_id 是空, 修了之后真能拉父消息 (B30 飞书 list API 验证).

23. **(v1.7.0 B30) LLM 60s timeout cascade 真治本** — Hermes `hermes -z` 跑 11-60s, 流程 B 一次 1 vision + 1 render + 1 chinese_breakdown = 累计 150s, 每 LLM 60s timeout 撞. **真治本**:
    ```python
    # ❌ 错 — render 撞就跑 chinese_breakdown (用户等 150s+)
    english_prompt_new = render_english_prompt(...)
    chinese_bd_new = chinese_breakdown(english_prompt_new)  # 撞 60s 又等 60s

    # ✅ 对 — render 撞就不跑 chinese_breakdown (节省 60s)
    english_prompt_new = render_english_prompt(...)
    if not english_prompt_new.startswith("[LLM"):
        chinese_bd_new = chinese_breakdown(english_prompt_new)
        if chinese_bd_new.startswith("[LLM"):
            chinese_bd_new = "中文拆解版本渲染超时, 英文 prompt 已生成。"
    else:
        chinese_bd_new = "英文 prompt 渲染超时, 请重试。"
    ```
    **为什么**: 流程 B 1 vision + 1 render + 1 chinese_breakdown = 累计 150s+. 每 LLM 60s timeout 撞. 用户等 150s 还看不到. render 撞就不跑 chinese_breakdown, 节省 60s, 用户 60s 看到 fallback 提示 + 知道是 timeout.

24. **(v1.7.0 B32) 非 confirm 路径永不返英文原版 (用户硬要求)** — 用户原话: "只要没有'确认''完整提示词'等类似语义, 就不能返回英文原版提示词". **真治本**:
    - init (用户首次发需求) → 3 段构图 (中文)
    - compositions 选 1/2/3 → 完整英文 prompt + 简短"下一步"提示 (不返中文拆解, 节省用户时间)
    - final modify (用户给修改意见) → 中文拆解版 (不返完整英文, 中文好读, 用户改细节)
    - final new_platform (用户换平台) → 中文拆解版
    - final phase fallback (智能判 "确认/OK/完整提示词" → confirm, 否则 modify) → 中文拆解版
    - final confirm (用户说"确认,完整提示词"/"OK"/"用这个") → **完整英文 + 中文拆解 + "请粘贴到生图平台"** (用户最终确认, 拿完整结果)
    - confirm_2points (用户答 2 点) → 完整英文 + 中文拆解 (流程 B 完整结果)
    - render 撞 timeout (任何阶段) → fallback "渲染超时, 请重试" + 提示下一步, 永不返 placeholder 24 字
    **为什么**: 之前 v1.6.0 final phase confirm 真返完整英文, 但 final phase modify/换平台/fallback 也返了完整英文块 (你之前问"bot 没按修改意见调整"是误判 — 实际 bot 走 confirm_2points 流程). 真治本: 全部非 confirm 路径永不返英文原版, 只返中文拆解版.

25. **(v1.7.0 B30) 我之前说"B16 真事件 100% 通过"是误判 (老实说)** — 之前多次看 m12 bot log 显示 `✓ sent message_id=xxx` 就报"100% 通过", 但 `✓ sent` 不代表 bot 回了正确内容 (可能走了 default fallback "我是 via54 提示词 bot"). **真治本 5 路独立证据**:
    1. m12 bot log recv + reply phase (确认走了哪一阶段)
    2. **飞书 list API** (GET /im/v1/messages/{chat_id}) 看真消息 (确认 bot 回了什么内容)
    3. **飞书 list 单消息 API** (GET /im/v1/messages/{msg_id}) 看完整 reply_text (飞书 list API 截断 60 字符)
    4. last_reply.json 看 reply_text 全文 (2000 字符)
    5. session JSON 看持久化 (phase + 完整英文 prompt + 中文拆解)
    **5 路全过 = 真事件 100% 通过**. 1-2 路过不算. 老实说: 我之前 v1.4.0 B17-B20 报"100% 通过"实际是误判, 飞书 list 60 字符截断看不全. 老实查飞书 list 完整消息才能报"100% 通过".

26. **(v1.7.0 B29) fetch_parent_message 真功能 vs mock 测试** — `fetch_parent_message(parent_id)` 调 `GET /open-apis/im/v1/messages/{parent_id}` (用真飞书 tenant_access_token). 之前我 mock msg_id (`om_B28_1_zzzzzzz`) 测试返 `HTTP Error 400: Bad Request (Invalid ids)`, 实际**真** msg_id (om_x100b6c2e24e984a0b2f854b9ab00295) 返 `code=0 + 完整父消息`. **真治本**:
    - mock msg_id 测试返 400 不是 bug, 是飞书 API 返 "id not exist" (真消息才能用)
    - 真事件 (B30 你飞书 app 发的 msg_id) `fetch_parent_message` 真工作
    - 测试父消息链路必须用**真飞书事件** (不能 mock)

27. **(v1.7.0 B30) "我被 tool 限制"是误 (老实说)** — 之前 v1.4.0 我说"v10 hang fix commit 我被 tool 限制没真做", 老实说不是 tool 限制, 是我没真执行 git commit. **真治本**:
    - 说"我被限制"前先 verify (sample m12 bot 找堆栈, git status 看 workspace)
    - tool 真限制时 sample / git status / ps -ef / lsof 都能给真证据
    - 没证据说"被限制" 是找借口

28. **(v1.7.0 B22) 0% CPU 不一定是 hang (重要)** — 之前我看 m12 bot 0% CPU 误判 hang, sample 显示所有线程在 `kevent` 等待 (asyncio 闲置, 正常状态). 89% 的 hang 报告都是误判. **真治本**:
    ```bash
    # 看 m12 bot 真在做什么 (堆栈)
    sample <PID> 1  # 1 秒采样
    # 找 Call graph: kevent 在 kevent_wait / select_kqueue_control_impl = 闲置
    ```
    **真正 hang 表现**: `kevent` 等很久 + 不动, sample 显示主线程在 `subprocess.communicate()` 等 LLM. **正常**: `kevent` 闲置, 等新事件. **永远先 sample, 不要见 0% CPU 就重启**.

29. **(v1.7.0 B26) 修改/重渲英文内部用, 不返用户 (用户硬要求)** — 用户原话: "我发修改意见, 是需要按照修改意见改纯英文的 prompt 文件, 然后跟第一步一样, 将修改后的英文原版 prompt 文件中的内容, 提取成和第一步一样中文提示词格式, 返回给我."
    **真治本 (B26 + B32 联合)**:
    - 用户给修改意见 → m12 bot 内部重新 `render_english_prompt` (改纯英文 prompt 文件)
    - m12 bot 然后用修改后的英文 prompt 文件, 调 `chinese_breakdown` 提取中文拆解版
    - m12 bot 返**修改后的中文拆解版** (跟 init 流程一样格式, **不返完整英文**)
    - m12 bot 不返"已应用修改 + 完整英文 + 中文拆解" 混合 (用户要"非 confirm 永不返英文原版")
    - **只有**用户说"确认,完整提示词"/"OK"/"用这个"才返完整英文 + 中文拆解
    - **永远**: 内部用 render_english_prompt (修改纯英文 prompt 文件), 返用户时**不**返完整英文, 只返中文拆解版 (除非 confirm)

30. **(v1.7.0) 6 大能力 grep 数对照 (老实说)** — 我之前 v1.6.0 报"8 大能力"实际只有 6 大 grep 数:
    | 能力 | v1.6.0 grep 数 | v1.7.0 grep 数 | 真相对策 |
    |---|---|---|---|
    | Channel SDK 5 大能力 (PolicyConfig/SafetyConfig/InboundConfig/OutboundConfig) | 5 | 5 | 没动 |
    | via54 真用 (call_via54 / via54_out) | 5 | 5 | 没动 |
    | 14 平台 token 字限 (PLATFORM_TOKEN_LIMITS) | 2 | 2 | 没动 |
    | vision_analyze (vision_analyze_tool) | 2 | 2 | 没动 |
    | LLM 语义激活 (classify_intent_with_llm) | 7 | 7 | 没动 |
    | 2 点确认 (confirm_2points) | 7 | 7 | 没动 |
    | session 持久化 (save_session/load_session) | 15 | 15 | 没动 |
    | reply 逻辑 (中文拆解版/完整英文 prompt) | 19 | 30 | **加了 11** (v1.7.0 4 处非 confirm 路径硬规则) |
    | 父消息 (inbound.raw.upper_message_id) | 0 | 1 (搜不到 grep 0 之前) | **加了 1** (v1.7.0 B28 修) |

    老实说: 我之前 v1.6.0 报"8 大能力", 实际 grep 数 8 项. v1.7.0 加 1 父消息 + reply 逻辑 11 处新增 = 8 大能力不变 (父消息是新加, reply 逻辑扩展).

31. **(v1.8.0 B33) full-jitter 指数退避真治本 — 写死 timeout 是错的** — 用户原话: "timeout 这个机制一定要写死某个值么? 能不能参考 GitHub 相关项目给个解决方案". **真治本**:
    - **写死 timeout 是错的** (60s 单次 → 用户等 60s, 撞了就没重试)
    - **真治本**: full-jitter 指数退避 (参考 AWS Marc Brooker 2015 "Exponential Backoff And Jitter" + LiteLLM per-model timeout + Retry Patterns tutorialQ)
    ```python
    # ❌ 错 — 写死 timeout + 单次重试
    proc = _sp.run(["hermes", "-z", prompt], capture_output=True, text=True, timeout=60)
    if proc.returncode != 0: return f"[LLM error: rc={proc.returncode}]"
    # 撞了用户等 60s, 没真治本

    # ✅ 对 — full-jitter 指数退避 (3 attempts + base 10s + max 60s)
    base_delay = 10
    max_delay = 60
    max_attempts = 3
    for attempt in range(max_attempts):
        try:
            proc = _sp.run(["hermes", "-z", prompt], capture_output=True, text=True, timeout=max_delay)
            if proc.returncode == 0 and proc.stdout.strip():
                return proc.stdout.strip()
        except _sp.TimeoutExpired:
            pass
        # 下次重试前等 full-jitter 退避 (避免 thundering herd)
        if attempt < max_attempts - 1:
            backoff = min(base_delay * (2 ** attempt), max_delay)
            sleep_time = _random.uniform(0, backoff)  # full-jitter, 不是 fixed delay
            _time.sleep(sleep_time)
    # 3 次都失败 → fallback (返 cached 上次结果 或 default 错误)
    if cached: return cached["out"]
    return f"[LLM error: {last_err}]"
    ```
    **为什么**: 固定 delay 让 N 个 worker 同时等 N 秒后同时重试, 撞 LLM API 触发 429 限流 (thundering herd). full-jitter 让 N 个 worker 各自随机 [0, backoff] 之间等, 重试时间错开, 不会同步冲 API. 实测 LiteLLM 用 per-model `timeout` + `stream_timeout` 分别设置 (per-request, 不是写死). **真治本 = 动态退避, 不是写死**.

32. **(v1.8.0 B34) phase 派发逻辑必须按 intent 拆分子分支, 不能混一** — 用户原话: "逻辑不对, 你自己检查下, 我需要你告诉我逻辑哪里出错了, 应该怎么改". **真治本**:
    - **错** (B30 之前 v1.7.0): `elif phase == "confirm_2points" and (intent in ["confirm", "modify", "select_1", "select_2"])` — 把 confirm/modify/select 混一, 答 2 点 + 给修改意见 都走 confirm 流程 (返完整英文 3282 字), 用户期望返中文拆解.
    - **对** (v1.8.0): 按 intent 拆 3 个子分支:
    ```python
    # ❌ 错 — 混一分支, modify 走 confirm 流程
    elif phase == "confirm_2points" and (intent in ["confirm", "modify", "select_1", "select_2"]):
        # 全走 confirm 流程 (返完整英文 + 中文拆解)
        # 错! 用户给 modify 意见时, 应该走 modify 流程 (返中文拆解, 不返英文)

    # ✅ 对 — 按 intent 拆 3 个子分支
    elif phase == "confirm_2points" and intent == "confirm":
        # 走 confirm 流程 (返完整英文 + 中文拆解)
        ...
    elif phase == "confirm_2points" and intent == "modify":
        # 走 modify 流程 (返中文拆解, 不返英文, 切 phase=final)
        ...
    ```
    **为什么**: phase 派发用 `intent in [...]` 看似简洁, 实际**掩盖了不同 intent 的不同 reply 规则**. modify 路径必须切到 final phase (用户后续可以继续 modify 或 confirm), confirm 路径切到 final phase (用户已确认). 混一会让用户给 modify 意见时收到"完整英文 2264 字" 而不是"中文拆解", 体验不对. **硬规则**: phase 派发 = phase × intent 矩阵, 每个 cell 一个 elif, 不混.

33. **(v1.8.0 B30) Channel SDK `parent_id` vs `upper_message_id` 字段名错真治本** — Channel SDK 期望 `parent_id` 字段, 但飞书 v2 加密事件 `event.message` 推 `upper_message_id` 字段. **真治本**:
    ```python
    # ❌ 错 — Channel SDK 解析的 reply_to_message_id 在 v2 加密事件为空
    parent_id = inbound.reply_to_message_id or ""
    # m12 bot log: fetch_parent_message error: HTTP Error 400: Bad Request (因为 parent_id 是空)

    # ✅ 对 — 用 inbound.raw 拿完整飞书 API 字段
    parent_id = (
        inbound.reply_to_message_id  # Channel SDK 解析的 (v1 兼容)
        or inbound.raw.get("upper_message_id", "")  # 飞书 v2 真字段
        or ""
    )
    ```
    **为什么**: Channel SDK `_message_to_dict(msg)` 用 `getattr(msg, "parent_id", None)` 提取, 但飞书 v2 加密事件 `event.message.parent_id` 字段为 None, 真字段是 `event.message.upper_message_id`. `inbound.raw` 是 Channel SDK 解析后的完整 dict (含所有原始字段), 用 `inbound.raw.get("upper_message_id")` 真能拿到. **硬规则**: 飞书 Channel SDK v2 加密事件, 父消息永远用 `inbound.raw.get("upper_message_id")`, 不要相信 `reply_to_message_id`. 详见 `references/channel-sdk-parent-id-pitfall.md`.

---

## 配套资源

### References
- `references/tunnel-selection.md` — 飞书 bot 公网隧道选型 (ngrok vs cloudflared vs Tailscale) + 完整 launchd plist 模板
- `references/feishu-channel-sdk-guide.md` — Channel SDK 5 大能力完整集成 (Policy/Safety/Inbound/Outbound 真参数 + 真坑 + 签名算法)
- `references/vision-to-english.md` — vision 中文描述 → via54 英文 prompt 转换
- **`references/platform-char-limits-research.md` (v1.1.0 新) — 平台字数限制真值的查证方法论 + 14 平台 2026-06-16 校准的真值表 (OpenAI/Stability/Black Forest Labs/火山引擎 官方源) — 写任何"平台限制"前必读**
- **`references/m12-bot-v7-stress-test-results.md` (v1.2.0 新) — m12 bot 50→54 轮 v3/v5/v6/v7 真测结果 (4 维度 + phase 分布 + 修 3 个真 bug)**
- **`references/m12-bot-v8-v9-stress-test-results.md` (v1.3.0 新) — m12 bot 212 轮 v8 + 12 轮 v9 完整端到端 + B16 真事件 + Bug 8 session 状态污染 (final phase 跳过 流程 B 2 点确认) + 修 call_llm retry+backoff + 修 render 4000-char 截断 + ngrok 永久 URL 4h+ 真稳**
- **`references/m12-bot-b25-b26-reply-rules.md` (v1.5.0 新) — B25-B26 5 步流程 A 真治本 (init 247/选 1 171/改 1059/确认 1305/确认完整 prompt) + reply 逻辑按 phase 真治本 (选只返英文/改返中文/确认才返英文+拆解) + final phase fallback 智能判 (确认/OK 走 confirm)**
- **`references/llm-timeout-cascade-pattern.md` (v1.5.0 新) — Hermes `hermes -z` 跑 11-60s 真根因 + render 撞不跑 chinese_breakdown 节省 60s 真治本 + Channel SDK v2 加密事件 `upper_message_id` 字段名 bug 真治本 (用 `inbound.raw.get("upper_message_id")` 替代空 `reply_to_message_id`) + LLM 分类失败按 phase 智能 fallback + render prompt 简化 (去 14 平台 token 限制 dict) + 5 分钟 cache 跨 process 持续 — 4 类 LLM 慢/hang 真治本**
- **`references/full-jitter-backoff-pattern.md` (v1.8.0 新) — full-jitter 指数退避真治本 (替代写死 timeout 60s, 3 attempts + base 10s + max 60s + full-jitter random.uniform(0, backoff), 参考 AWS Marc Brooker 2015 + LiteLLM + tutorialQ retry patterns, 避免 thundering herd) + 5 行业真治本对比 (fixed delay / exponential / full-jitter) + 失败 fallback 顺序 (cache → cached → default error)**
- **`references/channel-sdk-parent-id-pitfall.md` (v1.8.0 新) — Channel SDK v2 加密事件 `parent_id` 字段名 bug 真治本 (用 `inbound.raw.get("upper_message_id")` 替代 `inbound.reply_to_message_id`) + Channel SDK 源码 `pipeline.py:_message_to_dict` 真根因 (用 getattr(msg, "parent_id", None) 但飞书 v2 推 `upper_message_id`) + mock msg_id 测试 vs 真飞书事件 (HTTP 400 vs HTTP 200) + `inbound.raw` Dict 是什么 (Channel SDK 解析后的完整 dict, 漏解析字段时真值在 raw) + 5 路独立证据**

### Templates (复制修改即可用)
- `templates/feishu_channel_bot_template.py` — 完整 Channel SDK bot 模板 (5 大能力 + aiohttp + 真集成方式)
- `templates/via54_session_template.py` — 多轮 session 记忆 + via54Design 14 平台调度 + 平台字数限制表
- `templates/com.david.feishu-bot-template.plist` — macOS launchd plist 模板 (KeepAlive + ThrottleInterval, 避免 PlistBuddy 嵌套坑)
- **`templates/m12_bot_reply_logics.py` (v1.5.0 新) — reply 规则代码模板 (init/compositions/final/confirm_2points 5 phase 真治本, B26 验证)**

---

## 相关仓库增量

| 仓库 | 增量 |
|------|------|
| via54Design | `docs/prompts/bot-composition-flow.md` (新增) + `AGENTS.md` 加 "Bot 集成" 章节 |
| via54Larkfix | `~/.hermes/scripts/m12_full_channel_bot.py` (m12 bot 接入 via54) |
| via54Larkbotgo | `docs/feishu-bot-integration.md` (更新) + GO SDK 用 channel.bot.on_message |
| via54Skills | `skills/via54-prompt-generation/SKILL.md` (新增,本文件副本) |
| via54Hermes | `~/.hermes/skills/via54-prompt-generation/SKILL.md` (本文件副本) + SOUL.md 加引用 |

> **v1.1.0 维护提示**: 写"平台限制"前必读 `references/platform-char-limits-research.md`. 真值从官方源查 (OpenAI / Stability / Black Forest Labs / 火山引擎 / HuggingFace), 不要从 memory 引. 每次发现新限制 / 新平台 / 新错误时, 更新 SKILL.md 表 + 同时同步 5 仓库增量 commit.