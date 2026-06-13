// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ── Config ──

// Config holds the runtime configuration for the pipeline.
type Config struct {
	LLMEndpoint string
	LLMKey      string
	LLMModel    string
	Provider    string
}

// ConfigFromEnv loads configuration from environment variables and provider presets.
// VIA54_LLM_ENDPOINT, VIA54_LLM_KEY, VIA54_LLM_MODEL env vars are used.
func ConfigFromEnv(provider string) *Config {
	preset, ok := ProviderPresets[provider]
	if !ok {
		preset = ProviderPresets["openai"]
	}

	cfg := &Config{
		LLMEndpoint: preset.Endpoint,
		LLMModel:    preset.Model,
		Provider:    provider,
	}

	if ep := os.Getenv("VIA54_LLM_ENDPOINT"); ep != "" {
		cfg.LLMEndpoint = ep
	}
	if key := os.Getenv("VIA54_LLM_KEY"); key != "" {
		cfg.LLMKey = key
	}
	if model := os.Getenv("VIA54_LLM_MODEL"); model != "" {
		cfg.LLMModel = model
	}

	return cfg
}

// ProviderRequiresKey returns true if the current provider needs an API key.
func (c *Config) ProviderRequiresKey() bool {
	preset, ok := ProviderPresets[c.Provider]
	if !ok {
		return true
	}
	return preset.KeyRequired
}

// ── Pipeline orchestration ──

// Pipeline runs the full prompt enhancement pipeline.
// Steps:
//  1. Detect language → translate Chinese→English if needed.
//  2. Fill all dimension fields via LLM.
//  3. Build raw prompt from filled fields.
func Pipeline(scene, platform, provider, endpoint, apiKey, model string) (*PromptScaffold, error) {
	// Resolve provider preset if endpoint/model not explicitly set
	if provider != "" {
		preset, ok := ProviderPresets[provider]
		if ok {
			if endpoint == "" {
				endpoint = preset.Endpoint
			}
			if model == "" {
				model = preset.Model
			}
		}
	}

	scaffold := &PromptScaffold{
		Scene:         scene,
		Platform:      platform,
		Fields:        make(map[string]string),
		Negative:      []string{},
		OriginalScene: scene,
	}

	// Step 1: P1 — i18n Auto-Translate
	englishScene, _, err := TranslateToEnglish(scene, func(messages []ChatMessage) (string, error) {
		return CallLLM(messages, endpoint, apiKey, model)
	})
	if err != nil {
		return scaffold, fmt.Errorf("translation: %w", err)
	}
	scaffold.Scene = englishScene

	// Step 2: P0 — LLM Semantic Expansion
	expansion, err := ExpandWithLLM(englishScene, platform, apiKey, endpoint, model)
	if err != nil {
		return scaffold, fmt.Errorf("LLM expansion: %w", err)
	}
	scaffold.Fields = expansion.Fields
	if len(expansion.Negative) > 0 {
		scaffold.Negative = expansion.Negative
	}

	// Step 3: Build raw prompt
	scaffold.RawPrompt = BuildRawPrompt(scaffold)

	return scaffold, nil
}

// ExpandWithLLM fills all dimension fields via an LLM call.
func ExpandWithLLM(scene, platform, apiKey, endpoint, model string) (*ExpansionResult, error) {
	fieldsList := strings.Join(DimensionFields, ", ")
	userPrompt := fmt.Sprintf(
		"Platform: %s\nScene: %s\n\nFill these 36 dimensions with specific, vivid values:\n%s\n\nAlso provide 3-5 negative prompt terms as an array 'negative'.\nReturn a JSON object with keys \"fields\" (object) and \"negative\" (array).",
		platform, scene, fieldsList,
	)

	raw, err := CallLLM([]ChatMessage{
		{Role: "system", Content: SystemPromptFill},
		{Role: "user", Content: userPrompt},
	}, endpoint, apiKey, model)
	if err != nil {
		return nil, fmt.Errorf("LLM call: %w", err)
	}

	result := ParseLLMExpansion(raw)
	return &result, nil
}

// ReverseImage analyzes an image and returns filled scaffold via vision LLM.
func ReverseImage(imagePath, apiKey, endpoint, model string) (*ExpansionResult, error) {
	// Read and encode image
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(imagePath))
	mimeMap := map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".webp": "image/webp",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
	}
	mime := mimeMap[ext]
	if mime == "" {
		mime = "image/png"
	}

	// Base64 encode
	b64Data := encodeBase64(data)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mime, b64Data)

	fieldsList := strings.Join(DimensionFields, ", ")
	textContent := fmt.Sprintf(
		"Analyze this image and return JSON with 36 dimension fields:\n%s\n\nAlso provide 3-5 negative prompt terms as an array 'negative'.\nReturn a JSON object with keys \"fields\" (object) and \"negative\" (array).",
		fieldsList,
	)

	contentParts := []ContentPart{
		{Type: "text", Text: textContent},
		{Type: "image_url", ImageURL: &ImageURLPart{URL: dataURL}},
	}

	raw, err := CallLLMVision(SystemPromptReverse, contentParts, endpoint, apiKey, model)
	if err != nil {
		return nil, fmt.Errorf("vision LLM call: %w", err)
	}

	result := ParseLLMExpansion(raw)
	return &result, nil
}

// encodeBase64 base64-encodes binary data without external deps.
func encodeBase64(data []byte) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var b strings.Builder
	b.Grow((len(data) + 2) / 3 * 4)

	for i := 0; i < len(data); i += 3 {
		var val int
		if i < len(data) {
			val = int(data[i]) << 16
		}
		if i+1 < len(data) {
			val |= int(data[i+1]) << 8
		}
		if i+2 < len(data) {
			val |= int(data[i+2])
		}

		b.WriteByte(charset[(val>>18)&0x3F])
		b.WriteByte(charset[(val>>12)&0x3F])

		if i+1 < len(data) {
			b.WriteByte(charset[(val>>6)&0x3F])
		} else {
			b.WriteByte('=')
		}

		if i+2 < len(data) {
			b.WriteByte(charset[val&0x3F])
		} else {
			b.WriteByte('=')
		}
	}

	return b.String()
}

// ── JSON helpers ──

// ScaffoldToJSON serializes the scaffold to JSON bytes.
func ScaffoldToJSON(s *PromptScaffold) ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// ScaffoldFromJSON deserializes scaffold from JSON bytes.
func ScaffoldFromJSON(data []byte) (*PromptScaffold, error) {
	var s PromptScaffold
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal scaffold: %w", err)
	}
	if s.Fields == nil {
		s.Fields = make(map[string]string)
	}
	return &s, nil
}
