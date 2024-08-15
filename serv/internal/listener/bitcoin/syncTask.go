package bitcoin

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"bsquared.network/b2-message-channel-serv/internal/types"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

func (l *Listener) syncTask() {
	for {
		duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
		var tasks []models.SyncTask
		err := l.Db.Where("`chain_type`=? AND `chain_id`=? AND `status`=?", enums.ChainTypeBitcoin, l.Blockchain.ChainId, models.SyncTaskPending).Limit(20).Find(&tasks).Error
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
				err := l.HandleTask(task)
				if err != nil {
					log.Errorf("[Handler.SyncTask] HandleTask ID: %d , err: %s \n", task.Id, err)
				}
			}(task, &wg)
		}
		wg.Wait()
	}
}

func (l *Listener) HandleTask(task models.SyncTask) error {

	var (
		currentBlock   int64 // index current block number
		currentTxIndex int64 // index current block tx index
	)

	currentBlock = task.LatestBlock
	if currentBlock < task.StartBlock {
		currentBlock = task.StartBlock
	}
	currentTxIndex = task.LatestTx

	for {
		if l.LatestBlockNumber <= currentBlock {
			//<-ticker.C
			//ticker.Reset(NewBlockWaitTimeout)
			//
			//// update latest block
			//latestBlock, err = bis.txIdxr.LatestBlock()
			//if err != nil {
			//	bis.log.Errorw("bitcoin indexer latestBlock", "error", err.Error())
			//}
			time.Sleep(time.Second * 10)
			continue
		}
		if currentTxIndex == 0 {
			currentBlock++
		} else {
			currentTxIndex++
		}
		for i := currentBlock; i <= l.LatestBlockNumber; i++ {
			txResults, blockHeader, err := l.ParseBlock(i, currentTxIndex)
			if err != nil {
				if errors.Is(err, ErrTargetConfirmations) {
					//bis.log.Warnw("parse block confirmations", "error", err.Error(), "currentBlock", i, "currentTxIndex", currentTxIndex)
					time.Sleep(NewBlockWaitTimeout)
				} else {
					//bis.log.Errorw("parse block unknown err", "error", err.Error(), "currentBlock", i, "currentTxIndex", currentTxIndex)
				}
				if currentTxIndex == 0 {
					currentBlock = i - 1
				} else {
					currentBlock = i
					currentTxIndex--
				}

				break
			}
			if len(txResults) > 0 {
				currentBlock, currentTxIndex, err = l.HandleResults(txResults, task, blockHeader.Timestamp, i)
				if err != nil {
					//bis.log.Errorw("failed to handle results", "error", err,
					//	"currentBlock", currentBlock, "currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
					rollback := true
					// not duplicated key, rollback index
					if pgErr, ok := err.(*pgconn.PgError); ok {
						// 23505 duplicate key value violates unique constraint , continue
						if pgErr.Code == "23505" {
							rollback = false
						}
					}

					if rollback {
						if currentTxIndex == 0 {
							currentBlock = i - 1
						} else {
							currentBlock = i
							currentTxIndex--
						}
						break
					}
				}
			}
			currentBlock = i
			currentTxIndex = 0
			task.LatestBlock = currentBlock
			task.LatestTx = currentTxIndex
			if err := l.Db.Save(&task).Error; err != nil {
				//bis.log.Errorw("failed to save bitcoin index block", "error", err, "currentBlock", i,
				//	"currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
				// rollback
				currentBlock = i - 1
				break
			}
			//bis.log.Infow("bitcoin indexer parsed", "currentBlock", i,
			//	"currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
			time.Sleep(IndexBlockTimeout)
		}
	}
	return nil
}

func (l *Listener) ParseBlock(height int64, txIndex int64) ([]*types.BitcoinTxParseResult, *wire.BlockHeader, error) {
	blockResult, err := l.getBlockByHeight(height)
	if err != nil {
		return nil, nil, err
	}

	blockParsedResult := make([]*types.BitcoinTxParseResult, 0)
	for k, v := range blockResult.Transactions {
		if int64(k) < txIndex {
			continue
		}

		//b.logger.Debugw("parse block", "k", k, "height", height, "txIndex", txIndex, "tx", v.TxHash().String())

		parseTxs, err := l.parseTx(v, k)
		if err != nil {
			return nil, nil, err
		}
		if parseTxs != nil {
			blockParsedResult = append(blockParsedResult, parseTxs)
		}
	}

	return blockParsedResult, &blockResult.Header, nil
}

