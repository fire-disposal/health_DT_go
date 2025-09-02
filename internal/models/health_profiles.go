package models

import (
	"encoding/json"
	"time"
)

// HealthProfile 健康档案模型
// swagger:model HealthProfile
type HealthProfile struct {
	ID        int             `json:"id"`
	UserID    *int            `json:"user_id"` // 外键 app_users(id)
	Name      string          `json:"name"`
	Gender    string          `json:"gender"`
	BirthDate *time.Time      `json:"birth_date"`
	Metadata  json.RawMessage `json:"metadata"` // JSONB 灵活字段
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
