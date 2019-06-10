package webkit

import (
	"fmt"
	"testing"
)

func TestFileUpload(t *testing.T) {
	cli := NewWebClient("http://localhost:8080/")
	body, err := cli.FileUpload("file", "file", "README.md",
		map[string]interface{}{"path": "read.me"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestFileDownload(t *testing.T) {
	cli := NewWebClient("http://localhost:8080/")
	err := cli.FileDownload(
		"file?path=../web_client.go",
		"./tempFile", map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
