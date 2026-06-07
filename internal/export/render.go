package export

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RenderVideo 使用 Playwright 将 HTML 录制成 MP4
// 实际调用 npx playwright (保留 node 依赖)
// 输入: htmlPath 路径, duration 秒
func RenderVideo(htmlPath string, duration int, width, height int) (*RenderResult, error) {
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("HTML 文件不存在: %s", htmlPath)
	}
	absHTML, _ := filepath.Abs(htmlPath)
	outDir := filepath.Dir(absHTML)
	baseName := strings.TrimSuffix(filepath.Base(absHTML), ".html")

	// Playwright Node.js 脚本 —— 通过 stdlib os/exec 调 npx
	// 使用内联 Node.js 脚本，避免文件依赖
	script := fmt.Sprintf(`
const { chromium } = require('playwright');
(async () => {
	const browser = await chromium.launch({ headless: true });
	const page = await browser.newPage({
		viewport: { width: %d, height: %d }
	});
	await page.goto('file://%s', { waitUntil: 'networkidle' });
	await page.waitForTimeout(2000);
	const ctx = await browser.newContext({
		recordVideo: { dir: '%s', size: { width: %d, height: %d } }
	});
	const p2 = await ctx.newPage();
	await p2.goto('file://%s', { waitUntil: 'networkidle' });
	await p2.waitForTimeout(3000);
	await p2.waitForTimeout(%d * 1000);
	await p2.close();
	await ctx.close();
	await browser.close();
	console.log('DONE');
})();
`, width, height, absHTML, outDir, width, height, absHTML, duration)

	// 写入临时脚本
	tmpScript := filepath.Join(outDir, "_render_temp.mjs")
	os.WriteFile(tmpScript, []byte(script), 0644)
	defer os.Remove(tmpScript)

	cmd := exec.Command("npx", "playwright", "run", tmpScript)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Playwright 渲染失败: %w (需要: npm install playwright + npx playwright install chromium)", err)
	}

	// 查找生成的 webm 文件
	videoFile := filepath.Join(outDir, baseName+".mp4")

	return &RenderResult{
		VideoPath: videoFile,
		Duration:  duration,
		Width:     width,
		Height:    height,
	}, nil
}
