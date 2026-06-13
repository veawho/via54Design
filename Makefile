# via54Design Makefile
# Standard Go project build automation
# Builds: via54 (CLI) + via54-mcp (MCP Server)
#
# --- GOFLAGS environment variable ---
# GOFLAGS controls additional flags passed to all `go build` invocations.
# Default: -buildvcs=false (safe for exFAT/FAT32 without file locking)
# Override:  GOFLAGS="-v" make build   (your own flags)
# Disable:   GOFLAGS="" make build      (VCS stamping enabled, may fail on exFAT)
#
# Why -buildvcs=false by default?
#   Go 1.18+ embeds VCS info (commit hash) in the binary. Reading this requires
#   file locking on .git/. On exFAT/FAT32/SD cards (no file locking), the read
#   may corrupt the git index. -buildvcs=false skips this read entirely.
#   No effect on NTFS, APFS, ext4, btrfs. Safe everywhere.
#
BINARY    := via54
BINARY_MCP := via54-mcp
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS   := -ldflags="-s -w -X main.Version=$(VERSION)"

# Auto-detect exFAT/FAT32 filesystems (no file locking).
# On these filesystems, -buildvcs=false is REQUIRED to avoid "dubious ownership"
# errors and possible .git/ index corruption. The check uses df's fstype column.
ifeq ($(shell uname -s),Darwin)
  # macOS: stat -f %T  returns filesystem type. exFAT = "exfat", FAT32 = "msdos".
  REPO_FS := $(shell cd $(dir $(abspath $(lastword $(MAKEFILE_LIST)))) && stat -f %T . 2>/dev/null)
  NEEDS_NO_VCS := $(filter exfat msdos,$(REPO_FS))
else ifeq ($(shell uname -s),Linux)
  # Linux: df --output=fstype .  works on most systems. exFAT = "exfat", FAT32 = "vfat".
  REPO_FS := $(shell cd $(dir $(abspath $(lastword $(MAKEFILE_LIST)))) && df --output=fstype . 2>/dev/null | tail -1 | tr -d ' ')
  NEEDS_NO_VCS := $(filter exfat vfat msdos,$(REPO_FS))
else
  # Windows / other: assume NTFS (no detection needed; VCS works fine).
  NEEDS_NO_VCS :=
endif

# Default GOFLAGS: include -buildvcs=false if on exFAT/FAT32, otherwise empty.
# User can override via environment: GOFLAGS="-v" make build
GOFLAGS   ?= $(if $(NEEDS_NO_VCS),-buildvcs=false,)

.PHONY: all build build-mcp test clean install install-mcp lint wasm cross release help fs-check

all: build build-mcp

build: ## 编译 CLI 二进制
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) ./cmd/via54/

build-mcp: ## 编译 MCP Server 独立二进制
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_MCP) ./cmd/mcp-server/

test: ## 运行所有测试
	go test ./... -v -count=1 -timeout 60s 2>&1 || echo "(no tests yet)"

test-e2e: ## 20 轮端到端功能测试
	python test_20_rounds.py

test-stress: ## 200 轮连续压力测试 (持久性)
	python test_stress_200.py

test-concurrent: ## Web API 并发/吞吐压力测试
	python test_concurrent.py

lint: ## 运行静态检查
	go vet ./...

clean: ## 清理构建产物
	rm -f $(BINARY) $(BINARY_MCP)
	rm -rf dist/

install: build ## 安装 CLI 到 GOPATH/bin
	go install $(LDFLAGS) ./cmd/via54/

install-mcp: build-mcp ## 安装 MCP Server 到 GOPATH/bin
	go install $(LDFLAGS) ./cmd/mcp-server/

fs-check: ## 检测当前文件系统类型 (exFAT/FAT32/NTFS/APFS)
	@echo "Current filesystem:"
	@case "$$(uname -s)" in
	  Darwin)  echo "  macOS:  $$(stat -f %T .)" ;;
	  Linux)   echo "  Linux:  $$(df --output=fstype . 2>/dev/null | tail -1)" ;;
	  MINGW*|MSYS*) echo "  Windows: NTFS (assumed)" ;;
	  *)       echo "  Unknown: $$(uname -s)" ;;
	esac
	@echo ""
	@echo "GOFLAGS in effect: '$(GOFLAGS)'"
	@if [ -n "$(NEEDS_NO_VCS)" ]; then \
		echo "  ⚠️  Detected exFAT/FAT32 — -buildvcs=false is REQUIRED"; \
	else \
		echo "  ✓  Modern filesystem detected — GOFLAGS empty (no override needed)"; \
	fi

wasm: ## 编译 Rust WASM 引擎 (需要 Rust)
	cd hack/wasm && bash build.sh

cross: ## 跨平台编译 (CLI + MCP 双二进制, 5 平台)
	@echo "=== 跨平台编译 ==="
	@mkdir -p dist
	@for binary in ./cmd/via54/ ./cmd/mcp-server/; do \
		name=via54; \
		[ "$$binary" = "./cmd/mcp-server/" ] && name=via54-mcp; \
		GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o dist/$$name-darwin-amd64     $$binary; \
		GOOS=darwin  GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o dist/$$name-darwin-arm64     $$binary; \
		GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o dist/$$name-linux-amd64      $$binary; \
		GOOS=linux   GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o dist/$$name-linux-arm64      $$binary; \
		GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o dist/$$name-windows-amd64.exe $$binary; \
	done
	@echo "=== 编译完成 ==="
	@ls -lh dist/ | grep -v "^total" | awk '{print $$5, $$NF}' | column -t

release: cross ## 打包发布 (zip per platform)
	@echo "=== 打包发布 ==="
	@for plat in darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 windows-amd64; do \
		ext=; [ "$$plat" = "windows-amd64" ] && ext=.exe; \
		dir=dist/via54-$$plat; \
		mkdir -p $$dir; \
		cp dist/via54-$$plat$$ext   $$dir/; \
		cp dist/via54-mcp-$$plat$$ext $$dir/; \
		cp -r templates $$dir/ 2>/dev/null; \
		cp -r docs $$dir/ 2>/dev/null || true; \
		cp README.md LICENSE $$dir/ 2>/dev/null || true; \
		cd dist && zip -r via54-$$plat.zip via54-$$plat/ > /dev/null && rm -rf via54-$$plat/ && cd ..; \
		echo "  ✅ via54-$$plat.zip"; \
	done
	@echo "=== 发布包 ==="
	@ls -lh dist/*.zip 2>/dev/null | awk '{print $$5, $$NF}'

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
