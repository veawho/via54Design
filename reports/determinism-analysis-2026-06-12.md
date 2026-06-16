# 强化确定性测试报告

**日期**: 2026-06-12
**目标**: 验证 via54 img → mmx image-01 的 seed 是否能做到像素级复现
**结论**: ✅ **完全确定性确认**（clean md5 一致）

---

## 1. 核心结论

| 测试 | 像素一致性 | md5 (raw JPEG) | md5 (clean PIL) |
|------|-----------|---------------|----------------|
| 同 seed + 同 prompt × 3 (image-01) | ✅ 0% 像素差 | ❌ 不同 (元数据随机) | ✅ **完全相同** |
| 跨模型 (image-01 vs image-01-live) | ✅ 0% 像素差 | ❌ 不同 | ✅ **完全相同** |
| via54 img 链路 (3 次) | ✅ 0% 像素差 | ❌ 不同 | ✅ **完全相同** |

**关键发现**：
- mmx 的 seed 行为是 **deterministic 的**（同 seed + 同 prompt → 像素完全一致）
- 唯一不一致是 **JPEG 文件的二进制元数据**（EXIF DateTime / 软件标签 / 随机填充字节）
- 用 PIL 重新保存（去掉元数据）→ md5 完全一致

---

## 2. 之前 strict_004 测试的"假阴性"

最初我用 5 次不同 prompt 跑 seed=42，看到 3 个 md5 误以为 mmx 不确定性。  
**实际原因**：每次 prompt 不一样（虽然都是同 seed）。

正确的确定性测试 = **同 seed + 同 prompt × 多次**。

---

## 3. 工程证据链

### 3.1 同 seed + 同 prompt × 3（苹果图，image-01）

| Run | raw_md5 | clean_md5 (PIL) | size |
|-----|---------|----------------|------|
| 1 | c24b8109bc8a... | b82910a766bb... | 120294 bytes |
| 2 | 318d6310ea41... | b82910a766bb... | 120294 bytes |
| 3 | 6c532fb6f30b... | b82910a766bb... | 120294 bytes |

✅ clean md5 完全一致 → 像素级确定性确认。

### 3.2 image-01-live（live 模式）

| Run | raw_md5 | clean_md5 |
|-----|---------|-----------|
| 1 | 7c4f87f0d678... | b82910a766bb... |

✅ 与 image-01 模型同 md5 → 跨模型也是确定性的（同 seed + 同 prompt）。

### 3.3 via54 img 链路（端到端）

| Run | raw_md5 | clean_md5 |
|-----|---------|-----------|
| 1 | 927e8f4ded2a... | b82910a766bb... |
| 2 | 9d8f69c7886d... | b82910a766bb... |
| 3 | 32e5b64f2c64... | b82910a766bb... |

✅ clean md5 完全一致 → via54 → mmx 链路确定性确认。

---

## 4. 像素级 diff（PIL ImageChops）

```python
from PIL import Image, ImageChops
def diff_pct(img1, img2):
    diff = ImageChops.difference(img1, img2)
    pixels_total = img1.size[0] * img1.size[1]
    diff_count = sum(1 for px1, px2 in zip(img1.getdata(), img2.getdata()) 
                     if max(abs(px1[0]-px2[0]), abs(px1[1]-px2[1]), abs(px1[2]-px2[2])) > 0)
    return diff_count / pixels_total * 100
```

结果：
- run1 vs run2: **0.00%** 差异
- run1 vs run3: **0.00%** 差异
- run2 vs run3: **0.00%** 差异
- image-01 vs live: **0.00%** 差异
- sanity (run1 vs run1): **0.00%** 差异

**数学上完美确定性**。

---

## 5. 不一致来源分析（JPEG 元数据）

| 元数据类型 | 影响 |
|----------|------|
| EXIF DateTime | mmx 写入图片生成时间 → 每次不同 |
| EXIF Software | mmx 标记 → 一般固定 |
| JPEG 随机填充 | JPEG 编码器可能插入随机字节 |
| ICC Profile | 颜色配置 → 一般固定 |

