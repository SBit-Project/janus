package transformer

import (
	"encoding/json"

	"github.com/SBit-Project/janus/pkg/conversion"
	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/labstack/echo"
)

// ProxyETHGetLogs implements ETHProxy
type ProxyETHGetLogs struct {
	*sbit.Sbit
}

func (p *ProxyETHGetLogs) Method() string {
	return "eth_getLogs"
}

func (p *ProxyETHGetLogs) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetLogsRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	// TODO: Graph Node is sending the topic
	// if len(req.Topics) != 0 {
	// 	return nil, errors.New("topics is not supported yet")
	// }

	// Calls ToRequest in order transform ETH-Request to a Sbit-Request
	sbitreq, err := p.ToRequest(&req)
	if err != nil {
		return nil, err
	}

	return p.request(sbitreq)
}

func (p *ProxyETHGetLogs) request(req *sbit.SearchLogsRequest) (*eth.GetLogsResponse, eth.JSONRPCError) {
	receipts, err := conversion.SearchLogsAndFilterExtraTopics(p.Sbit, req)
	if err != nil {
		return nil, err
	}

	logs := make([]eth.Log, 0)
	for _, receipt := range receipts {
		r := sbit.TransactionReceipt(receipt)
		logs = append(logs, conversion.ExtractETHLogsFromTransactionReceipt(r, r.Log)...)
	}

	resp := eth.GetLogsResponse(logs)
	return &resp, nil
}

func (p *ProxyETHGetLogs) ToRequest(ethreq *eth.GetLogsRequest) (*sbit.SearchLogsRequest, eth.JSONRPCError) {
	//transform EthRequest fromBlock to SbitReq fromBlock:
	from, err := getBlockNumberByRawParam(p.Sbit, ethreq.FromBlock, true)
	if err != nil {
		return nil, err
	}

	//transform EthRequest toBlock to SbitReq toBlock:
	to, err := getBlockNumberByRawParam(p.Sbit, ethreq.ToBlock, true)
	if err != nil {
		return nil, err
	}

	//transform EthReq address to SbitReq address:
	var addresses []string
	if ethreq.Address != nil {
		if isBytesOfString(ethreq.Address) {
			var addr string
			if jsonErr := json.Unmarshal(ethreq.Address, &addr); jsonErr != nil {
				return nil, eth.NewInvalidParamsError(jsonErr.Error())
			}
			addresses = append(addresses, addr)
		} else {
			if jsonErr := json.Unmarshal(ethreq.Address, &addresses); jsonErr != nil {
				return nil, eth.NewInvalidParamsError(jsonErr.Error())
			}
		}
		for i := range addresses {
			addresses[i] = utils.RemoveHexPrefix(addresses[i])
		}
	}

	//transform EthReq topics to SbitReq topics:
	topics, topicsErr := eth.TranslateTopics(ethreq.Topics)
	if topicsErr != nil {
		return nil, eth.NewCallbackError(topicsErr.Error())
	}

	return &sbit.SearchLogsRequest{
		Addresses: addresses,
		FromBlock: from,
		ToBlock:   to,
		Topics:    sbit.NewSearchLogsTopics(topics),
	}, nil
}
