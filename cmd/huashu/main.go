package main
// via54Design — 衍生自 huashu-design (MIT) by alchaincyf


import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/veawho/via54Design/internal/export"
	"github.com/veawho/via54Design/internal/media"
	"github.com/veawho/via54Design/internal/mcp"
	"github.com/veawho/via54Design/internal/pattern"
	"github.com/veawho/via54Design/internal/quality"
	"github.com/veawho/via54Design/internal/template"
	"github.com/veawho/via54Design/internal/wasm"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "serve":      cmdServe()
	case "generate":   cmdGenerate()
	case "quality":    cmdQuality()
	case "pattern":    cmdPattern()
	case "list":       cmdList()
	case "media":      cmdMedia()
	case "export":     cmdExport()
	case "version":    cmdVersion()
	default:           help()
	}
}

func help() {
	fmt.Println("via54 — 设计模板引擎 (via54Design)")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  serve             启动 MCP Server (stdio)")
	fmt.Println("  generate ...       模板组合生成HTML")
	fmt.Println("  quality ...        质量门禁检查")
	fmt.Println("  pattern ...        从HTML提取设计模式")
	fmt.Println("  list               列出所有可用模板")
	fmt.Println("  media ...          媒体管线 (ffmpeg)")
	fmt.Println("    add-music     input.mp4 --mood=tech")
	fmt.Println("    convert       input.mp4")
	fmt.Println("    fetch         --query parrot --out ./img")
	fmt.Println("  export ...         导出 (Playwright/TTS)")
	fmt.Println("    render        input.html --duration 30")
	fmt.Println("    pdf           input.html")
	fmt.Println("    tts           --text '你好' --out audio.mp3")
	fmt.Println("  version           版本信息")
	fmt.Println()
	fmt.Println("MCP: Claude Desktop / Cursor / Copilot / VS Code / Hermes")
	fmt.Println("Docs: https://github.com/veawho/via54Design")
}

func baseDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil { return dir }
	if _, err := os.Stat(filepath.Join(dir, "template-registry.yaml")); err == nil { return dir }
	wd, _ := os.Getwd()
	return wd
}

func cmdVersion() {
	fmt.Println("via54Design v0.3.0")
	fmt.Printf("Go: %s %s/%s\n", strings.TrimPrefix(runtime.Version(), "go"), runtime.GOOS, runtime.GOARCH)
	we := wasm.NewEngine(baseDir())
	if we.Available() {
		fmt.Println("WASM: ✅ via54-engine loaded")
	} else {
		fmt.Println("WASM: ❌ (cd internal/wasm && bash build.sh)")
	}
	fmt.Println("Lang: Go + Rust + 1 JS (html2pptx)")
}

// ─── Template ───
func cmdGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	layout := fs.String("layout", "", "布局模板ID")
	color := fs.String("color", "", "配色模板ID")
	font := fs.String("font", "", "字体模板ID")
	title := fs.String("title", "via54Design", "页面标题")
	output := fs.String("output", "output.html", "输出路径")
	fs.Parse(os.Args[2:])
	if *layout == "" || *color == "" || *font == "" {
		fmt.Fprintln(os.Stderr, "请指定: --layout, --color, --font"); os.Exit(1)
	}
	eng, err := template.NewEngine(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1) }
	result, err := eng.Compose(*layout, *color, *font, *title)
	if err != nil { fmt.Fprintf(os.Stderr, "生成失败: %v\n", err); os.Exit(1) }
	result.SaveToFile(*output)
	fmt.Printf("✅ %s (%d bytes)\n   layout=%s color=%s font=%s\n", *output, len(result.HTML), result.LayoutID, result.ColorID, result.FontID)
}

