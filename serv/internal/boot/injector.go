//go:build wireinject
// +build wireinject

// go:build wireinject
package boot

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"bsquared.network/b2-message-channel-serv/internal/controllers"
	"bsquared.network/b2-message-channel-serv/internal/repository"
	"bsquared.network/b2-message-channel-serv/internal/services"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var appConfig = wire.NewSet(config.LoadConfig)

var MessageInitSet = wire.NewSet(NewMessageInitialization)

var MessageServiceSet = wire.NewSet(services.MessageServiceInit, wire.Bind(new(services.MessageService), new(*services.MessageServiceImpl)))

var MessageRepoSet = wire.NewSet(repository.MessageRepositoryInit, wire.Bind(new(repository.MessageRepository), new(*repository.MessageRepositoryImpl)))

var MessageCtrlSet = wire.NewSet(controllers.MessageControllerInit, wire.Bind(new(controllers.MessageController), new(*controllers.MessageControllerImpl)))

func Init(db *gorm.DB, cache *config.Cache, cfg config.AppConfig) *Initialization {
	wire.Build(NewInitialization,
		MessageInitSet,
		MessageCtrlSet, MessageServiceSet, MessageRepoSet,
	)
	return nil
}
