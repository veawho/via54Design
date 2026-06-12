# 🚨 紧急行动指南: Pexels API Key 仍公开 (需用户介入)

> **当前状态 (2026-06-12 17:00)**: 本地 + history 已净化, 但 **GitHub 远端 main 仍含 Pexels KEY** (git push 反复失败)
> **唯一彻底方案**: 用户手动删 repo + 重建 + force push
> **时间**: ~5 分钟 (用户操作)
>
> **★ 旧 KEY 内容已 redact, 完整 KEY 备份在用户本地 `C:\Users\via54\.env.personal` ★**

---

## ⚠️ 关键事实

| 项 | 状态 |
|---|---|
| 本地工作区 | ✅ 干净 (3 fetch_*.py 用 env 变量) |
| 本地 Git history (151 commits) | ✅ 干净 (filter-repo 净化) |
| **GitHub 远端 main** | ❌ **仍含 Pexels KEY** (521MB, 旧 SHA) |
| **GitHub Code Search** | ❌ **命中 3 个文件** (公开可见) |
| git push 尝试 | ❌ 失败 (size 521MB, 反复 timeout) |
| **Pexels 旧 KEY** | ⚠️ **仍有效** (用户应立刻 revoke) |

---

## 🛑 紧急操作 (用户立刻做)

### 1️⃣ 撤销 Pexels 旧 KEY (★ 最重要 ★)

- 登录 https://www.pexels.com/api/
- 找到 "Your API Key" 列表
- **点 "Regenerate" 或 "Delete"** 旧 key (56 字符, `PEXELS_KEY_REDACTED`, 本地备份在 `C:\Users\via54\.env.personal`)
- 旧 key 立即失效, 攻击者无法再用

### 2️⃣ 检查 Pexels 用量 (看是否被滥用)

- 登录 Pexels Dashboard
- 查看 "API Calls" 统计
- 如果有大量陌生请求, 警惕

---

## 🔥 彻底修复 (用户手动, 5 分钟)

### 方案 A: 删 repo + 重建 (★ 最干净 ★)

#### 步骤:

**a) GitHub UI: 删仓库**

1. 打开 https://github.com/veawho/via54Design/settings
2. 滚到最下 "Danger Zone"
3. 点 "Delete this repository"
4. 输入仓库名确认: `veawho/via54Design`
5. 输入密码 + 2FA 验证
6. 等待 1-2 分钟, 仓库消失

**b) GitHub UI: 重建空仓库**

1. https://github.com/new
2. Owner: veawho, Name: via54Design
3. 选择 Public/Private (按原样)
4. **不要**勾选 "Initialize with README" (留空)
5. 创建

**c) 本地 force push**

```bash
cd G:/agent/developments/via54Design

# 杀之前的 git push 进程
taskkill /F /IM git.exe 2>&1

# 推 (空 repo 接受 force push)
git push origin main --force
```

#### 预计时间:
- 删 repo: 1 分钟
- 重建: 30 秒
- push 521MB: 3-5 分钟 (空 repo 接受)
- 总计: **5-7 分钟**

---

### 方案 B: 不删, 接受风险

如果不想重建:
- **Pexels 旧 KEY 已 revoke** → 攻击者拿到 key 也不能用
- **风险**: 攻击者会知道你用过 Pexels (信息公开, 不会太敏感)
- **代码 history 仍含 key** → 任何人 clone 旧 commit 仍能看到 (但 key 已失效)

**★ 推荐: 至少做方案 B + revoke KEY ★**

---

## 📝 之后操作 (★ 推送后 ★)

### 更新远端后:

```bash
# 1) 确认 Pexels KEY 替换生效
curl -s -H "Accept: application/vnd.github+json" \
  https://api.github.com/repos/veawho/via54Design/contents/_scripts/fetch_market_v2.py | \
  python -c "import json, sys, base64; d = json.load(sys.stdin); c = base64.b64decode(d['content']).decode(); print('OLD KEY' if 'aHyfRPK9' in c else 'OK: env var')"

# 应输出: "OK: env var"
```

### 启用 GitHub Secret Scanning:

1. https://github.com/veawho/via54Design/settings/security_analysis
2. 启用 "Secret scanning" (免费, GitHub 自动扫)
3. 启用 "Push protection" (阻止含 secret 的 push)

### 装 gitleaks (本地 pre-commit):

```bash
# 用 winget 装
winget install gitleaks

# 或 pip
pip install pre-commit
# 然后在 .pre-commit-config.yaml 加 gitleaks hook
```

---

## 🔍 给用户的简易诊断命令

随时可跑, 看远端是否仍含 Pexels KEY:

```bash
# PowerShell
(Invoke-WebRequest -Uri "https://raw.githubusercontent.com/veawho/via54Design/main/_scripts/fetch_market_v2.py" -UseBasicParsing).Content | Select-String "PEXELS"

# bash (curl + grep)
curl -s https://raw.githubusercontent.com/veawho/via54Design/main/_scripts/fetch_market_v2.py | grep -E "PEXELS"

# GitHub API
curl -s "https://api.github.com/search/code?q=repo:veawho/via54Design+PEXELS" | jq .total_count
```

期望输出: `0` (KEY 已清) 或 `3` (KEY 仍公开)

---

## 📅 时间线

| 时间 | 事件 |
|---|---|
| 16:00 | 用户: "确保没有 api token 涉密残留" |
| 16:05 | 发现 Pexels KEY 3 文件 + 3 commit |
| 16:15 | 用户选 A: 替换 + history purge + force push |
| 16:20 | 替换 3 文件 (env 变量) |
| 16:25 | 装 git-filter-repo + 跑 history 重写 (151 commits) |
| 17:00 | **本指南**: 用户手动删 repo + 重建 |

---

**★ 立即行动 ★**: **撤销 Pexels 旧 KEY (30 秒) → 删 repo + 重建 (5 分钟) → 之后 push 干净**

**Last updated**: 2026-06-12 17:00
**Status**: 本地+history 修复完成, 远端清理待用户操作
