// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ── Provider presets ──

// ProviderPreset defines an LLM provider preset.
type ProviderPreset struct {
	Endpoint    string `json:"endpoint"`
	Model       string `json:"model"`
	KeyRequired bool   `json:"key_required"`
	Description string `json:"description"`
}

// ProviderPresets is the built-in provider table.
var ProviderPresets = map[string]ProviderPreset{
	"openai": {
		Endpoint:    "https://api.openai.com/v1",
		Model:       "gpt-4o-mini",
		KeyRequired: true,
		Description: "OpenAI GPT-4o / GPT-4o-mini (default)",
	},
	"deepseek": {
		Endpoint:    "https://api.deepseek.com/v1",
		Model:       "deepseek-chat",
		KeyRequired: true,
		Description: "DeepSeek Chat / DeepSeek V3",
	},
	"ollama": {
		Endpoint:    "http://localhost:11434/v1",
		Model:       "llama3.2",
		KeyRequired: false,
		Description: "Local Ollama (llama3.2, qwen2.5, etc.) — no API key needed",
	},
	"hermes": {
		Endpoint:    "http://localhost:18791/v1",
		Model:       "deepseek-v4-flash",
		KeyRequired: false,
		Description: "Hermes Agent gateway proxy (port 18791)",
	},
	"local": {
		Endpoint:    "http://localhost:8000/v1",
		Model:       "local-model",
		KeyRequired: false,
		Description: "Generic local OpenAI-compatible server (vLLM, llama.cpp, etc.)",
	},
}

// ── Dimension fields (26 image + 10 video) ──

// DimensionFields lists the 36 semantic expansion dimensions.
var DimensionFields = []string{
	// core subject
	"subject",   // 主体对象
	"secondary", // 辅助元素
	// style
	"art_movement", // 风格流派
	"artist_ref",   // 艺术家参考
	"medium",       // 媒介
	"genre",        // 类型/题材
	"hair",         // 发型/发色
	"pose",         // 姿态/动态
	// composition
	"camera_shot",      // 景别
	"composition_type", // 构图
	"depth_of_field",   // 景深
	"view",             // 视角/视点
	"format",           // 画幅
	// lighting
	"lighting", // 光线
	// color
	"color_palette", // 色彩
	// environment
	"environment", // 环境
	"weather",     // 天气
	"era",         // 时代
	"time",        // 时间
	// detail
	"texture",  // 纹理
	"effects",  // 效果
	"material", // 材质
	"face",     // 面部
	"detail",   // 细节
	// quality
	"quality_tags", // 质量标签
	"emotion",      // 情绪/氛围
	// video control
	"camera_movement",  // 运镜类型
	"motion_intensity", // 运动强度
	"frame_count",      // 帧数
	"fps",              // 帧率
	"shot_size",        // 景别尺寸
	"angle",            // 拍摄角度
	"duration_seconds", // 时长
	"keyframe",         // 关键帧
	"transition",       // 转场
	"motion_blur",      // 运动模糊
}

// ── System prompts ──

const SystemPromptFill = `You are a professional AI image and video prompt engineer. Given a scene description, fill the 36 dimensions below with specific, vivid, English-language values. Each value should be 2-8 words, descriptive and concrete. Return ONLY a JSON object with the field values.`

const SystemPromptTranslate = `You are a professional translator. Translate the following Chinese scene description to English. Preserve all artistic and technical nuance. Return ONLY the English translation, no commentary.`

const SystemPromptReverse = `You are a professional AI image and video prompt engineer specializing in reverse prompt engineering. Analyze this image and infer the 36 prompt dimensions below with specific, vivid, English-language values. Each value should be 2-8 words, descriptive and concrete. Return ONLY a JSON object with the field values.`

// ── ChatMessage types ──

// ChatMessage represents a message in a chat completion request.
type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string or []ContentPart for vision
}

