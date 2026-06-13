# via54Design SVG 规范 (参考: Claude Design visualize.md)
# 来源: system_prompts_leaks (⭐41,383)
# 版本: 1.0.0

## viewBox
- 固定: `viewBox="0 0 680 H"` (宽度680不可变)
- H = max_y + 40px buffer
- 安全区: x=40..640, y=40..(H-40)

## 字号
- 仅2种: 14px(标签, class=t), 12px(副标题, class=ts)
- 每个<text>必须带class(t/ts), 未分类文本继承错误字体
- 换行: 显式<tspan>, SVG <text>从不自动换行

## 颜色
- 浅色: 50填充 + 600描边 + 800标题/600副标题
- 深色: 800填充 + 200描边 + 100标题/200副标题
- 禁止: 在彩色背景上使用黑/灰文本
- 连接线: fill="none" 在所有path/polyline上

## 规范
- 箭头: <defs>中定义标准cheveron(使用context-stroke)
- 描边: 0.5px(细边显精致)
- 圆角: rx=4(默认), rx=8(最大)
- 无图标/插图在框内(纯文本)
- 无装饰编号, 无超大标题
- 动画: 包裹@media (prefers-reduced-motion: no-preference)
- 仅transform和opacity动画
- 循环<2s
