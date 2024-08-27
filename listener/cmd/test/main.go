package main

import (
	"bsquared.network/b2-message-channel-listener/internal/utils/ethereum/message"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func main() {

	var chainId int64 = 1
	var messageContract string = "0xb4a88c4086055b8f546911BaaED28b7CEE134EF9"
	var fromChainId int64 = 0
	fromId := big.NewInt(1)
	var fromSender string = common.HexToAddress("0x0").Hex()
	var toChainId int64 = 1
	var contractAddress string = "0xbf99D4875C47D07Db5381c5092D669D3A0eBfB72"
	var data string = "0xbee471d044748a02b1440b19e0e1b4caaccb14d3c9f6f20ec6c02a90445b34c700000000000000000000000000000000000000000000000000000000000000800000000000000000000000000b0c0149d5dbb9eec2c87cff703eadfe428df89100000000000000000000000000000000000000000000000000000000000186a0000000000000000000000000000000000000000000000000000000000000002a746231716b707464737573786b73636176756473326175646663367036377a7533726c307a6d6e377a7500000000000000000000000000000000000000000000"
	var key *ecdsa.PrivateKey
	key, err := crypto.ToECDSA(common.FromHex("0x67f20dc3b0842117c049b292dd88794b3321c95a1b607e735be88c34327420ba"))
	signer := crypto.PubkeyToAddress(key.PublicKey).Hex()
	fmt.Println("signer:", signer)

	signature, err := message.SignMessageSend(chainId, messageContract, fromChainId, fromId, fromSender, toChainId, contractAddress, data, key)
	if err != nil {
		panic(err)
	}
	verify, err := message.VerifyMessageSend(chainId, messageContract, fromChainId, fromId, fromSender, toChainId, contractAddress, data, signer, signature)
	if err != nil {
		panic(err)
	}
	fmt.Println("verify:", verify)
}
