package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func (l *Listener) checkBlock() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		var blocks []models.SyncBlock
		// only check last 1000 block and  event_count is 0 and Valid
		start := l.SyncedBlockNumber - 1000
		if start < 0 {
			start = 0
		}
		end := l.SyncedBlockNumber - 100
		if end < 0 {
			end = 0
		}
		if start < end {
			err := l.Db.Where("block_number BETWEEN ? AND ? AND status = ? AND event_count=? AND check_count<10",
				start, end, models.BlockValid, 0).Order("block_number").Order("check_count").Find(&blocks).Error
			if err != nil {
				log.Errorf("[Handle.CheckBlock] Find Block err: %s\n", errors.WithStack(err))
				time.Sleep(duration)
				continue
			}
			log.Errorf("[Handle.CheckBlock] Blocks Length is: %d", len(blocks))
			if len(blocks) == 0 {
				time.Sleep(duration)
				continue
			}
			for k, _ := range blocks {
				block := blocks[k]
				err = l.HandleCheckBlock(block)
				if err != nil {
					log.Errorf("[Handle.CheckBlock] Check Block err: %s\n", errors.WithStack(err))
				}
			}
		} else {
			time.Sleep(duration)
		}
	}
}

func (l *Listener) HandleCheckBlock(block models.SyncBlock) error {
	log.Infof("[Handle.CheckBlock] Check Block: %d, %s\n", block.BlockNumber, block.BlockHash)
	events, err := l.LogFilter(block, l.GetContracts(), [][]common.Hash{l.GetEvents()})
	if err != nil {
		log.Errorf("[Handle.CheckBlock] Log filter err: %s\n", err)
		return errors.WithStack(err)
	}
	eventCount := len(events)
	if eventCount > 0 && events[0].BlockHash != block.BlockHash {
		log.Errorf("[Handle.CheckBlock]Don't match block hash\n")
		return nil
	} else if eventCount > 0 && events[0].BlockHash == block.BlockHash {
		BatchEvents := make([]*models.SyncEvent, 0)
		for _, event := range events {
			var one models.SyncEvent
			err = l.Db.Where("sync_block_id=? AND block_log_indexed=? AND tx_hash=? AND event_hash=? ",
				block.Id, event.BlockLogIndexed, event.TxHash, event.EventHash).First(&one).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Errorf("[Handle.CheckBlock]Query SyncEvent Is err %s\n", err)
				return errors.WithStack(err)
			} else if err == gorm.ErrRecordNotFound {
				BatchEvents = append(BatchEvents, event)
			}
		}
		if len(BatchEvents) > 0 {
			err = l.Db.Transaction(func(tx *gorm.DB) error {
				err = tx.CreateInBatches(&BatchEvents, 100).Error
				if err != nil {
					log.Errorf("[Handle.CheckBlock]CreateInBatches Is err %s\n", err)
					return errors.WithStack(err)
				}
				block.Status = models.BlockValid
				block.EventCount = int64(eventCount)
				block.CheckCount = block.CheckCount + 1
				err = tx.Save(&block).Error
				if err != nil {
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
	block.CheckCount = block.CheckCount + 1
	//block
	err = l.Db.Save(&block).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
