package message

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"math/big"
)

func Send(fromChainId int64, fromId *big.Int, fromSender string, contractAddress string, toBytes string, signatures []string) []byte {
	// function send(uint256 from_chain_id, uint256 from_id, address from_sender, address contract_address, bytes calldata data, bytes[] calldata signatures) external
	Method := crypto.Keccak256([]byte("send(uint256,uint256,address,address,bytes,bytes[])"))[:4]
	FromChainId := common.BytesToHash(big.NewInt(fromChainId).Bytes()).Bytes()
	FromId := common.BytesToHash(fromId.Bytes()).Bytes()
	FromSender := common.BytesToHash(common.HexToAddress(fromSender).Bytes()).Bytes()
	ContractAddress := common.BytesToHash(common.HexToAddress(contractAddress).Bytes()).Bytes()

	ToBytes := common.FromHex(toBytes)
	ToBytesDataOffset := common.BytesToHash(big.NewInt(192).Bytes()).Bytes()
	ToBytesDataLength := common.BytesToHash(big.NewInt(int64(len(ToBytes))).Bytes()).Bytes()

	if len(ToBytes)%32 > 0 {
		ToBytes = append(ToBytes, make([]byte, 32-len(ToBytes)%32)...)
	}

	SignaturesDataOffset := common.BytesToHash(big.NewInt(int64(224 + len(ToBytes))).Bytes()).Bytes()
	SignaturesDataLength := common.BytesToHash(big.NewInt(int64(len(signatures))).Bytes()).Bytes()

	var streamOffsets []byte
	var streamData []byte
	streamIndex := int64(32 * len(signatures))
	for _, _signature := range signatures {
		signatureDataOffset := common.BytesToHash(big.NewInt(streamIndex).Bytes()).Bytes()
		streamOffsets = append(streamOffsets, signatureDataOffset...)

		signature := common.FromHex(_signature)
		signatureDataLength := common.BytesToHash(big.NewInt(int64(len(signature))).Bytes()).Bytes()
		streamData = append(streamData, signatureDataLength...)
		if len(signature)%32 > 0 {
			signature = append(signature, make([]byte, 32-len(signature)%32)...)
		}
		streamData = append(streamData, signature...)
		streamIndex = int64(32*len(signatures) + len(streamData))
	}

	var stream []byte
	stream = append(stream, Method...)
	stream = append(stream, FromChainId...)
	stream = append(stream, FromId...)
	stream = append(stream, FromSender...)
	stream = append(stream, ContractAddress...)
	stream = append(stream, ToBytesDataOffset...)
	stream = append(stream, SignaturesDataOffset...)
	stream = append(stream, ToBytesDataLength...)
	stream = append(stream, ToBytes...)
	stream = append(stream, SignaturesDataLength...)
	stream = append(stream, streamOffsets...)
	stream = append(stream, streamData...)
	return stream
}

func EncodeSendData(txId string, fromAddress string, toAddress string, amount decimal.Decimal) []byte {
	TxId := common.HexToHash(txId).Bytes()

	ToAddress := common.BytesToHash(common.HexToAddress(toAddress).Bytes()).Bytes()
	Amount := common.BytesToHash(amount.BigInt().Bytes()).Bytes()

	FromAddress := []byte(fromAddress)
	FromAddressOffset := common.BytesToHash(big.NewInt(128).Bytes()).Bytes()
	FromAddressLength := common.BytesToHash(big.NewInt(int64(len(FromAddress))).Bytes()).Bytes()
	if len(FromAddress)%32 > 0 {
		FromAddress = append(FromAddress, make([]byte, 32-len(FromAddress)%32)...)
	}

	var stream []byte
	stream = append(stream, TxId...)
	stream = append(stream, FromAddressOffset...)
	stream = append(stream, ToAddress...)
	stream = append(stream, Amount...)
	stream = append(stream, FromAddressLength...)
	stream = append(stream, FromAddress...)

	return stream
}

//function encodeData3(bytes32 tx_id, string memory from_address, address to_address, uint256 amount) public pure returns (bytes memory) {
//return abi.encode(tx_id,from_address,to_address,amount);
//}
