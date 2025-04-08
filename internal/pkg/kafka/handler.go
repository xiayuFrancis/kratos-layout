package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

// MessageHandler 基础消息处理器
type MessageHandler struct {
	log *log.Helper
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(logger log.Logger) *MessageHandler {
	return &MessageHandler{
		log: log.NewHelper(logger),
	}
}

// Setup 实现 ConsumerHandler 接口
func (h *MessageHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.log.Infof("消费者会话开始: %s", session.MemberID())
	return nil
}

// Cleanup 实现 ConsumerHandler 接口
func (h *MessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.log.Infof("消费者会话结束: %s", session.MemberID())
	return nil
}

// ConsumeClaim 实现 ConsumerHandler 接口
func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.log.Infof("收到消息: topic=%s, partition=%d, offset=%d, key=%s, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		
		// 标记消息为已处理
		session.MarkMessage(msg, "")
	}
	return nil
}

// JSONMessageHandler JSON 消息处理器
type JSONMessageHandler struct {
	MessageHandler
	processor func(map[string]interface{}) error
}

// NewJSONMessageHandler 创建 JSON 消息处理器
func NewJSONMessageHandler(logger log.Logger, processor func(map[string]interface{}) error) *JSONMessageHandler {
	return &JSONMessageHandler{
		MessageHandler: *NewMessageHandler(logger),
		processor:      processor,
	}
}

// ConsumeClaim 处理 JSON 消息
func (h *JSONMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.log.Infof("收到 JSON 消息: topic=%s, partition=%d, offset=%d",
			msg.Topic, msg.Partition, msg.Offset)
		
		// 解析 JSON 消息
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			h.log.Errorf("解析 JSON 消息失败: %v", err)
			session.MarkMessage(msg, "")
			continue
		}
		
		// 处理消息
		if err := h.processor(data); err != nil {
			h.log.Errorf("处理消息失败: %v", err)
		}
		
		// 标记消息为已处理
		session.MarkMessage(msg, "")
	}
	return nil
}

// EventHandler 事件处理器
type EventHandler struct {
	MessageHandler
	eventHandlers map[string]func(map[string]interface{}) error
}

// NewEventHandler 创建事件处理器
func NewEventHandler(logger log.Logger) *EventHandler {
	return &EventHandler{
		MessageHandler: *NewMessageHandler(logger),
		eventHandlers:  make(map[string]func(map[string]interface{}) error),
	}
}

// RegisterEventHandler 注册事件处理函数
func (h *EventHandler) RegisterEventHandler(eventType string, handler func(map[string]interface{}) error) {
	h.eventHandlers[eventType] = handler
}

// ConsumeClaim 处理事件消息
func (h *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.log.Infof("收到事件消息: topic=%s, partition=%d, offset=%d",
			msg.Topic, msg.Partition, msg.Offset)
		
		// 解析 JSON 消息
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			h.log.Errorf("解析 JSON 消息失败: %v", err)
			session.MarkMessage(msg, "")
			continue
		}
		
		// 获取事件类型
		eventType, ok := data["type"].(string)
		if !ok {
			h.log.Errorf("消息缺少事件类型字段")
			session.MarkMessage(msg, "")
			continue
		}
		
		// 查找对应的处理函数
		handler, ok := h.eventHandlers[eventType]
		if !ok {
			h.log.Warnf("未找到事件类型的处理函数: %s", eventType)
			session.MarkMessage(msg, "")
			continue
		}
		
		// 处理事件
		if err := handler(data); err != nil {
			h.log.Errorf("处理事件失败: %v", err)
		}
		
		// 标记消息为已处理
		session.MarkMessage(msg, "")
	}
	return nil
}
