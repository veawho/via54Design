# 飞书 bot 集成指南 (2026-06-15)

via54Design 作为"提示词生成引擎", 被飞书 bot (外部 Python 进程) 通过 CLI 调用。

## 调用方式

飞书 bot 跑 `feishu_ref_prompt_handler.py` (Python) 在 `~/.hermes/scripts/`, **通过 subprocess 调 via54 CLI**:

```bash
# 流程 A (纯文字): 用户中文需求
~/.local/bin/via54 prompt i2i \
    --scene "我想做一张湿疹主题的年轻女性广告海报" \
    --platform jimeng

# 流程 B (图+文字): vision 描述 + 用户需求
~/.local/bin/via54 prompt i2i \
    --scene "[参考图 vision 描述]: ...\n[用户需求]: 帮我改成清新绿色" \
    --platform midjourney
```

## ⚠️ 重要: via54 1.0.34 限制

**`prompt i2i` 1.0.34 是"格式化引擎"不"生成英文 prompt"**:
- 输入 `--scene "中文需求"`, **返** `FinalEN: "中文需求"` (原样返, **不**做翻译/扩写)
- 不满足"飞书 bot 需要 26 维度细节拉满英文 prompt"需求

**飞书 bot 实际用法 (v2.10.8.1)**:
- 跳过 via54 CLI 调
- 改用 LLM (`hermes -z`) 真生成 26 维度英文 prompt
- via54 CLI 仍保留作未来模板格式化/参考图分析等用

## 飞书 bot 代码位置

- Python handler: `~/.hermes/scripts/feishu_ref_prompt_handler.py` (Hermes 私有 scripts/, **不在 via54Design 仓库**)
- 飞书 bot daemon: `~/.hermes/scripts/feishu_bot_daemon.py`
- 飞书 inbox watcher: `~/.hermes/scripts/inbox_watcher.py`
- via54Larkfix 仓库 (飞书 bot mirror): https://github.com/veawho/via54Larkfix

## via54 CLI 命令参考

```bash
via54 prompt i2i --scene "..." --platform <platform> [--ref ref.jpg] [--ref-desc "..."] [--max-chars 1500]
via54 prompt list                    # 列所有平台
via54 ref --image ref.jpg             # 参考图分析
via54 gallery                         # 提示词模板市场
```

## 4 个相关 issue

1. **via54 1.0.34 不生成英文 prompt** — 飞书 bot 改用 LLM 绕过
2. **Go MCP server 没编译** — 飞书 bot 不用 MCP, 用 CLI 即可
3. **via54 v2.3 引擎 (via54_i2i_handler.py) 跟 v2.10.8 飞书 bot 逻辑不兼容** — 飞书 bot 不再调 v2.3 引擎
4. **lark SDK 1.6.8 wss 长连接死锁** — 跟 via54Design 无关, 飞书 bot 自己的 daemon 死循环

## 改进建议 (未来)

- via54Design 应加 `--generate-english` 模式, 真正从中文 scene 生成 26 维度英文 prompt
- via54Design 当前 v0.4.0 + v1.0.34 binary 仍可用作 YAML 模板格式化, 但飞书 bot 主流程不依赖
