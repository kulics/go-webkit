package webkit

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context 类型别名
type Context = *gin.Context

// NewWebServerDefault 构建web服务数据对象,domain可以为空
func NewWebServerDefault(listen string) *WebServer {
	srv := &WebServer{}
	srv.listen = listen
	srv.engine = gin.Default()
	return srv
}

// WebServer web服务辅助工具
type WebServer struct {
	listen string
	engine *gin.Engine
}

// Run 运行web服务
func (me *WebServer) Run() error {
	return me.engine.Run(me.listen)
}

// Handle_Func 监听函数
func (me *WebServer) HandleFunc(method Method, relativePath string, handle func(Context)) *WebServer {
	switch method {
	case GET:
		me.engine.GET(relativePath, handle)
	case POST:
		me.engine.POST(relativePath, handle)
	case PUT:
		me.engine.PUT(relativePath, handle)
	case DELETE:
		me.engine.DELETE(relativePath, handle)
	case PATCH:
		me.engine.PATCH(relativePath, handle)
	case OPTIONS:
		me.engine.OPTIONS(relativePath, handle)
	}
	return me
}

// HandleGET 监听Get
func (me *WebServer) HandleGET(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(GET, relativePath, handle)
}

// HandlePOST 监听Post
func (me *WebServer) HandlePOST(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(POST, relativePath, handle)
}

// HandlePUT 监听Put
func (me *WebServer) HandlePUT(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(PUT, relativePath, handle)
}

// HandleDELETE 监听Delete
func (me *WebServer) HandleDELETE(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(DELETE, relativePath, handle)
}

// HandlePATCH 监听Patch
func (me *WebServer) HandlePATCH(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(PATCH, relativePath, handle)
}

// HandleOPTIONS 监听Options
func (me *WebServer) HandleOPTIONS(relativePath string, handle func(Context)) *WebServer {
	return me.HandleFunc(OPTIONS, relativePath, handle)
}

/*
HandleStruct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
结构体的方法必须与 Method 类型的名称一致
*/
func (me *WebServer) HandleStruct(relativePath string, handle interface{}) *WebServer {
	me.handleStruct(relativePath, handle)
	return me
}

// handleStruct 使用反射遍历结构体的方法和字段，对http方法进行注册
func (me *WebServer) handleStruct(relativePath string, handle interface{}) {
	rfType := reflect.TypeOf(handle)
	rfValue := reflect.ValueOf(handle)
	// 只接受结构体、接口及指针
	switch rfType.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Struct:
		// 反射方法
		for i := 0; i <= rfType.NumMethod()-1; i++ {
			methodName := rfType.Method(i).Name
			if !isMethod(NewMethod(methodName)) {
				continue
			}
			handleFunc := func(ctx Context) {
				rfValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
			}
			me.HandleFunc(NewMethod(methodName), relativePath, handleFunc)
		}
		// 反射字段
		for i := 0; i <= rfType.NumField()-1; i++ {
			fieldName := rfType.Field(i).Name
			me.handleStruct(relativePath+"/"+strings.ToLower(fieldName[:1])+fieldName[1:], rfValue.FieldByName(fieldName).Interface())
		}
	}
}
