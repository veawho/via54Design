// SPDX-License-Identifier: MIT OR AGPL-3.0

// via54Design — JSON 结构化导出
// 从叙事 scaffold 或 HTML 提取结构化 scene 数据
package export

import (
	"encoding/json"
	"os"
)

// SceneData 结构化场景数据
type SceneData struct {
	Title       string            `json:"title"`
	Voiceover   string            `json:"voiceover"`
	Body        string            `json:"body"`
	Mood        string            `json:"mood"`
	BeatName    string            `json:"beat_name"`
	SceneNo     int               `json:"scene_no"`
	TotalScenes int               `json:"total_scenes"`
	Duration    int               `json:"duration_seconds"`
	Timing      SceneTiming       `json:"timing"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type SceneTiming struct {
	StartSec int `json:"start_sec"`
	EndSec   int `json:"end_sec"`
}

// ExportJSON 导出场景数据为 JSON 文件
func ExportJSON(scenes []SceneData, outputPath string) error {
	data, err := json.MarshalIndent(map[string]interface{}{
		"version": "1.0",
		"total_scenes": len(scenes),
		"scenes": scenes,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0644)
}
