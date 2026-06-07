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

package media

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// VTracer 版本
const vtracerVersion = "0.6.4"

// vtracerReleaseURL 下载地址模板
var vtracerReleaseURL = fmt.Sprintf(
	"https://github.com/visioncortex/vtracer/releases/download/v%s/vtracer-%%s-%%s.%%s",
	vtracerVersion)

// downloadURL 返回当前平台的 VTracer 下载 URL
func downloadURL() (string, string, error) {
	arch := runtime.GOARCH
	osName := runtime.GOOS

	var osPart, ext string
	switch osName {
	case "darwin":
		osPart = "apple-darwin"
		ext = "tar.gz"
	case "linux":
		osPart = "unknown-linux-musl"
		ext = "tar.gz"
	case "windows":
		osPart = "pc-windows-msvc"
		ext = "zip"
	default:
		return "", "", fmt.Errorf("不支持的系统: %s", osName)
	}

	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	default:
		return "", "", fmt.Errorf("不支持架构: %s", arch)
	}

	url := fmt.Sprintf(vtracerReleaseURL, arch, osPart, ext)
	binaryName := "vtracer"
	if osName == "windows" {
		binaryName = "vtracer.exe"
	}
	return url, binaryName, nil
}

// vtracerPath 查找或下载 VTracer 二进制
func vtracerPath() (string, error) {
	// 1. 检查 PATH
	if p, err := exec.LookPath("vtracer"); err == nil {
		return p, nil
	}

	// 2. 检查项目本地目录
	localDir := filepath.Join(".via54", "bin")
	os.MkdirAll(localDir, 0755)
	binaryName := "vtracer"
	if runtime.GOOS == "windows" {
		binaryName = "vtracer.exe"
	}
	localPath := filepath.Join(localDir, binaryName)
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// 3. 下载
	fmt.Fprintf(os.Stderr, "📥 下载 VTracer (v%s) ~1MB...\n", vtracerVersion)
	url, binName, err := downloadURL()
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载 VTracer 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("下载 VTracer 失败: HTTP %d", resp.StatusCode)
	}

	if strings.HasSuffix(url, ".tar.gz") {
		// tar.gz
		gzr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("解压失败: %w", err)
		}
		defer gzr.Close()
		tr := tar.NewReader(gzr)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("解压失败: %w", err)
			}
			if header.Name == binName || strings.HasSuffix(header.Name, "/"+binName) {
				out, err := os.Create(localPath)
				if err != nil {
					return "", err
				}
				defer out.Close()
				io.Copy(out, tr)
				os.Chmod(localPath, 0755)
				fmt.Fprintf(os.Stderr, "✅ VTracer 已下载到 %s\n", localPath)
				return localPath, nil
			}
		}
		return "", fmt.Errorf("压缩包中未找到 %s", binName)
	} else {
		// zip (Windows)
		out, err := os.Create(localPath)
		if err != nil {
			return "", err
		}
		defer out.Close()
		io.Copy(out, resp.Body)
		os.Chmod(localPath, 0755)
		fmt.Fprintf(os.Stderr, "✅ VTracer 已下载到 %s\n", localPath)
		return localPath, nil
	}
}

// TraceImageOpts VTracer 选项
type TraceImageOpts struct {
	Input       string // 输入图片路径
	Output      string // 输出 SVG 路径
	Hierarchical int   // 层次细节 (0-1, 默认0)
	Threshold    int   // 二值化阈值 (0-255, 默认180)
	PreDenoise   int   // 预降噪 (0-20, 默认1)
	ColorPrecision int // 颜色精度 (默认6)
	CornerThreshold int // 角点阈值 (默认98)
}

// DefaultTraceOpts 书法矢量化推荐参数
var DefaultTraceOpts = TraceImageOpts{
	Hierarchical:    0,
	Threshold:       180,
	PreDenoise:      2,
	ColorPrecision:  6,
	CornerThreshold: 98,
}

// TraceImage 将图片矢量为 SVG
func TraceImage(input string, opts *TraceImageOpts) (string, error) {
	if opts == nil {
		opts = &DefaultTraceOpts
	}

	// 检查输入文件
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return "", fmt.Errorf("输入文件不存在: %s", input)
	}

	// 确定输出路径
	output := opts.Output
	if output == "" {
		ext := filepath.Ext(input)
		output = input[:len(input)-len(ext)] + ".svg"
	}

	// 找到 VTracer 二进制
	vtracer, err := vtracerPath()
	if err != nil {
		return "", fmt.Errorf("VTracer 不可用: %w", err)
	}

	// 构建命令行参数
	args := []string{
		"--input", input,
		"--output", output,
		"--hierarchical", fmt.Sprintf("%d", opts.Hierarchical),
		"--threshold", fmt.Sprintf("%d", opts.Threshold),
		"--pre_denoise", fmt.Sprintf("%d", opts.PreDenoise),
		"--color_precision", fmt.Sprintf("%d", opts.ColorPrecision),
		"--corner_threshold", fmt.Sprintf("%d", opts.CornerThreshold),
	}

	cmd := exec.Command(vtracer, args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("VTracer 处理失败: %w", err)
	}

	// 验证输出
	if _, err := os.Stat(output); os.IsNotExist(err) {
		return "", fmt.Errorf("VTracer 未生成输出文件: %s", output)
	}

	return output, nil
}

// TraceResult SVG 追踪结果
type TraceResult struct {
	SVGPath string
	Width   int
	Height  int
}
