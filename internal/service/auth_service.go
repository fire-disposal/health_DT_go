// Package service 实现鉴权业务逻辑服务
package service

import (
	"errors"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
)

// AuthService 提供 Token 相关业务方法
type AuthService struct {
	repo postgres.AuthRepository
}

// NewAuthService 构造鉴权服务
func NewAuthService(repo postgres.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// GenerateToken 生成 Token
func (s *AuthService) GenerateToken(userID int64, expireDuration time.Duration) (string, error) {
	token, err := s.repo.CreateToken(userID, expireDuration)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken 校验 Token 有效性
func (s *AuthService) ValidateToken(token string) (*models.Auth, error) {
	auth, err := s.repo.GetToken(token)
	if err != nil {
		return nil, err
	}
	if auth == nil || auth.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token invalid or expired")
	}
	return auth, nil
}

// InvalidateToken 使 Token 失效
func (s *AuthService) InvalidateToken(token string) error {
	return s.repo.DeleteToken(token)
}
