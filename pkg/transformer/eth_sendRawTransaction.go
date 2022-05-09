package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHSendRawTransaction implements ETHProxy
type ProxyETHSendRawTransaction struct {
	*sbit.Sbit
}

var _ ETHProxy = (*ProxyETHSendRawTransaction)(nil)

func (p *ProxyETHSendRawTransaction) Method() string {
	return "eth_sendRawTransaction"
}

func (p *ProxyETHSendRawTransaction) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var params eth.SendRawTransactionRequest
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}
	if params[0] == "" {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("invalid parameter: raw transaction hexed string is empty")
	}

	return p.request(params)
}

func (p *ProxyETHSendRawTransaction) request(params eth.SendRawTransactionRequest) (eth.SendRawTransactionResponse, eth.JSONRPCError) {
	var (
		sbitHexedRawTx = utils.RemoveHexPrefix(params[0])
		req            = sbit.SendRawTransactionRequest([1]string{sbitHexedRawTx})
	)

	sbitresp, err := p.Sbit.SendRawTransaction(&req)
	if err != nil {
		if err == sbit.ErrVerifyAlreadyInChain {
			// already committed
			// we need to send back the tx hash
			rawTx, err := p.Sbit.DecodeRawTransaction(sbitHexedRawTx)
			if err != nil {
				p.GetErrorLogger().Log("msg", "Error decoding raw transaction for duplicate raw transaction", "err", err)
				return eth.SendRawTransactionResponse(""), eth.NewCallbackError(err.Error())
			}
			sbitresp = &sbit.SendRawTransactionResponse{Result: rawTx.Hash}
		} else {
			return eth.SendRawTransactionResponse(""), eth.NewCallbackError(err.Error())
		}
	} else {
		p.GenerateIfPossible()
	}

	resp := *sbitresp
	ethHexedTxHash := utils.AddHexPrefix(resp.Result)
	return eth.SendRawTransactionResponse(ethHexedTxHash), nil
}
