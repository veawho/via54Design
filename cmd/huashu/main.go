package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/veawho/via54Design/internal/mcp"
	"github.com/veawho/via54Design/internal/pattern"
	"github.com/veawho/via54Design/internal/quality"
	"github.com/veawho/via54Design/internal/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("via54Design — 结构化设计模板引擎")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  via54Design serve             启动 MCP Server (stdio)")
		fmt.Println("  via54Design generate ...       模板组合生成HTML")
		fmt.Println("  via54Design quality ...        质量门禁检查")
		fmt.Println("  via54Design pattern ...        从HTML提取设计模式")
		fmt.Println("  via54Design list               列出所有可用模板")
		fmt.Println("  via54Design version            版本信息")
		fmt.Println()
		fmt.Println("MCP 兼容: Claude Desktop / Cursor / Copilot / VS Code / Hermes")
		return
	}

	switch os.Args[1] {
	case "serve":
		cmdServe()
	case "generate":
		cmdGenerate()
	case "quality":
		cmdQuality()
	case "pattern":
		cmdPattern()
	case "list":
		cmdList()
	case "version":
		fmt.Println("via54Design v0.2.0")
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func baseDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil { return dir }
	if _, err := os.Stat(filepath.Join(dir, "template-registry.yaml")); err == nil { return dir }
	wd, _ := os.Getwd()
	return wd
}

func cmdServe() {
	srv, err := mcp.New(baseDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server init failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "via54Design MCP Server (stdio)...\n")
	if err := srv.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
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
		os.Exit(1)
	}
	eng, err := template.NewEngine(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "引擎初始化失败: %v\n", err); os.Exit(1) }
	result, err := eng.Compose(*layout, *color, *font, *title)
	if err != nil { fmt.Fprintf(os.Stderr, "生成失败: %v\n", err); os.Exit(1) }
	if err := result.SaveToFile(*output); err != nil {
		fmt.Fprintf(os.Stderr, "写入失败: %v\n", err); os.Exit(1)
	}
	fmt.Printf("✅ %s (%d bytes)\n   layout=%s color=%s font=%s\n",
		*output, len(result.HTML), result.LayoutID, result.ColorID, result.FontID)
}

func cmdQuality() {
	fs := flag.NewFlagSet("quality", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	verbose := fs.Bool("verbose", false, "显示 info 级别")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, err := os.ReadFile(*htmlFile)
	if err != nil { fmt.Fprintf(os.Stderr, "读取失败: %v\n", err); os.Exit(1) }
	report := quality.CheckHTML(string(data))
	fmt.Printf("\n=== 质量门禁: %s ===\n", report.Verdict)
	fmt.Printf("文件: %d bytes | CSS块: %d | 行: %d\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	fmt.Printf("问题: %d 错误 / %d 警告 / %d 信息\n\n", report.Summary["error"], report.Summary["warning"], report.Summary["info"])
	for _, iss := range report.Issues {
		if iss.Severity == "info" && !*verbose { continue }
		icon := map[string]string{"error":"❌","warning":"⚠️","info":"ℹ️"}[iss.Severity]
		fmt.Printf("  %s [%s] %s\n", icon, iss.Category, iss.Message)
	}
}

func cmdPattern() {
	fs := flag.NewFlagSet("pattern", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	name := fs.String("name", "unnamed", "作品名称")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, err := os.ReadFile(*htmlFile)
	if err != nil { fmt.Fprintf(os.Stderr, "读取失败: %v\n", err); os.Exit(1) }
	patterns, yaml := pattern.ExtractFromHTML(string(data), *name)
	fmt.Printf("\n=== 模式提取: %s ===\n", *name)
	fmt.Printf("  文件: %d bytes\n\n", len(data))
	fmt.Printf("  🎨 配色 (%d种):\n", patterns.Colors.TotalUnique)
	for i, c := range patterns.Colors.Palette {
		if i >= 6 { break }
		fmt.Printf("    %s (x%d)\n", c.Hex, c.Freq)
	}
	fmt.Printf("\n  🔤 字体:\n    Display: %s\n    Body: %s\n", patterns.Fonts.Display, patterns.Fonts.Body)
	fmt.Printf("\n  📐 布局: %v | %d sections\n", patterns.Layout.Types, patterns.Layout.Sections)
	fmt.Printf("  🎞️ 动画: %s (%v)\n", patterns.Animations.Complexity, patterns.Animations.Types)
	fmt.Printf("\n─── YAML 模板候选 ───\n%s\n", yaml)
}

func cmdList() {
	eng, err := template.NewEngine(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "引擎初始化失败: %v\n", err); os.Exit(1) }
	all := eng.Registry.ListAll()
	for cat, entries := range all {
		fmt.Printf("\n=== %s ===\n", cat)
		for _, e := range entries {
			fmt.Printf("  %-32s %s\n", e.ID, e.Name)
		}
	}
}
