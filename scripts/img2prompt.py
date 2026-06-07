#!/usr/bin/env python3
"""
via54Design — Image-to-Prompt Analyzer
Extract visual features from an image and generate structured prompts.
"""
import json, sys, os
from PIL import Image, ImageStat, ImageFilter
from collections import Counter

def analyze_image(path: str) -> dict:
    """Extract visual features from image."""
    if not os.path.exists(path):
        return {"error": f"file not found: {path}"}
    
    img = Image.open(path).convert("RGB")
    w, h = img.size
    
    # Basic info
    aspect = w / h
    aspect_label = "landscape" if aspect > 1.1 else "portrait" if aspect < 0.9 else "square"
    
    # Dominant colors (simple quantization)
    small = img.resize((64, 64))
    pixels = list(small.getdata())
    # Quantize to 16 colors
    quantized = [(r//32*32, g//32*32, b//32*32) for r,g,b in pixels]
    color_counts = Counter(quantized)
    top_colors = [{"rgb": list(c), "hex": f"#{c[0]:02x}{c[1]:02x}{c[2]:02x}", "count": n}
                  for c, n in color_counts.most_common(8) if n > 5]
    
    # Brightness analysis
    gray = img.convert("L")
    stat = ImageStat.Stat(gray)
    avg_brightness = stat.mean[0]
    brightness_label = "high-key" if avg_brightness > 200 else "low-key" if avg_brightness < 60 else "mid-key"
    
    # Contrast (std dev)
    std = stat.stddev[0] if stat.stddev else 0
    contrast = "high contrast" if std > 60 else "soft" if std < 25 else "moderate"
    
    # Colorfulness (variance across channels)
    r_stat = ImageStat.Stat(img.split()[0])
    g_stat = ImageStat.Stat(img.split()[1])
    b_stat = ImageStat.Stat(img.split()[2])
    color_var = (r_stat.stddev[0] or 0) + (g_stat.stddev[0] or 0) + (b_stat.stddev[0] or 0)
    colorful = "vibrant" if color_var > 120 else "muted" if color_var < 50 else "balanced"
    
    # Edge detection for details
    edges = img.filter(ImageFilter.FIND_EDGES).convert("L")
    edge_stat = ImageStat.Stat(edges)
    edge_density = edge_stat.mean[0]
    detail_level = "high detail" if edge_density > 40 else "smooth" if edge_density < 10 else "moderate detail"
    
    return {
        "width": w, "height": h,
        "aspect_ratio": f"{aspect:.2f}",
        "orientation": aspect_label,
        "brightness": round(avg_brightness, 1),
        "brightness_label": brightness_label,
        "contrast": contrast,
        "colorfulness": colorful,
        "detail": detail_level,
        "dominant_colors": top_colors[:5],
        "suggested_moods": suggest_moods(avg_brightness, std, color_var),
        "suggested_styles": suggest_styles(aspect_label, brightness_label, colorful),
    }

def suggest_moods(brightness: float, contrast: float, colorfulness: float) -> list:
    moods = []
    if brightness > 180: moods.extend(["bright", "airy", "ethereal"])
    elif brightness < 80: moods.extend(["moody", "dramatic", "noir"])
    else: moods.extend(["balanced", "natural"])
    if contrast > 50: moods.append("dramatic")
    if colorfulness > 100: moods.append("vibrant")
    elif colorfulness < 60: moods.append("subdued")
    return moods[:5]

def suggest_styles(orientation: str, brightness: str, colorfulness: str) -> list:
    styles = []
    if orientation == "landscape": styles.append("cinematic wide")
    if brightness == "high-key": styles.append("soft lighting")
    if brightness == "low-key": styles.append("chiaroscuro")
    if colorfulness == "vibrant": styles.append("color grading")
    styles.append("professional photography")
    return styles[:4]

def build_prompt_from_analysis(analysis: dict, user_prompt: str = "") -> str:
    """Build a structured prompt from analysis + user input."""
    parts = []
    if user_prompt:
        parts.append(user_prompt)
    
    # Add visual context
    ctx = []
    if analysis.get("brightness_label"):
        ctx.append(f"{analysis['brightness_label']} lighting")
    if analysis.get("contrast"):
        ctx.append(analysis["contrast"])
    if analysis.get("colorfulness"):
        ctx.append(f"{analysis['colorfulness']} colors")
    if analysis.get("detail"):
        ctx.append(analysis["detail"])
    if ctx:
        parts.append(f"({', '.join(ctx)})")
    
    # Add dominant colors as mood palette
    colors = analysis.get("dominant_colors", [])
    if colors:
        hexes = [c["hex"] for c in colors[:4]]
        parts.append(f"color palette: {', '.join(hexes)}")
    
    # Style hints
    styles = analysis.get("suggested_styles", [])
    if styles:
        parts.append(f"--style {', '.join(styles)}")
    
    return ". ".join(parts) + "."

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "usage: img2prompt.py <image_path> [user_prompt]"}))
        sys.exit(1)
    
    path = sys.argv[1]
    user = sys.argv[2] if len(sys.argv) > 2 else ""
    
    analysis = analyze_image(path)
    if "error" in analysis:
        print(json.dumps(analysis))
        sys.exit(1)
    
    prompt = build_prompt_from_analysis(analysis, user)
    analysis["generated_prompt"] = prompt
    
    print(json.dumps(analysis, indent=2))
