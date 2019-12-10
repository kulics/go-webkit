package webkit

import "sync"
import "time"

type taskFunc = func(int) error
type WebBenchmark struct {
	cli *WebClient
}

func (me *WebBenchmark) RunSingleAPI(tps int, rounds int, interval time.Duration, req func(*WebClient, int) error) (count BenchmarkCount) {
	api := func(index int) (e error) {
		return req(me.cli, index)
	}
	return RunBenchmark(tps, rounds, interval, api)
}
func NewWebBenchmark(cli *WebClient) (v *WebBenchmark) {
	return &WebBenchmark{cli}
}

type BenchmarkCount struct {
	Begin      time.Time
	End        time.Time
	RoundCount map[int]*RoundCount
}
type RoundCount struct {
	Begin      time.Time
	End        time.Time
	TaskCounts map[int]*TaskCount
}
type TaskCount struct {
	Begin  time.Time
	End    time.Time
	Status bool
}

func RunBenchmark(tps int, rounds int, interval time.Duration, task taskFunc) (count BenchmarkCount) {
	benchmarkBegin := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(rounds)
	roundCount := map[int]*RoundCount{}
	for r := 0; r <= rounds-1; r += 1 {
		go runRound(r, roundCount, tps, wg, task)
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	benchmarkEnd := time.Now()
	return BenchmarkCount{benchmarkBegin, benchmarkEnd, roundCount}
}
func runRound(index int, countMap map[int]*RoundCount, tps int, wg *sync.WaitGroup, task taskFunc) {
	roundWG := &sync.WaitGroup{}
	roundWG.Add(tps)
	taskCount := map[int]*TaskCount{}
	roundBegin := time.Now()
	for t := 0; t <= tps-1; t += 1 {
		go runTask(t, taskCount, roundWG, task)
	}
	roundWG.Wait()
	roundEnd := time.Now()
	countMap[index] = &RoundCount{roundBegin, roundEnd, taskCount}
	wg.Done()
}
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
