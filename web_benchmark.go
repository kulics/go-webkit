package webkit

import "sync"
import "time"

type taskFunc = func(int) error
type Web_Benchmark struct {
	cli *WebClient
}

func New_Web_Benchmark(cli *WebClient) (v *Web_Benchmark) {
	return &Web_Benchmark{cli}
}
func (me *Web_Benchmark) Run_Single_API(tps int, rounds int, interval time.Duration, req func(*WebClient, int) error) (count Benchmark_Count) {
	api := func(index int) (e error) {
		return req(me.cli, index)
	}
	return Run_Benchmark(tps, rounds, interval, api)
}

type Benchmark_Count struct {
	Begin       time.Time
	End         time.Time
	Round_Count map[int]*Round_Count
}
type Round_Count struct {
	Begin       time.Time
	End         time.Time
	Task_Counts map[int]*Task_Count
}
type Task_Count struct {
	Begin  time.Time
	End    time.Time
	Status bool
}

func Run_Benchmark(tps int, rounds int, interval time.Duration, task taskFunc) (count Benchmark_Count) {
	var benchmark_begin = time.Now()
	var wg = &sync.WaitGroup{}
	wg.Add(rounds)
	var round_count = map[int]*Round_Count{}
	for r := 0; r < rounds; r += 1 {
		go run_Round(r, round_count, tps, wg, task)
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	var benchmark_end = time.Now()
	return Benchmark_Count{benchmark_begin, benchmark_end, round_count}
}
func run_Round(index int, countMap map[int]*Round_Count, tps int, wg *sync.WaitGroup, task taskFunc) {
	var roundWG = &sync.WaitGroup{}
	roundWG.Add(tps)
	var taskCount = map[int]*Task_Count{}
	var roundBegin = time.Now()
	for t := 0; t < tps; t += 1 {
		go run_Task(t, taskCount, roundWG, task)
	}
	roundWG.Wait()
	var roundEnd = time.Now()
	countMap[index] = &Round_Count{roundBegin, roundEnd, taskCount}
	wg.Done()
}
func run_Task(index int, countMap map[int]*Task_Count, wg *sync.WaitGroup, task taskFunc) {
	var taskBegin = time.Now()
	var err = task(index)
	var taskEnd = time.Now()
	var isSuccess = true
	if err != nil {
		isSuccess = false
	}
	countMap[index] = &Task_Count{taskBegin, taskEnd, isSuccess}
	wg.Done()
}
