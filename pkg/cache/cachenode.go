package cache

import (
	"errors"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/wolfelee/gocomm/pkg/db"
	"github.com/wolfelee/gocomm/pkg/util/jsync"
	"math/rand"
	"time"
	"xorm.io/xorm"
)

var (
	sc = jsync.NewSharesCall()
)

type DataBase interface {
	Exec(exec func(session *xorm.Session) error, keys ...string) error
	QueryRow(key string, v interface{}, query func(session *xorm.Session, v interface{}) (bool, error)) (bool, error)
	QueryRowNotCache(v interface{}, query func(session *xorm.Session, v interface{}) (bool, error)) (bool, error)
	Take(key string, v interface{}, query func(v interface{}) (bool, error)) (bool, error)
	CloseSession() error
}

func NewDataBase(dbName string, cacheName string, opts ...Options) DataBase {
	c := &cacheNode{
		xorm:              db.Use(dbName).NewSession(),
		sharesCall:        sc,
		expire:            time.Hour * 24 * 7,
		rand:              0.1,
		keyNotExistExpire: time.Minute * 5,
		errNotFound:       errors.New("data not found"),
	}

	if cacheName != "" {
		c.redis = Use(cacheName)
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func NewDataBaseTest(x *xorm.Session, r *redis.Client) DataBase {
	c := &cacheNode{
		redis:             r,
		xorm:              x,
		sharesCall:        sc,
		expire:            time.Hour * 24 * 7,
		rand:              0.1,
		keyNotExistExpire: time.Minute * 5,
		errNotFound:       errors.New("data not found"),
	}

	return c
}

type cacheNode struct {
	// xorm session
	xorm *xorm.Session
	// redis client
	redis *redis.Client
	// 同一时间做同一件事 防止redis击穿
	sharesCall jsync.SharesCall
	// 缓存过期时间
	expire time.Duration
	// 缓存过期时间偏移量 防止雪崩 expire = (expire * (1-rand到1+rand的随机数))
	rand float64
	// db中没有这个数据 设置的短的redis缓存 时间不易太长 防止redis穿透
	keyNotExistExpire time.Duration
	//
	errNotFound error
}

func (c *cacheNode) Exec(exec func(session *xorm.Session) error, keys ...string) error {
	err := exec(c.xorm)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.redis.Del(keys...).Err()

}

// QueryRow 根据ID查询数据
// key redis key ,v xorm模型数据 传入指针
func (c *cacheNode) QueryRow(key string, v interface{}, query func(session *xorm.Session,
	v interface{}) (bool, error)) (bool, error) {
	return c.Take(key, v, func(v interface{}) (bool, error) {
		return query(c.xorm, v)
	})
}

func (c *cacheNode) QueryRowNotCache(v interface{}, query func(session *xorm.Session,
	v interface{}) (bool, error)) (bool, error) {
	return query(c.xorm, v)
}

func (c *cacheNode) Take(key string, v interface{}, query func(v interface{}) (bool, error)) (bool, error) {
	return c.doTake(key, v, query, func(v interface{}) error {
		return c.set(key, v)
	})
}

func (c *cacheNode) doTake(key string, v interface{}, query func(v interface{}) (bool, error),
	cacheVal func(v interface{}) error) (bool, error) {
	val, fresh, err := c.sharesCall.DoEx(key, func() (interface{}, error) {
		if err := c.doGetCache(key, v); err != nil {
			if err == errPlaceholder {
				return nil, c.errNotFound
			} else if err != c.errNotFound {
				return nil, err
			}

			if has, err := query(v); err == c.errNotFound || !has {
				if err := c.setCacheWithNotFound(key); err != nil {
					// 报警
					cacheAlarm()
				}
				return nil, c.errNotFound
			} else if err != nil {
				// 报警
				cacheAlarm()
				return nil, err
			}

			if err := cacheVal(v); err != nil {
				cacheAlarm()
				// 报警
			}
		}
		return jsoniter.MarshalToString(v)
	})

	if err == c.errNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if fresh {
		return true, nil
	}

	// 采样 查询个数

	err = jsoniter.UnmarshalFromString(val.(string), v)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *cacheNode) CloseSession() error {
	return c.xorm.Close()
}

// 设置过期时间偏移量 防止雪崩
func (c *cacheNode) getCacheExpire() time.Duration {
	rand.Seed(time.Now().UnixNano())
	max, min := 1+c.rand, 1-c.rand
	factor := min + rand.Float64()*(max-min)

	return time.Duration(float64(c.expire) * factor)
}
