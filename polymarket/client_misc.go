package polymarket

import (
	"encoding/json"
	"fmt"
)

// CreateReadonlyAPIKey 创建只读API密钥
// 需要L2认证
func (c *ClobClient) CreateReadonlyAPIKey() (*ReadonlyApiKeyResponse, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	requestArgs := &RequestArgs{
		Method:      "POST",
		RequestPath: CreateReadonlyAPIKey,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(CreateReadonlyAPIKey, headers, nil)
	if err != nil {
		return nil, err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return &ReadonlyApiKeyResponse{
		APIKey: getStringFromMap(respMap, "apiKey"),
	}, nil
}

// GetReadonlyAPIKeys 获取只读API密钥列表
// 需要L2认证
func (c *ClobClient) GetReadonlyAPIKeys() (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: GetReadonlyAPIKeys,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Get(GetReadonlyAPIKeys, headers)
}

// DeleteReadonlyAPIKey 删除只读API密钥
// 需要L2认证
func (c *ClobClient) DeleteReadonlyAPIKey(key string) (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	body := map[string]string{"key": key}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	bodyStr := string(bodyJSON)

	requestArgs := &RequestArgs{
		Method:        "DELETE",
		RequestPath:   DeleteReadonlyAPIKey,
		Body:          body,
		SerializedBody: &bodyStr,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Delete(DeleteReadonlyAPIKey, headers, bodyStr)
}

// ValidateReadonlyAPIKey 验证只读API密钥
// 公开端点，不需要认证
func (c *ClobClient) ValidateReadonlyAPIKey(address, key string) (interface{}, error) {
	path := fmt.Sprintf("%s?address=%s&key=%s", ValidateReadonlyAPIKey, address, key)
	return c.httpClient.Get(path, nil)
}

// IsOrderScoring 检查订单是否正在评分
// 需要L2认证
func (c *ClobClient) IsOrderScoring(params *OrderScoringParams) (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	url := AddOrderScoringParamsToURL(c.host+IsOrderScoring, params)
	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: IsOrderScoring,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Get(url[len(c.host):], headers)
}

// AreOrdersScoring 检查多个订单是否正在评分
// 需要L2认证
func (c *ClobClient) AreOrdersScoring(params *OrdersScoringParams) (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	bodyJSON, err := json.Marshal(params.OrderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order IDs: %w", err)
	}
	bodyStr := string(bodyJSON)

	requestArgs := &RequestArgs{
		Method:        "POST",
		RequestPath:   AreOrdersScoring,
		Body:          params.OrderIDs,
		SerializedBody: &bodyStr,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Post(AreOrdersScoring, headers, bodyStr)
}

// GetMarkets 获取市场列表
func (c *ClobClient) GetMarkets(nextCursor string) (interface{}, error) {
	if nextCursor == "" {
		nextCursor = "MA=="
	}
	path := fmt.Sprintf("%s?next_cursor=%s", GetMarkets, nextCursor)
	return c.httpClient.Get(path, nil)
}

// GetSimplifiedMarkets 获取简化市场列表
func (c *ClobClient) GetSimplifiedMarkets(nextCursor string) (interface{}, error) {
	if nextCursor == "" {
		nextCursor = "MA=="
	}
	path := fmt.Sprintf("%s?next_cursor=%s", GetSimplifiedMarkets, nextCursor)
	return c.httpClient.Get(path, nil)
}

// GetSamplingMarkets 获取采样市场列表
func (c *ClobClient) GetSamplingMarkets(nextCursor string) (interface{}, error) {
	if nextCursor == "" {
		nextCursor = "MA=="
	}
	path := fmt.Sprintf("%s?next_cursor=%s", GetSamplingMarkets, nextCursor)
	return c.httpClient.Get(path, nil)
}

// GetSamplingSimplifiedMarkets 获取采样简化市场列表
func (c *ClobClient) GetSamplingSimplifiedMarkets(nextCursor string) (interface{}, error) {
	if nextCursor == "" {
		nextCursor = "MA=="
	}
	path := fmt.Sprintf("%s?next_cursor=%s", GetSamplingSimplifiedMarkets, nextCursor)
	return c.httpClient.Get(path, nil)
}

// GetMarket 根据condition_id获取市场
func (c *ClobClient) GetMarket(conditionID string) (interface{}, error) {
	path := GetMarket + conditionID
	return c.httpClient.Get(path, nil)
}

// GetMarketTradesEvents 根据condition_id获取市场交易事件
func (c *ClobClient) GetMarketTradesEvents(conditionID string) (interface{}, error) {
	path := GetMarketTradesEvents + conditionID
	return c.httpClient.Get(path, nil)
}

// UpdateBalanceAllowance 更新余额和授权
// 需要L2认证
func (c *ClobClient) UpdateBalanceAllowance(params *BalanceAllowanceParams) (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	// 如果signature_type未设置，使用builder的签名类型
	if params.SignatureType == nil || (params.SignatureType != nil && *params.SignatureType < 0) {
		if c.builder != nil {
			sigType := c.builder.GetSigType()
			params.SignatureType = &sigType
		} else {
			// 默认使用0（EOA）
			defaultSigType := 0
			params.SignatureType = &defaultSigType
		}
	}

	url := AddBalanceAllowanceParamsToURL(c.host+UpdateBalanceAllowance, params)
	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: UpdateBalanceAllowance,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Get(url[len(c.host):], headers)
}

// GetOrderBookHash 获取订单簿哈希
func (c *ClobClient) GetOrderBookHash(orderbook *OrderBookSummary) string {
	return GenerateOrderBookSummaryHash(orderbook)
}

// GetBuilderTrades 获取Builder交易记录
// 需要Builder认证
func (c *ClobClient) GetBuilderTrades(params *TradeParams, nextCursor string) ([]interface{}, error) {
	// TODO: 实现Builder认证检查
	// 目前使用L2认证作为替代
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	if nextCursor == "" {
		nextCursor = "MA=="
	}

	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: GetBuilderTrades,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	var results []interface{}
	for nextCursor != EndCursor {
		url := AddQueryTradeParams(c.host+GetBuilderTrades, params, nextCursor)
		resp, err := c.httpClient.Get(url[len(c.host):], headers)
		if err != nil {
			return nil, err
		}

		respMap, ok := resp.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid response format")
		}

		if cursor, ok := respMap["next_cursor"].(string); ok {
			nextCursor = cursor
		} else {
			nextCursor = EndCursor
		}

		if data, ok := respMap["data"].([]interface{}); ok {
			results = append(results, data...)
		}
	}

	return results, nil
}

