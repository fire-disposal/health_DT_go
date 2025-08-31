package models

import (
	"encoding/json"
	"time"
)

// Alert 告警模型
type Alert struct {
	ID              int             `json:"id"`
	HealthProfileID *int            `json:"health_profile_id"`
	DeviceID        *int            `json:"device_id"`
	SourceEventID   *int            `json:"source_event_id"`
	RuleName        string          `json:"rule_name"`
	Level           string          `json:"level"`
	Message         string          `json:"message"`
	EventType       string          `json:"event_type"`
	Description     string          `json:"description"`
	Extra           json.RawMessage `json:"extra"`
	Status          string          `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	ResolvedAt      *time.Time      `json:"resolved_at"`
}
