// SPDX-License-Identifier: MIT OR AGPL-3.0

use crate::types::{ColorScheme, DesignSystem, Typography};
use std::collections::HashMap;

/// 从配色+字体生成完整 CSS 变量块
pub fn generate_css_variables(color: &ColorScheme, font: &Typography) -> String {
    let mut css = String::from(":root {\n");

    // 优先使用模板中定义的 css_variables
    if let Some(ref vars) = color.css_variables {
        css.push_str(vars);
    } else {
        for (role, hex) in &color.colors {
            css.push_str(&format!("  --{}: {};\n", role, hex));
        }
    }

    // 字号变量
    if let Some(ref sizes) = font.sizes {
        for (name, size) in sizes {
            css.push_str(&format!("  --size-{}: {};\n", name, size));
        }
    }

    css.push_str("}");
    css
}

/// Google Fonts 常用字体集合
fn is_google_font(family: &str) -> bool {
    matches!(family,
        "Inter" | "Geist" | "Geist Mono" | "JetBrains Mono" | "Fira Code" |
        "Fraunces" | "Newsreader" | "Playfair Display" | "Source Serif 4" |
        "Lora" | "EB Garamond" | "Cormorant Garamond" | "Space Grotesk" |
        "Space Mono" | "Manrope" | "Poppins" | "Nunito" | "Baloo 2" |
        "Archivo" | "Archivo Black" | "Anton" | "IBM Plex Sans" |
        "IBM Plex Mono" | "DM Serif Display" | "Noto Serif SC"
    )
}

/// 生成 Google Fonts 导入
pub fn generate_font_imports(font: &Typography) -> String {
    let mut seen = Vec::new();
    for family in font.fonts.values() {
        let primary = family.split(',')
            .next()
            .unwrap_or(family)
            .trim()
            .trim_matches(&['\'', '\"'] as &[_]);
        if seen.contains(&primary.to_string()) || !is_google_font(primary) {
            continue;
        }
        seen.push(primary.to_string());
        let weight = if primary == "Noto Serif SC" { "400;600;700;900" } else { "400;500;600;700" };
        return format!(
            "<link rel="preconnect" href="https://fonts.googleapis.com">\n             <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>\n             <link href="https://fonts.googleapis.com/css2?family={}:wght@{}&display=swap" rel="stylesheet">",
            primary.replace(' ', "+"), weight
        );
    }
    String::new()
}

/// 生成基础 CSS
pub fn generate_base_css(font: &Typography) -> String {
    let body = font.fonts.get("body").map_or("'Inter', sans-serif".to_string(), |s| s.clone());
    let display = font.fonts.get("display").map_or(body.clone(), |s| s.clone());
    let mono = font.fonts.get("mono").map_or("'JetBrains Mono', monospace".to_string(), |s| s.clone());

    format!(
        r#"* {{ box-sizing: border-box; margin: 0; padding: 0; }}
html {{ scroll-behavior: smooth; }}
body {{
  font-family: {};
  line-height: 1.7;
  color: var(--text-primary, #1A1A1A);
  background: var(--background, #FFFFFF);
  -webkit-font-smoothing: antialiased;
}}
h1,h2,h3,h4 {{ font-family: {}; line-height: 1.1; }}
code,pre {{ font-family: {}; }}
a {{ color: var(--accent, inherit); }}
img {{ max-width: 100%; height: auto; }}
.container {{ max-width: 1200px; margin: 0 auto; padding: 0 40px; }}
"#, body, display, mono)
}

/// 组装完整设计系统
pub fn build_design_system(
    color: &ColorScheme,
    font: &Typography,
    layout_css: &str,
) -> DesignSystem {
    DesignSystem {
        css_variables: generate_css_variables(color, font),
        font_imports: generate_font_imports(font),
        base_css: generate_base_css(font),
        layout_css: layout_css.to_string(),
    }
}
