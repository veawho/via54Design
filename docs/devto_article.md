---
title: "How I Built an Open Source Go CLI to Bridge AI Storytelling, Modern Layouts, and Google Vids"
published: true
description: "A deep dive into how via54Design compiles structured storyboards into dark-neon CSS layouts and dynamically injects speaker notes into PPTX (OOXML) for Google Vids integration."
tags: "go, opensource, webdev, ai"
---

As developers, we've all been there: you ask Claude or ChatGPT to write a video script, a story outline, or a presentation draft. It spits out a beautiful script. 

But then, the painful manual labor begins:
1. You copy-paste scene descriptions into layout tools.
2. You struggle with grid alignment, typography, and styling.
3. You copy-paste voiceover lines into the "speaker notes" or narrators.
4. When you convert it into a video draft (e.g. using tools like Google Vids), you have to map everything again.

I wanted a single, unified pipeline. I wanted to feed in a one-line seed (e.g., *"A rogue AI waking up in an antique museum"*), and have it compile all the way to **a responsive HTML storyboard, editable slide presentations, and ready-to-narrate Google Vids draft files**.

So I built **via54Design** in Go. Here is how it works under the hood.

---

## The Architecture: Seed to Screenplay to Layout

The pipeline is split into four decoupled components:

```
[ One-line Seed ] 
       │  (via54 narrate)
       ▼
[ Storyboard JSON / Fountain Screenplay ]
       │
       ├─► [ 26-Dimensional Image Prompts ] (via54 prompt)
       ├─► [ 10 responsive grid CSS slides ] (via54 generate)
       └─► [ OOXML PPTX Presentation with Speaker Notes ] (via54 export)
```

### 1. The Narrator: Aligning Seeds to Structures
The CLI parses the user's seed sentence and maps it onto structured screenplay models (like the *Hero's Journey*, *Three-Act Structure*, or *Cinematic Epic Trailer*). 
It generates a 12-shot camera/lighting timeline:
- **Beats**: Detailed timing, active setups, and speaker voiceovers.
- **Shot scale**: Orbit, Dolly zoom, Close-up, or Detail sequences.
- **Lighting/SFX**: Mood tags (e.g., *volumetric lighting, low-bass boom*) for GenAI tools.

### 2. The Layout Compiler: Raw JSON to Dark-Neon Glassmorphism
Instead of slow headless browsers, the layout compiler compiles layout files instantly. 
It supports 10 golden-ratio grid systems (Bento Grid, 3-pane dashboards, editorial magazines) and applies CSS themes (like `cinematic-neon`—obsidian dark panels, amber glows, and electric violet accents).

### 3. The OOXML Exporter: Injecting Speaker Notes for Google Vids
This was the trickiest part. Since **Google Vids** (Google Workspace's AI video tool) does not have a public developer API, the cleanest way to import external scripts is by uploading a Google Slides file and converting it to video.

When you import Slides, Google Vids maps:
- **Slide pages** ➜ **Video Scenes**
- **Speaker Notes** ➜ **AI Voiceover Narration Script**

To automate this, the Go-based PPTX exporter compiles slides and injects the script directly into the PowerPoint speaker notes. PPTX files are technically zipped XML folders. To make speaker notes appear, we must generate:
1. `ppt/notesSlides/notesSlide[N].xml` containing the XML text placeholder:
   ```xml
   <p:notes>
     <p:cSld>
       <p:spTree>
         <p:sp>
           <p:txBody>
             <a:p><a:r><a:t>SPEAKER_NOTES_CONTENT</a:t></a:r></a:p>
           </p:txBody>
         </p:sp>
       </p:spTree>
     </p:cSld>
   </p:notes>
   ```
2. Two-way relations files (`ppt/slides/_rels/slide[N].xml.rels` and `ppt/notesSlides/_rels/notesSlide[N].xml.rels`) to bind the slide XML to its notes page.
3. ContentType overrides in `[Content_Types].xml`.

By writing this from scratch in Go (zero Python/Node wrappers), the CLI generates widescreen PPTX presentations in milliseconds.

---

## How to Try It

The main CLI is open-source and has a master test harness checking stability across Windows, macOS, and Linux.

1. **Generate the storyboard**:
   ```bash
   via54 narrate --seed "A lone mechanic in a desert wasteland finding water" --model cinematic-epic --output script.json
   ```
2. **Export to PPTX**:
   ```bash
   via54 export pptx script.json --output project.pptx
   ```
3. **Open as Slides & Convert to Google Vids**:
   Upload `project.pptx` to Google Drive, open it in Google Slides, and select `File -> Convert to video`. Vids will automatically read the speaker notes as your video narration script!

---

## What's Next?
The project is under active development. The latest `v1.1.0-beta` contains the complete Google Vids integration flow. 

I'd love to hear your feedback! Check out the repository on GitHub:
👉 **[GitHub Repo: veawho/via54Design](https://github.com/veawho/via54Design)**
