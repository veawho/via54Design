#!/usr/bin/env python3
"""
via54Design — Storyboard → Video Pipeline
1. Single image → video opening/scene prompt
2. Full storyboard → narrative scaffold → video generation prompts
"""
import json, sys, os, re
from PIL import Image
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from img2prompt import analyze_image, build_prompt_from_analysis

NARRATIVE_MODELS = {
    "three-act": {
        "name": "三幕结构",
        "beats": [
            ("setup", "铺垫", 0.20, "平静"),
            ("inciting", "激励事件", 0.10, "好奇"),
            ("rising", "上升行动", 0.35, "紧张"),
            ("climax", "高潮", 0.15, "激动"),
            ("resolution", "结局", 0.20, "释然"),
        ]
    },
    "heros-journey": {
        "name": "英雄之旅",
        "beats": [
            ("ordinary", "平凡世界", 0.10, "平静"),
            ("call", "冒险召唤", 0.10, "好奇"),
            ("threshold", "跨越门槛", 0.15, "紧张"),
            ("trials", "考验/盟友/敌人", 0.30, "紧张"),
            ("ordeal", "最大考验", 0.10, "激动"),
            ("reward", "奖赏", 0.10, "喜悦"),
            ("return", "携宝而归", 0.15, "释然"),
        ]
    },
    "problem-solution": {
        "name": "问题→解决方案",
        "beats": [
            ("problem", "问题呈现", 0.20, "压抑"),
            ("pain", "痛点放大", 0.20, "紧张"),
            ("discovery", "发现方案", 0.15, "好奇"),
            ("solution", "方案展示", 0.30, "自信"),
            ("result", "成果愿景", 0.15, "希望"),
        ]
    },
}

def analyze_storyboard(image_paths: list, user_desc: str = "") -> dict:
    """Analyze multiple storyboard images and build a narrative."""
    frames = []
    for i, path in enumerate(image_paths):
        if not os.path.exists(path):
            frames.append({"index": i, "error": f"file not found: {path}"})
            continue
        analysis = analyze_image(path)
        prompt = build_prompt_from_analysis(analysis, user_desc)
        analysis["generated_prompt"] = prompt
        analysis["index"] = i
        frames.append(analysis)
    
    return {
        "total_frames": len(image_paths),
        "frames": frames,
        "narrative_scaffold": None,
        "video_prompts": None,
    }

