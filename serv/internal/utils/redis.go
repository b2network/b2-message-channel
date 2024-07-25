package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
)

type Manager struct {
	singleNodeClient *redis.Client
	clusterClient    *redis.ClusterClient
	Client           redis.Cmdable
}

var (
	ManagerClient *Manager
)

func GetClient() redis.Cmdable {
	return ManagerClient.Client
}

func InitRedis(redisIsClusterMode bool, redisClusterAddresses, redisPassword string, tlsInsecureSkip bool, redisHost, redisPort string, redisDB int) error {
	m := &Manager{}
	if redisIsClusterMode {
		//集群模式
		opt := &redis.ClusterOptions{
			Addrs:    strings.Split(redisClusterAddresses, ","),
			Password: redisPassword,
		}
		if tlsInsecureSkip {
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		cli := redis.NewClusterClient(opt)
		_, err := cli.Ping(context.Background()).Result()
		if err != nil {
			return err
		}
		m.clusterClient = cli
		m.Client = cli

	} else {

		//单节点模式
		opt := &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisPassword,
			DB:       redisDB,
		}
		if tlsInsecureSkip {
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		cli := redis.NewClient(opt)
		_, err := cli.Ping(context.Background()).Result()
		if err != nil {
			return err
		}
		m.singleNodeClient = cli
		m.Client = cli
	}

	ManagerClient = m
	return nil
}
