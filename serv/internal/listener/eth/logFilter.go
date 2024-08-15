package eth

import (
	"bsquared.network/b2-message-channel-serv/internal/event"
	"bsquared.network/b2-message-channel-serv/internal/models"
	rpc2 "bsquared.network/b2-message-channel-serv/internal/utils/rpc"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
)

func (l *Listener) LogBatchFilter(startBlock, endBlock int64, addresses []common.Address, topics [][]common.Hash) ([]*models.SyncEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(startBlock),
		ToBlock:   big.NewInt(endBlock),
		Topics:    topics,
		Addresses: addresses,
	}
	logs, err := l.RPC.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return l.LogsToEvents(logs, startBlock)
}

func (l *Listener) LogFilter(block models.SyncBlock, addresses []common.Address, topics [][]common.Hash) ([]*models.SyncEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(block.BlockNumber),
		ToBlock:   big.NewInt(block.BlockNumber),
		Topics:    topics,
		Addresses: addresses,
	}
	logs, err := l.RPC.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Infof("[CancelOrder.Handle] Cancel Pending List Length is %d ,block number is %d \n", len(logs), block.BlockNumber)
	return l.LogsToEvents(logs, block.Id)
}

func (l *Listener) LogsToEvents(logs []types.Log, syncBlockId int64) ([]*models.SyncEvent, error) {
	var events []*models.SyncEvent
	blockTimes := make(map[int64]int64)
	for _, vlog := range logs {
		eventHash := event.TopicToHash(vlog, 0)
		contractAddress := vlog.Address
		Event := l.GetEvent(eventHash, contractAddress)
		if Event == nil {
			log.Infof("[LogsToEvents] logs[txHash: %s, contractAddress:%s, eventHash: %s]\n", vlog.TxHash, strings.ToLower(contractAddress.Hex()), eventHash)
			continue
		}

		blockTime := blockTimes[int64(vlog.BlockNumber)]
		if blockTime == 0 {
			blockJson, err := rpc2.HttpPostJson("", l.Blockchain.RpcUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", vlog.BlockNumber)+"\", true],\"id\":1}")
			if err != nil {
				log.Errorf("[Handler.SyncBlock] Syncing block by number error: %s\n", errors.WithStack(err))
				time.Sleep(3 * time.Second)
				continue
			}
			block := rpc2.ParseJsonBlock(string(blockJson))
			//
			//block, err := ctx.RPC.BlockByNumber(context.Background(), big.NewInt(int64(vlog.BlockNumber)))
			//if err != nil {
			//	return nil, errors.WithStack(err)
			//}
			blockTime = int64(block.Timestamp())
			blockTimes[int64(vlog.BlockNumber)] = blockTime
		}

		data, err := Event.Data(vlog)
		if err != nil {
			log.Errorf("[LogsToEvents] logs[txHash: %s, contractAddress:%s, eventHash: %s]\n", vlog.TxHash, strings.ToLower(contractAddress.Hex()), eventHash)
			log.Errorf("[LogsToEvents] data err: %s\n", errors.WithStack(err))
			continue
		}

		events = append(events, &models.SyncEvent{
			ChainId:         l.Blockchain.ChainId,
			SyncBlockId:     syncBlockId,
			BlockTime:       blockTime,
			BlockNumber:     int64(vlog.BlockNumber),
			BlockHash:       vlog.BlockHash.Hex(),
			BlockLogIndexed: int64(vlog.Index),
			TxIndex:         int64(vlog.TxIndex),
			TxHash:          vlog.TxHash.Hex(),
			EventName:       Event.Name(),
			EventHash:       eventHash.Hex(),
			ContractAddress: strings.ToLower(contractAddress.Hex()),
			Data:            data,
			Status:          models.EventPending,
		})
	}
	return events, nil
}
