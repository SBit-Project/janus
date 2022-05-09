package transformer

import (
	"fmt"
	"testing"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEthValueToSbitAmount(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in":   "0xde0b6b3a7640000",
			"want": decimal.NewFromFloat(1),
		},
		{

			"in":   "0x6f05b59d3b20000",
			"want": decimal.NewFromFloat(0.5),
		},
		{
			"in":   "0x2540be400",
			"want": decimal.NewFromFloat(0.00000001),
		},
		{
			"in":   "0x1",
			"want": decimal.NewFromInt(0),
		},
	}
	for _, c := range cases {
		in := c["in"].(string)
		want := c["want"].(decimal.Decimal)
		got, err := EthValueToSbitAmount(in, MinimumGas)
		if err != nil {
			t.Error(err)
		}

		// TODO: Refactor to use new testing utilities?
		if !got.Equal(want) {
			t.Errorf("in: %s, want: %v, got: %v", in, want, got)
		}
	}
}

func TestSbitValueToEthAmount(t *testing.T) {
	cases := []decimal.Decimal{
		decimal.NewFromFloat(1),
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.00000001),
		MinimumGas,
	}
	for _, c := range cases {
		in := c
		eth := SbitDecimalValueToETHAmount(in)
		out := EthDecimalValueToSbitAmount(eth)

		// TODO: Refactor to use new testing utilities?
		if !in.Equals(out) {
			t.Errorf("in: %s, eth: %v, sbit: %v", in, eth, out)
		}
	}
}

func TestSbitAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.1), "0x16345785d8a0000"
	got, err := formatSbitAmount(in)
	if err != nil {
		t.Error(err)
	}

	internal.CheckTestResultUnspecifiedInputMarshal(in, want, got, t, false)
}

func TestLowestSbitAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.00000001), "0x2540be400"
	got, err := formatSbitAmount(in)
	if err != nil {
		t.Error(err)
	}

	internal.CheckTestResultUnspecifiedInputMarshal(in, want, got, t, false)
}

func TestAddressesConversion(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		sbitChain   string
		ethAddress  string
		sbitAddress string
	}{
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "6c89a1a6ca2ae7c00b248bb2832d6f480f27da68",
			sbitAddress: "qTTH1Yr2eKCuDLqfxUyBLCAjmomQ8pyrBt",
		},

		// NOTE: Ethereum addresses are without `0x` prefix, as it expects by conversion functions
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "7926223070547d2d15b2ef5e7383e541c338ffe9",
			sbitAddress: "sUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
		},
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "2352be3db3177f0a07efbe6da5857615b8c9901d",
			sbitAddress: "sLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf",
		},
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "69b004ac2b3993bf2fdf56b02746a1f57997420d",
			sbitAddress: "qTCCy8qy7pW94EApdoBjYc1vQ2w68UnXPi",
		},
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "8c647515f03daeefd09872d7530fa8d8450f069a",
			sbitAddress: "qWMi6ne9mDQFatRGejxdDYVUV9rQVkAFGp",
		},
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "2191744eb5ebeac90e523a817b77a83a0058003b",
			sbitAddress: "qLcshhsRS6HKeTKRYFdpXnGVZxw96QQcfm",
		},
		{
			sbitChain:   sbit.ChainTest,
			ethAddress:  "88b0bf4b301c21f8a47be2188bad6467ad556dcf",
			sbitAddress: "qW28njWueNpBXYWj2KDmtFG2gbLeALeHfV",
		},
	}

	for i, in := range inputs {
		var (
			in       = in
			testDesc = fmt.Sprintf("#%d", i)
		)
		// TODO: Investigate why this testing setup is so different
		t.Run(testDesc, func(t *testing.T) {
			sbitAddress, err := convertETHAddress(in.ethAddress, in.sbitChain)
			require.NoError(t, err, "couldn't convert Ethereum address to Sbit address")
			require.Equal(t, in.sbitAddress, sbitAddress, "unexpected converted Sbit address value")

			ethAddress, err := utils.ConvertSbitAddress(in.sbitAddress)
			require.NoError(t, err, "couldn't convert Sbit address to Ethereum address")
			require.Equal(t, in.ethAddress, ethAddress, "unexpected converted Ethereum address value")
		})
	}
}

func TestSendTransactionRequestHasDefaultGasPriceAndAmount(t *testing.T) {
	var req eth.SendTransactionRequest
	err := unmarshalRequest([]byte(`[{}]`), &req)
	if err != nil {
		t.Fatal(err)
	}
	defaultGasPriceInWei := req.GasPrice.Int
	defaultGasPriceInSBIT := EthDecimalValueToSbitAmount(decimal.NewFromBigInt(defaultGasPriceInWei, 1))

	// TODO: Refactor to use new testing utilities?
	if !defaultGasPriceInSBIT.Equals(MinimumGas) {
		t.Fatalf("Default gas price does not convert to SBIT minimum gas price, got: %s want: %s", defaultGasPriceInSBIT.String(), MinimumGas.String())
	}
	if eth.DefaultGasAmountForSbit.String() != req.Gas.Int.String() {
		t.Fatalf("Default gas amount does not match expected default, got: %s want: %s", req.Gas.Int.String(), eth.DefaultGasAmountForSbit.String())
	}
}
