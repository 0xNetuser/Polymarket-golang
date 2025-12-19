package polymarket

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
)

// ParseRawOrderBookSummary 解析原始订单簿摘要
func ParseRawOrderBookSummary(rawObs map[string]interface{}) (*OrderBookSummary, error) {
	bids := []OrderSummary{}
	if bidsRaw, ok := rawObs["bids"].([]interface{}); ok {
		for _, bidRaw := range bidsRaw {
			if bid, ok := bidRaw.(map[string]interface{}); ok {
				bids = append(bids, OrderSummary{
					Price: fmt.Sprintf("%v", bid["price"]),
					Size:  fmt.Sprintf("%v", bid["size"]),
				})
			}
		}
	}

	asks := []OrderSummary{}
	if asksRaw, ok := rawObs["asks"].([]interface{}); ok {
		for _, askRaw := range asksRaw {
			if ask, ok := askRaw.(map[string]interface{}); ok {
				asks = append(asks, OrderSummary{
					Price: fmt.Sprintf("%v", ask["price"]),
					Size:  fmt.Sprintf("%v", ask["size"]),
				})
			}
		}
	}

	obs := &OrderBookSummary{
		Market:       getString(rawObs, "market"),
		AssetID:      getString(rawObs, "asset_id"),
		Timestamp:    getString(rawObs, "timestamp"),
		MinOrderSize: getString(rawObs, "min_order_size"),
		NegRisk:      getBool(rawObs, "neg_risk"),
		TickSize:     getString(rawObs, "tick_size"),
		Bids:         bids,
		Asks:         asks,
		Hash:         getString(rawObs, "hash"),
	}

	return obs, nil
}

// GenerateOrderBookSummaryHash 生成订单簿摘要哈希
func GenerateOrderBookSummaryHash(orderbook *OrderBookSummary) string {
	// 临时清空hash
	originalHash := orderbook.Hash
	orderbook.Hash = ""

	// 序列化为JSON
	jsonData, err := json.Marshal(orderbook)
	if err != nil {
		orderbook.Hash = originalHash
		return ""
	}

	// SHA1哈希
	hash := sha1.Sum(jsonData)
	hashStr := fmt.Sprintf("%x", hash)

	// 恢复hash
	orderbook.Hash = hashStr
	return hashStr
}

// OrderToJSON 将订单转换为JSON格式
func OrderToJSON(order *SignedOrder, owner string, orderType OrderType) map[string]interface{} {
	// 将SignedOrder转换为字典
	orderDict := map[string]interface{}{
		"maker":         order.Maker,
		"taker":         order.Taker,
		"tokenId":       order.TokenId,
		"makerAmount":   order.MakerAmount,
		"takerAmount":   order.TakerAmount,
		"side":          order.Side,
		"feeRateBps":    order.FeeRateBps,
		"nonce":         order.Nonce,
		"signer":        order.Signer,
		"expiration":    order.Expiration,
		"signatureType": order.SignatureType,
		"signature":     order.Signature,
	}
	return map[string]interface{}{
		"order":    orderDict,
		"owner":    owner,
		"orderType": string(orderType),
	}
}

// IsTickSizeSmaller 检查tick size是否更小
func IsTickSizeSmaller(a, b TickSize) bool {
	aFloat, _ := strconv.ParseFloat(string(a), 64)
	bFloat, _ := strconv.ParseFloat(string(b), 64)
	return aFloat < bFloat
}

// PriceValid 检查价格是否有效
func PriceValid(price float64, tickSize TickSize) bool {
	tickSizeFloat, _ := strconv.ParseFloat(string(tickSize), 64)
	return price >= tickSizeFloat && price <= 1.0-tickSizeFloat
}

// 辅助函数
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

