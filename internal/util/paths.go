// via54Design — 公共工具函数
// SPDX-License-Identifier: AGPL-3.0-only

package util

import (
	"os"
	"path/filepath"
)

// FindBaseDir 从可执行文件位置向上查找，找到包含 templates/ 的目录
// 最多向上查找 5 层，兜底返回当前工作目录
func FindBaseDir() string {
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
