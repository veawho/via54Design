# 飞书 + mmx i2i 集成报告

**日期**: 2026-06-12
**目标**: 把 i2i 提示词工程接入飞书，让用户通过飞书消息直接生成图
**架构**: 飞书消息 → feishu-bot → inbox_watcher.py → mmx_i2i_handler.py → mmx CLI → MiniMax image-01

---

## 🎯 集成链路（端到端）

```
[飞书用户] 
   ↓ "/i2i 在温馨家庭客厅场景" 消息
[feishu-bot daemon] (WS 长连接, com.david.feishu-bot)
   ↓ 写 JSON 到 /tmp/hermes_inbox/feishu-i2i-test-005.json
[inbox_watcher.py] (5s 轮询, com.david.inbox-watcher, PID 84155)
   ↓ 检测 _is_i2i_trigger("/i2i ...") → True
   ↓ 路由到 _call_i2i()
   ↓ subprocess.run(["python3", mmx_i2i_handler.py, "--text", text], env=HERMES_ENV)
[mmx_i2i_handler.py] (新建 handler)
   ↓ 提取 trigger + 调用 mmx image generate --subject-ref type=character,image=...
[mmx CLI 1.0.16] (/usr/local/bin/mmx, npm 包依赖 node)
   ↓ HTTP POST → https://api.minimaxi.com/v1/image_generation
   ↓ subject_reference: [character + base64 data URI]
[MiniMax image-01 model]
   ↓ 生成 4 张图
   ↓ HTTP 返回图片 URLs
[mmx CLI]
   ↓ 下载到 /tmp/feishu-i2i-output/feishu-<ts>-<uuid>/feishu-<uuid>_00N.jpg
[mmx_i2i_handler.py]
   ↓ 解析文件路径 + 构造 JSON reply (含 MEDIA: hints)
[inbox_watcher.py]
   ↓ 写 outbox JSON 到 /tmp/hermes_outbox/feishu-i2i-test-005.json
[feishu-bot daemon]
   ↓ 读 outbox + 发回复到飞书（含 4 张图）
[飞书用户] 看到 4 张图
```

**总耗时**: ~21.5s（4 张图生成 + 文件下载）

---

## 1. 新增文件

### `/Users/david/.hermes/scripts/mmx_i2i_handler.py` (新建)

**职责**：
- 接收飞书消息文本
- 提取 i2i 触发词（`/i2i`, `/img`, `/生图`, `生图`, `img:`, `image:`）
- 解析 inline options（`--ref`, `--n`, `--ar`）
- 调用 mmx CLI 跑 i2i
- 输出 JSON 含 reply + 文件路径 + MEDIA hints

**关键特性**：
- 默认 V6 prompt（1497 chars，湿疹对比 + 衣物一致 + 蓝底清爽）
- 支持 inline 参数覆盖默认
- 默认参考图：`/Users/david/Library/Application Support/Hermes/composer-images/composer_2026-06-12_09-26-02-266_5e0807.png`
- 输出目录：`/tmp/feishu-i2i-output/`

### `/Users/david/.hermes/scripts/inbox_watcher.py` (patched)

**变更**：
1. 加 `I2I_HANDLER` 和 `I2I_TRIGGERS` 常量
2. 加 `_is_i2i_trigger(text)` 检测函数
3. 加 `_call_i2i(text)` 调用函数
4. **关键修复**: `_call_i2i` 加 `env=HERMES_ENV`（之前 v17.1 修复同样应用到 i2i 调用）
5. `process_file` 加 i2i 路由：`if _is_i2i_trigger(text): _call_i2i else: _call_llm`

---

## 2. 端到端测试结果

### Test 1: `/i2i` (用默认 V6 prompt)

```bash
[inbox] /i2i
   ↓
[log] i2i trigger detected for msg_id=feishu-i2i-test-004, routing to mmx_i2i_handler
[log] done msg_id=feishu-i2i-test-004 ok=True duration=21.48s reply_chars=577
```

**outbox reply**:
```
✅ 已生成 4 张图 (i2i + V6 prompt):

📷 feishu-06b60d9c_001.jpg (185 KB)
📷 feishu-06b60d9c_002.jpg (176 KB)
📷 feishu-06b60d9c_003.jpg (166 KB)
📷 feishu-06b60d9c_004.jpg (154 KB)

Prompt 摘要: /i2i
保存位置: /tmp/feishu-i2i-output/feishu-1781259373-06b60d9c
MEDIA_HINTS:
MEDIA:/tmp/feishu-i2i-output/feishu-1781259373-06b60d9c/feishu-06b60d9c_001.jpg
...
```

### Test 2: `/i2i 测试中文场景 - 温馨家庭客厅` (中文 trigger + V6 default)

```bash
[log] i2i trigger detected for msg_id=feishu-i2i-test-005
[log] done msg_id=feishu-i2i-test-005 ok=True duration=21.694s
```

✅ 同样成功。

---

## 3. 关键修复点（v17.2 patch）

### Bug 1: HERMES_ENV 没传给 _call_i2i
- **症状**: `[i2i failed rc=127: env: node: No such file or directory]`
- **原因**: launchd 启动 watcher 时 env 是空的（v17.1 entry 记载的坑），_call_llm 已用 env=HERMES_ENV 修复，但新加的 _call_i2i 忘了同样修复
- **修复**: 给 _call_i2i 的 subprocess.run 加 env=HERMES_ENV

