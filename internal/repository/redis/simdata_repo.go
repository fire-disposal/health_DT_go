// simdata_repo.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/simdata"
	"github.com/redis/go-redis/v9"
)

// SimDataRepository 管理模拟健康数据的 Redis 存取
type SimDataRepository struct {
	client *redis.Client
}

// NewSimDataRepository 构造
func NewSimDataRepository(client *redis.Client) *SimDataRepository {
	return &SimDataRepository{client: client}
}

// SaveHeartRate 保存心率模拟数据
func (r *SimDataRepository) SaveHeartRate(ctx context.Context, data simdata.HeartRateEventData) error {
	key := fmt.Sprintf("sim:heart_rate:%s:%d", data.UserID, data.Timestamp)
	val, _ := json.Marshal(data)
	return r.client.Set(ctx, key, val, time.Hour).Err()
}

// GetHeartRate 获取心率模拟数据
func (r *SimDataRepository) GetHeartRate(ctx context.Context, userID string, timestamp int64) (*simdata.HeartRateEventData, error) {
	key := fmt.Sprintf("sim:heart_rate:%s:%d", userID, timestamp)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data simdata.HeartRateEventData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// SaveBloodPressure 保存血压模拟数据
func (r *SimDataRepository) SaveBloodPressure(ctx context.Context, data simdata.BloodPressureEventData) error {
	key := fmt.Sprintf("sim:blood_pressure:%s:%d", data.UserID, data.Timestamp)
	val, _ := json.Marshal(data)
	return r.client.Set(ctx, key, val, time.Hour).Err()
}

// GetBloodPressure 获取血压模拟数据
func (r *SimDataRepository) GetBloodPressure(ctx context.Context, userID string, timestamp int64) (*simdata.BloodPressureEventData, error) {
	key := fmt.Sprintf("sim:blood_pressure:%s:%d", userID, timestamp)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data simdata.BloodPressureEventData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// SaveSpO2 保存血氧模拟数据
func (r *SimDataRepository) SaveSpO2(ctx context.Context, data simdata.SpO2EventData) error {
	key := fmt.Sprintf("sim:spo2:%s:%d", data.UserID, data.Timestamp)
	val, _ := json.Marshal(data)
	return r.client.Set(ctx, key, val, time.Hour).Err()
}

// GetSpO2 获取血氧模拟数据
func (r *SimDataRepository) GetSpO2(ctx context.Context, userID string, timestamp int64) (*simdata.SpO2EventData, error) {
	key := fmt.Sprintf("sim:spo2:%s:%d", userID, timestamp)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data simdata.SpO2EventData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// SaveTemperature 保存体温模拟数据
func (r *SimDataRepository) SaveTemperature(ctx context.Context, data simdata.TemperatureEventData) error {
	key := fmt.Sprintf("sim:temperature:%s:%d", data.UserID, data.Timestamp)
	val, _ := json.Marshal(data)
	return r.client.Set(ctx, key, val, time.Hour).Err()
}

// GetTemperature 获取体温模拟数据
func (r *SimDataRepository) GetTemperature(ctx context.Context, userID string, timestamp int64) (*simdata.TemperatureEventData, error) {
	key := fmt.Sprintf("sim:temperature:%s:%d", userID, timestamp)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data simdata.TemperatureEventData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
