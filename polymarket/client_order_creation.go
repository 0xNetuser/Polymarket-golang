package polymarket

import (
	"fmt"
	"strconv"

	"github.com/polymarket/go-order-utils/pkg/model"
	obuilder "github.com/0xNetuser/Polymarket-golang/polymarket/order_builder"
)

// ResolveTickSize 解析tick size
func (c *ClobClient) resolveTickSize(tokenID string, tickSize *TickSize) (TickSize, error) {
	minTickSize, err := c.GetTickSize(tokenID)
	if err != nil {
		return "", err
	}

	if tickSize != nil {
		if IsTickSizeSmaller(*tickSize, minTickSize) {
			return "", fmt.Errorf("invalid tick size (%s), minimum for the market is %s", *tickSize, minTickSize)
		}
		return *tickSize, nil
	}

	return minTickSize, nil
}

// ResolveFeeRate 解析手续费率
func (c *ClobClient) resolveFeeRate(tokenID string, userFeeRate int) (int, error) {
	marketFeeRateBps, err := c.GetFeeRateBps(tokenID)
	if err != nil {
		return 0, err
	}

	// 如果市场手续费率和用户提供的手续费率都不为零，验证它们是否匹配
	if marketFeeRateBps > 0 && userFeeRate > 0 && userFeeRate != marketFeeRateBps {
		return 0, fmt.Errorf("invalid user provided fee rate: (%d), fee rate for the market must be %d", userFeeRate, marketFeeRateBps)
	}

	return marketFeeRateBps, nil
}

// CreateOrder 创建并签名订单（限价订单）
// 需要L1认证
func (c *ClobClient) CreateOrder(orderArgs *OrderArgs, options *PartialCreateOrderOptions) (*SignedOrder, error) {
	if err := c.assertLevel1Auth(); err != nil {
		return nil, err
	}

	// 解析tick size
	var tickSizePtr *TickSize
	if options != nil && options.TickSize != nil {
		tickSizePtr = options.TickSize
	}
	tickSize, err := c.resolveTickSize(orderArgs.TokenID, tickSizePtr)
	if err != nil {
		return nil, err
	}

	// 验证价格
	if !PriceValid(orderArgs.Price, tickSize) {
		tickSizeFloat, _ := strconv.ParseFloat(string(tickSize), 64)
		return nil, fmt.Errorf("price (%.6f), min: %s - max: %.6f", orderArgs.Price, tickSize, 1.0-tickSizeFloat)
	}

	// 解析neg risk
	negRisk := false
	if options != nil && options.NegRisk != nil {
		negRisk = *options.NegRisk
	} else {
		negRisk, err = c.GetNegRisk(orderArgs.TokenID)
		if err != nil {
			return nil, err
		}
	}

	// 解析手续费率
	feeRateBps, err := c.resolveFeeRate(orderArgs.TokenID, orderArgs.FeeRateBps)
	if err != nil {
		return nil, err
	}
	orderArgs.FeeRateBps = feeRateBps

	// 获取舍入配置
	roundConfig, ok := obuilder.RoundingConfig[string(tickSize)]
	if !ok {
		return nil, fmt.Errorf("unsupported tick size: %s", tickSize)
	}

	// 获取订单金额
	side, makerAmount, takerAmount, err := c.builder.GetOrderAmounts(
		orderArgs.Side,
		orderArgs.Size,
		orderArgs.Price,
		roundConfig,
	)
	if err != nil {
		return nil, err
	}

	// 构建OrderData
	taker := orderArgs.Taker
	if taker == "" {
		taker = ZeroAddress
	}

	orderData := &model.OrderData{
		Maker:         c.builder.GetFunder(),
		Taker:         taker,
		TokenId:       orderArgs.TokenID,
		MakerAmount:   makerAmount.String(),
		TakerAmount:   takerAmount.String(),
		Side:          side,
		FeeRateBps:    strconv.Itoa(orderArgs.FeeRateBps),
		Nonce:         strconv.Itoa(orderArgs.Nonce),
		Signer:        c.signer.Address(),
		Expiration:    strconv.Itoa(orderArgs.Expiration),
		SignatureType: model.SignatureType(c.builder.GetSigType()),
	}

	// 获取合约配置
	contractConfig := getContractConfig(c.chainID, negRisk)

	// 构建并签名订单
	signedOrder, err := c.builder.BuildSignedOrder(orderData, contractConfig.Exchange, c.chainID, negRisk)
	if err != nil {
		return nil, err
	}

	return signedOrder, nil
}

