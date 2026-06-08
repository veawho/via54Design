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

package template

// ─── 布局模板类型 (v2 — 16:9 多端适配) ───

type LayoutTemplate struct {
	ID       string      `yaml:"id"`
	Name     interface{} `yaml:"name"`
	Version  string      `yaml:"version"`
	Category string      `yaml:"category"`
	Tags     []string    `yaml:"tags"`
	When     LayoutWhen  `yaml:"when"`

	Viewport   ViewportConfig  `yaml:"viewport"`
	Structure  LayoutStructure `yaml:"structure"`
	Spacing    SpacingScale    `yaml:"spacing,omitempty"`
	Responsive []BreakpointDef `yaml:"responsive"`
	Elements   []Element       `yaml:"elements"`
	CSS        string          `yaml:"css"`
}

// ─── 视口配置 ───

type ViewportConfig struct {
	Baseline  string `yaml:"baseline"`            // "16:9"
	MinHeight string `yaml:"min_height"`          // "100dvh"
	MaxWidth  string `yaml:"max_width,omitempty"` // 可选最大宽度（TV 用）

	// 演示模式:
	//   开启: 锁定 16:9, 生成 PPT/视频级输出
	//   关闭: 自由网页布局
	//   CLI 通过 --presentation 控制, PPT/PPTX/演示文稿 关键词自动激活
	PresentationMode bool `yaml:"presentation_mode,omitempty"`
	Presentation     struct {
		RatioLock string `yaml:"ratio_lock,omitempty"` // "16:9"
		FitMode   string `yaml:"fit_mode,omitempty"`   // "contain" / "cover" / "fill"
		BgColor   string `yaml:"bg_color,omitempty"`   // 演示背景色（超出16:9区域）
	} `yaml:"presentation,omitempty"`
}

// ─── 布局结构 ───

type LayoutStructure struct {
	Type  string `yaml:"type"`            // grid-2col / grid-3col / flex / bento
	Ratio string `yaml:"ratio,omitempty"` // "5,7" / "1fr 1fr"
	Cols  int    `yaml:"cols,omitempty"`  // 列数 (bento)
	Rows  int    `yaml:"rows,omitempty"`  // 行数 (bento)
	Gap   string `yaml:"gap,omitempty"`   // "24px"
}

// ─── 间距黄金比例 ───

type SpacingScale struct {
	Base     int               `yaml:"base"`               // 4
	Ratio    float64           `yaml:"ratio"`              // 1.618
	Steps    []int             `yaml:"steps,omitempty"`    // 自动生成
	Semantic map[string]string `yaml:"semantic,omitempty"` // {section: "step-7", card: "step-4"}
}

// ─── 响应式断点 ───
// 参考: extra-strength-responsive-grids (⭐254) 流体栅格 + CSS Container Queries

type BreakpointDef struct {
	Name     string `yaml:"name"`                // tv / desktop / tablet / phone
	MinWidth int    `yaml:"min_width"`           // 1920 / 1280 / 768 / 0
	MaxWidth int    `yaml:"max_width,omitempty"` // 可选上限

	// 栅格覆盖
	Columns string `yaml:"columns,omitempty"` // "5,7" / "1fr"
	Stack   bool   `yaml:"stack,omitempty"`   // 是否堆叠

	// 字体缩放
	FontScale float64 `yaml:"font_scale"` // 1.2 / 1.0 / 0.9 / 0.75

	// 安全区域 (TV overscan / 手机刘海)
	SafeArea []int `yaml:"safe_area,omitempty"` // [top, right, bottom, left] px

	// 间距缩放 (相对 base)
	SpacingScale float64 `yaml:"spacing_scale,omitempty"` // 1.0 / 0.8 / 0.6

	// 隐藏/显示元素 (按 role)
	HideRoles  []string `yaml:"hide_roles,omitempty"`
	StackOrder []string `yaml:"stack_order,omitempty"` // [text, image] 堆叠顺序

	// 布局特定设置
	FullBleed *bool  `yaml:"full_bleed,omitempty"`
	RatioLock string `yaml:"ratio_lock,omitempty"` // "16:9" / "4:3" / "auto"
}

// ─── 元素 ───

