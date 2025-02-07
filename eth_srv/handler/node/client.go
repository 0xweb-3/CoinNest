package node

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// 超时时间的定义
const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

// Logs 对获取的日志信息做包装
type Logs struct {
	Logs          []types.Log
	toBlockHeader *types.Header // 标记日志信息的来源
}

type EthClient interface {
	LatestBlockHeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	LatestBlockByNumber(context.Context, *big.Int) (types.Block, error)
	LatestSafeBlockByNumber(context.Context, *big.Int) (types.Header, error)
	LatestFinalizedBlockByNumber(context.Context, *big.Int) (types.Header, error)

	TxByHash(ctx context.Context, hash common.Hash) (*types.Transactions, error)
	TxReceiptByHash(ctx context.Context, hash common.Hash) (*types.Receipt, error)

	//FilterLogs 事件
	FilterLogs(ctx context.Context, q ethereum.FilterQuery, chainId *big.Int) (Logs, error)
}

type client struct {
	rpc RPC
}

func DialEthClient(ctx context.Context, rpcUrl string) (EthClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	clt, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}
	rpc := NewRPC(clt)
	return &client{
		rpc: rpc,
	}, nil
}

func (c *client) LatestBlockHeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return nil, nil
}

func (c *client) LatestBlockByNumber(context.Context, *big.Int) (types.Block, error) {
	return types.Block{}, nil
}

func (c *client) LatestSafeBlockByNumber(context.Context, *big.Int) (types.Header, error) {
	return types.Header{}, nil
}

func (c *client) LatestFinalizedBlockByNumber(context.Context, *big.Int) (types.Header, error) {
	return types.Header{}, nil
}

func (c *client) TxByHash(ctx context.Context, hash common.Hash) (*types.Transactions, error) {
	return nil, nil
}

func (c *client) TxReceiptByHash(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (c *client) FilterLogs(ctx context.Context, q ethereum.FilterQuery, chainId *big.Int) (Logs, error) {
	return Logs{}, nil
}

type RPC interface {
	Close()
	CallContext(ctx context.Context, result any, method string, args ...any) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

func NewRPC(client *rpc.Client) RPC {
	return &rpcClient{client}
}

type rpcClient struct {
	rpc *rpc.Client
}

func (r rpcClient) Close() {
	r.rpc.Close()
}

func (r rpcClient) CallContext(ctx context.Context, result any, method string, args ...any) error {
	err := r.rpc.CallContext(ctx, result, method, args...)
	return err
}

func (r rpcClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	err := r.rpc.BatchCallContext(ctx, b)
	return err
}
