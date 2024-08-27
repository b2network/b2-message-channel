package config

import (
	"bsquared.network/b2-message-channel-listener/internal/enums"
	"github.com/spf13/viper"
	"strings"
)

type AppConfig struct {
	Log       LogConfig
	Database  Database
	Bsquared  Blockchain
	Bitcoin   Blockchain
	Arbitrum  Blockchain
	Particle  Particle
	Validator Validator
	Builder   Builder
}

type Builder struct {
	Bsquared string
	Arbitrum string
}

type Validator struct {
	Bsquared string
	Arbitrum string
}

type LogConfig struct {
	Level uint32
}

type Database struct {
	UserName string
	Password string
	Host     string
	Port     int64
	DbName   string
	LogLevel int64
}

type Blockchain struct {
	Name              string
	ChainType         enums.ChainType
	ChainId           int64
	RpcUrl            string
	SafeBlockNumber   int64
	ListenAddress     string
	BlockInterval     int64
	Mainnet           bool
	ToChainId         int64
	ToContractAddress string
	BtcUser           string
	BtcPass           string
}

type Particle struct {
	AAPubKeyAPI string
	Url         string
	ChainId     int
	ProjectUuid string
	ProjectKey  string
}

func LoadConfig() AppConfig {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	var config AppConfig

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.SetEnvPrefix("app")

	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
