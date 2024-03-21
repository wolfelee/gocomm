package mqueue

import (
	"github.com/streadway/amqp"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/util/juuid"
	"sync"
	"time"
)

type Producer interface {
	Produce(data PublicData)
	SetRoutingKey(key string)
	WaitEmpty()
}

type PublicData struct {
	Message  []byte
	RouteKey string
	Delay    int
}
type producer struct {
	sync.Mutex
	workerStatus

	channel         *amqp.Channel
	errorChannel    chan<- error
	exchange        string
	mandatory       bool
	immediate       bool
	options         Options
	publishChannel  chan PublicData
	routingKey      string
	shutdownChannel chan struct{}
	ali             bool //标记是否是阿里云mq队列，用于区分延迟队列的处理：阿里mq直接支持延迟队列，自建mq需要通过死信方式支持
	emptyNotify     *EmptyWatcher
}

func newProducer(channel *amqp.Channel, errorChannel chan<- error, config ProducerConfig, isAli bool) *producer {
	return &producer{
		channel:         channel,
		errorChannel:    errorChannel,
		exchange:        config.Exchange,
		options:         config.Options,
		mandatory:       config.Mandatory,
		immediate:       config.Immediate,
		publishChannel:  make(chan PublicData, config.BufferSize),
		routingKey:      config.RoutingKey,
		shutdownChannel: make(chan struct{}),
		ali:             isAli,
		emptyNotify:     NewEmptyWatcher(false),
	}
}

func (producer *producer) worker() {
	producer.markAsRunning()

	for {
		select {
		case pubData := <-producer.publishChannel:
			err := producer.produce(pubData)
			if err != nil {
				producer.errorChannel <- err
			}

			if len(producer.publishChannel) == 0 {
				producer.emptyNotify.Set(true)
			}
		case <-producer.shutdownChannel:
			producer.closeChannel()

			return
		}
	}
}

// WaitEmpty 阻塞,直到发送队列为空
func (producer *producer) WaitEmpty() {
	//如果为空直接返回
	producer.emptyNotify.Wait()
}

func (producer *producer) setChannel(channel *amqp.Channel) {
	producer.Lock()
	producer.channel = channel
	producer.Unlock()
}

func (producer *producer) SetRoutingKey(key string) {
	producer.Lock()
	producer.routingKey = key
	producer.Unlock()
}

func (producer *producer) closeChannel() {
	producer.Lock()
	if err := producer.channel.Close(); err != nil {
		producer.errorChannel <- err
	}
	producer.Unlock()
}

func (producer *producer) Produce(data PublicData) {

	tm := time.NewTimer(1 * time.Second)
	select {
	case producer.publishChannel <- data:
		producer.emptyNotify.Set(false)
	case <-tm.C:
		jlog.Error("mq service error")
	}
}

func (producer *producer) produce(pubData PublicData) error {
	producer.Lock()
	defer producer.Unlock()

	var msg = amqp.Publishing{}

	msg.ContentType = "application/json"
	msg.DeliveryMode = 2 //2:持久化, 0,1:临时

	msg.Body = pubData.Message
	msg.MessageId = juuid.ShortUUID()
	routingKey := producer.routingKey
	if len(pubData.RouteKey) > 0 {
		routingKey = pubData.RouteKey
	}

	if pubData.Delay > 0 {
		delay := "delay"
		if !producer.ali {
			delay = "x-delay"
		}
		msg.Headers = amqp.Table{delay: pubData.Delay * 1000}
	}
	var fields = make([]jlog.Field, 0, 4)
	fields = append(fields,
		jlog.String("mqID", msg.MessageId),
		jlog.String("exchange", producer.exchange),
		jlog.String("routKey", routingKey),
	)
	err := producer.channel.Publish(producer.exchange, routingKey, producer.mandatory, producer.immediate, msg)
	if err != nil {
		fields = append(fields, jlog.String("error", err.Error()), jlog.String("data", string(msg.Body)))
		jlog.Error(err.Error(), fields...)
	} else {
		jlog.Info(string(msg.Body), fields...)
	}
	return err
}

func (producer *producer) Stop() {
	if producer.markAsStoppedIfCan() {
		producer.shutdownChannel <- struct{}{}
	}
}
