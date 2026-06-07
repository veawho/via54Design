// SPDX-License-Identifier: MIT OR AGPL-3.0

package main

import (
	"fmt"
	"github.com/veawho/via54Design/internal/mcp"
	"os"
)

func cmdServe() {
	srv, err := mcp.New(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "MCP 失败: %v\n", err); os.Exit(1) }
	fmt.Fprintf(os.Stderr, "via54Design MCP Server (stdio)...\n")
	if err := srv.ServeStdio(); err != nil { fmt.Fprintf(os.Stderr, "错误: %v\n", err); os.Exit(1) }
}


