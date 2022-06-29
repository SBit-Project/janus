package transformer

import (
	"math/big"
	"strings"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

type ProxyETHChainId struct {
	*sbit.Sbit
}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	chainId, err := getChainId(p.Sbit)
	if err != nil {
		return nil, err
	}
	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}

func getChainId(p *sbit.Sbit) (*big.Int, eth.JSONRPCError) {
	var sbitresp *sbit.GetBlockChainInfoResponse
	if err := p.Request(sbit.MethodGetBlockChainInfo, nil, &sbitresp); err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	var chainId *big.Int
	switch strings.ToLower(sbitresp.Chain) {
	case "main":
		chainId = big.NewInt(81)
	case "test":
		chainId = big.NewInt(8889)
	case "regtest":
		chainId = big.NewInt(8890)
	default:
		chainId = big.NewInt(8890)
		p.GetDebugLogger().Log("msg", "Unknown chain "+sbitresp.Chain)
	}

	return chainId, nil
}