// CreateMarketOrder 创建并签名市价订单
// 需要L1认证
func (c *ClobClient) CreateMarketOrder(orderArgs *MarketOrderArgs, options *PartialCreateOrderOptions) (*SignedOrder, error) {
	if err := c.assertLevel1Auth(); err != nil {
		return nil, err
	}

	// 解析tick size
	var tickSizePtr *TickSize
	if options != nil && options.TickSize != nil {
		tickSizePtr = options.TickSize
	}
	tickSize, err := c.resolveTickSize(orderArgs.TokenID, tickSizePtr)
	if err != nil {
		return nil, err
	}

	// 如果价格未设置或为0，计算市价
	if orderArgs.Price <= 0 {
		price, err := c.CalculateMarketPrice(orderArgs.TokenID, orderArgs.Side, orderArgs.Amount, orderArgs.OrderType)
		if err != nil {
			return nil, err
		}
		orderArgs.Price = price
	}

	// 验证价格
	if !PriceValid(orderArgs.Price, tickSize) {
		tickSizeFloat, _ := strconv.ParseFloat(string(tickSize), 64)
		return nil, fmt.Errorf("price (%.6f), min: %s - max: %.6f", orderArgs.Price, tickSize, 1.0-tickSizeFloat)
	}

	// 解析neg risk
	negRisk := false
	if options != nil && options.NegRisk != nil {
		negRisk = *options.NegRisk
	} else {
		negRisk, err = c.GetNegRisk(orderArgs.TokenID)
		if err != nil {
			return nil, err
		}
	}

	// 解析手续费率
	feeRateBps, err := c.resolveFeeRate(orderArgs.TokenID, orderArgs.FeeRateBps)
	if err != nil {
		return nil, err
	}
	orderArgs.FeeRateBps = feeRateBps

	// 获取舍入配置
	roundConfig, ok := obuilder.RoundingConfig[string(tickSize)]
	if !ok {
		return nil, fmt.Errorf("unsupported tick size: %s", tickSize)
	}

	// 获取订单金额
	side, makerAmount, takerAmount, err := c.builder.GetMarketOrderAmounts(
		orderArgs.Side,
		orderArgs.Amount,
		orderArgs.Price,
		roundConfig,
	)
	if err != nil {
		return nil, err
	}

	// 构建OrderData
	taker := orderArgs.Taker
	if taker == "" {
		taker = ZeroAddress
	}

	orderData := &model.OrderData{
		Maker:         c.builder.GetFunder(),
		Taker:         taker,
		TokenId:       orderArgs.TokenID,
		MakerAmount:   makerAmount.String(),
		TakerAmount:   takerAmount.String(),
		Side:          side,
		FeeRateBps:    strconv.Itoa(orderArgs.FeeRateBps),
		Nonce:         strconv.Itoa(orderArgs.Nonce),
		Signer:        c.signer.Address(),
		Expiration:    "0", // 市价订单无过期时间
		SignatureType: model.SignatureType(c.builder.GetSigType()),
	}

	// 获取合约配置
	contractConfig := getContractConfig(c.chainID, negRisk)

	// 构建并签名订单
	signedOrder, err := c.builder.BuildSignedOrder(orderData, contractConfig.Exchange, c.chainID, negRisk)
	if err != nil {
		return nil, err
	}

	return signedOrder, nil
}

// CreateAndPostOrder 创建并提交订单（便捷方法）
func (c *ClobClient) CreateAndPostOrder(orderArgs *OrderArgs, options *PartialCreateOrderOptions) (interface{}, error) {
	order, err := c.CreateOrder(orderArgs, options)
	if err != nil {
		return nil, err
	}
	return c.PostOrder(order, OrderTypeGTC)
}

// CalculateMarketPrice 计算市价
func (c *ClobClient) CalculateMarketPrice(tokenID, side string, amount float64, orderType OrderType) (float64, error) {
	book, err := c.GetOrderBook(tokenID)
	if err != nil {
		return 0, fmt.Errorf("no orderbook: %w", err)
	}

	if side == BUY {
		if len(book.Asks) == 0 {
			return 0, fmt.Errorf("no match")
		}
		return c.builder.CalculateBuyMarketPrice(ConvertOrderSummaries(book.Asks), amount, string(orderType))
	} else {
		if len(book.Bids) == 0 {
			return 0, fmt.Errorf("no match")
		}
		return c.builder.CalculateSellMarketPrice(ConvertOrderSummaries(book.Bids), amount, string(orderType))
	}
}

// ConvertOrderSummaries 转换OrderSummary为order_builder.OrderSummary接口（导出函数）
func ConvertOrderSummaries(summaries []OrderSummary) []interface{} {
	result := make([]interface{}, len(summaries))
	for i, s := range summaries {
		result[i] = &OrderSummaryWrapper{OrderSummary: s}
	}
	return result
}


