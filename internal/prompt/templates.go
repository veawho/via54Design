// SPDX-License-Identifier: AGPL-3.0-only
package prompt

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// NegativeBank 平台特定负面词库
var NegativeBank = map[string][]string{
	"midjourney": {
		"blurry, low quality, distorted, deformed, ugly, bad anatomy",
		"extra limbs, missing fingers, bad hands",
		"watermark, signature, text, logo",
		"oversaturated, underexposed, noise, grainy",
	},
	"kling": {
		"jittery, shaky, unstable motion, flickering",
		"distorted face, bad expression, unnatural movement",
		"watermark, text overlay, logo",
	},
	"jimeng": {
		"模糊, 变形, 畸形, 崩坏, 低质量",
		"多余肢体, 手指错误, 水印, 文字",
		"过度饱和, 噪点, 压缩痕迹",
	},
	"gemini": {
		"blurry, distorted, low quality, bad anatomy",
		"watermark, text, signature",
		"overexposed, noise, artifacts",
	},
}

func loadTemplate(platform, baseDir string) *PromptTemplate {
	candidates := []string{
		filepath.Join(baseDir, "templates", "prompts", platform+".yaml"),
		filepath.Join(baseDir, "templates", "prompts", platform+".yml"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var t PromptTemplate
		if err := yaml.Unmarshal(data, &t); err == nil && t.ID != "" {
			return &t
		}
	}
	return nil
}
