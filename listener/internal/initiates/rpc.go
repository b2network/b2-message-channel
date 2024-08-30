package initiates

import (
	"bsquared.network/b2-message-channel-listener/internal/config"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func InitEthereumRpc(config config.Blockchain) *ethclient.Client {
	rpc, err := ethclient.Dial(config.RpcUrl)
	if err != nil {
		log.Panicf("client dial error: %s\n", err)
	}
	return rpc
}

func InitBitcoinRpc(config config.Blockchain) *rpcclient.Client {
	rpc, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         config.RpcUrl,
		User:         config.BtcUser,
		Pass:         config.BtcPass,
		HTTPPostMode: true,              // Bitcoin core only supports HTTP POST mode
		DisableTLS:   config.DisableTLS, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		log.Panicf("client dial error: %s\n", err)
	}
	return rpc
}
