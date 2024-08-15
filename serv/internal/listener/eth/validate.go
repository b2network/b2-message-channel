package eth

import (
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"bsquared.network/b2-message-channel-serv/internal/utils/message"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (l *Listener) validate() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		list, err := l.validatingCallMessage(10)
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
				err = l.validateMessage(message)
				if err != nil {
					log.Errorf("Handle err: %v, %v\n", err, message)
				}
			}(&wg, message)
		}
		wg.Wait()
	}
}

func (l *Listener) validateMessage(msg models.Message) error {
	log.Infof("handle validateMessage\n")
	var signatures []string
	for _, validatorKey := range l.DataMap.ValidatorMap {
		_key, err := crypto.ToECDSA(common.FromHex(validatorKey))
		if err != nil {
			return errors.WithStack(err)
		}
		log.Infof("key:%v\n", _key)
		signature, err := message.SignMessageSend(l.Blockchain.ChainId, l.Blockchain.MessageAddress, msg.FromChainId, common.HexToHash(msg.FromId).Big(), msg.FromSender, msg.ToChainId, msg.ToContractAddress, msg.ToBytes, _key)
		if err != nil {
			return errors.WithStack(err)
		}
		//if signature
		signatures = append(signatures, signature)
	}

	_signatures, err := json.Marshal(&signatures)
	if err != nil {
		return errors.WithStack(err)
	}
	//log.Infof("validateMessage signatures:%s\n", string(_signatures))
	fmt.Println("signatures:", string(_signatures))

	err = l.Db.Model(models.Message{}).Where("`id`=? AND `status`=?", msg.Id, enums.MessageStatusValidating).
		Update("signatures", string(_signatures)).
		Update("status", enums.MessageStatusPending).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (l *Listener) validatingCallMessage(limit int) ([]models.Message, error) {
	var list []models.Message
	err := l.Db.Where("`to_chain_id`=? AND `type`=? AND `status`=?", l.Blockchain.ChainId, enums.MessageTypeCall, enums.MessageStatusValidating).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
