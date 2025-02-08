package withdraw

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/common/tasks"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/node"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type Withdraw struct {
	client         node.EthClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWithdraw(client node.EthClient, shutdown context.CancelCauseFunc) (*Withdraw, error) {
	resCtx, resCancel := context.WithCancel(context.Background())

	return &Withdraw{
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{
			HandleCrit: func(err error) {
				shutdown(fmt.Errorf("critical error in withdraw: %w", err))
			},
		},
	}, nil
}

func (w *Withdraw) Close() error {
	var result error
	w.resourceCancel()
	if err := w.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await withdraw %w"), err)
	}
	return nil
}

func (w *Withdraw) Start() error {
	log.Info("start withdraw......")
	tickerWithdrawWorker := time.NewTicker(time.Second * 5)
	w.tasks.Go(func() error {
		for range tickerWithdrawWorker.C {
			log.Info("Withdraw work task go")
		}
		return nil
	})
	return nil
}
