// via54Design — 公共工具函数
// SPDX-License-Identifier: AGPL-3.0-only

package util

import (
	"os"
	"path/filepath"
)

// FindBaseDir 解析 via54 的基础目录
//
// 优先级:
//  1. VIA54_BASE_DIR 环境变量 (Mac/Linux 安装到 /usr/local/bin 时必须)
//  2. 可执行文件位置向上查找, 找到含 templates/ 的目录 (最多 5 层)
//  3. 当前工作目录 (兜底)
func FindBaseDir() string {
	if env := os.Getenv("VIA54_BASE_DIR"); env != "" {
		return env
	}
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	dir := filepath.Dir(exe)
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	wd, _ := os.Getwd()
	return wd
}