func cmdQuality() {
	fs := flag.NewFlagSet("quality", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	verbose := fs.Bool("verbose", false, "显示 info")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, _ := os.ReadFile(*htmlFile)
	report := quality.CheckHTML(string(data))
	fmt.Printf("\n=== 质量门禁: %s ===\n", report.Verdict)
	fmt.Printf("文件: %d bytes | CSS块: %d | 行: %d\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	fmt.Printf("问题: %d 错误 / %d 警告 / %d 信息\n\n", report.Summary["error"], report.Summary["warning"], report.Summary["info"])
	for _, iss := range report.Issues {
		if iss.Severity == "info" && !*verbose { continue }
		icon := map[string]string{"error": "❌", "warning": "⚠️", "info": "ℹ️"}[iss.Severity]
		fmt.Printf("  %s [%s] %s\n", icon, iss.Category, iss.Message)
	}
}

func cmdPattern() {
	fs := flag.NewFlagSet("pattern", flag.ExitOnError)
	htmlFile := fs.String("html", "", "HTML 文件路径")
	name := fs.String("name", "unnamed", "作品名称")
	fs.Parse(os.Args[2:])
	if *htmlFile == "" { fmt.Fprintln(os.Stderr, "请指定 --html"); os.Exit(1) }
	data, _ := os.ReadFile(*htmlFile)
	p, yaml := pattern.ExtractFromHTML(string(data), *name)
	fmt.Printf("\n=== 模式提取: %s ===\n", *name)
	fmt.Printf("  🎨 %d colors | 📐 %v | 🎞️ %s\n", p.Colors.TotalUnique, p.Layout.Types, p.Animations.Complexity)
	fmt.Printf("  🔤 Display: %s  Body: %s\n", p.Fonts.Display, p.Fonts.Body)
	fmt.Printf("\n─── YAML 候选 ───\n%s\n", yaml)
}

func cmdList() {
	eng, err := template.NewEngine(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1) }
	for cat, entries := range eng.Registry.ListAll() {
		fmt.Printf("\n=== %s ===\n", cat)
		for _, e := range entries { fmt.Printf("  %-32s %s\n", e.ID, e.Name) }
	}
}

// ─── Media (Shell+Python 迁移) ───
func cmdMedia() {
	if len(os.Args) < 3 { fmt.Fprintln(os.Stderr, "用法: via54 media [add-music|convert|fetch]"); os.Exit(1) }
	switch os.Args[2] {
	case "add-music":
		fs := flag.NewFlagSet("add-music", flag.ExitOnError)
		mood := fs.String("mood", "tech", "配乐 mood")
		output := fs.String("output", "", "输出路径")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.mp4"); os.Exit(1) }
		if err := media.AddMusic(fs.Arg(0), *mood, *output); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1)
		}
		fmt.Println("✅ BGM 添加完成")

	case "convert":
		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.mp4"); os.Exit(1) }
		mp4, gif, err := media.ConvertFormats(fs.Arg(0))
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ 60fps: %s\n", mp4)
		if gif != "" { fmt.Printf("✅ GIF:   %s\n", gif) }

	case "fetch":
		fs := flag.NewFlagSet("fetch", flag.ExitOnError)
		query := fs.String("query", "", "关键词 (逗号分隔)")
		out := fs.String("out", "./img", "输出目录")
		count := fs.Int("count", 2, "每关键词张数")
		fs.Parse(os.Args[3:])
		if *query == "" { fmt.Fprintln(os.Stderr, "请指定 --query"); os.Exit(1) }
		queries := strings.Split(*query, ",")
		results, err := media.FetchImages(queries, *out, *count)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ 下载 %d 张到 %s\n", len(results), *out)

	default:
		fmt.Fprintf(os.Stderr, "未知 media 命令: %s (支持: add-music, convert, fetch)\n", os.Args[2])
		os.Exit(1)
	}
}

// ─── Export (JS 迁移) ───
func cmdExport() {
	if len(os.Args) < 3 { fmt.Fprintln(os.Stderr, "用法: via54 export [render|pdf|tts]"); os.Exit(1) }
	switch os.Args[2] {
	case "render":
		fs := flag.NewFlagSet("render", flag.ExitOnError)
		duration := fs.Int("duration", 10, "时长(秒)")
		width := fs.Int("width", 1920, "宽")
		height := fs.Int("height", 1080, "高")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.html"); os.Exit(1) }
		r, err := export.RenderVideo(fs.Arg(0), *duration, *width, *height)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ 视频: %s\n", r.VideoPath)

	case "pdf":
		fs := flag.NewFlagSet("pdf", flag.ExitOnError)
		output := fs.String("output", "", "输出路径")
		fs.Parse(os.Args[3:])
		if fs.NArg() < 1 { fmt.Fprintln(os.Stderr, "请指定 input.html"); os.Exit(1) }
		p, err := export.ExportPDF(fs.Arg(0), *output)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ PDF: %s\n", p)

	case "tts":
		fs := flag.NewFlagSet("tts", flag.ExitOnError)
		text := fs.String("text", "", "文本")
		output := fs.String("output", "output.mp3", "输出路径")
		voice := fs.String("voice", "", "音色")
		fs.Parse(os.Args[3:])
		if *text == "" { fmt.Fprintln(os.Stderr, "请指定 --text"); os.Exit(1) }
		r, err := export.Synthesize(*text, *output, "", *voice)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		fmt.Printf("✅ TTS: %s (%d chars)\n", r.AudioPath, r.CharCount)

	default:
		fmt.Fprintf(os.Stderr, "未知 export 命令: %s (支持: render, pdf, tts)\n", os.Args[2])
		os.Exit(1)
	}
}

func cmdServe() {
	srv, err := mcp.New(baseDir())
	if err != nil { fmt.Fprintf(os.Stderr, "MCP 失败: %v\n", err); os.Exit(1) }
	fmt.Fprintf(os.Stderr, "via54Design MCP Server (stdio)...\n")
	if err := srv.ServeStdio(); err != nil { fmt.Fprintf(os.Stderr, "错误: %v\n", err); os.Exit(1) }
}
