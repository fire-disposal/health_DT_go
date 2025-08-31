// JWT 工具
package auth

import (
	"sync"

	"github.com/fire-disposal/health_DT_go/config"
	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtSecret []byte
	once      sync.Once
)

// Claims 结构体，包含角色
type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 获取 JWT 密钥（只加载一次）
func JwtSecret() []byte {
	once.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			panic("无法加载配置: " + err.Error())
		}
		jwtSecret = []byte(cfg.JWTSecret)
	})
	return jwtSecret
}

// 生成 JWT token
func GenerateToken(userID int64, role string) (string, error) {
	claims := &Claims{
		UserID:           userID,
		Role:             role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 可根据需要设置过期时间等
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenObj.SignedString(JwtSecret())
}
