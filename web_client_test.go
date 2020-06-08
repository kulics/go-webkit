package webkit

import (
	"fmt"
	"testing"
)

// 测试上传
func TestUploadFile(t *testing.T) {
	cli := NewWebClient("http://localhost:8080/")
	body, err := cli.UploadFile("file", "file", "README.md", map[string]interface{}{"path": "read.me"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// 测试下载
func TestDownloadFile(t *testing.T) {
	cli := NewWebClient("http://localhost:8080/")
	err := cli.DownloadFile("file?path=../web_client.go", "./tempFile", map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
