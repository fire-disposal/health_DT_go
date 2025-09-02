// Package models 定义普通用户模型
package models

import (
	"time"
)

// AppUser 普通用户模型，便于扩展，可兼容 Ent 或标准 struct
// swagger:model AppUser
type AppUser struct {
	ID           int64     `json:"id"`            // 用户ID
	Username     string    `json:"username"`      // 用户名
	Email        string    `json:"email"`         // 邮箱
	Phone        string    `json:"phone"`         // 手机号
	PasswordHash string    `json:"password_hash"` // 密码哈希
	IsActive     bool      `json:"is_active"`     // 激活状态
	LastLogin    time.Time `json:"last_login"`    // 最后登录时间
	WechatOpenID string    `json:"wechat_openid"` // 微信openid
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}
