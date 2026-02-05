#!/usr/bin/env python3
"""Generate screenshots for a batch of themes."""
import subprocess
import sys
from pathlib import Path
from ansi2html import Ansi2HTMLConverter
from playwright.sync_api import sync_playwright

PROJECT_DIR = Path(__file__).parent.parent
THEMES_DIR = PROJECT_DIR / "assets" / "themes"

html_template = '''<!DOCTYPE html>
<html>
<head>
<style>
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap');
* {{ margin: 0; padding: 0; box-sizing: border-box; }}
body {{
    background: #0d1117;
    padding: 16px 20px;
    font-family: 'JetBrains Mono', 'Monaco', 'Menlo', monospace;
    font-size: 13px;
    line-height: 1.5;
    display: inline-block;
}}
pre {{ margin: 0; white-space: pre; font-family: inherit; }}
</style>
</head>
<body><pre>{content}</pre></body>
</html>
'''

themes = sys.argv[1:] if len(sys.argv) > 1 else []
if not themes:
    print("Usage: python gen_batch.py theme1 theme2 ...")
    sys.exit(1)

conv = Ansi2HTMLConverter(inline=True, dark_bg=True)

with sync_playwright() as p:
    browser = p.chromium.launch()
    page = browser.new_page()

    for theme in themes:
        print(f"Generating: {theme}")
        result = subprocess.run(
            ["./statusline", "--preview", theme],
            capture_output=True, text=True, cwd=PROJECT_DIR
        )
        html_content = conv.convert(result.stdout, full=False)
        page.set_content(html_template.format(content=html_content))
        page.wait_for_timeout(200)
        page.locator("body").screenshot(path=str(THEMES_DIR / f"{theme}.png"))
        print(f"  -> {theme}.png")

    browser.close()
print("Done!")
