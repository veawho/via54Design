#!/usr/bin/env python3
"""
via54_pipeline.py — Prompt orchestrator for the via54Design prompt engine.

Architecture:
    User Scene (Chinese/English)
      │
      ▼
    via54_pipeline.py
      │
      ├── P1: detect language → auto-translate Chinese→English if needed
      ├── P0: LLM call to fill 26 dimension fields with REAL content
      │       (replaces all placeholders)
      ├── P4: Reverse Image → Prompt (vision LLM analysis)
      ├── P5: Template Syntax — {variant|expansion}
      ├── P6: PromptArchive — JSONL-based prompt storage
      ├── P7: Browser MCP Integration (stub + instructions)
      └── P8: A1111 / ComfyUI Export Format
      │
      ▼
    via54 prompt --format json  ← feeds filled scaffold to Go binary
      │
      ▼
    Return: final prompt with ALL 26 dimensions populated

This module has zero external dependencies — it uses only Python stdlib
(urllib.request for HTTP, argparse for CLI, yaml optional via _load_yaml).
"""

from __future__ import annotations

import argparse
import base64
import json
import os
import random
import re
import subprocess
import sys
import tempfile
import urllib.error
import urllib.request
from dataclasses import dataclass, field, asdict
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

# The 26 dimension fields (updated v3.0.0 — from Flux-PG, Dynamic Prompts, ai-media-gen).
DIMENSION_FIELDS: List[str] = [
    # core subject
    "subject",          # 主体对象
    "secondary",        # 辅助元素
    # style
    "art_movement",     # 风格流派
    "artist_ref",       # 艺术家参考
    "medium",           # 媒介
    "genre",            # 类型/题材 (NEW v3)
    "hair",             # 发型/发色 (NEW v3)
    "pose",             # 姿态/动态 (NEW v3)
    # composition
    "camera_shot",      # 景别
    "composition_type",  # 构图
    "depth_of_field",   # 景深
    "view",             # 视角/视点 (NEW v3)
    "format",           # 画幅 (NEW v3)
    # lighting
    "lighting",         # 光线
    # color
    "color_palette",    # 色彩
    # environment
    "environment",      # 环境
    "weather",          # 天气 (NEW v3)
    "era",              # 时代
    "time",             # 时间 (NEW v3)
    # detail
    "texture",          # 纹理
    "effects",          # 效果
    "material",         # 材质 (NEW v3)
    "face",             # 面部 (NEW v3)
    "detail",           # 细节 (NEW v3)
    # quality
    "quality_tags",     # 质量标签
    "emotion",          # 情绪/氛围 (NEW v3)
]

DEFAULT_ENDPOINT = "https://api.openai.com/v1"
DEFAULT_MODEL = "gpt-4o-mini"

# Provider presets — endpoint + model + key_required for each provider.
# All use OpenAI-compatible /chat/completions format (Ollama, DeepSeek, Hermes
# all support this). Providers with key_required=False work with empty API key.
PROVIDER_PRESETS: Dict[str, Dict[str, Any]] = {
    "openai": {
        "endpoint": "https://api.openai.com/v1",
        "model": "gpt-4o-mini",
        "key_required": True,
        "description": "OpenAI GPT-4o / GPT-4o-mini (default)",
    },
    "deepseek": {
        "endpoint": "https://api.deepseek.com/v1",
        "model": "deepseek-chat",
        "key_required": True,
        "description": "DeepSeek Chat / DeepSeek V3",
    },
    "ollama": {
        "endpoint": "http://localhost:11434/v1",
        "model": "llama3.2",
        "key_required": False,
        "description": "Local Ollama (llama3.2, qwen2.5, etc.) — no API key needed",
    },
    "hermes": {
        "endpoint": "http://localhost:18791/v1",
        "model": "deepseek-v4-flash",
        "key_required": False,
        "description": "Hermes Agent gateway proxy (port 18791)",
    },
    "local": {
        "endpoint": "http://localhost:8000/v1",
        "model": "local-model",
        "key_required": False,
        "description": "Generic local OpenAI-compatible server (vLLM, llama.cpp, etc.)",
    },
}


# LLM system prompt for filling prompt dimensions.
SYSTEM_PROMPT_FILL = (
    "You are a professional AI image prompt engineer. "
    "Given a scene description, fill the 26 dimensions below with specific, vivid, "
    "English-language values. Each value should be 2-8 words, descriptive and concrete. "
    "Return ONLY a JSON object with the field values."
)

# LLM system prompt for Chinese→English translation.
SYSTEM_PROMPT_TRANSLATE = (
    "You are a professional translator. Translate the following Chinese scene description "
    "to English. Preserve all artistic and technical nuance. Return ONLY the English "
    "translation, no commentary."
)

# LLM system prompt for reverse image analysis.
SYSTEM_PROMPT_REVERSE = (
    "You are a professional AI image prompt engineer specializing in reverse prompt engineering. "
    "Analyze this image and infer the 26 prompt dimensions below with specific, vivid, "
    "English-language values. Each value should be 2-8 words, descriptive and concrete. "
    "Return ONLY a JSON object with the field values."
)

# ---------------------------------------------------------------------------
# Data models
# ---------------------------------------------------------------------------


@dataclass
class PromptScaffold:
    """Represents the full prompt scaffold with all 26 dimension fields."""

    scene: str = ""
    platform: str = "midjourney"
    fields: Dict[str, str] = field(default_factory=dict)
    negative: List[str] = field(default_factory=list)
    original_scene: str = ""  # preserved Chinese (if any)
    raw_prompt: str = ""  # combined final prompt string

    def to_dict(self) -> Dict[str, Any]:
        return asdict(self)

    @classmethod
    def from_json(cls, data: Dict[str, Any]) -> "PromptScaffold":
        return cls(
            scene=data.get("scene", ""),
            platform=data.get("platform", "midjourney"),
            fields=data.get("fields", {}),
            negative=data.get("negative", []),
            original_scene=data.get("original_scene", ""),
            raw_prompt=data.get("raw_prompt", ""),
        )


# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------


