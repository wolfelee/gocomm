package lock

import "context"

type Locker interface {
	Lock() error                           //不返回错误就上锁, 直到获取锁才返回
	Unlock() error                         //解锁
	LockWithCtx(ctx context.Context) error //直到超时或者获取锁
}
