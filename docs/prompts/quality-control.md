# 生图质量门禁 (Quality Control)

## 质量层级

| 层级 | 提示词标签 | 效果 |
|------|-----------|------|
| ⭐⭐⭐⭐⭐ | `masterpiece, best quality, highly detailed` | 最高质量 |
| ⭐⭐⭐⭐ | `high quality, detailed` | 高质量 |
| ⭐⭐⭐ | `detailed` | 中等 |
| ⭐⭐ | (不加质量词) | 基础 |
| ⭐ | (加负面词) | 快速生成 |

## 平台常用质量词

| 平台 | 正面质量词 | 负面过滤词 |
|------|-----------|-----------|
| Midjourney | `--style raw -s 250` | `--no blurry, low quality` |
| Stable Diffusion | `masterpiece, best quality` | `nsfw, lowres, bad anatomy` |
| DALL·E 3 | 自然语言描述 | 不需要负面词 |

## 负面词分层

| 层级 | 排除内容 |
|------|---------|
| L1 结构 | `deformed, bad anatomy, disfigured, poorly drawn` |
| L2 肢体 | `extra limbs, missing fingers, bad hands` |
| L3 瑕疵 | `blurry, low quality, noisy, grainy` |
| L4 水印 | `watermark, signature, text, logo` |
| L5 环境 | `oversaturated, underexposed, bad lighting` |

## 验证清单

- [ ] 主体描述是否精确?
- [ ] 风格与内容是否匹配?
- [ ] 光线是否强化了氛围?
- [ ] 负面词是否过滤了常见瑕疵?
- [ ] Token 权重是否合理?(1.0-1.5)
- [ ] 平台参数是否正确?(--ar, --v, etc.)
