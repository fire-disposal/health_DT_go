// Package health 实现血氧数据处理器，复用 BaseHealthHandler 并实现 HealthHandler 接口。
package health

import (
	"context"
	"errors"
	"fmt"
)

// SpO2EventData 表示血氧事件的数据结构，可扩展字段。
type SpO2EventData struct {
	UserID    string
	SpO2      int // 血氧饱和度（%）
	Timestamp int64
}

// SpO2Handler 血氧数据处理器，实现 HealthHandler，嵌入 BaseHealthHandler。
type SpO2Handler struct {
	BaseHealthHandler
}

// ValidateData 校验血氧数据的有效性。
func (h *SpO2Handler) ValidateData(data interface{}) error {
	eventData, ok := data.(SpO2EventData)
	if !ok {
		return errors.New("数据类型错误，需为 SpO2EventData")
	}
	if eventData.SpO2 < 70 || eventData.SpO2 > 100 {
		return fmt.Errorf("血氧值异常: %d", eventData.SpO2)
	}
	if eventData.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// HandleEvent 处理血氧事件，校验并执行业务逻辑。
func (h *SpO2Handler) HandleEvent(ctx context.Context, event HealthEvent) error {
	if event.Type != "spo2" {
		return errors.New("事件类型错误，仅支持 spo2")
	}
	if err := h.ValidateData(event.Data); err != nil {
		return err
	}
	// 业务处理逻辑（可扩展，如存储、通知等）
	// 示例：打印血氧数据
	eventData := event.Data.(SpO2EventData)
	fmt.Printf("用户 %s 血氧数据: SpO2 %d%% @ %d\n", eventData.UserID, eventData.SpO2, eventData.Timestamp)
	return nil
}
