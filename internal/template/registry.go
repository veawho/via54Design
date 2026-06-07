// SPDX-License-Identifier: MIT OR AGPL-3.0

package template

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

type Registry struct {
	BaseDir string
	Data    *TemplateRegistry
}

func NewRegistry(baseDir string) (*Registry, error) {
	r := &Registry{BaseDir: baseDir}
	if err := r.load(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) load() error {
	candidates := []string{
		filepath.Join(r.BaseDir, "templates", "registry.yaml"),
		filepath.Join(r.BaseDir, "core", "template-registry.yaml"),
		filepath.Join(r.BaseDir, "template-registry.yaml"),
	}
	var data []byte
	var err error
	for _, path := range candidates {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("register not found: %w", err)
	}
	r.Data = &TemplateRegistry{}
	return yaml.Unmarshal(data, r.Data)
}

func (r *Registry) ResolveLayout(id string) (string, error) {
	return r.resolve(r.Data.Layouts, id)
}

func (r *Registry) ResolveColorScheme(id string) (string, error) {
	return r.resolve(r.Data.ColorSchemes, id)
}

func (r *Registry) ResolveTypography(id string) (string, error) {
	return r.resolve(r.Data.Typography, id)
}

func (r *Registry) ResolveNarratology(id string) (string, error) {
	return r.resolve(r.Data.Narratology, id)
}

func (r *Registry) resolve(entries []RegistryEntry, id string) (string, error) {
	for _, e := range entries {
		if e.ID == id {
			for _, base := range []string{r.BaseDir, filepath.Join(r.BaseDir, "templates")} {
				p := filepath.Join(base, e.File)
				if _, err := os.Stat(p); err == nil {
					return p, nil
				}
			}
			return "", fmt.Errorf("template file not found: %s", e.File)
		}
	}
	return "", fmt.Errorf("template not found: %s", id)
}

func (r *Registry) ListAll() map[string][]RegistryEntry {
	return map[string][]RegistryEntry{
		"layouts":        r.Data.Layouts,
		"color_schemes":  r.Data.ColorSchemes,
		"typography":     r.Data.Typography,
		"narratology":    r.Data.Narratology,
	}
}

func (r *Registry) ListNarratology() []RegistryEntry {
	return r.Data.Narratology
}
