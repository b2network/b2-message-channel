package main

import (
	"bsquared.network/b2-message-channel-listener/internal/config"
	"bsquared.network/b2-message-channel-listener/internal/initiates"
	"bsquared.network/b2-message-channel-listener/internal/serves/builder"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	cfg := config.LoadConfig()

	db := initiates.InitDB(cfg.Database)
	bsquaredRpc := initiates.InitEthereumRpc(cfg.Bsquared)
	bsquaredBuilder := builder.NewBuilder(cfg.Builder.Bsquared, cfg.Log, cfg.Bsquared, db, bsquaredRpc)
	bsquaredBuilder.Start()

	arbitrumRpc := initiates.InitEthereumRpc(cfg.Arbitrum)
	arbitrumBuilder := builder.NewBuilder(cfg.Builder.Arbitrum, cfg.Log, cfg.Arbitrum, db, arbitrumRpc)
	arbitrumBuilder.Start()
	select {}
}
