package webkit

import "fmt"
import "testing"

func Test_upload_file(t *testing.T) {
	cli := New_Web_Client("http://localhost:8080/")
	body, err := cli.Upload_file("file", "file", "README.md", map[string]interface{}{"path": "read.me"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
func Test_download_file(t *testing.T) {
	cli := New_Web_Client("http://localhost:8080/")
	err := cli.Download_file("file?path=../web_client.go", "./tempFile", map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
