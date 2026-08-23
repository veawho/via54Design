// via54Design — MCP Server 独立入口 (Independent entry point for MCP Server)
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

// [CN] via54-mcp — Model Context Protocol 服务端。
// [EN] via54-mcp — Model Context Protocol Server.
// 独立二进制，供 Claude Desktop / Cursor / Copilot / VS Code / Hermes 调用
//
// 用法 (Usage):
//   via54-mcp                    # 启动 stdio 模式 (Start stdio mode)
//   via54-mcp --http :8080      # 启动 HTTP 模式 (开发用) (Start HTTP mode for dev)
//
// MCP 配置 (MCP Configuration):
//   Claude Desktop:
//     { "mcpServers": { "via54Design": { "command": "via54-mcp" } } }

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/veawho/via54Design/internal/mcp"
	"github.com/veawho/via54Design/internal/util"
)

func main() {
	httpAddr := flag.String("http", "", "HTTP 监听地址 (如 :8080)")
	flag.Parse()

	srv, err := mcp.New(util.FindBaseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "MCP 初始化失败: %v\n", err)
		os.Exit(1)
	}

	if *httpAddr != "" {
		fmt.Fprintf(os.Stderr, "via54-mcp HTTP server on %s\n", *httpAddr)
		if err := srv.ServeHTTP(*httpAddr); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP 错误: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "via54-mcp stdio (MCP protocol)...\n")
		if err := srv.ServeStdio(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP 错误: %v\n", err)
			os.Exit(1)
		}
	}
}
