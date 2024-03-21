package redislock

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

//解锁lua脚本
const unlockLua = `
if redis.call("GET", KEYS[1])==ARGV[1] then
	redis.call("DEL", KEYS[1])
	return true
else	
	return false
end
`

//续期使用的lua脚本
const renewalLua = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	redis.call("EXPIRE", KEYS[1], ARGV[2])
	return true
else 
	return false
end
`

type RedisLock struct {
	redisClient  *redis.Client //redis 客户端
	key          string        //lock key
	randomValue  string        //随机数值
	renewalTime  time.Duration //续期时间
	defaultTime  time.Duration //默认上锁时间
	intervalTime time.Duration //重新获取锁间隔
	closeChan    chan struct{} //解锁后立即通知续期线程退出
}

func (r *RedisLock) SetDefaultTime(duration time.Duration) {
	r.defaultTime = duration
	r.renewalTime = duration / 3 //重新续期时间为默认时长的3分之一
}

func (r *RedisLock) Lock() error {
	return r.LockWithCtx(context.Background())
}

func (r *RedisLock) Unlock() error {
	cmd := r.redisClient.Eval(unlockLua, []string{r.key}, r.randomValue)
	b, err := cmd.Bool()
	log.Println("unlock", b, err)
	close(r.closeChan)
	return err
}

func (r *RedisLock) LockWithCtx(ctx context.Context) error {
	timer := time.NewTicker(r.intervalTime)
	defer timer.Stop()
	for {
		ret := r.redisClient.SetNX(r.key, r.randomValue, r.defaultTime)
		isSet, err := ret.Result()
		if err != nil { //异常
			return err
		}
		if isSet { //设置成功
			go func() { //续期线程
				defer func() {
					if err := recover(); err != nil {
						log.Println("lock recover", err)
					}
				}()
				timer := time.NewTicker(r.renewalTime)
				defer timer.Stop()
				for {
					select {
					case <-r.closeChan:
						log.Println("unlock receive from closeChan:", r.key)
						return
					case <-timer.C:
					}
					cmd := r.redisClient.Eval(renewalLua, []string{r.key}, r.randomValue, int(r.defaultTime/time.Second))
					ret, err := cmd.Bool()
					log.Println("renewal:", r.key)
					if err != nil && err != redis.Nil {
						log.Println("lock renewal error:", r.key, ret, err)
						break
					}
					if !ret { //说明随机值已经不相等,或者key被删除
						log.Println("lock value not equal:", r.key)
						break
					}
				}
			}()
			return nil
		}
		select {
		case <-ctx.Done():
			//超时
			return ctx.Err()
		case <-timer.C:
		}
	}
}
