package polymarket

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer 签名器
type Signer struct {
	privateKey *ecdsa.PrivateKey
	chainID    int
	address    string
}

// NewSigner 创建新的签名器
func NewSigner(privateKeyHex string, chainID int) (*Signer, error) {
	if privateKeyHex == "" || chainID == 0 {
		return nil, fmt.Errorf("private key and chain ID are required")
	}

	// 移除0x前缀
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &Signer{
		privateKey: privateKey,
		chainID:    chainID,
		address:    address,
	}, nil
}

// Address 返回签名器的地址
func (s *Signer) Address() string {
	return s.address
}

// GetChainID 返回链ID
func (s *Signer) GetChainID() int {
	return s.chainID
}

// Sign 签名消息哈希
// 对于EIP-712，messageHash已经是最终的哈希值，不需要TextHash
// 对于普通消息签名，应该先使用TextHash
func (s *Signer) Sign(messageHash []byte) (string, error) {
	// 直接签名哈希值（EIP-712已经处理了前缀）
	signature, err := crypto.Sign(messageHash, s.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	// 添加恢复ID（v = 27 或 28）
	signature[64] += 27

	return hexutil.Encode(signature), nil
}

// GetPrivateKey 返回私钥（用于订单构建器）
func (s *Signer) GetPrivateKey() string {
	keyBytes := crypto.FromECDSA(s.privateKey)
	// 移除0x前缀（如果有）
	keyHex := hexutil.Encode(keyBytes)
	if len(keyHex) > 2 && keyHex[:2] == "0x" {
		return keyHex[2:]
	}
	return keyHex
}

