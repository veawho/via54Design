package export

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExportPDF 通过 Playwright 将 HTML 导出为 PDF
func ExportPDF(htmlPath string, output string) (string, error) {
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		return "", fmt.Errorf("HTML 文件不存在: %s", htmlPath)
	}
	absHTML, _ := filepath.Abs(htmlPath)
	if output == "" {
		output = strings.TrimSuffix(absHTML, ".html") + ".pdf"
	}

	script := fmt.Sprintf(`
const { chromium } = require('playwright');
(async () => {
	const browser = await chromium.launch({ headless: true });
	const page = await browser.newPage();
	await page.goto('file://%s', { waitUntil: 'networkidle' });
	await page.waitForTimeout(1000);
	await page.pdf({ path: '%s', format: 'A4', printBackground: true });
	await browser.close();
	console.log('DONE');
})();
`, absHTML, output)

	tmpScript := filepath.Join(filepath.Dir(absHTML), "_pdf_temp.mjs")
	os.WriteFile(tmpScript, []byte(script), 0644)
	defer os.Remove(tmpScript)

	cmd := exec.Command("npx", "playwright", "run", tmpScript)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PDF 导出失败: %w", err)
	}
	return output, nil
}
