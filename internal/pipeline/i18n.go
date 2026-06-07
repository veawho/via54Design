// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package pipeline

import (
	"fmt"
	"strings"
)

// ContainsChinese detects if text contains any CJK Unified Ideographs.
func ContainsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// TranslateToEnglish translates Chinese text to English via LLM.
// If the text contains no Chinese characters, returns it unchanged.
// llmFunc is a function that calls the LLM (typically CallLLM).
func TranslateToEnglish(text string, llmFunc func(messages []ChatMessage) (string, error)) (string, bool, error) {
	if !ContainsChinese(text) {
		return text, false, nil
	}

	prompt := fmt.Sprintf("Translate this Chinese scene description to English:\n\n%s", text)
	translated, err := llmFunc([]ChatMessage{
		{Role: "system", Content: SystemPromptTranslate},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return text, false, fmt.Errorf("translate: %w", err)
	}

	translated = strings.Trim(translated, "\"' ")
	return translated, true, nil
}
