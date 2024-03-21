package redislock

import "time"

type defaultTime struct {
	DefaultTime time.Duration
}

func (d defaultTime) apply(r *RedisLock) {
	r.SetDefaultTime(d.DefaultTime)
}

//WithDefaultTime 设置默认锁定时间
func WithDefaultTime(t time.Duration) LockOpt {
	return defaultTime{DefaultTime: t}
}