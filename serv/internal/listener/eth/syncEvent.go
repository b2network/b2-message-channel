package eth

import (
	"bsquared.network/b2-message-channel-serv/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

func (l *Listener) syncEvent() {
	for {
		duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
		var blocks []models.SyncBlock
		err := l.Db.Where("chain_id=? AND (status=? OR status=?)", l.Blockchain.ChainId, models.BlockPending, models.BlockRollback).Order("block_number").Limit(50).Find(&blocks).Error
		if err != nil {
			log.Errorf("[Handler.SyncEvent] Pending and rollback blocks err: %s\n", errors.WithStack(err))
			time.Sleep(duration)
			continue
		}
		if len(blocks) == 0 {
			log.Infof("[Handler.SyncEvent] Pending blocks count is 0\n")
			time.Sleep(duration)
			continue
		}

		var wg sync.WaitGroup
		for _, block := range blocks {
			wg.Add(1)
			go func(_wg *sync.WaitGroup, block models.SyncBlock) {
				defer _wg.Done()
				if block.Status == models.BlockPending {
					// add events & block.status= valid
					err = l.handlePendingBlock(block)
					if err != nil {
						log.Errorf("[Handler.SyncEvent] HandlePendingBlock err: %s\n", errors.WithStack(err))
					}
				} else if block.Status == models.BlockRollback {
					// event.status=rollback & block.status=invalid
					err = l.handleRollbackBlock(block)
					if err != nil {
						log.Errorf("[Handler.SyncEvent] HandleRollbackBlock err: %s\n", errors.WithStack(err))
					}
				}
			}(&wg, block)
		}
		wg.Wait()
	}
}

func (l *Listener) handlePendingBlock(block models.SyncBlock) error {
	log.Infof("[Handler.SyncEvent.PendingBlock]Start: %d, %s \n", block.BlockNumber, block.BlockHash)
	log.Infof("[Handler.SyncEvent.PendingBlock]GetContracts: %v\n", l.GetContracts())
	log.Infof("[Handler.SyncEvent.PendingBlock]GetEvents: %v\n", l.GetEvents())
	events, err := l.LogFilter(block, l.GetContracts(), [][]common.Hash{l.GetEvents()})
	log.Infof("[Handler.SyncEvent.PendingBlock] block %d, events number is %d:", block.BlockNumber, len(events))
	if err != nil {
		log.Errorf("[Handler.SyncEvent.PendingBlock] Log filter err: %s\n", err)
		return errors.WithStack(err)
	}
	eventCount := len(events)
	if eventCount > 0 && events[0].BlockHash != block.BlockHash {
		log.Infof("[Handler.SyncEvent.PendingBlock]Don't match block hash\n")
		return nil
	} else if eventCount > 0 && events[0].BlockHash == block.BlockHash {
		BatchEvents := make([]*models.SyncEvent, 0)
		for _, event := range events {
			var one models.SyncEvent
			log.Infof("[Handler.SyncEvent.PendingBlock]BlockLogIndexed: %d ,TxHash: %s,EventHash: %s\n", event.BlockLogIndexed, event.TxHash, event.EventHash)
			err = l.Db.Select("id").Where("sync_block_id=? AND block_log_indexed=? AND tx_hash=? AND event_hash=? ",
				block.Id, event.BlockLogIndexed, event.TxHash, event.EventHash).First(&one).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Errorf("[Handler.SyncEvent.PendingBlock]Query SyncEvent err: %s\n ", err)
				return errors.WithStack(err)
			} else if err == gorm.ErrRecordNotFound {
				BatchEvents = append(BatchEvents, event)
				log.Infof("[Handler.SyncEvent.PendingBlock]block %d, BatchEvents len is %d:", block.BlockNumber, len(BatchEvents))
			}
		}
		if len(BatchEvents) > 0 {
			err = l.Db.Transaction(func(tx *gorm.DB) error {
				err = tx.CreateInBatches(&BatchEvents, 200).Error
				if err != nil {
					log.Errorf("[Handler.SyncEvent.PendingBlock]CreateInBatches err: %s\n ", err)
					return errors.WithStack(err)
				}
				block.Status = models.BlockValid
				block.EventCount = int64(eventCount)
				err = tx.Save(&block).Error
				if err != nil {
					log.Errorf("[Handler.SyncEvent.PendingBlock]Batch Events Update SyncBlock Status err: %s\n ", err)
					return errors.WithStack(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	block.Status = models.BlockValid
	block.EventCount = int64(eventCount)
	err = l.Db.Save(&block).Error
	if err != nil {
		log.Errorf("[Handler.PendingBlock]Update SyncBlock Status err: %s\n ", err)
		return errors.WithStack(err)
	}
	log.Infof("[Handler.SyncEvent.PendingBlock]End: %d, %s \n", block.BlockNumber, block.BlockHash)
	return nil
}

func (l *Listener) handleRollbackBlock(block models.SyncBlock) error {
	log.Infof("[Handler.RollbackBlock] Start: %d, %s\n", block.BlockNumber, block.BlockHash)
	err := l.Db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// event.status=rollback by syncBlockId
		err := tx.Model(models.SyncEvent{}).Where("sync_block_id=?", block.Id).
			Updates(map[string]interface{}{"status": models.EventRollback, "updated_at": now}).Error
		if err != nil {
			log.Errorf("[Handler.RollbackBlock]Query SyncBlock Status err: %s ,id : %d \n", err, block.Id)
			return errors.WithStack(err)
		}
		block.Status = models.BlockInvalid
		err = tx.Save(&block).Error
		if err != nil {
			log.Errorf("[Handler.RollbackBlock]Save SyncBlock Status err: %s\n ", err)
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		log.Errorf("[Handler.RollbackBlock] err: %s\n ", err)
		return errors.WithStack(err)
	}
	return nil
}
