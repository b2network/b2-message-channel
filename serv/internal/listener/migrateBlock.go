package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (l *Listener) migrateBlock(duration time.Duration) {
	for {
		queryBlockNum := l.SyncedBlockNumber - 100000
		if queryBlockNum < 100 {
			time.Sleep(duration)
			continue
		}
		var blocks []models.SyncBlock
		err := l.Db.Model(models.SyncBlock{}).Where("chain_id=? AND (status=? OR status=?) AND block_number < ?",
			l.Blockchain.ChainId, models.BlockValid, models.BlockInvalid, queryBlockNum).
			Order("block_number").Limit(400).Find(&blocks).Error
		if err != nil {
			log.Errorf("[Handler.MigrateBlock] Find Blocks err: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(blocks) == 0 {
			time.Sleep(duration)
			continue
		}

		delIds := make([]int64, len(blocks), len(blocks))
		blocksHis := make([]*models.SyncBlockHistory, len(blocks), len(blocks))
		for i, v := range blocks {
			delIds[i] = v.Id
			blocksHis[i] = &models.SyncBlockHistory{
				Base: models.Base{
					Id:        v.Id,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				},
				ChainId:     v.ChainId,
				Miner:       v.Miner,
				BlockTime:   v.BlockTime,
				BlockNumber: v.BlockNumber,
				BlockHash:   v.BlockHash,
				TxCount:     v.TxCount,
				EventCount:  v.EventCount,
				ParentHash:  v.ParentHash,
				Status:      v.Status,
			}
		}

		err = l.Db.Transaction(func(tx *gorm.DB) error {
			err := tx.CreateInBatches(blocksHis, 100).Error
			if err != nil && strings.Index(err.Error(), "Duplicate") == -1 {
				log.Errorf("[Handler.MigrateBlock] CreateInBatches err: %s\n", err)
				return errors.WithStack(err)
			}
			err = tx.Where("id in ?", delIds).Delete(models.SyncBlock{}).Error
			if err != nil {
				log.Errorf("[Handler.MigrateBlock] Delete SyncBlock err: %s\n", err)
				return errors.WithStack(err)
			}
			return nil
		})

		if err != nil {
			log.Errorf("[Handler.MigrateBlock] err: %s\n", err)
			time.Sleep(duration)
			continue
		}

	}
}
