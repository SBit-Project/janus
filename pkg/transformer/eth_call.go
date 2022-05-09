package transformer

import (
	"math/big"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*sbit.Sbit
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Is this correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req)
}

func (p *ProxyETHCall) request(ethreq *eth.CallRequest) (interface{}, eth.JSONRPCError) {
	// eth req -> sbit req
	sbitreq, jsonErr := p.ToRequest(ethreq)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if sbitreq.GasLimit != nil && sbitreq.GasLimit.Cmp(big.NewInt(40000000)) > 0 {
		sbitresp := eth.CallResponse("0x")
		p.Sbit.GetLogger().Log("msg", "Caller gas above allowance, capping", "requested", sbitreq.GasLimit.Int64(), "cap", "40,000,000")
		return &sbitresp, nil
	}

	sbitresp, err := p.CallContract(sbitreq)
	if err != nil {
		if err == sbit.ErrInvalidAddress {
			sbitresp := eth.CallResponse("0x")
			return &sbitresp, nil
		}

		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.ToResponse(sbitresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*sbit.CallContractRequest, eth.JSONRPCError) {
	from := ethreq.From
	var err error
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, eth.NewCallbackError(err.Error())
		}
	}

	var gasLimit *big.Int
	if ethreq.Gas != nil {
		gasLimit = ethreq.Gas.Int
	}

	if gasLimit != nil && gasLimit.Int64() < MinimumGasLimit {
		p.GetLogger().Log("msg", "Gas limit is too low", "gasLimit", gasLimit.String())
	}

	return &sbit.CallContractRequest{
		To:       ethreq.To,
		From:     from,
		Data:     ethreq.Data,
		GasLimit: gasLimit,
	}, nil
}

func (p *ProxyETHCall) ToResponse(qresp *sbit.CallContractResponse) interface{} {
	if qresp.ExecutionResult.Output == "" {
		return eth.NewJSONRPCError(
			-32000,
			"Revert: executionResult output is empty",
			nil,
		)
	}

	data := utils.AddHexPrefix(qresp.ExecutionResult.Output)
	sbitresp := eth.CallResponse(data)
	return &sbitresp

}
