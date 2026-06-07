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

use crate::types::{ColorScheme, LayoutTemplate, Typography};
use serde_yaml;

/// 从 YAML 字符串解析布局模板
pub fn parse_layout(yaml: &str) -> Result<LayoutTemplate, String> {
    serde_yaml::from_str::<LayoutTemplate>(yaml)
        .map_err(|e| format!("布局模板解析失败: {}", e))
}

/// 从 YAML 字符串解析配色模板
pub fn parse_color_scheme(yaml: &str) -> Result<ColorScheme, String> {
    serde_yaml::from_str::<ColorScheme>(yaml)
        .map_err(|e| format!("配色模板解析失败: {}", e))
}

/// 从 YAML 字符串解析字体模板
pub fn parse_typography(yaml: &str) -> Result<Typography, String> {
    serde_yaml::from_str::<Typography>(yaml)
        .map_err(|e| format!("字体模板解析失败: {}", e))
}

/// 批量解析三个模板
pub fn parse_templates(
    layout_yaml: &str,
    color_yaml: &str,
    font_yaml: &str,
) -> Result<(LayoutTemplate, ColorScheme, Typography), String> {
    let layout = parse_layout(layout_yaml)?;
    let color = parse_color_scheme(color_yaml)?;
    let font = parse_typography(font_yaml)?;
    Ok((layout, color, font))
}
