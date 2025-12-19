# 钱包余额查询示例

这个示例程序演示如何使用 Polymarket Go SDK 查询钱包余额。

## 功能

- 查询抵押品（USDC）余额和授权
- 查询条件代币余额和授权
- 自动创建或派生API凭证

## 使用方法

### 1. 设置环境变量

```bash
# 必需：私钥（十六进制格式，带或不带0x前缀都可以）
export PRIVATE_KEY="your-private-key-hex"

# 可选：签名类型（0=EOA钱包, 1=Magic/Email钱包, 2=Browser代理）
# 对于Magic钱包或邮箱登录，设置为1
export SIGNATURE_TYPE="1"

# 可选：CLOB API端点（默认为 https://clob.polymarket.com）
export CLOB_API_URL="https://clob.polymarket.com"

# 可选：如果已有API凭证，可以直接设置（避免每次创建）
export CLOB_API_KEY="your-api-key"
export CLOB_SECRET="your-api-secret"
export CLOB_PASSPHRASE="your-passphrase"

# 可选：查询特定条件代币余额时使用
export TOKEN_ID="your-token-id"
```

### 2. 运行程序

```bash
cd examples/check_balance
go run main.go
```

### 3. 输出示例

```
=== Polymarket 钱包余额查询 ===
地址: 0x1234567890123456789012345678901234567890
链ID: 137

使用环境变量中的API凭证...

=== 查询抵押品余额 ===

抵押品 (USDC):
{
  "balance": "1000000000",
  "allowance": "1000000000"
}
  余额: 1000000000
  授权: 1000000000

=== 查询条件代币余额 ===
提示: 需要提供具体的 token_id 才能查询条件代币余额
示例: 从订单簿或市场信息中获取 token_id
未设置 TOKEN_ID 环境变量，跳过条件代币查询
设置方式: export TOKEN_ID="your-token-id"

=== 查询完成 ===
```

## 注意事项

1. **私钥安全**: 私钥是敏感信息，请妥善保管，不要提交到版本控制系统
2. **API凭证**: 首次运行时会自动创建或派生API凭证，请保存好这些凭证以便后续使用
3. **余额单位**: 余额以最小单位返回（USDC为6位小数，即余额1000000表示1 USDC）
4. **授权额度**: allowance表示已授权给交易所的额度，用于交易

## 余额单位转换

USDC使用6位小数，所以：
- 余额 `1000000` = 1 USDC
- 余额 `1000000000` = 1000 USDC

条件代币也使用6位小数。

## 错误处理

如果遇到以下错误：
- `L2AuthUnavailable`: 需要设置API凭证
- `invalid private key`: 私钥格式不正确
- `API request failed`: 网络或API错误

请检查环境变量设置和网络连接。

