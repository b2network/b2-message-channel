package config

import (
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func InitB2Rpc(rpcUrl string) *ethclient.Client {
	rpc, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Panicf("b2 client dial error: %s\n", err)
	}
	return rpc
}
