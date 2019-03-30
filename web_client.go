package go_webkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Form 表单类型
type Form map[string]interface{}

const (
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJSON = "application/json"
)

// WebClient web请求辅助工具
type WebClient struct {
	host string
	cli  *http.Client
}

// NewWebClient WebClient构造函数
func NewWebClient(domain string, port string) *WebClient {
	return &WebClient{fmt.Sprintf("%s:%s", domain, port), &http.Client{}}
}

// HTTPRequest http请求
func (sf *WebClient) HTTPRequest(method Method, relativePath string,
	contentType string, params io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method.String(), sf.host+relativePath, params)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	// req.Header.Set("Cookie", "name=anny")
	return sf.cli.Do(req)
}

// mapToParams 将map转换为参数
func mapToParams(params Form) (io.Reader, error) {
	data := make([]string, 0)
	for k, v := range params {
		if s, ok := v.(string); ok {
			data = append(data, fmt.Sprintf("%s=%v", k, s))
			continue
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = append(data, fmt.Sprintf("%s=%s", k, string(b)))
	}

	return strings.NewReader(strings.Join(data, "&")), nil
}

// FormRequest 表单请求
func (sf *WebClient) FormRequest(method Method, relativePath string, params Form) ([]byte, error) {
	reader, err := mapToParams(params)
	if err != nil {
		return nil, err
	}
	resp, err := sf.HTTPRequest(method, relativePath, ContentTypeForm, reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// FormGET 表单get
func (sf *WebClient) FormGET(relativePath string, params Form) ([]byte, error) {
	return sf.FormRequest(GET, relativePath, params)
}

// FormPOST 表单post
func (sf *WebClient) FormPOST(relativePath string, params Form) ([]byte, error) {
	return sf.FormRequest(POST, relativePath, params)
}

// FormPUT 表单put
func (sf *WebClient) FormPUT(relativePath string, params Form) ([]byte, error) {
	return sf.FormRequest(PUT, relativePath, params)
}

// FormDELETE 表单delete
func (sf *WebClient) FormDELETE(relativePath string, params Form) ([]byte, error) {
	return sf.FormRequest(DELETE, relativePath, params)
}

// JSONRequest JSON请求
func (sf *WebClient) JSONRequest(method Method, relativePath string,
	params interface{}, response interface{}) error {
	bt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	resp, err := sf.HTTPRequest(method, relativePath, ContentTypeJSON, bytes.NewReader(bt))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, response)
}

// JSONGET JSON get
func (sf *WebClient) JSONGET(relativePath string,
	params interface{}, response interface{}) error {
	return sf.JSONRequest(GET, relativePath, params, response)
}

// JSONPOST JSON post
func (sf *WebClient) JSONPOST(relativePath string,
	params interface{}, response interface{}) error {
	return sf.JSONRequest(POST, relativePath, params, response)
}

// JSONPUT JSON put
func (sf *WebClient) JSONPUT(relativePath string,
	params interface{}, response interface{}) error {
	return sf.JSONRequest(PUT, relativePath, params, response)
}

// JSONDELETE JSON delete
func (sf *WebClient) JSONDELETE(relativePath string,
	params interface{}, response interface{}) error {
	return sf.JSONRequest(DELETE, relativePath, params, response)
}
