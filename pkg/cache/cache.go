package cache

import (
	"github.com/wolfelee/gocomm/pkg/jlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	Redis = "redis"
)

var cacheCfg *CheConfig

type CheConfig struct {
	Redis map[string]*RedisConfig `yaml:"redis"`
}

func Init(cfgFile string) error {
	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		jlog.Error(cfgFile + "file read error")
		return err
	}
	err = yaml.Unmarshal(buf, &cacheCfg)
	if err != nil {
		jlog.Error(cfgFile + "file unmarshal error")
		return err
	}
	redisCfgGroup = cacheCfg.Redis
	if redisCfgGroup != nil {
		err = redisConnGroup()
		if err != nil {
			jlog.Error(err.Error())
			return err
		}
	}
	return nil
}
