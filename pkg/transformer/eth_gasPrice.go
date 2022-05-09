package transformer

import (
	"math/big"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHGasPrice struct {
	*sbit.Sbit
}

func (p *ProxyETHGasPrice) Method() string {
	return "eth_gasPrice"
}

func (p *ProxyETHGasPrice) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	sbitresp, err := p.Sbit.GetGasPrice()
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.response(sbitresp), nil
}

func (p *ProxyETHGasPrice) response(sbitresp *big.Int) string {
	// 34 GWEI is the minimum price that SBIT will confirm tx with
	return hexutil.EncodeBig(convertFromSatoshiToWei(sbitresp))
}