type Element struct {
	Role     string    `yaml:"role"`
	Position string    `yaml:"position,omitempty"`
	Behavior string    `yaml:"behavior,omitempty"`
	Tag      string    `yaml:"tag,omitempty"`
	Style    string    `yaml:"style,omitempty"`
	FontSize string    `yaml:"font_size,omitempty"`
	ZIndex   int       `yaml:"z_index,omitempty"`
	Padding  []int     `yaml:"padding,omitempty"` // [top, right, bottom, left]
	MaxWidth string    `yaml:"max_width,omitempty"`
	Children []Element `yaml:"children,omitempty"`

	// 元素级响应式 (key = breakpoint name: tv/desktop/tablet/phone)
	Responsive map[string]ElementResponsive `yaml:"responsive,omitempty"`
}

type ElementResponsive struct {
	Hide     bool   `yaml:"hide,omitempty"`      // 隐藏此元素
	Order    int    `yaml:"order,omitempty"`     // flex/grid order
	FontSize string `yaml:"font_size,omitempty"` // 覆盖字号
	Padding  []int  `yaml:"padding,omitempty"`   // 覆盖内边距
	Columns  int    `yaml:"columns,omitempty"`   // 跨越列数
}

// ─── 以下保留不动 ───

type LayoutWhen struct {
	ContentHas  []string `yaml:"content_has"`
	SuitableFor []string `yaml:"suitable_for"`
}

type PaletteItem struct {
	Role         string `yaml:"role"`
	Hex          string `yaml:"hex"`
	NameZh       string `yaml:"name_zh,omitempty"`
	CulturalNote string `yaml:"cultural_note,omitempty"`
}

type ColorSchemeTemplate struct {
	ID           string            `yaml:"id"`
	Name         interface{}       `yaml:"name"`
	Version      string            `yaml:"version"`
	Tags         []string          `yaml:"tags"`
	Source       string            `yaml:"source,omitempty"`
	Colors       map[string]string `yaml:"colors,omitempty"`
	Palette      []PaletteItem     `yaml:"palette,omitempty"`
	When         ColorWhen         `yaml:"when"`
	CSSVariables string            `yaml:"css_variables,omitempty"`
	Mood         []string          `yaml:"mood,omitempty"`
	Season       string            `yaml:"season,omitempty"`
}

type ColorWhen struct {
	BrandTone   []string `yaml:"brand_tone"`
	Audience    []string `yaml:"audience"`
	SuitableFor []string `yaml:"suitable_for"`
}

type TypographyTemplate struct {
	ID          string            `yaml:"id"`
	Name        interface{}       `yaml:"name"`
	Version     string            `yaml:"version"`
	Tags        []string          `yaml:"tags"`
	Fonts       map[string]string `yaml:"fonts"`
	Weights     map[string]string `yaml:"weights,omitempty"`
	Sizes       map[string]string `yaml:"sizes"`
	GoogleFonts []string          `yaml:"google_fonts,omitempty"`
	Mood        []string          `yaml:"mood,omitempty"`
	When        TypographyWhen    `yaml:"when"`
}

type TypographyWhen struct {
	BrandTone   []string `yaml:"brand_tone"`
	SuitableFor []string `yaml:"suitable_for"`
}

type TemplateRegistry struct {
	Layouts      []RegistryEntry `yaml:"layouts"`
	ColorSchemes []RegistryEntry `yaml:"color_schemes"`
	Typography   []RegistryEntry `yaml:"typography"`
	Narratology  []RegistryEntry `yaml:"narratology"`
	Combinations []Combination   `yaml:"combinations"`
}

type RegistryEntry struct {
	ID       string      `yaml:"id"`
	Name     interface{} `yaml:"name"`
	Version  string      `yaml:"version"`
	Category string      `yaml:"category"`
	Tags     []string    `yaml:"tags"`
	File     string      `yaml:"file"`
}

type Combination struct {
	Name     interface{} `yaml:"name"`
	Layout   string      `yaml:"layout"`
	Color    string      `yaml:"color"`
	Font     string      `yaml:"font"`
	Suitable []string    `yaml:"suitable"`
}

type GenerationResult struct {
	HTML             string
	LayoutID         string
	ColorID          string
	FontID           string
	Title            string
	FontImports      string
	CSSVars          string
	BaseCSS          string
	LayoutCSS        string
	LetteringSVG     string
	PresentationMode bool // 是否启用 16:9 演示锁定
}