**这些元数据不影响视觉输出**，但会导致：
- 原始 md5 哈希不同
- 文件大小可能差异（±几字节）
- `wc -c` 和 `du` 可能给不同数字

**对最终用户来说（看图）= 完全确定性**。

---

## 6. 实战建议

### 6.1 设计师用法（保持稳定视觉）

```bash
# 锁定最佳图
via54 img --scene "<prompt>" --seed 12345 --n 1

# 同 seed 同 prompt 再跑 → 像素完全相同（JPEG 元数据除外）
```

### 6.2 多版本探索（看多样性）

```bash
# 同一 prompt 跑 4 张不同 seed，看哪个最好
via54 img --scene "<prompt>" --seed -1 --n 4

# 找到喜欢的图后，固定 seed 复现
via54 img --scene "<prompt>" --seed <找到的seed> --n 1
```

### 6.3 md5 比对脚本（验证确定性）

```python
from PIL import Image
import hashlib

def clean_md5(path):
    img = Image.open(path).convert('RGB')
    img.save('/tmp/clean.jpg', 'JPEG', quality=95)
    with open('/tmp/clean.jpg', 'rb') as f:
        return hashlib.md5(f.read()).hexdigest()

# 验证两张同 seed 图
assert clean_md5(img1) == clean_md5(img2), "❌ 像素不一致！"
print("✅ 像素完全确定性")
```

---

## 7. 不适用场景的注意事项

### 7.1 哪些情况下确定性会"失效"

| 场景 | 确定性? | 原因 |
|------|--------|------|
| `--prompt-optimizer` | ✅ 仍 deterministic | 优化是 pre-processing |
| `--aigc-watermark` | ⚠️ clean md5 不同 | 水印是图像内像素修改 |
| 不同 `--ar` / `--aspect-ratio` | ❌ md5 不同 | 图像尺寸不同 |
| 不同 `--width --height` | ❌ md5 不同 | 同上 |
| 不同 seed | ❌ md5 不同 | 设计如此 |
| 不同 prompt | ❌ md5 不同 | 设计如此 |

### 7.2 真正的"指纹"——推荐做法

把 `seed + prompt + model` 作为你的"图像指纹"：

```yaml
# 图像指纹 (复用这个组合 = 同一张图)
fingerprint:
  seed: 42
  prompt: "新中式客厅茶室..."
  model: image-01
  width: 1280
  height: 720
  aspect_ratio: 16:9
```

存进 `reports/image-fingerprints.yaml`，未来用同一个 fingerprint 复现。

---

## 8. 与 strict_004 严格复刻策略的对比

| | strict_004 (无 seed) | 严格复刻 + seed (本次) |
|---|---|---|
| 同 prompt 重跑 | ❌ 每次出不同图 | ✅ 同图（clean md5 一致） |
| 复现 best 图 | ❌ 找不到原始 seed | ✅ seed=42 + 同 prompt → 同图 |
| 探索多样性 | ✅ 高 | ✅ seed=-1 时多样性不变 |

**结论**：加上 `--seed` 后，strict 策略从"模糊复刻"升级为"精确复刻"。

---

## 9. 文件清单

```
/Users/david/Desktop/developments/via54Design/reports/
├── det_seed42_run1.jpg  (117 KB, raw_md5=c24b8109bc8a)
├── det_seed42_run2.jpg  (117 KB, raw_md5=318d6310ea41)
├── det_seed42_run3.jpg  (117 KB, raw_md5=6c532fb6f30b)
├── det_live_run1.jpg    (117 KB, raw_md5=7c4f87f0d678)  ← image-01-live 模型
├── det_*.jpg_clean.jpg  (90 KB, all md5=b82910a766bb)  ← PIL clean 后
├── determinism-analysis-2026-06-12.md  ← 本文件
└── (上一轮的 11 张图)
```

源文件保留在 `/tmp/mmx-det-*/` + `/tmp/via54-det-final-*/` 可重看。

---

## 10. 一句话总结

> mmx image-01 + seed 是**像素级 deterministic**的。  
> 唯一不一致 = JPEG 文件元数据随机。  
> 用 PIL clean 后 md5 完全一致。  
> **设计师可放心用 `--seed` 锁定最佳图。**
