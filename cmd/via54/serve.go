// via54Design — 设计模板引擎 + 叙事引擎
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
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/veawho/via54Design/internal/mcp"
	"github.com/veawho/via54Design/internal/util"
)

func cmdServe() {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := fs.String("http", "", "HTTP 监听地址 (如 :8080)")
	metricsMode := fs.Bool("metrics", false, "Phase D: 启用 /metrics + /debug/pprof (Prometheus 文本格式 + stdlib pprof)")
	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "flag error: %v\n", err)
		os.Exit(1)
	}

	srv, err := mcp.New(util.FindBaseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "MCP 失败: %v\n", err)
		os.Exit(1)
	}

	if *httpAddr != "" {
		if *metricsMode {
			fmt.Fprintf(os.Stderr, "via54-mcp HTTP server (Phase D: metrics+pprof) on %s\n", *httpAddr)
			if err := srv.ServeWithMetrics(*httpAddr); err != nil {
				fmt.Fprintf(os.Stderr, "HTTP 错误: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "via54-mcp HTTP server on %s\n", *httpAddr)
			if err := srv.ServeHTTP(*httpAddr); err != nil {
				fmt.Fprintf(os.Stderr, "HTTP 错误: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "via54Design MCP Server (stdio)...\n")
		if err := srv.ServeStdio(); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	}
}
