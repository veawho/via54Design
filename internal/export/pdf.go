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
	if err := checkPlaywright(); err != nil {
		return "", err
	}
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

// checkPlaywright 检测 Node.js + Playwright 是否可用
// 跨平台通用，macOS/Linux/Windows 都返回一致的错误信息
func checkPlaywright() error {
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("需要 Node.js (https://nodejs.org)\n  macOS: brew install node\n  Linux: apt install nodejs\n  Windows: choco install nodejs")
	}
	if _, err := exec.LookPath("npx"); err != nil {
		return fmt.Errorf("需要 npx (随 Node.js 安装)")
	}
	// 检查 playwright 是否安装
	cmd := exec.Command("npx", "playwright", "--version")
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("需要 Playwright (npm install -g playwright)")
	}
	return nil
}

// CheckDependencies 检查所有外部运行时依赖，返回缺失列表
// 用于 via54 quality / via54 web 启动时的前置检查
func CheckDependencies() []string {
	var missing []string
	for _, dep := range []struct {
		name    string
		command string
		args    []string
		url     string
	}{
		{"Node.js", "node", []string{"--version"}, "https://nodejs.org"},
		{"ffmpeg", "ffmpeg", []string{"-version"}, "https://ffmpeg.org"},
		{"Playwright", "npx", []string{"playwright", "--version"}, "npm install -g playwright"},
	} {
		if _, err := exec.LookPath(dep.command); err != nil {
			missing = append(missing, fmt.Sprintf("%s (%s)", dep.name, dep.url))
		} else if dep.args != nil {
			cmd := exec.Command(dep.command, dep.args...)
			if err := cmd.Run(); err != nil {
				missing = append(missing, fmt.Sprintf("%s (%s)", dep.name, dep.url))
			}
		}
	}
	return missing
}
