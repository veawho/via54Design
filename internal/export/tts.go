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

package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"time"
)

// callDouBaoTTS 通过豆包 TTS API 合成语音
func callDouBaoTTS(text, apiKey, voice string) ([]byte, error) {
	if voice == "" {
		voice = "BV000_streaming_voice_comfort"
	}
	body := map[string]interface{}{
		"model": "doubao-tts",
		"input": map[string]string{"text": text},
		"voice": map[string]string{"voice_type": voice},
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.volcengine.com/tts", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TTS 请求失败: %w", err)
	}
	defer resp.Body.Close()

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取音频失败: %w", err)
	}
	return audio, nil
}

// Synthesize 合成 TTS 语音文件
func Synthesize(text, outputPath, apiKey, voice string) (*TTSResult, error) {
	if apiKey == "" {
		apiKey = os.Getenv("DOUBAO_TTS_API_KEY")
	}
	audio, err := callDouBaoTTS(text, apiKey, voice)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(outputPath, audio, 0644); err != nil {
		return nil, fmt.Errorf("写入音频失败: %w", err)
	}
	return &TTSResult{
		AudioPath: outputPath,
		CharCount: len(text),
	}, nil
}
