package grpc

import (
	"context"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"google.golang.org/grpc"
	"time"
)

func newGRPCClient(config *Config) (*grpc.ClientConn, error) {
	var ctx = context.Background()
	var dialOptions = config.DialOptions
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName))
	var cc *grpc.ClientConn
	var err error
	cc, err = grpc.DialContext(ctx, config.Address, dialOptions...)

	if err == context.DeadlineExceeded {
		for i := 0; i < 5; i++ {
			jlog.Errorf("dial grpc server address:%s err:%s, reproducing connection count:%d", config.Address, err, i+1)
			cc, err = func() (*grpc.ClientConn, error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
				defer cancel()
				return grpc.DialContext(ctx, config.Address, dialOptions...)
			}()
			if err != nil && err != context.DeadlineExceeded {
				break
			}
		}
	}

	if err != nil {
		jlog.Errorf("dial grpc server error -> name:%s address:%s err:%s,", config.Name, config.Address, err)
		return nil, err
	}

	jlog.Infof("start grpc client -> name:%s  address:%s", config.Name, config.Address)
	return cc, nil
}
