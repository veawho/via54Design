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

package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BGM 配乐选项
var bgmMoods = map[string]string{
	"tech":           "bgm-tech.mp3",
	"ad":             "bgm-ad.mp3",
	"educational":    "bgm-educational.mp3",
	"educational-alt":"bgm-educational-alt.mp3",
	"tutorial":       "bgm-tutorial.mp3",
	"tutorial-alt":   "bgm-tutorial-alt.mp3",
}

// AddMusic 给无声视频添加 BGM
func AddMusic(input, mood, output string) error {
	bgm, ok := bgmMoods[mood]
	if !ok {
		moods := make([]string, 0, len(bgmMoods))
		for k := range bgmMoods { moods = append(moods, k) }
		return fmt.Errorf("未知 mood: %s (可选: %v)", mood, moods)
	}
	// 查找 bgm 文件: 先找 assets/，再找 ../assets/
	bgmPath := findAsset(bgm)
	if bgmPath == "" {
		return fmt.Errorf("BGM 文件未找到: %s", bgm)
	}
	if output == "" {
		output = addSuffix(input, "-with-music")
	}
	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-i", bgmPath,
		"-filter_complex", "[1:a]adelay=500|500[a1];[0:a][a1]amix=inputs=2:duration=first[aout]",
		"-map", "0:v", "-map", "[aout]",
		"-c:v", "copy", "-c:a", "aac", "-b:a", "192k",
		"-shortest", output)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg AddMusic 失败: %w", err)
	}
	return nil
}

// ConvertFormats 从 MP4 派生产出 (60fps + GIF)
func ConvertFormats(input string) (string, string, error) {
	base := trimExt(input)
	output60 := base + "-60fps.mp4"
	gifOut := base + ".gif"

	// 60fps (帧复制，兼容性好)
	cmd60 := exec.Command("ffmpeg",
		"-i", input, "-r", "60",
		"-c:v", "libx264", "-preset", "fast", "-crf", "18",
		"-c:a", "aac", "-b:a", "128k", "-shortest", output60)
	cmd60.Stderr = os.Stderr
	if err := cmd60.Run(); err != nil {
		return "", "", fmt.Errorf("60fps 转换失败: %w", err)
	}

	// GIF (palette 优化)
	palette := base + "_palette.png"
	cmdPal := exec.Command("ffmpeg",
		"-i", input, "-vf", "fps=15,scale=960:-1:flags=lanczos,palettegen=stats_mode=diff",
		"-y", palette)
	cmdPal.Stderr = os.Stderr
	if err := cmdPal.Run(); err != nil {
		return output60, "", fmt.Errorf("palette 生成失败: %w", err)
	}

	cmdGif := exec.Command("ffmpeg",
		"-i", input, "-i", palette,
		"-lavfi", "fps=15,scale=960:-1:flags=lanczos [x]; [x][1:v] paletteuse=dither=bayer:bayer_scale=5",
		"-y", gifOut)
	cmdGif.Stderr = os.Stderr
	if err := cmdGif.Run(); err != nil {
		return output60, "", fmt.Errorf("GIF 生成失败: %w", err)
	}
	os.Remove(palette)
	return output60, gifOut, nil
}

// MixVoiceover 混合人声 + BGM
func MixVoiceover(video, voiceover, bgmMood, output string) error {
	bgm := findAsset(bgmMoods[bgmMood])
	cmd := exec.Command("ffmpeg",
		"-i", video, "-i", voiceover,
		"-i", bgm,
		"-filter_complex",
		"[1:a]adelay=200|200[a_vo];[2:a]volume=0.15[a_bgm];[a_vo][a_bgm]amix=inputs=2:duration=first[aout]",
		"-map", "0:v", "-map", "[aout]", "-c:v", "copy",
		"-c:a", "aac", "-b:a", "192k", output)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ─── helper ───
func findAsset(name string) string {
	dirs := []string{"assets", "../assets", "../../assets"}
	for _, d := range dirs {
		p := filepath.Join(d, name)
		if _, err := os.Stat(p); err == nil { return p }
	}
	return ""
}

func addSuffix(path, suffix string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + suffix + ext
}

func trimExt(path string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)]
}
