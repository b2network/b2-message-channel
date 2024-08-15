package bitcoin

import (
	"bsquared.network/b2-message-channel-serv/internal/contract/message"
	"bsquared.network/b2-message-channel-serv/internal/enums"
	"bsquared.network/b2-message-channel-serv/internal/models"
	"bsquared.network/b2-message-channel-serv/internal/utils/aa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

func (l *Listener) syncMessage() {
	for {
		var list []models.Deposit
		err := l.Db.Where("status=?", enums.DepositStatusPending).First(&list).Error
		if err != nil {
			continue
		}
		wg := sync.WaitGroup{}
		for _, one := range list {
			wg.Add(1)
			go func(one models.Deposit, wg *sync.WaitGroup) {
				defer wg.Done()
				err := l.handleMessage(one)
				if err != nil {
					log.Errorf("[Handler.syncMessage] Deposit ID: %d , err: %s \n", one.Id, err)
				}
			}(one, &wg)
		}
		wg.Wait()
		time.Sleep(time.Second * 10)
	}
}

func (l *Listener) handleMessage(deposit models.Deposit) error {
	toChainId := l.Blockchain.ToChainId
	toContractAddress := l.Blockchain.ToContractAddress

	depositAddress, err := l.GetDepositAddress(deposit)
	if err != nil {
		return err
	}

	data := message.EncodeSendData(deposit.BtcTxHash, deposit.BtcFrom, depositAddress, decimal.New(deposit.BtcValue, 0))
	msg := models.Message{
		ChainId:             l.Blockchain.ChainId,
		Type:                enums.MessageTypeCall,
		FromChainId:         l.Blockchain.ChainId,
		FromSender:          common.HexToAddress("0x0").Hex(),
		FromContractAddress: common.HexToAddress("0x0").Hex(),
		FromId:              common.HexToHash(deposit.BtcTxHash).Hex(),
		ToChainId:           toChainId,
		ToContractAddress:   toContractAddress,
		ToBytes:             hexutil.Encode(data),
		Signatures:          "{}",
		Status:              enums.MessageStatusValidating,
		Blockchain: models.Blockchain{
			EventId:     deposit.Id,
			BlockTime:   deposit.BtcBlockTime.Unix(),
			BlockNumber: deposit.BtcBlockNumber,
			LogIndex:    deposit.BtcTxIndex,
			TxHash:      common.HexToHash(deposit.BtcTxHash).Hex(),
		},
	}
	err = l.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&msg).Error
		if err != nil {
			return err
		}

		if deposit.BtcFromEvmAddress != "" {
			err = tx.Model(models.Deposit{}).
				Where("id=?", deposit.Id).
				Update("status", enums.MessageStatusValid).Error
		} else {
			err = tx.Model(models.Deposit{}).
				Where("id=?", deposit.Id).
				Update("btc_from_aa_address", depositAddress).
				Update("status", enums.MessageStatusValid).Error
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *Listener) GetDepositAddress(deposit models.Deposit) (string, error) {
	if deposit.BtcFromEvmAddress != "" && common.IsHexAddress(deposit.BtcFromEvmAddress) {
		return deposit.BtcFromEvmAddress, nil
	} else {
		evmAddress, err := aa.BitcoinAddressToEthAddress(l.Particle.AAPubKeyAPI, deposit.BtcFrom,
			l.Particle.Url, l.Particle.ChainId, l.Particle.ProjectUuid, l.Particle.ProjectKey)
		if err != nil {
			return "", err
		}
		return evmAddress, nil
	}
}
