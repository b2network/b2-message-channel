package bitcoin

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (l *Listener) syncLatestBlockNumber() {
	for {
		duration := time.Millisecond * time.Duration(l.Blockchain.BlockInterval)
		latest, err := l.Rpc.GetBlockCount()
		if err != nil {
			log.Errorf("[Handle.LatestBlackNumber][%d]Syncing latest block number error: %s\n", l.Blockchain.ChainId, err)
			time.Sleep(duration)
			continue
		}
		l.LatestBlockNumber = int64(latest) - 3
		log.Infof("[Handle.LatestBlackNumber][%d]Syncing latest block number: %d \n", l.Blockchain.ChainId, latest)
		time.Sleep(duration)
	}
}