@dataclass
class Config:
    """Runtime configuration, resolved from env vars / config file."""

    llm_endpoint: str = DEFAULT_ENDPOINT
    llm_key: str = ""
    llm_model: str = DEFAULT_MODEL
    via54_binary: str = "via54.exe"

    @classmethod
    def from_env(cls, provider: str = "openai") -> "Config":
        """Load config from environment variables and optional YAML file.

        Supports provider presets via --provider flag. Provider presets set
        default endpoint + model + key_required. Env vars override presets.
        YAML file overrides env defaults. CLI flags override everything.

        Supported providers:
          openai, deepseek, ollama, hermes, local
        """
        cfg = cls()

        # Apply provider preset defaults
        preset = PROVIDER_PRESETS.get(provider, PROVIDER_PRESETS["openai"])
        cfg.llm_endpoint = preset["endpoint"]
        cfg.llm_model = preset["model"]

        # Try loading from ~/.via54/config.yaml
        config_path = Path.home() / ".via54" / "config.yaml"
        yaml_data = _load_yaml(config_path)
        if yaml_data:
            cfg.llm_endpoint = yaml_data.get("llm_endpoint", cfg.llm_endpoint)
            cfg.llm_key = yaml_data.get("llm_key", cfg.llm_key)
            cfg.llm_model = yaml_data.get("llm_model", cfg.llm_model)
            cfg.via54_binary = yaml_data.get("via54_binary", cfg.via54_binary)

        # Environment variables override config file
        cfg.llm_endpoint = os.environ.get("VIA54_LLM_ENDPOINT", cfg.llm_endpoint)
        cfg.llm_key = os.environ.get("VIA54_LLM_KEY", cfg.llm_key)
        cfg.llm_model = os.environ.get("VIA54_LLM_MODEL", cfg.llm_model)
        cfg.via54_binary = os.environ.get("VIA54_BINARY", cfg.via54_binary)

        # No fallback to VIA54_LLM_KEY — provider-agnostic means we don't
        # assume any specific provider. Local providers (Ollama/Hermes/local)
        # don't need an API key. Cloud providers will error if key is missing.
        return cfg


def _load_yaml(path: Path) -> Optional[Dict[str, Any]]:
    """Load a YAML file without the PyYAML dependency (basic parser).

    Falls back to JSON parsing if the file extension is .json or .yaml
    contains only JSON-safe content. Returns None if file doesn't exist
    or cannot be parsed.
    """
    if not path.exists():
        return None

    try:
        raw = path.read_text(encoding="utf-8").strip()
        # Try parsing as JSON first (valid YAML subset)
        if raw.startswith("{"):
            return json.loads(raw)

        # Simple YAML key: value parser (no nested structures)
        result = {}
        for line in raw.splitlines():
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            if ":" in line:
                key, _, val = line.partition(":")
                result[key.strip()] = val.strip().strip('"').strip("'")
        return result
    except (OSError, json.JSONDecodeError):
        return None


# ---------------------------------------------------------------------------
# P1 — i18n Auto-Translate
# ---------------------------------------------------------------------------


def _contains_chinese(text: str) -> bool:
    """Detect if text contains any CJK Unified Ideographs (Chinese characters).

    Args:
        text: The input string to check.

    Returns:
        True if any Chinese character is found.
    """
    for ch in text:
        if "\u4e00" <= ch <= "\u9fff":
            return True
    return False


def translate_to_english(
    text: str,
    api_key: str,
    endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
) -> Tuple[str, bool]:
    """P1: Auto-translate Chinese text to English using an LLM call.

    If the text contains no Chinese characters, returns it unchanged.

    Args:
        text: The scene description (may be Chinese or English).
        api_key: OpenAI-compatible API key.
        endpoint: OpenAI-compatible API endpoint.
        model: Model name to use for translation.

    Returns:
        A tuple of (translated_text, was_translated).
    """
    if not _contains_chinese(text):
        return text, False

    prompt = f"Translate this Chinese scene description to English:\n\n{text}"
    translated = _llm_chat_completion(
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT_TRANSLATE},
            {"role": "user", "content": prompt},
        ],
        api_key=api_key,
        endpoint=endpoint,
        model=model,
    )
    translated = translated.strip().strip('"').strip("'")
    return translated, True


# ---------------------------------------------------------------------------
# P0 — LLM Semantic Expansion
# ---------------------------------------------------------------------------


def expand_with_llm(
    scene: str,
    platform: str = "midjourney",
    api_key: str = "",
    endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
) -> Dict[str, Any]:
    """P0: Expand a scene description into all 26 dimension fields via LLM.

    Calls an OpenAI-compatible API to fill every dimension with specific,
    vivid, English-language values. Returns a dict with 'fields' (the 26
    dimension key-value pairs) and 'negative' (a list of negative prompt
    terms).

    Args:
        scene: The scene description (English).
        platform: Target platform (midjourney, dalle, stable-diffusion, etc.).
        api_key: OpenAI-compatible API key.
        endpoint: OpenAI-compatible API endpoint.
        model: Model name.

    Returns:
        Dict with keys: 'fields' (dict of 26 dimensions) and 'negative' (list).
    """
    fields_list = ", ".join(DIMENSION_FIELDS)
    user_prompt = (
        f"Platform: {platform}\n"
        f"Scene: {scene}\n\n"
        f"Fill these 26 dimensions with specific, vivid values:\n"
        f"{fields_list}\n\n"
        f"Also provide 3-5 negative prompt terms as an array 'negative'.\n"
        f'Return a JSON object with keys "fields" (object) and "negative" (array).'
    )

    raw = _llm_chat_completion(
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT_FILL},
            {"role": "user", "content": user_prompt},
        ],
        api_key=api_key,
        endpoint=endpoint,
        model=model,
    )

    return _parse_llm_expansion(raw)


def _parse_llm_expansion(raw: str) -> Dict[str, Any]:
    """Parse the LLM response into a structured expansion dict.

    Attempts to extract JSON from the LLM output (handling markdown fences).

    Args:
        raw: Raw LLM response text.

    Returns:
        Dict with 'fields' and 'negative' keys. Missing fields are filled
        with empty strings; missing negative defaults to empty list.
    """
    # Strip markdown code fences if present
    text = raw.strip()
    if text.startswith("```"):
        # Remove opening fence (possibly with language hint)
        first_nl = text.find("\n")
        if first_nl != -1:
            text = text[first_nl + 1 :]
        # Remove closing fence
        if text.endswith("```"):
            text = text[:-3].strip()
        elif "```" in text:
            text = text[: text.rindex("```")].strip()

    # Try to parse as JSON
    result: Dict[str, Any] = {"fields": {}, "negative": []}

    try:
        parsed = json.loads(text)
    except json.JSONDecodeError:
        # Last resort: try to find a JSON object in the text
        import re

        match = re.search(r"\{.*\}", text, re.DOTALL)
        if match:
            try:
                parsed = json.loads(match.group(0))
            except json.JSONDecodeError:
                # Store raw and return empty fields
                for field in DIMENSION_FIELDS:
                    result["fields"][field] = ""
                result["_parse_error"] = "Could not parse LLM response as JSON"
                result["_raw"] = raw
                return result
        else:
            for field in DIMENSION_FIELDS:
                result["fields"][field] = ""
            result["_parse_error"] = "No JSON found in LLM response"
            result["_raw"] = raw
            return result

    # Extract fields
    fields_data = parsed.get("fields", parsed)  # handle both nested and flat
    if isinstance(fields_data, dict):
        for field in DIMENSION_FIELDS:
            result["fields"][field] = str(fields_data.get(field, ""))
    else:
        for field in DIMENSION_FIELDS:
            result["fields"][field] = ""

    # Extract negative terms
    neg = parsed.get("negative", [])
    if isinstance(neg, list):
        result["negative"] = [str(x) for x in neg if x]
    else:
        result["negative"] = [str(neg)] if neg else []

    return result


