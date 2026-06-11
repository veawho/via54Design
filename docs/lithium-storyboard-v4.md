# 锂电 30s 视频分镜脚本 v4 (★ 5 层公式 + 3 镜头分镜 ★)

> 基于 lanshu-awesome-ai-video-kit 方法论 (03-分镜时序 + 02-进阶公式 + 13-Hailuo 02)
> Hailuo 02 = "Director's AI" — 喜欢叙事流+时序+物理细节+克制
> 公式: Scene → Characters → Action → Camera → Audio & Style
> **3 镜头分镜 + 1 段 prompt 1 段 10s 视频** (mmx Hailuo-2.3 一段 = 5.875s,
> 后期 ffmpeg setpts 拉长 1.7x → 10s)

---

## 段 1: 赛博朋克科技新闻 (0-10s)

### 镜头 1 (0-3.5s): 演播室开场 — 全息数据可视化

```
Scene: A dark cyberpunk newsroom with floor-to-ceiling windows showing
       a neon-lit megacity at night. Holographic data interfaces float
       in mid-air at chest height.
Characters: An unseen news anchor's hands (mid-30s, slender, no person
            visible) gesture toward a glowing lithium battery cell hologram.
Action: A 3D lithium battery cell rotates slowly between the host's
        open palms. Electric blue current flows along the cell's
        cathode-anode layers. As the host's right hand pushes forward,
        the hologram explodes into a constellation of energy density
        numbers and stock ticker data.
Camera: Slow dolly-in from medium shot to close-up of the hands, rack
        focus shifts from cell to floating data stream. End with a 0.5s
        hold on the rotating cell.
Audio & Style: Low synthetic hum rising in pitch. Soft data-tick sound
        effect on each new data point. Cold teal-magenta color grade,
        rim light on hands. Anamorphic lens flare on hologram edges.
```

### 镜头 2 (3.5-7s): 数据图滚动 — 投资趋势柱状图

```
Scene: A massive curved LED wall (40m wide) inside the same newsroom
       shows a real-time lithium battery production data dashboard.
Characters: No visible characters, only the host's voice (deep male,
            professional tone) implied off-screen.
Action: A 3D bar chart materializes left-to-right, each bar representing
        a major battery manufacturer's 2026 production capacity. As the
        bars stack, the topmost value (CATL at 800GWh) glows electric
        blue and pulses twice. A line graph then arcs upward showing
        global market growth from 2020 to 2030.
Camera: Lateral tracking shot left-to-right following the bar chart
        growth. Hold on the glowing peak bar for 0.3s.
Audio & Style: Crisp tick sound per bar appearance. A rising synth pad
        builds tension. Same teal-magenta grade continues.
```

### 镜头 3 (7-10s): 城市切换 — 工业未来感过渡

```
Scene: Exterior night view, slow push-in through a futuristic Asian
       megacity's industrial district (think Shenzhen-Shanghai skyline).
Characters: None — purely establishing shot.
Action: Aerial view descends slowly. Rows of glowing battery factories
        appear below, their blue-lit rooftops forming a geometric grid.
        Two delivery drones cross the foreground mid-frame, leaving
        short light trails. The city lights pulse in slow rhythm.
Camera: Drone descending forward dolly. Subtle parallax as foreground
        drones pass over the city. Final frame holds on the factory
        rooftop grid.
Audio & Style: Deep ambient city drone. Distant whoosh of drones.
        A single low piano note rings at the end, foreshadowing the
        next scene. Cool blue with warm orange window accents.
```

---

## 段 2: 固态电池工厂 (10-20s)

### 镜头 1 (0-3.5s): 工厂全貌 — 蓝白洁净车间

```
Scene: Interior of a sterile white-and-blue solid-state battery assembly
       line. Bright clean-room lighting, robotic arms in precise rows.
Characters: Two engineers in white cleanroom suits (genderless, focus
            on suits and hands) inspecting a battery cell.
Action: A robotic arm picks up a transparent solid-state battery cell,
        rotates it 180° for inspection, and places it on a conveyor.
        The conveyor moves 2 cells per second toward a packaging
        station. The lead engineer's gloved hand reaches in to adjust
        a misaligned cell.
Camera: Wide shot, slow push-in to medium. End on the engineer's
        gloved hand detail.
Audio & Style: Soft mechanical servo whir. Conveyor hum. No dialogue.
        Cool teal accent lights on the assembly line.
```

### 镜头 2 (3.5-7s): 微观制造 — 电池层压

