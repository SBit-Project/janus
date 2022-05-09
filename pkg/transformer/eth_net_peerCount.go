package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/dcb9/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyNetPeerCount implements ETHProxy
type ProxyNetPeerCount struct {
	*sbit.Sbit
}

func (p *ProxyNetPeerCount) Method() string {
	return "net_peerCount"
}

func (p *ProxyNetPeerCount) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyNetPeerCount) request() (*eth.NetPeerCountResponse, eth.JSONRPCError) {
	peerInfos, err := p.GetPeerInfo()
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	resp := eth.NetPeerCountResponse(hexutil.EncodeUint64(uint64(len(peerInfos))))
	return &resp, nil
}
