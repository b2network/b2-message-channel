package listener

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

func (l *Listener) syncLatestBlockNumber(duration time.Duration) {
	for {
		latest, err := l.RPC.BlockNumber(context.Background())
		if err != nil {
			log.Errorf("[Handle.LatestBlackNumber][%d]Syncing latest block number error: %s\n", l.Blockchain.ChainId, err)
			time.Sleep(duration)
			continue
		}
		l.LatestBlockNumber = int64(latest)
		log.Infof("[Handle.LatestBlackNumber][%d]Syncing latest block number: %d \n", l.Blockchain.ChainId, latest)
		time.Sleep(duration)
	}
}
