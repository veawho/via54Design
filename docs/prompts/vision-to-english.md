# Vision 中文描述 → via54 英文 Prompt 转换

> 流程 B 的关键环节:把 vision_analyze_tool 返的中文描述 + 用户文字需求 → 通过 via54Design 渲染完整英文 prompt。

---

## 输入

| 来源 | 例子 |
|------|------|
| `vision_analyze_tool` 返中文 | "一幅采用黑白线描风格绘制的横向宣传图,主色调为灰白黑三色,左侧女性焦虑+湿疹..." |
| 用户文字需求 | "摄影写实,底色色调为蓝色" |

---

## 转换逻辑

### Step 1: 中文描述 → 结构化 prompt 参数

把 vision 中文描述解析成 via54 prompt `--scene` 的结构:

```
<scene 内容>
1. 主体 (Subject): <vision 提取>
2. 风格 (Style): <用户文字>
3. 色调 (Color): <vision 提取 + 用户修改>
4. 构图 (Composition): <vision 推断>
5. 光线 (Lighting): <vision 推断>
6. 背景 (Background): <vision 提取>
```

### Step 2: 调 via54 CLI

```bash
~/.local/bin/via54 prompt \
  --scene "1. 主体: <parsed subject>
2. 风格: <user_style>
3. 色调: <parsed color>
4. 构图: <inferred composition>
5. 光线: <inferred lighting>
6. 背景: <parsed background>" \
  --platform <user-specified or "通用">
```

### Step 3: via54 返英文 prompt

via54 内部走 14 平台 YAML 模板,把中文结构化描述 → 英文 prompt,带:
- 平台字数限制 (jimeng ≤1500, midjourney ≤4500)
- 平台特定关键词 (midjourney `--style raw --v 6`, flux 自然语言)
- 26 维度自动展开 (subject/secondary/style/artist/medium/pose/hair/genre/shot/composition/dof/view/format/lighting/color/etc)

### Step 4: LLM 中文拆解版

bot 拿英文 prompt → 调 LLM (`hermes -z`) 简化 + 中文拆解:

```
### 主体
- 主体 (Subject): <英文 prompt 第 1 句>
- 辅助 (Secondary): <英文 prompt 第 2-3 句>

### 风格
- 风格流派: <英文 prompt 里的 "Photorealistic" / "Cinematic" 等>
- 媒介: digital photography

### 构图
- 镜头: medium shot
- 构图: rule of thirds

### 光线
- 光线: cinematic lighting

### 色彩
- 色调: cool blue grading (按用户要求)
```

---

## 代码实现 (m12 bot 内)

```python
def vision_to_english_prompt(vision_zh: str, user_text: str, platform: str = "通用") -> str:
    """调 via54Design 把 vision 中文描述 + 用户文字 → 英文 prompt"""
    # 1. 解析 vision 中文 → 结构化
    scene = f"""
1. 主体: {vision_zh[:300]}
2. 风格: {user_text or '保持参考图风格'}
3. 色调: 从参考图推断
4. 构图: 从参考图推断
5. 光线: 从参考图推断
6. 背景: {vision_zh[300:600]}
"""
    # 2. 调 via54 CLI
    import subprocess
    result = subprocess.run(
        ["/Users/david/.local/bin/via54", "prompt",
         "--scene", scene.strip(),
         "--platform", platform],
        capture_output=True, text=True, timeout=60,
        env={**os.environ, "PATH": "/Users/david/.local/bin:/usr/bin:/bin"},
    )
    if result.returncode == 0:
        return result.stdout.strip()
    # 失败 fallback: 用 LLM 直接翻译
    return call_llm(f"请把以下中文描述翻译成完整的英文 prompt (适合 AI 绘图工具):\n\n{vision_zh}\n\n用户文字需求: {user_text}")
```

---

## 增量仓库

| 仓库 | 增量 |
|------|------|
| via54Design | `docs/prompts/vision-to-english.md` (新增,本文件) |
| via54Larkfix | `m12_full_channel_bot.py` 加 `vision_to_english_prompt()` 函数 |
| via54Larkbotgo | `docs/vision-flow.md` (新增) |
| via54Skills | `skills/via54-prompt-generation/references/vision-to-english.md` (本文件副本) |