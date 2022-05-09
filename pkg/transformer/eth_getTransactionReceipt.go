package transformer

import (
	"github.com/SBit-Project/janus/pkg/conversion"
	"github.com/SBit-Projectt/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

var STATUS_SUCCESS = "0x1"
var STATUS_FAILURE = "0x0"

// ProxyETHGetTransactionReceipt implements ETHProxy
type ProxyETHGetTransactionReceipt struct {
	*sbit.Sbit
}

func (p *ProxyETHGetTransactionReceipt) Method() string {
	return "eth_getTransactionReceipt"
}

func (p *ProxyETHGetTransactionReceipt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetTransactionReceiptRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}
	if req == "" {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("empty transaction hash")
	}
	var (
		txHash  = utils.RemoveHexPrefix(string(req))
		sbitReq = sbit.GetTransactionReceiptRequest(txHash)
	)
	return p.request(&sbitReq)
}

func (p *ProxyETHGetTransactionReceipt) request(req *sbit.GetTransactionReceiptRequest) (*eth.GetTransactionReceiptResponse, eth.JSONRPCError) {
	sbitReceipt, err := p.Sbit.GetTransactionReceipt(string(*req))
	if err != nil {
		ethTx, _, getRewardTransactionErr := getRewardTransactionByHash(p.Sbit, string(*req))
		if getRewardTransactionErr != nil {
			errCause := errors.Cause(err)
			if errCause == sbit.EmptyResponseErr {
				return nil, nil
			}
			p.Sbit.GetDebugLogger().Log("msg", "Transaction does not exist", "txid", string(*req))
			return nil, eth.NewCallbackError(err.Error())
		}
		if ethTx == nil {
			// unconfirmed tx, return nil
			// https://github.com/openethereum/parity-ethereum/issues/3482
			return nil, nil
		}
		return &eth.GetTransactionReceiptResponse{
			TransactionHash:  ethTx.Hash,
			TransactionIndex: ethTx.TransactionIndex,
			BlockHash:        ethTx.BlockHash,
			BlockNumber:      ethTx.BlockNumber,
			// TODO: This is higher than GasUsed in geth but does it matter?
			CumulativeGasUsed: NonContractVMGasLimit,
			EffectiveGasPrice: "0x0",
			GasUsed:           NonContractVMGasLimit,
			From:              ethTx.From,
			To:                ethTx.To,
			Logs:              []eth.Log{},
			LogsBloom:         eth.EmptyLogsBloom,
			Status:            STATUS_SUCCESS,
		}, nil
	}

	ethReceipt := &eth.GetTransactionReceiptResponse{
		TransactionHash:   utils.AddHexPrefix(sbitReceipt.TransactionHash),
		TransactionIndex:  hexutil.EncodeUint64(sbitReceipt.TransactionIndex),
		BlockHash:         utils.AddHexPrefix(sbitReceipt.BlockHash),
		BlockNumber:       hexutil.EncodeUint64(sbitReceipt.BlockNumber),
		ContractAddress:   utils.AddHexPrefixIfNotEmpty(sbitReceipt.ContractAddress),
		CumulativeGasUsed: hexutil.EncodeUint64(sbitReceipt.CumulativeGasUsed),
		EffectiveGasPrice: "0x0",
		GasUsed:           hexutil.EncodeUint64(sbitReceipt.GasUsed),
		From:              utils.AddHexPrefixIfNotEmpty(sbitReceipt.From),
		To:                utils.AddHexPrefixIfNotEmpty(sbitReceipt.To),

		// TODO: researching
		// ! Temporary accept this value to be always zero, as it is at eth logs
		LogsBloom: eth.EmptyLogsBloom,
	}

	status := STATUS_FAILURE
	if sbitReceipt.Excepted == "None" {
		status = STATUS_SUCCESS
	} else {
		p.Sbit.GetDebugLogger().Log("transaction", ethReceipt.TransactionHash, "msg", "transaction excepted", "message", sbitReceipt.Excepted)
	}
	ethReceipt.Status = status

	r := sbit.TransactionReceipt(*sbitReceipt)
	ethReceipt.Logs = conversion.ExtractETHLogsFromTransactionReceipt(&r, r.Log)

	sbitTx, err := p.Sbit.GetRawTransaction(sbitReceipt.TransactionHash, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", "couldn't get transaction", "err", err)
		return nil, eth.NewCallbackError("couldn't get transaction")
	}
	decodedRawSbitTx, err := p.Sbit.DecodeRawTransaction(sbitTx.Hex)
	if err != nil {
		p.GetDebugLogger().Log("msg", "couldn't decode raw transaction", "err", err)
		return nil, eth.NewCallbackError("couldn't decode raw transaction")
	}
	if decodedRawSbitTx.IsContractCreation() {
		ethReceipt.To = ""
	} else {
		ethReceipt.ContractAddress = ""
	}

	// TODO: researching
	// - The following code reason is unknown (see original comment)
	// - Code temporary commented, until an error occures
	// ! Do not remove
	// // contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
	// if status != "0x1" {
	// 	// if failure, should return null for contractAddress, instead of the zero address.
	// 	ethTxReceipt.ContractAddress = ""
	// }

	return ethReceipt, nil
}
