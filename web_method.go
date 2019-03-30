package go_webkit

// Method http方法类型
type Method string

const (
	Get     Method = "Get"
	Post    Method = "Post"
	Put     Method = "Put"
	Delete  Method = "Delete"
	Patch   Method = "Patch"
	Options Method = "Options"
)

// NewMethod Method构建函数
func NewMethod(v string) Method {
	return Method(v)
}

// isMethod 判断是否存在的方法
func isMethod(m Method) (b bool) {
	switch m {
	case Get:
		fallthrough
	case Post:
		fallthrough
	case Put:
		fallthrough
	case Delete:
		fallthrough
	case Patch:
		fallthrough
	case Options:
		b = true
	}
	return
}
