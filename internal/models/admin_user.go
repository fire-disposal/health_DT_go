// internal/models/admin_user.go
package models

import (
	"time"
)

// AdminUser 管理员用户模型，便于扩展，可兼容 Ent 或标准 struct
type AdminUser struct {
	ID           int64     `json:"id"`            // 管理员ID
	Username     string    `json:"username"`      // 用户名
	Email        string    `json:"email"`         // 邮箱
	Phone        string    `json:"phone"`         // 手机号
	PasswordHash string    `json:"password_hash"` // 密码哈希
	Role         string    `json:"role"`          // 角色（如：superadmin, admin）
	IsActive     bool      `json:"is_active"`     // 激活状态
	LastLogin    time.Time `json:"last_login"`    // 最后登录时间
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}
