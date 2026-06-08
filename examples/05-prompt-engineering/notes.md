# Prompt Engineering Example

## Use Case
Cross-platform AI image generation. Best for:
- Marketing team prototyping
- Designer exploring AI tools
- Multi-platform A/B testing

## 17 Supported Platforms

### Image
midjourney, flux, dalle3, sd3, stable_diffusion, ideogram, recraft

### Video
seedance, kling, pika, jimeng, veo, sora

### Generic / Camera / Keyframe
video_generic, video_camera, video_keyframe

### Foundation
gemini

## Quality Assessment

```bash
# 5-dimension quality scoring
./via54.exe prompt assess \
  --image generated.png \
  --prompt prompt-mj.md
```

Outputs: Subject (0-10), Composition, Style, Technical, Cultural.
