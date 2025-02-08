package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/common/global_const"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"sync"
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
	ToBlockHeader *types.Header // 标记日志信息的来源
}

type TransactionList struct {
	To   string `json:"to"`
	Hash string `json:"hash"`
}

type RpcBlock struct {
	Hash         common.Hash       `json:"hash"`
	Transactions []TransactionList `json:"transactions"`
	BaseFee      string            `json:"baseFeePerGas"`
}

type EthClient interface {
	BlockHeaderByNumber(*big.Int) (*types.Header, error)

	BlockByNumber(*big.Int) (*RpcBlock, error)

	LatestSafeBlockHeader() (*types.Header, error)
	LatestFinalizedBlockHeader() (*types.Header, error)
	BlockHeaderByHash(common.Hash) (*types.Header, error)
	BlockHeadersByRange(*big.Int, *big.Int, uint) ([]types.Header, error)

	TxByHash(common.Hash) (*types.Transaction, error)
	TxReceiptByHash(common.Hash) (*types.Receipt, error)

	StorageHash(common.Address, *big.Int) (common.Hash, error)
	FilterLogs(filterQuery ethereum.FilterQuery, chainId uint) (Logs, error)

	TxCountByAddress(common.Address) (hexutil.Uint64, error)

	SendRawTransaction(rawTx string) error

	SuggestGasPrice() (*big.Int, error)
	SuggestGasTipCap() (*big.Int, error)

	Close()
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

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	return rpc.BlockNumber(number.Int64()).String()
}

func (c *client) BlockHeaderByNumber(number *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		log.Error("Call eth_getBlockByNumber method fail", "err", err)
		return nil, err
	} else if header == nil {
		log.Warn("header not found")
		return nil, ethereum.NotFound
	}
	return header, nil
}

func (c *client) BlockByNumber(number *big.Int) (*RpcBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var block *RpcBlock
	err := c.rpc.CallContext(ctx, &block, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		log.Error("Call eth_getBlockByNumber method fail", "err", err)
		return nil, err
	} else if block == nil {
		log.Warn("header not found")
		return nil, ethereum.NotFound
	}

	return block, nil
}

func (c *client) LatestSafeBlockHeader() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", "safe", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	return header, nil
}

func (c *client) LatestFinalizedBlockHeader() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	return header, nil
}

func (c *client) BlockHeaderByHash(hash common.Hash) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	if header.Hash() != hash {
		return nil, errors.New("header mismatch")
	}

	return header, nil
}

// BlockHeadersByRange 根据起始区块高度 startHeight 和结束区块高度 endHeight 批量获取区块头信息。
// 如果 startHeight == endHeight，则直接查询单个区块头。
// 对于 ZkFairSepolia 和 ZkFair 链，采用并发分批查询（每批最多 100 个区块）。
func (c *client) BlockHeadersByRange(startHeight, endHeight *big.Int, chainId uint) ([]types.Header, error) {
	// 比较起始块儿和总止块儿是否一样
	if startHeight.Cmp(endHeight) == 0 {
		header, err := c.BlockHeaderByNumber(startHeight)
		if err != nil {
			return nil, err
		}
		return []types.Header{*header}, nil
	}

	// 计算需要查询的区块数量
	count := new(big.Int).Sub(endHeight, startHeight).Uint64() + 1
	// 预分配存储区块头的切片
	headers := make([]types.Header, count)
	// 预分配 RPC 批量查询的请求切片
	batchElems := make([]rpc.BatchElem, count)

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	// 如果目标链是 ZkFairSepolia 或 ZkFair，则采用并发分批查询的方式
	if chainId == uint(global_const.ZkFairSepoliaChainId) ||
		chainId == uint(global_const.ZkFairChainId) {
		groupSize := 100 // 每批最多查询 100 个区块
		var wg sync.WaitGroup
		numGroups := (int(count)-1)/groupSize + 1 // 计算总批次数
		wg.Add(numGroups)

		// 以 groupSize 为单位进行分批
		for i := 0; i < int(count); i += groupSize {
			start := i
			end := i + groupSize - 1
			if end > int(count) {
				end = int(count) - 1 // 防止越界
			}

			// 启动并发协程查询每个批次的区块头信息
			go func(start, end int) {
				defer wg.Done()
				for j := start; j <= end; j++ {
					// 计算当前区块高度
					height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(uint64(j)))
					// 构造 RPC 请求
					batchElems[j] = rpc.BatchElem{
						Method: "eth_getBlockByNumber",
						Result: new(types.Header),
						Error:  nil,
					}
					// 直接执行 RPC 调用（同步方式）
					header := new(types.Header)
					batchElems[j].Error = c.rpc.CallContext(ctx, header, batchElems[j].Method, toBlockNumArg(height), false)
					batchElems[j].Result = header
				}
			}(start, end)
		}
		// 等待所有批次执行完毕
		wg.Wait()
	} else {
		// 普通情况，使用批量 RPC 查询（一次性请求所有区块头）
		for i := uint64(0); i < count; i++ {
			height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(i))
			// 构造批量查询请求
			batchElems[i] = rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{toBlockNumArg(height), false},
				Result: &headers[i],
			}
		}
		// 执行批量 RPC 调用
		err := c.rpc.BatchCallContext(ctx, batchElems)
		if err != nil {
			return nil, err
		}

	}

	// 解析批量查询结果
	size := 0
	for i, batchElem := range batchElems {
		// 确保 RPC 响应结果能够正确转换为 `types.Header`
		header, ok := batchElem.Result.(*types.Header)
		if !ok {
			return nil, fmt.Errorf("unable to transform rpc response %v into types.Header", batchElem.Result)
		}
		headers[i] = *header
		size = size + 1
	}

	// 根据实际获取的区块头数量调整切片大小
	headers = headers[:size]

	return headers, nil

}

