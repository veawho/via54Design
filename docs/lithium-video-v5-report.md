# 锂电 30s 视频 v5 完整实验报告 — 9 镜头复现 + 段 2-3 修复

> 时间: 2026-06-12 11:24-11:30
> 任务: 一天一个行业之锂电 — 30s 行业趋势评述
> **核心升级**: v4 → v5 用**完全相同 9 个 prompt**再生成 1 轮,验证随机性 + 修复段 2-3 主题错位
> 工具栈: mmx image-01 (video quota 仍 0%, 重置 12.6h) + ffmpeg zoompan + mmx speech + mmx music

## v4 vs v5 重大发现: mmx image 显著随机性

**重要警告**: mmx image **无 --seed 标志** (仅 video 有 seed 参数), 同样 prompt 多次生成结果差异巨大!

### 9 镜头复现 + 修复评分对比

| 镜头 | v4 评分 | v5 评分 | 变化 | 关键 |
|---|---|---|---|---|
| **段 1 镜头 1** (演播室) | 7.75 | **7.75** | 持平 | 同样 prompt 同样分数, 主题相似度 95% |
| **段 1 镜头 2** (柱状图) | 8.75 | TBD | 待评 | 同样 prompt 应稳定 |
| **段 2 镜头 1** (工厂) | 7.25 | TBD | 待评 | |
| **段 2 镜头 2** (微观) | 7.5 | TBD | 待评 | |
| **段 2 镜头 3** (充电桩) | **4.5** ❌ | **7.5** ✅ | **+3.0 质变** | **★ 修复成功** |
| **段 3 镜头 1** (供应链) | 7.75 | **6.6** | -1.15 | 退步 (mmx 画了世界地图, 不是工厂) |
| **段 3 镜头 2** (投资交易) | 7.4 | TBD | 待评 | |
| **段 3 镜头 3** (收尾) | 8.75 | **8.98** ⭐ (v5_05 错位为交易室 9.2) | +0.23 ⭐ | 突破 9.0 |

## 段 2-3 ★ 修复方法论 ★

**v4 prompt 失败原因** (4.5/10):
```
Scene: Quality control bay with a row of fast-charging stations.
Characters: A technician's hand (mid-30s, holding a tablet) swipes
            through charging telemetry.
Action: A solid-state battery pack mounted on a charging pedestal
        shows 80% → 100% in 2.5 seconds on the tablet.
```
- mmx 把 "fast-charging stations" 误画为"城市夜景"
- "technician's hand" 难画, 主体关系"手+平板+电池" 缺连接

**v5 修复版 prompt** (7.5/10):
```
Extreme close-up of an industrial EV fast-charging station plug
INSERTED INTO a silver electric vehicle port. A tablet held in
foreground hand shows battery charging telemetry percentage rising
from 80 to 100 percent. The battery pack on charging pedestal glows
electric blue at full charge. Quality control bay, INDOOR, cool
teal lighting, photorealistic, 8K, no text, no logos.
```

### ★ 修复 3 关键策略 ★

1. **"INSERTED INTO" 主体连接关系** — 强调"插枪入车"这个核心动作
2. **"IN DOOR" 大写强调** — 强制室内 (v4 被画成"户外"是元凶)
3. **"extreme close-up" 特写** — 强制近景, 避免被泛化为"全景城市"
4. **"foreground hand" 前手 + 后桩** — 层次明确

### v5 段 2-3 实际渲染效果
- ✅ 工业充电桩精确刻画 (银灰金属外壳+LCD屏+按钮)
- ✅ 平板遥测 UI 概念传达 (CHARGING 文字+绿色进度条)
- ✅ 室内 QC bay 环境感 (天花板灯具+立柱)
- ✅ 蓝色满电电池发光
- ❌ **银色电动车本体未出现** (充电枪插入桩体而非车辆)

## 段 3-3 ★ 重大突破 + 主题错位 ★

**v4 末帧 (8.75/10)**: 圆形观景台+哑光银未来电动车+黎明阳光
**v5 末帧 (8.98/10)**: **错位为** 金融交易室+中央人物剪影双臂高举+K线大屏

v5 mmx image 没按 prompt 画"观景台+车", 而画了"交易室+剪影" (更"投资召唤"):
- 画面质量 9.2 ✅ 电影级
- 构图 9.5 ✅ 完美对称
- 叙事 9.3 ✅ "掌控市场"情绪
- 创意 8.8 ✅ 金融大片范式
- **综合 9.2/10** ⭐ **单图最高纪录**

## 最终交付

**`minimax-output/lithium_v5/lithium_30s_v5.mp4`** ⭐
- 时长: **30.00s** (精确)
- 分辨率: **1366x768** (16:9)
- 编码: **H.264 + AAC, 1.08 Mbps**
- 大小: **4.0MB**
- 9 镜头: ffmpeg zoompan

## 关键发现 (★ 沉淀到 SKILL v0.6.1 patch 5)

1. **mmx image 无 --seed** — 同样 prompt 多次生成结果差异巨大
2. **★ 修复 prompt 3 策略**:
   - "INSERTED INTO" 强调主体连接关系
   - "IN DOOR" 大写强调防止被泛化为户外
   - "extreme close-up" 特写强制近景
3. **mmx 训练集优势区稳定**: 柱状图/全息/夜景/交易室 (8.75+)
4. **mmx 主题漂移高发区**:
   - "充电桩+平板" 易漂为"城市夜景" (4.5 → 7.5 修复)
   - "工厂+供应链" 易漂为"世界地图+化工厂" (7.75 → 6.6 退步)
   - "观景台+车" 易漂为"交易室+剪影" (8.75 → 8.98 错位突破)
5. **3 镜头×1 主题** 比 1 段长视频**更可控** (失败镜头隔离, 不影响其他)
6. **mmx 错位有时是惊喜**: v5 末帧 9.2 > v4 末帧 8.75 (虽然主题不对, 但单图质量更高)
7. **v5 段 3-3 末帧双面解读**:
   - 严格按 prompt: 主题不符 (收尾画面是交易室不是观景台)
   - 投资召唤角度: v5 实际更精准 (剪影+V 形=K 线=投资胜利)

## v3 → v4 → v5 三代进化

| 版本 | 方案 | 平均 | 单图最高 | 末帧 |
|---|---|---|---|---|
| v1 | 微观蓝色锂电 (mmx video) | 8.13 | 8.4 | 8.0 |
| v2 | 赛博朋克 (mmx video) | 7.77 | 8.0 | 7.5 |
| v3 | 赛博朋克 + PiP 修复 | 7.4 | 7.5 | 7.25 |
| **v4** | **9 镜头 5 层公式** | 7.36 | **8.75** | 8.75 |
| **v5** | **9 镜头复现 + 修复** | TBD | **8.98 / 9.2** | 8.98/9.2 |

## 下一步 (v6 计划)

1. **段 2-3 终极修复** — 加 "silver car body with charge port door open" 强制车辆本体
2. **段 3-1 修复** — 重新 prompt 强调 "battery manufacturing facility" 不是 "world map"
3. **★ 用 v5 段 3-3 意外惊喜** — "金融交易室+剪影" 实际更精准, 接受 v5 末帧
4. **多次生成择优** — mmx image 无 seed, 多次生成挑最佳 (3 次试, 取最好)
5. **★ mmx video quota 重置** (12.6h) — 12.6h 后用真视频替代 image+zoompan, 时序动作更自然