// getBlockByHeight returns a raw block from the server given its height
func (l *Listener) getBlockByHeight(height int64) (*wire.MsgBlock, error) {
	blockhash, err := l.Rpc.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	msgBlock, err := l.Rpc.GetBlock(blockhash)
	if err != nil {
		return nil, err
	}
	return msgBlock, nil
}
func (l *Listener) parseTx(txResult *wire.MsgTx, index int) (*types.BitcoinTxParseResult, error) {
	listenAddress := false
	var totalValue int64
	tos := make([]types.BitcoinTo, 0)
	for _, v := range txResult.TxOut {
		pkAddress, err := l.parseAddress(v.PkScript)
		if err != nil {
			if errors.Is(err, ErrParsePkScript) {
				continue
			}
			// parse null data
			if errors.Is(err, ErrParsePkScriptNullData) {
				nullData, err := l.parseNullData(v.PkScript)
				if err != nil {
					continue
				}
				tos = append(tos, types.BitcoinTo{
					Type:     types.BitcoinToTypeNullData,
					NullData: nullData,
				})
			} else {
				return nil, err
			}
		} else {
			parseTo := types.BitcoinTo{
				Address: pkAddress,
				Value:   v.Value,
				Type:    types.BitcoinToTypeNormal,
			}
			tos = append(tos, parseTo)
		}
		// if pk address eq dest listened address, after parse from address by vin prev tx
		if pkAddress == l.listenAddress.EncodeAddress() {
			listenAddress = true
			totalValue += v.Value
		}
	}
	if listenAddress {
		fromAddress, err := l.parseFromAddress(txResult)
		if err != nil {
			return nil, fmt.Errorf("vin parse err:%w", err)
		}

		// TODO: temp fix, if from is listened address, continue
		if len(fromAddress) == 0 {
			//b.logger.Warnw("parse from address empty or nonsupport tx type",
			//	"txId", txResult.TxHash().String(),
			//	"listenAddress", b.listenAddress.EncodeAddress())
			return nil, nil
		}

		return &types.BitcoinTxParseResult{
			TxID:   txResult.TxHash().String(),
			TxType: TxTypeTransfer,
			Index:  int64(index),
			Value:  totalValue,
			From:   fromAddress,
			To:     l.listenAddress.EncodeAddress(),
			Tos:    tos,
		}, nil
	}
	return nil, nil
}

func (l *Listener) parseAddress(pkScript []byte) (string, error) {
	pk, err := txscript.ParsePkScript(pkScript)
	if err != nil {
		scriptClass := txscript.GetScriptClass(pkScript)
		if scriptClass == txscript.NullDataTy {
			return "", ErrParsePkScriptNullData
		}
		return "", fmt.Errorf("%w:%s", ErrParsePkScript, err.Error())
	}

	if pk.Class() == txscript.NullDataTy {
		return "", ErrParsePkScriptNullData
	}

	//  encodes the script into an address for the given chain.
	pkAddress, err := pk.Address(l.ChainParams)
	if err != nil {
		return "", fmt.Errorf("PKScript to address err:%w", err)
	}
	return pkAddress.EncodeAddress(), nil
}

// parseNullData from pkscript parse null data
func (l *Listener) parseNullData(pkScript []byte) (string, error) {
	if !txscript.IsNullData(pkScript) {
		return "", ErrParsePkScriptNotNullData
	}
	return hex.EncodeToString(pkScript[1:]), nil
}

func (l *Listener) parseFromAddress(txResult *wire.MsgTx) (fromAddress []types.BitcoinFrom, err error) {
	for _, vin := range txResult.TxIn {
		// get prev tx hash
		prevTxID := vin.PreviousOutPoint.Hash
		vinResult, err := l.Rpc.GetRawTransaction(&prevTxID)
		if err != nil {
			return nil, fmt.Errorf("vin get raw transaction err:%w", err)
		}
		if len(vinResult.MsgTx().TxOut) == 0 {
			return nil, fmt.Errorf("vin txOut is null")
		}
		vinPKScript := vinResult.MsgTx().TxOut[vin.PreviousOutPoint.Index].PkScript
		//  script to address
		vinPkAddress, err := l.parseAddress(vinPKScript)
		if err != nil {
			//b.logger.Errorw("vin parse address", "error", err)
			if errors.Is(err, ErrParsePkScript) || errors.Is(err, ErrParsePkScriptNullData) {
				continue
			}
			return nil, err
		}

		fromAddress = append(fromAddress, types.BitcoinFrom{
			Address: vinPkAddress,
			Type:    types.BitcoinFromTypeBtc,
		})
	}
	return fromAddress, nil
}

func (l *Listener) HandleResults(
	txResults []*types.BitcoinTxParseResult,
	syncTask models.SyncTask,
	btcBlockTime time.Time,
	currentBlock int64,
) (int64, int64, error) {
	for _, v := range txResults {
		// if from is listen address, skip
		if l.ToInFroms(v.From, v.To) {
			//bis.log.Infow("current transaction from is listen address", "currentBlock", currentBlock, "currentTxIndex", v.Index, "data", v)
			continue
		}

		syncTask.LatestBlock = currentBlock
		syncTask.LatestTx = v.Index
		// write db
		err := l.SaveParsedResult(
			v,
			currentBlock,
			models.DepositB2TxStatusPending,
			btcBlockTime,
			syncTask,
		)
		if err != nil {
			//bis.log.Errorw("failed to save bitcoin index tx", "error", err,
			//	"data", v)
			return currentBlock, v.Index, err
		}
		//bis.log.Infow("save bitcoin index tx success", "currentBlock", currentBlock, "currentTxIndex", v.Index, "data", v)
		time.Sleep(IndexTxTimeout)
	}
	return currentBlock, 0, nil
}

