package grpc

import (
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/util/xtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	Name         string // config's name
	BalancerName string
	Address      string
	Block        bool
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	Direct       bool
	OnDialError  string // panic | error
	KeepAlive    *keepalive.ClientParameters
	DialOptions  []grpc.DialOption

	SlowThreshold time.Duration

	Debug                     bool
	DisableTraceInterceptor   bool
	DisableAidInterceptor     bool
	DisableTimeoutInterceptor bool
	DisableMetricInterceptor  bool
	DisableAccessInterceptor  bool
	AccessInterceptorLevel    string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		DialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 10,
		ReadTimeout:            xtime.Duration("1s"),
		SlowThreshold:          xtime.Duration("600ms"),
		OnDialError:            "panic",
		AccessInterceptorLevel: "info",
		Block:                  true,
		KeepAlive: &keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 3,
			PermitWithoutStream: true,
		},
	}
}

// WithDialOption ...
func (config *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	if config.DialOptions == nil {
		config.DialOptions = make([]grpc.DialOption, 0)
	}
	config.DialOptions = append(config.DialOptions, opts...)
	return config
}

// Build ...
func (config *Config) Build() *grpc.ClientConn {
	config.DialOptions = append(config.DialOptions,
		grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor()),
	)
	cc, err := newGRPCClient(config)
	if err != nil {
		jlog.Error(err.Error())
	}
	return cc
}
