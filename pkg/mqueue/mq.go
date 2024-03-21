package mqueue

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/mqueue/aliutils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"
)

const (
	statusReadyForReconnect int32 = iota
	statusReconnecting
)

const (
	statusStopped int32 = iota
	statusRunning
)

var (
	mqCfg   *MQConfig
	MqGroup map[string]MQ
	AppName string
)

type (
	MQ interface {
		GetConsumer(name string) (Consumer, error)
		SetConsumerHandler(name string, handler ConsumerHandler) error
		GetProducer(name string) (Producer, error)
		Error() <-chan error
		Close()
	}
)

type mq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *Config
	appName string //应用名称

	errorChannel         chan error
	internalErrorChannel chan error
	consumers            *consumersRegistry
	producers            *producersRegistry
	reconnectStatus      int32
}

func initMqGroup() (map[string]MQ, error) {
	var groups = make(map[string]MQ)
	if mqCfg == nil {
		return nil, errors.New("mq config setting error, config is nil")
	}
	for g, e := range mqCfg.Queues {
		group, err := New(e)
		if err != nil {
			return nil, errors.New("mq config setting error" + err.Error())
		}
		groups[g] = group
		jlog.Info(fmt.Sprintf("%s MqGroup Opened", g))
	}
	return groups, nil
}

func Init(mqCfgFile string, appname string) error {
	AppName = appname
	buf, err := ioutil.ReadFile(mqCfgFile)
	if err != nil {
		jlog.Error(mqCfgFile + "file read error")
		return err
	}
	err = yaml.Unmarshal(buf, &mqCfg)
	if err != nil {
		jlog.Error(mqCfgFile + "file unmarshal error" + err.Error())
		return err
	}
	MqGroup, err = initMqGroup()
	if err != nil {
		jlog.Error(err.Error())
		return err
	}
	return nil
}

func Use(mqName string) MQ {
	if MqGroup == nil {
		var err error
		MqGroup, err = initMqGroup()
		if err != nil {
			jlog.Error(err.Error())
		}
	}
	if g, ok := MqGroup[mqName]; ok {
		return g
	} else {
		jlog.Error(mqName + " - mq does not exist.")
	}
	return nil
}

func Close() {
	for n, mq := range MqGroup {
		mq.Close()
		jlog.Info(fmt.Sprintf("%s mq Closed", n))
	}
}

func New(config *Config) (MQ, error) {
	config.normalize()
	mq := &mq{
		config:               config,
		errorChannel:         make(chan error),
		internalErrorChannel: make(chan error),
		consumers:            newConsumersRegistry(len(config.Consumers)),
		producers:            newProducersRegistry(len(config.Producers)),
		appName:              AppName,
	}
	if err := mq.connect(); err != nil {
		return nil, err
	}

	go mq.errorHandler()

	return mq, mq.initialSetup()
}

func (mq *mq) GetConsumer(name string) (consumer Consumer, err error) {
	consumer, ok := mq.consumers.Get(name)
	if !ok {
		err = fmt.Errorf("consumer '%s' is not registered. Check your configuration", name)
	}

	return
}

func (mq *mq) SetConsumerHandler(name string, handler ConsumerHandler) error {
	consumer, err := mq.GetConsumer(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	consumer.Consume(handler)

	return nil
}

func (mq *mq) GetProducer(name string) (producer Producer, err error) {
	producer, ok := mq.producers.Get(name)

	if !ok {
		err = fmt.Errorf("producer '%s' is not registered. Check your configuration", name)
	}
	return
}

func (mq *mq) Error() <-chan error {
	return mq.errorChannel
}

func (mq *mq) Close() {
	mq.stopProducersAndConsumers()

	if mq.channel != nil {
		mq.channel.Close()
	}

	if mq.conn != nil {
		mq.conn.Close()
	}
}

func (mq *mq) connect() error {
	userName := mq.config.Ak
	password := mq.config.Sk
	if mq.config.Ali {
		//阿里账号
		userName = aliutils.GetUserName(userName, mq.config.AliInstanceId)
		password = aliutils.GetPassword(password)
	}

	var dsn bytes.Buffer
	dsn.WriteString("amqp://")
	dsn.WriteString(userName)
	dsn.WriteString(":")
	dsn.WriteString(password)
	dsn.WriteString("@")
	dsn.WriteString(mq.config.DSN)
	url := dsn.String()

	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": mq.appName,
		},
	}

	connection, err := amqp.DialConfig(url, cfg)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return err
	}
	mq.conn = connection
	mq.channel = channel

	go mq.handleCloseEvent()

	return nil
}

func (mq *mq) handleCloseEvent() {
	err := <-mq.conn.NotifyClose(make(chan *amqp.Error))
	if err != nil {
		mq.internalErrorChannel <- err
	}
}

func (mq *mq) errorHandler() {
	for err := range mq.internalErrorChannel {
		select {
		case mq.errorChannel <- err:
		default:
		}
		mq.processError(err)
	}
}

func (mq *mq) processError(err interface{}) {
	switch err.(type) {
	case *net.OpError:
		jlog.Error("mq error, will reconnect", jlog.Any("error", err))
		go mq.reconnect()
	case *amqp.Error:
		//rmqErr, _ := err.(*amqp.Error)
		//if rmqErr.Server == false
		{
			jlog.Error("mq error, will reconnect--", jlog.Any("error", err))
			go mq.reconnect()
		}
	default:
	}
}

