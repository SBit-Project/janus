package transformer

import (
	"encoding/json"
	"testing"

	"github.com/SBit-Project/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
)

func initializeProxyETHGetTransactionByBlockNumberAndIndex(sbitClient *sbit.Sbit) ETHProxy {
	return &ProxyETHGetTransactionByBlockNumberAndIndex{sbitClient}
}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockNumberAndIndex,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x0"`)},
		internal.GetTransactionByHashResponseData,
	)
}
