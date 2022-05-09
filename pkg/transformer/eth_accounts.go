package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHAccounts implements ETHProxy
type ProxyETHAccounts struct {
	*sbit.Sbit
}

func (p *ProxyETHAccounts) Method() string {
	return "eth_accounts"
}

func (p *ProxyETHAccounts) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHAccounts) request() (eth.AccountsResponse, eth.JSONRPCError) {
	var accounts eth.AccountsResponse

	for _, acc := range p.Accounts {
		acc := sbit.Account{acc}
		addr := acc.ToHexAddress()

		accounts = append(accounts, utils.AddHexPrefix(addr))
	}

	return accounts, nil
}

func (p *ProxyETHAccounts) ToResponse(ethresp *sbit.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	sbitresp := eth.CallResponse(data)
	return &sbitresp
}
