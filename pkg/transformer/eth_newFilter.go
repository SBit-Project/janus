package transformer

import (
	"encoding/json"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/dcb9/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHNewFilter implements ETHProxy
type ProxyETHNewFilter struct {
	*sbit.Sbit
	filter *eth.FilterSimulator
}

func (p *ProxyETHNewFilter) Method() string {
	return "eth_newFilter"
}

func (p *ProxyETHNewFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.NewFilterRequest
	if err := json.Unmarshal(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req)
}

func (p *ProxyETHNewFilter) request(ethreq *eth.NewFilterRequest) (*eth.NewFilterResponse, eth.JSONRPCError) {

	from, err := getBlockNumberByRawParam(p.Sbit, ethreq.FromBlock, true)
	if err != nil {
		return nil, err
	}

	to, err := getBlockNumberByRawParam(p.Sbit, ethreq.ToBlock, true)
	if err != nil {
		return nil, err
	}

	filter := p.filter.New(eth.NewFilterTy, ethreq)
	filter.Data.Store("lastBlockNumber", from.Uint64())

	filter.Data.Store("toBlock", to.Uint64())

	if len(ethreq.Topics) > 0 {
		topics, err := eth.TranslateTopics(ethreq.Topics)
		if err != nil {
			return nil, eth.NewCallbackError(err.Error())
		}
		filter.Data.Store("topics", sbit.NewSearchLogsTopics(topics))
	}
	resp := eth.NewFilterResponse(hexutil.EncodeUint64(filter.ID))
	return &resp, nil
}
