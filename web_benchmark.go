package webkit

import (
	"time"

	"github.com/kulics/go-benchmark"
)

// WebBenchmark 基准测试类型
type WebBenchmark struct {
	cli *WebClient
}

// Run single API 单个API基准测试
func (me *WebBenchmark) RunSingleAPI(tps int, rounds int, interval time.Duration, req func(*WebClient, int) error) benchmark.BenchmarkCount {
	api := func(index int) (e error) {
		return req(me.cli, index)
	}
	return benchmark.RunBenchmark(tps, rounds, interval, api)
}

// NewWebBenchmark 构建基准测试函数
func NewWebBenchmark(cli *WebClient) *WebBenchmark {
	return &WebBenchmark{cli}
}
