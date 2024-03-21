package localstorage

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Store interface {
	Get(name string) interface{}                              //获取
	Set(name string, value interface{}, expire time.Duration) // 添加,更新, 时间<=0 还可以删除
	Del(name string) bool                                     //删除节点
}

type node struct {
	expire time.Time   //过期时间
	value  interface{} //数据节点
	index  int         //在堆上的索引
	name   string      //名称
}

func (n *node) String() string {
	return fmt.Sprintf("(%d, %v, %d)", n.index, n.value, n.expire.Unix())
}

//LocalStorage 本地存储
type LocalStorage struct {
	l          *sync.Mutex
	store      map[string]*node //数据管理
	timeHeap   *MinHeap         //淘汰管理
	lastAccess time.Time        //最后访问时间,用来识别清理时间,闲时清理过期数据
}

//Get 获取
func (l *LocalStorage) Get(name string) interface{} {
	l.l.Lock()
	defer l.l.Unlock()

	l.lastAccess = time.Now()
	if v, ok := l.store[name]; ok {
		if v.expire.Before(time.Now()) {
			delete(l.store, name)
			l.timeHeap.DeleteNode(v.index)
			log.Println("节点过期:", name, v.expire.Unix())
			return nil
		}
		return v.value
	}

	return nil
}

//Set 设置
func (l *LocalStorage) Set(name string, value interface{}, expire time.Duration) {
	// if expire <= 0 {
	// 	return
	// }
	l.l.Lock()
	defer l.l.Unlock()
	l.lastAccess = time.Now()
	if v, ok := l.store[name]; ok { //修改已有的值
		v.value = value
		v.expire = time.Now().Add(expire)
		l.timeHeap.SetNode(v)
		return
	}
	//创建新值
	n := &node{
		expire: time.Now().Add(expire),
		value:  value,
		name:   name,
	}
	l.timeHeap.AddNode(n)
	l.store[name] = n
}

func (l *LocalStorage) Del(name string) bool {
	l.l.Lock()
	defer l.l.Unlock()
	v, ok := l.store[name]
	if !ok {
		return false
	}
	l.timeHeap.DeleteNode(v.index)
	return true
}

//expired 淘汰过期节点
//	force强制回收, 如果非强制,会闲时回收
func (l *LocalStorage) expired(force bool) {
	l.l.Lock()
	defer l.l.Unlock()

	//非强制回收,如果已经3秒没有访问过了,则判断为空闲,开始回收
	if !force && time.Now().Unix()-l.lastAccess.Unix() < 3 {
		log.Println("还未处于闲置状态,暂不回收")
		return
	}

	for !l.timeHeap.Empty() {
		min := l.timeHeap.Peek()
		if min.expire.After(time.Now()) { //最小值还未过期,则没有要回收的数据
			break
		}
		l.timeHeap.Pop()          //从堆上删除
		delete(l.store, min.name) //从map上删除
		log.Println("回收节点:", min.name, min.expire.Unix())
	}
}

func (l *LocalStorage) run(forceExpire bool) {
	tick := time.Tick(time.Second)
	for range tick {
		l.expired(forceExpire)
	}
}

//NewLocalStorage 新建本地存储
//	initCapcity初始化容量
//	forceExpire是否强制回收
func NewLocalStorage(initCapcity int, forceExpire bool) Store {
	store := &LocalStorage{
		l:          &sync.Mutex{},
		store:      make(map[string]*node, initCapcity),
		timeHeap:   NewMinHeap(initCapcity),
		lastAccess: time.Now(),
	}
	go func() {
		defer func() {
			log.Println("error:过期线程已退出")
		}()
		store.run(forceExpire)
	}()
	return store
}