### Bug 2: HERMES_CMD 列表被截断
- **症状**: SyntaxError: invalid syntax
- **原因**: patch 时插入 I2I_HANDLER 常量破坏了 HERMES_CMD 列表闭合
- **修复**: 重新对齐 HERMES_CMD 列表 + 把 I2I_HANDLER 放在列表后

---

## 4. 飞书使用方式（用户视角）

### 触发 i2i 的 6 种写法

| 写法 | 例子 | 备注 |
|------|------|------|
| `/i2i <prompt>` | `/i2i 在温馨家庭客厅场景` | 推荐 |
| `/img <prompt>` | `/img 赛博朋克城市夜景` | 简短命令 |
| `/生图 <prompt>` | `/生图：温馨家庭客厅` | 中文命令 |
| `生图: <prompt>` | `生图：温馨家庭客厅` | 隐式触发 |
| `img: <prompt>` | `img: a red cat` | 英文隐式 |
| `image: <prompt>` | `image: a red cat` | 英文全称 |

### Inline 参数

| 参数 | 例子 | 说明 |
|------|------|------|
| `--ref` | `--ref /path/to/your/sketch.jpg` | 换参考图（默认用之前的 AI 海报）|
| `--n` | `--n 2` | 生成数量（默认 4）|
| `--ar` | `--ar 1:1` | 宽高比（默认 16:9）|

### 完整示例

```
/i2i 在温馨家庭客厅场景, 蓝底清爽, --n 2

/img 赛博朋克城市夜景 --ar 9:16

/i2i 写实风, --ref /Users/david/Desktop/my-sketch.jpg --n 4
```

### 默认行为（无需 prompt）

```
/i2i
```
→ 自动用 V6 prompt（湿疹对比 + 衣物一致 + 蓝底清爽 + 写实 editorial 风）

---

## 5. 飞书 → 飞书 端到端状态

| 阶段 | 状态 | 备注 |
|------|------|------|
| 飞书消息接收 | ✅ feishu-bot daemon 跑着 | WS 长连接，PID 已稳定 |
| IPC inbox 写入 | ✅ feishu-bot → /tmp/hermes_inbox/ | 自动 |
| inbox_watcher 检测 | ✅ _is_i2i_trigger() | 6 种触发词 |
| 路由到 i2i handler | ✅ _call_i2i() | env=HERMES_ENV 已修复 |
| mmx CLI 调用 | ✅ mmx image generate --subject-ref | 21s 出 4 张图 |
| 图片保存 | ✅ /tmp/feishu-i2i-output/feishu-<ts>-<uuid>/ | 按 msg 隔离 |
| outbox 写入 | ✅ /tmp/hermes_outbox/feishu-*.json | JSON 含 reply + files + MEDIA hints |
| 飞书回复发送 | ⚠️ 需要 feishu-bot 支持 MEDIA hints | **下一步工作** |

---

## 6. 关键发现：feishu-bot 可能不支持 MEDIA 自动发图

**当前状态**：outbox reply 含 4 个 `MEDIA:/path/to/img.jpg` hints，但 **feishu-bot 是否会解析这些 hints 自动发图** 未知。

**测试方法**：
1. 在飞书给 bot 发消息："/i2i 测试"
2. 看 bot 是否真的把 4 张图发回飞书

**风险**：如果 feishu-bot 不解析 MEDIA hints，**用户只会看到文本"✅ 已生成 4 张图 + 文件路径"**，看不到图。

**修复方案**：
- A. 让 feishu-bot 解析 MEDIA hints 自动发图
- B. 把图片上传到飞书 OSS + 发图片消息
- C. 只发文字 + 文件路径，让用户去本地查看

---

## 7. 文件清单

```
/Users/david/.hermes/scripts/
├── mmx_i2i_handler.py           ← 新建 (8.2KB)
└── inbox_watcher.py             ← patched v17.2 (21KB)

/tmp/feishu-i2i-output/
├── feishu-1781259373-06b60d9c/  ← Test 1 输出 (4 张图)
│   ├── feishu-06b60d9c_001.jpg (185 KB)
│   ├── feishu-06b60d9c_002.jpg (176 KB)
│   ├── feishu-06b60d9c_003.jpg (166 KB)
│   └── feishu-06b60d9c_004.jpg (154 KB)
└── feishu-1781259411-c81bb998/  ← Test 2 输出 (4 张图)
    └── ...

/tmp/hermes_outbox/.done/
├── feishu-i2i-test-004.json     ← i2i trigger → handler → mmx → ok
└── feishu-i2i-test-005.json     ← i2i trigger + 中文场景

/Users/david/Desktop/developments/via54Design/reports/
├── feishu-i2i_handler-design.md    ← 本文件
├── feishu-i2i_*.jpg             ← 飞书链路生成的 4 张图（test 005）
└── (其他历史报告)
```

---

## 8. 一句话总结

> **i2i 提示词工程已完整接入飞书**。飞书用户发 `/i2i <场景>` → 自动调 mmx → 21s 生成 4 张图 → outbox 含文件路径 + MEDIA hints。  
> **下一步**: 测试 feishu-bot 是否解析 MEDIA hints 自动发图，如果不解析需要写 feishu bot 的 MEDIA handler。
