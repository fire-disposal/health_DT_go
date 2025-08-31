// Package handlers 实现 MQTT 消息事件处理器
package handlers

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fire-disposal/health_DT_go/internal/app"
)

// HandleMQTTMessage 解析 MQTT 消息并分发到 Pipeline
func HandleMQTTMessage(pipeline *app.Pipeline) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		parts := strings.Split(topic, "/")
		if len(parts) < 4 || parts[0] != "device" {
			return
		}
		deviceID := parts[1]
		messageType := parts[2]
		dataType := parts[3]
		if messageType != "data" {
			return
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(payload, &raw); err != nil {
			return
		}
		dataField, ok := raw["data"].(map[string]interface{})
		if !ok {
			return
		}
		event := app.HealthEvent{
			DeviceID:  deviceID,
			EventType: dataType,
			Payload:   dataField,
			Source:    "mqtt",
		}
		go pipeline.ReceiveEvent(event)
	}
}
