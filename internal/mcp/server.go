package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/veawho/via54Design/internal/media"
	"github.com/veawho/via54Design/internal/pattern"
	"github.com/veawho/via54Design/internal/quality"
	"github.com/veawho/via54Design/internal/template"
)

type Server struct {
	mcp     *server.MCPServer
	engine  *template.Engine
	baseDir string
}

func New(baseDir string) (*Server, error) {
	eng, err := template.NewEngine(baseDir)
	if err != nil {
		return nil, fmt.Errorf("engine init: %w", err)
	}
	s := &Server{mcp: server.NewMCPServer("via54Design", "0.2.0"), engine: eng, baseDir: baseDir}
	s.registerTools()
	return s, nil
}

func (s *Server) registerTools() {
	s.mcp.AddTool(mcp.NewTool("compose_template",
		mcp.WithDescription("Generate HTML by composing layout + color + font templates"),
		mcp.WithString("layout", mcp.Required(), mcp.Description("Layout template ID")),
		mcp.WithString("color", mcp.Required(), mcp.Description("Color scheme ID")),
		mcp.WithString("font", mcp.Required(), mcp.Description("Font pair ID")),
		mcp.WithString("title", mcp.Description("Page title")),
	), s.handleCompose)

	s.mcp.AddTool(mcp.NewTool("quality_check",
		mcp.WithDescription("Run quality gate on HTML file or content"),
		mcp.WithString("html", mcp.Required(), mcp.Description("HTML file path or content")),
	), s.handleQuality)

	s.mcp.AddTool(mcp.NewTool("extract_patterns",
		mcp.WithDescription("Extract design patterns (colors/fonts/layout/animations) from HTML"),
		mcp.WithString("html", mcp.Required(), mcp.Description("HTML file path or content")),
		mcp.WithString("name", mcp.Description("Project name for YAML output")),
	), s.handlePattern)

	s.mcp.AddTool(mcp.NewTool("list_templates",
		mcp.WithDescription("List all available templates"),
	), s.handleList)

	s.mcp.AddTool(mcp.NewTool("trace_image",
		mcp.WithDescription("Convert image (photo of handwriting/calligraphy) to SVG vector paths"),
		mcp.WithString("input", mcp.Required(), mcp.Description("Input image path (JPG/PNG)")),
		mcp.WithString("output", mcp.Description("Output SVG path (optional)")),
	), s.handleTrace)
}

func (s *Server) handleCompose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	result, err := s.engine.Compose(getArg[string](args, "layout"), getArg[string](args, "color"), getArg[string](args, "font"), getArgDefault(args, "title", "via54Design"))
	if err != nil { return mcp.NewToolResultError(fmt.Sprintf("Compose failed: %v", err)), nil }
	return mcp.NewToolResultText(result.HTML), nil
}

func (s *Server) handleQuality(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	htmlArg := getArg[string](req.Params.Arguments, "html")
	var htmlContent string
	if data, err := os.ReadFile(htmlArg); err == nil {
		htmlContent = string(data)
	} else {
		htmlContent = htmlArg
	}
	report := quality.CheckHTML(htmlContent)
	result := fmt.Sprintf("Quality Gate: **%s**\n\n", report.Verdict)
	result += fmt.Sprintf("- Size: %d bytes / %d CSS blocks / %d lines\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	result += fmt.Sprintf("- Issues: %d errors / %d warnings / %d info\n\n", report.Summary["error"], report.Summary["warning"], report.Summary["info"])
	for _, iss := range report.Issues {
		if iss.Severity == "info" { continue }
		icon := map[string]string{"error":"❌","warning":"⚠️"}[iss.Severity]
		result += fmt.Sprintf("%s [%s] %s\n", icon, iss.Category, iss.Message)
	}
	return mcp.NewToolResultText(result), nil
}

func (s *Server) handlePattern(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	htmlArg := getArg[string](req.Params.Arguments, "html")
	name := getArgDefault(req.Params.Arguments, "name", "unnamed")
	var htmlContent string
	if data, err := os.ReadFile(htmlArg); err == nil {
		htmlContent = string(data)
	} else {
		htmlContent = htmlArg
	}
	p, yaml := pattern.ExtractFromHTML(htmlContent, name)
	result := fmt.Sprintf("## Pattern: %s\n\n", name)
	result += fmt.Sprintf("### Colors (%d unique)\n", p.Colors.TotalUnique)
	for _, c := range p.Colors.Palette[:min(6, len(p.Colors.Palette))] {
		result += fmt.Sprintf("- %s (freq=%d)\n", c.Hex, c.Freq)
	}
	result += fmt.Sprintf("\n### Typography\n- Display: %s\n- Body: %s\n", p.Fonts.Display, p.Fonts.Body)
	result += fmt.Sprintf("\n### Layout\n- Types: %v\n- Sections: %d\n", p.Layout.Types, p.Layout.Sections)
	result += fmt.Sprintf("\n### Animation: %s\n", p.Animations.Complexity)
	result += fmt.Sprintf("\n### Metrics\n- %d lines, %d images, %d SVGs\n", p.Metrics.TotalLines, p.Metrics.ImageCount, p.Metrics.SVGCount)
	result += fmt.Sprintf("\n### YAML Template\n```yaml\n%s\n```\n", yaml)
	return mcp.NewToolResultText(result), nil
}

func (s *Server) handleList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	all := s.engine.Registry.ListAll()
	result := ""
	for cat, entries := range all {
		result += fmt.Sprintf("\n### %s\n", cat)
		for _, e := range entries {
			result += fmt.Sprintf("- `%s` — %s\n", e.ID, e.Name)
		}
	}
	return mcp.NewToolResultText(result), nil
}

func (s *Server) ServeStdio() error { return server.ServeStdio(s.mcp) }

func getArg[T any](args any, key string) T {
	if m, ok := args.(map[string]interface{}); ok {
		if v, ok := m[key]; ok { if vt, ok := v.(T); ok { return vt } }
	}
	var zero T; return zero
}
func (s *Server) handleTrace(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := getArg[string](req.Params.Arguments, "input")
	output := getArgDefault(req.Params.Arguments, "output", "")
	svgPath, err := media.TraceImage(input, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Trace failed: %v", err)), nil
	}
	if output != "" {
		os.Rename(svgPath, output)
		svgPath = output
	}
	// Read SVG content
	data, _ := os.ReadFile(svgPath)
	return mcp.NewToolResultText(fmt.Sprintf("✅ SVG generated: %s\n\n```svg\n%s\n```", svgPath, string(data[:min(len(data), 2000)]))), nil
}

func getArgDefault(args any, key string, def string) string {
	if m, ok := args.(map[string]interface{}); ok {
		if v, ok := m[key]; ok { if vt, ok := v.(string); ok && vt != "" { return vt } }
	}
	return def
}
