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
