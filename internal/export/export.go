// SPDX-License-Identifier: MIT OR AGPL-3.0

package export

// RenderResult Playwright 渲染结果
type RenderResult struct {
	VideoPath string
	PDFPath   string
	ThumbPath string
	Duration  int
	Width     int
	Height    int
}

// TTSResult 语音合成结果
type TTSResult struct {
	AudioPath string
	DurationMs int
	CharCount  int
}
