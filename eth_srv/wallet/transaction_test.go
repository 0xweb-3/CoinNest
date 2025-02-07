package wallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
)

func TestOfflineSignTx(t *testing.T) {
	privateKeyHex := "0cbb2ff952da876c4779200c83f6b90d73ea85a8da82e06c2276a11499922720"
	nonce := uint64(58)
	toAddress := common.HexToAddress("0x35096AD62E57e86032a3Bb35aDaCF2240d55421D")
	amount := big.NewInt(1000000000000)
	gasLimit := uint64(21000)                      // 燃料上限，指定交易可以消耗的最大 Gas 量。
	maxPriorityFeePerGas := big.NewInt(2600000000) //最大优先费， 付给矿工的额外小费，以提高交易优先级。
	maxFeePerGas := big.NewInt(2900000000)         // 最大 Gas 费上限，交易可接受的最高 Gas 费用上限
	chainID := big.NewInt(1)
	dFeeTx := &types.DynamicFeeTx{
		ChainID: chainID,
		Nonce:   nonce,
		// GasTipCap 和 GasFeeCap 分别等于 maxPriorityFeePerGas 和 maxFeePerGas，控制交易的费用策略。
		GasTipCap: maxPriorityFeePerGas, // 交易中的小费上限，确定矿工可获得的最大小费金额，直接影响交易被打包的速度
		GasFeeCap: maxFeePerGas,         // 交易中的 Gas 费用上限，交易可承受的最大 Gas 费用，避免 Gas 价格波动导致交易花费超出预期。
		Gas:       gasLimit,             // 交易的 Gas 量，指定交易最多能使用的 Gas 量，防止恶意消耗资源。
		To:        &toAddress,
		Value:     amount,
		Data:      nil,
	}
	txHex, txHash, _ := OfflineSignTx(dFeeTx, privateKeyHex, chainID)
	fmt.Println("txHex===", txHex, "txHash==", txHash)
}
