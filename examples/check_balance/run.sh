#!/bin/bash
# 钱包余额查询脚本

# 检查是否设置了私钥
if [ -z "$PRIVATE_KEY" ]; then
    echo "错误: 请设置环境变量 PRIVATE_KEY"
    echo "使用方法:"
    echo "  export PRIVATE_KEY=\"your-private-key-hex\""
    echo "  ./run.sh"
    exit 1
fi

# 运行程序
go run main.go
