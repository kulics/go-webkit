package go_webkit

import "testing"

func TestBenchMarkSingleAPI(t *testing.T) {
	bm := NewWebBenchMark("http://baidu.com", "80")
	bm.RunSingleAPI("", 100, 10, 1000, func(index int) {
		t.Log(index)
	})
}
