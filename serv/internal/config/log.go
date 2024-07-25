package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLog(logLevel uint32) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.SetLevel(log.Level(logLevel))
}
