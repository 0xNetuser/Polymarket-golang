package order_builder

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/polymarket/go-order-utils/pkg/builder"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// Signer 签名器接口（避免循环导入）
type Signer interface {
	Address() string
	GetChainID() int
	GetPrivateKey() string
}

// OrderBuilder 订单构建器
type OrderBuilder struct {
	signer      Signer
	sigType     int
	funder      string
}

// NewOrderBuilder 创建新的订单构建器
func NewOrderBuilder(signer Signer, sigType int, funder string) (*OrderBuilder, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required")
	}

	if funder == "" {
		funder = signer.Address()
	}

	return &OrderBuilder{
		signer:  signer,
		sigType: sigType,
		funder:  funder,
	}, nil
}

// GetOrderAmounts 获取订单金额（限价订单）
func (ob *OrderBuilder) GetOrderAmounts(side string, size, price float64, roundConfig RoundConfig) (model.Side, *big.Int, *big.Int, error) {
	rawPrice := RoundNormal(price, roundConfig.Price)

	if side == "BUY" {
		rawTakerAmt := RoundDown(size, roundConfig.Size)
		rawMakerAmt := rawTakerAmt * rawPrice

		if DecimalPlaces(rawMakerAmt) > roundConfig.Amount {
			rawMakerAmt = RoundUp(rawMakerAmt, roundConfig.Amount+4)
			if DecimalPlaces(rawMakerAmt) > roundConfig.Amount {
				rawMakerAmt = RoundDown(rawMakerAmt, roundConfig.Amount)
			}
		}

		makerAmount := big.NewInt(ToTokenDecimals(rawMakerAmt))
		takerAmount := big.NewInt(ToTokenDecimals(rawTakerAmt))

		return model.BUY, makerAmount, takerAmount, nil
	} else if side == "SELL" {
		rawMakerAmt := RoundDown(size, roundConfig.Size)
		rawTakerAmt := rawMakerAmt * rawPrice

		if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = RoundUp(rawTakerAmt, roundConfig.Amount+4)
			if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
				rawTakerAmt = RoundDown(rawTakerAmt, roundConfig.Amount)
			}
		}

		makerAmount := big.NewInt(ToTokenDecimals(rawMakerAmt))
		takerAmount := big.NewInt(ToTokenDecimals(rawTakerAmt))

		return model.SELL, makerAmount, takerAmount, nil
	}

		return 0, nil, nil, fmt.Errorf("order_args.side must be 'BUY' or 'SELL'")
}

// GetMarketOrderAmounts 获取市价订单金额
func (ob *OrderBuilder) GetMarketOrderAmounts(side string, amount, price float64, roundConfig RoundConfig) (model.Side, *big.Int, *big.Int, error) {
	rawPrice := RoundNormal(price, roundConfig.Price)

	if side == "BUY" {
		rawMakerAmt := RoundDown(amount, roundConfig.Size)
		rawTakerAmt := rawMakerAmt / rawPrice

		if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = RoundUp(rawTakerAmt, roundConfig.Amount+4)
			if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
				rawTakerAmt = RoundDown(rawTakerAmt, roundConfig.Amount)
			}
		}

		makerAmount := big.NewInt(ToTokenDecimals(rawMakerAmt))
		takerAmount := big.NewInt(ToTokenDecimals(rawTakerAmt))

		return model.BUY, makerAmount, takerAmount, nil
	} else if side == "SELL" {
		rawMakerAmt := RoundDown(amount, roundConfig.Size)
		rawTakerAmt := rawMakerAmt * rawPrice

		if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = RoundUp(rawTakerAmt, roundConfig.Amount+4)
			if DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
				rawTakerAmt = RoundDown(rawTakerAmt, roundConfig.Amount)
			}
		}

		makerAmount := big.NewInt(ToTokenDecimals(rawMakerAmt))
		takerAmount := big.NewInt(ToTokenDecimals(rawTakerAmt))

		return model.SELL, makerAmount, takerAmount, nil
	}

		return 0, nil, nil, fmt.Errorf("order_args.side must be 'BUY' or 'SELL'")
}

