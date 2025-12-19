package polymarket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Request 发送HTTP请求
func (c *HTTPClient) Request(method, path string, headers map[string]string, body interface{}) (interface{}, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		if bodyStr, ok := body.(string); ok {
			// 预序列化的body（用于HMAC签名，保持一致性）
			reqBody = bytes.NewBufferString(bodyStr)
		} else {
			// JSON序列化（紧凑格式，无空格）
			// 参考: https://github.com/Polymarket/py-clob-client/issues/164
			// API要求JSON必须去掉所有空格，否则会返回401错误
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置默认headers
	req.Header.Set("User-Agent", "polymarket-sdk-go")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	// 注意：不手动设置 Accept-Encoding，让 Go http.Client 自动处理 gzip

	// 设置自定义headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// 尝试解析JSON
	var jsonData interface{}
	if err := json.Unmarshal(respBody, &jsonData); err != nil {
		// 如果不是JSON，返回原始字符串
		return string(respBody), nil
	}

	return jsonData, nil
}

// Get 发送GET请求
func (c *HTTPClient) Get(path string, headers map[string]string) (interface{}, error) {
	return c.Request("GET", path, headers, nil)
}

// Post 发送POST请求
func (c *HTTPClient) Post(path string, headers map[string]string, body interface{}) (interface{}, error) {
	return c.Request("POST", path, headers, body)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(path string, headers map[string]string, body interface{}) (interface{}, error) {
	return c.Request("DELETE", path, headers, body)
}

// Put 发送PUT请求
func (c *HTTPClient) Put(path string, headers map[string]string, body interface{}) (interface{}, error) {
	return c.Request("PUT", path, headers, body)
}
