package cache

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"net"
)

var (
	RCacheGroup   map[string]*redis.Client
	RCache        *redis.Client
	redisCfgGroup map[string]*RedisConfig
)

type RedisMaster struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type RedisConfig struct {
	Cluster  string       `yaml:"cluster"`
	Master   *RedisMaster `yaml:"master"`
	Sentinel []string     `yaml:"sentinel"`
}

func Use(cacheName string) *redis.Client {
	c := RCacheGroup[cacheName]
	if c == nil {
		jlog.Error("cache " + cacheName + " is nil")
	}
	return c
}

func redisConnGroup() error {
	RCacheGroup = make(map[string]*redis.Client)
	for g, redisCfg := range redisCfgGroup {
		if redisCfg.Master != nil {
			master := redisCfg.Master
			switch redisCfg.Cluster {
			case "", "standalone":
				group := redis.NewClient(&redis.Options{
					Network:  master.Protocol,
					Addr:     net.JoinHostPort(master.Host, master.Port),
					Password: master.Password,
					DB:       master.Db,
				})
				if g == "default" {
					RCache = group
				}
				RCacheGroup[g] = group
			case "sentinel":
				if redisCfg.Sentinel != nil {
					group := redis.NewFailoverClient(&redis.FailoverOptions{
						MasterName:    master.Host,
						Password:      master.Password,
						DB:            master.Db,
						SentinelAddrs: redisCfg.Sentinel,
					})
					if g == "default" {
						RCache = group
					}
					RCacheGroup[g] = group
				} else {
					return errors.New("cache config setting error")
				}
			default:
				return errors.New("cache config - cluster setting error")
			}
		} else {
			return errors.New("cache config - master must be setting")
		}
	}

	if len(RCacheGroup) > 0 {
		for name, cc := range RCacheGroup {
			_, err := cc.Ping().Result()
			if err != nil {
				jlog.Error(err.Error())
			} else {
				jlog.Info("successful connection to redis-server:" + name)
			}
		}
	}
	return nil
}

func Ping() error {
	for _, cc := range RCacheGroup {
		_, err := cc.Ping().Result()
		if err != nil {
			jlog.Error(err.Error())
			return err
		}
	}
	return nil
}
