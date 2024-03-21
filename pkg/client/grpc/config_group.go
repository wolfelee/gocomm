package grpc

import (
	"github.com/wolfelee/gocomm/pkg/jlog"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	grpcCfg   *ConfigGroup
	GrpcGroup map[string]*grpc.ClientConn
)

// Config ...
type ConfigGroup struct {
	DialTimeout int                    `yaml:"dialTimeout"`
	Debug       bool                   `yaml:"debug"`
	GrpcList    map[string]*grpcConfig `yaml:"grpc"`
}

type grpcConfig struct {
	Addr        string `yaml:"addr"`
	Port        string `yaml:"port"`
	ReadTimeout int    `yaml:"readTimeout"`
}

// StdConfig ...
func StdConfig(configPath string) *ConfigGroup {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		jlog.Warn(configPath + " file read error")
	}
	err = yaml.Unmarshal(buf, &grpcCfg)
	if err != nil {
		jlog.Warn(configPath + "file unmarshal error")
	}
	return grpcCfg
}

// Build ...
func (cg *ConfigGroup) BuildGroup() map[string]*grpc.ClientConn {
	var groups = make(map[string]*grpc.ClientConn)
	if grpcCfg == nil {
		jlog.Error("db config setting error")
	}
	for gName, configS := range grpcCfg.GrpcList {
		var config = DefaultConfig()
		config.Address = configS.Addr + ":" + configS.Port
		config.Debug = grpcCfg.Debug
		config.ReadTimeout = time.Duration(configS.ReadTimeout) * time.Second

		config.DialOptions = append(config.DialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor()),
		)
		cc, err := newGRPCClient(config)
		if err != nil {
			if config.OnDialError == "panic" {
				jlog.Error("dial grpc server" + err.Error())
			} else {
				jlog.Error("dial grpc server" + err.Error())
			}
			panic("grpc connect failed:" + err.Error())
		}
		groups[gName] = cc
	}
	return groups
}
