package kafka

import (
	"context"
	"crypto/tls"
	"kratosdemo/internal/conf"
	"sync"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

// Client Kafka 客户端封装
type Client struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	admin    sarama.ClusterAdmin
	config   *conf.Data_Kafka
	log      *log.Helper
	topics   []string
	handlers map[string]ConsumerHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// ConsumerHandler 消费者处理接口
type ConsumerHandler interface {
	// Setup 在消费者会话开始时调用
	Setup(sarama.ConsumerGroupSession) error
	// Cleanup 在消费者会话结束时调用
	Cleanup(sarama.ConsumerGroupSession) error
	// ConsumeClaim 处理消息
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

// NewClient 创建 Kafka 客户端
func NewClient(c *conf.Data_Kafka, logger log.Logger) (*Client, error) {
	log := log.NewHelper(logger)
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // 使用 Kafka 2.8.0 版本
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 如果启用了 TLS
	if c.EnableTls {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// 如果配置了 SASL 认证
	if c.Username != "" && c.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.Username
		config.Net.SASL.Password = c.Password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	// 设置客户端 ID
	if c.ClientId != "" {
		config.ClientID = c.ClientId
	}

	// 创建生产者
	producer, err := sarama.NewSyncProducer(c.Brokers, config)
	if err != nil {
		log.Errorf("failed to create kafka producer: %v", err)
		return nil, err
	}

	// 创建消费者组
	consumer, err := sarama.NewConsumerGroup(c.Brokers, c.GroupId, config)
	if err != nil {
		log.Errorf("failed to create kafka consumer group: %v", err)
		producer.Close()
		return nil, err
	}

	// 创建管理客户端
	admin, err := sarama.NewClusterAdmin(c.Brokers, config)
	if err != nil {
		log.Errorf("failed to create kafka admin client: %v", err)
		producer.Close()
		consumer.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		producer: producer,
		consumer: consumer,
		admin:    admin,
		config:   c,
		log:      log,
		handlers: make(map[string]ConsumerHandler),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Close 关闭 Kafka 客户端
func (c *Client) Close() error {
	c.cancel()
	
	if err := c.producer.Close(); err != nil {
		c.log.Errorf("failed to close kafka producer: %v", err)
	}
	
	if err := c.consumer.Close(); err != nil {
		c.log.Errorf("failed to close kafka consumer: %v", err)
	}
	
	if err := c.admin.Close(); err != nil {
		c.log.Errorf("failed to close kafka admin: %v", err)
	}
	
	return nil
}

// SendMessage 发送消息到指定主题
func (c *Client) SendMessage(topic string, key, value []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	
	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}
	
	return c.producer.SendMessage(msg)
}

// RegisterHandler 注册消费者处理器
func (c *Client) RegisterHandler(topic string, handler ConsumerHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.handlers[topic] = handler
	
	// 添加到主题列表
	found := false
	for _, t := range c.topics {
		if t == topic {
			found = true
			break
		}
	}
	
	if !found {
		c.topics = append(c.topics, topic)
	}
}

// StartConsuming 开始消费消息
func (c *Client) StartConsuming() error {
	if len(c.topics) == 0 {
		return nil
	}
	
	// 创建一个消费者处理器
	handler := &consumerGroupHandler{
		client: c,
	}
	
	// 启动消费循环
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumer.Consume(c.ctx, c.topics, handler); err != nil {
					c.log.Errorf("error from consumer: %v", err)
				}
			}
		}
	}()
	
	return nil
}

// CreateTopic 创建主题
func (c *Client) CreateTopic(topic string, numPartitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}
	
	return c.admin.CreateTopic(topic, topicDetail, false)
}

// DeleteTopic 删除主题
func (c *Client) DeleteTopic(topic string) error {
	return c.admin.DeleteTopic(topic)
}

// ListTopics 列出所有主题
func (c *Client) ListTopics() (map[string]sarama.TopicDetail, error) {
	return c.admin.ListTopics()
}

// consumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type consumerGroupHandler struct {
	client *Client
}

// Setup 在消费者会话开始时调用
func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.client.log.Info("consumer group session setup")
	
	// 调用每个处理器的 Setup 方法
	h.client.mu.RLock()
	defer h.client.mu.RUnlock()
	
	for _, handler := range h.client.handlers {
		if err := handler.Setup(session); err != nil {
			return err
		}
	}
	
	return nil
}

// Cleanup 在消费者会话结束时调用
func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.client.log.Info("consumer group session cleanup")
	
	// 调用每个处理器的 Cleanup 方法
	h.client.mu.RLock()
	defer h.client.mu.RUnlock()
	
	for _, handler := range h.client.handlers {
		if err := handler.Cleanup(session); err != nil {
			return err
		}
	}
	
	return nil
}

// ConsumeClaim 处理消息
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	
	h.client.mu.RLock()
	handler, ok := h.client.handlers[topic]
	h.client.mu.RUnlock()
	
	if !ok {
		h.client.log.Warnf("no handler registered for topic: %s", topic)
		// 标记消息为已处理，但不实际处理
		for msg := range claim.Messages() {
			session.MarkMessage(msg, "")
		}
		return nil
	}
	
	// 调用对应主题的处理器
	return handler.ConsumeClaim(session, claim)
}
