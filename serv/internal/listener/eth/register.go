package eth

import (
	"bsquared.network/b2-message-channel-serv/internal/event/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Event interface {
	Name() string
	EventHash() common.Hash
	Data(log types.Log) (string, error)
	ToObj(data string) error
}

func (l *Listener) AutoRegister() {
	events := strings.Split(l.Blockchain.Events, ",")
	for _, event := range events {
		switch common.HexToHash(event) {
		case common.BytesToHash(message.MessageCallHash):
			l.Register(&message.MessageCall{})
		case common.BytesToHash(message.MessageSendHash):
			l.Register(&message.MessageSend{})
		}
	}
	l.AddContract(l.Blockchain.MessageAddress)
}

func (l *Listener) loadAccounts() {
	senders := strings.Split(l.Blockchain.Senders, ",")
	for i, accountKey := range senders {
		accountAddress, err := l.ToAddress(accountKey)
		if err != nil {
			log.Panicf("senders[%d] invalid", i)
			continue
		}
		l.DataMap.SenderMap[accountAddress] = accountKey
	}
	validators := strings.Split(l.Blockchain.Validators, ",")
	for i, accountKey := range validators {
		accountAddress, err := l.ToAddress(accountKey)
		if err != nil {
			log.Panicf("validators[%d] invalid", i)
			continue
		}
		l.DataMap.ValidatorMap[accountAddress] = accountKey
	}
}

func (l *Listener) ToAddress(accountKey string) (string, error) {
	_key, err := crypto.ToECDSA(common.FromHex(accountKey))
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(_key.PublicKey).Hex(), nil
}

func (l *Listener) Register(event Event) {
	if len(l.DataMap.EventMap[event.EventHash()]) == 0 {
		l.DataMap.Events = append(l.DataMap.Events, event.EventHash())
	}
	l.DataMap.EventMap[event.EventHash()] = append(l.DataMap.EventMap[event.EventHash()], event)
}

func (l *Listener) AddContract(contract string) {
	l.DataMap.Contracts = append(l.DataMap.Contracts, common.HexToAddress(contract))
}

func (l *Listener) GetContracts() []common.Address {
	return l.DataMap.Contracts
}

func (l *Listener) GetEvents() []common.Hash {
	return l.DataMap.Events
}

func (l *Listener) GetEventStrs() []string {
	strs := make([]string, 0)
	for _, event := range l.DataMap.Events {
		strs = append(strs, event.Hex())
	}
	return strs
}

func (l *Listener) GetEventsByTask() []common.Hash {
	return l.DataMap.Events
}

func (l *Listener) GetEvent(eventHash common.Hash, contractAddress common.Address) Event {
	EventList := l.DataMap.EventMap[eventHash]
	Event := EventList[0]
	return Event
}
