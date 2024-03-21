package mqueue

import (
	"github.com/streadway/amqp"
	"sync"
)

//Context 消息上下文
type Context struct {
	delivery *amqp.Delivery //mq收到的消息投递
	done     bool           //ack成功后为true
}

//GetHeader 获取消息头
func (c *Context) GetHeader(key string) interface{} {
	if value, ok := c.delivery.Headers[key]; ok {
		return value
	}
	return nil
}

//GetBody 获取body体
func (c *Context) GetBody() []byte {
	return c.delivery.Body
}

//GetMsgId 获取消息id
func (c *Context) GetMsgId() string {
	return c.delivery.MessageId
}

//Ack 消息应答
func (c *Context) Ack() error {
	if c.done { //防止重复应答,重复应答会导致custom停止接收消息
		return nil
	}
	e := c.delivery.Ack(false)
	if e != nil {
		return e
	}
	c.done = true
	return nil
}

//NAck 应答消息
//	request=true  //告知server重新派发到其他消费者
//	request=false //告知服务器删除此消息,或者派发到服务器配置的死信队列
func (c *Context) Nack(requeue bool) error {
	if c.done { //防止重复应答,重复应答会导致custom停止接收消息
		return nil
	}
	e := c.delivery.Nack(false, requeue)
	if e != nil {
		return e
	}
	c.done = true
	return nil
}

func (c *Context) reset() {
	c.delivery = nil
	c.done = false
}

//context poll
var pool = sync.Pool{
	New: func() interface{} {
		return &Context{}
	},
}

func newContext(d *amqp.Delivery) *Context {
	c, _ := pool.Get().(*Context)
	c.reset()
	c.delivery = d
	return c
}

func releaseContext(c *Context) {
	pool.Put(c)
}
