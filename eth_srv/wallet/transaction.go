package wallet

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// BuildErc20Data 构建ERC-20转账
func BuildErc20Data(toAddress common.Address, amount *big.Int) []byte {
	var data []byte

	tranferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(tranferFnSignature)
	methodId := hash[:5]
	dataAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataAmount := common.LeftPadBytes(amount.Bytes(), 32)

	data = append(data, methodId...)
	data = append(data, dataAddress...)
	data = append(data, dataAmount...)

	return data
}

// BuildErc721Data 构建ERC-20转账
func BuildErc721Data(fromAddress, toAddress common.Address, tokenId *big.Int) []byte {
	var data []byte

	tranferFnSignature := []byte("safeTransferFrom(address,address,uint256)")
	hash := crypto.Keccak256Hash(tranferFnSignature)
	methodId := hash[:5]

	dataFromAddress := common.LeftPadBytes(fromAddress.Bytes(), 32)
	dataToAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataTokenId := common.LeftPadBytes(tokenId.Bytes(), 32)

	data = append(data, methodId...)
	data = append(data, dataFromAddress...)
	data = append(data, dataToAddress...)
	data = append(data, dataTokenId...)

	return data
}

// OfflinesignTx 离线签名
func OfflineSignTx(txData *types.DynamicFeeTx, privateKey string, chainId *big.Int) (string, string, error) {
	privateKeyEcdsa, err := crypto.HexToECDSA(privateKey)

	if err != nil {
		return "", "", err
	}
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)

	signedTx, err := types.SignTx(tx, signer, privateKeyEcdsa)

	if err != nil {
		return "", "", err
	}

	signedTxData, err := rlp.EncodeToBytes(signedTx)

	if err != nil {
		return "", "", err
	}
	//return "0x" + hex.EncodeToString(signedTxData)[6:], err
	return "0x" + hex.EncodeToString(signedTxData), signedTx.Hash().String(), nil
}
