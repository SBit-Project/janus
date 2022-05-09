package transformer

import (
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/notifier"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Transformer struct {
	sbitClient   *sbit.Sbit
	debugMode    bool
	logger       log.Logger
	transformers map[string]ETHProxy
}

// New creates a new Transformer
func New(sbitClient *sbit.Sbit, proxies []ETHProxy, opts ...Option) (*Transformer, error) {
	if sbitClient == nil {
		return nil, errors.New("sbitClient cannot be nil")
	}

	t := &Transformer{
		sbitClient: sbitClient,
		logger:     log.NewNopLogger(),
	}

	var err error
	for _, p := range proxies {
		if err = t.Register(p); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Register registers an ETHProxy to a Transformer
func (t *Transformer) Register(p ETHProxy) error {
	if t.transformers == nil {
		t.transformers = make(map[string]ETHProxy)
	}

	m := p.Method()
	if _, ok := t.transformers[m]; ok {
		return errors.Errorf("method already exist: %s ", m)
	}

	t.transformers[m] = p

	return nil
}

// Transform takes a Transformer and transforms the request from ETH request and returns the proxy request
func (t *Transformer) Transform(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	proxy, err := t.getProxy(req.Method)
	if err != nil {
		return nil, err
	}
	resp, err := proxy.Request(req, c)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *Transformer) getProxy(method string) (ETHProxy, eth.JSONRPCError) {
	proxy, ok := t.transformers[method]
	if !ok {
		return nil, eth.NewMethodNotFoundError(method)
	}
	return proxy, nil
}

func (t *Transformer) IsDebugEnabled() bool {
	return t.debugMode
}

// DefaultProxies are the default proxy methods made available
func DefaultProxies(sbitRPCClient *sbit.Sbit, agent *notifier.Agent) []ETHProxy {
	filter := eth.NewFilterSimulator()
	getFilterChanges := &ProxyETHGetFilterChanges{Sbit: sbitRPCClient, filter: filter}
	ethCall := &ProxyETHCall{Sbit: sbitRPCClient}

	ethProxies := []ETHProxy{
		ethCall,
		&ProxyNetListening{Sbit: sbitRPCClient},
		&ProxyETHPersonalUnlockAccount{},
		&ProxyETHChainId{Sbit: sbitRPCClient},
		&ProxyETHBlockNumber{Sbit: sbitRPCClient},
		&ProxyETHHashrate{Sbit: sbitRPCClient},
		&ProxyETHMining{Sbit: sbitRPCClient},
		&ProxyETHNetVersion{Sbit: sbitRPCClient},
		&ProxyETHGetTransactionByHash{Sbit: sbitRPCClient},
		&ProxyETHGetTransactionByBlockNumberAndIndex{Sbit: sbitRPCClient},
		&ProxyETHGetLogs{Sbit: sbitRPCClient},
		&ProxyETHGetTransactionReceipt{Sbit: sbitRPCClient},
		&ProxyETHSendTransaction{Sbit: sbitRPCClient},
		&ProxyETHAccounts{Sbit: sbitRPCClient},
		&ProxyETHGetCode{Sbit: sbitRPCClient},

		&ProxyETHNewFilter{Sbit: sbitRPCClient, filter: filter},
		&ProxyETHNewBlockFilter{Sbit: sbitRPCClient, filter: filter},
		getFilterChanges,
		&ProxyETHGetFilterLogs{ProxyETHGetFilterChanges: getFilterChanges},
		&ProxyETHUninstallFilter{Sbit: sbitRPCClient, filter: filter},

		&ProxyETHEstimateGas{ProxyETHCall: ethCall},
		&ProxyETHGetBlockByNumber{Sbit: sbitRPCClient},
		&ProxyETHGetBlockByHash{Sbit: sbitRPCClient},
		&ProxyETHGetBalance{Sbit: sbitRPCClient},
		&ProxyETHGetStorageAt{Sbit: sbitRPCClient},
		&ETHGetCompilers{},
		&ETHProtocolVersion{},
		&ETHGetUncleByBlockHashAndIndex{},
		&ETHGetUncleCountByBlockHash{},
		&ETHGetUncleCountByBlockNumber{},
		&Web3ClientVersion{},
		&Web3Sha3{},
		&ProxyETHSign{Sbit: sbitRPCClient},
		&ProxyETHGasPrice{Sbit: sbitRPCClient},
		&ProxyETHTxCount{Sbit: sbitRPCClient},
		&ProxyETHSignTransaction{Sbit: sbitRPCClient},
		&ProxyETHSendRawTransaction{Sbit: sbitRPCClient},

		&ETHSubscribe{Sbit: sbitRPCClient, Agent: agent},
		&ETHUnsubscribe{Sbit: sbitRPCClient, Agent: agent},

		&ProxySBITGetUTXOs{Sbit: sbitRPCClient},
		&ProxySBITGenerateToAddress{Sbit: sbitRPCClient},

		&ProxyNetPeerCount{Sbit: sbitRPCClient},
	}

	permittedSbitCalls := []string{
		sbit.MethodGetHexAddress,
		sbit.MethodFromHexAddress,
	}

	for _, sbitMethod := range permittedSbitCalls {
		ethProxies = append(
			ethProxies,
			&ProxySBITGenericStringArguments{
				Sbit:   sbitRPCClient,
				prefix: "dev",
				method: sbitMethod,
			},
		)
	}

	return ethProxies
}

func SetDebug(debug bool) func(*Transformer) error {
	return func(t *Transformer) error {
		t.debugMode = debug
		return nil
	}
}

func SetLogger(l log.Logger) func(*Transformer) error {
	return func(t *Transformer) error {
		t.logger = log.WithPrefix(l, "component", "transformer")
		return nil
	}
}
