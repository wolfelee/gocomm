package redislock

import (
	"context"
	"sync"
	"time"

	"github.com/wolfelee/gocomm/pkg/localstorage"
	"github.com/wolfelee/gocomm/pkg/lock"
)

var localStore localstorage.Store = localstorage.NewLocalStorage(30, false)
var globalLock = sync.Mutex{}

const lockExpired = time.Minute * 10

// getLocalLock 获取本地锁
func getLocalLock(name string) *sync.Mutex {
	globalLock.Lock()
	defer globalLock.Unlock()
	lk := localStore.Get(name)
	if v, ok := lk.(*sync.Mutex); ok {
		localStore.Set(name, lk, lockExpired)
		return v
	}
	ll := &sync.Mutex{}
	localStore.Set(name, ll, lockExpired)
	return ll
}

type redisLockV2 struct {
	disLock   lock.Locker
	localLock *sync.Mutex
}

func NewRedisLockV2(key string, opt ...LockOpt) (lock.Locker, error) {
	disLock, err := NewRedisLock(key, opt...)
	if err != nil {
		return nil, err
	}
	localLock := getLocalLock(key)
	return &redisLockV2{
		disLock:   disLock,
		localLock: localLock,
	}, nil
}

// Lock 加锁
func (l *redisLockV2) Lock() error {
	l.localLock.Lock()
	return l.disLock.Lock()
}

// Unlock 解锁
func (l *redisLockV2) Unlock() error {
	err := l.disLock.Unlock()
	l.localLock.Unlock()
	return err
}

// LockWithCtx 带超时时间的抢锁
//
//	可能返回会有延后, 因为会抢本地锁
func (l *redisLockV2) LockWithCtx(ctx context.Context) error {
	l.localLock.Lock()
	err := l.disLock.LockWithCtx(ctx)
	if err != nil {
		l.localLock.Unlock()
	}
	return err
}
