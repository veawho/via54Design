// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

package template

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"

	"github.com/veawho/via54Design/internal/util"
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
			return "", util.WrapNotFound(nil, "template file not found: %s", e.File)
		}
	}
	return "", util.WrapNotFound(nil, "template not found: %s", id)
}

func (r *Registry) ListAll() map[string][]RegistryEntry {
	return map[string][]RegistryEntry{
		"layouts":       r.Data.Layouts,
		"color_schemes": r.Data.ColorSchemes,
		"typography":    r.Data.Typography,
		"narratology":   r.Data.Narratology,
	}
}

func (r *Registry) ListNarratology() []RegistryEntry {
	return r.Data.Narratology
}
