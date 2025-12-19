package polymarket

import (
	"fmt"
	"strings"
)

// BuildQueryParams 构建查询参数
func BuildQueryParams(url, param, val string) string {
	urlWithParams := url
	last := urlWithParams[len(urlWithParams)-1:]
	if last == "?" {
		urlWithParams = fmt.Sprintf("%s%s=%s", urlWithParams, param, val)
	} else {
		urlWithParams = fmt.Sprintf("%s&%s=%s", urlWithParams, param, val)
	}
	return urlWithParams
}

// AddQueryTradeParams 添加交易查询参数
func AddQueryTradeParams(baseURL string, params *TradeParams, nextCursor string) string {
	if nextCursor == "" {
		nextCursor = "MA=="
	}

	url := baseURL
	if params != nil {
		url = url + "?"
		if params.Market != "" {
			url = BuildQueryParams(url, "market", params.Market)
		}
		if params.AssetID != "" {
			url = BuildQueryParams(url, "asset_id", params.AssetID)
		}
		if params.After > 0 {
			url = BuildQueryParams(url, "after", fmt.Sprintf("%d", params.After))
		}
		if params.Before > 0 {
			url = BuildQueryParams(url, "before", fmt.Sprintf("%d", params.Before))
		}
		if params.MakerAddress != "" {
			url = BuildQueryParams(url, "maker_address", params.MakerAddress)
		}
		if params.ID != "" {
			url = BuildQueryParams(url, "id", params.ID)
		}
		if nextCursor != "" {
			url = BuildQueryParams(url, "next_cursor", nextCursor)
		}
	}
	return url
}

// AddQueryOpenOrdersParams 添加开放订单查询参数
func AddQueryOpenOrdersParams(baseURL string, params *OpenOrderParams, nextCursor string) string {
	if nextCursor == "" {
		nextCursor = "MA=="
	}

	url := baseURL
	if params != nil {
		url = url + "?"
		if params.Market != "" {
			url = BuildQueryParams(url, "market", params.Market)
		}
		if params.AssetID != "" {
			url = BuildQueryParams(url, "asset_id", params.AssetID)
		}
		if params.ID != "" {
			url = BuildQueryParams(url, "id", params.ID)
		}
		if nextCursor != "" {
			url = BuildQueryParams(url, "next_cursor", nextCursor)
		}
	}
	return url
}

// DropNotificationsQueryParams 添加删除通知查询参数
func DropNotificationsQueryParams(baseURL string, params *DropNotificationParams) string {
	url := baseURL
	if params != nil && len(params.IDs) > 0 {
		url = url + "?"
		idsStr := strings.Join(params.IDs, ",")
		url = BuildQueryParams(url, "ids", idsStr)
	}
	return url
}

// AddBalanceAllowanceParamsToURL 添加余额和授权查询参数
func AddBalanceAllowanceParamsToURL(baseURL string, params *BalanceAllowanceParams) string {
	url := baseURL
	if params != nil {
		url = url + "?"
		if params.AssetType != "" {
			url = BuildQueryParams(url, "asset_type", string(params.AssetType))
		}
		if params.TokenID != "" {
			url = BuildQueryParams(url, "token_id", params.TokenID)
		}
		if params.SignatureType != nil && *params.SignatureType >= 0 {
			url = BuildQueryParams(url, "signature_type", fmt.Sprintf("%d", *params.SignatureType))
		}
	}
	return url
}

// AddOrderScoringParamsToURL 添加订单评分查询参数
func AddOrderScoringParamsToURL(baseURL string, params *OrderScoringParams) string {
	url := baseURL
	if params != nil {
		url = url + "?"
		if params.OrderID != "" {
			url = BuildQueryParams(url, "order_id", params.OrderID)
		}
	}
	return url
}

// AddOrdersScoringParamsToURL 添加多个订单评分查询参数
func AddOrdersScoringParamsToURL(baseURL string, params *OrdersScoringParams) string {
	url := baseURL
	if params != nil && len(params.OrderIDs) > 0 {
		url = url + "?"
		orderIDsStr := strings.Join(params.OrderIDs, ",")
		url = BuildQueryParams(url, "order_ids", orderIDsStr)
	}
	return url
}

