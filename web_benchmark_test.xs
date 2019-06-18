"webkit" {
	"encoding/json"
	"fmt"
	"testing"
}

# 测试单api #
Test Benchmark Single API(t: ?testing.T) -> () {
	bm := New Web Benchmark(NewWebClient("http://baidu.com/"))
	count := bm.Run Single API(10, 3, 1000, (cli:?WebClient, index:Int) -> (err:error) {
		body, err := cli.FormGET("ping", nil)
		? err >< () {
			<- (err)
		}
		fmt.Println(string(body))
		<- (())
	})
	bts, err := json.Marshal(count)
	? err >< () {
		fmt.Println(err)
	}
	fmt.Println(string(bts))
}
