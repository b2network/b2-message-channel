package main

import (
	"bsquared.network/b2-message-channel-listener/internal/config"
	"bsquared.network/b2-message-channel-listener/internal/initiates"
	"bsquared.network/b2-message-channel-listener/internal/serves/builder"
	"bsquared.network/b2-message-channel-listener/internal/serves/listener/bitcoin"
	"bsquared.network/b2-message-channel-listener/internal/serves/listener/ethereum"
	"bsquared.network/b2-message-channel-listener/internal/serves/validator"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	cfg := config.LoadConfig()

	db := initiates.InitDB(cfg.Database)
	bitcoinRpc := initiates.InitBitcoinRpc(cfg.Bitcoin)
	bsquaredRpc := initiates.InitEthereumRpc(cfg.Bsquared)
	arbitrumRpc := initiates.InitEthereumRpc(cfg.Arbitrum)

	// builder
	bsquaredBuilder := builder.NewBuilder(cfg.Builder.Bsquared, cfg.Log, cfg.Bsquared, db, bsquaredRpc)
	bsquaredBuilder.Start()

	arbitrumBuilder := builder.NewBuilder(cfg.Builder.Arbitrum, cfg.Log, cfg.Arbitrum, db, arbitrumRpc)
	arbitrumBuilder.Start()

	// listener
	//bsquaredRpc := initiates.InitEthereumRpc(cfg.Bsquared)
	bsquaredListener := ethereum.NewListener(cfg.Log, cfg.Bsquared, bsquaredRpc, db)
	bsquaredListener.Start()

	arbitrumListener := ethereum.NewListener(cfg.Log, cfg.Arbitrum, arbitrumRpc, db)
	arbitrumListener.Start()

	bitcoinListener := bitcoin.NewListener(cfg.Log, cfg.Bitcoin, cfg.Particle, bitcoinRpc, db)
	bitcoinListener.Start()

	// validator
	bsquaredValidator := validator.NewValidator(cfg.Validator.Bsquared, cfg.Log, cfg.Bsquared, db, bsquaredRpc)
	bsquaredValidator.Start()

	arbitrumValidator := validator.NewValidator(cfg.Validator.Arbitrum, cfg.Log, cfg.Arbitrum, db, arbitrumRpc)
	arbitrumValidator.Start()
	select {}
}
