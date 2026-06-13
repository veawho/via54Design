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

pub mod types;
pub mod parser;
pub mod cssgen;
pub mod html;

// 非 WASM 编译时提供 CLI 入口
#[cfg(not(target_arch = "wasm32"))]
pub fn run_cli() {
    println!("via54Engine Rust Core v{}", env!("CARGO_PKG_VERSION"));
    println!();
    println!("编译为 WASM 后供 Go/browser 调用:");
    println!("  cargo build --target wasm32-unknown-unknown --release");
    println!();
    println!("暴露的 WASM 函数:");
    println!("  compose(layout_yaml, color_yaml, font_yaml, title) -> String");
    println!("  css_variables(color_yaml, font_yaml) -> String");
    println!("  font_imports(font_yaml) -> String");
}

// WASM 导出
#[cfg(target_arch = "wasm32")]
mod wasm_bridge {
    use wasm_bindgen::prelude::*;

    /// 完整模板组合：YAML → HTML
    /// 从 JS/Go WASM 运行时调用
    #[wasm_bindgen]
    pub fn compose(
        layout_yaml: &str,
        color_yaml: &str,
        font_yaml: &str,
        title: &str,
    ) -> Result<String, JsValue> {
        crate::html::compose_from_yaml(layout_yaml, color_yaml, font_yaml, title)
            .map_err(|e| JsValue::from_str(&e))
    }

    /// 仅生成 CSS 变量块（用于渐进式增强）
    #[wasm_bindgen]
    pub fn css_variables(color_yaml: &str, font_yaml: &str) -> Result<String, JsValue> {
        let color = crate::parser::parse_color_scheme(color_yaml)
            .map_err(|e| JsValue::from_str(&e))?;
        let font = crate::parser::parse_typography(font_yaml)
            .map_err(|e| JsValue::from_str(&e))?;
        Ok(crate::cssgen::generate_css_variables(&color, &font))
    }

    /// 仅生成字体导入链接
    #[wasm_bindgen]
    pub fn font_imports(font_yaml: &str) -> Result<String, JsValue> {
        let font = crate::parser::parse_typography(font_yaml)
            .map_err(|e| JsValue::from_str(&e))?;
        Ok(crate::cssgen::generate_font_imports(&font))
    }
}
