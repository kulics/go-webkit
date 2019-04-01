package go_webkit

import (
	"fmt"
	"testing"
)

func TestFileUpload(t *testing.T) {
	cli := NewWebClient("http://localhost:8080/")
	body, err := cli.FileUpload("file", "file", "README.md",
		map[string]string{"path": "read.me"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
