// Package health 实现血压数据处理器，复用 BaseHealthHandler 并实现 HealthHandler 接口。
package health

import (
	"context"
	"errors"
	"fmt"

	"github.com/fire-disposal/health_DT_go/internal/app"
	"go.uber.org/zap"
)

// BloodPressureEventData 表示血压事件的数据结构，可扩展字段。
type BloodPressureEventData struct {
	UserID    string
	Systolic  int // 收缩压
	Diastolic int // 舒张压
	Timestamp int64
}

// BloodPressureHandler 血压数据处理器，实现 HealthHandler，嵌入 BaseHealthHandler。
type BloodPressureHandler struct {
	BaseHealthHandler
}

// 适配 app.Pipeline 的 HealthDataProcessor 接口
func (h *BloodPressureHandler) Handle(event app.HealthEvent) {
	if event.EventType != "blood_pressure" {
		return
	}
	data, ok := event.Payload.(BloodPressureEventData)
	if !ok {
		return
	}
	_ = h.HandleEvent(context.Background(), HealthEvent{
		Type: "blood_pressure",
		Data: data,
	})
}

// ValidateData 校验血压数据的有效性。
func (h *BloodPressureHandler) ValidateData(data interface{}) error {
	eventData, ok := data.(BloodPressureEventData)
	if !ok {
		return errors.New("数据类型错误，需为 BloodPressureEventData")
	}
	if eventData.Systolic < 60 || eventData.Systolic > 250 {
		return fmt.Errorf("收缩压值异常: %d", eventData.Systolic)
	}
	if eventData.Diastolic < 40 || eventData.Diastolic > 150 {
		return fmt.Errorf("舒张压值异常: %d", eventData.Diastolic)
	}
	if eventData.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// HandleEvent 处理血压事件，校验并执行业务逻辑。
func (h *BloodPressureHandler) HandleEvent(ctx context.Context, event HealthEvent) error {
	if event.Type != "blood_pressure" {
		return errors.New("事件类型错误，仅支持 blood_pressure")
	}
	if err := h.ValidateData(event.Data); err != nil {
		return err
	}
	// 业务处理逻辑（可扩展，如存储、通知等）
	// 示例：打印血压数据
	eventData := event.Data.(BloodPressureEventData)
	zap.L().Info("血压数据",
		zap.String("user_id", eventData.UserID),
		zap.Int("systolic", eventData.Systolic),
		zap.Int("diastolic", eventData.Diastolic),
		zap.Int64("timestamp", eventData.Timestamp),
	)
	return nil
}
