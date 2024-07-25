package services

import (
	"bsquared.network/b2-message-channel-serv/internal/repository"
)

type MessageService interface {
}

type MessageServiceImpl struct {
	MessageRepository repository.MessageRepository
}

func MessageServiceInit(MessageRepository repository.MessageRepository) *MessageServiceImpl {
	return &MessageServiceImpl{
		MessageRepository: MessageRepository,
	}
}
