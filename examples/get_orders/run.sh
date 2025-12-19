#!/bin/bash

# Polymarket 获取订单示例运行脚本

# 基本配置
export PRIVATE_KEY="${PRIVATE_KEY:-}"
export CHAIN_ID="${CHAIN_ID:-137}"
export CLOB_HOST="${CLOB_HOST:-https://clob.polymarket.com}"

# 账户类型和代理地址
export SIGNATURE_TYPE="${SIGNATURE_TYPE:-0}"  # 0=EOA, 1=Magic/Email, 2=Browser
export FUNDER="${FUNDER:-}"  # 代理钱包地址（可选）

# API凭证（可选，如果已有）
export CLOB_API_KEY="${CLOB_API_KEY:-}"
export CLOB_SECRET="${CLOB_SECRET:-}"
export CLOB_PASSPHRASE="${CLOB_PASSPHRASE:-}"

# 过滤参数（可选）
export MARKET="${MARKET:-}"           # 市场 condition_id
export ASSET_ID="${ASSET_ID:-}"       # 资产 token_id
export ORDER_ID="${ORDER_ID:-}"       # 订单 ID（用于过滤）
export SINGLE_ORDER_ID="${SINGLE_ORDER_ID:-}"  # 单个订单 ID（获取详情）

# 运行程序
go run main.go

