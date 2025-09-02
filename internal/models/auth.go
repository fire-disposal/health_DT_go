// Package models 定义鉴权相关数据结构。
package models

import (
	"time"
)

// Auth 鉴权模型，兼容 Ent 及标准 struct 扩展。
// swagger:model Auth
type Auth struct {
	ID        int64     `json:"id"`         // 认证ID
	UserID    int64     `json:"user_id"`    // 用户ID
	Token     string    `json:"token"`      // 认证Token
	ExpiresAt time.Time `json:"expires_at"` // Token过期时间
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
