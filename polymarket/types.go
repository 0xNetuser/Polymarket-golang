package polymarket

import (
	"github.com/polymarket/go-order-utils/pkg/model"
)

// ApiCreds API凭证
type ApiCreds struct {
	APIKey     string `json:"apiKey"`
	APISecret  string `json:"secret"`
	APIPassphrase string `json:"passphrase"`
}

// ReadonlyApiKeyResponse 只读API密钥响应
type ReadonlyApiKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// RequestArgs 请求参数
type RequestArgs struct {
	Method        string
	RequestPath   string
	Body          interface{}
	SerializedBody *string
}

// BookParams 订单簿参数
type BookParams struct {
	TokenID string `json:"token_id"`
	Side    string `json:"side,omitempty"`
}

// OrderArgs 限价订单参数
type OrderArgs struct {
	TokenID     string  `json:"token_id"`      // 条件代币资产ID
	Price       float64 `json:"price"`        // 订单价格
	Size        float64 `json:"size"`         // 条件代币数量
	Side        string  `json:"side"`         // BUY 或 SELL
	FeeRateBps  int     `json:"fee_rate_bps"` // 手续费率（基点）
	Nonce       int     `json:"nonce"`        // 用于链上取消的nonce
	Expiration  int     `json:"expiration"`    // 订单过期时间戳
	Taker       string  `json:"taker"`         // 订单接受者地址，零地址表示公开订单
}

// MarketOrderArgs 市价订单参数
type MarketOrderArgs struct {
	TokenID     string    `json:"token_id"`      // 条件代币资产ID
	Amount      float64   `json:"amount"`       // BUY: 美元金额, SELL: 份额数量
	Side        string    `json:"side"`          // BUY 或 SELL
	Price       float64   `json:"price"`        // 订单价格（可选）
	FeeRateBps  int       `json:"fee_rate_bps"` // 手续费率（基点）
	Nonce       int       `json:"nonce"`        // 用于链上取消的nonce
	Taker       string    `json:"taker"`         // 订单接受者地址
	OrderType   OrderType `json:"order_type"`   // 订单类型
}

// TradeParams 交易查询参数
type TradeParams struct {
	ID           string `json:"id,omitempty"`
	MakerAddress string `json:"maker_address,omitempty"`
	Market       string `json:"market,omitempty"`
	AssetID      string `json:"asset_id,omitempty"`
	Before       int    `json:"before,omitempty"`
	After        int    `json:"after,omitempty"`
}

// OpenOrderParams 开放订单查询参数
type OpenOrderParams struct {
	ID      string `json:"id,omitempty"`
	Market  string `json:"market,omitempty"`
	AssetID string `json:"asset_id,omitempty"`
}

// DropNotificationParams 删除通知参数
type DropNotificationParams struct {
	IDs []string `json:"ids,omitempty"`
}

// OrderSummary 订单摘要
type OrderSummary struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// OrderBookSummary 订单簿摘要
type OrderBookSummary struct {
	Market        string         `json:"market"`
	AssetID       string         `json:"asset_id"`
	Timestamp     string         `json:"timestamp"`
	Bids          []OrderSummary `json:"bids"`
	Asks          []OrderSummary `json:"asks"`
	MinOrderSize  string         `json:"min_order_size"`
	NegRisk       bool           `json:"neg_risk"`
	TickSize      string         `json:"tick_size"`
	Hash          string         `json:"hash"`
}

// AssetType 资产类型
type AssetType string

const (
	AssetTypeCollateral  AssetType = "COLLATERAL"  // 抵押品（如USDC）
	AssetTypeConditional AssetType = "CONDITIONAL" // 条件代币
)

// BalanceAllowanceParams 余额和授权查询参数
type BalanceAllowanceParams struct {
	AssetType     AssetType `json:"asset_type,omitempty"`
	TokenID       string    `json:"token_id,omitempty"`
	SignatureType *int      `json:"signature_type,omitempty"` // 指针类型，允许nil表示未设置
}

// BalanceAllowanceResponse 余额和授权响应
type BalanceAllowanceResponse struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}

// OrderScoringParams 订单评分参数
type OrderScoringParams struct {
	OrderID string `json:"order_id"`
}

// OrdersScoringParams 多个订单评分参数
type OrdersScoringParams struct {
	OrderIDs []string `json:"order_ids"`
}

// CreateOrderOptions 创建订单选项
type CreateOrderOptions struct {
	TickSize TickSize `json:"tick_size"`
	NegRisk  bool     `json:"neg_risk"`
}

// PartialCreateOrderOptions 部分创建订单选项
type PartialCreateOrderOptions struct {
	TickSize *TickSize `json:"tick_size,omitempty"`
	NegRisk  *bool     `json:"neg_risk,omitempty"`
	RawOrder bool      `json:"raw_order,omitempty"` // 跳过 tick_size 获取和价格舍入，直接使用原始值
}

// RoundConfig 舍入配置
type RoundConfig struct {
	Price  int // 价格小数位数
	Size   int // 数量小数位数
	Amount int // 金额小数位数
}

// ContractConfig 合约配置
type ContractConfig struct {
	Exchange         string `json:"exchange"`          // 交易所合约地址
	Collateral       string `json:"collateral"`         // 抵押品代币地址
	ConditionalTokens string `json:"conditional_tokens"` // 条件代币合约地址
}

// PostOrdersArgs 批量下单参数
type PostOrdersArgs struct {
	Order    *model.SignedOrder `json:"order"`
	OrderType OrderType         `json:"orderType"`
}

// SignedOrder 已签名的订单（包装go-order-utils的SignedOrder）
type SignedOrder = model.SignedOrder

