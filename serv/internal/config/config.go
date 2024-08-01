package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"strings"
)

type Server struct {
	Port            string
	LogLevel        uint32
	BybitPrivateKey string
}

type Database struct {
	UserName string
	Password string
	Host     string
	Port     int64
	DbName   string
	LogLevel int64
}

type Redis struct {
	Host             string
	Port             string
	Password         string
	DB               int
	TlsInsecureSkip  bool
	IsClusterMode    bool
	ClusterAddresses string
}

type Blockchain struct {
	ChainId         int64
	RpcUrl          string
	InitBlockNumber int64
	InitBlockHash   string
	MessageAddress  string
	Events          string
	Senders         string
	Validators      string
	BlockInterval   int64
}

type SentryConfig struct {
	Url        string
	SampleRate float64
	Env        string
	Release    string
}

type AppConfig struct {
	Server     Server
	Database   Database
	Redis      Redis
	Blockchain []Blockchain
	Sentry     SentryConfig
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
	err := v.BindEnv("server.port")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("server.loglevel")

	err = v.BindEnv("api.secret")
	if err != nil {
		return AppConfig{}
	}

	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("database.host")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("database.dbname")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("database.username")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("database.password")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("database.loglevel")
	if err != nil {
		return AppConfig{}
	}

	err = v.BindEnv("redis.host")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.port")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.password")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.db")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.tls_insecure_skip")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.is_cluster_mode")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("redis.cluster_addresses")
	if err != nil {
		return AppConfig{}
	}

	err = v.BindEnv("sentry.url")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("sentry.SampleRate")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("sentry.env")
	if err != nil {
		return AppConfig{}
	}
	err = v.BindEnv("sentry.release")
	if err != nil {
		return AppConfig{}
	}

	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}
	var blockchains []Blockchain
	blockchainsJson := v.GetString("BLOCKCHAINS")
	err = json.Unmarshal([]byte(blockchainsJson), &blockchains)
	if err != nil {
		panic(err)
	}
	config.Blockchain = blockchains

	return config
}