// CreateOrder 创建并签名订单（限价订单）
func (ob *OrderBuilder) CreateOrder(orderArgs interface{}, options interface{}) (*model.SignedOrder, error) {
	// 这里需要从主包传入类型，暂时使用interface{}
	// 实际使用时需要类型断言
	return nil, fmt.Errorf("CreateOrder需要从主包调用，传入正确的类型")
}

// CreateMarketOrder 创建并签名市价订单
func (ob *OrderBuilder) CreateMarketOrder(orderArgs interface{}, options interface{}) (*model.SignedOrder, error) {
	// 这里需要从主包传入类型，暂时使用interface{}
	// 实际使用时需要类型断言
	return nil, fmt.Errorf("CreateMarketOrder需要从主包调用，传入正确的类型")
}

// OrderSummary 订单摘要接口（避免循环导入）
type OrderSummary interface {
	GetPrice() string
	GetSize() string
}

// CalculateBuyMarketPrice 计算买入市价
func (ob *OrderBuilder) CalculateBuyMarketPrice(positions []interface{}, amountToMatch float64, orderType string) (float64, error) {
	if len(positions) == 0 {
		return 0, fmt.Errorf("no match")
	}

	sum := 0.0
	for i := len(positions) - 1; i >= 0; i-- {
		pos, ok := positions[i].(OrderSummary)
		if !ok {
			continue
		}

		price, _ := strconv.ParseFloat(pos.GetPrice(), 64)
		size, _ := strconv.ParseFloat(pos.GetSize(), 64)
		sum += size * price

		if sum >= amountToMatch {
			return price, nil
		}
	}

	if orderType == "FOK" {
		return 0, fmt.Errorf("no match")
	}

	// 返回第一个价格
	if pos, ok := positions[0].(OrderSummary); ok {
		price, _ := strconv.ParseFloat(pos.GetPrice(), 64)
		return price, nil
	}

	return 0, fmt.Errorf("invalid position format")
}

// CalculateSellMarketPrice 计算卖出市价
func (ob *OrderBuilder) CalculateSellMarketPrice(positions []interface{}, amountToMatch float64, orderType string) (float64, error) {
	if len(positions) == 0 {
		return 0, fmt.Errorf("no match")
	}

	sum := 0.0
	for i := len(positions) - 1; i >= 0; i-- {
		pos, ok := positions[i].(OrderSummary)
		if !ok {
			continue
		}

		size, _ := strconv.ParseFloat(pos.GetSize(), 64)
		sum += size

		if sum >= amountToMatch {
			price, _ := strconv.ParseFloat(pos.GetPrice(), 64)
			return price, nil
		}
	}

	if orderType == "FOK" {
		return 0, fmt.Errorf("no match")
	}

	// 返回第一个价格
	if pos, ok := positions[0].(OrderSummary); ok {
		price, _ := strconv.ParseFloat(pos.GetPrice(), 64)
		return price, nil
	}

	return 0, fmt.Errorf("invalid position format")
}

// BuildSignedOrder 构建已签名订单（导出方法，供主包使用）
func (ob *OrderBuilder) BuildSignedOrder(orderData *model.OrderData, exchangeAddr string, chainID int, negRisk bool) (*model.SignedOrder, error) {
	// 解析私钥
	privateKeyHex := ob.signer.GetPrivateKey()
	// 移除0x前缀（如果有）
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// 创建订单构建器
	chainIDBig := big.NewInt(int64(chainID))
	orderBuilder := builder.NewExchangeOrderBuilderImpl(chainIDBig, nil)

	// VerifyingContract是int类型：CTFExchange=0, NegRiskCTFExchange=1
	var contract model.VerifyingContract
	if negRisk {
		contract = model.NegRiskCTFExchange
	} else {
		contract = model.CTFExchange
	}

	// 构建并签名订单
	signedOrder, err := orderBuilder.BuildSignedOrder(privateKey, orderData, contract)
	if err != nil {
		return nil, fmt.Errorf("failed to build signed order: %w", err)
	}

	return signedOrder, nil
}

// GetSigType 获取签名类型
func (ob *OrderBuilder) GetSigType() int {
	return ob.sigType
}

// GetFunder 获取资金持有者地址
func (ob *OrderBuilder) GetFunder() string {
	return ob.funder
}
