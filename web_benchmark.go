package go_webkit

import (
	"fmt"
	"sync"
	"time"
)

type taskFunc = func(index int) error

// WebBenchMark 基准测试类型
type WebBenchMark struct {
	host string
}

// NewWebBenchMark 构建基准测试函数
func NewWebBenchMark(host string) *WebBenchMark {
	return &WebBenchMark{host}
}

// RunSingleAPI 单个API基准测试
func (sf *WebBenchMark) RunSingleAPI(relativePath string, tps int, rounds int, interval time.Duration,
	task taskFunc) BenchMarkCount {
	apiFunc := func(index int) error {
		cli := NewWebClient(sf.host)
		body, err := cli.FormGET(relativePath, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(body))
		return task(index)
	}
	return RunBenchMark(tps, rounds, interval, apiFunc)
}

// BenchMarkCount 基准测试统计
type BenchMarkCount struct {
	Begin      time.Time
	End        time.Time
	RoundCount map[int]*RoundCount
}

// RoundCount 单轮统计
type RoundCount struct {
	Begin      time.Time
	End        time.Time
	TaskCounts map[int]*TaskCount
}

// TaskCount 单次统计
type TaskCount struct {
	Begin  time.Time
	End    time.Time
	Status bool
}

// RunBenchMark 一次基准测试
func RunBenchMark(tps int, rounds int, interval time.Duration, task taskFunc) BenchMarkCount {
	benchMarkBegin := time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(rounds)
	roundCount := make(map[int]*RoundCount)
	for r := 0; r < rounds; r++ {
		go runRound(r, roundCount, tps, wg, task)
		// 延时等待
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	benchMarkEnd := time.Now()
	return BenchMarkCount{benchMarkBegin, benchMarkEnd, roundCount}
}

// runRound 一轮并发
func runRound(index int, countMap map[int]*RoundCount, tps int, wg *sync.WaitGroup, task taskFunc) {
	roundWG := new(sync.WaitGroup)
	roundWG.Add(tps)
	taskCount := make(map[int]*TaskCount)
	roundBegin := time.Now()
	for t := 0; t < tps; t++ {
		go runTask(t, taskCount, roundWG, task)
	}
	roundWG.Wait()
	roundEnd := time.Now()
	countMap[index] = &RoundCount{roundBegin, roundEnd, taskCount}
	wg.Done()
}

// runTask 单个任务
func runTask(index int, countMap map[int]*TaskCount, wg *sync.WaitGroup, task taskFunc) {
	taskBegin := time.Now()
	err := task(index)
	taskEnd := time.Now()
	isSuccess := true
	if err != nil {
		isSuccess = false
	}
	countMap[index] = &TaskCount{taskBegin, taskEnd, isSuccess}
	wg.Done()
}
