package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/veawho/via54Design/internal/mcp"
	"github.com/veawho/via54Design/internal/quality"
	"github.com/veawho/via54Design/internal/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("via54Design — 结构化设计模板引擎")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  via54Design serve            启动 MCP Server (stdio)")
		fmt.Println("  via54Design generate ...      模板组合生成HTML")
		fmt.Println("  via54Design quality ...       质量门禁检查")
		fmt.Println("  via54Design list              列出所有可用模板")
		fmt.Println("  via54Design version           版本信息")
		fmt.Println()
		fmt.Println("MCP 兼容: Claude Desktop / Cursor / Copilot / VS Code / Hermes")
		fmt.Println("文档: https://github.com/veawho/via54Design")
		return
	}

	switch os.Args[1] {
	case "serve":
		cmdServe()
	case "generate":
		cmdGenerate()
	case "quality":
		cmdQuality()
	case "list":
		cmdList()
	case "version":
		fmt.Println("via54Design v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func baseDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil {
		return dir
	}
	if _, err := os.Stat(filepath.Join(dir, "template-registry.yaml")); err == nil {
		return dir
	}
	wd, _ := os.Getwd()
	return wd
}

func cmdServe() {
	srv, err := mcp.New(baseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server 初始化失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "via54Design MCP Server 启动 (stdio)...\n")
	if err := srv.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "Server 错误: %v\n", err)
		os.Exit(1)
	}
}

func cmdGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	layout := fs.String("layout", "", "布局模板ID")
	color := fs.String("color", "", "配色模板ID")
	font := fs.String("font", "", "字体模板ID")
	title := fs.String("title", "via54Design", "页面标题")
	output := fs.String("output", "output.html", "输出路径")
	fs.Parse(os.Args[2:])

	if *layout == "" || *color == "" || *font == "" {
		fmt.Fprintln(os.Stderr, "请指定: --layout, --color, --font")
		fmt.Fprintln(os.Stderr, "可用模板: via54Design list")
		os.Exit(1)
	}

	eng, err := template.NewEngine(baseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "引擎初始化失败: %v\n", err)
		os.Exit(1)
	}

	result, err := eng.Compose(*layout, *color, *font, *title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}

	if err := result.SaveToFile(*output); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 已生成: %s (%d bytes)\n", *output, len(result.HTML))
	fmt.Printf("   布局: %s  |  配色: %s  |  字体: %s\n",
		result.LayoutID, result.ColorID, result.FontID)
}

func cmdQuality() {
	fs := flag.NewFlagSet("quality", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	verbose := fs.Bool("verbose", false, "显示 info 级别")
	fs.Parse(os.Args[2:])

	if *htmlFile == "" {
		fmt.Fprintln(os.Stderr, "请指定 --html")
		os.Exit(1)
	}

	data, err := os.ReadFile(*htmlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
		os.Exit(1)
	}

	report := quality.CheckHTML(string(data))
	fmt.Printf("\n=== 质量门禁: %s ===\n", report.Verdict)
	fmt.Printf("文件: %d bytes | CSS块: %d | 行: %d\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	fmt.Printf("问题: %d 错误 / %d 警告 / %d 信息\n\n",
		report.Summary["error"], report.Summary["warning"], report.Summary["info"])

	for _, iss := range report.Issues {
		if iss.Severity == "info" && !*verbose {
			continue
		}
		icon := map[string]string{"error": "❌", "warning": "⚠️", "info": "ℹ️"}[iss.Severity]
		fmt.Printf("  %s [%s] %s\n", icon, iss.Category, iss.Message)
	}
}

func cmdList() {
	eng, err := template.NewEngine(baseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "引擎初始化失败: %v\n", err)
		os.Exit(1)
	}
	all := eng.Registry.ListAll()
	for cat, entries := range all {
		fmt.Printf("\n=== %s ===\n", cat)
		for _, e := range entries {
			fmt.Printf("  %-32s %s\n", e.ID, e.Name)
		}
	}
}
