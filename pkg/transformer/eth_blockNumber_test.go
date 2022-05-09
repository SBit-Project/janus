package transformer

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
)

func TestBlockNumberRequest(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	sbitClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing client response
	getBlockCountResponse := sbit.GetBlockCountResponse{Int: big.NewInt(11284900)}
	err = mockedClientDoer.AddResponseWithRequestID(2, sbit.MethodGetBlockCount, getBlockCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHBlockNumber{sbitClient}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.BlockNumberResponse("0xac31a4")

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
