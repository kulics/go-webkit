package webkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBenchMarkSingleAPI(t *testing.T) {
	bm := New_Web_Benchmark(NewWebClient("http://baidu.com/"))
	count := bm.Run_Single_API(10, 3, 1000, func(cli *WebClient, index int) error {
		body, err := cli.FormGET("ping", nil)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	})
	bts, err := json.Marshal(count)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bts))
}
