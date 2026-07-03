// via54Design — 设计模板引擎 + 叙事引擎
//
// via54 reimagine — 截图 → HTML (LLM Vision API 驱动)
// 借鉴: abi/screenshot-to-code (73k⭐ MIT) — 但只用 Go + 复用 via54Design 的 layout/color/font 三元组
//
// 用法:
//   via54 reimagine --screenshot shot.png --provider openai --output out.html
//   via54 reimagine --screenshot shot.png --provider gemini --model gemini-2.0-flash --layout bento-grid-2x2
//
// 支持 provider: openai (gpt-4o / gpt-4-vision) / anthropic (claude-sonnet-4) / gemini (gemini-2.0-flash)
// 不引入 Python/FastAPI/前端, 保持 via54 单二进制零依赖路线.

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func cmdReimagine() {
	fs := flag.NewFlagSet("reimagine", flag.ExitOnError)
	screenshot := fs.String("screenshot", "", "截图文件路径 (PNG/JPG/WebP)")
	provider := fs.String("provider", "openai", "LLM provider: openai/anthropic/gemini")
	model := fs.String("model", "", "覆盖默认模型 (默认按 provider 选最新 vision 模型)")
	layoutHint := fs.String("layout-hint", "", "可选: 指定 via54 layout ID 强引导 LLM (如 bento-grid-2x2)")
	colorHint := fs.String("color-hint", "", "可选: 指定 via54 color ID")
	fontHint := fs.String("font-hint", "", "可选: 指定 via54 font ID")
	output := fs.String("output", "reimagined.html", "输出 HTML 路径")
	apiKey := fs.String("api-key", os.Getenv("VIA54_LLM_API_KEY"), "API key (默认读 env VIA54_LLM_API_KEY)")
	timeout := fs.Int("timeout", 90, "HTTP timeout 秒")
	fs.Parse(os.Args[2:])

	if *screenshot == "" {
		fmt.Fprintln(os.Stderr, "必须指定 --screenshot <path>")
		os.Exit(1)
	}
	if _, err := os.Stat(*screenshot); err != nil {
		fmt.Fprintf(os.Stderr, "截图文件不存在: %s\n", *screenshot)
		os.Exit(1)
	}
	if *apiKey == "" {
		// 尝试按 provider 各自的 env 兜底
		switch *provider {
		case "openai":
			*apiKey = os.Getenv("OPENAI_API_KEY")
		case "anthropic":
			*apiKey = os.Getenv("ANTHROPIC_API_KEY")
		case "gemini":
			*apiKey = os.Getenv("GEMINI_API_KEY")
		}
	}
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "需要 API key: --api-key 或 env (OPENAI_API_KEY / ANTHROPIC_API_KEY / GEMINI_API_KEY / VIA54_LLM_API_KEY)")
		os.Exit(1)
	}

	// 1. 读图, base64
	imgData, mimeType, err := readImageAsBase64(*screenshot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取截图失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📸 截图: %s (%d bytes, %s)\n", *screenshot, len(imgData), mimeType)

	// 2. 构造 prompt (中文, 引导 LLM 复刻视觉结构)
	prompt := buildReimaginePrompt(*layoutHint, *colorHint, *fontHint)

	// 3. 按 provider 调 LLM
	fmt.Printf("🤖 调用 %s vision API ...\n", *provider)
	html, err := callVisionAPI(*provider, *model, *apiKey, mimeType, imgData, prompt, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LLM 调用失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ LLM 返回 %d bytes HTML\n", len(html))

	// 4. 写到文件
	if err := os.WriteFile(*output, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("💾 %s (%d bytes)\n", *output, len(html))
	fmt.Println("提示: 这是 LLM 原始输出, 可能需要微调. 可用 'via54 generate --layout X --color Y --font Z' 改风格.")
}

func readImageAsBase64(path string) (string, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return "", "", err
	}
	ext := strings.ToLower(filepath.Ext(path))
	mime := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".webp":
		mime = "image/webp"
	case ".gif":
		mime = "image/gif"
	}
	return base64.StdEncoding.EncodeToString(data), mime, nil
}

