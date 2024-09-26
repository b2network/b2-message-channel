package main

import (
	"bsquared.network/b2-message-channel-applications/internal/config"
	"bsquared.network/b2-message-channel-applications/internal/initiates"
	"bsquared.network/b2-message-channel-applications/internal/serves/builder"
	"bsquared.network/b2-message-channel-applications/internal/utils/log"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	var fileName string
	flag.StringVar(&fileName, "f", "builder", "-f config filename, default: builder")
	flag.Parse()
	cfg := config.LoadConfig(fileName)
	logger := log.NewLogger(fmt.Sprintf("builder-common"), cfg.Log.Level)
	value, err := json.Marshal(cfg)
	if err != nil {
		logger.Panicf("json marshal err: %s", err)
	}
	logger.Infof("config: %s", value)
	logger.Info("------------------------------------------------------")

	db, err := initiates.InitDB(cfg.Database)
	if err != nil {
		logger.Panicf("init db err: %s", err)
	}

	go func() {
		logger := log.NewLogger(fmt.Sprintf("builder-%s", cfg.Bsquared.Name), cfg.Log.Level)
		rpc, err := initiates.InitEthereumRpc(cfg.Bsquared.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		builder.NewBuilder(cfg.Bsquared.Builders, cfg.Bsquared, db, rpc, logger).Start()
	}()

	go func() {
		logger := log.NewLogger(fmt.Sprintf("builder-%s", cfg.Arbitrum.Name), cfg.Log.Level)
		rpc, err := initiates.InitEthereumRpc(cfg.Arbitrum.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		builder.NewBuilder(cfg.Bsquared.Builders, cfg.Arbitrum, db, rpc, logger).Start()
	}()
	logger.Info("======================================================")
	select {}
}
