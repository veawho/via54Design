// SPDX-License-Identifier: AGPL-3.0-only
package prompt

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ─── Bilingual Option parsing ────────────────────────────────────────────────
//
// YAML option values use the canonical "中文 (English)" form, e.g.:
//   value: "摄影写实 (Photorealistic)"
//   value: "赛博朋克 (Cyberpunk)"
//   value: "warm amber tones"          // pure English is also accepted
//   value: "暖琥珀色"                  // pure Chinese (label-only)
//
// i18n.go splits these into BilingualOption{ZH, EN, Raw} so the engine can
// render either language, and round-trip if a user reply matches either side.

// reBilingual matches "中文 (English)" with a single space between the ZH
// chunk and the parenthesised EN chunk. We deliberately allow the EN chunk
// to contain nested parens (e.g. "主观镜头 (POV, first person)") by using a
// non-greedy match on the parenthesised group.
var reBilingual = regexp.MustCompile(`^(.+?)\s*\(([^()]+)\)\s*$`)

// BilingualOption is the parsed view of a YAML option.
type BilingualOption struct {
	ZH  string // Chinese label, may be empty for label-only options
	EN  string // English value, may be empty
	Raw string // Original value, preserved as the single source of truth
}

// ParseBilingual splits a YAML option value into ZH + EN.
// Round-trip-safe: input "warm amber tones" → ZH="" EN="warm amber tones".
func ParseBilingual(raw string) BilingualOption {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return BilingualOption{Raw: raw}
	}
	m := reBilingual.FindStringSubmatch(raw)
	if m == nil {
		// No parens. Decide by Unicode script: if it contains CJK ideographs
		// treat as pure Chinese; otherwise pure English.
		if hasCJK(raw) {
			return BilingualOption{ZH: raw, EN: "", Raw: raw}
		}
		return BilingualOption{ZH: "", EN: raw, Raw: raw}
	}
	zh := strings.TrimSpace(m[1])
	en := strings.TrimSpace(m[2])
	// If the supposed "Chinese" half has no CJK characters, the parens were
	// probably part of an English phrase (e.g. "POV (first person)") — in
	// that case, neither half is Chinese; the whole raw is English.
	if !hasCJK(zh) {
		// Could be:
		//   (a) "POV (first person)" — neither half is Chinese, raw is English
		//   (b) "Photorealistic (写实)" — EN half is Chinese
		// Distinguish by checking the EN half: if it has CJK, swap so ZH gets
		// the Chinese half (case b). Otherwise, both halves are non-Chinese,
		// so the entire raw is the EN value (case a).
		if hasCJK(en) {
			zh, en = en, zh
		} else {
			// Neither half is Chinese — return raw as English.
			return BilingualOption{ZH: "", EN: raw, Raw: raw}
		}
	}
	return BilingualOption{ZH: zh, EN: en, Raw: raw}
}

// RenderZH returns the Chinese label, falling back to EN, falling back to Raw.
func (b BilingualOption) RenderZH() string {
	if b.ZH != "" {
		return b.ZH
	}
	if b.EN != "" {
		return b.EN
	}
	return b.Raw
}

// RenderEN returns the English value, falling back to ZH, falling back to Raw.
func (b BilingualOption) RenderEN() string {
	if b.EN != "" {
		return b.EN
	}
	if b.ZH != "" {
		return b.ZH
	}
	return b.Raw
}

// hasCJK returns true if s contains any CJK Unified Ideograph.
// Used to disambiguate "中文 (English)" from "POV (first person)".
func hasCJK(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// ─── Bilingual value normalisation for final-prompt composition ──────────────

// NormalizeFieldZH takes a field value and returns its Chinese rendering.
// If the value is in "中文 (English)" form, the ZH half is returned.
// If it's pure English, the original is returned (we don't translate here —
// that requires an LLM; use TranslateFreeTextZH for that).
func NormalizeFieldZH(value string) string {
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "（LLM填充") {
		// Unfilled placeholder — leave as-is so the renderer can mark it
		return value
	}
	return ParseBilingual(value).RenderZH()
}

// NormalizeFieldEN returns the English rendering of a bilingual field value.
func NormalizeFieldEN(value string) string {
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "（LLM填充") {
		return value
	}
	return ParseBilingual(value).RenderEN()
}

// ─── Translation result for free-text fields ─────────────────────────────────

// FreeTextTranslation is the result of translating one free-text field
// (subject, secondary, environment) from the user's language to the
// target language.
type FreeTextTranslation struct {
	EN  string `json:"en"`
	ZH  string `json:"zh"`
	Note string `json:"note,omitempty"` // optional translator's note
}

// ToJSON serialises for shell-pipe consumption by the Go engine.
func (f FreeTextTranslation) ToJSON() string {
	b, _ := json.Marshal(f)
	return string(b)
}

// FreeTextTranslateRequest is what gets piped to `hermes -z` for translation.
// HermeS returns a single JSON object on stdout.
type FreeTextTranslateRequest struct {
	FieldID  string `json:"field_id"`
	UserText string `json:"user_text"`
	RefDesc  string `json:"ref_desc,omitempty"`
	Platform string `json:"platform"`
	Mode     string `json:"mode"` // "en_zh" or "zh_en"
}

// ─── Shell-callable translator ───────────────────────────────────────────────
//
// The Go engine doesn't import any LLM client library — it shells out to
// `hermes -z` and parses the JSON response. This keeps the engine dependency-
// free and lets the user pick the LLM at runtime via Hermes config.

// TranslateFreeText shells out to `hermes -z` with a translation prompt.
// The prompt is constructed to be hermes-friendly (single JSON-line reply).
//
// Returns the translated text, or the original on error (best-effort).
func TranslateFreeText(userText, refDesc, platform, mode string) (string, error) {
	// We don't actually call hermes here — the engine never inlines a LLM
	// call. Instead, this function is a marker for tests and a reminder to
	// the Go engine's caller (feishu bot / inbox_watcher) to do the
	// translation externally. See TranslateViaHermes in
	// cmd/via54/prompt_cmd.go for the actual shell-out wrapper.
	return userText, fmt.Errorf("TranslateFreeText must be invoked via the CLI wrapper (cmd/via54/prompt_cmd.go TranslateViaHermes)")
}
