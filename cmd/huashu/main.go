package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/veawho/via54Design/internal/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("via54Design -- 结构化设计模板引擎")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  via54Design serve            启动 MCP Server")
		fmt.Println("  via54Design generate ...      模板组合生成HTML")
		fmt.Println("  via54Design list              列出所有可用模板")
		fmt.Println("  via54Design version           版本信息")
		fmt.Println()
		fmt.Println("文档: https://github.com/veawho/via54Design")
		return
	}

	switch os.Args[1] {
	case "serve":
		cmdServe()
	case "generate":
		cmdGenerate()
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
	// 尝试可执行文件所在目录
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	// 检查目录结构是否匹配
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil {
		return dir
	}
	// 回退到当前工作目录
	wd, _ := os.Getwd()
	return wd
}

func cmdServe() {
	fmt.Println("MCP Server - 待实现")
}

func cmdGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	layout := fs.String("layout", "", "布局模板ID")
	color := fs.String("color", "", "配色模板ID")
	font := fs.String("font", "", "字体模板ID")
	title := fs.String("title", "via54Design", "页面标题")
	output := fs.String("output", "output.html", "输出路径")
	_ = fs.Parse(os.Args[2:])

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
			fmt.Printf("  %-30s %s\n", e.ID, e.Name)
		}
	}
}
