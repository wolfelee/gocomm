package cache

import (
	"github.com/wolfelee/gocomm/pkg/util/jsync"
	"time"
)

type Options func(c *cacheNode)

func WithCacheSharesCall(sc jsync.SharesCall) Options {
	return func(c *cacheNode) {
		c.sharesCall = sc
	}
}

func WithExpire(expire time.Duration) Options {
	return func(c *cacheNode) {
		c.expire = expire
	}
}

func WithRand(r float64) Options {
	if r <= 0 || r >= 1 {
		panic("cache rand should be between 0 and 1")
	}
	return func(c *cacheNode) {
		c.rand = r
	}
}

func WithKeyNotExistExpire(keyNotExistExpire time.Duration) Options {
	return func(c *cacheNode) {
		c.keyNotExistExpire = keyNotExistExpire
	}
}

func WithErrNotFound(err error) Options {
	return func(c *cacheNode) {
		c.errNotFound = err
	}
}
