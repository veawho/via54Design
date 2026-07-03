# 🚀 via54Design 推广与增长策略案 (Promotion & Launch Toolkit)

本案基于 Reddit、Hacker News (HN)、Twitter/X 的主流开发者社区运作规律，量身定制了一套 **“价值导向/故事叙述 (Value-First Storytelling)”** 的推广文案包，以最大化点赞与转化率（避免生硬的硬广被管理员封锁）。

---

## 1. Reddit 推广文案模板 (Value-First Storytelling)

**建议发布板块**：`r/sideproject` (首选，支持发布自己的项目), `r/golang` (适合分享纯 Go 实现的 PPTX OOXML 生成原理), `r/AItools`

### 📝 英文帖模板 (Reddit Post Template)
*   **Title**: Show r/sideproject: I built via54Design — An open-source Go CLI that turns a 1-line seed into structured cinematic storyboard JSON, custom HTML layout slides, and PPTX with speaker notes for Google Vids.
*   **Content**:
    Hi everyone,

    I was tired of copy-pasting Claude/GPT outputs into slide editors, struggling with alignment, and then manually rewriting voiceovers when converting them into videos. 

    So I built **via54Design** in Go. It’s an open-source CLI/engine that builds a complete aesthetic media pipeline.

    ### 🛠️ What it does:
    1. **Seed ➜ Storyboard JSON**: You write one sentence (e.g. *"A rogue AI waking up in an antique museum"*), it outputs a 12-shot camera/lighting timeline JSON and a Fountain screenplay.
    2. **High-Fidelity Layout Compile**: Compiles the JSON beats into beautiful responsive HTML slides (supports Bento Grid, Dashboard 3-pane, pricing layouts) with Linear/Vercel-like dark-neon aesthetics.
    3. **OOXML PPTX Exporter**: Export to editable 16:9 widescreen PPTX. **The cool part**: It injects your script voiceover directly into the slide's Speaker Notes. 
    4. **Google Vids & Audio Bridge**: You upload the PPTX to Google Drive, click "Convert to video" in Google Slides, and Google Vids automatically picks up your speaker notes as the AI voiceover narration script.

    ### 💻 Codebase Details:
    * Written in 100% Go (zero Node.js/python dependencies for the main engine).
    * Dual-lingual docs with beautiful high-fidelity screenshots.
    * 110 automated tests passing on Windows/macOS/Linux.

    I’d love to hear your thoughts on the slide note bridging approach! What features would you like to see next?
    
    👉 GitHub: https://github.com/veawho/via54Design

---

## 2. Hacker News "Show HN" 启动方案

**发布格式**：Show HN: via54Design - Open-source Go CLI mapping storyboards to aesthetic layouts & Vids
**发布文案**：
> Show HN: via54Design
>
> Hello HN,
>
> I wanted to share via54Design, a Go CLI utility that bridges the gap between text-based AI storytelling and final slide layout compiling.
>
> Instead of using web-based layout builders, this tool parses narrative arcs (Hero's Journey, Cinematic Epic, etc.) into structured storyboard JSON files. It then compiles these scenes into responsive CSS-grid HTML or generates standard Office Open XML (PPTX) packages.
>
> To support the new Google Vids workflow without proprietary API hooks, we implemented native PPTX Speaker Notes XML injection (`ppt/notesSlides/notesSlide%d.xml`). When you convert the slide deck into Google Vids, the script voiceovers are automatically parsed as narration.
>
> The engine is written entirely in Go with zero external presentation dependencies.
>
> I'd love to hear feedback on our templating design.
>
> GitHub: https://github.com/veawho/via54Design

---

## 3. Twitter / X 宣传推文链 (Thread Template)

*   **Tweet 1**: 
    Tired of manually copying LLM screenplays into presentation layouts? 
    I built via54Design, an open-source Go engine that turns one-line inspiration into 12-shot cinematic storyboard JSON, responsive HTML pages, and editable PPTX. 🧵👇
    [Attach docs/images/cinematic_storyboard.jpg]
*   **Tweet 2**:
    Why Go? The CLI compiles down to a standalone binary with zero dependencies. 
    It supports 10 golden-ratio layouts (Bento, Dashboard 3-pane, Magazine grid) styled with modern glassmorphism and ambient dark-neon colors. 🎨✨
    [Attach docs/images/design_showcase.jpg]
*   **Tweet 3**:
    🎥 Integration with Google Vids:
    The PPTX exporter automatically writes voiceover notes into slide Speaker Notes XML. Upload to Google Drive ➜ Convert to Video in Google Slides ➜ Google Vids extracts notes as narration script automatically! 🤯
*   **Tweet 4**:
    It’s open-source, dual-lingual (EN/ZH), and runs natively on Windows, macOS, and Linux. 
    Give it a spin and let me know what you think! 🌟
    GitHub: https://github.com/veawho/via54Design

---

## 4. GitHub Awesome List 提议提交 (PR Submission)

您可以向 `awesome-go` 发起 Pull Request：
*   **分类**：`Video` / `GUI` 或 `Advanced CLI`
*   **提交描述**：
    > `[via54Design](https://github.com/veawho/via54Design) - Command line tool and layout engine that compiles narrative storyboards into editable PPTX slide notes and CSS-grid responsive HTML.`
