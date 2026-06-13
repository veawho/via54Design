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
		"version":      "1.0",
		"total_scenes": len(scenes),
		"scenes":       scenes,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0644)
}
