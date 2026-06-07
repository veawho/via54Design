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
	if voice == "" { voice = "BV000_streaming_voice_comfort" }
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
	if err != nil { return nil, fmt.Errorf("TTS 请求失败: %w", err) }
	defer resp.Body.Close()

	audio, err := io.ReadAll(resp.Body)
	if err != nil { return nil, fmt.Errorf("读取音频失败: %w", err) }
	return audio, nil
}

// Synthesize 合成 TTS 语音文件
func Synthesize(text, outputPath, apiKey, voice string) (*TTSResult, error) {
	if apiKey == "" {
		apiKey = os.Getenv("DOUBAO_TTS_API_KEY")
	}
	audio, err := callDouBaoTTS(text, apiKey, voice)
	if err != nil { return nil, err }
	if err := os.WriteFile(outputPath, audio, 0644); err != nil {
		return nil, fmt.Errorf("写入音频失败: %w", err)
	}
	return &TTSResult{
		AudioPath:  outputPath,
		CharCount:  len(text),
	}, nil
}