def build_narrative_scaffold(image_analyses: list, model_id: str = "three-act",
                              duration: int = 30, user_desc: str = "") -> dict:
    """
    Build a narrative scaffold from analyzed storyboard images.
    Uses the beat structure defined in NARRATIVE_MODELS.
    """
    model = NARRATIVE_MODELS.get(model_id, NARRATIVE_MODELS["three-act"])
    total_beats = len(model["beats"])
    
    # Distribute images across beats
    beats = []
    remaining_images = list(image_analyses)
    images_per_beat = max(1, len(remaining_images) // total_beats)
    
    cur_time = 0
    for i, (bid, bname, weight, mood) in enumerate(model["beats"]):
        beat_dur = max(3, int(duration * weight))
        if i == total_beats - 1:
            beat_dur = max(3, duration - cur_time)
        
        # Assign images to this beat
        n_images = min(images_per_beat, len(remaining_images))
        if i == total_beats - 1:
            n_images = len(remaining_images)  # last beat gets all remaining
        beat_images = remaining_images[:n_images] if n_images > 0 else []
        remaining_images = remaining_images[n_images:]
        
        # Build visual context from images
        visual_context = ""
        if beat_images:
            colors = set()
            styles = set()
            for img in beat_images:
                if img.get("dominant_colors"):
                    for c in img["dominant_colors"][:2]:
                        colors.add(c.get("hex", ""))
                if img.get("suggested_styles"):
                    styles.update(img.get("suggested_styles", []))
            if colors:
                visual_context = f"palette: {', '.join(list(colors)[:4])}"
            if styles:
                visual_context += f" | style: {', '.join(list(styles)[:3])}"
        
        beat_prompt = ""
        if beat_images and beat_images[0].get("generated_prompt"):
            beat_prompt = beat_images[0]["generated_prompt"][:120]
        
        beats.append({
            "id": bid, "name": bname, "mood": mood,
            "start_time": cur_time, "duration": beat_dur,
            "voiceover": f"({bname}场景) {user_desc[:80] if user_desc else ''}",
            "visual_context": visual_context,
            "image_hint": beat_prompt,
            "image_count": len(beat_images),
            "translation": f"{bname} — {mood}",
        })
        cur_time += beat_dur
    
    # Build video prompts
    video_prompts = build_video_prompts(beats, user_desc)
    
    return {
        "model": model_id,
        "model_name": model["name"],
        "total_duration": duration,
        "beats": beats,
        "video_prompts": video_prompts,
        "narrative_text": "\n\n".join([
            f"## {b['name']} ({b['start_time']}s-{b['start_time']+b['duration']}s)\n"
            f"情绪: {b['mood']} | {b['visual_context']}\n"
            f"旁白: {b['voiceover']}"
            for b in beats
        ])
    }

def build_video_prompts(beats: list, user_desc: str) -> list:
    """Build ComfyUI-compatible video prompts from beats."""
    prompts = []
    for i, beat in enumerate(beats):
        is_first = (i == 0)
        prefix = "Opening: " if is_first else f"Scene {i+1}: "
        
        prompt_text = (
            f"{prefix}{beat['name']} — {user_desc[:60] if user_desc else beat['voiceover'][:60]}"
        )
        if beat.get("image_hint"):
            prompt_text += f". Visual reference: {beat['image_hint'][:80]}"
        
        prompts.append({
            "scene": i + 1,
            "beat_id": beat["id"],
            "prompt": prompt_text,
            "negative": "blurry, low quality, distorted, ugly",
            "mood": beat["mood"],
            "duration": beat["duration"],
            "workflow": "animatediff_txt2vid" if beat["duration"] > 5 else "sdxl_txt2img",
        })
    return prompts

def process_storyboard(image_dir_or_paths: list, model: str = "three-act",
                       duration: int = 30, desc: str = "", single_image_mode: bool = False) -> dict:
    """Main entry point: analyze images, build narrative, output structured result."""
    # Validate paths
    valid_paths = []
    for p in image_dir_or_paths:
        if os.path.isfile(p):
            valid_paths.append(p)
    
    if not valid_paths:
        return {"error": "no valid image files found", "paths_checked": image_dir_or_paths}
    
    # Analyze all images
    result = analyze_storyboard(valid_paths, desc)
    
    if single_image_mode or len(valid_paths) == 1:
        # Single image → generate opening scene prompt
        img = result["frames"][0]
        prompt = img.get("generated_prompt", desc)
        result["mode"] = "single_image"
        result["opening_prompt"] = prompt
        result["video_prompt"] = {
            "scene": 1,
            "prompt": f"Opening scene: {prompt[:200]}",
            "negative": "blurry, low quality",
            "workflow": "sdxl_txt2img",
            "duration": 5,
            "visual_context": {
                "colors": [c.get("hex") for c in (img.get("dominant_colors") or [])[:5]],
                "brightness": img.get("brightness_label"),
                "moods": img.get("suggested_moods", []),
            }
        }
    else:
        # Multiple images → full storyboard narrative
        result["mode"] = "storyboard"
        result["narrative_scaffold"] = build_narrative_scaffold(
            result["frames"], model, duration, desc
        )
        result["video_prompts"] = result["narrative_scaffold"]["video_prompts"]
    
    result["config"] = {
        "model": model,
        "duration": duration,
        "image_count": len(valid_paths),
    }
    return result

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="Storyboard → Video Pipeline")
    parser.add_argument("images", nargs="+", help="Image paths (1 for opening, multiple for storyboard)")
    parser.add_argument("--model", default="three-act", choices=list(NARRATIVE_MODELS.keys()))
    parser.add_argument("--duration", type=int, default=30, help="Target video duration in seconds")
    parser.add_argument("--desc", default="", help="User description / story seed")
    parser.add_argument("--single", action="store_true", help="Single image mode")
    args = parser.parse_args()
    
    result = process_storyboard(
        args.images, args.model, args.duration, args.desc, args.single
    )
    print(json.dumps(result, indent=2, ensure_ascii=False))
