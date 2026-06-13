// via54Design — 交互式菜单（傻瓜式操作）
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func cmdInteractive() {
	reader := bufio.NewReader(os.Stdin)
	for {
		showMenu()
		fmt.Print("请输入选项 [0-9]: ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "0":
			fmt.Println("👋 再见！")
			return
		case "1":
			interactiveGenerate(reader)
		case "2":
			interactivePresent(reader)
		case "3":
			interactivePrompt(reader)
		case "4":
			interactiveNarrate(reader)
		case "5":
			interactiveWeb()
		case "6":
			interactiveCheckDeps()
		default:
			fmt.Println("❌ 无效选项，请重新输入")
		}
	}
}

func showMenu() {
	clearScreen()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  via54Design  —  设计模板引擎 + 叙事引擎")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("🎨  设  计")
	fmt.Println("  1.  生成 HTML 页面 — 选择布局/配色/字体，生成网页")
	fmt.Println("  2.  生成演示文稿 — 导出为 PPTX 演示文件")
	fmt.Println()
	fmt.Println("📝  提示词")
	fmt.Println("  3.  写 AI 生图提示词 — 描述场景，自动生成结构化提示词")
	fmt.Println()
	fmt.Println("📖  叙  事")
	fmt.Println("  4.  生成故事框架 — 写一句话，AI 扩展为完整故事")
	fmt.Println()
	fmt.Println("🌐  Web UI")
	fmt.Println("  5.  启动可视化面板 — 浏览器操作界面（推荐新手）")
	fmt.Println()
	fmt.Println("🔧  工  具")
	fmt.Println("  6.  检查系统依赖 — 检测 Node.js / ffmpeg / Go")
	fmt.Println()
	fmt.Println("  0.  退出")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// ── 1. 生成 HTML 页面 ──

func interactiveGenerate(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("━━━ 🎨 生成 HTML 页面 ━━━")
	fmt.Println()
	fmt.Println("您想要生成一个什么样的页面？")
	fmt.Println()

	title := ask(reader, "页面标题", "我的设计")
	mode := ask(reader, "模式 (网页/演示)", "网页")

	presFlag := ""
	if strings.Contains(mode, "演示") {
		presFlag = "--presentation"
	}

	layout := choose(reader, "选择布局", []string{
		"hero-split-16-9  左右分割，适合首页",
		"bento-grid-2x2   网格布局，适合展示",
		"gallery-waterfall 瀑布流，适合画廊",
	}, "1")

	color := choose(reader, "选择配色", []string{
		"ink-wash          水墨风 — 中国传统",
		"crimson-elegance  朱砂 — 典雅中国红",
		"moon-white        雨过天青 — 清新淡雅",
		"dark-terminal-blue 科技蓝 — 现代",
		"neon-dark         霓虹暗色 — 潮流",
	}, "1")

	font := choose(reader, "选择字体", []string{
		"ming-hei-editorial  明黑 — 编辑风格",
		"kai-rounded-friendly 楷体 — 亲和",
		"sans-geometric-tech  无衬线几何 — 科技",
		"elegant-didone       迪多 — 时尚",
	}, "1")

	fmt.Println()
	fmt.Println("⏳ 正在生成...")
	fmt.Println()

	args := []string{"generate",
		"--layout", parseKey(layout),
		"--color", parseKey(color),
		"--font", parseKey(font),
		"--title", title,
	}
	if presFlag != "" {
		args = append(args, presFlag)
	}

	out, err := exec.Command(via54Path(), args...).CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n%s\n", err, out)
	} else {
		fmt.Printf("✅ 生成成功！(%d bytes)\n", len(out))
		fmt.Println()
		ask(reader, "按 Enter 返回菜单", "")
	}
}

// ── 2. 生成演示文稿 ──

