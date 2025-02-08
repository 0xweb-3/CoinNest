package collection_cold

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/common/tasks"
	"github.com/0xweb-3/CoinNest/eth_srv/handler/node"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type CollectionCold struct {
	client         node.EthClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewCollectionCold(client node.EthClient, shutdown context.CancelCauseFunc) (*CollectionCold, error) {
	resCtx, resCancel := context.WithCancel(context.Background())

	return &CollectionCold{
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{
			HandleCrit: func(err error) {
				shutdown(fmt.Errorf("critical error in collection cold: %w", err))
			},
		},
	}, nil
}

func (cc *CollectionCold) Close() error {
	var result error
	cc.resourceCancel()
	if err := cc.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await collection cold %w"), err)
	}
	return nil
}

func (cc *CollectionCold) Start() error {
	log.Info("start collection cold......")
	tickerCollectionColdWorker := time.NewTicker(time.Second * 5)
	cc.tasks.Go(func() error {
		for range tickerCollectionColdWorker.C {
			log.Info("collection cold work task go")
		}
		return nil
	})
	return nil
}
