package webkit

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
	"path/filepath"
	"strings"
)

// Form 表单类型
type Form map[string]interface{}
type responseHandle = func(resp http.Response) error

const (
	ContentType_Form = "application/x-www-form-urlencoded"
	ContentType_JSON = "application/json"
)

// WebClient web请求辅助工具
type WebClient struct {
	host    string
	headers map[string]string
	cli     *http.Client
}

// NewWebClient WebClient构造函数
func NewWebClient(host string) *WebClient {
	jar, _ := cookiejar.New(nil)
	return &WebClient{host, make(map[string]string), &http.Client{Jar: jar}}
}

// Set_Token 设置token
func (me *WebClient) SetHeader(key string, value string) {
	me.headers[key] = value
}

// GetCookies 根据域名获取cookies
func (me *WebClient) GetCookies(rawURL string) ([]*http.Cookie, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return me.cli.Jar.Cookies(u), nil
}

// SetCookie 根据域名设置单条cookie
func (me *WebClient) SetCookie(rawURL string, name string, value string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	me.cli.Jar.SetCookies(u,
		[]*http.Cookie{&http.Cookie{Name: name, Value: value, HttpOnly: true}})
	return nil
}

// SetCookies 根据域名设置cookies
func (me *WebClient) SetCookies(rawURL string, cookies []*http.Cookie) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	me.cli.Jar.SetCookies(u, cookies)
	return nil
}

// HTTP_request http请求
func (me *WebClient) HTTP_request(method Method, relativePath string,
	contentType string, params io.Reader,
	header map[string]interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method.String(),
		me.host+relativePath, params)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	for k, v := range me.headers {
		req.Header.Set(k, v)
	}

	for k, v := range header {
		req.Header.Set(k, fmt.Sprint(v))
	}

	return me.cli.Do(req)
}

func (me *WebClient) processRequest(method Method, relativePath string,
	contentType string, params io.Reader, handles ...responseHandle) ([]byte, error) {
	resp, err := me.HTTP_request(method, relativePath, contentType, params, nil)
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

// Form_request 表单请求
func (me *WebClient) Form_request(method Method, relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	forms := url.Values{}
	for k, v := range params {
		forms.Add(k, fmt.Sprint(v))
	}
	return me.processRequest(method, relativePath, ContentType_Form,
		strings.NewReader(forms.Encode()), handles...)
}

// Form_GET 表单get
func (me *WebClient) Form_GET(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return me.Form_request(GET, relativePath, params, handles...)
}

// Form_POST 表单post
func (me *WebClient) Form_POST(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return me.Form_request(POST, relativePath, params, handles...)
}

// Form_PUT 表单put
func (me *WebClient) Form_PUT(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return me.Form_request(PUT, relativePath, params, handles...)
}

// Form_DELETE 表单delete
func (me *WebClient) Form_DELETE(relativePath string, params Form,
	handles ...responseHandle) ([]byte, error) {
	return me.Form_request(DELETE, relativePath, params, handles...)
}

// JSON_request JSON请求
func (me *WebClient) JSON_request(method Method, relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	bt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body, err := me.processRequest(method, relativePath, ContentType_JSON, bytes.NewReader(bt), handles...)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, response)
}

// JSON_GET JSON get
func (me *WebClient) JSON_GET(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return me.JSON_request(GET, relativePath, params, response, handles...)
}

// JSON_POST JSON post
func (me *WebClient) JSON_POST(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return me.JSON_request(POST, relativePath, params, response, handles...)
}

// JSON_PUT JSON put
func (me *WebClient) JSON_PUT(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return me.JSON_request(PUT, relativePath, params, response, handles...)
}

// JSON_DELETE JSON delete
func (me *WebClient) JSON_DELETE(relativePath string,
	params interface{}, response interface{},
	handles ...responseHandle) error {
	return me.JSON_request(DELETE, relativePath, params, response, handles...)
}

// File_Info 发送文件类型
type File_Info struct {
	Field  string
	Path   string
	Params map[string]string
}

// UploadFile 文件上传方法
func (me *WebClient) UploadFile(relativePath string, field string,
	path string, params Form,
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
		err = bodyWriter.WriteField(k, fmt.Sprint(v))
		if err != nil {
			return nil, err
		}
	}
	err = bodyWriter.Close()
	if err != nil {
		return nil, err
	}

	return me.processRequest(POST, relativePath, bodyWriter.FormDataContentType(),
		fileBuffer, handles...)
}

// DownloadFile 文件下载方法
func (me *WebClient) DownloadFile(relativePath string, savePath string,
	params Form, handles ...responseHandle) error {
	body, err := me.Form_GET(relativePath, params)
	if err != nil {
		return err
	}
	// 创建文件夹
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return err
	}
	// 写入临时文件
	f, err := ioutil.TempFile(filepath.Dir(savePath), filepath.Base(savePath)+"_temp")
	if err != nil {
		return err
	}
	if _, err := f.Write(body); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	// 将临时文件重命名为目标文件
	return os.Rename(f.Name(), savePath)
}
