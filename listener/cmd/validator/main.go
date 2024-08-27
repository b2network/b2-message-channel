package main

import (
	"bsquared.network/b2-message-channel-listener/internal/config"
	"bsquared.network/b2-message-channel-listener/internal/initiates"
	"bsquared.network/b2-message-channel-listener/internal/serves/validator"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	cfg := config.LoadConfig()

	db := initiates.InitDB(cfg.Database)

	bsquaredRpc := initiates.InitEthereumRpc(cfg.Bsquared)
	bsquaredValidator := validator.NewValidator(cfg.Validator.Bsquared, cfg.Log, cfg.Bsquared, db, bsquaredRpc)
	bsquaredValidator.Start()

	arbitrumRpc := initiates.InitEthereumRpc(cfg.Arbitrum)
	arbitrumValidator := validator.NewValidator(cfg.Validator.Arbitrum, cfg.Log, cfg.Arbitrum, db, arbitrumRpc)
	arbitrumValidator.Start()

	select {}
}
