package go_webkit

import "fmt"

// WebClient web请求辅助工具
type WebClient struct {
	host string
}

// NewWebClient WebClient构造函数
func NewWebClient(domain string, port string) *WebClient {
	return &WebClient{fmt.Sprintf("%s:%s", domain, port)}
}

func (sf *WebClient) httpRequest() ([]byte, error) {
	return nil, nil
}
