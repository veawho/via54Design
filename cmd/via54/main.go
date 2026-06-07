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
	"github.com/veawho/via54Design/internal/narrate"
	"github.com/veawho/via54Design/internal/pattern"
	"github.com/veawho/via54Design/internal/quality"
	"github.com/veawho/via54Design/internal/template"
	"github.com/veawho/via54Design/internal/wasm"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "serve":      cmdServe()
	case "generate":   cmdGenerate()
	case "narrate":    cmdNarrate()
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
	fmt.Println("  narrate ...        叙事脚手架 — 一句话到故事大纲+剧本+分镜")
	fmt.Println("                     --seed \"一句话\" --model three-act [--duration 30]")
	fmt.Println("                     [--format markdown|json] [--output file]")
	fmt.Println("  quality ...        质量门禁检查")
	fmt.Println("  pattern ...        从HTML提取设计模式")
	fmt.Println("  list               列出所有可用模板")
	fmt.Println("  media ...          媒体管线 (ffmpeg + trace)")
	fmt.Println("    add-music     input.mp4 --mood=tech")
	fmt.Println("    convert       input.mp4")
	fmt.Println("    fetch         --query parrot --out ./img")
	fmt.Println("    trace         --input photo.jpg --output logo.svg")
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
	if _, err := os.Stat(filepath.Join(dir, "templates/registry.yaml")); err == nil { return dir }
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
	letteringSVG := fs.String("lettering-svg", "", "SVG lettering 文件路径 (手写/书法文字)")
	fromNarrative := fs.String("from-narrative", "", "叙事脚手架 JSON 路径 (via54 narrate --format json 的输出)")
	fs.Parse(os.Args[2:])

	bd := baseDir()
	eng, err := template.NewEngine(bd)
	if err != nil { fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1) }

	// ── 叙事驱动生成 ──
	if *fromNarrative != "" {
		generateFromNarrative(eng, *fromNarrative, *layout, *color, *font, *output)
		return
	}

	if (*layout == "" || *color == "" || *font == "") && *letteringSVG == "" {
		fmt.Fprintln(os.Stderr, "请指定: --layout, --color, --font")
		fmt.Fprintln(os.Stderr, "或者: --from-narrative scaffold.json (叙事驱动模式)")
		os.Exit(1)
	}
	var result *template.GenerationResult
	if *letteringSVG != "" {
		data, _ := os.ReadFile(*letteringSVG)
		if *layout == "" { *layout = "hero-split-left-image" }
		if *color == "" { *color = "ink-wash" }
		if *font == "" { *font = "serif-sans-editorial" }
		result, err = eng.ComposeWithSVG(*layout, *color, *font, *title, string(data))
	} else {
		result, err = eng.Compose(*layout, *color, *font, *title)
	}
	if err != nil { fmt.Fprintf(os.Stderr, "生成失败: %v\n", err); os.Exit(1) }
	result.SaveToFile(*output)
	fmt.Printf("✅ %s (%d bytes)\n   layout=%s color=%s font=%s\n", *output, len(result.HTML), result.LayoutID, result.ColorID, result.FontID)
}

