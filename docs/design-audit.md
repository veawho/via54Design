╔════════════════════════════════════════════════════════════════════╗
║  via54Design × system_prompts_leaks (⭐41,383) 优化清单        ║
╚════════════════════════════════════════════════════════════════════╝

审计来源: Anthropic Claude Design / Claude visualize / Google Gemini / OpenAI Canvas
验证项目: slidev(47k⭐) / reveal.js(71k⭐) / marp(11.9k⭐) / guizang(15.5k⭐) / garden-skills(7.5k⭐)

┌────────────────────────────────────────────────────────────────────┐
│ 1. 暗色模式 (Dark Mode) — Claude Design: "Mandatory"                │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 缺失                                                   │
│ 参考标准: Claude Design要求所有颜色必须在light/dark双模式下工作       │
│ 验证项目: reveal.js(71k⭐)内置5暗色主题, guizang(15.5k⭐)双模式      │
│ 修复方案: 在engine.go的CSS模板中添加@media (prefers-color-scheme)    │
│ 优化结果: ✅ 已添加到checker_v2.RunAllV2()检测                       │
├────────────────────────────────────────────────────────────────────┤
│ 2. Viewport Meta — 响应式基础                                        │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 缺失                                                   │
│ 参考标准: Claude Design要求viewport meta用于移动端适配              │
│ 验证项目: 所有主流框架(slidev/reveal/marp)默认包含viewport meta     │
│ 修复方案: engine.go CSS生成时添加<meta name="viewport">            │
│ 优化结果: ✅ quality checker已检测(checkResponsive)                  │
├────────────────────────────────────────────────────────────────────┤
│ 3. @media 媒体查询 — 响应式布局                                      │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 缺失                                                   │
│ 参考标准: Claude Design grid/table overflow方案                     │
│ 验证项目: reveal.js(71k⭐)内置responsive支持                        │
│ 修复方案: layout模板添加响应式断点                                   │
│ 优化结果: ✅ quality checker已检测(checkResponsive)                  │
├────────────────────────────────────────────────────────────────────┤
│ 4. CSS变量系统 — 仅4/30配色含变量                                    │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ⚠️ 4/30                                                  │
│ 参考标准: Claude 9-Ramp × 7-Stop = 63色设计系统                     │
│ 验证项目: garden-skills(7.5k⭐) 23主题全token化                     │
│ 修复方案: 30配色逐步添加CSS变量                                      │
│ 优化结果: ✅ quality checker已检测(checkColorCompliance)            │
├────────────────────────────────────────────────────────────────────┤
│ 5. minmax(0, 1fr) — Grid溢出修复                                     │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 缺失(使用1fr而非minmax(0,1fr))                          │
│ 参考标准: Claude Design明确规定grid溢出修复方案                      │
│ 验证项目: OpenAI Canvas(Tailwind最佳实践)                           │
│ 修复方案: layout模板将1fr替换为minmax(0,1fr)                        │
│ 优化结果: ⏳ 待修复                                                │
├────────────────────────────────────────────────────────────────────┤
│ 6. 无障碍(Accessibility) — alt文本/焦点样式/reduced-motion          │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 缺失(0 alt标签, 0焦点样式)                              │
│ 参考标准: Claude Design: 所有img需alt, 所有元素需focus样式          │
│          动画必须包裹在prefers-reduced-motion中                      │
│ 验证项目: OpenAI Canvas使用shadcn/ui(内置无障碍)                     │
│ 修复方案:                                             │
│   - 布局模板添加alt属性占位                                          │
│   - CSS添加:focus-visible样式                                        │
│   - 动画添加@media (prefers-reduced-motion: no-preference)包裹     │
│ 优化结果: ✅ quality checker已检测(checkAccessibility)               │
├────────────────────────────────────────────────────────────────────┤
│ 7. font-weight 约束 — Claude Design: 仅400和500                     │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ⚠️ 未检测                                                │
│ 参考标准: Claude Design明确禁止600/700, 仅用400(regular)和500(bold)  │
│ 验证项目: Anthropic品牌规范                                          │
│ 修复方案: quality checker添加font-weight检测                        │
│ 优化结果: ✅ checker已检测(当前0个违规)                             │
├────────────────────────────────────────────────────────────────────┤
│ 8. Antipattern: !important                                           │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ✅ 已修复(engine.go移除!important)                         │
│ 参考标准: Claude Design禁止!important                                │
│ 优化结果: ✅ quality checker v2已检测                                │
├────────────────────────────────────────────────────────────────────┤
│ 9. Antipattern: scrollIntoView                                      │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ✅ 不存在                                                │
│ 参考标准: Claude Design明确禁止scrollIntoView                        │
│ 优化结果: ✅ checker已检测(0个违规)                                 │
├────────────────────────────────────────────────────────────────────┤
│ 10. Antipattern: position: fixed                                    │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ✅ 不存在                                                │
│ 参考标准: Claude Design禁止在iframe widget中使用position:fixed       │
│ 优化结果: ✅ checker已检测(0个违规)                                 │
├────────────────────────────────────────────────────────────────────┤
│ 11. 色阶系统 — 缺少Claude 9-Ramp参考                               │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 0/9色阶                                               │
│ 参考标准: Claude Design: 9色阶×7停位=63色系统, 每色阶明确语义        │
│ 验证项目: garden-skills(7.5k⭐) 23主题全token化                     │
│ 修复方案: 建立9-Ramp参考YAML(参考Claude Prompt中的精确色值)         │
│ 优化结果: ⏳ 待修复                                                │
├────────────────────────────────────────────────────────────────────┤
│ 12. SVG规范 — Claude visualize标准                                  │
├────────────────────────────────────────────────────────────────────┤
│ 当前状态: ❌ 未实现                                                │
│ 参考标准: Claude Design: SVG viewBox=680, text-anchor, fill=none    │
│          仅2种字号(14px/12px), 每次须带class(t/ts/th)              │
│ 验证项目: banan-slides(14.8k⭐) 内置SVG图表                         │
│ 修复方案: 创建SVG生成器模块(参考Claude visualize prompt)            │
│ 优化结果: ⏳ 待修复                                                │
└────────────────────────────────────────────────────────────────────┘

优化总结:
  ✅ 已修复: !important, quality checker v2 (8项检测)
  ⏳ 待优化: minmax(0,1fr), 9-Ramp色阶库, SVG规范, viewport
  未发现违规: scrollIntoView, position:fixed, font-weight 600/700

参考项目验证:
  system_prompts_leaks ⭐41,383 — 设计系统最全面
  reveal.js ⭐71,632 — HTML演示框架标配
  slidev ⭐47,028 — 主题包生态标杆
  guizang-ppt-skill ⭐15,557 — AGPL-3.0 设计系统最接近
  garden-skills ⭐7,516 — 主题token架构
  banan-slides ⭐14,857 — AI PPT生成
