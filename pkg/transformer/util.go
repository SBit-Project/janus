package transformer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/SBit-Projectt/janus/pkg/sbit"
	"github.com/btcsuite/btcutil/base58"

	"github.com/SBit-Projectt/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var ZeroSatoshi = decimal.NewFromInt(0)
var OneSatoshi = decimal.NewFromFloat(0.00000001)
var MinimumGas = decimal.NewFromFloat(0.0000004)

type EthGas interface {
	GasHex() string
	GasPriceHex() string
}

func EthGasToSbit(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = g.(*eth.SendTransactionRequest).Gas.Int

	gasPriceDecimal, err := EthValueToSbitAmount(g.GasPriceHex(), MinimumGas)
	if err != nil {
		return nil, "0.0", err
	}
	if gasPriceDecimal.LessThan(MinimumGas) {
		gasPriceDecimal = MinimumGas
	}
	gasPrice = fmt.Sprintf("%v", gasPriceDecimal)

	return
}

func SbitGasToEth(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = g.(*eth.SendTransactionRequest).Gas.Int

	gasPriceDecimal, err := EthValueToSbitAmount(g.GasPriceHex(), MinimumGas)
	if err != nil {
		return nil, "0.0", err
	}
	if gasPriceDecimal.LessThan(MinimumGas) {
		gasPriceDecimal = MinimumGas
	}
	gasPrice = fmt.Sprintf("%v", gasPriceDecimal)

	return
}

func EthValueToSbitAmount(val string, defaultValue decimal.Decimal) (decimal.Decimal, error) {
	if val == "" {
		return defaultValue, nil
	}

	ethVal, err := utils.DecodeBig(val)
	if err != nil {
		return ZeroSatoshi, err
	}

	ethValDecimal, err := decimal.NewFromString(ethVal.String())
	if err != nil {
		return ZeroSatoshi, errors.New("decimal.NewFromString was not a success")
	}

	return EthDecimalValueToSbitAmount(ethValDecimal), nil
}

func EthDecimalValueToSbitAmount(ethValDecimal decimal.Decimal) decimal.Decimal {
	// Convert Wei to Sbit
	// 10000000000
	// one satoshi is 0.00000001
	// we need to drop precision for values smaller than that
	// 1e-8?
	maximumPrecision := ethValDecimal.Mul(decimal.NewFromFloat(float64(1e-9))).Floor()
	// was 1e-10
	amount := maximumPrecision.Mul(decimal.NewFromFloat(float64(1e-9)))

	return amount
}

func SbitValueToETHAmount(val string, defaultValue decimal.Decimal) (decimal.Decimal, error) {
	if val == "" {
		return defaultValue, nil
	}

	sbitVal, err := utils.DecodeBig(val)
	if err != nil {
		return ZeroSatoshi, err
	}

	sbitValDecimal, err := decimal.NewFromString(sbitVal.String())
	if err != nil {
		return ZeroSatoshi, errors.New("decimal.NewFromString was not a success")
	}

	return SbitDecimalValueToETHAmount(sbitValDecimal), nil
}

func SbitDecimalValueToETHAmount(sbitValDecimal decimal.Decimal) decimal.Decimal {
	// Computes inverse of EthDecimalValueToSbitAmount
	amount := sbitValDecimal.Div(decimal.NewFromFloat(float64(1e-18)))

	return amount
}

func formatSbitAmount(amount decimal.Decimal) (string, error) {
	decimalAmount := amount.Mul(decimal.NewFromFloat(float64(1e18)))

	//convert decimal to Integer
	result := decimalAmount.BigInt()

	if !decimalAmount.Equals(decimal.NewFromBigInt(result, 0)) {
		return "0x0", errors.New("decimal.BigInt() was not a success")
	}

	return hexutil.EncodeBig(result), nil
}

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "Invalid RPC input")
	}
	return nil
}

