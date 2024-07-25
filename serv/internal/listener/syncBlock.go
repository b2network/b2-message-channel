package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/models"
	rpc2 "bsquared.network/b2-message-channel-serv/internal/utils/rpc"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func (l *Listener) syncBlock(duration time.Duration) {
	time.Sleep(duration)
	var syncedBlock models.SyncBlock
	err := l.Db.Where("chain_id =? AND (status = ? or status = ?) ", l.Blockchain.ChainId, models.BlockValid, models.BlockPending).Order("block_number desc").First(&syncedBlock).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	} else if err == gorm.ErrRecordNotFound {
		l.SyncedBlockNumber = l.Blockchain.InitBlockNumber
		l.SyncedBlockHash = common.HexToHash(l.Blockchain.InitBlockHash)
	} else {
		l.SyncedBlockNumber = syncedBlock.BlockNumber
		l.SyncedBlockHash = common.HexToHash(syncedBlock.BlockHash)
	}
	log.Infof("[Handler.SyncBlock] blockNumber: %d, blockHash:%s \n", l.SyncedBlockNumber, l.SyncedBlockHash)

	for {
		syncingBlockNumber := l.SyncedBlockNumber + 1
		log.Infof("[Handler.SyncBlock] Try to sync block number: %d\n", syncingBlockNumber)

		if syncingBlockNumber > l.LatestBlockNumber {
			time.Sleep(3 * time.Second)
			continue
		}

		//block, err := ctx.RPC.BlockByNumber(context.Background(), big.NewInt(syncingBlockNumber))
		blockJson, err := rpc2.HttpPostJson("", l.Blockchain.RpcUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", syncingBlockNumber)+"\", true],\"id\":1}")
		if err != nil {
			log.Errorf("[Handler.SyncBlock] Syncing block by number error: %s\n", errors.WithStack(err))
			time.Sleep(3 * time.Second)
			continue
		}
		block := rpc2.ParseJsonBlock(string(blockJson))
		log.Infof("[Handler.SyncBlock] Syncing block number: %d, hash: %v, parent hash: %v \n", block.Number(), block.Hash(), block.ParentHash())
		// 回滚判断
		fmt.Println("block.ParentHash", block.ParentHash())
		fmt.Println("SyncedBlockHash", l.SyncedBlockHash.String())

		if common.HexToHash(block.ParentHash()) != l.SyncedBlockHash {
			log.Errorf("[Handler.SyncBlock] ParentHash of the block being synchronized is inconsistent: %s \n", l.SyncedBlockHash)
			l.rollbackBlock()
			continue
		}

		/* Create SyncBlock start */
		err = l.Db.Create(&models.SyncBlock{
			Miner:       block.Result.Miner,
			ChainId:     l.Blockchain.ChainId,
			BlockTime:   int64(block.Timestamp()),
			BlockNumber: block.Number(),
			BlockHash:   block.Hash(),
			TxCount:     int64(len(block.Result.Transactions)),
			EventCount:  0,
			ParentHash:  block.ParentHash(),
			Status:      models.BlockPending,
		}).Error
		if err != nil {
			log.Errorf("[Handler.SyncBlock] Db Create SyncBlock error: %s\n", errors.WithStack(err))
			time.Sleep(1 * time.Second)
			continue
		}
		/* Create SyncBlock end */
		l.SyncedBlockNumber = block.Number()
		l.SyncedBlockHash = common.HexToHash(block.Hash())
	}
}

func (l *Listener) rollbackBlock() {
	for {
		rollbackBlockNumber := l.SyncedBlockNumber

		log.Infof("[Handler.SyncBlock.RollRackBlock] Try to rollback block number: %d\n", rollbackBlockNumber)

		//rollbackBlock, err := ctx.RPC.BlockByNumber(context.Background(), big.NewInt(rollbackBlockNumber))
		blockJson, err := rpc2.HttpPostJson("", l.Blockchain.RpcUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", rollbackBlockNumber)+"\", true],\"id\":1}")
		if err != nil {
			log.Errorf("[Handler.SyncBlock.RollRackBlock]Rollback block by number error: %s\n", errors.WithStack(err))
			continue
		}
		rollbackBlock := rpc2.ParseJsonBlock(string(blockJson))
		log.Errorf("[Handler.SyncBlock.RollRackBlock] rollbackBlock: %s, syncedBlockHash: %s \n", rollbackBlock.Hash(), l.SyncedBlockHash)

		if common.HexToHash(rollbackBlock.Hash()) == l.SyncedBlockHash {
			err = l.Db.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(models.SyncBlock{}).Where("`chain_id`=？ AND (status = ? or status = ?) AND block_number>?", models.BlockValid, models.BlockPending, l.SyncedBlockNumber).Update("status", models.BlockRollback).Error
				if err != nil {
					log.Errorf("[Handler.SyncBlock.RollRackBlock] Rollback Block err: %s\n", errors.WithStack(err))
					return err
				}
				return nil
			})
			if err != nil {
				log.Errorf("[Handler.SyncBlock.RollRackBlock] Rollback db transaction err: %s\n", errors.WithStack(err))
				continue
			}
			log.Infof("[Handler.SyncBlock.RollRackBlock] Rollback blocks is Stop\n")
			return
		}
		var previousBlock models.SyncBlock
		rest := l.Db.Where("`chain_id`=? AND `block_number`=? AND (status=? or status=?) ", l.Blockchain.ChainId, rollbackBlockNumber-1, models.BlockValid, models.BlockPending).First(&previousBlock)
		if rest.Error != nil {
			log.Errorf("[Handler.RollRackBlock] Previous block by number error: %s\n", errors.WithStack(rest.Error))
			continue
		}
		l.SyncedBlockNumber = previousBlock.BlockNumber
		l.SyncedBlockHash = common.HexToHash(previousBlock.BlockHash)
	}
}
