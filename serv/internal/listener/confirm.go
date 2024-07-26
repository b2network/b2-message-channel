package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

func (l *Listener) confirm() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		list, err := l.pendingSendMessage(10)
		if err != nil {
			log.Errorf("Get pending call message err: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(list) == 0 {
			log.Infof("Get pending call message length is 0\n")
			time.Sleep(duration)
			continue
		}
		var wg sync.WaitGroup
		for _, message := range list {
			wg.Add(1)
			go func(_wg *sync.WaitGroup, message models.Message) {
				defer _wg.Done()
				err = l.confirmMessage(message)
				if err != nil {
					log.Errorf("Handle err: %v, %v\n", err, message)
				}
			}(&wg, message)
		}
		wg.Wait()
	}
}

func (l *Listener) confirmMessage(message models.Message) error {
	err := l.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(models.Message{}).Where("type=? AND from_chain_id=? AND from_id=? AND status=?",
			enums.MessageTypeCall, message.FromChainId, message.FromId, enums.MessageStatusBroadcast).
			Update("status", enums.MessageStatusValid).Error
		if err != nil {
			return err
		}
		err = tx.Model(models.Message{}).Where("id=? AND status=?", message.Id, message.Status).
			Update("status", enums.MessageStatusValid).Error
		if err != nil {
			return err
		}
		err = tx.Model(models.Signature{}).Where("chain_id=? AND refer_id=?", message.ChainId, message.FromId).
			Updates(map[string]interface{}{
				"status":       enums.SignatureStatusSuccess,
				"event_id":     message.EventId,
				"block_time":   message.BlockTime,
				"block_number": message.BlockNumber,
				"log_index":    message.LogIndex,
				"tx_hash":      message.TxHash,
			}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *Listener) pendingSendMessage(limit int) ([]models.Message, error) {
	var list []models.Message
	err := l.Db.Where("`to_chain_id`=? AND `type`=? AND status=?", l.Blockchain.ChainId, enums.MessageTypeSend, enums.MessageStatusPending).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
