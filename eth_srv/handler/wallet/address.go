package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthAddress struct {
	PrivateKey string `json:"private_key"` // 地址密钥
	PublicKey  string `json:"public_key"`  // 地址公钥
	Address    string `json:"address"`     // 地址
}

// CreateAddressFromPrivateKey 通过密钥生成地址
func CreateAddressFromPrivateKey(priKey *ecdsa.PrivateKey) (string, string, error) {
	return hex.EncodeToString(priKey.D.Bytes()), crypto.PubkeyToAddress(priKey.PublicKey).String(), nil
}

// CreateAddressByKeyPairs 直接生成地址
func CreateAddressByKeyPairs() (*EthAddress, error) {
	prvKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := &EthAddress{
		PrivateKey: hex.EncodeToString(crypto.FromECDSA(prvKey)),
		PublicKey:  hex.EncodeToString(crypto.FromECDSAPub(&prvKey.PublicKey)),
		Address:    crypto.PubkeyToAddress(prvKey.PublicKey).String(),
	}
	return address, nil
}

// PublicKeyToAddress 将公钥转换为地址
func PublicKeyToAddress(publicKey string) (string, error) {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", nil
	}

	addressCommon := common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:])
	return addressCommon.String(), err
}
