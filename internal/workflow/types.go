// via54Design — ComfyUI Workflow Types
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package workflow

// WorkflowTemplate represents a ComfyUI workflow template loaded from YAML.
type WorkflowTemplate struct {
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Model       string            `yaml:"model"`
	VAE         string            `yaml:"vae,omitempty"`
	Type        string            `yaml:"type"` // txt2img, img2img, txt2vid
	Params      map[string]interface{} `yaml:"params"`
	Nodes       map[string]string `yaml:"nodes"` // logical name -> ComfyUI class_type
	Models      []ModelDownload   `yaml:"models,omitempty"`
}

// ModelDownload describes a model to auto-download.
type ModelDownload struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}

// WorkflowRegistryEntry is an entry in the workflow registry.
type WorkflowRegistryEntry struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Model       string `yaml:"model"`
	Params      string `yaml:"params"`
}

// WorkflowRegistry is the top-level registry structure.
type WorkflowRegistry struct {
	Version   int                      `yaml:"version"`
	Workflows []WorkflowRegistryEntry  `yaml:"workflows"`
}

// ComfyUINode represents a single node in the ComfyUI API JSON format.
type ComfyUINode struct {
	ClassType string                 `json:"class_type"`
	Inputs    map[string]interface{} `json:"inputs"`
}

// BuildResult holds the built workflow JSON bytes and node metadata.
type BuildResult struct {
	JSON       []byte
	Injected   int // number of CLIPTextEncode nodes injected
	TemplateID string
}
