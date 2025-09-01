// Package health 实现心率数据处理器，复用 BaseHealthHandler 并实现 HealthHandler 接口。
package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/app"
	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
	"github.com/fire-disposal/health_DT_go/internal/repository/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// 适配 app.Pipeline 的 HealthDataProcessor 接口
func (h *HeartRateHandler) Handle(event app.HealthEvent) {
	if event.EventType != "heart_rate" {
		return
	}
	data, ok := event.Payload.(HeartRateEventData)
	if !ok {
		return
	}
	_ = h.HandleEvent(context.Background(), HealthEvent{
		Type: "heart_rate",
		Data: data,
	})
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
	eventData := event.Data.(HeartRateEventData)

	// Redis缓存分支
	{
		redisClient := redis.GetRedisClient()
		cacheKey := fmt.Sprintf("health_data:%s:heart_rate", eventData.UserID)
		cacheValue, _ := json.Marshal(eventData)
		redisClient.Set(ctx, cacheKey, cacheValue, 5*time.Minute)
	}

	// 实际落库逻辑
	db, ok := ctx.Value("db").(*sql.DB)
	if !ok || db == nil {
		return errors.New("数据库连接未注入")
	}
	repo := postgres.NewHealthDataRepository(db)
	record := &models.HealthDataRecord{
		HealthProfileID: 0, // 可根据业务补充
		DeviceID:        nil,
		SchemaType:      "heart_rate",
		RecordedAt:      time.Unix(eventData.Timestamp, 0),
		Payload:         nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	payload, _ := json.Marshal(eventData)
	record.Payload = payload
	_, err := repo.Create(record)
	if err != nil {
		return err
	}

	// 发布事件通知（如有 eventbus，可扩展通知逻辑）
	// bus, ok := ctx.Value("eventbus").(*eventbus.EventBus)
	// if ok && bus != nil {
	// 	bus.Publish("heart_rate_created", eventData)
	// }

	zap.L().Info("心率数据已入库",
		zap.String("user_id", eventData.UserID),
		zap.Int("heart_rate", eventData.HeartRate),
		zap.Int64("timestamp", eventData.Timestamp),
	)
	return nil
}

// CreateHeartRateHandler 创建心率数据
func CreateHeartRateHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.HealthDataRecord
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		repo := postgres.NewHealthDataRepository(db)
		id, err := repo.Create(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id, "message": "created"})
	}
}

// GetHeartRateHandler 查询心率数据
func GetHeartRateHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		repo := postgres.NewHealthDataRepository(db)
		record, err := repo.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
	}
}

// UpdateHeartRateHandler 更新心率数据
func UpdateHeartRateHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var req models.HealthDataRecord
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		repo := postgres.NewHealthDataRepository(db)
		err := repo.Update(id, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
	}
}

// DeleteHeartRateHandler 删除心率数据
func DeleteHeartRateHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		repo := postgres.NewHealthDataRepository(db)
		err := repo.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}
