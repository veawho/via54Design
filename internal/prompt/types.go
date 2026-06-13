// SPDX-License-Identifier: AGPL-3.0-only
package prompt

// PromptTemplate 16维度提示词模板 (从 YAML 加载)
type PromptTemplate struct {
	ID          string            `yaml:"id"`
	Name        map[string]string `yaml:"name"`
	Platform    string            `yaml:"platform"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	Sections    []PromptSection   `yaml:"sections"`
	Format      string            `yaml:"format"`
	Params      map[string]string `yaml:"params,omitempty"`
	Negative    []string          `yaml:"negative,omitempty"`
}

type PromptSection struct {
	ID        string   `yaml:"id"`
	Label     string   `yaml:"label"`
	Category  string   `yaml:"category,omitempty"`
	Hint      string   `yaml:"hint,omitempty"`
	Default   string   `yaml:"default,omitempty"`
	Options   []Option `yaml:"options,omitempty"`
	Weighted  bool     `yaml:"weighted,omitempty"`
	Weight    float64  `yaml:"weight,omitempty"`
	VideoOnly bool     `yaml:"video_only,omitempty"`
}

type Option struct {
	Value  string  `yaml:"value"`
	Label  string  `yaml:"label"`
	Weight float64 `yaml:"weight,omitempty"`
}

type PromptScaffold struct {
	Seed        string             `yaml:"seed" json:"seed"`
	Platform    string             `yaml:"platform" json:"platform"`
	Model       string             `yaml:"model" json:"model"`
	Fields      map[string]string  `yaml:"fields" json:"fields"`
	Weights     map[string]float64 `yaml:"weights,omitempty" json:"weights,omitempty"`
	Negative    []string           `yaml:"negative" json:"negative"`
	FinalPrompt string             `yaml:"final_prompt" json:"final_prompt"`
	Params      map[string]string  `yaml:"params,omitempty" json:"params,omitempty"`
	Expanded    string             `yaml:"expanded,omitempty" json:"expanded,omitempty"`
	RefImage    string             `yaml:"ref_image,omitempty" json:"ref_image,omitempty"`
}

type QualityReport struct {
	ImagePath        string   `json:"image_path"`
	OverallScore     float64  `json:"overall_score"`
	ClarityScore     float64  `json:"clarity_score"`
	CompositionScore float64  `json:"composition_score"`
	ColorScore       float64  `json:"color_score"`
	PromptMatch      float64  `json:"prompt_match"`
	Issues           []string `json:"issues"`
}

type PromptVersion struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Seed      string `json:"seed"`
	Platform  string `json:"platform"`
	Prompt    string `json:"prompt"`
	Previous  string `json:"previous,omitempty"`
}