# ---------------------------------------------------------------------------
# P4 — Reverse Image → Prompt (Vision Analysis)
# ---------------------------------------------------------------------------


def reverse_image(
    image_path: str,
    api_key: str,
    endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
) -> Dict[str, Any]:
    """P4: Analyze an image and return filled scaffold dict via vision LLM.

    Reads an image file, converts to base64, sends to a vision-capable LLM,
    and asks it to infer the 26 dimension fields from the image content.

    Args:
        image_path: Path to the image file (png, jpg, webp, etc.).
        api_key: OpenAI-compatible API key.
        endpoint: OpenAI-compatible API endpoint.
        model: Model name with vision support (e.g. gpt-4o, gpt-4o-mini).

    Returns:
        Dict with 'fields' (dict of 26 dimensions) and 'negative' (list),
        same format as expand_with_llm().
    """
    # Read and encode image
    with open(image_path, "rb") as f:
        image_data = f.read()
    image_b64 = base64.b64encode(image_data).decode("utf-8")

    # Guess mime type from extension
    ext = Path(image_path).suffix.lower()
    mime_map = {
        ".png": "image/png",
        ".jpg": "image/jpeg",
        ".jpeg": "image/jpeg",
        ".webp": "image/webp",
        ".gif": "image/gif",
        ".bmp": "image/bmp",
    }
    mime = mime_map.get(ext, "image/png")

    # Build the vision prompt — content is a list of parts for multimodal
    fields_list = ", ".join(DIMENSION_FIELDS)
    user_content: List[Dict[str, Any]] = [
        {
            "type": "text",
            "text": (
                f"Analyze this image and return JSON with 26 dimension fields:\n"
                f"{fields_list}\n\n"
                f"Also provide 3-5 negative prompt terms as an array 'negative'.\n"
                f'Return a JSON object with keys "fields" (object) and "negative" (array).'
            ),
        },
        {
            "type": "image_url",
            "image_url": {"url": f"data:{mime};base64,{image_b64}"},
        },
    ]

    raw = _llm_chat_completion_vision(
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT_REVERSE},
            {"role": "user", "content": user_content},
        ],
        api_key=api_key,
        endpoint=endpoint,
        model=model,
    )

    return _parse_llm_expansion(raw)


# ---------------------------------------------------------------------------
# P5 — Template Syntax ({variant|syntax} expansion)
# ---------------------------------------------------------------------------


def expand_variants(scene: str, count: int = 1) -> List[str]:
    """P5: Expand {opt1|opt2|opt3} template syntax into one or more variants.

    Parses the scene for curly-brace pipe-delimited option groups like
    {red|blue|green}. When count=1, randomly picks one option per group.
    When count>1, generates N unique combinations (up to the total number
    of possible combinations).

    Args:
        scene: Scene text potentially containing {option|syntax} patterns.
        count: Number of unique variants to generate (default: 1).

    Returns:
        List of expanded scene strings.
    """
    # Find all variant groups
    pattern = re.compile(r"\{([^}]+)\}")
    matches = list(pattern.finditer(scene))

    if not matches:
        # No variants found; return the original scene (repeated if count>1)
        return [scene] * count

    # Parse each group into a list of options
    option_lists: List[List[str]] = []
    for m in matches:
        options = [opt.strip() for opt in m.group(1).split("|") if opt.strip()]
        if options:
            option_lists.append(options)

    if not option_lists:
        return [scene] * count

    # Compute total possible combinations
    total = 1
    for opts in option_lists:
        total *= len(opts)

    # Generate unique combinations
    results: List[str] = []
    seen: set = set()

    # Helper: pick a specific combination by index (for deterministic generation)
    def _make_variant(indices: Tuple[int, ...]) -> str:
        result = scene
        # Replace matches from left to right, using the corresponding index
        for i, m in enumerate(reversed(list(matches))):
            # Build replacement
            opts = option_lists[i] if i < len(option_lists) else [m.group(1)]
            if i < len(indices):
                replacement = opts[indices[i] % len(opts)]
            else:
                replacement = random.choice(opts)
            # Work backwards to avoid offset issues
            result = result[: m.start()] + replacement + result[m.end():]
            # Remove the match from the list so we don't double-replace
            # Actually, scanning forward: let me do forward replacement
        return result

    # Forward replacement approach
    def _make_variant_forward(indices: Tuple[int, ...]) -> str:
        result = scene
        # We need to track offsets as we replace
        offset = 0
        for i, m in enumerate(pattern.finditer(result)):
            if i >= len(option_lists):
                break
            opts = option_lists[i]
            idx = indices[i] % len(opts) if i < len(indices) else random.randint(0, len(opts) - 1)
            replacement = opts[idx]
            # Adjust for previous replacements
            actual_start = m.start() + offset
            actual_end = m.end() + offset
            result = result[:actual_start] + replacement + result[actual_end:]
            offset += len(replacement) - (m.end() - m.start())
        return result

    if count == 1:
        indices = tuple(random.randint(0, len(opts) - 1) for opts in option_lists)
        return [_make_variant_forward(indices)]

    # Generate multiple unique variants
    attempts = 0
    max_attempts = max(count * 3, total * 3)
    while len(results) < count and attempts < max_attempts:
        indices = tuple(random.randint(0, len(opts) - 1) for opts in option_lists)
        if indices not in seen:
            seen.add(indices)
            results.append(_make_variant_forward(indices))
        attempts += 1

    return results


# ---------------------------------------------------------------------------
# P6 — Prompt Archive (JSONL-based)
# ---------------------------------------------------------------------------


