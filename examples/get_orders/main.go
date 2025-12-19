package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/0xNetuser/Polymarket-golang/polymarket"
)

func main() {
	// 从环境变量获取配置
	host := os.Getenv("CLOB_HOST")
	if host == "" {
		host = "https://clob.polymarket.com" // 默认值
	}

	chainIDStr := os.Getenv("CHAIN_ID")
	chainID := 137 // 默认Polygon主网
	if chainIDStr != "" {
		fmt.Sscanf(chainIDStr, "%d", &chainID)
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatalf("错误: 必须设置 PRIVATE_KEY 环境变量")
	}

	funder := os.Getenv("FUNDER") // 可选，用于代理钱包

	// 获取签名类型（Magic钱包使用1，EOA使用0）
	signatureTypeStr := os.Getenv("SIGNATURE_TYPE")
	signatureType := 0 // 默认EOA
	if signatureTypeStr != "" {
		fmt.Sscanf(signatureTypeStr, "%d", &signatureType)
	}
	sigTypePtr := &signatureType

	// 创建客户端
	client, err := polymarket.NewClobClient(
		host,
		chainID,
		privateKey,
		nil,        // 初始时没有API凭证
		sigTypePtr, // 签名类型（0=EOA, 1=Magic/Email, 2=Browser proxy）
		funder,
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("=== Polymarket 获取订单示例 ===")
	fmt.Printf("地址: %s\n", client.GetAddress())
	fmt.Printf("链ID: %d\n", chainID)
	fmt.Printf("签名类型: %d (0=EOA, 1=Magic/Email, 2=Browser)\n", signatureType)
	fmt.Println()

	// 检查是否需要创建或派生API凭证
	apiKey := os.Getenv("CLOB_API_KEY")
	apiSecret := os.Getenv("CLOB_SECRET")
	apiPassphrase := os.Getenv("CLOB_PASSPHRASE")

	if apiKey != "" && apiSecret != "" && apiPassphrase != "" {
		fmt.Println("使用环境变量中的API凭证...")
		creds := &polymarket.ApiCreds{
			APIKey:        apiKey,
			APISecret:     apiSecret,
			APIPassphrase: apiPassphrase,
		}
		client.SetAPICreds(creds)
	} else {
		fmt.Println("未找到API凭证，正在创建或派生...")
		nonce := 0
		creds, err := client.DeriveAPIKey(&nonce)
		if err != nil {
			fmt.Println("派生失败，尝试创建新的API密钥...")
			creds, err = client.CreateAPIKey(&nonce)
			if err != nil {
				log.Fatalf("创建API密钥失败: %v", err)
			}
			fmt.Println("⚠️  新API密钥已创建，请保存以下凭证：")
			fmt.Printf("   API Key: %s\n", creds.APIKey)
			fmt.Printf("   Secret: %s\n", creds.APISecret)
			fmt.Printf("   Passphrase: %s\n", creds.APIPassphrase)
			fmt.Println()
		} else {
			fmt.Println("✓ 成功派生API密钥")
		}
	}

	// 获取所有订单（不带过滤条件）
	fmt.Println("\n=== 获取所有订单 ===")
	orders, err := client.GetOrders(nil, "")
	if err != nil {
		log.Fatalf("获取订单失败: %v", err)
	}

	if len(orders) == 0 {
		fmt.Println("没有找到任何订单")
	} else {
		fmt.Printf("找到 %d 个订单:\n\n", len(orders))
		for i, order := range orders {
			printOrder(i+1, order)
		}
	}

	// 如果指定了 MARKET 或 ASSET_ID，使用过滤条件获取订单
	market := os.Getenv("MARKET")
	assetID := os.Getenv("ASSET_ID")
	orderID := os.Getenv("ORDER_ID")

	if market != "" || assetID != "" || orderID != "" {
		fmt.Println("\n=== 使用过滤条件获取订单 ===")
		params := &polymarket.OpenOrderParams{
			ID:      orderID,
			Market:  market,
			AssetID: assetID,
		}
		fmt.Printf("过滤条件: ID=%s, Market=%s, AssetID=%s\n", orderID, market, assetID)

		filteredOrders, err := client.GetOrders(params, "")
		if err != nil {
			log.Fatalf("获取过滤订单失败: %v", err)
		}

		if len(filteredOrders) == 0 {
			fmt.Println("没有找到匹配的订单")
		} else {
			fmt.Printf("找到 %d 个匹配订单:\n\n", len(filteredOrders))
			for i, order := range filteredOrders {
				printOrder(i+1, order)
			}
		}
	}

	// 如果指定了特定订单ID，获取单个订单详情
	singleOrderID := os.Getenv("SINGLE_ORDER_ID")
	if singleOrderID != "" {
		fmt.Println("\n=== 获取单个订单详情 ===")
		fmt.Printf("订单ID: %s\n", singleOrderID)

		order, err := client.GetOrder(singleOrderID)
		if err != nil {
			log.Fatalf("获取订单详情失败: %v", err)
		}

		fmt.Println("\n订单详情:")
		printJSON(order)
	}

	fmt.Println("\n=== 完成 ===")
}

// printOrder 打印订单信息
func printOrder(index int, order interface{}) {
	orderMap, ok := order.(map[string]interface{})
	if !ok {
		fmt.Printf("订单 %d: %+v\n", index, order)
		return
	}

	fmt.Printf("--- 订单 %d ---\n", index)
	
	// 打印关键字段
	if id, ok := orderMap["id"].(string); ok {
		fmt.Printf("  ID: %s\n", id)
	}
	if status, ok := orderMap["status"].(string); ok {
		fmt.Printf("  状态: %s\n", status)
	}
	if side, ok := orderMap["side"].(string); ok {
		fmt.Printf("  方向: %s\n", side)
	}
	if price, ok := orderMap["price"].(string); ok {
		fmt.Printf("  价格: %s\n", price)
	}
	if originalSize, ok := orderMap["original_size"].(string); ok {
		fmt.Printf("  原始数量: %s\n", originalSize)
	}
	if sizeMatched, ok := orderMap["size_matched"].(string); ok {
		fmt.Printf("  已成交数量: %s\n", sizeMatched)
	}
	if assetID, ok := orderMap["asset_id"].(string); ok {
		// 只显示前20个字符
		if len(assetID) > 20 {
			fmt.Printf("  资产ID: %s...\n", assetID[:20])
		} else {
			fmt.Printf("  资产ID: %s\n", assetID)
		}
	}
	if market, ok := orderMap["market"].(string); ok {
		fmt.Printf("  市场: %s\n", market)
	}
	if orderType, ok := orderMap["type"].(string); ok {
		fmt.Printf("  类型: %s\n", orderType)
	}
	if createdAt, ok := orderMap["created_at"].(string); ok {
		fmt.Printf("  创建时间: %s\n", createdAt)
	}
	
	fmt.Println()
}

// printJSON 打印JSON格式的数据
func printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("原始数据: %+v\n", data)
		return
	}
	fmt.Println(string(jsonData))
}

