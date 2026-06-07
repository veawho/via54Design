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
