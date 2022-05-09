package transformer

import (
	"runtime"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/params"
	"github.com/labstack/echo"
)

// Web3ClientVersion implements web3_clientVersion
type Web3ClientVersion struct {
	// *sbit.Sbit
}

func (p *Web3ClientVersion) Method() string {
	return "web3_clientVersion"
}

func (p *Web3ClientVersion) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return "Janus/" + params.VersionWithGitSha + "/" + runtime.GOOS + "-" + runtime.GOARCH + "/" + runtime.Version(), nil
}

// func (p *Web3ClientVersion) ToResponse(ethresp *sbit.CallContractResponse) *eth.CallResponse {
// 	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
// 	sbitresp := eth.CallResponse(data)
// 	return &sbitresp
// }
