package models

import (
	"encoding/json"
	"time"
)

// Event 事件模型
// swagger:model Event
type Event struct {
	ID              int             `json:"id"`
	EventType       string          `json:"event_type"`
	HealthProfileID int             `json:"health_profile_id"`
	DeviceID        *int            `json:"device_id"`
	SourceRecordID  *int            `json:"source_record_id"`
	Timestamp       time.Time       `json:"timestamp"`
	Data            json.RawMessage `json:"data"`
	Metadata        json.RawMessage `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
