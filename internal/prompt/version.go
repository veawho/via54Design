// SPDX-License-Identifier: AGPL-3.0-only
package prompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func SaveVersion(dir string, s *PromptScaffold) (string, error) {
	os.MkdirAll(dir, 0755)
	v := PromptVersion{
		Version: fmt.Sprintf("v%d", len(ListVersions(dir))+1),
		Timestamp: time.Now().Format(time.RFC3339),
		Seed: s.Seed, Platform: s.Platform, Prompt: s.FinalPrompt,
	}
	data, _ := json.Marshal(v)
	fp := filepath.Join(dir, fmt.Sprintf("prompt-%s.json", v.Version))
	return fp, os.WriteFile(fp, data, 0644)
}

func ListVersions(dir string) []string {
	var v []string
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "prompt-v") && strings.HasSuffix(e.Name(), ".json") { v = append(v, e.Name()) }
	}
	sort.Strings(v)
	return v
}

func DiffVersions(dir, v1, v2 string) (string, error) {
	var pv1, pv2 PromptVersion
	d1, _ := os.ReadFile(filepath.Join(dir, "prompt-"+v1+".json"))
	d2, _ := os.ReadFile(filepath.Join(dir, "prompt-"+v2+".json"))
	json.Unmarshal(d1, &pv1); json.Unmarshal(d2, &pv2)
	if pv1.Prompt == pv2.Prompt { return "相同 (no changes)\n", nil }
	return fmt.Sprintf("差异: %s → %s\n之前: %s\n之后: %s\n", v1, v2, pv1.Prompt, pv2.Prompt), nil
}
