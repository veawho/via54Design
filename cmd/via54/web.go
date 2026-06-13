// via54Design — Web UI 命令
//
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/veawho/via54Design/web"
)

func cmdWeb() {
	fs := flag.NewFlagSet("web", flag.ExitOnError)
	port := fs.Int("port", 8080, "HTTP server port")
	openBrowser := fs.Bool("open", false, "Open browser automatically")
	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "flag error: %v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost:%d", *port)

	if *openBrowser {
		go func() {
			switch runtime.GOOS {
			case "darwin":
				exec.Command("open", url).Start()
			case "windows":
				exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			default:
				exec.Command("xdg-open", url).Start()
			}
		}()
	}

	fmt.Fprintf(os.Stderr, "🌐 via54Design Web UI → %s\n", url)
	fmt.Fprintf(os.Stderr, "   Templates: /api/templates\n")
	fmt.Fprintf(os.Stderr, "   Build:     POST /api/build\n")
	fmt.Fprintf(os.Stderr, "   Health:    /api/health\n")
	fmt.Fprintf(os.Stderr, "   Press Ctrl+C to stop\n")

	if err := http.ListenAndServe(addr, web.Handler(baseDir())); err != nil {
		fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
		os.Exit(1)
	}
}
