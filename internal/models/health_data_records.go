package models

import (
	"encoding/json"
	"time"
)

// HealthDataRecord 健康数据记录模型
// swagger:model HealthDataRecord
type HealthDataRecord struct {
	ID              int             `json:"id"`
	HealthProfileID int             `json:"health_profile_id"`
	DeviceID        *int            `json:"device_id"`
	SchemaType      string          `json:"schema_type"`
	RecordedAt      time.Time       `json:"recorded_at"`
	Payload         json.RawMessage `json:"payload"` // JSONB
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
