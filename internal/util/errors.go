// via54Design — 公共错误定义
// SPDX-License-Identifier: AGPL-3.0-only

package util

import (
	"errors"
	"fmt"
)

// Sentinel errors — 使用 errors.Is 检查
var (
	// ErrNotFound 资源未找到 (模板/文件/平台)
	ErrNotFound = errors.New("not found")

	// ErrInvalidArgument 无效参数 (布局名/颜色名/字体名)
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrTemplateParse YAML 模板解析失败
	ErrTemplateParse = errors.New("template parse error")

	// ErrIO I/O 错误 (读/写文件)
	ErrIO = errors.New("I/O error")

	// ErrNotImplemented 功能未实现
	ErrNotImplemented = errors.New("not implemented")
)

// Wrap 包装错误并附带上下文 (使用 %w 保留错误链)
//   - 如果 err != nil: 返回带 sentinel + 上下文的 error
//   - 如果 err == nil: 仍返回带 sentinel 的 error (不返回 nil!)
func Wrap(err error, format string, args ...any) error {
	if err == nil {
		return fmt.Errorf("%s", fmt.Sprintf(format, args...))
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// WrapInvalid 包装为 ErrInvalidArgument
func WrapInvalid(err error, format string, args ...any) error {
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrInvalidArgument, err)
}

// WrapNotFound 包装为 ErrNotFound
// 即使 err == nil, 也会返回带 ErrNotFound 的 error (用于"未找到"语义)
func WrapNotFound(err error, format string, args ...any) error {
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrNotFound, err)
}

// WrapIO 包装为 ErrIO
func WrapIO(err error, format string, args ...any) error {
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrIO, err)
}

// WrapTemplate 包装为 ErrTemplateParse
func WrapTemplate(err error, format string, args ...any) error {
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrTemplateParse, err)
}

// NewNotFound 创建带 ErrNotFound sentinel 的新错误
// 用于"明确标识"未找到的场景 (e.g., resource not found by ID)
func NewNotFound(format string, args ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), ErrNotFound)
}

// NewInvalid 创建带 ErrInvalidArgument sentinel 的新错误
func NewInvalid(format string, args ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), ErrInvalidArgument)
}

// IsNotFound 检查是否为 NotFound 错误
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalid 检查是否为 InvalidArgument 错误
func IsInvalid(err error) bool {
	return errors.Is(err, ErrInvalidArgument)
}

// IsTemplate 检查是否为 TemplateParse 错误
func IsTemplate(err error) bool {
	return errors.Is(err, ErrTemplateParse)
}

// IsIO 检查是否为 I/O 错误
func IsIO(err error) bool {
	return errors.Is(err, ErrIO)
}
