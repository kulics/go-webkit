package go_webkit

import (
	"fmt"
	"sync"
	"time"
)

type task = func(index int)

// WebBenchMark 基准测试类型
type WebBenchMark struct {
	host string
}

// NewWebBenchMark 构建基准测试函数
func NewWebBenchMark(domain string, port string) *WebBenchMark {
	return &WebBenchMark{fmt.Sprintf("%s:%s", domain, port)}
}

// RunSingleAPI 单个API基准测试
func (sf *WebBenchMark) RunSingleAPI(relativePath string, tps int, rounds int, interval time.Duration, task task) {
	apiFunc := func(index int) {
		cli := NewWebClient(sf.host, "")
		body, err := cli.FormGET("", nil)
		if err != nil {
			fmt.Println(err)
			// return nil, err
			return
		}
		fmt.Println(string(body))
		task(index)
	}
	runBenchMark(tps, rounds, interval, apiFunc)
}

// runBenchMark 一次基准测试
func runBenchMark(tps int, rounds int, interval time.Duration, task task) {
	benchMarkBegin := time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(rounds)
	for r := 0; r < rounds; r++ {
		go runRound(tps, wg, task)
		// 延时等待
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	benchMarkEnd := time.Now()
	fmt.Println(fmt.Sprintf("begin:%s  ->  end:%s", benchMarkBegin.String(), benchMarkEnd.String()))
}

// runRound 一轮并发
func runRound(tps int, wg *sync.WaitGroup, task task) {
	roundWG := new(sync.WaitGroup)
	roundWG.Add(tps)
	for t := 0; t < tps; t++ {
		go runTask(t, roundWG, task)
	}
	roundWG.Wait()
	wg.Done()
}

// runTask 单个任务
func runTask(index int, wg *sync.WaitGroup, task task) {
	task(index)
	wg.Done()
}
