// internal/models/admin_user.go
package models

import (
	"time"
)

// AdminUser 管理员用户模型，便于扩展，可兼容 Ent 或标准 struct
type AdminUser struct {
	ID        int64     `json:"id"`         // 管理员ID
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码哈希
	Role      string    `json:"role"`       // 角色（如：superadmin, admin）
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
