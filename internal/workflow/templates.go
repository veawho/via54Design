// via54Design — ComfyUI Workflow Template Loader
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadWorkflowTemplate loads a single workflow template by ID.
func LoadWorkflowTemplate(id, baseDir string) (*WorkflowTemplate, error) {
	candidates := []string{
		filepath.Join(baseDir, "templates", "workflows", id+".yaml"),
		filepath.Join(baseDir, "templates", "workflows", id+".yml"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var t WorkflowTemplate
		if err := yaml.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("parse workflow %s: %w", id, err)
		}
		if t.ID == "" {
			return nil, fmt.Errorf("workflow %s has empty id", id)
		}
		return &t, nil
	}
	return nil, fmt.Errorf("workflow template %q not found in templates/workflows/", id)
}

// ListWorkflowTemplates scans the workflows directory and returns all template IDs.
func ListWorkflowTemplates(baseDir string) ([]string, error) {
	dir := filepath.Join(baseDir, "templates", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read workflows directory: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "registry.yaml" || name == "registry.yml" {
			continue
		}
		if filepath.Ext(name) == ".yaml" || filepath.Ext(name) == ".yml" {
			id := name[:len(name)-len(filepath.Ext(name))]
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// LoadRegistry loads the workflow registry file.
func LoadRegistry(baseDir string) (*WorkflowRegistry, error) {
	candidates := []string{
		filepath.Join(baseDir, "templates", "workflows", "registry.yaml"),
		filepath.Join(baseDir, "templates", "workflows", "registry.yml"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var reg WorkflowRegistry
		if err := yaml.Unmarshal(data, &reg); err != nil {
			return nil, fmt.Errorf("parse registry: %w", err)
		}
		return &reg, nil
	}
	return nil, fmt.Errorf("workflow registry not found in templates/workflows/")
}
