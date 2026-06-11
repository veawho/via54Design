// via54Design — 9-Ramp Color Scale Library
// Copyright (C) 2026  via54 (veawho)
//
// Reference: docs/design-audit.md §11 (色阶系统 — Claude 9-Ramp × 7-Stop = 63色)
// Standard : Claude Design / Anthropic brand — 9 luminance ramps × semantic role
// Validation: garden-skills (7.5k⭐) 23-theme token architecture
//
// SPDX-License-Identifier: AGPL-3.0-only

// Package nineramp implements the Claude 9-Ramp color system used by via54Design.
//
// Each ramp defines 9 luminance stops (0..8) of a single hue. The 9-Ramp
// system pairs with 7 semantic stops (50/100/200/300/400/500/600/700) to
// produce a 63-color token palette. Roles: bg, surface, border, text, accent.
package nineramp

import "fmt"

// Stop is a single color stop on a 0..8 luminance scale.
type Stop struct {
	Index int    // 0 = darkest, 8 = lightest
	Hex   string // #RRGGBB
}

// NineRamp is a 9-stop color ramp of a single hue.
type NineRamp struct {
	Name  string // e.g. "claude-orange"
	Stops [9]Stop
}

// Hex returns the hex string for a given stop index (0..8), or "" if out of range.
func (r *NineRamp) Hex(i int) string {
	if i < 0 || i > 8 {
		return ""
	}
	return r.Stops[i].Hex
}

// Token resolves a 7-stop semantic alias (e.g. 500) onto the 9-Ramp.
// Mapping: 50→1, 100→2, 200→3, 300→4, 400→5, 500→6, 600→7.
func (r *NineRamp) Token(semantic int) string {
	idx := -1
	switch semantic {
	case 50:
		idx = 1
	case 100:
		idx = 2
	case 200:
		idx = 3
	case 300:
		idx = 4
	case 400:
		idx = 5
	case 500:
		idx = 6
	case 600:
		idx = 7
	}
	if idx < 0 {
		return ""
	}
	return r.Hex(idx)
}

// Claude9Ramps is the canonical via54Design 9-Ramp palette
// (values match garden-skills Claude-compatible ramp set).
var Claude9Ramps = []NineRamp{
	{Name: "claude-orange", Stops: stops("#3D1308", "#5C1F11", "#7A2B1A", "#9B3D26", "#C25538", "#D97559", "#E89679", "#F2B89D", "#FAD9C5")},
	{Name: "claude-cream", Stops: stops("#1F1B16", "#33291F", "#4A3B2A", "#66523B", "#856D50", "#A78B69", "#C5AB89", "#DECBAF", "#F0E6D4")},
	{Name: "claude-sage", Stops: stops("#0F1A14", "#19271E", "#243629", "#324837", "#45634C", "#5A7E63", "#7A9A82", "#9FB7A4", "#C5D5C7")},
	{Name: "claude-blue", Stops: stops("#0B1622", "#13243A", "#1B3354", "#26487A", "#3462A3", "#4A7DBF", "#6E9AD2", "#94B8E0", "#BCD3EC")},
	{Name: "claude-red", Stops: stops("#2A0A0A", "#451313", "#5F1D1D", "#7E2A2A", "#A23A3A", "#C25252", "#D87575", "#E89C9C", "#F4C2C2")},
	{Name: "claude-green", Stops: stops("#0A1A12", "#13291D", "#1C3A29", "#295039", "#3A6B4D", "#508763", "#6CA47D", "#8EBF9A", "#B3D7B8")},
	{Name: "claude-purple", Stops: stops("#150F1F", "#231A33", "#322547", "#463260", "#5C427D", "#755896", "#8E72AE", "#A990C2", "#C4B1D5")},
	{Name: "claude-neutral", Stops: stops("#0E0E0E", "#1C1C1C", "#2A2A2A", "#3D3D3D", "#555555", "#707070", "#909090", "#B5B5B5", "#D8D8D8")},
	{Name: "claude-amber", Stops: stops("#1F1505", "#33240B", "#4A3514", "#66481F", "#85602E", "#A77A45", "#C59763", "#DDB58A", "#EDD2B3")},
}

func stops(v ...string) [9]Stop {
	var s [9]Stop
	for i, h := range v {
		s[i] = Stop{Index: i, Hex: h}
	}
	return s
}

// ByName returns a ramp by its name, or false if not found.
func ByName(name string) (NineRamp, bool) {
	for _, r := range Claude9Ramps {
		if r.Name == name {
			return r, true
		}
	}
	return NineRamp{}, false
}

// CSSVars emits a `--ramp-<name>-<idx>` block for embedding in stylesheets.
func CSSVars(r NineRamp) string {
	out := ""
	for _, s := range r.Stops {
		out += fmt.Sprintf("  --ramp-%s-%d: %s;\n", r.Name, s.Index, s.Hex)
	}
	return out
}
