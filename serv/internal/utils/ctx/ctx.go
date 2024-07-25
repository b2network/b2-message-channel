package ctx

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type Blockchain struct {
	RPC               *ethclient.Client
	LatestBlockNumber int64
	SyncedBlockNumber int64
	SyncedBlockHash   common.Hash
}

type ServiceContext struct {
	Config      config.AppConfig
	Db          *gorm.DB
	Cache       *config.Cache
	Blockchains []Blockchain
}

func NewServiceContext(db *gorm.DB, cache *config.Cache, cfg config.AppConfig) *ServiceContext {
	blockchains := make([]Blockchain, 0)
	for _, blockchain := range cfg.Blockchain {
		b2Rpc := config.InitB2Rpc(blockchain.RpcUrl)
		blockchains = append(blockchains, Blockchain{
			RPC:               b2Rpc,
			LatestBlockNumber: blockchain.InitBlockNumber,
			SyncedBlockNumber: 0,
			SyncedBlockHash:   common.HexToHash("0x0"),
		})
	}
	return &ServiceContext{
		Config:      cfg,
		Db:          db,
		Cache:       cache,
		Blockchains: blockchains,
	}
}
