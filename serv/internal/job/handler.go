package job

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"bsquared.network/b2-message-channel-serv/internal/listener"
	svc "bsquared.network/b2-message-channel-serv/internal/utils/ctx"
	"gorm.io/gorm"
)

var (
	ctx *svc.ServiceContext
)

func Run(db *gorm.DB, cache *config.Cache, cfg config.AppConfig) {
	listeners := make([]*listener.Listener, 0)
	for _, blockchain := range cfg.Blockchain {
		listeners = append(listeners, listener.NewListener(db, cache, blockchain))
	}
	for _, listener := range listeners {
		listener.Run()
	}
}
