package btc

//import "github.com/btcsuite/btcd/rpcclient"
//
//type Config struct {
//	RPCHost    string
//	RPCPort    string
//	RPCUser    string
//	RPCPass    string
//	DisableTLS bool
//}
//
//func NewBtcClient(bitcoinCfg Config) (*rpcclient.Client, error) {
//	bclient, err := rpcclient.New(&rpcclient.ConnConfig{
//		Host:         bitcoinCfg.RPCHost + ":" + bitcoinCfg.RPCPort,
//		User:         bitcoinCfg.RPCUser,
//		Pass:         bitcoinCfg.RPCPass,
//		HTTPPostMode: true,                  // Bitcoin core only supports HTTP POST mode
//		DisableTLS:   bitcoinCfg.DisableTLS, // Bitcoin core does not provide TLS by default
//	}, nil)
//	return bclient, err
//}
//
//func NewBitcoinIndexer() {
//	bidxer, err := bitcoin.NewBitcoinIndexer(bidxLogger, bclient, bitcoinParam, bitcoinCfg.IndexerListenAddress, bitcoinCfg.IndexerListenTargetConfirmations)
//	if err != nil {
//		logger.Errorw("failed to new bitcoin indexer indexer", "error", err.Error())
//		return err
//	}
//}