// generateFromNarrative 从叙事脚手架 JSON 生成多场景 HTML 动画
func generateFromNarrative(eng *template.Engine, narrativePath, layoutOverride, colorOverride, fontOverride, output string) {
	data, err := os.ReadFile(narrativePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取叙事文件失败: %v\n", err)
		os.Exit(1)
	}

	var scaffold narrate.NarrativeScaffold
	if err := yaml.Unmarshal(data, &scaffold); err != nil {
		fmt.Fprintf(os.Stderr, "解析叙事文件失败: %v\n", err)
		os.Exit(1)
	}

	// 叙事映射到视觉模板
	moodColor := map[string]string{
		"mysterious":   "dark-terminal-blue",
		"aspirational": "cosmic-retro",
		"confident":    "neo-brutalist-vibrant",
		"urgent":       "neon-dark",
		"calm":         "moon-white",
		"curious":      "ink-wash",
		"excited":      "candy-duolingo",
		"inspiring":    "warm-editorial-cream",
		"informative":  "bauhaus-primary",
		"focused":      "crimson-elegance",
		"warm":         "daylily-warmth",
		"practical":    "earth-terracotta",
		"insightful":   "ultramarine-deep",
		"hopeful":      "pine-spring",
		"frustrated":   "dark-terminal-blue",
		"tense":        "neon-dark",
		"triumphant":   "candy-duolingo",
	}

	// 默认视觉模板
	lID := layoutOverride
	if lID == "" { lID = "hero-split-left-image" }
	cID := colorOverride
	if cID == "" { cID = "ink-wash" }
	fID := fontOverride
	if fID == "" { fID = "ming-hei-editorial" }

	// 为每个 beat 生成场景 HTML，拼合成多场景页面
	var scenesHTML strings.Builder
	scenesHTML.WriteString(`<!DOCTYPE html><html lang="zh-CN"><head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: system-ui; overflow-x: hidden; }
.scene { min-height: 100vh; display: flex; flex-direction: column; justify-content: center;
  align-items: center; padding: 4rem 2rem; position: relative; transition: opacity 0.5s; }
.scene .time { position: absolute; top: 1rem; right: 1rem; font-size: 0.8rem; opacity: 0.5; }
.scene .mood-tag { position: absolute; top: 1rem; left: 1rem; font-size: 0.7rem;
  text-transform: uppercase; letter-spacing: 0.1em; padding: 0.2rem 0.6rem;
  border-radius: 4px; background: rgba(255,255,255,0.1); }
.scene h2 { font-size: 1.8rem; margin-bottom: 1rem; text-align: center; max-width: 40rem; }
.scene .voiceover { font-size: 2.4rem; font-weight: 700; text-align: center;
  margin: 1.5rem 0; max-width: 36rem; line-height: 1.4; }
.scene .sub { font-size: 1rem; opacity: 0.7; text-align: center;
  margin-top: 1rem; max-width: 30rem; }
.scene .act-label { font-size: 0.9rem; opacity: 0.6; margin-bottom: 0.5rem;
  text-transform: uppercase; letter-spacing: 0.15em; }
.scene-nav { position: fixed; bottom: 2rem; left: 50%; transform: translateX(-50%);
  display: flex; gap: 0.5rem; z-index: 100; }
.scene-nav button { padding: 0.5rem 1.2rem; border: 1px solid currentColor;
  background: transparent; cursor: pointer; border-radius: 4px; font-size: 0.9rem;
  transition: all 0.2s; }
.scene-nav button:hover { background: rgba(255,255,255,0.15); }
.scene-nav .progress { display: flex; align-items: center; gap: 0.3rem;
  font-size: 0.8rem; opacity: 0.5; }
</style></head><body>
`)

	for i, beat := range scaffold.Beats {
		// 根据情绪选配色
		sceneColor := cID
		if mc, ok := moodColor[beat.Mood]; ok && colorOverride == "" {
			sceneColor = mc
		}

		// 每个场景用独立配色生成
		sceneResult, err := eng.Compose(lID, sceneColor, fID,
			fmt.Sprintf("%s — %s", scaffold.Seed, beat.Act))
		if err != nil {
			// fallback: plain scene
			sceneResult = &template.GenerationResult{
				HTML: fmt.Sprintf("<div class=\"scene\"><h2>%s</h2><p class=\"voiceover\">%s</p></div>",
					beat.Act, beat.Voiceover),
			}
		}

		// 提取 body 内容
		html := sceneResult.HTML
		bodyStart := strings.Index(html, "<body")
		bodyEnd := strings.Index(html, "</body>")
		bodyContent := ""
		if bodyStart > 0 && bodyEnd > 0 {
			contentStart := strings.Index(html[bodyStart:], ">")
			if contentStart > 0 {
				bodyContent = html[bodyStart+contentStart+1 : bodyEnd]
			}
		}
		if bodyContent == "" {
			bodyContent = fmt.Sprintf("<h2>%s</h2><p>%s</p>", beat.Act, beat.Voiceover)
		}

		scenesHTML.WriteString(fmt.Sprintf(`<div class="scene" id="scene-%d" data-beat="%s" data-mood="%s" data-duration="%d">
  <div class="act-label">%s</div>
  <div class="time">%ds → %ds</div>
  <div class="mood-tag">%s</div>
  %s
  <div class="voiceover">%s</div>
</div>
`, i+1, beat.Act, beat.Mood, beat.Duration, beat.Act,
			beat.StartTime, beat.StartTime+beat.Duration, beat.Mood,
			bodyContent, beat.Voiceover))
	}

	// 导航 + JS 自动播放
	scenesHTML.WriteString(`<div class="scene-nav">`)
	for i := range scaffold.Beats {
		scenesHTML.WriteString(fmt.Sprintf(`<button onclick="location.href='#scene-%d'" data-idx="%d">%d</button>`, i+1, i, i+1))
	}
	scenesHTML.WriteString(`<div class="progress"><span id="progress-text">1/`)
	scenesHTML.WriteString(fmt.Sprintf("%d", len(scaffold.Beats)))
	scenesHTML.WriteString(`</span></div></div>
<script>
let currentScene = 0;
const scenes = document.querySelectorAll('.scene');
scenes.forEach((s, i) => { s.style.display = i === 0 ? 'flex' : 'none'; });
document.querySelectorAll('.scene-nav button').forEach(btn => {
  btn.addEventListener('click', () => {
    const idx = parseInt(btn.dataset.idx);
    scenes.forEach((s, i) => { s.style.display = i === idx ? 'flex' : 'none'; });
    currentScene = idx;
  });
});
</script></body></html>`)

	finalHTML := scenesHTML.String()
	os.WriteFile(output, []byte(finalHTML), 0644)
	fmt.Printf("✅ 叙事驱动动画: %s (%d 场景, %d 字节)\n   via54 narrate → via54 generate --from-narrative\n",
		output, len(scaffold.Beats), len(finalHTML))
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
		// 叙事学分组显示
		if cat == "narratology" {
			fmt.Println("  ── guides ──")
			for _, e := range entries {
				if e.Category == "narratology/guide" || e.Category == "" {
					fmt.Printf("    %-32s %s\n", e.ID, e.Name)
				}
			}
			fmt.Println("  ── models ──")
			for _, e := range entries {
				if e.Category == "narratology/model" {
					extra := ""
					if len(e.Tags) > 0 {
						extra = fmt.Sprintf("  [%s]", strings.Join(e.Tags, ", "))
					}
					fmt.Printf("    %-32s %s%s\n", e.ID, e.Name, extra)
				}
			}
			continue
		}
		for _, e := range entries {
			fmt.Printf("  %-32s %s\n", e.ID, e.Name)
		}
	}
}

