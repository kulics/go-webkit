package webkit

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context 类型别名
type Context = *gin.Context

// Web_Server web服务辅助工具
type Web_Server struct {
	listen string
	engine *gin.Engine
}

// NewWebServerDefault 构建web服务数据对象
// domain可以为空
func New_Web_Server_Default(listen string) *Web_Server {
	return &Web_Server{listen, gin.Default()}
}

// Run 运行web服务
func (sf *Web_Server) Run() error {
	return sf.engine.Run(sf.listen)
}

// Handle_Func 监听函数
func (sf *Web_Server) Handle_Func(method Method, relativePath string, handle func(ctx Context)) *Web_Server {
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

// Handle_GET 监听Get
func (sf *Web_Server) Handle_GET(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(GET, relativePath, handle)
}

// Handle_POST 监听Post
func (sf *Web_Server) Handle_POST(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(POST, relativePath, handle)
}

// Handle_PUT 监听Put
func (sf *Web_Server) Handle_PUT(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(PUT, relativePath, handle)
}

// Handle_DELETE 监听Delete
func (sf *Web_Server) Handle_DELETE(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(DELETE, relativePath, handle)
}

// Handle_PATCH 监听Patch
func (sf *Web_Server) Handle_PATCH(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(PATCH, relativePath, handle)
}

// Handle_OPTIONS 监听Options
func (sf *Web_Server) Handle_OPTIONS(relativePath string, handle func(ctx Context)) *Web_Server {
	return sf.Handle_Func(OPTIONS, relativePath, handle)
}

// Handle_Struct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
// 结构体的方法必须与 Method 类型的名称一致
func (sf *Web_Server) Handle_Struct(relativePath string, handle interface{}) *Web_Server {
	sf.handle_struct(relativePath, handle)
	return sf
}

// handle_struct 使用反射遍历结构体的方法和字段，对http方法进行注册
func (sf *Web_Server) handle_struct(relativePath string, handle interface{}) {
	rfType := reflect.TypeOf(handle)
	rfValue := reflect.ValueOf(handle)
	// 只接受结构体、接口及指针
	switch rfType.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Struct:
		// 反射方法
		for i := 0; i < rfType.NumMethod(); i++ {
			methodName := rfType.Method(i).Name
			if !is_Method(New_Method(methodName)) {
				continue
			}
			handleFunc := func(ctx Context) {
				rfValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
			}
			sf.Handle_Func(New_Method(methodName), relativePath, handleFunc)
		}
		// 反射字段
		for i := 0; i < rfType.NumField(); i++ {
			fieldName := rfType.Field(i).Name
			sf.handle_struct(relativePath+"/"+strings.ToLower(fieldName[:1])+fieldName[1:],
				rfValue.FieldByName(fieldName).Interface())
		}
	}
}
