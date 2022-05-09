package transformer

import (
	"encoding/json"
	"testing"

	"github.com/SBit-Project/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
)

type ETHProxyInitializer = func(*sbit.Sbit) ETHProxy

func testETHProxyRequest(t *testing.T, initializer ETHProxyInitializer, requestParams []json.RawMessage, want interface{}) {
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	sbitClient, err := internal.CreateMockedClient(mockedClientDoer)

	internal.SetupGetBlockByHashResponses(t, mockedClientDoer)

	//preparing proxy & executing request
	proxyEth := initializer(sbitClient)
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatalf("Failed to process request on %T.Request(%s): %s", proxyEth, requestParams, jsonErr)
	}

	internal.CheckTestResultEthRequestRPC(*request, want, got, t, false)
}
