package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/labstack/echo"
)

type ProxySBITGenericStringArguments struct {
	*sbit.Sbit
	prefix string
	method string
}

var _ ETHProxy = (*ProxySBITGenericStringArguments)(nil)

func (p *ProxySBITGenericStringArguments) Method() string {
	return p.prefix + "_" + p.method
}

func (p *ProxySBITGenericStringArguments) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var params eth.StringsArguments
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("couldn't unmarshal request parameters")
	}

	if len(params) != 1 {
		return nil, eth.NewInvalidParamsError("require 1 argument: the base58 Sbit address")
	}

	return p.request(params)
}

func (p *ProxySBITGenericStringArguments) request(params eth.StringsArguments) (*string, eth.JSONRPCError) {
	var response string
	err := p.Client.Request(p.method, params, &response)
	if err != nil {
		return nil, eth.NewInvalidRequestError(err.Error())
	}

	return &response, nil
}
