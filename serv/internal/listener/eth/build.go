package eth

import (
	msg "bsquared.network/b2-message-channel-serv/internal/contract/message"
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"time"
)

func (l *Listener) build() {
	duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
	for {
		list, err := l.pendingCallMessage(10)
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
				err = l.buildMessage(message)
				if err != nil {
					log.Errorf("Handle err: %v, %v\n", err, message)
				}
			}(&wg, message)
		}
		wg.Wait()
	}
}

func (l *Listener) buildMessage(message models.Message) error {
	UserAddress, err := l.BorrowAccount()
	if err != nil {
		log.Errorf("borrow account err: %s\n", err)
		return errors.WithStack(err)
	}
	lock, err := l.LockUser(UserAddress, time.Minute*2)
	if err != nil {
		log.Errorf("lock err: %s\n", err)
		return errors.WithStack(err)
	}
	if !lock {
		log.Infof("load result: %v\n", lock)
		return errors.WithStack(err)
	}
	defer l.UnlockUser(UserAddress)

	gasPrice, err := l.GasPrice()
	if err != nil {
		return errors.WithStack(err)
	}
	//gasPrice := big.NewInt(1000)
	log.Debugf("gasPrice: %v\n", gasPrice)

	toAddress := common.HexToAddress(l.Blockchain.MessageAddress)
	log.Debugf("toAddress: %v\n", toAddress)

	var signatures []string
	err = json.Unmarshal([]byte(message.Signatures), &signatures)
	if err != nil {
		return errors.WithStack(err)
	}

	data := msg.Send(message.FromChainId, common.HexToHash(message.FromId).Big(), message.FromSender, message.ToContractAddress, message.ToBytes, signatures)
	log.Debugf("data: %x\n", data)
	gasLimit, err := l.RPC.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(UserAddress),
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     data,
	})
	if err != nil {
		log.Errorf("Get gasLimit err: %s\n", err)
		return errors.WithStack(err)
	}
	log.Debugf("gasLimit: %v\n", gasLimit)
	err = l.Db.Transaction(func(tx *gorm.DB) error {
		// nonce
		nonce, err := l.GetNonce(UserAddress)
		if err != nil {
			log.Errorf("Get nonce err: %s\n", err)
			return errors.WithStack(err)
		}
		log.Debugf("nonce: %v\n", nonce)
		// signTx
		_signature, err := l.SignTx(UserAddress, nonce, toAddress.Hex(), big.NewInt(0), gasLimit, gasPrice, data, l.Blockchain.ChainId)
		if err != nil {
			log.Errorf("Sign Tx err: %s\n", err)
			return errors.WithStack(err)
		}
		_txHash := crypto.Keccak256Hash(_signature)
		log.Debugf("txHash: %s, signature: %s\n", _txHash, hex.EncodeToString(_signature))

		// create signature
		err = l.CreateSignature(tx, message.ToChainId, message.FromId, UserAddress, int64(nonce), enums.MessageTypeSend, hex.EncodeToString(data), decimal.Zero, hex.EncodeToString(_signature), _txHash.Hex())
		if err != nil {
			log.Errorf("Create signature err: %s\n", err)
			return errors.WithStack(err)
		}
		log.Infof("create signature success\n")
		// broadcast
		message.Status = enums.MessageStatusBroadcast
		err = tx.Save(&message).Error
		if err != nil {
			log.Errorf("Broadcast approve err: %s\n", err)
			return errors.WithStack(err)
		}
		log.Infof("handle approve success\n")
		return nil
	})
	if err != nil {
		log.Errorf("Handle approve err: %s\n", err)
		return errors.WithStack(err)
	}
	return nil
}

func (l *Listener) pendingCallMessage(limit int) ([]models.Message, error) {
	var list []models.Message
	err := l.Db.Where("`to_chain_id`=? AND `type`=? AND `status`=?", l.Blockchain.ChainId, enums.MessageTypeCall, enums.MessageStatusPending).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *Listener) SignTx(accountAddress string, nonce uint64, toAddress string, value *big.Int, gasLimit uint64, gasPrice *big.Int, bytecode []byte, chainID int64) ([]byte, error) {
	senderKey, err := l.getKeyByAddress(accountAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	_signature, err := l._signTx(accountAddress, senderKey, nonce, toAddress, value, gasLimit, gasPrice, bytecode, chainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return _signature, nil
}

func (l *Listener) _signTx(accountAddress, accountKey string, nonce uint64, toAddress string, value *big.Int, gasLimit uint64, gasPrice *big.Int, bytecode []byte, chainID int64) ([]byte, error) {
	_key, err := crypto.ToECDSA(common.FromHex(accountKey))
	if crypto.PubkeyToAddress(_key.PublicKey) != common.HexToAddress(accountAddress) {
		return nil, errors.New(" address and index do not match ")
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		value,
		gasLimit,
		gasPrice,
		bytecode,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), _key)
	if err != nil {
		return nil, err
	}
	ts := types.Transactions{signedTx}
	var rawTxBytes bytes.Buffer
	ts.EncodeIndex(0, &rawTxBytes)
	return rawTxBytes.Bytes(), nil
}

func (l *Listener) GetNonce(userAddress string) (uint64, error) {
	var signature models.Signature
	err := l.Db.Raw("SELECT * FROM signatures WHERE `status`!=? AND `address`=? ORDER BY nonce DESC FOR UPDATE ",
		enums.SignatureStatusInvalid, userAddress).First(&signature).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		nonce, err := l.RPC.NonceAt(context.Background(), common.HexToAddress(userAddress), nil)
		if err != nil {
			return 0, err
		}
		log.Infof("==========nonce: %v\n", nonce)
		return nonce, nil
	} else {
		if signature.Status != enums.SignatureStatusSuccess && signature.Status != enums.SignatureStatusFailed {
			return 0, errors.New("The current user has pending transactions ")
		}
		nonce, err := l.RPC.NonceAt(context.Background(), common.HexToAddress(userAddress), nil)
		if err != nil {
			return 0, err
		}
		log.Infof("==========nonce: %v\n", nonce)
		return nonce, nil
	}
}

func (l *Listener) CreateSignature(tx *gorm.DB, chainId int64, referId string, address string, nonce int64, signatureType enums.MessageType, data string, value decimal.Decimal, signature string, txHash string) error {
	err := tx.Create(&models.Signature{
		ChainId:   chainId,
		ReferId:   referId,
		Address:   address,
		Nonce:     nonce,
		Type:      signatureType,
		Data:      data,
		Value:     value,
		Signature: signature,
		Status:    enums.SignatureStatusPending,
		Blockchain: models.Blockchain{
			EventId:     0,
			BlockTime:   0,
			BlockNumber: 0,
			LogIndex:    0,
			TxHash:      txHash,
		},
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *Listener) GasPrice() (*big.Int, error) {
	gasPrice, err := l.RPC.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return gasPrice, nil
}

func (l *Listener) LockUser(key string, duration time.Duration) (bool, error) {
	return l.Cache.Client.SetNX(context.Background(), key, true, duration).Result()
}

func (l *Listener) UnlockUser(key string) error {
	_, err := l.Cache.Client.Del(context.Background(), key).Result()
	return err
}

func (l *Listener) BorrowAccount() (string, error) {
	for address, _ := range l.DataMap.SenderMap {
		return address, nil
	}
	return "", errors.New("account not found")
}

func (l *Listener) getKeyByAddress(accountAddress string) (string, error) {
	if key, ok := l.DataMap.SenderMap[accountAddress]; ok {
		return key, nil
	} else {
		return "", errors.New("account not found")
	}
}