func interactivePresent(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("━━━ 📊 生成演示文稿 ━━━")
	fmt.Println()

	seed := ask(reader, "输入故事种子（一句话描述你的演示）", "")
	if seed == "" {
		fmt.Println("❌ 故事种子不能为空")
		ask(reader, "按 Enter 返回菜单", "")
		return
	}

	model := choose(reader, "选择叙事模型", []string{
		"three-act        三幕剧 — 适合产品发布",
		"heros-journey    英雄之旅 — 适合品牌故事",
		"cognitive-arc    认知弧 — 适合科普教程",
		"problem-solution 问题解决 — 适合销售视频",
	}, "1")

	duration := ask(reader, "演示时长（秒）", "30")
	dur, _ := strconv.Atoi(duration)
	if dur <= 0 {
		dur = 30
	}

	fmt.Println()
	fmt.Println("⏳ 正在生成叙事框架...")

	args := []string{"narrate",
		"--seed", seed,
		"--model", parseKey(model),
		"--duration", strconv.Itoa(dur),
		"--format", "markdown",
	}

	out, err := exec.Command(via54Path(), args...).CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n%s\n", err, out)
	} else {
		fmt.Println("✅ 叙事框架生成成功！")
		fmt.Println()
		fmt.Println(string(out))
		fmt.Println()
		export := choose(reader, "是否导出为 PPTX？", []string{
			"yes  是 — 导出为 PPTX 文件",
			"no   否 — 只看文本",
		}, "2")
		if strings.HasPrefix(export, "yes") || export == "1" {
			args2 := []string{"export", "pptx", seed, "--output", seed + ".pptx"}
			out2, err2 := exec.Command(via54Path(), args2...).CombinedOutput()
			if err2 != nil {
				fmt.Printf("❌ PPTX 导出失败: %v\n", err2)
			} else {
				fmt.Printf("✅ PPTX 已导出: %s (%s)\n", seed+".pptx", string(out2))
			}
		}
	}
	ask(reader, "按 Enter 返回菜单", "")
}

// ── 3. 写提示词 ──

func interactivePrompt(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("━━━ 📝 写 AI 生图提示词 ━━━")
	fmt.Println()

	scene := ask(reader, "描述你想生成的画面（越详细越好）", "")
	if scene == "" {
		fmt.Println("❌ 场景描述不能为空")
		ask(reader, "按 Enter 返回菜单", "")
		return
	}

	platform := choose(reader, "选择目标平台", []string{
		"midjourney    Midjourney — 艺术感最强",
		"flux          Flux — 真实感最强",
		"dalle3        DALL-E 3 — 文字理解最好",
		"sd3           SD3 — 开源灵活",
		"kling         可灵 — 视频生成",
	}, "1")

	fmt.Println()
	fmt.Println("⏳ 正在生成提示词...")

	args := []string{"prompt",
		"--scene", scene,
		"--platform", parseKey(platform),
		"--format", "markdown",
	}

	out, err := exec.Command(via54Path(), args...).CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n%s\n", err, out)
	} else {
		fmt.Println("✅ 提示词生成成功！")
		fmt.Println()
		fmt.Println(string(out))
		fmt.Println()
		fmt.Println("💡 提示词已就绪，复制到对应平台使用")
	}
	ask(reader, "按 Enter 返回菜单", "")
}

// ── 4. 叙事框架 ──

func interactiveNarrate(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("━━━ 📖 生成故事框架 ━━━")
	fmt.Println()

	seed := ask(reader, "写一句话种子（故事的核心灵感）", "")
	if seed == "" {
		fmt.Println("❌ 故事种子不能为空")
		ask(reader, "按 Enter 返回菜单", "")
		return
	}

	model := choose(reader, "选择叙事模型", []string{
		"three-act        三幕剧 — 设问→解答→号召",
		"heros-journey    英雄之旅 — 日常→相遇→蜕变→回归",
		"cognitive-arc    认知弧 — 钩子→基础→核心→案例",
		"problem-solution 问题解决 — 痛点→方案→证明→行动",
	}, "1")

	duration := ask(reader, "视频/演示时长（秒）", "60")
	dur, _ := strconv.Atoi(duration)
	if dur <= 0 {
		dur = 60
	}

	fmt.Println()
	fmt.Println("⏳ 正在生成叙事框架...")

	args := []string{"narrate",
		"--seed", seed,
		"--model", parseKey(model),
		"--duration", strconv.Itoa(dur),
		"--format", "markdown",
	}

	out, err := exec.Command(via54Path(), args...).CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n%s\n", err, out)
	} else {
		fmt.Println("✅ 叙事框架生成成功！")
		fmt.Println()
		fmt.Println(string(out))
	}
	ask(reader, "按 Enter 返回菜单", "")
}

