"webkit" {
	"encoding/json"
	"fmt"
	"testing"
}

# 测试单api #
Test Benchmark single API(t: ?testing.T) -> () {
	bm := New Web Benchmark(New Web Client("http://baidu.com/"))
	count := bm.Run single API(10, 3, 1000, (cli:?Web Client, index:Int) -> (err:error) {
		(body, err) := cli.Form GET("ping", nil)
		? err >< () {
			<- (err)
		}
		fmt.Println(string(body))
		<- (())
	})
	(bts, err) := json.Marshal(count)
	? err >< () {
		fmt.Println(err)
	}
	fmt.Println(string(bts))
}
