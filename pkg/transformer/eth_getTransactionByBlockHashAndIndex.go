package transformer

import (
	"encoding/json"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHGetTransactionByBlockHashAndIndex implements ETHProxy
type ProxyETHGetTransactionByBlockHashAndIndex struct {
	*sbit.Sbit
}

func (p *ProxyETHGetTransactionByBlockHashAndIndex) Method() string {
	return "eth_getTransactionByBlockHashAndIndex"
}

func (p *ProxyETHGetTransactionByBlockHashAndIndex) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetTransactionByBlockHashAndIndex
	if err := json.Unmarshal(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}
	if req.BlockHash == "" {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("invalid argument 0: empty hex string")
	}

	return p.request(&req)
}

func (p *ProxyETHGetTransactionByBlockHashAndIndex) request(req *eth.GetTransactionByBlockHashAndIndex) (interface{}, eth.JSONRPCError) {
	transactionIndex, err := hexutil.DecodeUint64(req.TransactionIndex)
	if err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("invalid argument 1")
	}

	// Proxy eth_getBlockByHash and return the transaction at requested index
	getBlockByNumber := ProxyETHGetBlockByHash{p.Sbit}
	blockByNumber, jsonErr := getBlockByNumber.request(&eth.GetBlockByHashRequest{BlockHash: req.BlockHash, FullTransaction: true})

	if jsonErr != nil {
		return nil, jsonErr
	}

	if blockByNumber == nil {
		return nil, nil
	}

	if len(blockByNumber.Transactions) <= int(transactionIndex) {
		return nil, nil
	}

	return blockByNumber.Transactions[int(transactionIndex)], nil
}