func (mq *mq) initialSetup() error {
	if err := mq.setupExchanges(); err != nil {
		return err
	}

	if err := mq.setupQueues(); err != nil {
		return err
	}

	if err := mq.setupProducers(); err != nil {
		return err
	}

	return mq.setupConsumers()
}

func (mq *mq) setupAfterReconnect() error {
	if err := mq.setupExchanges(); err != nil {
		return err
	}

	if err := mq.setupQueues(); err != nil {
		return err
	}

	mq.producers.GoEach(func(producer *producer) {
		if err := mq.reconnectProducer(producer); err != nil {
			mq.internalErrorChannel <- err
		}
	})

	mq.consumers.GoEach(func(consumer *consumer) {
		if err := mq.reconnectConsumer(consumer); err != nil {
			mq.internalErrorChannel <- err
		}
	})

	return nil
}

func (mq *mq) setupExchanges() error {
	for _, config := range mq.config.Exchanges {
		if err := mq.declareExchange(config); err != nil {
			return err
		}
	}

	return nil
}

func (mq *mq) declareExchange(config ExchangeConfig) error {
	return mq.channel.ExchangeDeclare(config.Name, config.Type, true, false, false, false, nil)
}

func (mq *mq) setupQueues() error {
	for _, config := range mq.config.Queues {
		if err := mq.declareQueue(config); err != nil {
			return err
		}
	}

	return nil
}

func (mq *mq) declareQueue(config QueueConfig) error {
	if _, err := mq.channel.QueueDeclare(config.Name, true, false, false, false, nil); err != nil {
		return err
	}
	return mq.channel.QueueBind(config.Name, config.RoutingKey, config.Exchange, false, nil)
}

func (mq *mq) setupProducers() error {
	for _, config := range mq.config.Producers {
		//检查exchange
		if err := mq.channel.ExchangeDeclarePassive(config.Exchange, "topic", true, false, false, false, nil); err != nil {
			return err
		}

		if err := mq.registerProducer(config, mq.config.Ali); err != nil {
			return err
		}
	}

	return nil
}

func (mq *mq) registerProducer(config ProducerConfig, isAli bool) error {
	if _, ok := mq.producers.Get(config.Name); ok {
		return fmt.Errorf(`producer with name "%s" is already registered`, config.Name)
	}

	channel, err := mq.conn.Channel()
	if err != nil {
		return err
	}

	producer := newProducer(channel, mq.internalErrorChannel, config, isAli)

	go producer.worker()
	mq.producers.Set(config.Name, producer)

	return nil
}

func (mq *mq) reconnectProducer(producer *producer) error {
	channel, err := mq.conn.Channel()
	if err != nil {
		return err
	}

	producer.setChannel(channel)
	go producer.worker()

	return nil
}

func (mq *mq) setupConsumers() error {
	for _, config := range mq.config.Consumers {
		if err := mq.registerConsumer(config); err != nil {
			return err
		}
	}
	return nil
}

func (mq *mq) registerConsumer(config ConsumerConfig) error {
	if _, ok := mq.consumers.Get(config.Name); ok {
		return fmt.Errorf(`consumer with name "%s" is already registered`, config.Name)
	}
	if config.Workers == 0 {
		config.Workers = 1
	}

	consumer := newConsumer(config, mq)
	consumer.prefetchCount = 1
	consumer.prefetchSize = 0

	for i := 0; i < config.Workers; i++ {
		worker := newWorker(mq.internalErrorChannel)

		consumer.workers[i] = worker
	}

	mq.consumers.Set(config.Name, consumer)

	return nil
}

func (mq *mq) reconnectConsumer(consumer *consumer) error {
	for _, worker := range consumer.workers {
		if consumer.handler == nil { //如果配置文件写了消费者,但是代码里没有设置消费者则跳过
			continue
		}
		if err := mq.initializeConsumersWorker(consumer, worker); err != nil {
			return err
		}

		go worker.Run(consumer.handler)
	}

	return nil
}

func (mq *mq) initializeConsumersWorker(consumer *consumer, worker *worker) error {
	channel, err := mq.conn.Channel()
	if err != nil {
		return err
	}

	if err := channel.Qos(consumer.prefetchCount, consumer.prefetchSize, false); err != nil {
		return err
	}

	deliveries, err := channel.Consume(consumer.queue, consumer.name, false, false, false, false, nil)
	if err != nil {
		return err
	}

	worker.setChannel(channel)
	worker.deliveries = deliveries

	return nil
}

func (mq *mq) reconnect() {
	notBusy := atomic.CompareAndSwapInt32(&mq.reconnectStatus, statusReadyForReconnect, statusReconnecting)
	if !notBusy {
		return
	}

	defer func() {
		atomic.StoreInt32(&mq.reconnectStatus, statusReadyForReconnect)
	}()
	time.Sleep(mq.config.ReconnectDelay)

	mq.stopProducersAndConsumers()
	if err := mq.connect(); err != nil {
		mq.internalErrorChannel <- err
		return
	}

	if err := mq.setupAfterReconnect(); err != nil {
		mq.internalErrorChannel <- err
	}
}

func (mq *mq) stopProducersAndConsumers() {
	mq.producers.GoEach(func(producer *producer) {
		producer.Stop()
	})

	mq.consumers.GoEach(func(consumer *consumer) {
		consumer.Stop()
	})
}

type workerStatus struct {
	value int32
}

func (status *workerStatus) markAsRunning() {
	atomic.StoreInt32(&status.value, statusRunning)
}

func (status *workerStatus) markAsStoppedIfCan() bool {
	return atomic.CompareAndSwapInt32(&status.value, statusRunning, statusStopped)
}
