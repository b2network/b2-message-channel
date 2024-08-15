package job

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"bsquared.network/b2-message-channel-serv/internal/listener/bitcoin"
	"bsquared.network/b2-message-channel-serv/internal/listener/eth"
	svc "bsquared.network/b2-message-channel-serv/internal/utils/ctx"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"gorm.io/gorm"
)

var (
	ctx *svc.ServiceContext
)

func Run(db *gorm.DB, cache *config.Cache, cfg config.AppConfig) {
	for _, blockchain := range cfg.Blockchain {
		if blockchain.ChainId == 1 {
			bclient, err := rpcclient.New(&rpcclient.ConnConfig{
				Host:         blockchain.RpcUrl,
				User:         blockchain.BtcUser,
				Pass:         blockchain.BtcPass,
				HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
				DisableTLS:   false, // Bitcoin core does not provide TLS by default
			}, nil)
			if err != nil {
				fmt.Println("new: ", err.Error())
				return
			}
			listener := bitcoin.NewListener(db, cache, blockchain, cfg.Particle, bclient)
			listener.Run()
		} else {
			listener := eth.NewListener(db, cache, blockchain)
			listener.Run()
		}
	}
}
