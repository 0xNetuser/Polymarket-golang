package order_builder

import (
	"math"
	"strconv"
	"strings"
)

// RoundConfig 舍入配置
type RoundConfig struct {
	Price  int // 价格小数位数
	Size   int // 数量小数位数
	Amount int // 金额小数位数
}

// RoundingConfig 舍入配置映射
var RoundingConfig = map[string]RoundConfig{
	"0.1":    {Price: 1, Size: 2, Amount: 3},
	"0.01":   {Price: 2, Size: 2, Amount: 4},
	"0.001":  {Price: 3, Size: 2, Amount: 5},
	"0.0001": {Price: 4, Size: 2, Amount: 6},
}

// RoundDown 向下舍入
func RoundDown(x float64, sigDigits int) float64 {
	return math.Floor(x*math.Pow(10, float64(sigDigits))) / math.Pow(10, float64(sigDigits))
}

// RoundNormal 正常舍入
func RoundNormal(x float64, sigDigits int) float64 {
	return math.Round(x*math.Pow(10, float64(sigDigits))) / math.Pow(10, float64(sigDigits))
}

// RoundUp 向上舍入
func RoundUp(x float64, sigDigits int) float64 {
	return math.Ceil(x*math.Pow(10, float64(sigDigits))) / math.Pow(10, float64(sigDigits))
}

// ToTokenDecimals 转换为代币小数位（6位）
func ToTokenDecimals(x float64) int64 {
	f := (math.Pow(10, 6)) * x
	if DecimalPlaces(f) > 0 {
		f = RoundNormal(f, 0)
	}
	return int64(f)
}

// DecimalPlaces 计算小数位数
func DecimalPlaces(x float64) int {
	s := strconv.FormatFloat(x, 'f', -1, 64)
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		return len(parts[1])
	}
	return 0
}

