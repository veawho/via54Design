// via54Design — Phase D: Prometheus /metrics + pprof 端点
//
// 设计: AGENTS.md 规定 Go stdlib 优先 (仅允许 mcp-go + yaml.v3)
//       → 不引 prometheus/client_golang, 自己写最小 Prometheus 文本格式
//       → pprof 用 stdlib net/http/pprof
//
// 用法: ServeHTTP 改成 ServeWithMetrics, 包一层 http.ServeMux:
//       /metrics         → Prometheus 文本格式 (Go runtime + 自定义)
//       /debug/pprof/*   → stdlib pprof
//       /sse             → 转发给 mcp-go SSEServer
//       /message         → 转发给 mcp-go SSEServer
package mcp

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"runtime/metrics"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// Metrics 自定义业务指标计数器
type Metrics struct {
	StartTime    time.Time
	HTTPRequests atomic.Int64
	MCPCalls     atomic.Int64
	ToolCalls    sync.Map // tool name → *atomic.Int64
}

// NewMetrics 创建指标实例
func NewMetrics() *Metrics {
	return &Metrics{StartTime: time.Now()}
}

// IncToolCall 原子加一
func (m *Metrics) IncToolCall(tool string) {
	v, _ := m.ToolCalls.LoadOrStore(tool, &atomic.Int64{})
	v.(*atomic.Int64).Add(1)
}

// IncMCPCall 原子加一
func (m *Metrics) IncMCPCall() {
	m.MCPCalls.Add(1)
}

// IncHTTPRequest 原子加一
func (m *Metrics) IncHTTPRequest() {
	m.HTTPRequests.Add(1)
}

// WritePrometheus 写 Prometheus 文本格式到 w
//
// 输出示例:
//   # HELP via54_uptime_seconds Process uptime in seconds
//   # TYPE via54_uptime_seconds gauge
//   via54_uptime_seconds 123.45
//   ...
func (m *Metrics) WritePrometheus(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	uptime := time.Since(m.StartTime).Seconds()
	fmt.Fprintf(w, "# HELP via54_uptime_seconds Process uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE via54_uptime_seconds gauge\n")
	fmt.Fprintf(w, "via54_uptime_seconds %.2f\n\n", uptime)

	fmt.Fprintf(w, "# HELP via54_http_requests_total Total HTTP requests\n")
	fmt.Fprintf(w, "# TYPE via54_http_requests_total counter\n")
	fmt.Fprintf(w, "via54_http_requests_total %d\n\n", m.HTTPRequests.Load())

	fmt.Fprintf(w, "# HELP via54_mcp_calls_total Total MCP tool calls\n")
	fmt.Fprintf(w, "# TYPE via54_mcp_calls_total counter\n")
	fmt.Fprintf(w, "via54_mcp_calls_total %d\n\n", m.MCPCalls.Load())

	// 业务: 每个 tool 单独一个 metric
	fmt.Fprintf(w, "# HELP via54_tool_calls_total Total calls per MCP tool\n")
	fmt.Fprintf(w, "# TYPE via54_tool_calls_total counter\n")
	m.ToolCalls.Range(func(k, v any) bool {
		fmt.Fprintf(w, `via54_tool_calls_total{tool="%s"} %d`+"\n", k, v.(*atomic.Int64).Load())
		return true
	})
	fmt.Fprintln(w)

	// Go runtime (使用 runtime/metrics 拿 50+ 系统指标)
	fmt.Fprintln(w, "# HELP go_runtime Go runtime metrics")
	fmt.Fprintln(w, "# TYPE go_runtime gauge")
	samples := []metrics.Sample{
		{Name: "/memory/classes/heap/objects:bytes"},
		{Name: "/memory/classes/heap/bytes:bytes"},
		{Name: "/gc/heap/objects:objects"},
		{Name: "/goroutine/sched/goroutines:goroutines"},
	}
	metrics.Read(samples)
	for _, s := range samples {
		switch s.Value.Kind() {
		case metrics.KindUint64:
			fmt.Fprintf(w, "go_runtime{name=%q} %d\n", s.Name, s.Value.Uint64())
		case metrics.KindFloat64:
			fmt.Fprintf(w, "go_runtime{name=%q} %.2f\n", s.Name, s.Value.Float64())
		}
	}

	// Go runtime/memstats (传统指标)
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Fprintf(w, "\n# HELP go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP go_memstats_heap_alloc_bytes Bytes of allocated heap objects\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_alloc_bytes %d\n", ms.HeapAlloc)

	fmt.Fprintf(w, "# HELP go_gc_duration_seconds Total GC time\n")
	fmt.Fprintf(w, "# TYPE go_gc_duration_seconds counter\n")
	fmt.Fprintf(w, "go_gc_duration_seconds %.6f\n", float64(ms.PauseTotalNs)/1e9)
}

// ServeWithMetrics 启动 HTTP server, 包含 /metrics + /debug/pprof/* + /sse 代理
func (s *Server) ServeWithMetrics(addr string) error {
	metrics := NewMetrics()
	mux := http.NewServeMux()

	// Prometheus
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.IncHTTPRequest()
		metrics.WritePrometheus(w)
	})

	// pprof (stdlib)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		metrics.IncHTTPRequest()
		fmt.Fprintln(w, `{"status":"ok","uptime_sec":`, int(time.Since(metrics.StartTime).Seconds()), `}`)
	})

	// MCP SSE 代理
	sseServer := server.NewSSEServer(s.mcp, server.WithBaseURL("http://"+addr))
	mux.Handle("/sse", sseServer)
	mux.Handle("/message", sseServer)

	fmt.Fprintf(os.Stderr, "via54-mcp HTTP server on %s\n", addr)
	fmt.Fprintf(os.Stderr, "  Prometheus:  http://%s/metrics\n", addr)
	fmt.Fprintf(os.Stderr, "  pprof:       http://%s/debug/pprof/\n", addr)
	fmt.Fprintf(os.Stderr, "  SSE:         http://%s/sse\n", addr)
	fmt.Fprintf(os.Stderr, "  health:      http://%s/health\n", addr)
	return http.ListenAndServe(addr, mux)
}