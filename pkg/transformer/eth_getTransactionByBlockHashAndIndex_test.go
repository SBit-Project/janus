package transformer

import (
	"encoding/json"
	"testing"

	"github.com/SBit-Project/janus/pkg/internal"
	"github.com/SBit-Projectt/janus/pkg/sbit"
)

func initializeProxyETHGetTransactionByBlockHashAndIndex(sbitClient *sbit.Sbit) ETHProxy {
	return &ProxyETHGetTransactionByBlockHashAndIndex{sbitClient}
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockHashAndIndex,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockHash + `"`), []byte(`"0x0"`)},
		internal.GetTransactionByHashResponseData,
	)
}
