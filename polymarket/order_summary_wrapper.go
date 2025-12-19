package polymarket

// OrderSummaryWrapper 包装OrderSummary以实现order_builder.OrderSummary接口
type OrderSummaryWrapper struct {
	OrderSummary
}

// GetPrice 获取价格
func (w *OrderSummaryWrapper) GetPrice() string {
	return w.OrderSummary.Price
}

// GetSize 获取数量
func (w *OrderSummaryWrapper) GetSize() string {
	return w.OrderSummary.Size
}

// convertOrderSummaries 转换OrderSummary为order_builder.OrderSummary接口
func convertOrderSummaries(summaries []OrderSummary) []interface{} {
	result := make([]interface{}, len(summaries))
	for i, s := range summaries {
		result[i] = &OrderSummaryWrapper{OrderSummary: s}
	}
	return result
}

