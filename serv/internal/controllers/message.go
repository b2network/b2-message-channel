package controllers

import (
	"bsquared.network/b2-message-channel-serv/internal/config"
	"bsquared.network/b2-message-channel-serv/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageController interface {
	Records(c *gin.Context)
}

type MessageControllerImpl struct {
	BasicController
	svc   services.MessageService
	db    *gorm.DB
	cache *config.Cache
	cfg   config.AppConfig
}

func MessageControllerInit(MessageService services.MessageService, db *gorm.DB, cfg config.AppConfig) *MessageControllerImpl {
	return &MessageControllerImpl{
		svc:   MessageService,
		db:    db,
		cache: config.InitCache(cfg),
		cfg:   cfg,
	}
}

func (c *MessageControllerImpl) Records(ctx *gin.Context) {

}
