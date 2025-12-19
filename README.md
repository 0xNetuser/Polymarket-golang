# Polymarket Go SDK

[English](README.md) | [中文](README_zh.md)

A comprehensive Go SDK for the Polymarket CLOB (Central Limit Order Book) API, fully implementing all core features from [py-clob-client](https://github.com/Polymarket/py-clob-client).

Follow at X:  @netu5er



## Features

- ✅ **Complete API Coverage**: All endpoints from py-clob-client
- ✅ **Three Authentication Levels**: L0 (read-only), L1 (private key), L2 (full access)
- ✅ **Order Management**: Create, submit, cancel, and query orders
- ✅ **Market Data**: Order books, prices, spreads, and market information
- ✅ **RFQ Support**: Request for Quote functionality
- ✅ **Type Safety**: Strong typing with Go's type system
- ✅ **EIP-712 Signing**: Full support for Ethereum message signing

## Installation

```bash
go get github.com/0xNetuser/Polymarket-golang
```

## Quick Start

### Level 0 (Read-Only Mode)

```go
package main

import (
    "fmt"
    "github.com/0xNetuser/Polymarket-golang/polymarket"
)

func main() {
    // Create read-only client
    client, err := polymarket.NewClobClient(
        "https://clob.polymarket.com",
        137, // Polygon chain ID
        "",  // No private key
        nil, // No API credentials
        nil, // No signature type
        "",  // No funder address
    )
    if err != nil {
        panic(err)
    }

    // Health check
    ok, err := client.GetOK()
    fmt.Println("Server OK:", ok)

    // Get order book
    orderBook, err := client.GetOrderBook("token-id")
    if err != nil {
        panic(err)
    }
    fmt.Printf("OrderBook: %+v\n", orderBook)
}
```

### Level 1 (Private Key Authentication)

```go
client, err := polymarket.NewClobClient(
    "https://clob.polymarket.com",
    137,
    "your-private-key-hex",
    nil,
    nil,
    "",
)

// Create or derive API credentials
creds, err := client.CreateOrDeriveAPIKey(nil)
if err != nil {
    panic(err)
}
fmt.Printf("API Key: %s\n", creds.APIKey)
```

### Level 2 (Full Functionality)

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

// Query balance
balance, err := client.GetBalanceAllowance(&polymarket.BalanceAllowanceParams{
    AssetType: polymarket.AssetTypeCollateral,
})
if err != nil {
    panic(err)
}
fmt.Printf("Balance: %+v\n", balance)
```

## Examples

### Check Balance

```bash
cd examples/check_balance
export PRIVATE_KEY="your-private-key"
export CHAIN_ID="137"  # Optional, defaults to 137 (Polygon)
export SIGNATURE_TYPE="0"  # Optional, 0=EOA, 1=Magic/Email, 2=Browser
export FUNDER=""  # Optional, for proxy wallets
go run main.go
```

### Create and Post Order

```go
// Create order
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

// Post order
result, err := client.PostOrder(order, polymarket.OrderTypeGTC)
if err != nil {
    panic(err)
}
fmt.Printf("Order posted: %+v\n", result)
```

## Project Structure

```
polymarket/
├── client.go                  # Main client structure
├── client_api.go              # API methods (health check, API keys, market data)
├── client_orders.go           # Order management methods (submit, cancel, query)
├── client_order_creation.go   # Order creation methods (CreateOrder, CreateMarketOrder)
├── client_misc.go             # Other features (readonly API keys, order scoring, market queries)
├── rfq_client.go              # RFQ client convenience methods
├── config.go                  # Contract configuration
├── constants.go               # Constants
├── endpoints.go               # API endpoint constants
├── http_client.go             # HTTP client
├── http_helpers.go            # HTTP helper functions (query parameter building)
├── signer.go                  # Signer
├── signing_internal.go        # Signing implementation (EIP-712, HMAC)
├── types.go                   # Type definitions
├── utilities.go               # Utility functions
├── order_summary_wrapper.go   # OrderSummary wrapper
├── headers/                   # Authentication headers (wrapper functions)
│   └── headers.go
├── order_builder/             # Order builder
│   ├── order_builder.go       # Order builder implementation
│   └── helpers.go             # Order builder helper functions
└── rfq/                       # RFQ client
    ├── rfq_client.go          # RFQ client implementation
    └── types.go               # RFQ type definitions
```

## Implemented Features

### ✅ Core Features
- [x] Basic type definitions (all py-clob-client types)
- [x] Signer - EIP-712 and HMAC signing
- [x] L0/L1/L2 authentication header generation
- [x] HTTP client wrapper
- [x] Client base structure (supports three modes)

### ✅ API Methods
- [x] **Health Check**: `GetOK()`, `GetServerTime()`
- [x] **API Key Management**: `CreateAPIKey()`, `DeriveAPIKey()`, `CreateOrDeriveAPIKey()`, `GetAPIKeys()`, `DeleteAPIKey()`
- [x] **Market Data**: `GetMidpoint()`, `GetMidpoints()`, `GetPrice()`, `GetPrices()`, `GetSpread()`, `GetSpreads()`
- [x] **Order Book**: `GetOrderBook()`, `GetOrderBooks()`, `GetOrderBookHash()`
- [x] **Market Info**: `GetTickSize()`, `GetNegRisk()`, `GetFeeRateBps()` (with caching)
- [x] **Last Trade Price**: `GetLastTradePrice()`, `GetLastTradesPrices()`

### ✅ Order Management
- [x] **Order Submission**: `PostOrder()`, `PostOrders()`
- [x] **Order Cancellation**: `Cancel()`, `CancelOrders()`, `CancelAll()`, `CancelMarketOrders()`
- [x] **Order Query**: `GetOrders()`, `GetOrder()`
- [x] **Trade Query**: `GetTrades()`
- [x] **Balance Query**: `GetBalanceAllowance()`
- [x] **Notification Management**: `GetNotifications()`, `DropNotifications()`

### ✅ Order Building and Creation
- [x] Complete order builder implementation (using go-order-utils)
- [x] Order creation methods: `CreateOrder()`, `CreateMarketOrder()`, `CreateAndPostOrder()`
- [x] Market price calculation: `CalculateMarketPrice()`
- [x] Rounding configuration and amount calculation

### ✅ Readonly API Key Management
- [x] `CreateReadonlyAPIKey()` - Create readonly API key
- [x] `GetReadonlyAPIKeys()` - Get readonly API key list
- [x] `DeleteReadonlyAPIKey()` - Delete readonly API key
- [x] `ValidateReadonlyAPIKey()` - Validate readonly API key

### ✅ RFQ Client Features
- [x] `CreateRfqRequest()` - Create RFQ request
- [x] `CancelRfqRequest()` - Cancel RFQ request
- [x] `GetRfqRequests()` - Get RFQ request list
- [x] `CreateRfqQuote()` - Create RFQ quote
- [x] `CancelRfqQuote()` - Cancel RFQ quote
- [x] `GetRfqQuotes()` - Get RFQ quote list
- [x] `GetRfqBestQuote()` - Get best RFQ quote
- [x] `AcceptQuote()` - Accept quote
- [x] `ApproveOrder()` - Approve order
- [x] `GetRfqConfig()` - Get RFQ configuration

### ✅ Other Features
- [x] Order scoring: `IsOrderScoring()`, `AreOrdersScoring()`
- [x] Market queries: `GetMarkets()`, `GetSimplifiedMarkets()`, `GetSamplingMarkets()`, `GetSamplingSimplifiedMarkets()`
- [x] Market details: `GetMarket()`, `GetMarketTradesEvents()`
- [x] Balance update: `UpdateBalanceAllowance()`
- [x] Order book hash: `GetOrderBookHash()`
- [x] Builder trades: `GetBuilderTrades()`

## Feature Comparison

| Feature | py-clob-client | Go SDK | Status |
|---------|---------------|--------|--------|
| Basic client | ✅ | ✅ | Complete |
| L0/L1/L2 auth | ✅ | ✅ | Complete |
| API key management | ✅ | ✅ | Complete |
| Market data queries | ✅ | ✅ | Complete |
| Order book queries | ✅ | ✅ | Complete |
| Order submit/cancel | ✅ | ✅ | Complete |
| Order queries | ✅ | ✅ | Complete |
| Trade queries | ✅ | ✅ | Complete |
| Balance queries | ✅ | ✅ | Complete |
| Order builder | ✅ | ✅ | Fully implemented |
| RFQ features | ✅ | ✅ | Fully implemented |
| Readonly API keys | ✅ | ✅ | Fully implemented |

## Environment Variables

The example programs use the following environment variables:

- `PRIVATE_KEY` (required): Your Ethereum private key in hex format
- `CHAIN_ID` (optional): Chain ID (default: 137 for Polygon)
- `SIGNATURE_TYPE` (optional): Signature type (0=EOA, 1=Magic/Email, 2=Browser, default: 0)
- `FUNDER` (optional): Funder address for proxy wallets
- `CLOB_HOST` (optional): CLOB API host (default: https://clob.polymarket.com)
- `CLOB_API_KEY` (optional): API key for L2 authentication
- `CLOB_SECRET` (optional): API secret for L2 authentication
- `CLOB_PASSPHRASE` (optional): API passphrase for L2 authentication
- `TOKEN_ID` (optional): Token ID for conditional token balance queries

## References

- [py-clob-client](https://github.com/Polymarket/py-clob-client) - Official Python implementation
- [Polymarket CLOB Documentation](https://docs.polymarket.com/developers/CLOB)
- [go-order-utils](https://github.com/polymarket/go-order-utils) - Order building utility library

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
