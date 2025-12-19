package polymarket

import (
	"github.com/0xNetuser/Polymarket-golang/polymarket/rfq"
)

// CreateRfqRequest 创建RFQ请求（便捷方法）
func (c *ClobClient) CreateRfqRequest(request *rfq.RfqUserRequest) (interface{}, error) {
	return c.rfq.CreateRfqRequest(request)
}

// CancelRfqRequest 取消RFQ请求（便捷方法）
func (c *ClobClient) CancelRfqRequest(params *rfq.CancelRfqRequestParams) (interface{}, error) {
	return c.rfq.CancelRfqRequest(params)
}

// GetRfqRequests 获取RFQ请求列表（便捷方法）
func (c *ClobClient) GetRfqRequests(params *rfq.GetRfqRequestsParams) (interface{}, error) {
	return c.rfq.GetRfqRequests(params)
}

// CreateRfqQuote 创建RFQ报价（便捷方法）
func (c *ClobClient) CreateRfqQuote(quote *rfq.RfqUserQuote) (interface{}, error) {
	return c.rfq.CreateRfqQuote(quote)
}

// CancelRfqQuote 取消RFQ报价（便捷方法）
func (c *ClobClient) CancelRfqQuote(params *rfq.CancelRfqQuoteParams) (interface{}, error) {
	return c.rfq.CancelRfqQuote(params)
}

// GetRfqQuotes 获取RFQ报价列表（便捷方法）
func (c *ClobClient) GetRfqQuotes(params *rfq.GetRfqQuotesParams) (interface{}, error) {
	return c.rfq.GetRfqQuotes(params)
}

// GetRfqBestQuote 获取最佳RFQ报价（便捷方法）
func (c *ClobClient) GetRfqBestQuote(params *rfq.GetRfqBestQuoteParams) (interface{}, error) {
	return c.rfq.GetRfqBestQuote(params)
}

// AcceptRfqQuote 接受RFQ报价（便捷方法）
func (c *ClobClient) AcceptRfqQuote(params *rfq.AcceptQuoteParams) (interface{}, error) {
	return c.rfq.AcceptQuote(params)
}

// ApproveRfqOrder 批准RFQ订单（便捷方法）
func (c *ClobClient) ApproveRfqOrder(params *rfq.ApproveOrderParams) (interface{}, error) {
	return c.rfq.ApproveOrder(params)
}

// GetRfqConfig 获取RFQ配置（便捷方法）
func (c *ClobClient) GetRfqConfig() (interface{}, error) {
	return c.rfq.GetRfqConfig()
}