// ─── Narrate (叙事引擎) ───
func cmdNarrate() {
	fs := flag.NewFlagSet("narrate", flag.ExitOnError)
	seed := fs.String("seed", "", "一句话种子 (必填)")
	model := fs.String("model", "three-act", "叙事模型ID")
	duration := fs.Int("duration", 30, "目标视频时长(秒)")
	output := fs.String("output", "", "输出文件路径 (默认: stdout)")
	format := fs.String("format", "markdown", "输出格式: markdown / json")
	listModels := fs.Bool("list", false, "列出所有叙事模型")
	fs.Parse(os.Args[2:])

	bd := baseDir()

	if *listModels {
		out, err := narrate.ListModels(bd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "失败: %v\n", err); os.Exit(1)
		}
		fmt.Print(out)
		return
	}
	if *seed == "" {
		fmt.Fprintln(os.Stderr, "请指定 --seed \"一句话描述\"")
		fmt.Fprintln(os.Stderr, "或使用 --list 查看可用叙事模型")
		os.Exit(1)
	}

	scaffold, err := narrate.GenerateScaffold(*seed, *model, *duration, bd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "叙事生成失败: %v\n", err)
		os.Exit(1)
	}

	if *format == "json" {
		js, err := scaffold.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON 序列化失败: %v\n", err)
			os.Exit(1)
		}
		if *output != "" {
			os.WriteFile(*output, []byte(js), 0644)
			fmt.Printf("✅ 叙事脚手架 JSON 已保存: %s\n", *output)
		} else {
			fmt.Print(js)
		}
		return
	}

	md, err := scaffold.RenderMarkdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "渲染失败: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		os.WriteFile(*output, []byte(md), 0644)
		fmt.Printf("✅ 叙事脚手架已保存: %s\n", *output)
	} else {
		fmt.Print(md)
	}
}

// ─── Media (Shell+Python 迁移) ───
func cmdMedia() {
	if len(os.Args) < 3 { fmt.Fprintln(os.Stderr, "用法: via54 media [add-music|convert|fetch|trace]"); os.Exit(1) }
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

	case "trace":
		fs := flag.NewFlagSet("trace", flag.ExitOnError)
		input := fs.String("input", "", "输入图片路径 (JPG/PNG)")
		output := fs.String("output", "", "输出 SVG 路径 (可选)")
		fs.Parse(os.Args[3:])
		if *input == "" { fmt.Fprintln(os.Stderr, "请指定 --input"); os.Exit(1) }
		svgPath, err := media.TraceImage(*input, nil)
		if err != nil { fmt.Fprintf(os.Stderr, "❌ %v\n", err); os.Exit(1) }
		if *output != "" {
			os.Rename(svgPath, *output)
			svgPath = *output
		}
		fmt.Printf("✅ SVG: %s\n", svgPath)
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
