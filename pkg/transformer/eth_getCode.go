package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHGetCode implements ETHProxy
type ProxyETHGetCode struct {
	*sbit.Sbit
}

func (p *ProxyETHGetCode) Method() string {
	return "eth_getCode"
}

func (p *ProxyETHGetCode) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetCodeRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req)
}

func (p *ProxyETHGetCode) request(ethreq *eth.GetCodeRequest) (eth.GetCodeResponse, eth.JSONRPCError) {
	sbitreq := sbit.GetAccountInfoRequest(utils.RemoveHexPrefix(ethreq.Address))

	sbitresp, err := p.GetAccountInfo(&sbitreq)
	if err != nil {
		if err == sbit.ErrInvalidAddress {
			/**
			// correct response for an invalid address
			{
				"jsonrpc": "2.0",
				"id": 123,
				"result": "0x"
			}
			**/
			return "0x", nil
		} else {
			return "", eth.NewCallbackError(err.Error())
		}
	}

	// sbit res -> eth res
	return eth.GetCodeResponse(utils.AddHexPrefix(sbitresp.Code)), nil
}
