package polymarket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignClobAuthMessage 签名CLOB认证消息（L1认证）
// 使用EIP-712标准签名
func SignClobAuthMessage(signer *Signer, timestamp int, nonce int) (string, error) {
	// 构建EIP-712域
	// 根据EIP-712标准，域分隔符的构建方式：
	// keccak256(0x1901 || keccak256("EIP712Domain(string name,string version,uint256 chainId)") || nameHash || versionHash || chainId)
	domainNameHash := crypto.Keccak256Hash([]byte(CLOBDomainName))
	domainVersionHash := crypto.Keccak256Hash([]byte(CLOBVersion))
	chainID := big.NewInt(int64(signer.GetChainID()))

	// EIP712Domain类型哈希：keccak256("EIP712Domain(string name,string version,uint256 chainId)")
	eip712DomainTypeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId)"))

	// 编码域数据
	domainData := make([]byte, 96) // 3个字段，每个32字节
	copy(domainData[0:32], domainNameHash[:])      // name hash
	copy(domainData[32:64], domainVersionHash[:])  // version hash
	copy(domainData[64:96], common.LeftPadBytes(chainID.Bytes(), 32)) // chainId

	// 构建域分隔符哈希
	domainHashBytes := append(eip712DomainTypeHash[:], domainData...)
	domainSeparator := crypto.Keccak256Hash(domainHashBytes)

	// 构建ClobAuth结构体
	address := common.HexToAddress(signer.Address())
	timestampStr := strconv.Itoa(timestamp)
	nonceBig := big.NewInt(int64(nonce))

	// EIP-712类型哈希：keccak256("ClobAuth(address address,string timestamp,uint256 nonce,string message)")
	// 注意：EIP-712规范要求类型字符串格式为 "type name"，必须包含字段名
	typeHash := crypto.Keccak256Hash([]byte("ClobAuth(address address,string timestamp,uint256 nonce,string message)"))

	// 根据EIP-712标准，字符串需要先进行keccak256哈希
	timestampHash := crypto.Keccak256Hash([]byte(timestampStr))
	messageStrHash := crypto.Keccak256Hash([]byte(MsgToSign))

	// 构建编码后的结构体值
	// address: 32字节（左对齐）
	// timestamp: 32字节（字符串的keccak256哈希）
	// nonce: 32字节（uint256）
	// message: 32字节（字符串的keccak256哈希）
	encoded := make([]byte, 128) // 4个字段，每个32字节
	
	// address (左对齐到32字节)
	copy(encoded[0:32], common.LeftPadBytes(address.Bytes(), 32))
	
	// timestamp hash (32字节)
	copy(encoded[32:64], timestampHash[:])
	
	// nonce (32字节，左对齐)
	copy(encoded[64:96], common.LeftPadBytes(nonceBig.Bytes(), 32))
	
	// message hash (32字节)
	copy(encoded[96:128], messageStrHash[:])

	// 构建结构体哈希：keccak256(typeHash || encodedValues)
	structHashBytes := append(typeHash[:], encoded...)
	structHash := crypto.Keccak256Hash(structHashBytes)

	// 构建signable_bytes（对应Python的signable_bytes方法）
	// 根据EIP-712标准和poly_eip712_structs的实现：
	// signable_bytes返回 "\x19\x01" || domainSeparator || structHash
	// 注意：domainSeparator本身已经是keccak256哈希，不需要再次哈希
	prefix := []byte("\x19\x01")
	signableBytes := append(prefix, domainSeparator[:]...)
	signableBytes = append(signableBytes, structHash.Bytes()...)

	// Python代码：keccak(signable_bytes).hex()
	// 对signable_bytes进行keccak256哈希
	authStructHash := crypto.Keccak256Hash(signableBytes)

	// Python代码：signer.sign(auth_struct_hash)
	// Account._sign_hash接收hex字符串，内部会解码为字节并签名
	// 我们直接对哈希值进行签名（等价于解码hex字符串后签名）
	signature, err := signer.Sign(authStructHash.Bytes())
	if err != nil {
		return "", err
	}

	// 添加0x前缀（如果还没有）
	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}

	return signature, nil
}

// BuildHMACSignature 构建HMAC签名（L2认证）
func BuildHMACSignature(secret string, timestamp int, method string, requestPath string, body interface{}) (string, error) {
	// Base64解码secret
	base64Secret, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	// 构建消息
	message := strconv.Itoa(timestamp) + method + requestPath
	if body != nil {
		// 将body转换为JSON字符串
		// 注意：必须使用带空格的JSON格式，否则API会返回401错误
		// 参考: https://github.com/Polymarket/py-clob-client/issues/164
		var bodyStr string
		if bodyStrPtr, ok := body.(*string); ok {
			// 如果已经是字符串指针，直接使用
			bodyStr = *bodyStrPtr
		} else if bodyStrVal, ok := body.(string); ok {
			// 如果已经是字符串，直接使用
			bodyStr = bodyStrVal
		} else {
			// 序列化为JSON（紧凑格式，无空格）
			// 参考: https://github.com/Polymarket/py-clob-client/issues/164
			// API要求JSON必须去掉所有空格，否则会返回401错误
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return "", fmt.Errorf("failed to marshal body for HMAC signing: %w", err)
			}
			bodyStr = string(bodyBytes)
		}
		message += bodyStr
	}

	// HMAC-SHA256
	mac := hmac.New(sha256.New, base64Secret)
	mac.Write([]byte(message))
	digest := mac.Sum(nil)

	// Base64编码
	return base64.URLEncoding.EncodeToString(digest), nil
}

// CreateLevel1Headers 创建L1认证头
func CreateLevel1Headers(signer *Signer, nonce *int) (map[string]string, error) {
	timestamp := int(time.Now().Unix())

	n := 0
	if nonce != nil {
		n = *nonce
	}

	signature, err := SignClobAuthMessage(signer, timestamp, n)
	if err != nil {
		return nil, err
	}

	// 确保地址使用checksummed格式（与Python的eth_account一致）
	address := common.HexToAddress(signer.Address())
	
	headers := map[string]string{
		PolyAddress:   address.Hex(), // 使用checksummed格式
		PolySignature: signature,
		PolyTimestamp: strconv.Itoa(timestamp),
		PolyNonce:     strconv.Itoa(n),
	}

	// 调试输出（可以通过环境变量控制）
	if os.Getenv("DEBUG_HEADERS") == "1" {
		fmt.Fprintf(os.Stderr, "=== L1 Headers ===\n")
		for k, v := range headers {
			fmt.Fprintf(os.Stderr, "%s: %s\n", k, v)
		}
		fmt.Fprintf(os.Stderr, "==================\n")
	}

	return headers, nil
}

// CreateLevel2Headers 创建L2认证头
func CreateLevel2Headers(signer *Signer, creds *ApiCreds, requestArgs *RequestArgs) (map[string]string, error) {
	timestamp := int(time.Now().Unix())

	// 优先使用预序列化的body
	var bodyForSig interface{}
	if requestArgs.SerializedBody != nil {
		bodyForSig = *requestArgs.SerializedBody
	} else {
		bodyForSig = requestArgs.Body
	}

	hmacSig, err := BuildHMACSignature(
		creds.APISecret,
		timestamp,
		requestArgs.Method,
		requestArgs.RequestPath,
		bodyForSig,
	)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		PolyAddress:    signer.Address(),
		PolySignature:  hmacSig,
		PolyTimestamp:  strconv.Itoa(timestamp),
		PolyAPIKey:     creds.APIKey,
		PolyPassphrase: creds.APIPassphrase,
	}, nil
}

