package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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
	s := &Server{
		mcp:     server.NewMCPServer("via54Design", "0.1.0"),
		engine:  eng,
		baseDir: baseDir,
	}
	s.registerTools()
	return s, nil
}

func (s *Server) registerTools() {
	composeTool := mcp.NewTool("compose_template",
		mcp.WithDescription("Generate HTML by composing layout + color + font templates"),
		mcp.WithString("layout", mcp.Required(), mcp.Description("Layout template ID")),
		mcp.WithString("color", mcp.Required(), mcp.Description("Color scheme ID")),
		mcp.WithString("font", mcp.Required(), mcp.Description("Font pair ID")),
		mcp.WithString("title", mcp.Description("Page title")),
	)
	s.mcp.AddTool(composeTool, s.handleCompose)

	qualityTool := mcp.NewTool("quality_check",
		mcp.WithDescription("Run quality gate checks on an HTML file"),
		mcp.WithString("html", mcp.Required(), mcp.Description("HTML file path or content")),
	)
	s.mcp.AddTool(qualityTool, s.handleQuality)

	listTool := mcp.NewTool("list_templates",
		mcp.WithDescription("List all available templates"),
	)
	s.mcp.AddTool(listTool, s.handleList)
}

func (s *Server) handleCompose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	layout := getArg[string](args, "layout")
	color := getArg[string](args, "color")
	font := getArg[string](args, "font")
	title := getArgDefault(args, "title", "via54Design")

	result, err := s.engine.Compose(layout, color, font, title)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Compose failed: %v", err)), nil
	}
	return mcp.NewToolResultText(result.HTML), nil
}

func (s *Server) handleQuality(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	htmlArg := getArg[string](req.Params.Arguments, "html")

	var htmlContent string
	data, err := os.ReadFile(htmlArg)
	if err == nil {
		htmlContent = string(data)
	} else {
		htmlContent = htmlArg
	}

	report := quality.CheckHTML(htmlContent)
	result := fmt.Sprintf("Quality Gate: **%s**\n\n", report.Verdict)
	result += fmt.Sprintf("- Size: %d bytes, %d CSS blocks, %d lines\n", report.HTMLSize, report.CSSBlocks, report.TotalLines)
	result += fmt.Sprintf("- Issues: %d errors / %d warnings / %d info\n\n", report.Summary["error"], report.Summary["warning"], report.Summary["info"])

	for _, iss := range report.Issues {
		if iss.Severity == "info" {
			continue
		}
		icon := map[string]string{"error": "❌", "warning": "⚠️"}[iss.Severity]
		result += fmt.Sprintf("%s [%s] %s\n", icon, iss.Category, iss.Message)
	}
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

func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.mcp)
}

// Helper: extract arg from map with type assertion
func getArg[T any](args any, key string) T {
	if m, ok := args.(map[string]interface{}); ok {
		if v, ok := m[key]; ok {
			if vt, ok := v.(T); ok {
				return vt
			}
		}
	}
	var zero T
	return zero
}

func getArgDefault(args any, key string, def string) string {
	if m, ok := args.(map[string]interface{}); ok {
		if v, ok := m[key]; ok {
			if vt, ok := v.(string); ok && vt != "" {
				return vt
			}
		}
	}
	return def
}