class PromptArchive:
    """P6: JSONL-based prompt archive with save, search, list, and delete.

    Stores prompt records as one JSON object per line in
    ~/.via54/archive.jsonl. Zero external dependencies.

    Each record has: id, scene, platform, fields, final_prompt, created_at, tags.
    """

    def __init__(self, path: Optional[str] = None):
        """Initialize the archive.

        Args:
            path: Path to the JSONL archive file (default: ~/.via54/archive.jsonl).
        """
        if path is None:
            path = str(Path.home() / ".via54" / "archive.jsonl")
        self.path = Path(path)
        self.path.parent.mkdir(parents=True, exist_ok=True)

    def save(self, scaffold: PromptScaffold, tags: Optional[List[str]] = None) -> str:
        """Save a prompt scaffold to the archive.

        Args:
            scaffold: The PromptScaffold to archive.
            tags: Optional list of tags for search.

        Returns:
            The record ID string.
        """
        import uuid

        record_id = uuid.uuid4().hex[:8]
        record = {
            "id": record_id,
            "scene": scaffold.scene,
            "platform": scaffold.platform,
            "fields": scaffold.fields,
            "negative": scaffold.negative,
            "final_prompt": scaffold.raw_prompt,
            "created_at": datetime.now().isoformat(),
            "tags": tags or [],
        }
        with open(self.path, "a", encoding="utf-8") as f:
            f.write(json.dumps(record, ensure_ascii=False) + "\n")
        return record_id

    def search(self, query: str, limit: int = 10) -> List[Dict[str, Any]]:
        """Search archive for prompts matching query in scene or tags.

        Performs a simple case-insensitive substring match on scene and tags.

        Args:
            query: Search string.
            limit: Maximum number of results to return (default: 10).

        Returns:
            List of matching record dicts.
        """
        results: List[Dict[str, Any]] = []
        if not self.path.exists():
            return results
        query_lower = query.lower()
        with open(self.path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                try:
                    record = json.loads(line)
                except json.JSONDecodeError:
                    continue
                if query_lower in record.get("scene", "").lower() or any(
                    query_lower in t.lower() for t in record.get("tags", [])
                ):
                    results.append(record)
                    if len(results) >= limit:
                        break
        return results

    def list(self, recent: int = 20) -> List[Dict[str, Any]]:
        """List recent archive entries.

        Args:
            recent: Number of most recent records to return (default: 20).

        Returns:
            List of record dicts, newest first.
        """
        records: List[Dict[str, Any]] = []
        if not self.path.exists():
            return records
        with open(self.path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                try:
                    records.append(json.loads(line))
                except json.JSONDecodeError:
                    continue
        return records[-recent:][::-1]  # newest first

    def delete(self, record_id: str) -> bool:
        """Delete an archive entry by its id.

        Args:
            record_id: The id of the record to delete.

        Returns:
            True if the record was found and deleted, False otherwise.
        """
        if not self.path.exists():
            return False
        kept: List[Dict[str, Any]] = []
        found = False
        with open(self.path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                try:
                    record = json.loads(line)
                except json.JSONDecodeError:
                    continue
                if record.get("id") == record_id:
                    found = True
                    continue
                kept.append(record)
        if found:
            with open(self.path, "w", encoding="utf-8") as f:
                for record in kept:
                    f.write(json.dumps(record, ensure_ascii=False) + "\n")
        return found


# ---------------------------------------------------------------------------
# P7 — Browser MCP Integration (stub + instructions)
# ---------------------------------------------------------------------------


def browser_submit(
    scaffold: PromptScaffold,
    platform: str = "midjourney",
) -> str:
    """P7: Generate executable instructions for Browser Use / Playwright MCP.

    Does NOT perform actual browser automation — that requires the Browser Use
    tool or Playwright MCP server running alongside. Instead, outputs the
    precise command sequence a user or orchestrator would run.

    Args:
        scaffold: The PromptScaffold with final prompt to submit.
        platform: Target platform (midjourney, dalle, etc.).

    Returns:
        Formatted command string for browser automation.
    """
    prompt = scaffold.raw_prompt or _build_raw_prompt(scaffold)

    instructions = [
        f"# P7 Browser MCP — Submit to {platform}",
        "",
        f"## 1. Final prompt to submit:",
        prompt,
        "",
        "## 2. Browser Use MCP command sequence:",
        "",
    ]

    if platform == "midjourney":
        instructions.extend([
            "# Navigate to Midjourney Discord (or Midjourney web app)",
            "browser goto https://www.midjourney.com/imagine",
            "",
            "# Type prompt into the input field",
            f'browser type "#imagine {prompt}"',
            "",
            "# Or use the Discord channel:",
            'browser type "/imagine"',
            f'browser type "{prompt}"',
            'browser press Enter',
        ])
    elif platform == "dalle":
        instructions.extend([
            "browser goto https://www.bing.com/create",
            f'browser type "{prompt}"',
            "browser click 'Create' button",
        ])
    elif platform in ("stable-diffusion", "flux"):
        instructions.extend([
            "# For Automatic1111 WebUI:",
            "browser goto http://127.0.0.1:7860",
            f'browser type "#txt2img_prompt" "{prompt}"',
            "browser click '#txt2img_generate'",
            "",
            "# For ComfyUI: Paste via workflow node",
            f"# Load prompt into CLIPTextEncode node: {prompt}",
        ])
    else:
        instructions.extend([
            f"# Generic browser-based submission for {platform}",
            "# 1. Open the platform's web interface",
            "# 2. Paste the prompt into the text input",
            "# 3. Click the generate/submit button",
            "",
            f"Prompt: {prompt}",
        ])

    instructions.extend([
        "",
        "## 3. Playwright (Python) equivalent:",
        "# Requires: pip install playwright",
        "# from playwright.sync_api import sync_playwright",
        "",
        "## 4. Browser Use MCP tool call:",
        '# Submit via tool: browser_submit with content:',
        f'# {prompt[:200]}{"..." if len(prompt) > 200 else ""}',
    ])

    return "\n".join(instructions)


# ---------------------------------------------------------------------------
# P8 — A1111 / SD.Next / ComfyUI Export Format
# ---------------------------------------------------------------------------


def export_a1111(scaffold: PromptScaffold) -> str:
    """P8: Convert via54's structured prompt to Automatic1111 / SD.Next format.

    Produces a text block readable by A1111 with:
      - Positive prompt: comma-separated dimension values
      - Negative prompt: from scaffold.negative
      - Settings: Steps, CFG Scale, Sampler, Seed defaults

    Args:
        scaffold: The PromptScaffold with populated fields and negative terms.

    Returns:
        Multi-line string in A1111-compatible format.
    """
    # Build positive prompt from non-empty fields
    positive_parts = []
    for field in DIMENSION_FIELDS:
        val = scaffold.fields.get(field, "")
        if val:
            positive_parts.append(val)

    positive = ", ".join(positive_parts)
    negative = ", ".join(scaffold.negative) if scaffold.negative else ""

    lines = [
        "# Automatic1111 / SD.Next Prompt Export",
        "# Generated by via54Design pipeline",
        "",
        f"Positive prompt:",
        positive,
        "",
    ]

    if negative:
        lines += [
            f"Negative prompt:",
            negative,
            "",
        ]

    lines += [
        "Settings:",
        "    Steps: 20",
        "    CFG Scale: 7.0",
        "    Sampler: DPM++ 2M Karras",
        "    Seed: -1",
        "    Size: 1024x1024",
        "",
        "# To use in A1111:",
        "# 1. Copy Positive prompt into the txt2img prompt field",
        "# 2. Copy Negative prompt into the negative prompt field",
        "# 3. Apply the settings above",
        "# 4. Click Generate",
    ]

    return "\n".join(lines)


def export_comfyui_clip(scaffold: PromptScaffold, node_id: str = "6") -> str:
    """P8: Export prompt as A1111-format text for ComfyUI CLIPTextEncode nodes.

    Produces a text block suitable for pasting into a ComfyUI CLIPTextEncode
    node's text field. Includes the positive prompt. Negative prompt is
    described separately for a second CLIPTextEncode node.

    Args:
        scaffold: The PromptScaffold with populated fields and negative terms.
        node_id: The ComfyUI node ID (default "6" for positive CLIPTextEncode).

    Returns:
        Multi-line string in ComfyUI-compatible format.
    """
    # Build positive prompt from non-empty fields
    positive_parts = []
    for field in DIMENSION_FIELDS:
        val = scaffold.fields.get(field, "")
        if val:
            positive_parts.append(val)

    positive = ", ".join(positive_parts)
    negative = ", ".join(scaffold.negative) if scaffold.negative else ""

    pos_node = node_id
    neg_node = str(int(node_id) + 1) if node_id.isdigit() else "7"

    lines = [
        "# ComfyUI CLIPTextEncode Export",
        "# Generated by via54Design pipeline",
        "",
        f"## Node {pos_node} — CLIPTextEnode (positive):",
        positive,
        "",
    ]

    if negative:
        lines += [
            f"## Node {neg_node} — CLIPTextEncode (negative):",
            negative,
            "",
        ]

    lines += [
        "# To use in ComfyUI:",
        "# 1. Add two CLIPTextEncode nodes",
        f"# 2. Paste positive prompt into node {pos_node}",
        f"# 3. Paste negative prompt into node {neg_node}",
        "# 4. Connect node outputs to your model's CLIP input",
        "# 5. Connect KSampler accordingly",
    ]

    return "\n".join(lines)


# ---------------------------------------------------------------------------
# LLM Helpers
# ---------------------------------------------------------------------------


def _llm_chat_completion(
    messages: List[Dict[str, str]],
    api_key: str,
    endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
) -> str:
    """Call an OpenAI-compatible chat completion API (text-only).

    Uses stdlib urllib.request — no external dependencies.

    Args:
        messages: List of message dicts with 'role' and 'content' (strings).
        api_key: API key for authentication.
        endpoint: Base URL of the API (e.g. https://api.openai.com/v1).
        model: Model identifier.

    Returns:
        The text content of the assistant's response.

    Raises:
        RuntimeError: If the API call fails or returns an error status.
    """
    url = endpoint.rstrip("/") + "/chat/completions"

    payload = json.dumps(
        {
            "model": model,
            "messages": messages,
            "temperature": 0.7,
            "max_tokens": 2048,
        }
    ).encode("utf-8")

    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"

    req = urllib.request.Request(url, data=payload, headers=headers, method="POST")

    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            body = resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        error_body = e.read().decode("utf-8", errors="replace")
        raise RuntimeError(
            f"LLM API HTTP {e.code} error: {error_body}"
        ) from e
    except urllib.error.URLError as e:
        raise RuntimeError(f"LLM API connection failed: {e.reason}") from e

    data = json.loads(body)

    if "choices" not in data or not data["choices"]:
        raise RuntimeError(f"LLM API returned unexpected response: {body[:500]}")

    return data["choices"][0]["message"]["content"]


def _llm_chat_completion_vision(
    messages: List[Dict[str, Any]],
    api_key: str,
    endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
) -> str:
    """Call an OpenAI-compatible chat completion API (vision/multimodal).

    Same as _llm_chat_completion but accepts vision-style message content
    (list of content parts with image_url and text types).

    Args:
        messages: List of message dicts; 'content' may be a string or list of parts.
        api_key: API key for authentication.
        endpoint: Base URL of the API.
        model: Model identifier (must support vision, e.g. gpt-4o).

    Returns:
        The text content of the assistant's response.

    Raises:
        RuntimeError: If the API call fails or returns an error status.
    """
    url = endpoint.rstrip("/") + "/chat/completions"

    payload = json.dumps(
        {
            "model": model,
            "messages": messages,
            "temperature": 0.7,
            "max_tokens": 2048,
        }
    ).encode("utf-8")

    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"

    req = urllib.request.Request(url, data=payload, headers=headers, method="POST")

    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            body = resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        error_body = e.read().decode("utf-8", errors="replace")
        raise RuntimeError(
            f"LLM API HTTP {e.code} error: {error_body}"
        ) from e
    except urllib.error.URLError as e:
        raise RuntimeError(f"LLM API connection failed: {e.reason}") from e

    data = json.loads(body)

    if "choices" not in data or not data["choices"]:
        raise RuntimeError(f"LLM API returned unexpected response: {body[:500]}")

    return data["choices"][0]["message"]["content"]


# ---------------------------------------------------------------------------
# Main Pipeline
# ---------------------------------------------------------------------------


def pipeline(
    scene: str,
    platform: str = "midjourney",
    ref_image: Optional[str] = None,
    api_key: str = "",
    api_endpoint: str = DEFAULT_ENDPOINT,
    model: str = DEFAULT_MODEL,
    binary: str = "via54.exe",
    output_path: Optional[str] = None,
) -> PromptScaffold:
    """Run the full via54 prompt enhancement pipeline.

    Steps:
        1. Detect language → translate Chinese→English if needed.
        2. Fill all 26 dimension fields via LLM (P0).
        3. Create a temporary JSON scaffold file.
        4. Call ``via54 prompt --format json`` to process through Go binary.
        5. Read the result and enhance it with the LLM-filled fields.
        6. Save to output file if requested, return enhanced PromptScaffold.

    Args:
        scene: User scene description (Chinese or English).
        platform: Target image generation platform.
        ref_image: Optional reference image path.
        api_key: OpenAI-compatible API key.
        api_endpoint: OpenAI-compatible API endpoint URL.
        model: LLM model name.
        binary: Path to the via54 binary.
        output_path: Optional file path to save the result JSON.

    Returns:
        An enhanced PromptScaffold with all 26 dimensions populated.
    """
    scaffold = PromptScaffold()
    scaffold.scene = scene
    scaffold.platform = platform
    scaffold.original_scene = scene

    # Step 1: P1 — i18n Auto-Translate
    english_scene, was_translated = translate_to_english(
        scene, api_key=api_key, endpoint=api_endpoint, model=model
    )
    scaffold.scene = english_scene

    # Step 2: P0 — LLM Semantic Expansion
    expansion = expand_with_llm(
        english_scene,
        platform=platform,
        api_key=api_key,
        endpoint=api_endpoint,
        model=model,
    )
    scaffold.fields = expansion.get("fields", {})
    scaffold.negative = expansion.get("negative", [])

    # Step 3: Create a temporary JSON scaffold for the Go binary
    scaffold_data = scaffold.to_dict()
    tmp_dir = Path(tempfile.mkdtemp(prefix="via54_"))
    scaffold_path = tmp_dir / "scaffold.json"

    try:
        scaffold_path.write_text(
            json.dumps(scaffold_data, indent=2, ensure_ascii=False),
            encoding="utf-8",
        )

        # Step 4: Call via54 prompt binary
        filled_path = tmp_dir / "filled.json"
        cmd = [
            binary,
            "prompt",
            "--scene",
            english_scene,
            "--platform",
            platform,
            "--format",
            "json",
            "--output",
            str(filled_path),
        ]
        if ref_image:
            cmd.extend(["--ref-image", ref_image])

        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=120,
        )

        if result.returncode != 0:
            raise RuntimeError(
                f"via54 binary exited with code {result.returncode}:\n"
                f"stdout: {result.stdout[:500]}\n"
                f"stderr: {result.stderr[:500]}"
            )

        # Step 5: Read result and enhance with LLM-filled fields
        if filled_path.exists():
            go_output = json.loads(filled_path.read_text(encoding="utf-8"))
            # Merge: Go binary's raw_prompt takes precedence if present
            if "raw_prompt" in go_output:
                scaffold.raw_prompt = go_output["raw_prompt"]
            if "fields" in go_output and isinstance(go_output["fields"], dict):
                # Merge Go binary fields, but let LLM fields override
                go_fields = go_output["fields"]
                for key in DIMENSION_FIELDS:
                    if key not in scaffold.fields or not scaffold.fields.get(key):
                        scaffold.fields[key] = go_fields.get(key, "")
            if "negative" in go_output:
                # Combine negative lists (deduplicate)
                combined = list(scaffold.negative)
                for n in go_output["negative"]:
                    if n not in combined:
                        combined.append(n)
                scaffold.negative = combined
        else:
            # If Go binary didn't produce a file, use what we have
            scaffold.raw_prompt = _build_raw_prompt(scaffold)

        # Step 6: Save to output if requested
        if output_path:
            output_path_obj = Path(output_path)
            output_path_obj.parent.mkdir(parents=True, exist_ok=True)
            output_path_obj.write_text(
                json.dumps(scaffold.to_dict(), indent=2, ensure_ascii=False),
                encoding="utf-8",
            )

    finally:
        # Clean up temp directory
        try:
            for f in tmp_dir.iterdir():
                f.unlink()
            tmp_dir.rmdir()
        except OSError:
            pass

    return scaffold


def _build_raw_prompt(scaffold: PromptScaffold) -> str:
    """Build a combined prompt string from the scaffold fields.

    Constructs a Midjourney-style prompt from the 26 dimension fields.

    Args:
        scaffold: The PromptScaffold with populated fields.

    Returns:
        A combined prompt string.
    """
    parts = []
    for field in DIMENSION_FIELDS:
        val = scaffold.fields.get(field, "")
        if val:
            parts.append(val)

    prompt = ", ".join(parts)
    if scaffold.negative:
        neg_str = ", ".join(scaffold.negative)
        prompt += f" --negative {neg_str}"
    return prompt


# ---------------------------------------------------------------------------
# CLI Interface
# ---------------------------------------------------------------------------


def parse_args(argv: Optional[List[str]] = None) -> argparse.Namespace:
    """Parse command-line arguments for the pipeline CLI.

    Supports two modes:
      1. Main pipeline mode (default): --scene, --platform, etc.
      2. Archive subcommands: archive {save,list,search,delete} [...]

    Args:
        argv: Argument list (defaults to sys.argv[1:]).

    Returns:
        Parsed arguments namespace. The 'command' field distinguishes mode.
    """
    parser = argparse.ArgumentParser(
        description="via54Design Prompt Pipeline — LLM expansion + i18n + Go binary orchestration",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            "Examples:\n"
            "  python via54_pipeline.py --scene \"a serene mountain landscape\" --platform midjourney\n"
            "  python via54_pipeline.py --scene \"一只猫坐在沙发上\" --platform dalle --output result.json\n"
            "  python via54_pipeline.py --scene \"forest\" --variants 3\n"
            "  python via54_pipeline.py --image photo.jpg\n"
            "  python via54_pipeline.py --scene \"portrait\" --export-a1111\n"
            "  python via54_pipeline.py archive save --tags reference,nature\n"
            "  python via54_pipeline.py archive search \"mountain\"\n"
            "  python via54_pipeline.py archive list\n"
            "  python via54_pipeline.py --scene \"test\" --export-comfyui\n"
        ),
    )

    # Subcommands
    subparsers = parser.add_subparsers(dest="command")
    subparsers.required = False

    # --- archive subcommand ---
    archive_parser = subparsers.add_parser("archive", help="Manage prompt archive (JSONL)")
    archive_subparsers = archive_parser.add_subparsers(dest="archive_cmd", required=True)

    # archive save
    save_parser = archive_subparsers.add_parser("save", help="Save a prompt to the archive")
    save_parser.add_argument("--scene", "-s", type=str, default="", help="Scene description")
    save_parser.add_argument("--platform", "-p", type=str, default="midjourney", help="Target platform")
    save_parser.add_argument("--tags", type=str, default="", help="Comma-separated tags")
    save_parser.add_argument("--key", type=str, default=None, help="API key")
    save_parser.add_argument("--endpoint", type=str, default=None, help="API endpoint")
    save_parser.add_argument("--model", type=str, default=None, help="LLM model")
    save_parser.add_argument("--binary", type=str, default=None, help="Path to via54 binary")
    save_parser.add_argument("--output", "-o", type=str, default=None, help="Output JSON file path")
    save_parser.set_defaults(archive_mode="save")

    # archive list
    list_parser = archive_subparsers.add_parser("list", help="List recent archive entries")
    list_parser.add_argument("--limit", "-n", type=int, default=20, help="Number of entries to show")
    list_parser.set_defaults(archive_mode="list")

    # archive search
    search_parser = archive_subparsers.add_parser("search", help="Search archive entries")
    search_parser.add_argument("query", type=str, help="Search term")
    search_parser.add_argument("--limit", "-n", type=int, default=10, help="Max results")
    search_parser.set_defaults(archive_mode="search")

    # archive delete
    delete_parser = archive_subparsers.add_parser("delete", help="Delete an archive entry by ID")
    delete_parser.add_argument("id", type=str, help="Record ID to delete")
    delete_parser.set_defaults(archive_mode="delete")

    # --- Main pipeline arguments ---
    parser.add_argument(
        "--scene",
        "-s",
        type=str,
        default=None,
        help="Scene description (Chinese or English)",
    )

    parser.add_argument(
        "--platform",
        "-p",
        type=str,
        default="midjourney",
        choices=["midjourney", "flux", "dalle3", "sd3", "stable_diffusion", "ideogram", "recraft", "seedance", "gemini", "veo", "sora", "kling", "pika", "jimeng"],
        help="Target image generation platform (default: midjourney)",
    )

    parser.add_argument(
        "--ref-image",
        "-r",
        type=str,
        default=None,
        help="Optional reference image path",
    )

    parser.add_argument(
        "--output",
        "-o",
        type=str,
        default=None,
        help="Output JSON file path for the enhanced prompt scaffold",
    )

    parser.add_argument(
        "--endpoint",
        type=str,
        default=None,
        help="OpenAI-compatible API endpoint (overrides VIA54_LLM_ENDPOINT)",
    )

    parser.add_argument(
        "--key",
        type=str,
        default=None,
        help="API key (overrides VIA54_LLM_KEY / VIA54_LLM_KEY)",
    )

    parser.add_argument(
        "--model",
        type=str,
        default=None,
        help="LLM model name (overrides VIA54_LLM_MODEL)",
    )

    parser.add_argument(
        "--binary",
        type=str,
        default=None,
        help="Path to via54 binary (overrides VIA54_BINARY)",
    )
    parser.add_argument(
        "--provider",
        type=str,
        default="openai",
        choices=["openai", "deepseek", "ollama", "hermes", "local"],
        help="LLM provider preset (openai, deepseek, ollama, hermes, local). Sets default endpoint + model.",
    )

    # --- P4: Image reverse engineering ---
    parser.add_argument(
        "--image",
        type=str,
        default=None,
        help="P4: Path to image for reverse prompt engineering (vision LLM analysis)",
    )

    # --- P5: Template variant expansion ---
    parser.add_argument(
        "--variants",
        type=int,
        default=None,
        help="P5: Generate N variant expansions from {option|syntax} patterns",
    )

    # --- P8: Export flags ---
    parser.add_argument(
        "--export-a1111",
        action="store_true",
        default=False,
        help="P8: Output in Automatic1111 / SD.Next prompt format",
    )

    parser.add_argument(
        "--export-comfyui",
        action="store_true",
        default=False,
        help="P8: Output in ComfyUI CLIPTextEncode prompt format",
    )

    # --- Verbose ---
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        default=False,
        help="Print detailed progress information",
    )

    # Parse
    args = parser.parse_args(argv)

    # If no subcommand but --scene is missing and --image is missing, show help
    if args.command is None and args.scene is None and args.image is None:
        parser.print_help()
        sys.exit(1)

    return args


