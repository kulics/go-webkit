package webkit

// Method http方法类型
type Method string

const GET Method = "GET"
const POST Method = "POST"
const PUT Method = "PUT"
const DELETE Method = "DELETE"
const PATCH Method = "PATCH"
const OPTIONS Method = "OPTIONS"

// NewMethod Method构建函数
func NewMethod(v string) Method {
	return Method(v)
}
func (me Method) String() string {
	return string(me)
}

// isMethod 判断是否存在的方法
func isMethod(m Method) bool {
	b := false
	switch m {
	case GET, POST, PUT, DELETE, PATCH, OPTIONS:
		b = true
	}
	return b
}
