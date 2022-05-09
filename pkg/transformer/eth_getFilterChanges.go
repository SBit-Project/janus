package transformer

import (
	"encoding/json"
	"math/big"

	"github.com/labstack/echo"

	"github.com/SBit-Project/janus/pkg/conversion"
	"github.com/SBit-Projectt/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
)

// ProxyETHGetFilterChanges implements ETHProxy
type ProxyETHGetFilterChanges struct {
	*sbit.Sbit
	filter *eth.FilterSimulator
}

func (p *ProxyETHGetFilterChanges) Method() string {
	return "eth_getFilterChanges"
}

func (p *ProxyETHGetFilterChanges) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {

	filter, err := processFilter(p, rawreq)
	if err != nil {
		return nil, err
	}

	switch filter.Type {
	case eth.NewFilterTy:
		return p.requestFilter(filter)
	case eth.NewBlockFilterTy:
		return p.requestBlockFilter(filter)
	case eth.NewPendingTransactionFilterTy:
		fallthrough
	default:
		return nil, eth.NewInvalidParamsError("Unknown filter type")
	}
}

func (p *ProxyETHGetFilterChanges) requestBlockFilter(filter *eth.Filter) (sbitresp eth.GetFilterChangesResponse, err eth.JSONRPCError) {
	sbitresp = make(eth.GetFilterChangesResponse, 0)

	_lastBlockNumber, ok := filter.Data.Load("lastBlockNumber")
	if !ok {
		return sbitresp, eth.NewCallbackError("Could not get lastBlockNumber")
	}
	lastBlockNumber := _lastBlockNumber.(uint64)

	blockCountBigInt, blockErr := p.GetBlockCount()
	if blockErr != nil {
		return sbitresp, eth.NewCallbackError(blockErr.Error())
	}
	blockCount := blockCountBigInt.Uint64()

	differ := blockCount - lastBlockNumber

	hashes := make(eth.GetFilterChangesResponse, differ)
	for i := range hashes {
		blockNumber := new(big.Int).SetUint64(lastBlockNumber + uint64(i) + 1)

		resp, err := p.GetBlockHash(blockNumber)
		if err != nil {
			return sbitresp, eth.NewCallbackError(err.Error())
		}

		hashes[i] = utils.AddHexPrefix(string(resp))
	}

	sbitresp = hashes
	filter.Data.Store("lastBlockNumber", blockCount)
	return
}

func (p *ProxyETHGetFilterChanges) requestFilter(filter *eth.Filter) (sbitresp eth.GetFilterChangesResponse, err eth.JSONRPCError) {
	sbitresp = make(eth.GetFilterChangesResponse, 0)

	_lastBlockNumber, ok := filter.Data.Load("lastBlockNumber")
	if !ok {
		return sbitresp, eth.NewCallbackError("Could not get lastBlockNumber")
	}
	lastBlockNumber := _lastBlockNumber.(uint64)

	blockCountBigInt, blockErr := p.GetBlockCount()
	if blockErr != nil {
		return sbitresp, eth.NewCallbackError(blockErr.Error())
	}
	blockCount := blockCountBigInt.Uint64()

	differ := blockCount - lastBlockNumber

	if differ == 0 {
		return eth.GetFilterChangesResponse{}, nil
	}

	searchLogsReq, err := p.toSearchLogsReq(filter, big.NewInt(int64(lastBlockNumber+1)), big.NewInt(int64(blockCount)))
	if err != nil {
		return nil, err
	}

	return p.doSearchLogs(searchLogsReq)
}

func (p *ProxyETHGetFilterChanges) doSearchLogs(req *sbit.SearchLogsRequest) (eth.GetFilterChangesResponse, eth.JSONRPCError) {
	resp, err := conversion.SearchLogsAndFilterExtraTopics(p.Sbit, req)
	if err != nil {
		return nil, err
	}

	receiptToResult := func(receipt *sbit.TransactionReceipt) []interface{} {
		logs := conversion.ExtractETHLogsFromTransactionReceipt(receipt, receipt.Log)
		res := make([]interface{}, len(logs))
		for i := range res {
			res[i] = logs[i]
		}
		return res
	}
	results := make(eth.GetFilterChangesResponse, 0)
	for _, receipt := range resp {
		r := sbit.TransactionReceipt(receipt)
		results = append(results, receiptToResult(&r)...)
	}

	return results, nil
}

func (p *ProxyETHGetFilterChanges) toSearchLogsReq(filter *eth.Filter, from, to *big.Int) (*sbit.SearchLogsRequest, eth.JSONRPCError) {
	ethreq := filter.Request.(*eth.NewFilterRequest)
	var err error
	var addresses []string
	if ethreq.Address != nil {
		if isBytesOfString(ethreq.Address) {
			var addr string
			if err = json.Unmarshal(ethreq.Address, &addr); err != nil {
				// TODO: Correct error code?
				return nil, eth.NewInvalidParamsError(err.Error())
			}
			addresses = append(addresses, addr)
		} else {
			if err = json.Unmarshal(ethreq.Address, &addresses); err != nil {
				// TODO: Correct error code?
				return nil, eth.NewInvalidParamsError(err.Error())
			}
		}
		for i := range addresses {
			addresses[i] = utils.RemoveHexPrefix(addresses[i])
		}
	}

	sbitreq := &sbit.SearchLogsRequest{
		Addresses: addresses,
		FromBlock: from,
		ToBlock:   to,
	}

	topics, ok := filter.Data.Load("topics")
	if ok {
		sbitreq.Topics = topics.([]sbit.SearchLogsTopic)
	}

	return sbitreq, nil
}
