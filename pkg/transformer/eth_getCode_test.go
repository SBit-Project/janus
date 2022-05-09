package transformer

import (
	"encoding/json"
	"testing"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Project/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/btcsuite/btcutil"
)

func TestGetAccountInfoRequest(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	sbitClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	sbitClient.Accounts = append(sbitClient.Accounts, account)

	//prepare responses
	getAccountInfoResponse := sbit.GetAccountInfoResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Balance: 12431243,
		// Storage json.RawMessage `json:"storage"`,
		Code: "606060405236156100ad576000357c0100000000000000000...",
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, sbit.MethodGetAccountInfo, getAccountInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetCode{sbitClient}
	got, jsonErr := proxyEth.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetCodeResponse("0x606060405236156100ad576000357c0100000000000000000...")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}

func TestGetCodeInvalidAddressRequest(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x0000000000000000000000000000000000000000"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	sbitClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare responses
	getAccountInfoErrorResponse := sbit.GetErrorResponse(sbit.ErrInvalidAddress)
	if getAccountInfoErrorResponse == nil {
		panic("mocked error response is nil")
	}
	err = mockedClientDoer.AddError(sbit.MethodGetAccountInfo, getAccountInfoErrorResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetCode{sbitClient}
	got, jsonErr := proxyEth.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetCodeResponse("0x")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}
