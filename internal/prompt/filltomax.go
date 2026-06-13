// SPDX-License-Identifier: AGPL-3.0-only
//
// filltomax.go — v2.3 限额最优填充算法
//
// Goal: given a partially-resolved prompt (FinalEN shorter than MaxChars),
//       generate instructions for an LLM to expand it to ~95% of MaxChars
//       by filling the unfilled "（LLM填充）" sections with vivid detail,
//       while preserving the ref-lock anchor if present.
//
// The algorithm:
//   1. Sort unfilled sections by priority (Weighted + Section.Weight desc)
//   2. Compute remaining budget = MaxChars - len(FinalEN) - 50 (safety margin)
//   3. For each unfilled section, allocate a target length proportional to
//      its weight (high-weight = more chars)
//   4. Emit per-section expansion instructions keyed to the budget
//   5. Caller invokes hermes -z with these instructions, gets English
//      expansion text back, splices into Overrides, re-calls GenerateI2I
//
// Why not call LLM here?
//   - Keeps the Go engine hermes-independent at the binary level
//   - Caller (feishu bot / inbox_watcher) already has an LLM session;
//     routing a second hermes subprocess for one field wastes 8-15s
//     of subprocess startup vs. piggy-backing on the existing session.
//   - The Go engine stays deterministic + testable (no LLM in unit tests)
package prompt

import (
	"os"
	"sort"
	"strings"
)

// FillTarget describes one section that the LLM should expand.
type FillTarget struct {
	SectionID   string `json:"section_id"`
	Hint        string `json:"hint"`
	TargetChars int    `json:"target_chars"`
	RefAnchor   string `json:"ref_anchor,omitempty"` // present if this section is ref-locked
	Order       int    `json:"order"`               // 1 = fill first
	Weight      float64 `json:"weight"`
}

// FillPlan is what the caller hands to the LLM to drive expansion.
type FillPlan struct {
	Platform        string       `json:"platform"`
	MaxChars        int          `json:"max_chars"`
	CurrentChars    int          `json:"current_chars"`
	RemainingBudget int          `json:"remaining_budget"`
	Targets         []FillTarget `json:"targets"`
	RefDescription  string       `json:"ref_description,omitempty"`
	StyleDefaults   string       `json:"style_defaults,omitempty"` // comma-joined already-filled EN sections
}

// PlanFill computes which sections to expand and by how much.
//
//   result, _ := prompt.GenerateI2I(req)
//   plan := prompt.PlanFill(result, req.RefDescription)
//   // hand plan to LLM; LLM returns map[section_id] = expanded_text
//   // splice into req.Overrides[section_id] and re-call GenerateI2I
func PlanFill(result *I2IResult, refDesc string) *FillPlan {
	plan := &FillPlan{
		Platform:       result.Platform,
		MaxChars:       result.MaxChars,
		CurrentChars:   result.FinalChars,
		RefDescription: refDesc,
	}
	plan.RemainingBudget = result.MaxChars - result.FinalChars - 50
	if plan.RemainingBudget < 100 {
		plan.RemainingBudget = 0
		// Already at or over budget — no fill needed.
		return plan
	}
	// Collect filled sections to give the LLM style context
	for _, val := range result.Sections {
		if val != "" && !strings.HasPrefix(val, "（LLM") {
			if plan.StyleDefaults == "" {
				plan.StyleDefaults = val
			} else if len(plan.StyleDefaults) < 200 {
				plan.StyleDefaults += ", " + val
			}
		}
	}

	// Build priority list from unfilled
	type pri struct {
		id     string
		weight float64
		hint   string
	}
	var pris []pri
	for _, sec := range loadTemplate(result.Platform, baseDirGuess()).Sections {
		matched := false
		for _, u := range result.Unfilled {
			if sec.ID == u {
				matched = true
				break
			}
		}
		if matched {
			pris = append(pris, pri{id: sec.ID, weight: sec.Weight, hint: sec.Hint})
		}
	}
	// Sort by weight desc (high-weight = expand more)
	sort.Slice(pris, func(i, j int) bool { return pris[i].weight > pris[j].weight })

	// Allocate budget proportional to weight
	totalW := 0.0
	for _, p := range pris {
		if p.weight == 0 {
			p.weight = 1.0
		}
		totalW += p.weight
	}
	order := 1
	for _, p := range pris {
		if p.weight == 0 {
			p.weight = 1.0
		}
		share := int(float64(plan.RemainingBudget) * p.weight / totalW)
		if share < 30 {
			share = 30 // floor
		}
		if share > 350 {
			share = 350 // ceiling per field
		}
		tgt := FillTarget{
			SectionID:   p.id,
			Hint:        p.hint,
			TargetChars: share,
			Order:       order,
			Weight:      p.weight,
		}
		if refDesc != "" {
			tgt.RefAnchor = refDesc
		}
		plan.Targets = append(plan.Targets, tgt)
		order++
	}
	return plan
}

// baseDirGuess returns a best-effort guess for template lookup. The engine
// normally receives BaseDir via I2IRequest, but PlanFill can be called
// standalone (e.g. from tests). When empty, we look for the templates dir
// in CWD and parent dirs. Returns "" if not found.
func baseDirGuess() string {
	// Try CWD first
	if _, err := os.ReadDir("templates/prompts"); err == nil {
		return "."
	}
	// Try parents up to 3 levels
	for _, d := range []string{"..", "../..", "../../.."} {
		if _, err := os.ReadDir(d + "/templates/prompts"); err == nil {
			return d
		}
	}
	return ""
}