func buildReimaginePrompt(layoutHint, colorHint, fontHint string) string {
	hintParts := []string{}
	if layoutHint != "" {
		hintParts = append(hintParts, fmt.Sprintf("- 整体布局参考 via54Design 的 `%s` layout 骨架", layoutHint))
	}
	if colorHint != "" {
		hintParts = append(hintParts, fmt.Sprintf("- 配色采用 via54Design 的 `%s` 配色变量 (--bg/--text-primary/--accent/--border)", colorHint))
	}
	if fontHint != "" {
		hintParts = append(hintParts, fmt.Sprintf("- 字体用 via54Design 的 `%s` 字体方案", fontHint))
	}
	hintBlock := ""
	if len(hintParts) > 0 {
		hintBlock = "\n\nvia54Design 风格约束:\n" + strings.Join(hintParts, "\n")
	}

	return fmt.Sprintf(`你是一个资深前端工程师. 仔细观察这张 UI 截图, 用纯 HTML + 内联 CSS 复刻它的视觉结构.

要求:
1. 输出单一 HTML 文件 (含 <style> 块, 内联 CSS)
2. 使用现代 CSS: flexbox / grid / clamp() / CSS variables
3. 响应式: 至少支持 mobile (375px) 和 desktop (1280px)
4. 不使用外部 JS 框架 (Vanilla JS 可, 但推荐纯 CSS 实现交互)
5. 字体用系统字体栈或 Google Fonts CDN
6. 颜色用 CSS 变量 :root { --xxx: ... }
7. 保留原图的所有视觉细节: 间距、圆角、阴影、字号比例
8. 如果原图有图片, 用 <img> 标签 + placeholder 或 unsplash 链接
9. 不需要复刻品牌 logo 文字, 可用 placeholder%s

直接输出 HTML, 不要再加解释.`, hintBlock)
}

// ─── Provider: OpenAI ───

func callOpenAI(model, apiKey, mime, base64Img, prompt string, timeoutSec int) (string, error) {
	if model == "" {
		model = "gpt-4o"
	}
	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "text", "text": prompt},
					{"type": "image_url", "image_url": map[string]string{
						"url": fmt.Sprintf("data:%s;base64,%s", mime, base64Img),
					}},
				},
			},
		},
		"max_tokens": 4096,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	return doHTTP(req, timeoutSec)
}

// ─── Provider: Anthropic ───

func callAnthropic(model, apiKey, mime, base64Img, prompt string, timeoutSec int) (string, error) {
	if model == "" {
		model = "claude-sonnet-4-5"
	}
	url := "https://api.anthropic.com/v1/messages"
	body := map[string]interface{}{
		"model":      model,
		"max_tokens": 4096,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "image", "source": map[string]string{
						"type":       "base64",
						"media_type": mime,
						"data":       base64Img,
					}},
					{"type": "text", "text": prompt},
				},
			},
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	return doHTTP(req, timeoutSec)
}

// ─── Provider: Gemini ───

func callGemini(model, apiKey, mime, base64Img, prompt string, timeoutSec int) (string, error) {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
					{"inline_data": map[string]string{
						"mime_type": mime,
						"data":      base64Img,
					}},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 4096,
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	return doHTTP(req, timeoutSec)
}

// ─── HTTP dispatcher ───

func callVisionAPI(provider, model, apiKey, mime, base64Img, prompt string, timeoutSec int) (string, error) {
	switch provider {
	case "openai":
		return callOpenAI(model, apiKey, mime, base64Img, prompt, timeoutSec)
	case "anthropic":
		return callAnthropic(model, apiKey, mime, base64Img, prompt, timeoutSec)
	case "gemini":
		return callGemini(model, apiKey, mime, base64Img, prompt, timeoutSec)
	default:
		return "", fmt.Errorf("未知 provider %q (支持: openai/anthropic/gemini)", provider)
	}
}

func doHTTP(req *http.Request, timeoutSec int) (string, error) {
	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body)[:minBytes(500, len(body))])
	}

	// 解析各 provider 的响应
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// OpenAI
	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		return result.Choices[0].Message.Content, nil
	}
	// Anthropic / Gemini
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", fmt.Errorf("LLM 返回空 content. Raw: %s", string(body)[:minBytes(500, len(body))])
}

func minBytes(a, b int) int {
	if a < b {
		return a
	}
	return b
}