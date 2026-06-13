// SPDX-License-Identifier: AGPL-3.0-only
// `via54 img` — 通过 mmx (MiniMax) CLI 调用真生图 API
// mmx 未装时: 输出 prompt 给用户手动贴 web
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/veawho/via54Design/internal/prompt"
)

func cmdImg() {
	fs := flag.NewFlagSet("img", flag.ExitOnError)
	scene := fs.String("scene", "", "场景描述 (必填)")
	platform := fs.String("platform", "minimax", "目标平台 (minimax/midjourney/kling/...); 默认 minimax 走 mmx CLI")
	ar := fs.String("ar", "16:9", "宽高比: 1:1/16:9/9:16/4:3/3:4 (传给 mmx --aspect-ratio)")
	n := fs.Int("n", 1, "批量生成数 (1-4, 传给 mmx --num-images)")
	out := fs.String("out", "./minimax-output/", "输出目录 (mmx 默认 cwd/minimax-output/)")
	seed := fs.Int("seed", -1, "固定 seed 复现 (-1 = 随机, 0..2^31-1 固定)")
	prefix := fs.String("prefix", "image", "输出文件名前缀 (避免覆盖)")
	dryRun := fs.Bool("dry-run", false, "只生成 prompt 不真调用 mmx (用于调试)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `via54 img — 生图 (通过 mmx MiniMax CLI)

用法:
  via54 img --scene "a red cat on a vintage book" [--ar 16:9] [--n 1] [--out ./out/]

平台:
  minimax (默认) — 调 mmx image generate (需 mmx CLI + API key)
  其他 via54 平台 — 仅生成 prompt, 不真生图

需要 mmx CLI:
  npm i -g @MiniMax-AI/cli   (或: npx skills add MiniMax-AI/cli -y -g)
  mmx auth login             (API key 在 platform.minimaxi.com / platform.minimax.io)

示例:
  via54 img --scene "cyberpunk city at night, 16:9" --ar 16:9
  via54 img --scene "a cat on books" --platform midjourney --dry-run
`)
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[2:])
	if *scene == "" {
		fmt.Fprintln(os.Stderr, "❌ --scene 不能为空")
		fs.Usage()
		os.Exit(1)
	}

	// 1. 走 via54 prompt 翻译器拿最终 prompt (minimax 平台)
	promptOut, err := generateViaPrompt(*scene, *platform)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ prompt 翻译失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📋 平台: %s\n", *platform)
	fmt.Printf("📝 Prompt: %s\n\n", promptOut)

	// 2. 非 minimax 平台 — dry-run 默认行为 (只给 prompt, 让用户贴 MJ/Kling 网页)
	if *platform != "minimax" {
		fmt.Printf("ℹ️  平台 %q 不通过 mmx 调用, 请手动复制上面 prompt 到对应平台.\n", *platform)
		return
	}

	// 3. minimax — 检查 mmx
	mmxPath, err := exec.LookPath("mmx")
	if err != nil || *dryRun {
		if *dryRun {
			fmt.Println("🧪 --dry-run 模式: 跳过 mmx 调用")
		} else {
			fmt.Fprintln(os.Stderr, "❌ mmx CLI 未找到 (PATH 里没有 'mmx')")
			fmt.Fprintln(os.Stderr, "   安装: npm i -g @MiniMax-AI/cli")
			fmt.Fprintln(os.Stderr, "   登录: mmx auth login")
			fmt.Fprintln(os.Stderr, "   文档: https://platform.minimaxi.com/docs/token-plan/minimax-cli")
		}
		return
	}
	fmt.Printf("🔧 mmx 路径: %s\n", mmxPath)

	// 4. 准备 mmx 调用
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 创建输出目录失败: %v\n", err)
		os.Exit(1)
	}
	absOut, _ := filepath.Abs(*out)

	// mmx image generate --prompt "<prompt 主体>" --aspect-ratio <ar> --n <count>
	// (flag 名称: --prompt / --aspect-ratio / --n / --seed / --out / --out-dir)
	// 关键: prompt 来自 via54 翻译器的 FinalPrompt — 可能已含 "--ar ..." 字符串
	//       这些是模板注释残留,不是真 mmx flag,需要剥掉避免重复/冲突
	promptClean := stripTemplateFlags(promptOut, "--ar", "--n", "--v", "--style", "--s")
	args := []string{
		"image", "generate",
		"--prompt", promptClean,
		"--aspect-ratio", *ar,
		"--n", fmt.Sprintf("%d", *n),
		"--out-prefix", *prefix,
	}
	if *seed >= 0 {
		args = append(args, "--seed", fmt.Sprintf("%d", *seed))
	}
	fmt.Printf("🚀 调用: %s %s\n", mmxPath, strings.Join(args, " "))
	fmt.Printf("📁 输出目录: %s\n\n", absOut)

	cmd := exec.Command(mmxPath, args...)
	cmd.Dir = absOut
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ mmx 调用失败: %v\n", err)
		fmt.Fprintln(os.Stderr, "   排查: 1) mmx auth status  2) mmx image generate --help  3) mmx quota")
		os.Exit(1)
	}
	fmt.Printf("\n✅ 完成, 文件应在: %s\n", absOut)
}

// generateViaPrompt 通过 via54 prompt 引擎拿最终 prompt 字符串
func generateViaPrompt(scene, platform string) (string, error) {
	s, err := prompt.GeneratePrompt(scene, platform, "", baseDir())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s.FinalPrompt), nil
}

// stripTemplateFlags 剥离 via54 prompt 模板拼出的 --xxx yyy 片段,避免传 mmx 冲突
// 例: "a cat --ar 16:9 --n 1 --v 6.1" → "a cat"
func stripTemplateFlags(s string, flags ...string) string {
	fields := strings.Fields(s)
	out := make([]string, 0, len(fields))
	skipNext := false
	for _, f := range fields {
		if skipNext {
			skipNext = false
			continue
		}
		isFlag := false
		for _, fl := range flags {
			if f == fl {
				isFlag = true
				skipNext = true
				break
			}
		}
		if !isFlag {
			out = append(out, f)
		}
	}
	return strings.Join(out, " ")
}
