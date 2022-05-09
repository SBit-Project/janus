package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHNewBlockFilter implements ETHProxy
type ProxyETHNewBlockFilter struct {
	*sbit.Sbit
	filter *eth.FilterSimulator
}

func (p *ProxyETHNewBlockFilter) Method() string {
	return "eth_newBlockFilter"
}

func (p *ProxyETHNewBlockFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHNewBlockFilter) request() (eth.NewBlockFilterResponse, eth.JSONRPCError) {
	blockCount, err := p.GetBlockCount()
	if err != nil {
		return "", eth.NewCallbackError(err.Error())
	}

	filter := p.filter.New(eth.NewBlockFilterTy)
	filter.Data.Store("lastBlockNumber", blockCount.Uint64())

	p.GenerateIfPossible()

	return eth.NewBlockFilterResponse(hexutil.EncodeUint64(filter.ID)), nil
}
