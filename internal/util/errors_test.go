// via54Design — 错误处理测试
// SPDX-License-Identifier: AGPL-3.0-only

package util

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrInvalidArgument", ErrInvalidArgument},
		{"ErrTemplateParse", ErrTemplateParse},
		{"ErrIO", ErrIO},
		{"ErrNotImplemented", ErrNotImplemented},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("%s should have non-empty message", tt.name)
			}
		})
	}
}

func TestWrapFunctions(t *testing.T) {
	inner := errors.New("inner error")

	tests := []struct {
		name      string
		wrapped   error
		checkFunc func(error) bool
	}{
		{"Wrap", Wrap(inner, "outer"), func(e error) bool { return e != nil }},
		{"WrapInvalid", WrapInvalid(inner, "bad arg"), IsInvalid},
		{"WrapNotFound", WrapNotFound(inner, "missing"), IsNotFound},
		{"WrapIO", WrapIO(inner, "disk fail"), IsIO},
		{"WrapTemplate", WrapTemplate(inner, "yaml fail"), IsTemplate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.checkFunc(tt.wrapped) {
				t.Errorf("%s: check failed, got %v", tt.name, tt.wrapped)
			}
			if !errors.Is(tt.wrapped, inner) {
				t.Errorf("%s: should still match inner error, got %v", tt.name, tt.wrapped)
			}
			if tt.wrapped.Error() == "" {
				t.Errorf("%s: message should not be empty", tt.name)
			}
		})
	}
}

func TestWrapNil(t *testing.T) {
	// WrapXxx(nil, ...) 现在返回带 sentinel 的 error (用于"未找到"语义)
	// 旧行为是返回 nil, 这导致 'not found' 错误被吞掉
	if Wrap(nil, "should not be nil") == nil {
		t.Error("Wrap(nil, ...) should return non-nil")
	}
	if WrapInvalid(nil, "should not be nil") == nil {
		t.Error("WrapInvalid(nil, ...) should return non-nil")
	}
	if WrapNotFound(nil, "should not be nil") == nil {
		t.Error("WrapNotFound(nil, ...) should return non-nil")
	}
	// NewNotFound / NewInvalid 应返回带 sentinel 的 error
	if NewNotFound("test") == nil {
		t.Error("NewNotFound should return non-nil")
	}
	if !IsNotFound(NewNotFound("test")) {
		t.Error("NewNotFound result should match errors.Is(ErrNotFound)")
	}
}

func TestErrorChaining(t *testing.T) {
	// 测试: errors.Is 可穿透 %w 链
	inner := ErrNotFound
	outer := fmt.Errorf("operation failed: %w", inner)
	if !errors.Is(outer, ErrNotFound) {
		t.Error("errors.Is should match wrapped sentinel")
	}

	// 多层包装
	doubleWrapped := fmt.Errorf("layer 1: %w", fmt.Errorf("layer 2: %w", inner))
	if !errors.Is(doubleWrapped, ErrNotFound) {
		t.Error("errors.Is should match 2-level wrap")
	}
}

func TestErrorMessages(t *testing.T) {
	// 中文错误消息 (项目惯例)
	inner := errors.New("inner")
	wrapped := WrapNotFound(inner, "模板不存在: hero-split-16-9")
	expected := "模板不存在: hero-split-16-9: not found: inner"
	if wrapped.Error() != expected {
		t.Errorf("got %q, want %q", wrapped.Error(), expected)
	}
}
