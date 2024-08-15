package models

import "bsquared.network/b2-message-channel-serv/internal/enums"

type Message struct {
	Base
	ChainId             int64               `json:"chain_id"`
	Type                enums.MessageType   `json:"type"`
	FromChainId         int64               `json:"from_chain_id"`
	FromSender          string              `json:"from_sender"`
	FromContractAddress string              `json:"from_contract_address"`
	FromId              string              `json:"from_id"`
	ToChainId           int64               `json:"to_chain_id"`
	ToContractAddress   string              `json:"to_contract_address"`
	ToBytes             string              `json:"to_bytes"`
	Signatures          string              `json:"signatures"`
	Status              enums.MessageStatus `json:"status"`
	Blockchain
}

func (Message) TableName() string {
	return "`messages`"
}
