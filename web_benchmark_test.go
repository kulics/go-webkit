package go_webkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBenchMarkSingleAPI(t *testing.T) {
	bm := NewWebBenchMark(NewWebClient("http://localhost:8080/"))
	count := bm.RunSingleAPI("ping", 10, 3, 1000, func(cli *WebClient, index int) error {
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
