package models

import (
	"time"
)

// DeviceAssignment 设备与健康档案绑定模型
type DeviceAssignment struct {
	ID              int        `json:"id"`
	DeviceID        int        `json:"device_id"`
	HealthProfileID int        `json:"health_profile_id"`
	AssignedAt      time.Time  `json:"assigned_at"`
	UnassignedAt    *time.Time `json:"unassigned_at"`
}
