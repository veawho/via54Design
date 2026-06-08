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
func Wrap(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// WrapInvalid 包装为 ErrInvalidArgument
func WrapInvalid(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrInvalidArgument, err)
}

// WrapNotFound 包装为 ErrNotFound
func WrapNotFound(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrNotFound, err)
}

// WrapIO 包装为 ErrIO
func WrapIO(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrIO, err)
}

// WrapTemplate 包装为 ErrTemplateParse
func WrapTemplate(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w: %w", fmt.Sprintf(format, args...), ErrTemplateParse, err)
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
