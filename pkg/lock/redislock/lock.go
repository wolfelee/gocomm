package redislock

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/wolfelee/gocomm/pkg/lock"
)

var gCli *redis.Client

var lockPrefix string = "lock"

// SetRedisClient 设置redis客户端
func SetRedisClient(cli *redis.Client) {
	gCli = cli
}

// SetGlobalPrefix 设置redis锁key前缀
func SetGlobalPrefix(prefix string) {
	lockPrefix = prefix
}

func NewRedisLock(key string, opt ...LockOpt) (lock.Locker, error) {
	if gCli == nil {
		return nil, errors.New("请先设置redis客户端,使用:SetRedisClient")
	}
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	dT := 30 * time.Second //默认30秒
	rl := &RedisLock{
		redisClient:  gCli,
		key:          fmt.Sprintf("%s:%s", lockPrefix, key),
		randomValue:  u.String(),
		renewalTime:  dT / 3,                 //锁续期时间
		defaultTime:  dT,                     //锁时间
		intervalTime: 100 * time.Millisecond, //重新获取时间
		closeChan:    make(chan struct{}),
	}

	for _, op := range opt {
		op.apply(rl)
	}
	return rl, nil
}
