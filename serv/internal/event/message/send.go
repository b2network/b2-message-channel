package message

import (
	"bsquared.network/b2-message-channel-serv/internal/event"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	MessageSendName = "message#send"
	MessageSendHash = crypto.Keccak256([]byte("Send(uint256,uint256,address,uint256,address,bytes)"))
)

type MessageSend struct {
	FromChainId     int64  `json:"from_chain_id"`
	FromId          int64  `json:"from_id"`
	FromSender      string `json:"from_sender"`
	ToChainId       int64  `json:"to_chain_id"`
	ContractAddress string `json:"contract_address"`
	Bytes           string `json:"bytes"`
}

func (*MessageSend) Name() string {
	return MessageSendName
}

func (*MessageSend) EventHash() common.Hash {
	return common.BytesToHash(MessageSendHash)
}

func (t *MessageSend) ToObj(data string) error {
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return err
	}
	return nil
}

func (*MessageSend) Data(log types.Log) (string, error) {
	transfer := &MessageSend{
		FromChainId:     event.DataToInt64(log, 0),
		FromId:          event.DataToInt64(log, 1),
		FromSender:      event.DataToAddress(log, 2).Hex(),
		ToChainId:       event.DataToInt64(log, 3),
		ContractAddress: event.DataToAddress(log, 4).Hex(),
		Bytes:           event.DataToBytes(log, 5),
	}
	data, err := event.ToJson(transfer)
	if err != nil {
		return "", err
	}
	return data, nil
}
