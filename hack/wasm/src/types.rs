// via54Design — WASM 模板引擎 (Rust)
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

use serde::{Deserialize, Serialize};

/// 布局模板
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LayoutTemplate {
    pub id: String,
    pub name: String,
    pub css: Option<String>,
    pub elements: Option<Vec<Element>>,
}

/// 布局元素
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Element {
    pub role: Option<String>,
    pub position: Option<String>,
    pub behavior: Option<String>,
    pub tag: Option<String>,
    pub style: Option<String>,
    pub font_size: Option<String>,
    pub children: Option<Vec<Element>>,
}

/// 配色模板
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ColorScheme {
    pub id: String,
    pub name: String,
    pub colors: std::collections::HashMap<String, String>,
    pub css_variables: Option<String>,
}

/// 字体模板
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Typography {
    pub id: String,
    pub name: String,
    pub fonts: std::collections::HashMap<String, String>,
    pub sizes: Option<std::collections::HashMap<String, String>>,
}

/// 组合后的设计系统
#[derive(Debug, Clone, Serialize)]
pub struct DesignSystem {
    pub css_variables: String,
    pub font_imports: String,
    pub base_css: String,
    pub layout_css: String,
}

/// 生成的 HTML 结果
#[derive(Debug, Clone, Serialize)]
pub struct GenerationResult {
    pub html: String,
    pub css_variables: String,
    pub font_imports: String,
    pub base_css: String,
}
