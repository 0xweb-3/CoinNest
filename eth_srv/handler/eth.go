package handler

import (
	"context"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/collection_cold"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/deposit"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/node"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/withdraw"
	"github.com/0xweb-3/CoinNest/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync/atomic"
)

type EthWallet struct {
	ethClient      node.EthClient
	collectionCold *collection_cold.CollectionCold
	deposit        *deposit.Deposit
	withdraw       *withdraw.Withdraw

	shoutDown context.CancelCauseFunc
	stopped   atomic.Bool
}

func NewEthWallet(ctx context.Context, shoutDown context.CancelCauseFunc) (*EthWallet, error) {
	ethClient, err := node.DialEthClient(ctx, "")
	if err != nil {
		return nil, err
	}

	withdraw, err := withdraw.NewWithdraw(ethClient, shoutDown)
	if err != nil {
		return nil, err
	}
	deposit, err := deposit.NewDeposit(ethClient, shoutDown)
	if err != nil {
		return nil, err
	}
	collectionCold, err := collection_cold.NewCollectionCold(ethClient, shoutDown)
	if err != nil {
		return nil, err
	}

	out := &EthWallet{
		ethClient:      ethClient,
		collectionCold: collectionCold,
		deposit:        deposit,
		withdraw:       withdraw,
		shoutDown:      shoutDown,
	}

	return out, nil
}

func (ew *EthWallet) Start(ctx context.Context) error {
	err := ew.deposit.Start()
	if err != nil {
		return err
	}
	//err = ew.withdraw.Start()
	//if err != nil {
	//	return err
	//}
	err = ew.collectionCold.Start()
	if err != nil {
		return err
	}
	return nil
}

func (ew *EthWallet) Stop(ctx context.Context) error {
	err := ew.deposit.Close()
	if err != nil {
		return err
	}
	//err = ew.withdraw.Close()
	//if err != nil {
	//	return err
	//}

	err = ew.collectionCold.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ew *EthWallet) Stopped() bool {
	return ew.stopped.Load()
}

func NewUserRepo(db *gorm.DB) *EthRepo {
	return &EthRepo{
		db:  db,
		log: zap.S(),
	}
}

// GetUserById 获取账号信息
func (r *EthRepo) GetUserById(ctx context.Context, userId uint64) (*proto.UserInfo, error) {
	return nil, nil
}
