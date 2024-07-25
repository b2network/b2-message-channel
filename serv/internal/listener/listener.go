package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"time"
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
	go l.syncLatestBlockNumber(time.Second * 1)
	go l.syncBlock(time.Second * 1)
	go l.syncEvent(time.Second * 1)
	go l.syncTask(time.Second * 1)
	go l.checkBlock(time.Second * 10)
	go l.migrateBlock(time.Second * 10)
	go l.migrateEvent(time.Second * 10)
	go l.consume(time.Second * 1)
	go l.confirm(time.Second * 1)
	go l.broadcast(time.Second * 1)
	go l.build(time.Second * 1)
}
