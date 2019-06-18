"webkit" {
	"sync"
	"time"
}

taskFunc => (Int)->(error) 

# Web Benchmark 基准测试类型 #
Web Benchmark -> {
	cli: ?WebClient
}

# New Web Benchmark 构建基准测试函数 #
New Web Benchmark(cli: ?WebClient) -> (v: ?Web Benchmark) {
	<- (Web Benchmark{cli}?)
}

# Run Single API 单个API基准测试 #
(me: ?Web Benchmark) Run Single API(tps: Int, rounds: Int, interval: time.Duration, req: (?WebClient, Int)->(error) ) -> (count:Benchmark Count) {
	api(index: Int) -> (e:error) {
		<- (req(me.cli, index))
	}
	<- (Run Benchmark(tps, rounds, interval, api))
}

# Benchmark Count 基准测试统计 #
Benchmark Count -> {
	Begin:          time.Time
	End:            time.Time
	Round Count:    [Int]?Round Count
}

# Round Count 单轮统计 #
Round Count -> {
	Begin:          time.Time
	End:            time.Time
	Task Counts:    [Int]?Task Count
}

# Task Count 单次统计 #
Task Count -> {
	Begin:      time.Time
	End:        time.Time
	Status:     Bool
}

# Run Benchmark 一次基准测试 #
Run Benchmark(tps: Int, rounds: Int, interval: time.Duration, task: taskFunc) -> (count: Benchmark Count) {
	benchmark begin := time.Now()
	wg := sync.WaitGroup{}?
	wg.Add(rounds)
	round count := [Int]?Round Count{}
	[0 < rounds] @ r {
		run Round(r, round count, tps, wg, task) <~
		# 延时等待 #
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	benchmark end := time.Now()
	<- (Benchmark Count{benchmark begin, benchmark end, round count})
}

# runRound 一轮并发 #
run Round(index:Int, countMap:[Int]?Round Count, tps:Int, wg:?sync.WaitGroup, task:taskFunc) -> () {
	roundWG := sync.WaitGroup{}?
	roundWG.Add(tps)
	taskCount := [Int]?Task Count{}
	roundBegin := time.Now()
	[0 < tps] @ t {
		run Task(t, taskCount, roundWG, task) <~
	}
	roundWG.Wait()
	roundEnd := time.Now()
	countMap[index] = Round Count{roundBegin, roundEnd, taskCount}?
	wg.Done()
}

# runTask 单个任务 #
run Task(index:Int, countMap:[Int]?Task Count, wg:?sync.WaitGroup, task:taskFunc) -> () {
	taskBegin := time.Now()
	err := task(index)
	taskEnd := time.Now()
	isSuccess := true
	? err >< () {
		isSuccess = false
	}
	countMap[index] = Task Count{taskBegin, taskEnd, isSuccess}?
	wg.Done()
}
