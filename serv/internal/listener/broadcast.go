package listener

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"context"
	"encoding/hex"
	_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (l *Listener) broadcast() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		var signatures []models.Signature
		err := l.Db.Where("`chain_id`=? AND `status`=?", l.Blockchain.ChainId, enums.SignatureStatusPending).Order("id").Limit(100).Find(&signatures).Error
		if err != nil {
			log.Errorf("Get signatures  err: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(signatures) == 0 {
			log.Infof("Get signatures length is 0\n")
			time.Sleep(duration)
			continue
		}

		var wg sync.WaitGroup
		for _, signature := range signatures {
			wg.Add(1)
			go func(_wg *sync.WaitGroup, signature models.Signature) {
				defer _wg.Done()
				err = l.broadcastSignature(signature)
				if err != nil {
					log.Errorf("Broadcast signature err[%d]: %s\n", signature.Id, err)
				}
			}(&wg, signature)
		}
		wg.Wait()
	}
}

func (l *Listener) broadcastSignature(signature models.Signature) error {
	log.Infof("signature: %v\n", signature)
	err := l._broadcast(signature.Signature)
	if err != nil && err.Error() != "nonce too low" && err.Error() != "already known" {
		log.Errorf("Broadcast err[%d]: %s\n", signature.Id, err)
		return err
	} else if err != nil && err.Error() == "nonce too low" {
		// _, err := ctx.RPC.TransactionReceipt(context.Background(), common.HexToHash(signature.TxHash))
		// if err != nil {
		//	log.Errorf("Get TransactionReceipt err[%d]: %s\n", signature.Id, err)
		//	return err
		// }
		err = l.Db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			log.Errorf("update signature status broadcast err[%d]: %s\n", signature.Id, err)
			return err
		}
	} else if err != nil && err.Error() == "already known" {
		err = l.Db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			log.Errorf("update signature status broadcast err[%d]: %s\n", signature.Id, err)
			return err
		}
	} else {
		err = l.Db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			log.Errorf("Broadcast signature err[%d]: %s\n", signature.Id, err)
			return err
		}
		log.Infof("signature success\n")
		return nil
	}
	return nil
}

func (l *Listener) _broadcast(signature string) error {
	rawTxBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Errorf("[Broadcast]Decode string err: %s\n", err)
		return err
	}
	tx := new(_types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)
	err = l.RPC.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Errorf("[Broadcast]Send transaction err: %s\n", err)
		return err
	}
	return nil
}
