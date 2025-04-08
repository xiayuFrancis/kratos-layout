package data

import (
	"kratosdemo/internal/conf"
	"kratosdemo/internal/pkg/kafka"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// KafkaData 消息数据结构体
type KafkaData struct {
	client *kafka.Client
	log    *log.Helper
}

// NewKafkaData 创建 Kafka 数据结构体
func NewKafkaData(c *conf.Data, logger log.Logger) (*KafkaData, func(), error) {
	log := log.NewHelper(logger)
	
	// 如果没有配置 Kafka 或者没有配置 Brokers，则跳过 Kafka 客户端的创建
	if c.Kafka == nil || len(c.Kafka.Brokers) == 0 {
		log.Warn("Kafka 配置不存在或 brokers 为空，跳过 Kafka 客户端的创建")
		return &KafkaData{
			client: nil,
			log:    log,
		}, func() {}, nil
	}
	
	// 创建 Kafka 客户端
	client, err := kafka.NewClient(c.Kafka, logger)
	if err != nil {
		log.Errorf("创建 Kafka 客户端失败: %v", err)
		return nil, nil, err
	}
	
	kd := &KafkaData{
		client: client,
		log:    log,
	}
	
	// 清理函数
	cleanup := func() {
		log.Info("关闭 Kafka 客户端资源")
		if kd.client != nil {
			if err := kd.client.Close(); err != nil {
				log.Error(err)
			}
		}
	}
	
	return kd, cleanup, nil
}

// Client 获取 Kafka 客户端
func (kd *KafkaData) Client() *kafka.Client {
	return kd.client
}

// NewKafkaClient 提供 Kafka 客户端实例
func NewKafkaClient(kd *KafkaData) *kafka.Client {
	return kd.client
}

// KafkaProviderSet 是 Kafka 提供者集合
var KafkaProviderSet = wire.NewSet(NewKafkaData, NewKafkaClient)
