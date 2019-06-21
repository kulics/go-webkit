"webkit" {
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
}

# Context 类型别名 #
Context => ?gin.Context

# Web_Server web服务辅助工具 #
Web Server -> {
	listen:Str
	engine:?gin.Engine
}

# NewWebServerDefault 构建web服务数据对象,domain可以为空 #
New Web Server Default(listen:Str) ->(v:?Web Server) {
	<- (Web Server{listen=listen, engine=gin.Default()}?)
}

# Run 运行web服务 #
(me: ?Web Server) Run() -> (e:error) {
	<- (me.engine.Run(me.listen))
}

# Handle_Func 监听函数 #
(me: ?Web Server) Handle Func(method:Method, relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server) {
    method? GET {
        me.engine.GET(relativePath, handle)
    } POST {
        me.engine.POST(relativePath, handle)
    } PUT {
        me.engine.PUT(relativePath, handle)
    } DELETE {
        me.engine.DELETE(relativePath, handle)
    } PATCH {
        me.engine.PATCH(relativePath, handle)
    } OPTIONS {
		me.engine.OPTIONS(relativePath, handle)
	}
	<- (me)
}

# Handle_GET 监听Get #
(me: ?Web Server) Handle GET(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server) {
	<- (me.Handle Func(GET, relativePath, handle))
}

# Handle_POST 监听Post #
(me: ?Web_Server) Handle POST(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server)  {
	<- (me.Handle Func(POST, relativePath, handle))
}

# Handle_PUT 监听Put #
(me: ?Web Server) Handle PUT(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server)  {
	<- (me.Handle Func(PUT, relativePath, handle))
}

# Handle_DELETE 监听Delete #
(me: ?Web Server) Handle DELETE(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server)  {
	<- (me.Handle Func(DELETE, relativePath, handle))
}

# Handle_PATCH 监听Patch #
(me: ?Web Server) Handle PATCH(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server)  {
	<- (me.Handle Func(PATCH, relativePath, handle))
}

# Handle_OPTIONS 监听Options #
(me: ?Web Server) Handle OPTIONS(relativePath:Str, handle:(ctx:Context)->()) -> (v:?Web_Server)  {
	return me.Handle_Func(OPTIONS, relativePath, handle)
}

# Handle_Struct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
结构体的方法必须与 Method 类型的名称一致 #
(me: ?Web Server) Handle Struct(relativePath:Str, handle:Any) -> (v:?Web_Server) {
	me.handle struct(relativePath, handle)
	<- (me)
}

# handle_struct 使用反射遍历结构体的方法和字段，对http方法进行注册 #
(me: ?Web_Server) handle struct(relativePath:Str, handle:Any) -> () {
	rfType := reflect.TypeOf(handle)
	rfValue := reflect.ValueOf(handle)
	# 只接受结构体、接口及指针 #
	rfType.Kind() ? reflect.Ptr, reflect.Interface, reflect.Struct {
		# 反射方法 #
		[0 < rfType.NumMethod()] @ i {
			methodName := rfType.Method(i).Name
			? ~is Method(New Method(methodName)) {
				-> @
			}
			handleFunc := (ctx:Context) -> () {
				rfValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
			}
			me.Handle Func(New Method(methodName), relativePath, handleFunc)
		}
		# 反射字段 #
		[0 < rfType.NumField()] @ i {
			fieldName := rfType.Field(i).Name
			me.handle struct(relativePath+"/"+strings.ToLower(fieldName[<1])+fieldName[1<=],
				rfValue.FieldByName(fieldName).Interface())
		}
	}
}
