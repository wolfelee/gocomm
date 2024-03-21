package redislock

type LockOpt interface {
	apply(*RedisLock)
}
