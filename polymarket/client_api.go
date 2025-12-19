package polymarket

import (
	"fmt"
)

// GetOK 健康检查：确认服务器是否运行
// 不需要认证
func (c *ClobClient) GetOK() (interface{}, error) {
	return c.httpClient.Get("/", nil)
}

// GetServerTime 返回服务器当前时间戳
// 不需要认证
func (c *ClobClient) GetServerTime() (interface{}, error) {
	return c.httpClient.Get(Time, nil)
}

// CreateAPIKey 创建新的CLOB API密钥
// 需要L1认证
func (c *ClobClient) CreateAPIKey(nonce *int) (*ApiCreds, error) {
	if err := c.assertLevel1Auth(); err != nil {
		return nil, err
	}

	headers, err := CreateLevel1Headers(c.signer, nonce)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(CreateAPIKey, headers, nil)
	if err != nil {
		return nil, err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	creds := &ApiCreds{
		APIKey:        getStringFromMap(respMap, "apiKey"),
		APISecret:     getStringFromMap(respMap, "secret"),
		APIPassphrase: getStringFromMap(respMap, "passphrase"),
	}

	// 自动设置到客户端
	c.SetAPICreds(creds)

	return creds, nil
}

// DeriveAPIKey 派生已存在的CLOB API密钥
// 需要L1认证
func (c *ClobClient) DeriveAPIKey(nonce *int) (*ApiCreds, error) {
	if err := c.assertLevel1Auth(); err != nil {
		return nil, err
	}

	headers, err := CreateLevel1Headers(c.signer, nonce)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Get(DeriveAPIKey, headers)
	if err != nil {
		return nil, err
	}

	// 处理数组响应（API 可能返回数组）
	if respArr, ok := resp.([]interface{}); ok {
		if len(respArr) > 0 {
			if respMap, ok := respArr[0].(map[string]interface{}); ok {
				resp = respMap
			}
		}
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	creds := &ApiCreds{
		APIKey:        getStringFromMap(respMap, "apiKey"),
		APISecret:     getStringFromMap(respMap, "secret"),
		APIPassphrase: getStringFromMap(respMap, "passphrase"),
	}

	// 自动设置到客户端
	c.SetAPICreds(creds)

	return creds, nil
}

// CreateOrDeriveAPIKey 创建或派生API凭证
// 先尝试创建，如果失败则派生
func (c *ClobClient) CreateOrDeriveAPIKey(nonce *int) (*ApiCreds, error) {
	creds, err := c.CreateAPIKey(nonce)
	if err != nil {
		// 如果创建失败，尝试派生
		return c.DeriveAPIKey(nonce)
	}
	return creds, nil
}

// GetAPIKeys 获取可用的API密钥列表
// 需要L2认证
func (c *ClobClient) GetAPIKeys() (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: GetAPIKeys,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Get(GetAPIKeys, headers)
}

// GetClosedOnlyMode 获取closed only模式标志
// 需要L2认证
func (c *ClobClient) GetClosedOnlyMode() (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	requestArgs := &RequestArgs{
		Method:      "GET",
		RequestPath: ClosedOnly,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Get(ClosedOnly, headers)
}

// DeleteAPIKey 删除API密钥
// 需要L2认证
func (c *ClobClient) DeleteAPIKey() (interface{}, error) {
	if err := c.assertLevel2Auth(); err != nil {
		return nil, err
	}

	requestArgs := &RequestArgs{
		Method:      "DELETE",
		RequestPath: DeleteAPIKey,
	}

	headers, err := CreateLevel2Headers(c.signer, c.creds, requestArgs)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Delete(DeleteAPIKey, headers, nil)
}

// GetMidpoint 获取中点价格
func (c *ClobClient) GetMidpoint(tokenID string) (interface{}, error) {
	path := fmt.Sprintf("%s?token_id=%s", MidPoint, tokenID)
	return c.httpClient.Get(path, nil)
}

// GetMidpoints 获取多个token的中点价格
func (c *ClobClient) GetMidpoints(params []BookParams) (interface{}, error) {
	body := make([]map[string]string, len(params))
	for i, p := range params {
		body[i] = map[string]string{"token_id": p.TokenID}
	}
	return c.httpClient.Post(MidPoints, nil, body)
}

// GetPrice 获取市场价格
func (c *ClobClient) GetPrice(tokenID, side string) (interface{}, error) {
	path := fmt.Sprintf("%s?token_id=%s&side=%s", Price, tokenID, side)
	return c.httpClient.Get(path, nil)
}

// GetPrices 获取多个token的市场价格
func (c *ClobClient) GetPrices(params []BookParams) (interface{}, error) {
	body := make([]map[string]string, len(params))
	for i, p := range params {
		body[i] = map[string]string{
			"token_id": p.TokenID,
			"side":     p.Side,
		}
	}
	return c.httpClient.Post(GetPrices, nil, body)
}

// GetSpread 获取价差
func (c *ClobClient) GetSpread(tokenID string) (interface{}, error) {
	path := fmt.Sprintf("%s?token_id=%s", GetSpread, tokenID)
	return c.httpClient.Get(path, nil)
}

// GetSpreads 获取多个token的价差
func (c *ClobClient) GetSpreads(params []BookParams) (interface{}, error) {
	body := make([]map[string]string, len(params))
	for i, p := range params {
		body[i] = map[string]string{"token_id": p.TokenID}
	}
	return c.httpClient.Post(GetSpreads, nil, body)
}

// GetTickSize 获取tick size（带缓存）
func (c *ClobClient) GetTickSize(tokenID string) (TickSize, error) {
	c.mu.RLock()
	if tickSize, ok := c.tickSizes[tokenID]; ok {
		c.mu.RUnlock()
		return tickSize, nil
	}
	c.mu.RUnlock()

	path := fmt.Sprintf("%s?token_id=%s", GetTickSize, tokenID)
	resp, err := c.httpClient.Get(path, nil)
	if err != nil {
		return "", err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	tickSizeStr := getStringFromMap(respMap, "minimum_tick_size")
	tickSize := TickSize(tickSizeStr)

	c.mu.Lock()
	c.tickSizes[tokenID] = tickSize
	c.mu.Unlock()

	return tickSize, nil
}

// GetNegRisk 获取neg risk标志（带缓存）
func (c *ClobClient) GetNegRisk(tokenID string) (bool, error) {
	c.mu.RLock()
	if negRisk, ok := c.negRisk[tokenID]; ok {
		c.mu.RUnlock()
		return negRisk, nil
	}
	c.mu.RUnlock()

	path := fmt.Sprintf("%s?token_id=%s", GetNegRisk, tokenID)
	resp, err := c.httpClient.Get(path, nil)
	if err != nil {
		return false, err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid response format")
	}

	negRisk := getBoolFromMap(respMap, "neg_risk")

	c.mu.Lock()
	c.negRisk[tokenID] = negRisk
	c.mu.Unlock()

	return negRisk, nil
}

// GetFeeRateBps 获取手续费率（基点）（带缓存）
func (c *ClobClient) GetFeeRateBps(tokenID string) (int, error) {
	c.mu.RLock()
	if feeRate, ok := c.feeRates[tokenID]; ok {
		c.mu.RUnlock()
		return feeRate, nil
	}
	c.mu.RUnlock()

	path := fmt.Sprintf("%s?token_id=%s", GetFeeRate, tokenID)
	resp, err := c.httpClient.Get(path, nil)
	if err != nil {
		return 0, err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	feeRate := 0
	if baseFee, ok := respMap["base_fee"]; ok {
		if feeRateFloat, ok := baseFee.(float64); ok {
			feeRate = int(feeRateFloat)
		}
	}

	c.mu.Lock()
	c.feeRates[tokenID] = feeRate
	c.mu.Unlock()

	return feeRate, nil
}

// GetOrderBook 获取订单簿
func (c *ClobClient) GetOrderBook(tokenID string) (*OrderBookSummary, error) {
	path := fmt.Sprintf("%s?token_id=%s", GetOrderBook, tokenID)
	resp, err := c.httpClient.Get(path, nil)
	if err != nil {
		return nil, err
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return ParseRawOrderBookSummary(respMap)
}

// GetOrderBooks 获取多个订单簿
func (c *ClobClient) GetOrderBooks(params []BookParams) ([]*OrderBookSummary, error) {
	body := make([]map[string]string, len(params))
	for i, p := range params {
		body[i] = map[string]string{"token_id": p.TokenID}
	}

	resp, err := c.httpClient.Post(GetOrderBooks, nil, body)
	if err != nil {
		return nil, err
	}

	respArray, ok := resp.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	orderBooks := make([]*OrderBookSummary, len(respArray))
	for i, item := range respArray {
		if itemMap, ok := item.(map[string]interface{}); ok {
			obs, err := ParseRawOrderBookSummary(itemMap)
			if err != nil {
				return nil, err
			}
			orderBooks[i] = obs
		}
	}

	return orderBooks, nil
}

// GetLastTradePrice 获取最后成交价格
func (c *ClobClient) GetLastTradePrice(tokenID string) (interface{}, error) {
	path := fmt.Sprintf("%s?token_id=%s", GetLastTradePrice, tokenID)
	return c.httpClient.Get(path, nil)
}

// GetLastTradesPrices 获取多个token的最后成交价格
func (c *ClobClient) GetLastTradesPrices(params []BookParams) (interface{}, error) {
	body := make([]map[string]string, len(params))
	for i, p := range params {
		body[i] = map[string]string{"token_id": p.TokenID}
	}
	return c.httpClient.Post(GetLastTradesPrices, nil, body)
}

// 辅助函数
func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getBoolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}
