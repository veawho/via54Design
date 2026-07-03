// via54Design — 设计模板引擎 + 叙事引擎
//
// via54 ai — 统一 LLM 文本生成接口 (在线 API, 零依赖)
// 借鉴: Anil-matcha/Open-Generative-AI (21k⭐ MIT) 的 "统一多模型接口" 思想
//
// 用法:
//   via54 ai "扩展这句话: ..." --provider openai
//   via54 ai "Translate to English: ..." --provider anthropic --model claude-sonnet-4-5
//   via54 ai --file prompt.md --provider gemini --output result.txt
//   via54 ai --provider mmx --model MiniMax-Text-01 "生成中文文案"
//
// 支持 provider: openai / anthropic / gemini / mmx / replicate
// 不引入 Python/FastAPI/前端, 保持 via54 单二进制零依赖路线.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func cmdAI() {
	fs := flag.NewFlagSet("ai", flag.ExitOnError)
	prompt := fs.String("prompt", "", "Prompt 字符串 (与 --file 二选一)")
	file := fs.String("file", "", "Prompt 文件路径")
	provider := fs.String("provider", "openai", "LLM provider: openai/anthropic/gemini/mmx/replicate")
	model := fs.String("model", "", "覆盖默认模型")
	system := fs.String("system", "你是 via54Design 的 AI 助手, 帮助用户做创意工作.", "System prompt")
	maxTokens := fs.Int("max-tokens", 2048, "最大输出 token 数")
	temperature := fs.Float64("temperature", 0.7, "Temperature (0.0-2.0)")
	apiKey := fs.String("api-key", os.Getenv("VIA54_LLM_API_KEY"), "API key")
	output := fs.String("output", "", "输出文件路径 (默认 stdout)")
	timeout := fs.Int("timeout", 60, "HTTP timeout 秒")
	fs.Parse(os.Args[2:])

	// 读 prompt
	userPrompt := *prompt
	if *file != "" {
		data, err := os.ReadFile(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
			os.Exit(1)
		}
		userPrompt = string(data)
	}
	if userPrompt == "" {
		// 取剩余 args 作为 prompt
		args := fs.Args()
		if len(args) > 0 {
			userPrompt = strings.Join(args, " ")
		}
	}
	if userPrompt == "" {
		fmt.Fprintln(os.Stderr, "必须指定 prompt (--prompt / --file / 末尾参数)")
		os.Exit(1)
	}

	// API key 兜底
	if *apiKey == "" {
		switch *provider {
		case "openai":
			*apiKey = os.Getenv("OPENAI_API_KEY")
		case "anthropic":
			*apiKey = os.Getenv("ANTHROPIC_API_KEY")
		case "gemini":
			*apiKey = os.Getenv("GEMINI_API_KEY")
		case "mmx":
			*apiKey = os.Getenv("MMX_API_KEY")
			if *apiKey == "" {
				*apiKey = os.Getenv("MINIMAX_API_KEY")
			}
		case "replicate":
			*apiKey = os.Getenv("REPLICATE_API_KEY")
		}
	}
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "需要 API key: --api-key 或 env")
		os.Exit(1)
	}

	// 调 LLM
	fmt.Fprintf(os.Stderr, "🤖 %s ", *provider)
	if *model != "" {
		fmt.Fprintf(os.Stderr, "(%s) ", *model)
	}
	fmt.Fprintf(os.Stderr, "...\n")

	result, err := callLLM(*provider, *model, *apiKey, *system, userPrompt, *maxTokens, *temperature, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LLM 调用失败: %v\n", err)
		os.Exit(1)
	}

	// 输出
	if *output != "" {
		if err := os.WriteFile(*output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "写文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "💾 %s (%d bytes)\n", *output, len(result))
	} else {
		fmt.Print(result)
	}
}

// ─── Provider: OpenAI ───

func llmOpenAI(model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}
	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": userPrompt},
		},
		"max_tokens":  maxTokens,
		"temperature": temp,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	return llmDoHTTP(req, timeoutSec, "openai")
}

// ─── Provider: Anthropic ───

func llmAnthropic(model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	if model == "" {
		model = "claude-sonnet-4-5"
	}
	url := "https://api.anthropic.com/v1/messages"
	body := map[string]interface{}{
		"model":      model,
		"max_tokens": maxTokens,
		"system":     system,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
		"temperature": temp,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	return llmDoHTTP(req, timeoutSec, "anthropic")
}

// ─── Provider: Gemini ───

func llmGemini(model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
	body := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{{"text": userPrompt}},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": maxTokens,
			"temperature":     temp,
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	return llmDoHTTP(req, timeoutSec, "gemini")
}

// ─── Provider: MiniMax ───

func llmMMX(model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	if model == "" {
		model = "MiniMax-Text-01"
	}
	url := "https://api.MiniMax.chat/v1/text/chatcompletion_v2"
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": userPrompt},
		},
		"max_tokens":  maxTokens,
		"temperature": temp,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	return llmDoHTTP(req, timeoutSec, "mmx")
}

// ─── Provider: Replicate ───

func llmReplicate(model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	// Replicate 走 prediction API, 较复杂: 先 create prediction → poll → get output
	if model == "" {
		model = "meta/meta-llama-3-70b-instruct"
	}
	url := "https://api.replicate.com/v1/predictions"
	body := map[string]interface{}{
		"version": model, // 默认假设传入的是 version hash
		"input": map[string]interface{}{
			"prompt":          userPrompt,
			"system_prompt":   system,
			"max_new_tokens":  maxTokens,
			"temperature":     temp,
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Authorization", "Token "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	return llmDoHTTP(req, timeoutSec, "replicate")
}

// ─── Dispatcher ───

func callLLM(provider, model, apiKey, system, userPrompt string, maxTokens int, temp float64, timeoutSec int) (string, error) {
	switch provider {
	case "openai":
		return llmOpenAI(model, apiKey, system, userPrompt, maxTokens, temp, timeoutSec)
	case "anthropic":
		return llmAnthropic(model, apiKey, system, userPrompt, maxTokens, temp, timeoutSec)
	case "gemini":
		return llmGemini(model, apiKey, system, userPrompt, maxTokens, temp, timeoutSec)
	case "mmx":
		return llmMMX(model, apiKey, system, userPrompt, maxTokens, temp, timeoutSec)
	case "replicate":
		return llmReplicate(model, apiKey, system, userPrompt, maxTokens, temp, timeoutSec)
	default:
		return "", fmt.Errorf("未知 provider %q (支持: openai/anthropic/gemini/mmx/replicate)", provider)
	}
}

func llmDoHTTP(req *http.Request, timeoutSec int, providerName string) (string, error) {
	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("[%s] HTTP %d: %s", providerName, resp.StatusCode, string(body)[:minBytes(500, len(body))])
	}

	// 各 provider 响应解析
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Text string `json:"text"`
		} `json:"choices"`
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Output []string `json:"output"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// OpenAI / mmx (chat completion)
	if len(result.Choices) > 0 {
		if result.Choices[0].Message.Content != "" {
			return result.Choices[0].Message.Content, nil
		}
		if result.Choices[0].Text != "" {
			return result.Choices[0].Text, nil
		}
	}
	// Anthropic / Gemini
	if len(result.Content) > 0 && result.Content[0].Text != "" {
		return result.Content[0].Text, nil
	}
	// Replicate
	if len(result.Output) > 0 {
		return strings.Join(result.Output, ""), nil
	}

	return "", fmt.Errorf("[%s] LLM 返回空 content. Raw: %s", providerName, string(body)[:minBytes(500, len(body))])
}