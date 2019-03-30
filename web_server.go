package go_webkit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context 类型别名
type Context = *gin.Context

// WebServer web服务辅助工具
type WebServer struct {
	listen string
	engine *gin.Engine
}

// NewWebServerDefault 构建web服务数据对象
// host可以为空
func NewWebServerDefault(host string, port string) *WebServer {
	return &WebServer{fmt.Sprintf("%s:%s", host, port), gin.Default()}
}

// Run 运行web服务
func (sf *WebServer) Run() error {
	return sf.engine.Run(sf.listen)
}

// HandleFunc 监听函数
func (sf *WebServer) HandleFunc(method Method, url string, handle func(ctx Context)) *WebServer {
	switch method {
	case Get:
		sf.engine.GET(url, handle)
	case Post:
		sf.engine.POST(url, handle)
	case Put:
		sf.engine.PUT(url, handle)
	case Delete:
		sf.engine.DELETE(url, handle)
	case Patch:
		sf.engine.PATCH(url, handle)
	case Options:
		sf.engine.OPTIONS(url, handle)
	}
	return sf
}

// HandleFuncGet 监听Get
func (sf *WebServer) HandleFuncGet(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Get, url, handle)
}

// HandleFuncPost 监听Post
func (sf *WebServer) HandleFuncPost(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Post, url, handle)
}

// HandleFuncPut 监听Put
func (sf *WebServer) HandleFuncPut(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Put, url, handle)
}

// HandleFuncDelete 监听Delete
func (sf *WebServer) HandleFuncDelete(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Delete, url, handle)
}

// HandleFuncPatch 监听Patch
func (sf *WebServer) HandleFuncPatch(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Patch, url, handle)
}

// HandleFuncOptions 监听Options
func (sf *WebServer) HandleFuncOptions(url string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Options, url, handle)
}

// HandleStruct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
// 结构体的方法必须与 Method 类型的名称一致
func (sf *WebServer) HandleStruct(url string, handle interface{}) *WebServer {
	sf.handleStruct(url, handle)
	return sf
}

// handleStruct 使用反射遍历结构体的方法和字段，对http方法进行注册
func (sf *WebServer) handleStruct(url string, handle interface{}) {
	rfType := reflect.TypeOf(handle)
	rfValue := reflect.ValueOf(handle)
	// 只接受结构体、接口及指针
	switch rfType.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Struct:
		// 反射方法
		for i := 0; i < rfType.NumMethod(); i++ {
			methodName := rfType.Method(i).Name
			if !isMethod(NewMethod(methodName)) {
				continue
			}
			handleFunc := func(ctx Context) {
				rfValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
			}
			sf.HandleFunc(NewMethod(methodName), url, handleFunc)
		}
		// 反射字段
		for i := 0; i < rfType.NumField(); i++ {
			fieldName := rfType.Field(i).Name
			sf.handleStruct(url+"/"+strings.ToLower(fieldName[:1])+fieldName[1:],
				rfValue.FieldByName(fieldName).Interface())
		}
	}
}
