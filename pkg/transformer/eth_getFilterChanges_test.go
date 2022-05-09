package transformer

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
)

func TestGetFilterChangesRequest_EmptyResult(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1"`)}
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

	//preparing client response
	getBlockCountResponse := sbit.GetBlockCountResponse{Int: big.NewInt(657660)}
	err = mockedClientDoer.AddResponseWithRequestID(2, sbit.MethodGetBlockCount, getBlockCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	searchLogsResponse := sbit.SearchLogsResponse{
		//TODO: add
	}
	err = mockedClientDoer.AddResponseWithRequestID(2, sbit.MethodSearchLogs, searchLogsResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing filter
	filterSimulator := eth.NewFilterSimulator()
	filterRequest := eth.NewFilterRequest{}
	filterSimulator.New(eth.NewFilterTy, &filterRequest)
	_filter, _ := filterSimulator.Filter(1)
	filter := _filter.(*eth.Filter)
	filter.Data.Store("lastBlockNumber", uint64(657655))

	//preparing proxy & executing request
	proxyEth := ProxyETHGetFilterChanges{sbitClient, filterSimulator}
	got, jsonErr := proxyEth.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetFilterChangesResponse{}

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}

func TestGetFilterChangesRequest_NoNewBlocks(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1"`)}
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

	//preparing client response
	getBlockCountResponse := sbit.GetBlockCountResponse{Int: big.NewInt(657655)}
	err = mockedClientDoer.AddResponseWithRequestID(2, sbit.MethodGetBlockCount, getBlockCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing filter
	filterSimulator := eth.NewFilterSimulator()
	filterSimulator.New(eth.NewFilterTy, nil)
	_filter, _ := filterSimulator.Filter(1)
	filter := _filter.(*eth.Filter)
	filter.Data.Store("lastBlockNumber", uint64(657655))

	//preparing proxy & executing request
	proxyEth := ProxyETHGetFilterChanges{sbitClient, filterSimulator}
	got, jsonErr := proxyEth.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetFilterChangesResponse{}

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}

func TestGetFilterChangesRequest_NoSuchFilter(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1"`)}
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

	//preparing proxy & executing request
	filterSimulator := eth.NewFilterSimulator()
	proxyEth := ProxyETHGetFilterChanges{sbitClient, filterSimulator}
	_, got := proxyEth.Request(requestRPC, nil)

	want := eth.NewCallbackError("Invalid filter id")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}
