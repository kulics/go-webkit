package webkit

import "encoding/json"
import "fmt"
import "testing"

func Test_Benchmark_single_API(t *testing.T) {
	bm := New_Web_Benchmark(New_Web_Client("http://baidu.com/"))
	count := bm.Run_single_API(10, 3, 1000, func(cli *Web_Client, index int) (err error) {
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
