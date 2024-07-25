package main

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"bsquared.network/b2-message-channel-serv/internal/job"
	"bsquared.network/b2-message-channel-serv/internal/validators"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	cfg := config.LoadConfig()
	config.InitLog(cfg.Server.LogLevel)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterValidators(v)
	}
	db := config.ConnectToDB(cfg)
	cache := config.InitCache(cfg)
	job.Run(db, cache, cfg)
	select {}
}
