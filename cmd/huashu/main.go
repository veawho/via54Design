package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("via54Design — 结构化设计模板引擎")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  via54Design serve         启动 MCP Server")
		fmt.Println("  via54Design generate ...   生成 HTML")
		fmt.Println("  via54Design quality ...    质量门禁检查")
		fmt.Println("  via54Design learn ...      从参考图学习模板")
		fmt.Println()
		fmt.Println("文档: https://github.com/veawho/via54Design")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "serve":
		cmdServe()
	case "generate":
		cmdGenerate()
	case "quality":
		cmdQuality()
	case "learn":
		cmdLearn()
	case "version":
		fmt.Println("via54Design v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func cmdServe() {
	fmt.Println("MCP Server 模式 — 通过 stdio JSON-RPC 暴露工具:")
	fmt.Println("  - compose_template  (模板组合生成HTML)")
	fmt.Println("  - quality_check     (质量门禁检查)")
	fmt.Println("  - extract_patterns  (从HTML提取设计模式)")
	fmt.Println("  - learn_reference   (从参考图学习模板)")
}

func cmdGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	layout := fs.String("layout", "", "布局模板 ID")
	color := fs.String("color", "", "配色模板 ID")
	font := fs.String("font", "", "字体模板 ID")
	output := fs.String("output", "output.html", "输出路径")
	fs.Parse(os.Args[2:])

	if *layout == "" || *color == "" || *font == "" {
		fmt.Fprintln(os.Stderr, "请指定 --layout, --color, --font")
		os.Exit(1)
	}
	fmt.Printf("生成中: %s + %s + %s → %s\n", *layout, *color, *font, *output)
}

func cmdQuality() {
	fs := flag.NewFlagSet("quality", flag.ExitOnError)
	html := fs.String("html", "", "HTML 文件路径")
	fs.Parse(os.Args[2:])
	if *html == "" {
		fmt.Fprintln(os.Stderr, "请指定 --html")
		os.Exit(1)
	}
	fmt.Printf("质量检查: %s\n", *html)
}

func cmdLearn() {
	fs := flag.NewFlagSet("learn", flag.ExitOnError)
	dir := fs.String("dir", "", "参考图目录")
	name := fs.String("name", "unnamed", "模板名称")
	fs.Parse(os.Args[2:])
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "请指定 --dir")
		os.Exit(1)
	}
	fmt.Printf("学习: %s → 模板 '%s'\n", *dir, *name)
}
