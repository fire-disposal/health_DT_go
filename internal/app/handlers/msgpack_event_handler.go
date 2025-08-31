// Package handlers 提供 msgpack 业务事件处理
package handlers

import (
	"github.com/fire-disposal/health_DT_go/internal/app"
)

func HandleMsgpackPayload(pipeline *app.Pipeline) func(payload map[string]interface{}) {
	return func(payload map[string]interface{}) {
		deviceSN, ok := payload["sn"].(string)
		if !ok || deviceSN == "" {
			// 可加日志
			return
		}
		event := app.HealthEvent{
			DeviceID:  deviceSN,
			EventType: "mattress",
			Payload:   payload,
			Source:    "msgpack",
		}
		go pipeline.ReceiveEvent(event)
	}
}
