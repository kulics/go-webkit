"webkit" {
	"fmt"
	"testing"
}

Test upload file(t: ?testing.T) -> () {
	cli := New Web Client("http://localhost:8080/")
	(body, err) := cli.Upload_file("file", "file", "README.md",
		[Str]{}{["path"]"read.me"})
	? err >< () {
		fmt.Println(err)
		<- ()
	}
	fmt.Println(string(body))
}

Test download file(t:?testing.T) -> () {
	cli := New Web Client("http://localhost:8080/")
	err := cli.Download_file(
		"file?path=../web_client.go",
		"./tempFile", [Str]{}{})
	? err >< () {
		fmt.Println(err)
		<- ()
	}
}
