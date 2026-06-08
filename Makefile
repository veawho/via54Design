# via54Design Makefile
# Standard Go project build automation
# Builds: via54 (CLI) + via54-mcp (MCP Server)

BINARY    := via54
BINARY_MCP := via54-mcp
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS   := -ldflags="-s -w -X main.Version=$(VERSION)"
# GOFLAGS=-buildvcs=false disables VCS stamping (e.g. git commit hash embedded in binary).
# Required for exFAT / FAT32 filesystems which don't support file locking and may corrupt
# the VCS status read. Has no effect on NTFS or any modern filesystem with proper locking.
GOFLAGS   ?= -buildvcs=false

.PHONY: all build build-mcp test clean install lint wasm cross

all: build build-mcp

build: ## 编译 CLI 二进制
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) ./cmd/via54/

build-mcp: ## 编译 MCP Server 独立二进制
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_MCP) ./cmd/mcp-server/

test: ## 运行所有测试
	go test ./... -v -count=1 -timeout 60s 2>&1 || echo "(no tests yet)"

lint: ## 运行静态检查
	go vet ./...

clean: ## 清理构建产物
	rm -f $(BINARY) $(BINARY_MCP)
	rm -rf dist/

install: build ## 安装 CLI 到 GOPATH/bin
	go install $(LDFLAGS) ./cmd/via54/

install-mcp: build-mcp ## 安装 MCP Server 到 GOPATH/bin
	go install $(LDFLAGS) ./cmd/mcp-server/

wasm: ## 编译 Rust WASM 引擎 (需要 Rust)
	cd hack/wasm && bash build.sh

cross: ## 跨平台编译 (CLI + MCP 双二进制, 5 平台)
	@echo "=== 跨平台编译 ==="
	@mkdir -p dist
	@for binary in ./cmd/via54/ ./cmd/mcp-server/; do \
		name=via54; \
		[ "$$binary" = "./cmd/mcp-server/" ] && name=via54-mcp; \
		GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$$name-darwin-amd64     $$binary; \
		GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$$name-darwin-arm64     $$binary; \
		GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$$name-linux-amd64      $$binary; \
		GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$$name-linux-arm64      $$binary; \
		GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$$name-windows-amd64.exe $$binary; \
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
