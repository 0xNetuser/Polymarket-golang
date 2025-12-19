package rfq

// RfqUserRequest RFQ用户请求
type RfqUserRequest struct {
	TokenID string  `json:"token_id"`
	Side    string  `json:"side"`    // BUY 或 SELL
	Size    float64 `json:"size"`
	Price   float64 `json:"price,omitempty"`
}

// RfqUserQuote RFQ用户报价
type RfqUserQuote struct {
	RequestID string  `json:"request_id"`
	TokenID   string  `json:"token_id"`
	Side      string  `json:"side"`    // BUY 或 SELL
	Size      float64 `json:"size"`
	Price     float64 `json:"price"`
}

// CancelRfqRequestParams 取消RFQ请求参数
type CancelRfqRequestParams struct {
	RequestID string `json:"request_id"`
}

// CancelRfqQuoteParams 取消RFQ报价参数
type CancelRfqQuoteParams struct {
	QuoteID string `json:"quote_id"`
}

// AcceptQuoteParams 接受报价参数
type AcceptQuoteParams struct {
	RequestID  string `json:"request_id"`
	QuoteID    string `json:"quote_id"`
	Expiration int    `json:"expiration"`
}

// ApproveOrderParams 批准订单参数
type ApproveOrderParams struct {
	RequestID  string `json:"request_id"`
	QuoteID    string `json:"quote_id"`
	Expiration int    `json:"expiration"`
}

// GetRfqRequestsParams 获取RFQ请求参数
type GetRfqRequestsParams struct {
	TokenID string `json:"token_id,omitempty"`
	Side    string `json:"side,omitempty"`
	Status  string `json:"status,omitempty"`
}

// GetRfqQuotesParams 获取RFQ报价参数
type GetRfqQuotesParams struct {
	RequestID string `json:"request_id,omitempty"`
	TokenID   string `json:"token_id,omitempty"`
	Side      string `json:"side,omitempty"`
	Status    string `json:"status,omitempty"`
}

// GetRfqBestQuoteParams 获取最佳RFQ报价参数
type GetRfqBestQuoteParams struct {
	TokenID string `json:"token_id"`
	Side    string `json:"side"` // BUY 或 SELL
	Size    float64 `json:"size"`
}

// MatchType 匹配类型
type MatchType string

const (
	MatchTypeFull          MatchType = "FULL"
	MatchTypePartial       MatchType = "PARTIAL"
	MatchTypeComplementary MatchType = "COMPLEMENTARY"
	MatchTypeMint          MatchType = "MINT"
	MatchTypeMerge         MatchType = "MERGE"
)

// RfqQuoteResponse RFQ报价响应
type RfqQuoteResponse struct {
	QuoteID   string    `json:"quoteId"`
	RequestID string    `json:"requestId"`
	Token     string    `json:"token"`
	Complement string   `json:"complement,omitempty"`
	Side      string    `json:"side"`
	Price     string    `json:"price"`
	SizeIn    string    `json:"sizeIn"`
	SizeOut   string    `json:"sizeOut"`
	MatchType MatchType `json:"matchType"`
	Status    string    `json:"status"`
}

// COLLATERAL_TOKEN_DECIMALS USDC小数位数
const COLLATERAL_TOKEN_DECIMALS = 6

