package transformer

import (
	"math/big"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHTxCount struct {
	*sbit.Sbit
}

func (p *ProxyETHTxCount) Method() string {
	return "eth_getTransactionCount"
}

func (p *ProxyETHTxCount) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {

	/* not sure we need this. Need to figure out how to best unmarshal this in the future. For now this will work.
	var req eth.GetTransactionCountRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}*/
	sbitresp, err := p.Sbit.GetTransactionCount("", "")
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.response(sbitresp), nil
}

func (p *ProxyETHTxCount) response(sbitresp *big.Int) string {
	return hexutil.EncodeBig(sbitresp)
}