// ── 5. 启动 Web UI ──

// interactiveWeb 启动 Web UI 子进程，等待用户 Ctrl+C 退出
func interactiveWeb() {
	clearScreen()
	fmt.Println("━━━ 🌐 启动 Web UI ━━━")
	fmt.Println()
	fmt.Println("⏳ 正在启动 Web 服务器...")
	fmt.Println("   地址: http://localhost:8080")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println()
	fmt.Println("💡 提示: Web UI 启动后会持续运行。")
	fmt.Println("   在浏览器打开 http://localhost:8080 即可使用。")
	fmt.Println("   停止后会自动返回此菜单。")
	fmt.Println()

	cmd := exec.Command(via54Path(), "web", "--port", "8080")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // 传递 Ctrl+C
	if err := cmd.Run(); err != nil {
		// exit code 通常来自 Ctrl+C / 用户终止，不是真实错误
		fmt.Printf("\nℹ️ Web UI 已停止 (退出码非 0 通常是 Ctrl+C)\n")
	}
	fmt.Println()
	fmt.Println("🔙 返回菜单...")
}

// ── 6. 检查依赖 ──

func interactiveCheckDeps() {
	clearScreen()
	fmt.Println("━━━ 🔧 检查系统依赖 ━━━")
	fmt.Println()

	deps := []struct {
		name    string
		cmd     string
		install string
	}{
		{"Go (编译环境)", "go", "https://go.dev/dl/"},
		{"Node.js (PDF/视频)", "node", "https://nodejs.org"},
		{"ffmpeg (视频渲染)", "ffmpeg", "brew install ffmpeg / apt install ffmpeg"},
		{"Playwright (PDF导出)", "npx", "npm install -g playwright"},
	}

	allOk := true
	for _, dep := range deps {
		_, err := exec.LookPath(dep.cmd)
		if err != nil {
			fmt.Printf("  ❌ %s — 未安装\n     安装: %s\n", dep.name, dep.install)
			allOk = false
		} else {
			fmt.Printf("  ✅ %s — 已安装\n", dep.name)
		}
	}

	// 检查 via54.exe 自身
	self := via54Path()
	if _, err := os.Stat(self); os.IsNotExist(err) {
		fmt.Printf("  ⚠️ via54.exe — 未找到 (请先编译: go build -o via54.exe ./cmd/via54/)\n")
	} else {
		fmt.Printf("  ✅ via54.exe — 就绪\n")
	}

	fmt.Println()
	if allOk {
		fmt.Println("🎉 所有依赖已就绪，开始使用吧！")
	} else {
		fmt.Println("💡 部分依赖缺失，核心功能（生成/提示词/叙事）仍可使用")
		fmt.Println("   PDF导出和视频渲染需要 Node.js + ffmpeg")
	}
	ask(bufio.NewReader(os.Stdin), "按 Enter 返回菜单", "")
}

// ── 辅助函数 ──

// via54Path 返回当前可执行文件路径（跨平台）
func via54Path() string {
	exe, err := os.Executable()
	if err != nil {
		return "via54"
	}
	return exe
}

func ask(reader *bufio.Reader, label string, defaultVal string) string {
	promptText := "  " + label
	if defaultVal != "" {
		promptText += " [" + defaultVal + "]"
	}
	fmt.Print(promptText + ": ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

func choose(reader *bufio.Reader, label string, options []string, defaultIdx string) string {
	fmt.Println()
	fmt.Println("  " + label + ":")
	for i, opt := range options {
		mark := " "
		if strconv.Itoa(i+1) == defaultIdx {
			mark = ">"
		}
		fmt.Printf("  %s %d. %s\n", mark, i+1, opt)
	}
	fmt.Printf("  选择 [1-%d, 默认 %s]: ", len(options), defaultIdx)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		line = defaultIdx
	}
	if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(options) {
		return options[n-1]
	}
	return options[0]
}

func parseKey(option string) string {
	// 从 "key — description" 中提取 key
	if idx := strings.Index(option, " "); idx > 0 {
		return option[:idx]
	}
	return option
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
