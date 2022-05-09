package transformer

import (
	"math"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

//ProxyETHGetHashrate implements ETHProxy
type ProxyETHHashrate struct {
	*sbit.Sbit
}

func (p *ProxyETHHashrate) Method() string {
	return "eth_hashrate"
}

func (p *ProxyETHHashrate) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHHashrate) request() (*eth.HashrateResponse, eth.JSONRPCError) {
	sbitresp, err := p.Sbit.GetHashrate()
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.ToResponse(sbitresp), nil
}

func (p *ProxyETHHashrate) ToResponse(sbitresp *sbit.GetHashrateResponse) *eth.HashrateResponse {
	hexVal := hexutil.EncodeUint64(math.Float64bits(sbitresp.Difficulty))
	ethresp := eth.HashrateResponse(hexVal)
	return &ethresp
}
