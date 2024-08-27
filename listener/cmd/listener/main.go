package main

import (
	"bsquared.network/b2-message-channel-listener/internal/config"
	"bsquared.network/b2-message-channel-listener/internal/initiates"
	"bsquared.network/b2-message-channel-listener/internal/serves/listener/bitcoin"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	cfg := config.LoadConfig()
	value, _ := json.Marshal(cfg)
	fmt.Println(string(value))
	db := initiates.InitDB(cfg.Database)

	//bsquaredRpc := initiates.InitEthereumRpc(cfg.Bsquared)
	//bsquaredListener := ethereum.NewListener(cfg.Log, cfg.Bsquared, bsquaredRpc, db)
	//bsquaredListener.Start()
	//
	//arbitrumRpc := initiates.InitEthereumRpc(cfg.Arbitrum)
	//arbitrumListener := ethereum.NewListener(cfg.Log, cfg.Arbitrum, arbitrumRpc, db)
	//arbitrumListener.Start()

	bitcoinRpc := initiates.InitBitcoinRpc(cfg.Bitcoin)
	bitcoinListener := bitcoin.NewListener(cfg.Log, cfg.Bitcoin, cfg.Particle, bitcoinRpc, db)
	bitcoinListener.Start()
	select {}
}
