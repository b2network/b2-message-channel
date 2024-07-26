package main

import (
	"bsquared.network/b2-message-channel-serv/internal/contract/message"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func main() {
	//fromChainId int64, fromId int64, fromSender string, contractAddress string, toBytes []byte, signatures [][]byte
	var fromChainId int64 = 1123
	var fromId int64 = 1
	var fromSender = "0x9cc4669bb997c40579f89E08980B99218abaE3FE"
	var contractAddress = "0x1c66cBEE6d4660459Fda5aa936e727398175E981"
	var toBytes = "0x1234"
	var signatures = make([]string, 0)
	signatures = append(signatures, "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020003", "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020004")

	data := message.Send(fromChainId, fromId, fromSender, contractAddress, toBytes, signatures)
	fmt.Println(hexutil.Encode(data))
}
