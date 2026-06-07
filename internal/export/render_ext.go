// via54Design — 多格式视频导出
// 扩展 render.go, 支持 webm/hevc/frames/apng
package export

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RenderFormats 支持的视频格式
var RenderFormats = map[string]struct {
	Codec  string
	Ext    string
	Muxer  string
	Desc   string
}{
	"mp4":   {"libx264", ".mp4", "mp4", "H.264 通用格式"},
	"webm":  {"libvpx-vp9", ".webm", "webm", "VP9 开源格式"},
	"hevc":  {"libx265", ".mp4", "mp4", "H.265 高压缩比"},
	"frames":{"", "", "image2", "PNG 序列帧 (输出目录)"},
	"apng":  {"apng", ".png", "apng", "APNG 动图"},
}

// RenderVideoExt 扩展版: 支持多格式
func RenderVideoExt(htmlPath string, duration int, width, height int, format string) (*RenderResult, error) {
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("HTML 文件不存在: %s", htmlPath)
	}

	fmtInfo, ok := RenderFormats[format]
	if !ok {
		return nil, fmt.Errorf("不支持的格式: %s (可选: mp4/webm/hevc/frames/apng)", format)
	}

	absHTML, _ := filepath.Abs(htmlPath)
	outDir := filepath.Dir(absHTML)
	baseName := strings.TrimSuffix(filepath.Base(absHTML), ".html")

	// 先用 Playwright 录 raw video
	rawWebm := filepath.Join(outDir, baseName+"_raw.webm")
	script := fmt.Sprintf(`const { chromium } = require('playwright');
(async () => {
  const browser = await chromium.launch({ headless: true });
  const ctx = await browser.newContext({
    recordVideo: { dir: '%s', size: { width: %d, height: %d } }
  });
  const page = await ctx.newPage();
  await page.goto('file://%s', { waitUntil: 'networkidle' });
  await page.waitForTimeout(2000);
  await page.waitForTimeout(%d * 1000);
  await page.close();
  await ctx.close();
  await browser.close();
  // 重命名 raw 文件
  const fs = require('fs');
  const files = fs.readdirSync('%s').filter(f => f.endsWith('.webm'));
  if (files.length) fs.renameSync('%s/'+files[0], '%s');
})();`, outDir, width, height, absHTML, duration+2, outDir, outDir, rawWebm)

	tmpScript := filepath.Join(outDir, "_render_temp.mjs")
	os.WriteFile(tmpScript, []byte(script), 0644)
	defer os.Remove(tmpScript)

	cmd := exec.Command("npx", "playwright", "run", tmpScript)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Playwright 录制失败: %w", err)
	}

	outputFile := filepath.Join(outDir, baseName+fmtInfo.Ext)

	if format == "mp4" || format == "hevc" {
		// ffmpeg 转码
		ffmpegCmd := exec.Command("ffmpeg", "-y",
			"-i", rawWebm,
			"-c:v", fmtInfo.Codec,
			"-pix_fmt", "yuv420p",
			outputFile)
		ffmpegCmd.Stderr = os.Stderr
		if err := ffmpegCmd.Run(); err != nil {
			// fallback: 直接返回 raw
			os.Rename(rawWebm, outputFile)
		} else {
			os.Remove(rawWebm)
		}
	} else if format == "webm" {
		os.Rename(rawWebm, outputFile)
	} else if format == "frames" {
		os.MkdirAll(outputFile, 0755)
		ffmpegCmd := exec.Command("ffmpeg", "-y",
			"-i", rawWebm,
			filepath.Join(outputFile, "frame-%04d.png"))
		ffmpegCmd.Stderr = os.Stderr
		ffmpegCmd.Run()
		os.Remove(rawWebm)
	} else if format == "apng" {
		ffmpegCmd := exec.Command("ffmpeg", "-y",
			"-i", rawWebm,
			"-plays", "0",
			outputFile)
		ffmpegCmd.Stderr = os.Stderr
		if err := ffmpegCmd.Run(); err != nil {
			os.Rename(rawWebm, outputFile)
		} else {
			os.Remove(rawWebm)
		}
	}

	return &RenderResult{
		VideoPath: outputFile,
		Duration:  duration,
		Width:     width,
		Height:    height,
	}, nil
}
