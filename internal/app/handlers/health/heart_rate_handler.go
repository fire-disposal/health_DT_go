// Package health 实现心率数据处理器，复用 BaseHealthHandler 并实现 HealthHandler 接口。
package health

import (
	"context"
	"errors"
	"fmt"
)

// HeartRateEventData 表示心率事件的数据结构，可扩展字段。
type HeartRateEventData struct {
	UserID    string
	HeartRate int
	Timestamp int64
}

// HeartRateHandler 心率数据处理器，实现 HealthHandler，嵌入 BaseHealthHandler。
type HeartRateHandler struct {
	BaseHealthHandler
}

// ValidateData 校验心率数据的有效性。
func (h *HeartRateHandler) ValidateData(data interface{}) error {
	eventData, ok := data.(HeartRateEventData)
	if !ok {
		return errors.New("数据类型错误，需为 HeartRateEventData")
	}
	if eventData.HeartRate < 30 || eventData.HeartRate > 220 {
		return fmt.Errorf("心率值异常: %d", eventData.HeartRate)
	}
	if eventData.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// HandleEvent 处理心率事件，校验并执行业务逻辑。
func (h *HeartRateHandler) HandleEvent(ctx context.Context, event HealthEvent) error {
	if event.Type != "heart_rate" {
		return errors.New("事件类型错误，仅支持 heart_rate")
	}
	if err := h.ValidateData(event.Data); err != nil {
		return err
	}
	// 业务处理逻辑（可扩展，如存储、通知等）
	// 示例：打印心率数据
	eventData := event.Data.(HeartRateEventData)
	fmt.Printf("用户 %s 心率数据: %d @ %d\n", eventData.UserID, eventData.HeartRate, eventData.Timestamp)
	return nil
}
