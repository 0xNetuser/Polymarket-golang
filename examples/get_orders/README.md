# Get Orders 示例

这个示例展示如何使用 Polymarket Go SDK 获取订单信息。

## 功能

- 获取所有订单列表
- 使用过滤条件获取订单（按市场、资产ID或订单ID）
- 获取单个订单详情

## 环境变量

### 必需参数

| 变量 | 说明 | 示例 |
|------|------|------|
| `PRIVATE_KEY` | 私钥（带或不带0x前缀） | `0xabc123...` |

### 可选参数

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CHAIN_ID` | 链ID | `137` (Polygon) |
| `CLOB_HOST` | API主机地址 | `https://clob.polymarket.com` |
| `SIGNATURE_TYPE` | 签名类型：0=EOA, 1=Magic, 2=Browser | `0` |
| `FUNDER` | 代理钱包地址（用于proxy账户） | 空 |

### API凭证（可选，如果已有）

| 变量 | 说明 |
|------|------|
| `CLOB_API_KEY` | API Key |
| `CLOB_SECRET` | API Secret |
| `CLOB_PASSPHRASE` | API Passphrase |

### 过滤参数（可选）

| 变量 | 说明 |
|------|------|
| `MARKET` | 市场 condition_id |
| `ASSET_ID` | 资产 token_id |
| `ORDER_ID` | 订单 ID（用于过滤列表） |
| `SINGLE_ORDER_ID` | 单个订单 ID（获取详情） |

## 使用方法

### 1. 设置私钥并运行

```bash
export PRIVATE_KEY="your_private_key"
./run.sh
```

### 2. 获取所有订单

```bash
export PRIVATE_KEY="your_private_key"
go run main.go
```

### 3. 按市场过滤订单

```bash
export PRIVATE_KEY="your_private_key"
export MARKET="0x1234..."  # condition_id
go run main.go
```

### 4. 按资产ID过滤订单

```bash
export PRIVATE_KEY="your_private_key"
export ASSET_ID="12345678901234567890..."  # token_id
go run main.go
```

### 5. 获取单个订单详情

```bash
export PRIVATE_KEY="your_private_key"
export SINGLE_ORDER_ID="0xabcd1234..."  # 订单ID
go run main.go
```

### 6. 使用代理钱包

```bash
export PRIVATE_KEY="your_private_key"
export SIGNATURE_TYPE="1"
export FUNDER="0x4E45700535694bEE11F50c82010611AF4c9b76bD"
go run main.go
```

## 输出示例

```
=== Polymarket 获取订单示例 ===
地址: 0x974d5981685F986e76500765EB77FBe5453143e7
链ID: 137
签名类型: 1 (0=EOA, 1=Magic/Email, 2=Browser)

使用环境变量中的API凭证...

=== 获取所有订单 ===
找到 2 个订单:

--- 订单 1 ---
  ID: 0x65652efd7d5c9ac6034cf934c0db742cb45a1d2dcae59c95ec5c4b3c6f68c667
  状态: MATCHED
  方向: BUY
  价格: 0.01
  原始数量: 100
  已成交数量: 100
  资产ID: 10417355721474453...
  类型: GTC
  创建时间: 2025-12-20T02:49:57Z

--- 订单 2 ---
  ID: 0xf469ca2539bb25fee3ea66bd2bc7cbce76b1d3bc5eb5abdaa3945f0130e6b214
  状态: MATCHED
  方向: BUY
  价格: 0.01
  原始数量: 100
  已成交数量: 92.224
  资产ID: 10417355721474453...
  类型: GTC
  创建时间: 2025-12-20T02:48:58Z

=== 完成 ===
```

## 订单状态说明

| 状态 | 说明 |
|------|------|
| `LIVE` | 订单正在挂单中 |
| `MATCHED` | 订单已完全成交 |
| `PARTIALLY_MATCHED` | 订单部分成交 |
| `CANCELLED` | 订单已取消 |
| `DELAYED` | 订单延迟处理 |

## 注意事项

1. 获取订单列表需要 L2 认证（API凭证）
2. 如果没有提供API凭证，程序会自动派生或创建
3. 订单列表会自动分页获取所有结果

