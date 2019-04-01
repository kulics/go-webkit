package main

import (
	"fmt"
	"net/http"
	"github.com/kulics/go_webkit"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

func main() {
	err := go_webkit.NewWebServerDefault("localhost:8080").
		HandleFuncGet("ping", func(ctx go_webkit.Context) {
			ctx.String(http.StatusOK, "pong")
		}).
		HandleStruct("api", testRouter{}).
		HandleStruct("file", fileRouter{}).
		Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type testRouter struct {
	Item testItem
}

func (testRouter) GET(ctx go_webkit.Context) {
	ctx.String(http.StatusOK, "get")
}

func (testRouter) POST(ctx go_webkit.Context) {
	ctx.String(http.StatusOK, "post")
}

type testItem struct{}

func (testItem) GET(ctx go_webkit.Context) {
	ctx.String(http.StatusOK, "get")
}

type fileRouter struct{}

func (fileRouter) GET(ctx go_webkit.Context) {
	fPath := ctx.Query("path")
	fmt.Println(fPath)
	ctx.File(fPath)
}

func (fileRouter) POST(ctx go_webkit.Context) {
	filePath := filepath.Clean(ctx.PostForm("path"))
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "upload file err:" + err.Error(),
		})
		return
	}
	if err := ctx.SaveUploadedFile(file, `./`+
		filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "upload file err:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}