// Function for getting the sender address of a non-contract transaction by ID.
// Does not handle OP_SENDER addresses, because it is only present in contract TXs
//
// TODO: Investigate if limitations on Sbit RPC command GetRawTransaction can cause issues here
// Brief explanation: A default config Sbit node can only serve this command for transactions in the mempool, so it will likely break for SOME setup at SOME point.
// However the same info can be found with getblock verbosity = 2, so maybe use that instead?
func getNonContractTxSenderAddress(p *sbit.Sbit, tx *sbit.DecodedRawTransactionResponse) (string, error) {
	// Fetch raw Tx struct, which contains address data for Vins
	rawTx, err := p.GetRawTransaction(tx.ID, false)

	if err != nil {
		return "", errors.New("Couldn't get raw Transaction data from Transaction ID: " + err.Error())
	}

	// If Tx has no vins it's either a reward transaction or invalid/corrupt (Right?). This is outside the intended scope of this function, so throw an error
	if len(rawTx.Vins) == 0 {
		return "", errors.New("Transaction has 0 Vins and thus no valid sender address")
	}

	// Take the address of the first Vin as sender address, as per design decision
	// TODO: Make this not loop, it's not necessary and can in theory produce unintended behavior without causing an error
	// TODO (research): Is the raw TX Vin list always in the "correct" order? It has to be for this function to produce correct behavior
	for _, in := range rawTx.Vins {
		if len(in.Address) == 0 {
			continue
		}
		hexAddress, err := utils.ConvertSbitAddress(in.Address)
		if err != nil {
			return "", err
		}
		return utils.AddHexPrefix(hexAddress), nil
	}

	// If we get here, we have no Vins with a valid address, so search for sender address in previous Tx's vouts
	hexAddr, err := searchSenderAddressInPreviousTransactions(p, rawTx)
	if err != nil {
		return "", errors.New("Couldn't find sender address in previous transactions: " + err.Error())
	}

	return utils.AddHexPrefix(hexAddr), nil
}

// Searchs recursively for the sender address in previous transactions
func searchSenderAddressInPreviousTransactions(p *sbit.Sbit, rawTx *sbit.GetRawTransactionResponse) (string, error) {
	// search within current rawTx for vin containing opcode OP_SPEND
	var vout int64 = -1
	var txid string = ""
	for _, vin := range rawTx.Vins {
		if vin.ScriptSig.Asm == "OP_SPEND" {
			vout = vin.VoutN
			txid = vin.ID
			break
		}
	}
	if vout == -1 {
		return "", errors.New("Couldn't find OP_SPEND in transaction Vins")
	}
	// fetch previous transaction using txid found in vin above
	prevRawTx, err := p.GetRawTransaction(txid, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", "Failed to GetRawTransaction", "tx", txid, "err", err)
		return "", errors.New("Couldn't get raw transaction: " + err.Error())
	}
	// check opcodes contained in vout found in previous transaction
	prevVout := prevRawTx.Vouts[vout]
	script := strings.Split(prevVout.Details.Asm, " ")
	finalOp := script[len(script)-1]
	switch finalOp {
	// If the vout is an OP_SPEND recurse and keep fetching until we find an OP_CREATE or OP_CALL
	case "OP_SPEND":
		return searchSenderAddressInPreviousTransactions(p, prevRawTx)
	// If we find an OP_CREATE, compute the contract address and set that as the "from"
	case "OP_CREATE":
		createInfo, err := sbit.ParseCreateSenderASM(script)
		if err != nil {
			// Check for OP_CREATE without OP_SENDER
			createInfo, err = sbit.ParseCreateASM(script)
			if err != nil {
				return "", errors.WithMessage(err, "couldn't parse create sender ASM")
			}
		}
		return createInfo.From, nil
	// If it's an OP_CALL, extract the contract address and use that as the "from" address
	case "OP_CALL":
		callInfo, err := sbit.ParseCallSenderASM(script)
		if err != nil {
			// Check for OP_CALL without OP_SENDER
			callInfo, err = sbit.ParseCallASM(script)
			if err != nil {
				return "", errors.WithMessage(err, "couldn't parse call sender ASM")
			}
		}
		return callInfo.From, nil
	}
	return "", errors.New("couldn't find sender address")
}