func (l *Listener) ToInFroms(a []types.BitcoinFrom, s string) bool {
	for _, i := range a {
		if i.Address == s {
			return true
		}
	}
	return false
}

func (l *Listener) SaveParsedResult(
	parseResult *types.BitcoinTxParseResult,
	btcBlockNumber int64,
	b2TxStatus int,
	btcBlockTime time.Time,
	syncTask models.SyncTask,
) error {
	// write db
	err := l.Db.Transaction(func(tx *gorm.DB) error {
		if len(parseResult.From) == 0 {
			return fmt.Errorf("parse result from empty")
		}

		if len(parseResult.To) == 0 {
			return fmt.Errorf("parse result to empty")
		}

		if len(parseResult.Tos) == 0 {
			return fmt.Errorf("parse result to empty")
		}

		//bis.log.Infow("parseResult:", "result", parseResult)
		existsEvmAddressData := false // The evm address is processed only if it exists. Otherwise, aa is used
		parsedEvmAddress := ""        // evm address
		for _, v := range parseResult.Tos {
			// only handle first null data
			if existsEvmAddressData {
				continue
			}
			if v.Type == types.BitcoinToTypeNullData {
				decodeNullData, err := hex.DecodeString(v.NullData)
				if err != nil {
					//bis.log.Errorw("decode null data err", "error", err, "nullData", v.NullData)
					continue
				}
				evmAddress := bytes.TrimSpace(decodeNullData[1:])
				if common.IsHexAddress(string(evmAddress)) {
					existsEvmAddressData = true
					parsedEvmAddress = string(evmAddress)
					for k := range parseResult.From {
						parseResult.From[k].Type = types.BitcoinFromTypeEvm
						parseResult.From[k].EvmAddress = parsedEvmAddress
					}
				}
			}
		}
		froms, err := json.Marshal(parseResult.From)
		if err != nil {
			return err
		}
		tos, err := json.Marshal(parseResult.Tos)
		if err != nil {
			return err
		}
		// if existed, update deposit record
		var deposit models.Deposit
		err = tx.
			Set("gorm:query_option", "FOR UPDATE").
			First(&deposit,
				fmt.Sprintf("%s = ?", models.Deposit{}.Column().BtcTxHash),
				parseResult.TxID).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			deposit := models.Deposit{
				BtcBlockNumber: btcBlockNumber,
				BtcTxIndex:     parseResult.Index,
				BtcTxHash:      parseResult.TxID,
				BtcFrom:        parseResult.From[0].Address,
				BtcTos:         string(tos),
				BtcTo:          parseResult.To,
				BtcValue:       parseResult.Value,
				BtcFroms:       string(froms),
				B2TxStatus:     b2TxStatus,
				BtcBlockTime:   btcBlockTime,
				B2TxRetry:      0,
				ListenerStatus: models.ListenerStatusSuccess,
				CallbackStatus: models.CallbackStatusPending,
			}
			if existsEvmAddressData {
				deposit.BtcFromEvmAddress = parsedEvmAddress
			}
			err = tx.Create(&deposit).Error
			if err != nil {
				//bis.log.Errorw("failed to save tx parsed result", "error", err)
				return err
			}
		} else if deposit.CallbackStatus == models.CallbackStatusSuccess &&
			deposit.ListenerStatus == models.ListenerStatusPending {
			if deposit.BtcValue != parseResult.Value || deposit.BtcFrom != parseResult.From[0].Address {
				return fmt.Errorf("invalid parameter")
			}
			// if existed, update deposit record
			updateFields := map[string]interface{}{
				models.Deposit{}.Column().BtcBlockNumber: btcBlockNumber,
				models.Deposit{}.Column().BtcTxIndex:     parseResult.Index,
				models.Deposit{}.Column().BtcFroms:       string(froms),
				models.Deposit{}.Column().BtcTos:         string(tos),
				models.Deposit{}.Column().BtcBlockTime:   btcBlockTime,
				models.Deposit{}.Column().ListenerStatus: models.ListenerStatusSuccess,
			}
			if existsEvmAddressData {
				updateFields[models.Deposit{}.Column().BtcFromEvmAddress] = parsedEvmAddress
			}
			err = tx.Model(&models.Deposit{}).Where("id = ?", deposit.Id).Updates(updateFields).Error
			if err != nil {
				//bis.log.Errorw("failed to update tx parsed result", "error", err)
				return err
			}
		}

		if err := tx.Save(&syncTask).Error; err != nil {
			//bis.log.Errorw("failed to save bitcoin tx index", "error", err)
			return err
		}
		return nil
	})
	return err
}
