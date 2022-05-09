package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Project/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHNetVersion implements ETHProxy
type ProxyETHNetVersion struct {
	*sbit.Sbit
}

func (p *ProxyETHNetVersion) Method() string {
	return "net_version"
}

func (p *ProxyETHNetVersion) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHNetVersion) request() (*eth.NetVersionResponse, eth.JSONRPCError) {
	networkID, err := getChainId(p.Sbit)
	if err != nil {
		return nil, err
	}
	response := eth.NetVersionResponse(hexutil.EncodeBig(networkID))
	return &response, nil
}