func (c *client) TxByHash(hash common.Hash) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var tx *types.Transaction
	err := c.rpc.CallContext(ctx, &tx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, ethereum.NotFound
	}

	return tx, nil
}

func (c *client) TxReceiptByHash(hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var txReceipt *types.Receipt
	err := c.rpc.CallContext(ctx, &txReceipt, "eth_getTransactionReceipt", hash)

	if err != nil {
		return nil, err
	} else if txReceipt == nil {
		return nil, ethereum.NotFound
	}

	return txReceipt, nil
}

func (c *client) StorageHash(address common.Address, blockNumber *big.Int) (common.Hash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	proof := struct{ StorageHash common.Hash }{}

	err := c.rpc.CallContext(ctx, &proof, "eth_getProof", address, nil, toBlockNumArg(blockNumber))
	if err != nil {
		return common.Hash{}, err
	}

	return proof.StorageHash, nil

}

// toFilterArg 将 `ethereum.FilterQuery` 结构体转换为适用于 RPC 调用的过滤参数
func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	// 初始化参数，包含 `address` 和 `topics`
	arg := map[string]interface{}{
		"address": q.Addresses, // 监听的合约地址（可以是多个）
		"topics":  q.Topics,    // 过滤的事件主题（可以是多个）
	}

	// 如果指定了 `BlockHash`，则必须仅使用 `BlockHash` 进行过滤，而不能使用 `FromBlock` 或 `ToBlock`
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, errors.New("cannot specify both BlockHash and FromBlock/ToBlock") // 避免冲突
		}
	} else {
		// 处理 `FromBlock` 和 `ToBlock`，指定区块范围
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0" // 默认从创世区块开始
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock) // 指定 `ToBlock`
	}

	return arg, nil
}

// FilterLogs 根据指定的过滤条件查询区块链日志
func (c *client) FilterLogs(query ethereum.FilterQuery, chainId uint) (Logs, error) {
	// 将查询条件转换为 RPC 需要的参数
	arg, err := toFilterArg(query)
	if err != nil {
		return Logs{}, err
	}

	var logs []types.Log    // 用于存储查询到的日志
	var header types.Header // 用于存储 `ToBlock` 区块的头信息

	// 构造批量请求的元素，包括：
	// 1. 查询 `ToBlock` 对应的区块头信息
	// 2. 查询符合条件的日志
	batchElems := make([]rpc.BatchElem, 2)
	batchElems[0] = rpc.BatchElem{Method: "eth_getBlockByNumber", Args: []interface{}{toBlockNumArg(query.ToBlock), false}, Result: &header}
	batchElems[1] = rpc.BatchElem{Method: "eth_getLogs", Args: []interface{}{arg}, Result: &logs}

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	// 对于特定链（ZkFairSepolia 和 ZkFair），使用单独的 RPC 请求，而不是批量调用
	if chainId == uint(global_const.ZkFairSepoliaChainId) ||
		chainId == uint(global_const.ZkFairChainId) {

		// 分别查询区块头信息和日志数据
		batchElems[0].Error = c.rpc.CallContext(ctx, &header, batchElems[0].Method, toBlockNumArg(query.ToBlock), false)
		batchElems[1].Error = c.rpc.CallContext(ctx, &logs, batchElems[1].Method, arg)
	} else {
		// 其他链使用批量调用的方式，减少 RPC 请求次数，提高效率
		err = c.rpc.BatchCallContext(ctx, batchElems)
		if err != nil {
			return Logs{}, err
		}
	}

	// 检查查询 `ToBlock` 区块头信息是否出错
	if batchElems[0].Error != nil {
		return Logs{}, fmt.Errorf("unable to query for the `FilterQuery#ToBlock` header: %w", batchElems[0].Error)
	}
	// 检查查询日志是否出错
	if batchElems[1].Error != nil {
		return Logs{}, fmt.Errorf("unable to query logs: %w", batchElems[1].Error)
	}

	// 返回查询到的日志和区块头信息
	return Logs{Logs: logs, ToBlockHeader: &header}, nil
}

func (c *client) TxCountByAddress(address common.Address) (hexutil.Uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var nonce hexutil.Uint64
	err := c.rpc.CallContext(ctx, &nonce, "eth_getTransactionCount", address, "latest")
	if err != nil {
		log.Error("Call eth_getTransactionCount method fail", "err", err)
		return 0, err
	}
	log.Info("get nonce by address success", "nonce", nonce)
	return nonce, err
}

func (c *client) SendRawTransaction(rawTx string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	if err := c.rpc.CallContext(ctx, nil, "eth_sendRawTransaction", rawTx); err != nil {
		return err
	}
	log.Info("send tx to ethereum success")
	return nil
}

func (c *client) SuggestGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var hex hexutil.Big
	if err := c.rpc.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *client) SuggestGasTipCap() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var hex hexutil.Big
	if err := c.rpc.CallContext(ctx, &hex, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *client) Close() {
	c.rpc.Close()
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
