// SPDX-License-Identifier: MIT OR AGPL-3.0

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
