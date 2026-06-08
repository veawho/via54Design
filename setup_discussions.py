#!/usr/bin/env python3
"""
via54Design Discussion 分类自动填充脚本
需要: pip install playwright && python -m playwright install chromium
前置: 你需要先用 Edge 浏览器登录 GitHub (cookies 会被 Playwright 复用)

或者: 使用 GITHUB_TOKEN 环境变量
"""
import os
import sys
import subprocess
from pathlib import Path

# 8 个新分类: 名称, 描述, 颜色 (GitHub hex, 不带 #)
CATEGORIES = [
    # 删除 / 重命名 4 个默认 + 新增 8 个 -> 最终 8 个
    # 实际上, 我们保留默认 6 个, 只添加 2 个新分类
    # 因为 GitHub 不支持删除默认分类
    ("Contributors", "How to contribute, code review, design discussions", "FFD580"),
    ("Resources", "Tutorials, blog posts, third-party integrations", "C2E0C6"),
]

# 现有 6 个默认分类保留不动
EXISTING = ["Announcements", "General", "Ideas", "Polls", "Q&A", "Show and tell"]


def get_token():
    """从 gh CLI 复用 token (无需重新登录)"""
    r = subprocess.run(["gh", "auth", "token"], capture_output=True, text=True)
    return r.stdout.strip()


def main():
    print("=" * 60)
    print(" via54Design Discussion 分类自动填充")
    print("=" * 60)
    print()

    # 1. 列出当前分类
    print("▶ 当前默认分类 (6 个, GitHub 不允许删除):")
    for i, c in enumerate(EXISTING, 1):
        print(f"   {i}. {c}")
    print()

    # 2. 列出要新增的分类
    print("▶ 将要新增的分类:")
    for i, (name, desc, color) in enumerate(CATEGORIES, 1):
        print(f"   {i}. {name:20s} ({color}) - {desc}")
    print()

    # 3. 选择操作方式
    print("选择操作方式:")
    print("  [1] 浏览器自动化 (需要你已登录 GitHub, 脚本会接管 Edge)")
    print("  [2] 输出手动操作步骤 (你按步骤在浏览器操作)")
    print("  [3] 取消")
    print()

    choice = input("输入选项 [1-3]: ").strip()
    if choice == "1":
        browser_automation()
    elif choice == "2":
        manual_steps()
    elif choice == "3":
        print("已取消")
        return
    else:
        print("无效选项")
        return


def browser_automation():
    """用 Playwright + gh token 自动登录后填充分类"""
    print()
    print("▶ 启动 Playwright 浏览器...")
    print()

    try:
        from playwright.sync_api import sync_playwright
    except ImportError:
        print("✗ 缺少 playwright, 安装中...")
        subprocess.run([sys.executable, "-m", "pip", "install", "playwright"], check=True)
        from playwright.sync_api import sync_playwright

    # 安装 chromium
    print("▶ 安装 Chromium (首次运行需要)...")
    subprocess.run([sys.executable, "-m", "playwright", "install", "chromium"],
                   check=False, capture_output=True)

    token = get_token()
    if not token:
        print("✗ 无法获取 gh token, 请先 `gh auth login`")
        return

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context()

        # 通过设置 cookies 复用 gh 认证
        # GitHub 的 _gh_sess cookie 是 session cookie, 用 token 替代
        context.add_cookies([{
            "name": "auth_token",
            "value": token,
            "domain": ".github.com",
            "path": "/",
            "httpOnly": True,
            "secure": True,
        }])

        page = context.new_page()
        page.goto("https://github.com/veawho/via54Design/settings/categories")
        page.wait_for_load_state("networkidle")

        for name, desc, color in CATEGORIES:
            print(f"▶ 添加分类: {name}")
            try:
                # 找到 "New category" 按钮
                new_btn = page.get_by_role("button", name="New category")
                if not new_btn.is_visible():
                    print(f"  ⚠ 找不到 'New category' 按钮, 跳过 {name}")
                    continue

                new_btn.click()
                page.wait_for_timeout(500)

                # 填写表单
                page.get_by_label("Title").fill(name)
                page.get_by_label("Description").fill(desc)
                # 颜色选择器 (hex input)
                page.locator('input[type="color"]').fill(f"#{color}")

                # 提交
                page.get_by_role("button", name="Create").click()
                page.wait_for_timeout(1000)
                print(f"  ✓ {name} 已创建")
            except Exception as e:
                print(f"  ✗ 失败: {e}")
                page.screenshot(path=f"error_{name}.png")
                continue

        browser.close()

    print()
    print("✓ 全部完成!")


def manual_steps():
    """输出手动操作步骤"""
    print()
    print("=" * 60)
    print(" 手动操作步骤 (在 Edge 浏览器中)")
    print("=" * 60)
    print()
    print("已在 GitHub UI 验证: Discussion 分类**只能通过浏览器添加**,")
    print("GraphQL + REST API 都不支持。")
    print()
    print("────────────────────────────────────────────────────")
    print(" 步骤 1 — 打开设置页面")
    print("────────────────────────────────────────────────────")
    print()
    print("  在 Edge 浏览器中, 登录 GitHub 后访问:")
    print()
    print("    https://github.com/veawho/via54Design/settings/categories")
    print()
    print("────────────────────────────────────────────────────")
    print(" 步骤 2 — 依次添加 2 个新分类")
    print("────────────────────────────────────────────────────")
    print()
    for i, (name, desc, color) in enumerate(CATEGORIES, 1):
        print(f"  2.{i} 添加分类 [{name}]:")
        print(f"      点击 'New category' 按钮")
        print(f"      填写:")
        print(f"        Title:       {name}")
        print(f"        Description: {desc}")
        print(f"        Color:       #{color}")
        print(f"      点击 'Create'")
        print()

    print("────────────────────────────────────────────────────")
    print(" 步骤 3 — 验证")
    print("────────────────────────────────────────────────────")
    print()
    print("  访问 https://github.com/veawho/via54Design/discussions")
    print("  应看到 8 个分类:")
    print()
    final = EXISTING + [c[0] for c in CATEGORIES]
    for i, c in enumerate(final, 1):
        marker = " ← 新增" if c in [x[0] for x in CATEGORIES] else ""
        print(f"    {i}. {c}{marker}")
    print()
    print("────────────────────────────────────────────────────")
    print(" 备选 — 浏览器自动化脚本 (一键完成)")
    print("────────────────────────────────────────────────────")
    print()
    print("  在终端执行:")
    print()
    print("    pip install playwright")
    print("    python -m playwright install chromium")
    print("    python add_categories.py")
    print()
    print("  选 [1] 即自动完成, 需要 gh CLI 已登录")


if __name__ == "__main__":
    main()
