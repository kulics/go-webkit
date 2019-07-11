package webkit

type Method string

const GET Method = "GET"
const POST Method = "POST"
const PUT Method = "PUT"
const DELETE Method = "DELETE"
const PATCH Method = "PATCH"
const OPTIONS Method = "OPTIONS"

func NewMethod(v string) (r Method) {
	return Method(v)
}
func (me Method) String() (r string) {
	return string(me)
}
func isMethod(m Method) (r bool) {
	b := false
	switch m {
	case GET:
		{
			b = true
		}
	case POST:
		{
			b = true
		}
	case PUT:
		{
			b = true
		}
	case DELETE:
		{
			b = true
		}
	case PATCH:
		{
			b = true
		}
	case OPTIONS:
		{
			b = true
		}

	}
	return b
}
