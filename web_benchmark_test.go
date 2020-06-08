package webkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

//  测试单api
func TestBenchmarkSingleAPI(t *testing.T) {
	bm := NewWebBenchmark(NewWebClient("http://baidu.com/"))
	count := bm.RunSingleAPI(10, 3, 1000, func(cli *WebClient, index int) error {
		body, err := cli.Form_GET("ping", nil)
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
