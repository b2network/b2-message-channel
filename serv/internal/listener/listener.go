package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type DataMap struct {
	Events    []common.Hash
	Contracts []common.Address
	EventMap  map[common.Hash][]Event
	SenderMap map[string]string
}

type Listener struct {
	DataMap           DataMap
	Blockchain        config.Blockchain
	Db                *gorm.DB
	Cache             *config.Cache
	RPC               *ethclient.Client
	LatestBlockNumber int64
	SyncedBlockNumber int64
	SyncedBlockHash   common.Hash
}

func NewListener(db *gorm.DB, cache *config.Cache, blockchain config.Blockchain) *Listener {
	rpc := config.InitB2Rpc(blockchain.RpcUrl)
	return &Listener{
		Blockchain:        blockchain,
		Db:                db,
		Cache:             cache,
		RPC:               rpc,
		LatestBlockNumber: blockchain.InitBlockNumber,
		SyncedBlockNumber: 0,
		SyncedBlockHash:   common.HexToHash("0x0"),
		DataMap: DataMap{
			Events:    make([]common.Hash, 0),
			Contracts: make([]common.Address, 0),
			EventMap:  make(map[common.Hash][]Event, 0),
			SenderMap: make(map[string]string, 0),
		},
	}
}

func (l *Listener) Run() {
	l.loadAccounts()
	l.AutoRegister()
	go l.syncLatestBlockNumber()
	go l.syncBlock()
	go l.syncEvent()
	go l.syncTask()
	go l.checkBlock()
	go l.migrateBlock()
	go l.migrateEvent()
	go l.consume()
	go l.confirm()
	go l.broadcast()
	go l.build()
}
