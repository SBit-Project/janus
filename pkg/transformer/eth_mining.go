package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Project/janus/pkg/sbit"
	"github.com/labstack/echo"
)

//ProxyETHGetHashrate implements ETHProxy
type ProxyETHMining struct {
	*sbit.Sbit
}

func (p *ProxyETHMining) Method() string {
	return "eth_mining"
}

func (p *ProxyETHMining) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHMining) request() (*eth.MiningResponse, eth.JSONRPCError) {
	sbitresp, err := p.Sbit.GetMining()
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.ToResponse(sbitresp), nil
}

func (p *ProxyETHMining) ToResponse(sbitresp *sbit.GetMiningResponse) *eth.MiningResponse {
	ethresp := eth.MiningResponse(sbitresp.Staking)
	return &ethresp
}