```
Scene: Macro close-up of a battery cell's internal layers being
       pressed together in a hydraulic press.
Characters: None.
Action: A solid-state electrolyte layer (semi-transparent blue gel)
        is sandwiched between two electrode foils. The press descends
        with controlled force, the gel slightly bulging at the edges
        as pressure equalizes. Tiny bubbles escape. Final clamp locks
        with a satisfying click.
Camera: Static macro shot with very slow rack focus from top to bottom.
Audio & Style: Hydraulic press hiss. Gel squelch. A precise metallic
        click on lock. Subtle industrial music bed with plucked strings.
```

### 镜头 3 (7-10s): 成品测试 — 充电桩数据

```
Scene: Quality control bay with a row of fast-charging stations.
Characters: A technician's hand (mid-30s, holding a tablet) swipes
            through charging telemetry.
Action: A solid-state battery pack mounted on a charging pedestal
        shows 80% → 100% in 2.5 seconds on the tablet. The technician
        swipes to the next screen showing 500Wh/kg density. The pack
        glows a soft electric blue at full charge.
Camera: Close-up of tablet, then slow pull-back to medium revealing
        the technician and the pedestal in the same frame.
Audio & Style: Charging station beep. Subtle synth confirmation tone.
        Voice-over resumes: "能量密度突破五百瓦时每公斤."
```

---

## 段 3: 投资趋势 — 全球市场 (20-30s)

### 镜头 1 (0-3.5s): 工厂鸟瞰 — 蓝橙夜景

```
Scene: Aerial night view of a mega battery manufacturing facility
       ("Future Battery City") at dusk-to-night. Building exteriors
       light up sequentially.
Characters: None.
Action: The facility's LED facade projects a 3D holographic global
        supply chain map. Energy flows as light streams from lithium
        mines (Australia, Chile) → refining (China) → gigafactories
        (US, EU, China) → EV assembly. Each node lights up as energy
        passes through, creating a pulsing network visualization.
Camera: Slow orbit from west to east, gaining altitude. End on a
        high-angle wide shot.
Audio & Style: Orchestral swell with pulsing synth bass. 110 BPM
        ambient electronic bed. Cinematic teal-orange grade.
```

### 镜头 2 (3.5-7s): 数据中心 — 投资仪表盘

```
Scene: A financial trading floor in a glass-walled high-rise,
       night-time cityscape beyond.
Characters: Multiple traders' hands (diverse) typing on backlit
            keyboards. A central trader (silhouette only) gestures
            toward a massive curved monitor showing battery stock
            performance.
Action: The central monitor shows real-time tickers: CATL +4.2%,
        BYD +6.1%, LG Energy +3.8%. A portfolio total counter
        rapidly climbs: $2.5T → $3.1T. The central trader points
        to a peak.
Camera: Rack focus from foreground typing hands to background monitor.
        Slight handheld drift for tension.
Audio & Style: Rapid keyboard clatter. Soft cash register sound on
        each percentage gain. Subtle drum heartbeat builds.
```

### 镜头 3 (7-10s): 收尾 — 投资窗口召唤

```
Scene: A circular high-rise observation deck, dawn breaking. A
       circular infinity pool overlooks a battery factory complex
       in the distance.
Characters: No visible people, but a single futuristic EV (matte
            silver, LED accent line) parked in the foreground.
Action: The EV's LED strip pulses on, headlights ignite. Camera
        holds. The sun rises slightly behind the factory, casting
        long golden rays across the deck. Final frame: the EV
        silhouette against a sun-drenched industrial skyline.
Camera: Static wide shot, the only motion is the sunrise and EV
        lights activating. A 1.5s hold on the final composition.
Audio & Style: A single sustained synth note rises. Wind ambient.
        A deep male voice-over (the same anchor from Segment 1)
        delivers: "投资窗口,正在打开." Final 0.5s silence.
```

---

## 关键设计原则 (★ 沉淀)

1. **每个 prompt = 1 段 5.875s 视频, 后期 setpts 拉长 1.7x → 10s**
2. **3 镜头分镜, 每个镜头 3-3.5s** (符合 mmx 短时序 sweet spot)
3. **5 层公式**: Scene → Characters → Action → Camera → Audio & Style
4. **Hailuo 02 物理细节**: 触手可及 (液压/电芯/发光/伺服) 而非空泛 (史诗/震撼)
5. **同色温串联**: 段 1 冷青-品红, 段 2 冷青+白, 段 3 暖橙+冷青, 整片"科幻+金融"调性
6. **每段末帧 hold 0.3-0.5s** 便于后期接转场
7. **旁白只在段 1 起 + 段 2 末 + 段 3 末** (3 处出声, 不抢画面)
