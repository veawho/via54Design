// SPDX-License-Identifier: MIT OR AGPL-3.0

use crate::types::{DesignSystem, GenerationResult};

/// 将设计系统组装为完整 HTML
pub fn assemble_html(ds: &DesignSystem, title: &str) -> GenerationResult {
    let html = format!(
        r#"<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{title}</title>
{font_imports}
<style>
/* === via54Engine (Rust) === */
{css_vars}

/* === Base === */
{base_css}

/* === Layout === */
{layout_css}
</style>
</head>
<body>
<main>
<section class="hero-split">
  <div class="hero-split__image"><!-- img --></div>
  <div class="hero-split__text">
    <p class="hero-split__eyebrow">EYEBROW</p>
    <h1 class="hero-split__headline">标题</h1>
    <p class="hero-split__body">副标题</p>
    <a class="hero-split__cta" href="#">CTA</a>
  </div>
</section>
</main>
</body>
</html>"#,
        title = title,
        font_imports = ds.font_imports,
        css_vars = ds.css_variables,
        base_css = ds.base_css,
        layout_css = ds.layout_css,
    );

    GenerationResult {
        html,
        css_variables: ds.css_variables.clone(),
        font_imports: ds.font_imports.clone(),
        base_css: ds.base_css.clone(),
    }
}

/// 完整流程：YAML → DesignSystem → HTML
pub fn compose_from_yaml(
    layout_yaml: &str,
    color_yaml: &str,
    font_yaml: &str,
    title: &str,
) -> Result<String, String> {
    let (layout, color, font) = crate::parser::parse_templates(layout_yaml, color_yaml, font_yaml)?;
    let layout_css = layout.css.unwrap_or_default();
    let ds = crate::cssgen::build_design_system(&color, &font, &layout_css);
    let result = assemble_html(&ds, title);
    Ok(result.html)
}