// NOTE:
// 	- is not for reward transactions
// 	- returning address already has 0x prefix
//
// 	TODO: researching
// 	- Vout[0].Addresses[i] != "" - temporary solution
func findNonContractTxReceiverAddress(vouts []*sbit.DecodedRawTransactionOutV) (string, error) {
	for _, vout := range vouts {
		for _, address := range vout.ScriptPubKey.Addresses {
			if address != "" {
				hex, err := utils.ConvertSbitAddress(address)
				if err != nil {
					return "", err
				}
				return utils.AddHexPrefix(hex), nil
			}
		}
	}
	return "", errors.New("not found")
}

func getBlockNumberByHash(p *sbit.Sbit, hash string) (uint64, error) {
	block, err := p.GetBlock(hash)
	if err != nil {
		return 0, errors.WithMessage(err, "couldn't get block")
	}
	p.GetDebugLogger().Log("function", "getBlockNumberByHash", "hash", hash, "block", block.Height)
	return uint64(block.Height), nil
}

func getTransactionIndexInBlock(p *sbit.Sbit, txHash string, blockHash string) (int64, error) {
	block, err := p.GetBlock(blockHash)
	if err != nil {
		return -1, errors.WithMessage(err, "couldn't get block")
	}
	for i, blockTx := range block.Txs {
		if txHash == blockTx {
			p.GetDebugLogger().Log("function", "getTransactionIndexInBlock", "msg", "Found transaction index in block", "txHash", txHash, "blockHash", blockHash, "index", i)
			return int64(i), nil
		}
	}
	p.GetDebugLogger().Log("function", "getTransactionIndexInBlock", "msg", "Could not find transaction index for hash in block", "txHash", txHash, "blockHash", blockHash)
	return -1, errors.New("not found")
}

func formatSbitNonce(nonce int) string {
	var (
		hexedNonce     = strconv.FormatInt(int64(nonce), 16)
		missedCharsNum = 16 - len(hexedNonce)
	)
	for i := 0; i < missedCharsNum; i++ {
		hexedNonce = "0" + hexedNonce
	}
	return "0x" + hexedNonce
}

// Returns Sbit block number. Result depends on a passed raw param. Raw param's slice of bytes should
// has one of the following values:
// 	- hex string representation of a number of a specific block
//  - integer - returns the value
// 	- string "latest" - for the latest mined block
// 	- string "earliest" for the genesis block
// 	- string "pending" - for the pending state/transactions
// Uses defaultVal to differntiate from a eth_getBlockByNumber req and eth_getLogs/eth_newFilter
func getBlockNumberByRawParam(p *sbit.Sbit, rawParam json.RawMessage, defaultVal bool) (*big.Int, eth.JSONRPCError) {
	var param string
	if isBytesOfString(rawParam) {
		param = string(rawParam[1 : len(rawParam)-1]) // trim \" runes
	} else {
		integer, err := strconv.ParseInt(string(rawParam), 10, 64)
		if err == nil {
			return big.NewInt(integer), nil
		}
		return nil, eth.NewInvalidParamsError("invalid parameter format - string or integer is expected")
	}

	return getBlockNumberByParam(p, param, defaultVal)
}

