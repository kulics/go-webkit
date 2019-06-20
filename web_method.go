package webkit

// Method http方法类型
type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	OPTIONS Method = "OPTIONS"
)

// New_Method Method构建函数
func New_Method(v string) Method {
	return Method(v)
}

func (sf Method) String() string {
	return string(sf)
}

// isMethod 判断是否存在的方法
func is_Method(m Method) (b bool) {
	switch m {
	case GET:
		fallthrough
	case POST:
		fallthrough
	case PUT:
		fallthrough
	case DELETE:
		fallthrough
	case PATCH:
		fallthrough
	case OPTIONS:
		b = true
	}
	return
}
