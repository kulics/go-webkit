package go_webkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"path/filepath"
)

// Form 表单类型
type Form map[string]interface{}
type responseHandle = func(resp http.Response) error

const (
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJSON = "application/json"
)

// WebClient web请求辅助工具
type WebClient struct {
	host  string
	token string
	cli   *http.Client
}

// NewWebClient WebClient构造函数
func NewWebClient(host string) *WebClient {
	jar, _ := cookiejar.New(nil)
	return &WebClient{host, "", &http.Client{Jar: jar}}
}

// SetToken 设置token
func (sf *WebClient) SetToken(token string) {
	sf.token = token
}

// Cookies 根据域名获取cookies
func (sf *WebClient) Cookies(rawURL string) ([]*http.Cookie, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return sf.cli.Jar.Cookies(u), nil
}

// SetCookie 根据域名设置单条cookie
func (sf *WebClient) SetCookie(rawURL string, name string, value string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	sf.cli.Jar.SetCookies(u,
		[]*http.Cookie{&http.Cookie{Name: name, Value: value, HttpOnly: true}})
	return nil
}

// SetCookies 根据域名设置cookies
func (sf *WebClient) SetCookies(rawURL string, cookies []*http.Cookie) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	sf.cli.Jar.SetCookies(u, cookies)
	return nil
}

// HTTPRequest http请求
func (sf *WebClient) HTTPRequest(method Method, relativePath string,
	contentType string, params io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method.String(),
		sf.host+relativePath, params)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	if sf.token != "" {
		req.Header.Set("X-Access-Token", sf.token)
	}

	// req.Header.Set("Cookie", "name=anny")
	return sf.cli.Do(req)
}

func (sf *WebClient) processRequest(method Method, relativePath string,
	contentType string, params io.Reader, handles ...responseHandle) ([]byte, error) {
	resp, err := sf.HTTPRequest(method, relativePath, contentType, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	for _, v := range handles {
		err = v(*resp)
		if err != nil {
			return nil, err
		}
	}
	return ioutil.ReadAll(resp.Body)
}

// FormRequest 表单请求
func (sf *WebClient) FormRequest(method Method, relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	forms := url.Values{}
	for k, v := range params {
		forms.Add(k, fmt.Sprint(v))
	}
	return sf.processRequest(method, relativePath, ContentTypeForm,
		strings.NewReader(forms.Encode()), handles...)
}

// FormGET 表单get
func (sf *WebClient) FormGET(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return sf.FormRequest(GET, relativePath, params, handles...)
}

// FormPOST 表单post
func (sf *WebClient) FormPOST(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return sf.FormRequest(POST, relativePath, params, handles...)
}

// FormPUT 表单put
func (sf *WebClient) FormPUT(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return sf.FormRequest(PUT, relativePath, params, handles...)
}

// FormDELETE 表单delete
func (sf *WebClient) FormDELETE(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return sf.FormRequest(DELETE, relativePath, params, handles...)
}

// JSONRequest JSON请求
func (sf *WebClient) JSONRequest(method Method, relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	bt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body, err := sf.processRequest(method, relativePath, ContentTypeJSON, bytes.NewReader(bt), handles...)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, response)
}

// JSONGET JSON get
func (sf *WebClient) JSONGET(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return sf.JSONRequest(GET, relativePath, params, response, handles...)
}

// JSONPOST JSON post
func (sf *WebClient) JSONPOST(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return sf.JSONRequest(POST, relativePath, params, response, handles...)
}

// JSONPUT JSON put
func (sf *WebClient) JSONPUT(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return sf.JSONRequest(PUT, relativePath, params, response, handles...)
}

// JSONDELETE JSON delete
func (sf *WebClient) JSONDELETE(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return sf.JSONRequest(DELETE, relativePath, params, response, handles...)
}

// FileInfo 发送文件类型
type FileInfo struct {
	Field  string
	Path   string
	Params map[string]string
}

// FileUpload 文件上传方法
func (sf *WebClient) FileUpload(relativePath string, field string,
	path string, params map[string]string,
	handles ...responseHandle) ([]byte, error) {
	fileBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(fileBuffer)
	fileWriter, err := bodyWriter.CreateFormFile(field, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, err
	}
	for k, v := range params {
		err = bodyWriter.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}
	err = bodyWriter.Close()
	if err != nil {
		return nil,err
	}

	return sf.processRequest(POST, relativePath, bodyWriter.FormDataContentType(),
		fileBuffer, handles...)
}
