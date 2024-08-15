package eth

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

func (l *Listener) syncTask() {
	for {
		duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
		var tasks []models.SyncTask
		err := l.Db.Where("`chain_type`=? AND chain_id=? AND status=?", enums.ChainTypeEthereum, l.Blockchain.ChainId, models.SyncTaskPending).Limit(20).Find(&tasks).Error
		if err != nil {
			time.Sleep(duration)
			continue
		}
		if len(tasks) == 0 {
			log.Infof("[Handler.syncTask] Pending tasks count is 0\n")
			time.Sleep(duration)
			continue
		}
		wg := sync.WaitGroup{}
		for _, task := range tasks {
			wg.Add(1)
			go func(take models.SyncTask, wg *sync.WaitGroup) {
				defer wg.Done()
				err := l.HandleTake(take)
				if err != nil {
					log.Errorf("[Handler.SyncTask] HandleTask ID: %d , err: %s \n", task.Id, err)
				}
			}(task, &wg)

		}
		wg.Wait()
	}
}

func GetContracts(Contracts string) []common.Address {
	list := strings.Split(Contracts, ",")
	AddressList := make([]common.Address, 0)
	for _, one := range list {
		if one != "" {
			AddressList = append(AddressList, common.HexToAddress(one))
		}
	}
	return AddressList
}

func (l *Listener) HandleTake(task models.SyncTask) error {
	start := task.LatestBlock
	if task.StartBlock > start {
		start = task.StartBlock
	}

	if task.EndBlock > 0 && start > task.EndBlock {
		log.Infof("[Handler.SyncTask]  Handle task has done")
		task.Status = models.SyncTaskDone
		l.Db.Save(&task)
		return nil
	}
	end := start
	if task.HandleNum > 0 {
		end = start + task.HandleNum - 1
	}
	if task.EndBlock > 0 && end > task.EndBlock {
		end = task.EndBlock
	}
	if end > l.LatestBlockNumber {
		end = l.LatestBlockNumber
	}

	Contracts := GetContracts(task.Contracts)
	if len(Contracts) == 0 {
		Contracts = l.GetContracts()
		//log.Infof("[Handler.SyncTask]  Contracts invalid")
		////task.UpdateTime = time.Now()
		//task.Status = models.SyncTaskInvalid
		//ctx.Db.Save(&task)
		//return nil
	}
	events, err := l.LogBatchFilter(start, end, Contracts, [][]common.Hash{l.GetEventsByTask()})
	if err != nil {
		log.Infof("[Handler.SyncTask]  Log filter err: %s\n", err)
		return errors.WithStack(err)
	}
	BatchCreateEvents := make([]*models.SyncEvent, 0)
	BatchUpdateEventIds := make([]int64, 0)
	for _, event := range events {
		var one models.SyncEvent
		err = l.Db.Select("id").Where("block_number=? AND block_log_indexed=? AND tx_hash=? AND event_hash=? ",
			event.BlockNumber, event.BlockLogIndexed, event.TxHash, event.EventHash).First(&one).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("[Handler.SyncTask]Query SyncEvent err: %s\n", err)
			return errors.WithStack(err)
		} else if err == gorm.ErrRecordNotFound {
			BatchCreateEvents = append(BatchCreateEvents, event)
		} else {
			BatchUpdateEventIds = append(BatchUpdateEventIds, one.Id)
		}
	}

	err = l.Db.Transaction(func(tx *gorm.DB) error {
		if len(BatchCreateEvents) > 0 {
			err = tx.CreateInBatches(&BatchCreateEvents, 100).Error
			if err != nil {
				log.Errorf("[Handler.SyncEvent]CreateInBatches err: %s\n", err)
				return errors.WithStack(err)
			}
		}
		if len(BatchUpdateEventIds) > 0 {
			err = tx.Model(models.SyncEvent{}).
				Where("id in ?", BatchUpdateEventIds).
				Update("status", models.EventPending).Error
			if err != nil {
				log.Errorf("[Handler.SyncEvent]BatchUpdateEvents err: %s\n", err)
				return errors.WithStack(err)
			}
		}
		task.LatestBlock = end + 1
		err = tx.Save(&task).Error
		if err != nil {
			log.Errorf("[Handler.SyncEvent]Update SyncTask err: %s\n", err)
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		log.Errorf("[Handler.SyncEvent] err: %s\n", err)
		return err
	}
	return nil
}
