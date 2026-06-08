// via54Design — 公共工具函数
// SPDX-License-Identifier: AGPL-3.0-only

package util

import (
	"os"
	"path/filepath"
)

// FindBaseDir解析模板/资源根目录,优先级:
//1. $VIA54_BASE_DIR 环境变量(显式覆盖)
//2. 从可执行文件位置向上查找 templates/ (最多5层)
//3. 当前工作目录兜底
//
// 加 env override 是为了让把 via54移到 /usr/local/bin之类的全局位置
// 后,用户能用一个 env var指向 templates/所在仓库目录(无需从源码目录运行).
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
	for i :=0; i <5; i++ {
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
