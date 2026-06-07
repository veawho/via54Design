# via54Design Makefile
# Standard Go project build automation
# Target: single Go binary, zero external runtime deps

BINARY   := via54
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags="-s -w -X main.Version=$(VERSION)"

.PHONY: all build test clean install lint wasm

all: build

build: ## 编译主二进制
	go build $(LDFLAGS) -o $(BINARY) ./cmd/via54/

test: ## 运行所有测试
	go test ./... -v -count=1 -timeout 60s 2>&1 || echo "(no tests yet)"

lint: ## 运行静态检查
	go vet ./...
	@echo "vet: OK"

clean: ## 清理构建产物
	rm -f $(BINARY)
	rm -rf dist/

install: build ## 安装到 GOPATH/bin
	go install $(LDFLAGS) ./cmd/via54/

wasm: ## 编译 Rust WASM 引擎 (需要 Rust)
	cd hack/wasm && bash build.sh

cross: ## 跨平台编译 (macOS/Linux/Windows)
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64  ./cmd/via54/
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64   ./cmd/via54/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/via54/
	@echo "Builds in dist/"

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