func getBlockNumberByParam(p *sbit.Sbit, param string, defaultVal bool) (*big.Int, eth.JSONRPCError) {
	if len(param) < 1 {
		if defaultVal {
			res, err := p.GetBlockChainInfo()
			if err != nil {
				return nil, eth.NewCallbackError(err.Error())
			}
			p.GetDebugLogger().Log("function", "getBlockNumberByParam", "msg", "returning default value ("+strconv.Itoa(int(res.Blocks))+")")
			return big.NewInt(res.Blocks), nil
		} else {
			return nil, eth.NewInvalidParamsError("empty parameter value")
		}

	}

	switch param {
	case "latest":
		res, err := p.GetBlockChainInfo()
		if err != nil {
			return nil, eth.NewCallbackError(err.Error())
		}
		p.GetDebugLogger().Log("latest", res.Blocks, "msg", "Got latest block")
		return big.NewInt(res.Blocks), nil

	case "earliest":
		// TODO: discuss
		// ! Genesis block cannot be retreived
		return big.NewInt(0), nil

	case "pending":
		// TODO: discuss
		// 	! Researching
		return nil, eth.NewInvalidRequestError("TODO: tag is in implementation")

	default: // hex number
		if !strings.HasPrefix(param, "0x") {
			return nil, eth.NewInvalidParamsError("quantity values must start with 0x")
		}
		n, err := utils.DecodeBig(param)
		if err != nil {
			p.GetDebugLogger().Log("function", "getBlockNumberByParam", "msg", "Failed to decode hex parameter", "value", param)
			return nil, eth.NewInvalidParamsError("couldn't decode hex number to big int")
		}
		return n, nil
	}
}

func isBytesOfString(v json.RawMessage) bool {
	dQuote := []byte{'"'}
	if !bytes.HasPrefix(v, dQuote) && !bytes.HasSuffix(v, dQuote) {
		return false
	}
	if bytes.Count(v, dQuote) != 2 {
		return false
	}
	// TODO: decide
	// ? Should we iterate over v to check if v[1:len(v)-2] is in a range of a-A, z-Z, 0-9
	return true
}

// Converts Ethereum address to a Sbit address, where `address` represents
// Ethereum address without `0x` prefix and `chain` represents target Sbit
// chain
func convertETHAddress(address string, chain string) (sbitAddress string, _ error) {
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't decode hexed address - %q", address)
	}

	var prefix []byte
	switch chain {
	case sbit.ChainMain:
		chainPrefix, err := sbit.PrefixMainChainAddress.AsBytes()
		if err != nil {
			return "", errors.WithMessagef(err, "couldn't convert %q Sbit chain prefix to slice of bytes", chain)
		}
		prefix = chainPrefix

	case sbit.ChainTest, sbit.ChainRegTest:
		chainPrefix, err := sbit.PrefixTestChainAddress.AsBytes()
		if err != nil {
			return "", errors.WithMessagef(err, "couldn't convert %q Sbit chain prefix to slice of bytes", chain)
		}
		prefix = chainPrefix

	default:
		return "", errors.Errorf("unsupported %q Sbit chain", chain)
	}

	var (
		prefixedAddrBytes = append(prefix, addrBytes...)
		checksum          = sbit.CalcAddressChecksum(prefixedAddrBytes)
		sbitAddressBytes  = append(prefixedAddrBytes, checksum...)
	)
	return base58.Encode(sbitAddressBytes), nil
}

func processFilter(p *ProxyETHGetFilterChanges, rawreq *eth.JSONRPCRequest) (*eth.Filter, eth.JSONRPCError) {
	var req eth.GetFilterChangesRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	filterID, err := hexutil.DecodeUint64(string(req))
	if err != nil {
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	_filter, ok := p.filter.Filter(filterID)
	if !ok {
		return nil, eth.NewCallbackError("Invalid filter id")
	}
	filter := _filter.(*eth.Filter)

	return filter, nil
}

// Converts a satoshis to sbit balance
func convertFromSatoshisToSbit(inSatoshis decimal.Decimal) decimal.Decimal {
	return inSatoshis.Div(decimal.NewFromFloat(float64(1e8)))
}

// Converts a sbit balance to satoshis
func convertFromSbitToSatoshis(inSbit decimal.Decimal) decimal.Decimal {
	return inSbit.Mul(decimal.NewFromFloat(float64(1e8)))
}

func convertFromSatoshiToWei(inSatoshis *big.Int) *big.Int {
	return inSatoshis.Mul(inSatoshis, big.NewInt(1e10))
}
