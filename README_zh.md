# Polymarket Go SDK

[English](README.md) | [中文](README_zh.md)

Polymarket CLOB（中央限价订单簿）API 的 Go 语言 SDK，完整实现了 [py-clob-client](https://github.com/Polymarket/py-clob-client) 的所有核心功能。


Follow at X:  @netu5er


## 特性

- ✅ **完整的 API 覆盖**: 实现 py-clob-client 的所有端点
- ✅ **三种认证级别**: L0（只读）、L1（私钥）、L2（完整访问）
- ✅ **订单管理**: 创建、提交、取消和查询订单
- ✅ **市场数据**: 订单簿、价格、价差和市场信息
- ✅ **RFQ 支持**: 报价请求功能
- ✅ **类型安全**: 使用 Go 的类型系统提供强类型支持
- ✅ **EIP-712 签名**: 完整支持以太坊消息签名

## 安装

```bash
go get github.com/0xNetuser/Polymarket-golang
```

## 快速开始

### Level 0（只读模式）

```go
package main

import (
    "fmt"
    "github.com/0xNetuser/Polymarket-golang/polymarket"
)

func main() {
    // 创建只读客户端
    client, err := polymarket.NewClobClient(
        "https://clob.polymarket.com",
        137, // Polygon 链 ID
        "",  // 无私钥
        nil, // 无 API 凭证
        nil, // 无签名类型
        "",  // 无 funder 地址
    )
    if err != nil {
        panic(err)
    }

    // 健康检查
    ok, err := client.GetOK()
    fmt.Println("服务器状态:", ok)

    // 获取订单簿
    orderBook, err := client.GetOrderBook("token-id")
    if err != nil {
        panic(err)
    }
    fmt.Printf("订单簿: %+v\n", orderBook)
}
```

### Level 1（私钥认证）

```go
client, err := polymarket.NewClobClient(
    "https://clob.polymarket.com",
    137,
    "your-private-key-hex",
    nil,
    nil,
    "",
)

// 创建或派生 API 凭证
creds, err := client.CreateOrDeriveAPIKey(nil)
if err != nil {
    panic(err)
}
fmt.Printf("API Key: %s\n", creds.APIKey)
```

### Level 2（完整功能）

```go
creds := &polymarket.ApiCreds{
    APIKey:      "your-api-key",
    APISecret:   "your-api-secret",
    APIPassphrase: "your-passphrase",
}

client, err := polymarket.NewClobClient(
    "https://clob.polymarket.com",
    137,
    "your-private-key-hex",
    creds,
    nil,
    "",
)

// 查询余额
balance, err := client.GetBalanceAllowance(&polymarket.BalanceAllowanceParams{
    AssetType: polymarket.AssetTypeCollateral,
})
if err != nil {
    panic(err)
}
fmt.Printf("余额: %+v\n", balance)
```

## 示例

### 查询余额

```bash
cd examples/check_balance
export PRIVATE_KEY="your-private-key"
export CHAIN_ID="137"  # 可选，默认为 137 (Polygon)
export SIGNATURE_TYPE="0"  # 可选，0=EOA, 1=Magic/Email, 2=Browser
export FUNDER=""  # 可选，用于代理钱包
go run main.go
```

### 创建并提交订单

```go
// 创建订单
orderArgs := &polymarket.OrderArgs{
    TokenID:    "token-id",
    Price:      0.5,
    Size:       100.0,
    Side:       "BUY",
    FeeRateBps: 0,
    Nonce:      1,
    Expiration: 1234567890,
}

order, err := client.CreateOrder(orderArgs, nil)
if err != nil {
    panic(err)
}

// 提交订单
result, err := client.PostOrder(order, polymarket.OrderTypeGTC)
if err != nil {
    panic(err)
}
fmt.Printf("订单已提交: %+v\n", result)
```

## 项目结构

```
polymarket/
├── client.go                  # 主客户端结构
├── client_api.go              # API 方法（健康检查、API 密钥、市场数据等）
├── client_orders.go           # 订单管理方法（提交、取消、查询）
├── client_order_creation.go   # 订单创建方法（CreateOrder, CreateMarketOrder）
├── client_misc.go             # 其他功能（只读 API 密钥、订单评分、市场查询等）
├── rfq_client.go              # RFQ 客户端便捷方法
├── config.go                  # 合约配置
├── constants.go               # 常量定义
├── endpoints.go               # API 端点常量
├── http_client.go             # HTTP 客户端
├── http_helpers.go            # HTTP 辅助函数（查询参数构建）
├── signer.go                  # 签名器
├── signing_internal.go        # 签名实现（EIP-712, HMAC）
├── types.go                   # 类型定义
├── utilities.go               # 工具函数
├── order_summary_wrapper.go   # OrderSummary 包装器
├── headers/                   # 认证头（包装函数）
│   └── headers.go
├── order_builder/             # 订单构建器
│   ├── order_builder.go       # 订单构建器实现
│   └── helpers.go             # 订单构建辅助函数
└── rfq/                       # RFQ 客户端
    ├── rfq_client.go          # RFQ 客户端实现
    └── types.go               # RFQ 类型定义
```

## 已实现功能

### ✅ 核心功能
- [x] 基础类型定义（所有 py-clob-client 的类型）
- [x] 签名器 - EIP-712 和 HMAC 签名
- [x] L0/L1/L2 认证头生成
- [x] HTTP 客户端封装
- [x] 客户端基础结构（支持三种模式）

### ✅ API 方法
- [x] **健康检查**: `GetOK()`, `GetServerTime()`
- [x] **API 密钥管理**: `CreateAPIKey()`, `DeriveAPIKey()`, `CreateOrDeriveAPIKey()`, `GetAPIKeys()`, `DeleteAPIKey()`
- [x] **市场数据**: `GetMidpoint()`, `GetMidpoints()`, `GetPrice()`, `GetPrices()`, `GetSpread()`, `GetSpreads()`
- [x] **订单簿**: `GetOrderBook()`, `GetOrderBooks()`, `GetOrderBookHash()`
- [x] **市场信息**: `GetTickSize()`, `GetNegRisk()`, `GetFeeRateBps()` (带缓存)
- [x] **最后成交价**: `GetLastTradePrice()`, `GetLastTradesPrices()`

### ✅ 订单管理
- [x] **订单提交**: `PostOrder()`, `PostOrders()`
- [x] **订单取消**: `Cancel()`, `CancelOrders()`, `CancelAll()`, `CancelMarketOrders()`
- [x] **订单查询**: `GetOrders()`, `GetOrder()`
- [x] **交易查询**: `GetTrades()`
- [x] **余额查询**: `GetBalanceAllowance()`
- [x] **通知管理**: `GetNotifications()`, `DropNotifications()`

### ✅ 订单构建和创建
- [x] 订单构建器完整实现（使用 go-order-utils）
- [x] 订单创建方法：`CreateOrder()`, `CreateMarketOrder()`, `CreateAndPostOrder()`
- [x] 市价计算：`CalculateMarketPrice()`
- [x] 舍入配置和金额计算

### ✅ 只读 API 密钥管理
- [x] `CreateReadonlyAPIKey()` - 创建只读 API 密钥
- [x] `GetReadonlyAPIKeys()` - 获取只读 API 密钥列表
- [x] `DeleteReadonlyAPIKey()` - 删除只读 API 密钥
- [x] `ValidateReadonlyAPIKey()` - 验证只读 API 密钥

### ✅ RFQ 客户端功能
- [x] `CreateRfqRequest()` - 创建 RFQ 请求
- [x] `CancelRfqRequest()` - 取消 RFQ 请求
- [x] `GetRfqRequests()` - 获取 RFQ 请求列表
- [x] `CreateRfqQuote()` - 创建 RFQ 报价
- [x] `CancelRfqQuote()` - 取消 RFQ 报价
- [x] `GetRfqQuotes()` - 获取 RFQ 报价列表
- [x] `GetRfqBestQuote()` - 获取最佳 RFQ 报价
- [x] `AcceptQuote()` - 接受报价
- [x] `ApproveOrder()` - 批准订单
- [x] `GetRfqConfig()` - 获取 RFQ 配置

### ✅ 其他功能
- [x] 订单评分：`IsOrderScoring()`, `AreOrdersScoring()`
- [x] 市场查询：`GetMarkets()`, `GetSimplifiedMarkets()`, `GetSamplingMarkets()`, `GetSamplingSimplifiedMarkets()`
- [x] 市场详情：`GetMarket()`, `GetMarketTradesEvents()`
- [x] 余额更新：`UpdateBalanceAllowance()`
- [x] 订单簿哈希：`GetOrderBookHash()`
- [x] Builder 交易：`GetBuilderTrades()`

## 功能对比

| 功能 | py-clob-client | Go SDK | 状态 |
|------|---------------|--------|------|
| 基础客户端 | ✅ | ✅ | 完成 |
| L0/L1/L2 认证 | ✅ | ✅ | 完成 |
| API 密钥管理 | ✅ | ✅ | 完成 |
| 市场数据查询 | ✅ | ✅ | 完成 |
| 订单簿查询 | ✅ | ✅ | 完成 |
| 订单提交/取消 | ✅ | ✅ | 完成 |
| 订单查询 | ✅ | ✅ | 完成 |
| 交易查询 | ✅ | ✅ | 完成 |
| 余额查询 | ✅ | ✅ | 完成 |
| 订单构建器 | ✅ | ✅ | 完整实现 |
| RFQ 功能 | ✅ | ✅ | 完整实现 |
| 只读 API 密钥 | ✅ | ✅ | 完整实现 |

## 环境变量

示例程序使用以下环境变量：

- `PRIVATE_KEY` (必需): 您的以太坊私钥（十六进制格式）
- `CHAIN_ID` (可选): 链 ID（默认: 137，Polygon）
- `SIGNATURE_TYPE` (可选): 签名类型（0=EOA, 1=Magic/Email, 2=Browser，默认: 0）
- `FUNDER` (可选): 代理钱包的资金持有者地址
- `CLOB_HOST` (可选): CLOB API 主机（默认: https://clob.polymarket.com）
- `CLOB_API_KEY` (可选): L2 认证的 API 密钥
- `CLOB_SECRET` (可选): L2 认证的 API 密钥
- `CLOB_PASSPHRASE` (可选): L2 认证的 API 密钥
- `TOKEN_ID` (可选): 条件代币余额查询的 token ID

## 参考

- [py-clob-client](https://github.com/Polymarket/py-clob-client) - Python 官方实现
- [Polymarket CLOB 文档](https://docs.polymarket.com/developers/CLOB)
- [go-order-utils](https://github.com/polymarket/go-order-utils) - 订单构建工具库

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎贡献！请随时提交 Pull Request。

