package message

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"
	"github.com/storyicon/sigverify"
)

const MessageSendTypedData = `{
    "types":{
        "EIP712Domain":[
            {
                "name":"name",
                "type":"string"
            },
            {
                "name":"version",
                "type":"string"
            },
            {
                "name":"chainId",
                "type":"uint256"
            },
            {
                "name":"verifyingContract",
                "type":"address"
            }
        ],
        "Send":[
            {
                "name":"from_chain_id",
                "type":"uint256"
            },
            {
                "name":"from_id",
                "type":"uint256"
            },
            {
                "name":"from_sender",
                "type":"address"
            },
            {
                "name":"to_chain_id",
                "type":"uint256"
            },
            {
                "name":"contract_address",
                "type":"address"
            },
			{
                "name":"data",
                "type":"bytes"
            }
        ]
    },
    "domain":{
        "name":"B2MessageBridge",
        "version":"1",
        "chainId":"%d",
        "verifyingContract":"%s"
    },
    "primaryType":"Send",
    "message":{
        "from_chain_id":"%s",
        "from_id":"%s",
        "from_sender":"%s",
        "to_chain_id":"%s",
        "contract_address":"%s",
        "data":%s
    }
}`

func SignMessageSend(chainId int64, messageContract string, fromChainId int64, fromId int64, fromSender string, toChainId int64, contractAddress string, data string, key *ecdsa.PrivateKey) (string, error) {
	_data := fmt.Sprintf(MessageSendTypedData, chainId, messageContract, fromChainId, fromId, fromSender, toChainId, contractAddress, data)
	var typedData apitypes.TypedData
	if err := json.Unmarshal([]byte(_data), &typedData); err != nil {
		return "", errors.WithStack(err)
	}
	_, originHash, err := sigverify.HashTypedData(typedData)
	fmt.Println("originHash", common.Bytes2Hex(originHash))
	if err != nil {
		return "", errors.WithStack(err)
	}
	sig, err := crypto.Sign(originHash, key)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return common.Bytes2Hex(sig), nil
}