// ContentPart is a part of a multimodal message.
type ContentPart struct {
	Type     string        `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageURL *ImageURLPart `json:"image_url,omitempty"`
}

// ImageURLPart holds a base64-encoded image URL.
type ImageURLPart struct {
	URL string `json:"url"`
}

// chatCompletionRequest is the OpenAI-compatible request body.
type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

// chatCompletionResponse is the OpenAI-compatible response body.
type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ── LLM API call ──

// CallLLM sends a text-only chat completion request to an OpenAI-compatible API.
// If apiKey is empty, no Authorization header is sent (for local providers).
func CallLLM(messages []ChatMessage, endpoint, apiKey, model string) (string, error) {
	return callLLMInternal(messages, endpoint, apiKey, model, 120)
}

// callLLMInternal is the internal implementation with configurable timeout.
func callLLMInternal(messages []ChatMessage, endpoint, apiKey, model string, timeoutSec int) (string, error) {
	url := endpoint
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "chat/completions"

	body := chatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   2048,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM API connection failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API HTTP %d error: %s", resp.StatusCode, string(respBody[:min(len(respBody), 500)]))
	}

	var data chatCompletionResponse
	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if data.Error != nil && data.Error.Message != "" {
		return "", fmt.Errorf("LLM API error: %s", data.Error.Message)
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("LLM API returned no choices: %s", string(respBody[:min(len(respBody), 500)]))
	}

	return data.Choices[0].Message.Content, nil
}

// CallLLMVision sends a vision/multimodal chat completion request.
// contentParts should contain text and image_url parts.
func CallLLMVision(systemPrompt string, contentParts []ContentPart, endpoint, apiKey, model string) (string, error) {
	url := endpoint
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "chat/completions"

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: contentParts},
	}

	body := chatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   2048,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM API connection failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API HTTP %d error: %s", resp.StatusCode, string(respBody[:min(len(respBody), 500)]))
	}

	var data chatCompletionResponse
	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if data.Error != nil && data.Error.Message != "" {
		return "", fmt.Errorf("LLM API error: %s", data.Error.Message)
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("LLM API returned no choices: %s", string(respBody[:min(len(respBody), 500)]))
	}

	return data.Choices[0].Message.Content, nil
}

// ── Expansion parsing ──

// ExpansionResult holds the parsed LLM expansion output.
type ExpansionResult struct {
	Fields   map[string]string `json:"fields"`
	Negative []string          `json:"negative"`
	Raw      string            `json:"raw,omitempty"`
	Error    string            `json:"error,omitempty"`
}

// ParseLLMExpansion parses the LLM response into a structured expansion.
func ParseLLMExpansion(raw string) ExpansionResult {
	result := ExpansionResult{
		Fields:   make(map[string]string),
		Negative: []string{},
		Raw:      raw,
	}

	text := raw

	// Strip markdown code fences
	text = stripCodeFences(text)

	// Try to parse as JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		// Try regex fallback: find first JSON object
		parsed = findJSONObject(text)
		if parsed == nil {
			for _, f := range DimensionFields {
				result.Fields[f] = ""
			}
			result.Error = "Could not parse LLM response as JSON"
			return result
		}
	}

	// Extract fields (supports both {"fields": {...}} and flat JSON)
	fieldsData := parsed
	if f, ok := parsed["fields"]; ok {
		if fm, ok := f.(map[string]interface{}); ok {
			fieldsData = fm
		}
	}

	for _, field := range DimensionFields {
		if val, ok := fieldsData[field]; ok {
			result.Fields[field] = fmt.Sprintf("%v", val)
		} else {
			result.Fields[field] = ""
		}
	}

	// Extract negative terms
	if neg, ok := parsed["negative"]; ok {
		switch v := neg.(type) {
		case []interface{}:
			for _, x := range v {
				if s := fmt.Sprintf("%v", x); s != "" {
					result.Negative = append(result.Negative, s)
				}
			}
		case string:
			if v != "" {
				result.Negative = []string{v}
			}
		}
	}

	return result
}

// stripCodeFences removes markdown code fences from a string.
func stripCodeFences(text string) string {
	if len(text) < 3 {
		return text
	}
	if text[:3] == "```" {
		firstNL := -1
		for i, c := range text {
			if c == '\n' {
				firstNL = i
				break
			}
		}
		if firstNL != -1 {
			text = text[firstNL+1:]
		}
		// Remove closing fence
		if len(text) >= 3 && text[len(text)-3:] == "```" {
			text = text[:len(text)-3]
		} else if idx := lastIndex(text, "```"); idx != -1 {
			text = text[:idx]
		}
	}
	return text
}

// lastIndex returns the last index of substr in s.
func lastIndex(s, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// findJSONObject attempts to find a JSON object in a string.
func findJSONObject(text string) map[string]interface{} {
	start := -1
	depth := 0
	for i, c := range text {
		if c == '{' {
			if start == -1 {
				start = i
			}
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 && start != -1 {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(text[start:i+1]), &parsed); err == nil {
					return parsed
				}
				start = -1
			}
		}
	}
	return nil
}

// min returns the minimum of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
