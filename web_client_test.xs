"webkit" {
	"fmt"
	"testing"
}

Test upload file(t: ?testing.T) -> () {
	cli := New Web Client("http://localhost:8080/")
	(body, err) := cli.Upload_file("file", "file", "README.md",
		[Str]Any{["path"]"read.me"})
	? err >< Nil {
		fmt.Println(err)
		<- ()
	}
	fmt.Println(string(body))
}

Test download file(t:?testing.T) -> () {
	cli := New Web Client("http://localhost:8080/")
	err := cli.Download_file(
		"file?path=../web_client.go",
		"./tempFile", [Str]Any{})
	? err >< Nil {
		fmt.Println(err)
		<- ()
	}
}
