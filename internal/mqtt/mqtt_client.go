// Package mqtt 提供基础 MQTT 客户端功能，便于健康数据采集与事件驱动扩展。
package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ClientConfig 配置结构体，便于后续扩展
type ClientConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// MQTTClient 基础客户端结构体
type MQTTClient struct {
	client mqtt.Client
	config ClientConfig
}

// NewMQTTClient 创建并初始化 MQTTClient
func NewMQTTClient(cfg ClientConfig) *MQTTClient {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetAutoReconnect(true).
		SetConnectTimeout(5 * time.Second)

	client := mqtt.NewClient(opts)
	return &MQTTClient{
		client: client,
		config: cfg,
	}
}

// Connect 连接到 MQTT Broker
func (mc *MQTTClient) Connect() error {
	token := mc.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connect error: %w", token.Error())
	}
	return nil
}

// Subscribe 订阅主题
func (mc *MQTTClient) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := mc.client.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT subscribe error: %w", token.Error())
	}
	return nil
}

// Publish 发布消息
func (mc *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := mc.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT publish error: %w", token.Error())
	}
	return nil
}

// Disconnect 断开连接
func (mc *MQTTClient) Disconnect(quiesce uint) {
	mc.client.Disconnect(quiesce)
}
