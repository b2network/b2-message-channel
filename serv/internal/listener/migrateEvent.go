package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (l *Listener) migrateEvent() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		queryBlockNum := l.SyncedBlockNumber - 100000
		if queryBlockNum < 100 {
			time.Sleep(duration)
			continue
		}
		var events []models.SyncEvent
		err := l.Db.Model(models.SyncEvent{}).Where("chain_id=? AND (status=? OR status=?) AND block_number < ?",
			l.Blockchain.ChainId, models.EventValid, models.EventInvalid, queryBlockNum).
			Order("block_number").Order("block_log_indexed").Limit(400).Find(&events).Error
		if err != nil {
			log.Errorf("[Handler.MigrateEvent] Find Events err: %s\n", errors.WithStack(err))
			time.Sleep(duration)
			continue
		}
		if len(events) == 0 {
			time.Sleep(duration)
			continue
		}

		delIds := make([]int64, len(events), len(events))
		eventsHis := make([]*models.SyncEventHistory, len(events), len(events))
		for i, v := range events {
			delIds[i] = v.Id
			eventsHis[i] = &models.SyncEventHistory{
				Base: models.Base{
					Id:        v.Id,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				},
				SyncBlockId:     v.SyncBlockId,
				ChainId:         v.ChainId,
				BlockTime:       v.BlockTime,
				BlockNumber:     v.BlockNumber,
				BlockHash:       v.BlockHash,
				BlockLogIndexed: v.BlockLogIndexed,
				TxIndex:         v.TxIndex,
				TxHash:          v.TxHash,
				EventName:       v.EventName,
				EventHash:       v.EventHash,
				ContractAddress: v.ContractAddress,
				Data:            v.Data,
				Status:          v.Status,
			}
		}

		err = l.Db.Transaction(func(tx *gorm.DB) error {
			err := tx.CreateInBatches(eventsHis, 100).Error
			if err != nil && strings.Index(err.Error(), "Duplicate") == -1 {
				log.Errorf("[Handler.MigrateEvent] CreateInBatches err: %s\n", err)
				return errors.WithStack(err)
			}

			err = tx.Where("id in ?", delIds).Delete(models.SyncEvent{}).Error
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})

		if err != nil {
			log.Errorf("[Handler.MigrateEvent] err: %s\n", err)
			time.Sleep(duration)
			continue
		}

	}
}
