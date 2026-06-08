#!/usr/bin/env python3
"""通过 Playwright 在 GitHub UI 添加 2 个缺失的 Discussion 分类。
流程: 设置 → Discussions → Categories → New category
"""
import subprocess
import json
import sys
import time
from pathlib import Path

# 1. 从 gh CLI 提取认证信息 (storageState 等)
#    然后用持久化 context 复用 cookie

# 简单方案: 调用 gh auth 命令获取 token, 然后手动构造 cookie
result = subprocess.run(['gh', 'auth', 'token'], capture_output=True, text=True, cwd='.')
TOKEN = result.stdout.strip()
print(f"Got gh token: {TOKEN[:20]}...")

# 2. 使用 Playwright + Edge + 持久化 storage state
from playwright.sync_api import sync_playwright

# 通过 gh auth login --with-token 风格构造的 storage state
# 但更简单: 通过 GitHub OAuth token 直接添加 cookie
storage_state = {
    "cookies": [
        {
            "name": "user_session",
            "value": TOKEN,
            "domain": ".github.com",
            "path": "/",
            "expires": -1,
            "httpOnly": True,
            "secure": True,
            "sameSite": "None"
        }
    ],
    "origins": [
        {
            "origin": "https://github.com",
            "localStorage": [
                {"name": "github_csrf_token", "value": TOKEN[:32]},
                {"name": "user-login", "value": "veawho"}
            ]
        }
    ]
}

categories = [
    {
        "name": "Contributors",
        "description": "How to contribute, design discussions, code review",
        "color": "FFD580"
    },
    {
        "name": "Resources",
        "description": "Tutorials, blog posts, integrations, third-party tools",
        "color": "C2E0C6"
    }
]

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True, channel="msedge")
    context = browser.new_context(storage_state=storage_state)
    page = context.new_page()

    print("\n1. Navigating to category settings...")
    page.goto("https://github.com/veawho/via54Design/settings/categories", timeout=30000)
    page.wait_for_load_state("networkidle")
    print(f"  Current URL: {page.url}")
    print(f"  Title: {page.title()}")

    # Check if logged in
    if "login" in page.url.lower() or page.url.endswith("categories") is False:
        print(f"  ⚠️  Not logged in. Page snapshot:")
        print(page.content()[:500])
        sys.exit(1)

    # Take a screenshot for debugging
    page.screenshot(path="categories_page.png", full_page=True)
    print("  Screenshot saved: categories_page.png")

    # Try to find the "New category" button
    for cat in categories:
        print(f"\n2. Creating category: {cat['name']}")
        try:
            # Click "New category" button
            page.get_by_role("button", name="New category").first.click()
            page.wait_for_timeout(2000)

            # Fill the form
            page.get_by_label("Title").fill(cat["name"])
            page.get_by_label("Description").fill(cat["description"])

            # Color is a color picker - skip or use default
            # Save
            page.get_by_role("button", name="Save").click()
            page.wait_for_timeout(2000)
            print(f"  ✓ Created: {cat['name']}")
        except Exception as e:
            print(f"  ✗ Failed: {e}")

    browser.close()

print("\nDone. Check categories_page.png for verification.")
