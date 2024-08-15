package types

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
)

// BITCOINTxIndexer defines the interface of custom bitcoin tx indexer.
type BITCOINTxIndexer interface {
	// ParseBlock parse bitcoin block tx
	ParseBlock(int64, int64) ([]*BitcoinTxParseResult, *wire.BlockHeader, error)
	// LatestBlock get latest block height in the longest block chain.
	LatestBlock() (int64, error)
	// CheckConfirmations get tx detail info
	CheckConfirmations(txHash string) error
	// ParseTx parse bitcoin tx
	ParseTx(txHash string) (*BitcoinTxParseResult, error)
}

type BitcoinTxParseResult struct {
	// from is l2 user address, by parse bitcoin get the address
	From []BitcoinFrom
	// to is listening address
	To string
	// value is from transfer amount
	Value int64
	// tx_id is the btc transaction id
	TxID string
	// tx_type is the type of the transaction, eg. "brc20_transfer","transfer"
	TxType string
	// index is the index of the transaction in the block
	Index int64
	// tos tx all to info
	Tos []BitcoinTo
}

const (
	BitcoinFromTypeBtc = 0
	BitcoinFromTypeEvm = 1
)

type BitcoinFrom struct {
	Address    string
	Type       int
	EvmAddress string
}

const (
	BitcoinToTypeNormal   = 0
	BitcoinToTypeNullData = 1
)

type BitcoinTo struct {
	Address  string
	Value    int64
	Type     int
	NullData string
}

func ParseEvmAddressFromNullData(parseResult *BitcoinTxParseResult) (*BitcoinTxParseResult, bool, string, error) {
	existsEvmAddressData := false // The evm address is processed only if it exists. Otherwise, aa is used
	parsedEvmAddress := ""        // evm address
	for _, v := range parseResult.Tos {
		// only handle first null data
		if existsEvmAddressData {
			continue
		}
		if v.Type == BitcoinToTypeNullData {
			decodeNullData, err := hex.DecodeString(v.NullData)
			if err != nil {
				continue
			}
			evmAddress := bytes.TrimSpace(decodeNullData[1:])
			if common.IsHexAddress(string(evmAddress)) {
				existsEvmAddressData = true
				parsedEvmAddress = string(evmAddress)
				for k := range parseResult.From {
					parseResult.From[k].Type = BitcoinFromTypeEvm
					parseResult.From[k].EvmAddress = parsedEvmAddress
				}
			}
		}
	}
	return parseResult, existsEvmAddressData, parsedEvmAddress, nil
}
