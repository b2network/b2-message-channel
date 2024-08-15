package bitcoin

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"gorm.io/gorm"
)

type Listener struct {
	Blockchain        config.Blockchain
	Particle          config.Particle
	Db                *gorm.DB
	Cache             *config.Cache
	Rpc               *rpcclient.Client
	LatestBlockNumber int64
	ChainParams       *chaincfg.Params // bitcoin network params, e.g. mainnet, testnet, etc.
	listenAddress     btcutil.Address  // need listened bitcoin address
}

func NewListener(db *gorm.DB, cache *config.Cache, blockchain config.Blockchain, particle config.Particle, rpc *rpcclient.Client) *Listener {
	listenAddress, _ := btcutil.DecodeAddress(blockchain.ListenBtcAddress, &chaincfg.MainNetParams)
	var ChainParams *chaincfg.Params
	if blockchain.ChainId == 1 {
		ChainParams = &chaincfg.MainNetParams
	} else {
		ChainParams = &chaincfg.TestNet3Params
	}
	return &Listener{
		Blockchain:        blockchain,
		Particle:          particle,
		Db:                db,
		Cache:             cache,
		Rpc:               rpc,
		LatestBlockNumber: 0,
		ChainParams:       ChainParams,
		listenAddress:     listenAddress,
	}
}

func (l *Listener) Run() {
	//go l.syncLatestBlockNumber()
	//go l.syncTask()
	//go l.syncMessage()
}
