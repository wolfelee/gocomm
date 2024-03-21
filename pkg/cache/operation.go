package cache

import (
	"errors"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

const notFoundPlaceholder = "*"

// indicates there is no such value associate with the key
var errPlaceholder = errors.New("placeholder")

func (c *cacheNode) set(key string, v interface{}) error {
	str, err := jsoniter.MarshalToString(v)
	if err != nil {
		return err
	}
	return c.redis.Set(key, str, c.getCacheExpire()).Err()
}

func (c *cacheNode) doGetCache(key string, v interface{}) error {
	// 添加监控命中率等
	data, err := c.redis.Get(key).Result()
	if err == redis.Nil {
		return c.errNotFound
	}
	if err != nil {
		return err
	}

	if data == notFoundPlaceholder {
		return errPlaceholder
	}

	return jsoniter.UnmarshalFromString(data, v)
}

func (c cacheNode) setCacheWithNotFound(key string) error {
	return c.redis.Set(key, notFoundPlaceholder, c.keyNotExistExpire).Err()
}

// redis报警
func cacheAlarm() {

}
