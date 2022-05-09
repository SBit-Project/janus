package transformer

import (
	"math/big"

	"github.com/SBit-Project-ethereum/common/hexutil"
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Project/janus/pkg/sbit"
	"github.com/SBit-Project/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHGetBalance implements ETHProxy
type ProxyETHGetBalance struct {
	*sbit.Sbit
}

func (p *ProxyETHGetBalance) Method() string {
	return "eth_getBalance"
}

func (p *ProxyETHGetBalance) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetBalanceRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	addr := utils.RemoveHexPrefix(req.Address)
	{
		// is address a contract or an account?
		sbitreq := sbit.GetAccountInfoRequest(addr)
		sbitresp, err := p.GetAccountInfo(&sbitreq)

		// the address is a contract
		if err == nil {
			// the unit of the balance Satoshi
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "is a contract")
			return hexutil.EncodeUint64(uint64(sbitresp.Balance)), nil
		}
	}

	{
		// try account
		base58Addr, err := p.FromHexAddress(addr)
		if err != nil {
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "error parsing address", "error", err)
			return nil, eth.NewCallbackError(err.Error())
		}

		sbitreq := sbit.GetAddressBalanceRequest{Address: base58Addr}
		sbitresp, err := p.GetAddressBalance(&sbitreq)
		if err != nil {
			if err == sbit.ErrInvalidAddress {
				// invalid address should return 0x0
				return "0x0", nil
			}
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "error getting address balance", "error", err)
			return nil, eth.NewCallbackError(err.Error())
		}

		// 1 SBIT = 10 ^ 8 Satoshi
		balance := new(big.Int).SetUint64(sbitresp.Balance)

		//Balance for ETH response is represented in Weis (1 SBIT Satoshi = 10 ^ 10 Wei)
		balance = balance.Mul(balance, big.NewInt(10000000000))

		return hexutil.EncodeBig(balance), nil
	}
}
