package rfq

import (
	"fmt"
)

// HTTPClientInterface HTTP客户端接口
type HTTPClientInterface interface {
	Get(path string, headers map[string]string) (interface{}, error)
	Post(path string, headers map[string]string, body interface{}) (interface{}, error)
	Delete(path string, headers map[string]string, body interface{}) (interface{}, error)
}

// ClobClientInterface CLOB客户端接口（避免循环导入）
type ClobClientInterface interface {
	AssertLevel2Auth() error
	GetHTTPClient() HTTPClientInterface
	GetHost() string
	CreateLevel2HeadersInternal(method, path string, body interface{}) (map[string]string, error)
}

// RfqClient RFQ客户端
type RfqClient struct {
	parent ClobClientInterface
}

// NewRfqClient 创建新的RFQ客户端
func NewRfqClient(parent ClobClientInterface) *RfqClient {
	return &RfqClient{
		parent: parent,
	}
}

// ensureL2Auth 确保L2认证
func (r *RfqClient) ensureL2Auth() error {
	return r.parent.AssertLevel2Auth()
}

// getL2Headers 获取L2认证头
func (r *RfqClient) getL2Headers(method, endpoint string, body interface{}) (map[string]string, error) {
	return r.parent.CreateLevel2HeadersInternal(method, endpoint, body)
}

// CreateRfqRequest 创建RFQ请求
func (r *RfqClient) CreateRfqRequest(request *RfqUserRequest) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	headers, err := r.getL2Headers("POST", "/rfq/request", request)
	if err != nil {
		return nil, err
	}

	httpClient := r.parent.GetHTTPClient()
	return httpClient.Post("/rfq/request", headers, request)
}

// CancelRfqRequest 取消RFQ请求
func (r *RfqClient) CancelRfqRequest(params *CancelRfqRequestParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	headers, err := r.getL2Headers("DELETE", "/rfq/request", params)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Delete("/rfq/request", headers, params)
}

// GetRfqRequests 获取RFQ请求列表
func (r *RfqClient) GetRfqRequests(params *GetRfqRequestsParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	// 构建查询参数
	path := "/rfq/data/requests"
	if params != nil {
		path += "?"
		if params.TokenID != "" {
			path += fmt.Sprintf("token_id=%s&", params.TokenID)
		}
		if params.Side != "" {
			path += fmt.Sprintf("side=%s&", params.Side)
		}
		if params.Status != "" {
			path += fmt.Sprintf("status=%s&", params.Status)
		}
		// 移除末尾的&
		if len(path) > 0 && path[len(path)-1] == '&' {
			path = path[:len(path)-1]
		}
	}

	headers, err := r.getL2Headers("GET", "/rfq/data/requests", nil)
	if err != nil {
		return nil, err
	}

	httpClient := r.parent.GetHTTPClient()
	return httpClient.Get(path, headers)
}

// CreateRfqQuote 创建RFQ报价
func (r *RfqClient) CreateRfqQuote(quote *RfqUserQuote) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	headers, err := r.getL2Headers("POST", "/rfq/quote", quote)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Post("/rfq/quote", headers, quote)
}

// CancelRfqQuote 取消RFQ报价
func (r *RfqClient) CancelRfqQuote(params *CancelRfqQuoteParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	headers, err := r.getL2Headers("DELETE", "/rfq/quote", params)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Delete("/rfq/quote", headers, params)
}

// GetRfqQuotes 获取RFQ报价列表
func (r *RfqClient) GetRfqQuotes(params *GetRfqQuotesParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	// 构建查询参数
	path := "/rfq/data/quotes"
	if params != nil {
		path += "?"
		if params.RequestID != "" {
			path += fmt.Sprintf("request_id=%s&", params.RequestID)
		}
		if params.TokenID != "" {
			path += fmt.Sprintf("token_id=%s&", params.TokenID)
		}
		if params.Side != "" {
			path += fmt.Sprintf("side=%s&", params.Side)
		}
		if params.Status != "" {
			path += fmt.Sprintf("status=%s&", params.Status)
		}
		// 移除末尾的&
		if len(path) > 0 && path[len(path)-1] == '&' {
			path = path[:len(path)-1]
		}
	}

	headers, err := r.getL2Headers("GET", "/rfq/data/quotes", nil)
	if err != nil {
		return nil, err
	}

	httpClient := r.parent.GetHTTPClient()
	return httpClient.Get(path, headers)
}

// GetRfqBestQuote 获取最佳RFQ报价
func (r *RfqClient) GetRfqBestQuote(params *GetRfqBestQuoteParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/rfq/data/best-quote?token_id=%s&side=%s&size=%.6f", 
		params.TokenID, params.Side, params.Size)

	headers, err := r.getL2Headers("GET", "/rfq/data/best-quote", nil)
	if err != nil {
		return nil, err
	}

	httpClient := r.parent.GetHTTPClient()
	return httpClient.Get(path, headers)
}

// AcceptQuote 接受报价（请求方）
// 此方法会获取报价详情，创建签名订单，然后提交接受请求
func (r *RfqClient) AcceptQuote(params *AcceptQuoteParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	// 构建接受请求的body
	body := map[string]interface{}{
		"requestId": params.RequestID,
		"quoteId":   params.QuoteID,
	}

	headers, err := r.getL2Headers("POST", "/rfq/request/accept", body)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Post("/rfq/request/accept", headers, body)
}

// ApproveOrder 批准订单（报价方）
// 此方法会获取报价详情，创建签名订单，然后提交批准请求
func (r *RfqClient) ApproveOrder(params *ApproveOrderParams) (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	// 构建批准请求的body
	body := map[string]interface{}{
		"requestId": params.RequestID,
		"quoteId":   params.QuoteID,
	}

	headers, err := r.getL2Headers("POST", "/rfq/quote/approve", body)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Post("/rfq/quote/approve", headers, body)
}

// GetRfqConfig 获取RFQ配置
func (r *RfqClient) GetRfqConfig() (interface{}, error) {
	if err := r.ensureL2Auth(); err != nil {
		return nil, err
	}

	headers, err := r.getL2Headers("GET", "/rfq/config", nil)
	if err != nil {
		return nil, err
	}

	return r.parent.GetHTTPClient().Get("/rfq/config", headers)
}
