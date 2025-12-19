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

	// 创建客户端（需要L2认证才能查询余额）
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

	fmt.Println("=== Polymarket 钱包余额查询 ===")
	fmt.Printf("地址: %s\n", client.GetAddress())
	fmt.Printf("链ID: %d\n", chainID)
	fmt.Printf("签名类型: %d (0=EOA, 1=Magic/Email, 2=Browser)\n", signatureType)
	fmt.Println()

	// 检查是否需要创建或派生API凭证
	// 如果环境变量中已有API凭证，直接使用
	apiKey := os.Getenv("CLOB_API_KEY")
	apiSecret := os.Getenv("CLOB_SECRET")
	apiPassphrase := os.Getenv("CLOB_PASSPHRASE")

	var creds *polymarket.ApiCreds
	if apiKey != "" && apiSecret != "" && apiPassphrase != "" {
		fmt.Println("使用环境变量中的API凭证...")
		creds = &polymarket.ApiCreds{
			APIKey:        apiKey,
			APISecret:     apiSecret,
			APIPassphrase: apiPassphrase,
		}
		client.SetAPICreds(creds)
	} else {
		fmt.Println("未找到API凭证，正在创建或派生...")
		// 尝试派生已存在的API密钥（使用nonce=0）
		nonce := 0
		var err error
		creds, err = client.DeriveAPIKey(&nonce)
		if err != nil {
			// 如果派生失败，尝试创建新的
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
			// creds 已通过 DeriveAPIKey 自动设置到客户端
			// 为了代码清晰，我们显式使用 creds（虽然它已经被 SetAPICreds 使用）
			if creds == nil {
				log.Fatalf("派生API密钥失败: creds 为空")
			}
		}
	}

	// 查询抵押品余额
	fmt.Println("\n=== 查询抵押品余额 ===")
	collateralParams := &polymarket.BalanceAllowanceParams{
		AssetType: polymarket.AssetTypeCollateral,
		TokenID:   "", // 空字符串表示查询所有代币
	}

	collateralBalance, err := client.GetBalanceAllowance(collateralParams)
	if err != nil {
		log.Fatalf("查询抵押品余额失败: %v", err)
	}

	printBalance("抵押品 (USDC)", collateralBalance)

	// 查询条件代币余额（需要指定token_id）
	// 这里演示如何查询，实际使用时需要提供具体的token_id
	fmt.Println("\n=== 查询条件代币余额 ===")
	fmt.Println("提示: 需要提供具体的 token_id 才能查询条件代币余额")
	fmt.Println("示例: 从订单簿或市场信息中获取 token_id")

	// 如果提供了token_id环境变量，则查询
	tokenID := os.Getenv("TOKEN_ID")
	if tokenID != "" {
		conditionalParams := &polymarket.BalanceAllowanceParams{
			AssetType: polymarket.AssetTypeConditional,
			TokenID:   tokenID,
		}

		conditionalBalance, err := client.GetBalanceAllowance(conditionalParams)
		if err != nil {
			log.Printf("查询条件代币余额失败: %v", err)
		} else {
			printBalance(fmt.Sprintf("条件代币 (TokenID: %s)", tokenID), conditionalBalance)
		}
	} else {
		fmt.Println("未设置 TOKEN_ID 环境变量，跳过条件代币查询")
		fmt.Println("设置方式: export TOKEN_ID=\"your-token-id\"")
	}

	fmt.Println("\n=== 查询完成 ===")
}

// printBalance 格式化打印余额信息
func printBalance(title string, balance map[string]interface{}) {
	fmt.Printf("\n%s:\n", title)

	// 格式化JSON输出
	jsonData, err := json.MarshalIndent(balance, "  ", "  ")
	if err != nil {
		fmt.Printf("  原始数据: %+v\n", balance)
	} else {
		fmt.Println(string(jsonData))
	}

	// 尝试提取并转换余额（USDC使用6位小数）
	if balanceStr, ok := balance["balance"].(string); ok {
		balanceFormatted := formatUSDC(balanceStr)
		fmt.Printf("  余额: %s (原始值: %s)\n", balanceFormatted, balanceStr)
	}
	if allowanceStr, ok := balance["allowance"].(string); ok {
		allowanceFormatted := formatUSDC(allowanceStr)
		fmt.Printf("  授权: %s (原始值: %s)\n", allowanceFormatted, allowanceStr)
	}
}

// formatUSDC 格式化USDC余额（6位小数）
func formatUSDC(valueStr string) string {
	// 尝试解析为整数
	var value int64
	fmt.Sscanf(valueStr, "%d", &value)

	// 转换为USDC（6位小数）
	usdc := float64(value) / 1000000.0

	// 格式化输出
	if usdc >= 1000 {
		return fmt.Sprintf("%.2f USDC", usdc)
	} else if usdc >= 1 {
		return fmt.Sprintf("%.6f USDC", usdc)
	} else {
		return fmt.Sprintf("%.6f USDC", usdc)
	}
}
