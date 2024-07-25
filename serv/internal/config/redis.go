package config

import (
	"bsquared.network/b2-message-channel-serv/internal/utils"
	localcache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitCache(cfg AppConfig) *Cache {
	err := utils.InitRedis(cfg.Redis.IsClusterMode, cfg.Redis.ClusterAddresses, cfg.Redis.Password,
		cfg.Redis.TlsInsecureSkip, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	if err != nil {
		panic(err)
	}
	c, err := newCache(utils.GetClient())
	if err != nil {
		panic(err)
	}
	return c
	return nil
}

type Cache struct {
	Client redis.Cmdable
	local  *localcache.Cache
}

func newCache(client redis.Cmdable) (*Cache, error) {
	return &Cache{
		Client: client,
		local:  localcache.New(15*time.Minute, 20*time.Minute),
	}, nil
}
