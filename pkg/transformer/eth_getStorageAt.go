package transformer

import (
	"fmt"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHGetStorageAt implements ETHProxy
type ProxyETHGetStorageAt struct {
	*sbit.Sbit
}

func (p *ProxyETHGetStorageAt) Method() string {
	return "eth_getStorageAt"
}

func (p *ProxyETHGetStorageAt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetStorageRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	sbitAddress := utils.RemoveHexPrefix(req.Address)
	blockNumber, err := getBlockNumberByParam(p.Sbit, req.BlockNumber, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", fmt.Sprintf("Failed to get block number by param for '%s'", req.BlockNumber), "err", err)
		return nil, err
	}

	return p.request(&sbit.GetStorageRequest{
		Address:     sbitAddress,
		BlockNumber: blockNumber,
	}, utils.RemoveHexPrefix(req.Index))
}

func (p *ProxyETHGetStorageAt) request(ethreq *sbit.GetStorageRequest, index string) (*eth.GetStorageResponse, eth.JSONRPCError) {
	sbitresp, err := p.Sbit.GetStorage(ethreq)
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// sbit res -> eth res
	return p.ToResponse(sbitresp, index), nil
}

func (p *ProxyETHGetStorageAt) ToResponse(sbitresp *sbit.GetStorageResponse, slot string) *eth.GetStorageResponse {
	// the value for unknown anything
	storageData := eth.GetStorageResponse("0x0000000000000000000000000000000000000000000000000000000000000000")
	if len(slot) != 64 {
		slot = leftPadStringWithZerosTo64Bytes(slot)
	}
	for _, outerValue := range *sbitresp {
		sbitStorageData, ok := outerValue[slot]
		if ok {
			storageData = eth.GetStorageResponse(utils.AddHexPrefix(sbitStorageData))
			return &storageData
		}
	}

	return &storageData
}

// left pad a string with leading zeros to fit 64 bytes
func leftPadStringWithZerosTo64Bytes(hex string) string {
	return fmt.Sprintf("%064v", hex)
}
