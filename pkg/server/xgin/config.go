package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wolfelee/gocomm/pkg/conf"
)

// ModName ..
const ModName = "server.gin"

// Config HTTP config
type Config struct {
	Host           string
	Port           int
	Deployment     string
	Mode           string
	ServiceAddress string

	SlowQueryThresholdInMilli int64
}

func DefaultConfig() *Config {
	return &Config{
		Host:                      "0.0.0.0",
		Port:                      80,
		Mode:                      gin.ReleaseMode,
		SlowQueryThresholdInMilli: 500,
	}
}

func StdConfig() *Config {
	var config = DefaultConfig()

	mode := conf.GetString("httpMode")
	if len(mode) > 0 {
		config.Mode = mode
	}

	port := conf.GetInt64("httpPort")
	if port > 0 {
		config.Port = int(port)
	}

	return config
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	server := newServer(config)
	server.Use(gin.Recovery())
	//server.Use(recoverMiddleware(config.logger, config.SlowQueryThresholdInMilli))
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
