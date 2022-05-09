package transformer

import (
	"encoding/json"
	"testing"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
)

func TestGetTransactionReceiptForNonVMTransaction(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{[]byte(`"0x8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950"`)}
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
	err = mockedClientDoer.AddResponseWithRequestID(2, sbit.MethodGetTransactionReceipt, []byte("[]"))
	if err != nil {
		t.Fatal(err)
	}

	rawTransactionResponse := &sbit.GetRawTransactionResponse{
		BlockHash: internal.GetTransactionByHashBlockHash,
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, sbit.MethodGetRawTransaction, rawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	err = mockedClientDoer.AddResponseWithRequestID(4, sbit.MethodGetBlock, internal.GetBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetTransactionReceipt{sbitClient}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetTransactionReceiptResponse{
		TransactionHash:   "0x8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950",
		TransactionIndex:  "0x1",
		BlockHash:         "0xbba11e1bacc69ba535d478cf1f2e542da3735a517b0b8eebaf7e6bb25eeb48c5",
		BlockNumber:       "0xf8f",
		GasUsed:           NonContractVMGasLimit,
		Logs:              []eth.Log{},
		EffectiveGasPrice: "0x0",
		CumulativeGasUsed: NonContractVMGasLimit,
		To:                utils.AddHexPrefix(sbit.ZeroAddress),
		From:              utils.AddHexPrefix(sbit.ZeroAddress),
		LogsBloom:         eth.EmptyLogsBloom,
		Status:            STATUS_SUCCESS,
	}

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
