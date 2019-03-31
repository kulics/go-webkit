package main

import (
	"fmt"
	"net/http"

	"github.com/kulics/go_webkit"
)

func main() {
	err := go_webkit.NewWebServerDefault("localhost:8080").
		HandleFuncGet("ping", func(ctx go_webkit.Context) {
			ctx.String(http.StatusOK, "pong")
		}).
		HandleStruct("api", testRouter{}).
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
