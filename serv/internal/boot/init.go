package boot

import (
	"bsquared.network/b2-message-channel-serv/internal/controllers"
	"bsquared.network/b2-message-channel-serv/internal/repository"
	"bsquared.network/b2-message-channel-serv/internal/services"
)

type Initialization struct {
	MessageInit *MessageInitialization
}

type MessageInitialization struct {
	MessageRepo repository.MessageRepository
	MessageSvc  services.MessageService
	MessageCtrl controllers.MessageController
}

func NewMessageInitialization(MessageRepo repository.MessageRepository, MessageSvc services.MessageService, MessageCtrl controllers.MessageController) *MessageInitialization {
	return &MessageInitialization{
		MessageRepo: MessageRepo,
		MessageSvc:  MessageSvc,
		MessageCtrl: MessageCtrl,
	}
}

func NewInitialization(MessageInit *MessageInitialization) *Initialization {
	return &Initialization{
		MessageInit: MessageInit,
	}
}
