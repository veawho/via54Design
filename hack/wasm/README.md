# via54Engine — Rust 核心引擎

## 定位

Rust 版核心引擎，编译为 WASM 后可在以下场景运行：

| 场景 | 调用方式 | 价值 |
|------|----------|------|
| 浏览器内预览 | JS 直调 wasm | 不上传服务器，隐私安全 |
| Go 主程序 | wazero 运行时 | 零 CGo 依赖 |
| 边缘计算 | Cloudflare Workers | 全球分发 |
| 独立 CLI | 原生编译 | 调试/测试 |

## 架构

```
src/
├── lib.rs       WASM 入口 + CLI
├── types.rs     模板类型定义 (serde)
├── parser.rs    YAML→内部结构
├── cssgen.rs    CSS 生成
└── html.rs      HTML 组装
```

## 编译

```bash
# WASM (浏览器/Go/Node)
bash build.sh

# 原生 CLI (测试用)
cargo build --release
```

## WASM API

```js
// 完整组合
const html = compose(layoutYaml, colorYaml, fontYaml, "标题");

// 仅CSS变量
const vars = cssVariables(colorYaml, fontYaml);

// 仅字体导入
const fonts = fontImports(fontYaml);
```

## 和 Go 引擎的关系

- **Rust** — 纯计算：模板解析 + CSS/HTML 生成
- **Go** — IO 编排：文件读取、MCP、CLI、子进程管理
- **Shell** — 媒体管线：ffmpeg、Playwright

Rust 层提供 `compose()` 这个纯函数 — 输入 YAML 字符串，输出 HTML 字符串。无副作用，适合 WASM。
