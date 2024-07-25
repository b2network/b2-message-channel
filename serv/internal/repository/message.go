package repository

import (
	"gorm.io/gorm"
)

type MessageRepository interface {
}

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func MessageRepositoryInit(db *gorm.DB) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{db: db}
}
