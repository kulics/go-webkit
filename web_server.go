package webkit

import (
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
// domain可以为空
func NewWebServerDefault(listen string) *WebServer {
	return &WebServer{listen, gin.Default()}
}

// Run 运行web服务
func (sf *WebServer) Run() error {
	return sf.engine.Run(sf.listen)
}

// HandleFunc 监听函数
func (sf *WebServer) HandleFunc(method Method, relativePath string, handle func(ctx Context)) *WebServer {
	switch method {
	case GET:
		sf.engine.GET(relativePath, handle)
	case POST:
		sf.engine.POST(relativePath, handle)
	case PUT:
		sf.engine.PUT(relativePath, handle)
	case DELETE:
		sf.engine.DELETE(relativePath, handle)
	case PATCH:
		sf.engine.PATCH(relativePath, handle)
	case OPTIONS:
		sf.engine.OPTIONS(relativePath, handle)
	}
	return sf
}

// HandleFuncGet 监听Get
func (sf *WebServer) HandleFuncGet(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(GET, relativePath, handle)
}

// HandleFuncPost 监听Post
func (sf *WebServer) HandleFuncPost(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(POST, relativePath, handle)
}

// HandleFuncPut 监听Put
func (sf *WebServer) HandleFuncPut(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(PUT, relativePath, handle)
}

// HandleFuncDelete 监听Delete
func (sf *WebServer) HandleFuncDelete(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(DELETE, relativePath, handle)
}

// HandleFuncPatch 监听Patch
func (sf *WebServer) HandleFuncPatch(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(PATCH, relativePath, handle)
}

// HandleFuncOptions 监听Options
func (sf *WebServer) HandleFuncOptions(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(OPTIONS, relativePath, handle)
}

// HandleStruct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
// 结构体的方法必须与 Method 类型的名称一致
func (sf *WebServer) HandleStruct(relativePath string, handle interface{}) *WebServer {
	sf.handleStruct(relativePath, handle)
	return sf
}

// handleStruct 使用反射遍历结构体的方法和字段，对http方法进行注册
func (sf *WebServer) handleStruct(relativePath string, handle interface{}) {
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
			sf.HandleFunc(NewMethod(methodName), relativePath, handleFunc)
		}
		// 反射字段
		for i := 0; i < rfType.NumField(); i++ {
			fieldName := rfType.Field(i).Name
			sf.handleStruct(relativePath+"/"+strings.ToLower(fieldName[:1])+fieldName[1:],
				rfValue.FieldByName(fieldName).Interface())
		}
	}
}
