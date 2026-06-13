// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package pipeline

import (
	"math/rand"
	"regexp"
	"strings"
)

var variantPattern = regexp.MustCompile(`\{([^}]+)\}`)

// ExpandVariants expands {opt1|opt2|opt3} template syntax into variants.
// When count=1, picks one random option per group.
// When count>1, generates N unique combinations (up to total possibilities).
func ExpandVariants(scene string, count int) []string {
	matches := variantPattern.FindAllStringSubmatchIndex(scene, -1)
	if len(matches) == 0 {
		// No variants found; return the original scene (repeated if count>1)
		result := make([]string, count)
		for i := range result {
			result[i] = scene
		}
		return result
	}

	// Parse each group into a list of options
	var optionLists [][]string
	for _, m := range matches {
		content := scene[m[2]:m[3]]
		opts := strings.Split(content, "|")
		var cleaned []string
		for _, o := range opts {
			o = strings.TrimSpace(o)
			if o != "" {
				cleaned = append(cleaned, o)
			}
		}
		if len(cleaned) > 0 {
			optionLists = append(optionLists, cleaned)
		}
	}

	if len(optionLists) == 0 {
		result := make([]string, count)
		for i := range result {
			result[i] = scene
		}
		return result
	}

	// Compute total possible combinations
	total := 1
	for _, opts := range optionLists {
		total *= len(opts)
	}

	if count == 1 {
		indices := make([]int, len(optionLists))
		for i, opts := range optionLists {
			indices[i] = rand.Intn(len(opts))
		}
		return []string{makeVariant(scene, matches, optionLists, indices)}
	}

	// Generate multiple unique variants
	var results []string
	seen := make(map[string]bool)
	attempts := 0
	maxAttempts := count*3 + total*3

	for len(results) < count && attempts < maxAttempts {
		indices := make([]int, len(optionLists))
		key := strings.Builder{}
		for i, opts := range optionLists {
			idx := rand.Intn(len(opts))
			indices[i] = idx
			key.WriteString(string(rune('0' + idx)))
			key.WriteRune(',')
		}
		keyStr := key.String()
		if !seen[keyStr] {
			seen[keyStr] = true
			results = append(results, makeVariant(scene, matches, optionLists, indices))
		}
		attempts++
	}

	return results
}

// makeVariant builds a single variant by replacing each {group} with the
// option at the corresponding index.
func makeVariant(scene string, matches [][]int, optionLists [][]string, indices []int) string {
	var b strings.Builder
	lastEnd := 0

	for i, m := range matches {
		start, end := m[0], m[1]
		// Write text before this match
		b.WriteString(scene[lastEnd:start])

		// Write the chosen option
		if i < len(indices) && i < len(optionLists) {
			opts := optionLists[i]
			idx := indices[i] % len(opts)
			b.WriteString(opts[idx])
		} else {
			// Fallback: write the original content
			b.WriteString(scene[m[2]:m[3]])
		}

		lastEnd = end
	}

	// Write remaining text
	b.WriteString(scene[lastEnd:])

	return b.String()
}
