# 分镜模板 (Storyboard Template)

基于 Fountain 剧本格式生成 shot-by-shot 分镜表。

## 字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| Shot | 镜头编号 | 001 |
| Time | 时间码 | 0:00 - 0:03 |
| Duration | 时长 | 3s |
| Type | 景别 | Wide / Medium / Close-up / Detail / POV |
| Camera | 运镜 | Static / Pan / Tilt / Zoom / Dolly / Crane |
| Visual | 画面描述 | 产品在黑色背景上旋转 |
| Audio | 音效/对白 | logo-reveal sfx |
| Voiceover | 旁白 | "终于有一个产品真正理解你" |
| Mood | 情绪 | mysterious / warm / confident / exciting |

## 生成的分镜表示例

```
SHOT 001  0:00-0:03  3s  WIDE      Static    暗场中逐渐亮起产品轮廓    ambient pad     "你是否曾经..."                mysterious
SHOT 002  0:03-0:06  3s  MEDIUM    Slow zoom  人物在使用翻译机          keyboard sfx    "因为语言不通而错过..."        empathetic
SHOT 003  0:06-0:10  4s  CLOSE-UP  Dolly in   产品屏幕显示实时翻译      logo-reveal     "现在有了XX"                  aspirational
SHOT 004  0:10-0:16  6s  DETAIL    Static     翻译界面展示多语种切换    click sfx       "精准到每一个术语"             confident
SHOT 005  0:16-0:20  4s  MEDIUM    Pan         跨国团队会议场景          ambience        "跨国团队协作从未如此简单"     warm
SHOT 006  0:20-0:24  4s  WIDE      Crane up    数据增长可视化图表        whoosh          "效率提升300%"                excited
SHOT 007  0:24-0:28  4s  MEDIUM    Static      Logo+Slogan 出现         brand-stamp     "让世界听懂你"                inspiring
SHOT 008  0:28-0:30  2s  CLOSE-UP  Static      CTA 按钮 + 网址          click           "立即体验"                    urgent
```

## 从故事到分镜的转换规则

```
故事 beat  →  多个 shot
每个 shot  →  一种景别 + 一种运镜 + 一段旁白
情绪变化   →  景别和运镜的节奏变化
高潮点     →  close-up + dolly in + 强音效
```

## 参考

- Fountain screenplay format: https://fountain.io
- huobao-drama (⭐12623): 一句话生成完整短剧
