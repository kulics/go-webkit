package go_webkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBenchMarkSingleAPI(t *testing.T) {
	bm := NewWebBenchMark("http://localhost:8080/")
	count := bm.RunSingleAPI("ping", 10, 3, 1000, func(index int) error {
		fmt.Println(index)
		return nil
	})
	bts, err := json.Marshal(count)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bts))
}
