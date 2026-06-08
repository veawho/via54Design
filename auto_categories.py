#!/usr/bin/env python3
"""用 Playwright 验证登录态, 然后填充分类"""
import subprocess
import sys
from pathlib import Path

def get_gh_token():
    r = subprocess.run(["gh", "auth", "token"], capture_output=True, text=True)
    return r.stdout.strip()

def main():
    from playwright.sync_api import sync_playwright

    token = get_gh_token()
    if not token:
        print("✗ gh CLI 未登录, 请先运行: gh auth login")
        return 1

    print(f"✓ gh token 长度: {len(token)}")

    CATEGORIES = [
        ("Contributors", "How to contribute, code review, design discussions", "FFD580"),
        ("Resources", "Tutorials, blog posts, third-party integrations", "C2E0C6"),
    ]

    with sync_playwright() as p:
        # Use headless=False so we can see what happens
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()

        # GitHub 个人 token 可以用 Authorization header 注入
        # 但访问 /settings 需要 session cookie
        # 改为: 用 gh 的 user_session cookie
        context.add_cookies([{
            "name": "user_session",
            # GitHub uses this cookie for auth; gh token alone doesn't work for all endpoints
            # Try with auth_token (works for REST/GraphQL API)
            "value": token,
            "domain": ".github.com",
            "path": "/",
            "httpOnly": True,
            "secure": True,
        }])

        page = context.new_page()

        # 1. 先验证登录态
        print("▶ 访问 categories 设置页...")
        page.goto("https://github.com/veawho/via54Design/settings/categories",
                  wait_until="networkidle")

        # 看是否登录
        if "Sign in" in page.content() and page.get_by_text("Sign in").first.is_visible():
            print("✗ 未登录, token cookie 不工作 (GitHub 使用 user_session 加密 cookie)")
            print("  请用 'gh auth login --web' 登录后, gh token 才可用于 API")
            print("  但 settings 页面需要 user_session cookie, 只能通过浏览器登录")
            return 1

        # 2. 添加分类
        for name, desc, color in CATEGORIES:
            print(f"▶ 添加 {name}...")
            try:
                new_btn = page.get_by_role("button", name="New category")
                new_btn.first.click()
                page.wait_for_timeout(800)

                page.get_by_label("Title").first.fill(name)
                page.get_by_label("Description").first.fill(desc)
                # color picker
                color_input = page.locator('input[type="color"]').first
                if color_input.is_visible():
                    color_input.evaluate(f'el => el.value = "#{color}"')

                # submit
                page.get_by_role("button", name="Create").first.click()
                page.wait_for_timeout(1500)
                print(f"  ✓ {name}")
            except Exception as e:
                print(f"  ✗ {name}: {e}")
                page.screenshot(path=f"error_{name}.png")

        browser.close()
        print()
        print("✓ 完成!")

if __name__ == "__main__":
    sys.exit(main() or 0)
