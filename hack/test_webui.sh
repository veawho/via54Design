#!/usr/bin/env bash
# via54Design WebUI 测试脚本 (Shell 版)
# 用法: bash hack/test_webui.sh [--watch]
set -e

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
PASS=0; FAIL=0
ok()   { PASS=$((PASS+1)); echo -e "  ${GREEN}✅${NC} $1"; }
fail() { FAIL=$((FAIL+1)); echo -e "  ${RED}❌${NC} $1"; }

REPO="$(cd "$(dirname "$0")/.." && pwd)"
EXE="$REPO/via54.exe"
PORT=19994

echo -e "\n${CYAN}╔══ via54Design WebUI Test Suite ══╗${NC}\n"

# 1: Build
echo "═══ [1/9] Build ═══"
if go build -o "$EXE" ./cmd/via54/ 2>/dev/null; then ok "go build compiles"
else fail "go build"; exit 1; fi

# 2: CLI commands
echo -e "\n═══ [2/9] CLI Commands ═══"
for cmd in serve generate narrate quality pattern list media export prompt comfyui forge web version; do
    if $EXE 2>&1 | grep -q "$cmd"; then ok "CLI: $cmd"; else fail "CLI: $cmd"; fi
done

# 3: Start server
echo -e "\n═══ [3/9] Server Start ═══"
$EXE web --port $PORT &
SPID=$!
sleep 2
if kill -0 $SPID 2>/dev/null; then ok "Server starts"; else fail "Server start"; fi

# 4: API endpoints
echo -e "\n═══ [4/9] API Endpoints ═══"
curl -sf http://localhost:$PORT/api/health > /dev/null && ok "GET /api/health" || fail "GET /api/health"
curl -sf http://localhost:$PORT/api/templates > /dev/null && ok "GET /api/templates" || fail "GET /api/templates"
curl -sf -X POST http://localhost:$PORT/api/prompt -d '{"scene":"cat","platform":"midjourney"}' > /dev/null && ok "POST /api/prompt" || fail "POST /api/prompt"
curl -sf -X POST http://localhost:$PORT/api/narrate -d '{"seed":"test","model":"three-act","duration":15}' > /dev/null && ok "POST /api/narrate" || fail "POST /api/narrate"
curl -sf -X POST http://localhost:$PORT/api/generate -d '{"layout":"hero-split-16-9","color":"ink-wash","font":"ming-hei-editorial","title":"t"}' > /dev/null && ok "POST /api/generate" || fail "POST /api/generate"
curl -sf -X POST http://localhost:$PORT/api/build -d '{"workflow_id":"sdxl_txt2img","prompt":"cat"}' > /dev/null && ok "POST /api/build" || fail "POST /api/build"
curl -sf -X POST http://localhost:$PORT/api/export -d '{"type":"json","source":"test","output":"/tmp/test.json"}' > /dev/null && ok "POST /api/export" || fail "POST /api/export"
curl -sf -X POST http://localhost:$PORT/api/media -d '{"action":"list"}' > /dev/null && ok "POST /api/media" || fail "POST /api/media"

# 5: HTML structure
echo -e "\n═══ [5/9] HTML Structure ═══"
HTML=$(curl -sf http://localhost:$PORT/)
[ -n "$HTML" ] && ok "Page loads" || fail "Page loads"
echo "$HTML" | grep -q 'toggleLang' && ok "i18n toggle" || fail "i18n toggle"
echo "$HTML" | grep -q '/api/prompt' && ok "Prompt API ref" || fail "Prompt API ref"
echo "$HTML" | grep -q '/api/narrate' && ok "Narrate API ref" || fail "Narrate API ref"
echo "$HTML" | grep -q '/api/build' && ok "Build API ref" || fail "Build API ref"
echo "$HTML" | grep -q 'minmax(0, 1fr)' && ok "minmax grid fix" || fail "minmax grid fix"
echo "$HTML" | grep -q ':focus-visible' && ok "focus-visible styles" || fail "focus-visible styles"
echo "$HTML" | grep -q 'prefers-reduced-motion' && ok "reduced motion" || fail "reduced motion"
echo "$HTML" | grep -q -v 'googleapis\|cdnjs\|unpkg\|jsdelivr' && ok "Zero CDN" || fail "Zero CDN"

# 6: Index page size
echo -e "\n═══ [6/9] Index Page ═══"
SIZE=$(echo "$HTML" | wc -c)
echo "$HTML" | grep -qi '<!DOCTYPE html>' && ok "Valid HTML ($SIZE bytes)" || fail "Valid HTML"

# 7: Stress test
echo -e "\n═══ [7/9] Stress Test ═══"
ERRORS=0
for i in $(seq 1 50); do
    curl -sf http://localhost:$PORT/api/health > /dev/null 2>&1 || ERRORS=$((ERRORS+1))
done
[ $ERRORS -eq 0 ] && ok "50 requests: 0 errors" || fail "50 requests: $ERRORS errors"

# 8: Error boundaries
echo -e "\n═══ [8/9] Error Boundaries ═══"
curl -sf -X POST http://localhost:$PORT/api/prompt -d '{}' | grep -q 'error' && ok "Missing params" || fail "Missing params"
curl -sf -X POST http://localhost:$PORT/api/build -d '{"workflow_id":"invalid","prompt":"cat"}' | grep -q 'error' && ok "Invalid workflow" || fail "Invalid workflow"
curl -sf http://localhost:$PORT/api/build | grep -q 'error' && ok "GET on POST" || fail "GET on POST"
curl -sf http://localhost:$PORT/api/nonexistent -w '%{http_code}' | grep -q 404 && ok "404 route" || fail "404 route"

# Kill server
kill $SPID 2>/dev/null; wait $SPID 2>/dev/null

# Summary
TOTAL=$((PASS+FAIL))
PCT=$((PASS*100/TOTAL))
echo -e "\n${CYAN}╔══ Summary ═══════════════════════╗${NC}"
echo -e "${CYAN}║  Total: $TOTAL  Pass: $PASS ($PCT%)  Fail: $FAIL${NC}"
echo -e "${CYAN}╚══════════════════════════════════╝${NC}"
[ $FAIL -eq 0 ] && exit 0 || exit 1
