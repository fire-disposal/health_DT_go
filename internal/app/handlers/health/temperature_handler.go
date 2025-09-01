// Package health 实现体温数据处理器，复用 BaseHealthHandler 并实现 HealthHandler 接口。
package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/app"
	"github.com/fire-disposal/health_DT_go/internal/repository/redis"
	"go.uber.org/zap"
)

// TemperatureEventData 表示体温事件的数据结构，可扩展字段。
type TemperatureEventData struct {
	UserID      string
	Temperature float64 // 体温（℃）
	Timestamp   int64
}

// TemperatureHandler 体温数据处理器，实现 HealthHandler，嵌入 BaseHealthHandler。
type TemperatureHandler struct {
	BaseHealthHandler
}

// 适配 app.Pipeline 的 HealthDataProcessor 接口
func (h *TemperatureHandler) Handle(event app.HealthEvent) {
	if event.EventType != "temperature" {
		return
	}
	data, ok := event.Payload.(TemperatureEventData)
	if !ok {
		return
	}
	_ = h.HandleEvent(context.Background(), HealthEvent{
		Type: "temperature",
		Data: data,
	})
}

// ValidateData 校验体温数据的有效性。
func (h *TemperatureHandler) ValidateData(data interface{}) error {
	eventData, ok := data.(TemperatureEventData)
	if !ok {
		return errors.New("数据类型错误，需为 TemperatureEventData")
	}
	if eventData.Temperature < 34.0 || eventData.Temperature > 42.0 {
		return fmt.Errorf("体温值异常: %.1f", eventData.Temperature)
	}
	if eventData.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// HandleEvent 处理体温事件，校验并执行业务逻辑。
func (h *TemperatureHandler) HandleEvent(ctx context.Context, event HealthEvent) error {
	if event.Type != "temperature" {
		return errors.New("事件类型错误，仅支持 temperature")
	}
	if err := h.ValidateData(event.Data); err != nil {
		return err
	}
	eventData := event.Data.(TemperatureEventData)

	// Redis缓存分支
	{
		redisClient := redis.GetRedisClient()
		cacheKey := fmt.Sprintf("health_data:%s:temperature", eventData.UserID)
		cacheValue, _ := json.Marshal(eventData)
		redisClient.Set(ctx, cacheKey, cacheValue, 5*time.Minute)
	}

	// 业务处理逻辑（可扩展，如存储、通知等）
	zap.L().Info("体温数据",
		zap.String("user_id", eventData.UserID),
		zap.Float64("temperature", eventData.Temperature),
		zap.Int64("timestamp", eventData.Timestamp),
	)
	return nil
}
