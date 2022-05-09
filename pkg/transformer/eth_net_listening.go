package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Project/janus/pkg/sbit"
	"github.com/labstack/echo"
)

// ProxyETHGetCode implements ETHProxy
type ProxyNetListening struct {
	*sbit.Sbit
}

func (p *ProxyNetListening) Method() string {
	return "net_listening"
}

func (p *ProxyNetListening) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	networkInfo, err := p.GetNetworkInfo()
	if err != nil {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "Failed to query network info", "err", err)
		return false, eth.NewCallbackError(err.Error())
	}

	p.GetDebugLogger().Log("method", p.Method(), "network active", networkInfo.NetworkActive)
	return networkInfo.NetworkActive, nil
}
