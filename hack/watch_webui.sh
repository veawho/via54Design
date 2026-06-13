#!/usr/bin/env bash
# via54Design WebUI 变更检测脚本 (Shell 版)
# 由 Hermes cron 每10分钟触发, 用于 hack/test_webui.sh
set -e

REPO="$(cd "$(dirname "$0")/.." && pwd)"
CHECKSUM_FILE="$REPO/.webui_checksums.json"

# Compute checksums for watched files
compute_checksums() {
    find "$REPO/web" "$REPO/cmd" "$REPO/internal" "$REPO/templates/workflows" \
        -type f \( -name "*.go" -o -name "*.html" -o -name "*.yaml" -o -name "*.json" -o -name "*.sh" \) \
        ! -path "*/vendor/*" ! -name "*.bak" 2>/dev/null | \
    while read f; do
        echo "$(md5sum "$f" | cut -d' ' -f1)  $f"
    done
}

# Load saved checksums
load_saved() {
    if [ -f "$CHECKSUM_FILE" ]; then
        cat "$CHECKSUM_FILE"
    fi
}

# Save new checksums
save_checksums() {
    compute_checksums > "$CHECKSUM_FILE"
}

# Detect changes
detect_changes() {
    local changed=0
    local tmpfile
    tmpfile=$(mktemp)
    compute_checksums > "$tmpfile"
    
    if [ ! -f "$CHECKSUM_FILE" ]; then
        changed=1
        echo "🔄 First run — all files will be checksummed"
    elif ! diff -q "$CHECKSUM_FILE" "$tmpfile" > /dev/null 2>&1; then
        changed=1
        local count=$(diff "$CHECKSUM_FILE" "$tmpfile" | grep "^>" | wc -l)
        echo "🔄 $count file(s) changed"
    fi
    
    mv "$tmpfile" "$CHECKSUM_FILE"
    return $changed
}

# Main
cd "$REPO"

if detect_changes; then
    echo "   No changes detected"
    exit 0
fi

echo "   Running tests..."
if bash "$REPO/hack/test_webui.sh"; then
    save_checksums
    echo "✅ All tests passed"
    exit 0
else
    echo "❌ Tests failed"
    exit 1
fi
