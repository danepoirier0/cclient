package cclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tls "github.com/refraction-networking/utls"
)

type CclientHelper struct {
	client http.Client
}

// 新建一个CclientHelper,如果proxyUrl为空就不使用代理
func NewHelper(clientHelloId tls.ClientHelloID, proxyUrl string) (*CclientHelper, error) {
	httpClient, err := NewClient(clientHelloId, proxyUrl)
	if err != nil {
		return nil, err
	}
	clientHelper := CclientHelper{
		client: httpClient,
	}

	return &clientHelper, nil
}

// 发送GET请求
//
// 返回的err如果不是nil,记得 defer resp.Body.Close()
func (ch *CclientHelper) Get(reqUrl string, headers, cookies map[string]string) (*http.Response, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if cookies == nil {
		cookies = map[string]string{}
	}

	// 1.构造GET请求，拼接Headers和Cookies参数
	getReq, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		getReq.Header.Add(k, v)
	}
	for k, v := range cookies {
		getReq.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	// 2.发送请求
	resp, err := ch.client.Do(getReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 发送Post请求
//
// 如果返回不是错误，记得defer response.Body.Close()
func (ch *CclientHelper) Post(reqUrl string, headers, cookies map[string]string, data map[string]interface{}) (*http.Response, error) {

	if headers == nil {
		headers = map[string]string{}
	}
	if cookies == nil {
		cookies = map[string]string{}
	}
	if data == nil {
		data = map[string]interface{}{}
	}

	// 1. 判断ContentType来使用不同方式构造请求数据
	contentType := headers["Content-Type"]
	bodyStr := ""
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// 表单提交
		for key, value := range data {
			bodyStr += key + "=" + fmt.Sprintf("%v", value) + "&"
		}

	} else {
		// 非表单提交的都认为是json提交
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		bodyStr = string(bodyBytes)
	}

	// 2. 构造POST请求体, 并添加请求头和Cookies
	postReq, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		postReq.Header.Add(key, value)
	}
	for k, v := range cookies {
		postReq.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	// 3. 发送请求
	resp, err := ch.client.Do(postReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
