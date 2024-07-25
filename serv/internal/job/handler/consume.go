package handler

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/event/message"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"bsquared.network/b2-message-channel-serv/internal/utils/ctx"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func Consume(ctx *ctx.ServiceContext) {
	for {
		time.Sleep(3 * time.Second)
		events, err := ListPendingEvent(ctx.Db, 100)
		if err != nil {
			log.Errorf("list pending event err: %v\n", err)
			continue
		}
		if len(events) == 0 {
			log.Infof("list pending event length is 0\n")
			continue
		}
		valids := make([]int64, 0)
		invalids := make([]int64, 0)
		messages := make([]models.Message, 0)
		handles := make(map[string]bool)
		var Type enums.MessageType
		var FromChainId int64
		var FromSender string
		var FromId int64
		var ToChainId int64
		var ToContractAddress string
		var ToBytes string

		for _, event := range events {
			key := fmt.Sprintf("%s#%d", event.TxHash, event.BlockLogIndexed)
			if handles[key] {
				invalids = append(invalids, event.Id)
				continue
			}

			if event.EventName == message.MessageCallName {
				var messageCall message.MessageCall
				err := (&messageCall).ToObj(event.Data)
				if err != nil {
					log.Errorf("event to data err: %v, data: %v\n", err, event)
					continue
				}
				FromChainId = messageCall.FromChainId
				FromSender = messageCall.FromSender
				FromId = messageCall.FromId
				ToChainId = messageCall.ToChainId
				ToContractAddress = messageCall.ContractAddress
				ToBytes = messageCall.Bytes
				Type = enums.MessageTypeCall
			} else if event.EventName == message.MessageSendName {
				var messageSend message.MessageSend
				err := (&messageSend).ToObj(event.Data)
				if err != nil {
					log.Errorf("event to data err: %v, data: %v\n", err, event)
					continue
				}
				FromChainId = messageSend.FromChainId
				FromSender = messageSend.FromSender
				FromId = messageSend.FromId
				ToChainId = messageSend.ToChainId
				ToContractAddress = messageSend.ContractAddress
				ToBytes = messageSend.Bytes
				Type = enums.MessageTypeSend
			}

			var message models.Message
			err = ctx.Db.Where("tx_hash=? AND log_index=?", event.TxHash, event.BlockLogIndexed).First(&message).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Errorf("get message err: %v, data: %v\n", err, event)
				continue
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				handles[key] = true
				message = models.Message{
					Type:                Type,
					FromChainId:         FromChainId,
					FromSender:          FromSender,
					FromContractAddress: event.ContractAddress,
					FromId:              FromId,
					ToChainId:           ToChainId,
					ToContractAddress:   ToContractAddress,
					ToBytes:             ToBytes,
					Blockchain: models.Blockchain{
						EventId:     event.Id,
						BlockTime:   event.BlockTime,
						BlockNumber: event.BlockNumber,
						LogIndex:    event.BlockLogIndexed,
						TxHash:      event.TxHash,
					},
				}
				messages = append(messages, message)
				valids = append(valids, event.Id)
			} else {
				invalids = append(invalids, event.Id)
			}

		}
		err = ctx.Db.Transaction(func(tx *gorm.DB) error {
			if len(valids) > 0 {
				err = tx.Model(models.SyncEvent{}).Where("id in (?)", valids).Update("status", models.EventValid).Error
				if err != nil {
					log.Errorf("update valid Event  err: %v, data: %v\n", err, valids)
					return err
				}
			}
			if len(invalids) > 0 {
				err = tx.Model(models.SyncEvent{}).Where("id in (?)", invalids).Update("status", models.EventInvalid).Error
				if err != nil {
					log.Errorf("update invalid Event  err: %v, data: %v\n", err, invalids)
					return err
				}
			}
			if len(messages) > 0 {
				err = tx.CreateInBatches(messages, 100).Error
				if err != nil {
					log.Errorf("create in batches err: %v\n", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Errorf("consume events err: %v\n", err)
		}
	}
}

func ListPendingEvent(db *gorm.DB, limit int) ([]models.SyncEvent, error) {
	var list []models.SyncEvent
	err := db.Model(models.SyncEvent{}).Where("`event_name` in ? AND status=?",
		[]string{message.MessageCallName, message.MessageSendName},
		models.EventPending).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}
