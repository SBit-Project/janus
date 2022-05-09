package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/labstack/echo"
)

// ProxyETHPersonalUnlockAccount implements ETHProxy
type ProxyETHPersonalUnlockAccount struct{}

func (p *ProxyETHPersonalUnlockAccount) Method() string {
	return "personal_unlockAccount"
}

func (p *ProxyETHPersonalUnlockAccount) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return eth.PersonalUnlockAccountResponse(true), nil
}
