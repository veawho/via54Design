# via54Design — AI Agent Instructions

> GitHub Copilot 自动加载此文件。主上下文文件: [`AGENTS.md`](../AGENTS.md)

## 关键约束
- Go 核心引擎，标准库优先，仅 mcp-go + yaml.v3 两个外部依赖
- YAML 是唯一设计数据源（layouts/color-schemes/typography/prompts）
- 所有 map 遍历必须用 sortedKeys() 保证确定性输出
- Python 管道 (hack/via54_pipeline.py) 仅负责 LLM 编排，零外部依赖
- SPARQL-License-Identifier 头部必须出现在每个 .go 文件

## 常用命令
```
go build -o via54.exe ./cmd/via54/
go vet ./...
./via54.exe prompt --scene "..." --platform midjourney
./via54.exe prompt list
python test_20_rounds.py             # 20轮端到端测试
```

## 架构速览
```
User → Python管道(i18n+LLM) → Go CLI(结构化执行) → YAML模板 → 输出
```
