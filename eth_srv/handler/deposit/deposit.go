package deposit

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/common/tasks"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/node"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type Deposit struct {
	client         node.EthClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewDeposit(client node.EthClient, shutdown context.CancelCauseFunc) (*Deposit, error) {
	resCtx, resCancel := context.WithCancel(context.Background())

	return &Deposit{
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{
			HandleCrit: func(err error) {
				shutdown(fmt.Errorf("critical error in deposit: %w", err))
			},
		},
	}, nil
}

func (d *Deposit) Close() error {
	var result error
	d.resourceCancel()
	if err := d.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (d *Deposit) Start() error {
	log.Info("start withdraw......")
	tickerDepositWorker := time.NewTicker(time.Second * 5)
	d.tasks.Go(func() error {
		for range tickerDepositWorker.C {
			log.Info("Deposit work task go")
		}
		return nil
	})
	return nil
}