def _cmd_archive(args: argparse.Namespace, config: Config) -> int:
    """Handle archive subcommands."""
    archive = PromptArchive()

    if args.archive_mode == "save":
        scene = args.scene
        if not scene:
            print("Error: --scene is required for archive save", file=sys.stderr)
            return 1

        # Resolve API key
        api_key = args.key or config.llm_key
        preset = PROVIDER_PRESETS.get(args.provider, PROVIDER_PRESETS["openai"])
        if not api_key and preset.get("key_required", True):
            print(
                f"Error: Provider '{args.provider}' requires an API key. "
                "Set VIA54_LLM_KEY or pass --key.",
                file=sys.stderr,
            )
            return 1

        endpoint = args.endpoint or config.llm_endpoint
        model = args.model or config.llm_model
        binary = args.binary or config.via54_binary

        # Run the pipeline
        try:
            result = pipeline(
                scene=scene,
                platform=args.platform,
                api_key=api_key,
                api_endpoint=endpoint,
                model=model,
                binary=binary,
                output_path=args.output,
            )
        except RuntimeError as e:
            print(f"Error: {e}", file=sys.stderr)
            return 1

        # Parse tags
        tags = [t.strip() for t in args.tags.split(",") if t.strip()] if args.tags else []

        # Save to archive
        record_id = archive.save(result, tags=tags)
        print(f"Saved to archive: id={record_id}")
        if result.raw_prompt:
            print(f"Prompt: {result.raw_prompt[:100]}{'...' if len(result.raw_prompt) > 100 else ''}")
        return 0

    elif args.archive_mode == "list":
        records = archive.list(recent=args.limit)
        if not records:
            print("Archive is empty.")
            return 0
        print(f"Recent {len(records)} archive entries:\n")
        for r in records:
            tags_str = ", ".join(r.get("tags", [])) if r.get("tags") else "(no tags)"
            scene_preview = r.get("scene", "")[:60]
            print(f"  [{r['id']}] {r['created_at'][:19]} | {scene_preview}")
            print(f"        Tags: {tags_str}")
            print()
        return 0

    elif args.archive_mode == "search":
        results = archive.search(args.query, limit=args.limit)
        if not results:
            print(f"No results for query: {args.query}")
            return 0
        print(f"Found {len(results)} result(s) for '{args.query}':\n")
        for r in results:
            tags_str = ", ".join(r.get("tags", [])) if r.get("tags") else "(no tags)"
            scene_preview = r.get("scene", "")[:80]
            print(f"  [{r['id']}] {r['created_at'][:19]}")
            print(f"        Scene: {scene_preview}")
            print(f"        Tags: {tags_str}")
            print()
        return 0

    elif args.archive_mode == "delete":
        found = archive.delete(args.id)
        if found:
            print(f"Deleted record: {args.id}")
            return 0
        else:
            print(f"Record not found: {args.id}", file=sys.stderr)
            return 1

    return 0


