package models

import (
	"time"
)

// Device 设备模型
type Device struct {
	ID           int       `json:"id"`
	SerialNumber string    `json:"serial_number"`
	Name         string    `json:"name"`
	DeviceType   string    `json:"device_type"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
