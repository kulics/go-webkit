package go_webkit

import (
	"testing"
	"fmt"
)

func TestBenchMarkSingleAPI(t *testing.T) {
	bm := NewWebBenchMark("http://localhost:8080/")
	bm.RunSingleAPI("ping", 10, 3, 1000, func(index int) {
		fmt.Println(index)
	})
}