def main(argv: Optional[List[str]] = None) -> int:
    """CLI entry point: parse args, load config, run pipeline, print result.

    Args:
        argv: Command-line arguments (defaults to sys.argv[1:]).

    Returns:
        Exit code (0 on success, 1 on error).
    """
    args = parse_args(argv)
    config = Config.from_env(provider=args.provider)

    # Handle archive subcommand
    if args.command == "archive":
        return _cmd_archive(args, config)

    # CLI flags override config
    endpoint = args.endpoint or config.llm_endpoint
    api_key = args.key or config.llm_key
    model = args.model or config.llm_model
    binary = args.binary or config.via54_binary

    # Validate API key (needed for P0, P1, P4)
    if not api_key:
        # Check if we're in a mode that doesn't need API keys
        # Variant expansion (--variants) works without LLM
        needs_llm = bool(args.image)
        if not needs_llm and args.scene and args.variants:
            needs_llm = False  # variants mode is local-only
        elif args.scene and not args.variants:
            needs_llm = True   # scene mode needs LLM for expansion
        preset = PROVIDER_PRESETS.get(args.provider, PROVIDER_PRESETS["openai"])
        if needs_llm and preset.get("key_required", True):
            print(
                f"Error: Provider '{args.provider}' requires an API key. "
                "Set VIA54_LLM_KEY or pass --key.",
                file=sys.stderr,
            )
            return 1

    if args.verbose:
        if args.scene:
            print(f"[via54-pipeline] Scene: {args.scene[:80]}{'...' if len(args.scene) > 80 else ''}", file=sys.stderr)
        if args.image:
            print(f"[via54-pipeline] Image: {args.image}", file=sys.stderr)
        print(f"[via54-pipeline] Platform: {args.platform}", file=sys.stderr)
        print(f"[via54-pipeline] Endpoint: {endpoint}", file=sys.stderr)
        print(f"[via54-pipeline] Model: {model}", file=sys.stderr)
        print(f"[via54-pipeline] Binary: {binary}", file=sys.stderr)

    # ------------------------------------------------------------------
    # P4 — Reverse Image → Prompt
    # ------------------------------------------------------------------
    if args.image:
        if args.verbose:
            print(f"[via54-pipeline] P4: Analyzing image with vision LLM...", file=sys.stderr)
        try:
            expansion = reverse_image(
                args.image,
                api_key=api_key,
                endpoint=endpoint,
                model=model,
            )
        except RuntimeError as e:
            print(f"Error: {e}", file=sys.stderr)
            return 1
        except FileNotFoundError:
            print(f"Error: Image not found — '{args.image}'", file=sys.stderr)
            return 1

        scaffold = PromptScaffold()
        scaffold.scene = f"[Reverse from: {args.image}]"
        scaffold.platform = args.platform
        scaffold.fields = expansion.get("fields", {})
        scaffold.negative = expansion.get("negative", [])
        scaffold.raw_prompt = _build_raw_prompt(scaffold)

        # If --export-a1111 or --export-comfyui, output that format instead
        if args.export_a1111:
            print(export_a1111(scaffold))
        elif args.export_comfyui:
            print(export_comfyui_clip(scaffold))
        else:
            print(json.dumps(scaffold.to_dict(), indent=2, ensure_ascii=False))

        if args.verbose:
            print(f"[via54-pipeline] P4 Done.", file=sys.stderr)
        return 0

    # ------------------------------------------------------------------
    # Main pipeline (requires --scene)
    # ------------------------------------------------------------------
    if args.scene is None:
        print("Error: --scene is required (or use --image for reverse engineering)", file=sys.stderr)
        return 1

    try:
        # P5 — Template variant expansion (pre-process scene before pipeline)
        if args.variants is not None and args.variants > 0:
            if args.verbose:
                print(f"[via54-pipeline] P5: Expanding {args.variants} variants...", file=sys.stderr)
            expanded_scenes = expand_variants(args.scene, count=args.variants)

            # If no API key, output variants directly (no LLM expansion)
            if not api_key:
                print(f"Variants ({len(expanded_scenes)}):")
                for i, s in enumerate(expanded_scenes, 1):
                    print(f"  {i}. {s}")
                return 0
        else:
            expanded_scenes = [args.scene]

        results: List[PromptScaffold] = []
        for i, scene_variant in enumerate(expanded_scenes):
            if args.verbose and len(expanded_scenes) > 1:
                print(f"[via54-pipeline] Variant {i+1}/{len(expanded_scenes)}", file=sys.stderr)

            result = pipeline(
                scene=scene_variant,
                platform=args.platform,
                ref_image=args.ref_image,
                api_key=api_key,
                api_endpoint=endpoint,
                model=model,
                binary=binary,
                output_path=args.output,
            )
            results.append(result)

    except RuntimeError as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1
    except subprocess.TimeoutExpired:
        print("Error: via54 binary timed out (120s)", file=sys.stderr)
        return 1
    except FileNotFoundError as e:
        print(
            f"Error: Binary not found — '{binary}'. "
            f"Set VIA54_BINARY or ensure via54.exe is in PATH.",
            file=sys.stderr,
        )
        return 1

    # ------------------------------------------------------------------
    # Output
    # ------------------------------------------------------------------
    for i, result in enumerate(results):
        if len(results) > 1:
            print(f"--- Variant {i+1} ---")

        if args.export_a1111:
            # P8 — A1111 format
            print(export_a1111(result))
        elif args.export_comfyui:
            # P8 — ComfyUI format
            print(export_comfyui_clip(result))
        elif result.raw_prompt:
            print(result.raw_prompt)
        else:
            print(json.dumps(result.to_dict(), indent=2, ensure_ascii=False))

    if args.verbose:
        total = len(results)
        dims = len(results[0].fields) if results else 0
        print(f"[via54-pipeline] Done. {total} result(s), {dims} dimensions populated.", file=sys.stderr)
        if results and results[0].original_scene != results[0].scene:
            preview = results[0].original_scene[:60]
            translated = results[0].scene[:60]
            print(
                f"[via54-pipeline] Translated: '{preview}' → '{translated}'",
                file=sys.stderr,
            )

    return 0


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------

if __name__ == "__main__":
    sys.exit(main())
