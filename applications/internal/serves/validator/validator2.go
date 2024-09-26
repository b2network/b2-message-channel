package validator

//
//import (
//	"bsquared.network/b2-message-channel-applications/internal/config"
//	"bsquared.network/b2-message-channel-applications/internal/enums"
//	"bsquared.network/b2-message-channel-applications/internal/models"
//	"bsquared.network/b2-message-channel-applications/internal/utils/ethereum/message"
//	"bsquared.network/b2-message-channel-applications/internal/utils/log"
//	"encoding/json"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/pkg/errors"
//	"gorm.io/gorm"
//	"sync"
//	"time"
//)
//
//type Validator struct {
//	config config.Blockchain
//	db     *gorm.DB
//	rpc    *ethclient.Client
//	key    string
//	logger *log.Logger
//}
//
//func NewValidator(key string, logConfig config.LogConfig, config config.Blockchain, db *gorm.DB, rpc *ethclient.Client) *Validator {
//	return &Validator{
//		config: config,
//		db:     db,
//		rpc:    rpc,
//		key:    key,
//		logger: log.NewLogger(config.Name, logConfig.Level),
//	}
//}
//
//func (v *Validator) Start() {
//	go v.validate()
//}
//
//func (v *Validator) validate() {
//	duration := time.Millisecond * time.Duration(v.config.BlockInterval)
//	for {
//		list, err := v.validatingCallMessage(10)
//		if err != nil {
//			v.logger.Errorf("Get pending call message error: %v\n", err)
//			time.Sleep(duration)
//			continue
//		}
//		if len(list) == 0 {
//			v.logger.Error("Get pending call message list length is 0\n")
//			time.Sleep(duration)
//			continue
//		}
//
//		var wg sync.WaitGroup
//		for _, message := range list {
//			wg.Add(1)
//			go func(_wg *sync.WaitGroup, message models.Message) {
//				defer _wg.Done()
//				err = v.validateMessage(message)
//				if err != nil {
//					v.logger.Errorf("Validate message error: %v\n", err)
//				}
//			}(&wg, message)
//		}
//		wg.Wait()
//	}
//}
//
//func (v *Validator) validateMessage(msg models.Message) error {
//	var signatures []string
//	//var ValidatorMap map[string]string
//	//for _, validatorKey := range ValidatorMap {
//	_key, err := crypto.ToECDSA(common.FromHex(v.key))
//	if err != nil {
//		v.logger.Errorf("ToECDSA error: %v\n", err)
//		return errors.WithStack(err)
//	}
//	signature, err := message.SignMessageSend(v.config.ChainId, v.config.ListenAddress, msg.FromChainId, common.HexToHash(msg.FromId).Big(), msg.FromSender, msg.ToChainId, msg.ToContractAddress, msg.ToBytes, _key)
//	if err != nil {
//		v.logger.Errorf("Sign message error: %v\n", err)
//		return errors.WithStack(err)
//	}
//	signatures = append(signatures, signature)
//	//}
//
//	_signatures, err := json.Marshal(&signatures)
//	if err != nil {
//		v.logger.Errorf("Marshal signatures error: %v\n", err)
//		return errors.WithStack(err)
//	}
//
//	err = v.db.Model(models.Message{}).Where("`id`=? AND `status`=?", msg.Id, enums.MessageStatusValidating).
//		Update("signatures", string(_signatures)).
//		Update("status", enums.MessageStatusPending).Error
//	if err != nil {
//		v.logger.Errorf("Update message error: %v\n", err)
//		return errors.WithStack(err)
//	}
//	return nil
//}
//
//func (v *Validator) validatingCallMessage(limit int) ([]models.Message, error) {
//	var list []models.Message
//	err := v.db.Where("`to_chain_id`=? AND `type`=? AND `status`=?", v.config.ChainId, enums.MessageTypeCall, enums.MessageStatusValidating).Limit(limit).Find(&list).Error
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return list, nil
//}
