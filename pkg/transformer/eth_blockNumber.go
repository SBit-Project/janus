package transformer

import (
	"strings"
	"time"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*sbit.Sbit
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c, 5)
}

func (p *ProxyETHBlockNumber) request(c echo.Context, retries int) (*eth.BlockNumberResponse, eth.JSONRPCError) {
	sbitresp, err := p.Sbit.GetBlockCount()
	if err != nil {
		if retries > 0 && strings.Contains(err.Error(), sbit.ErrTryAgain.Error()) {
			ctx := c.Request().Context()
			t := time.NewTimer(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				return nil, eth.NewCallbackError(err.Error())
			case <-t.C:
				// fallthrough
			}
			return p.request(c, retries-1)
		}
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.ToResponse(sbitresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(sbitresp *sbit.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(sbitresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
