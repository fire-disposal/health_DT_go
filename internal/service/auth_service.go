// Package service 实现鉴权业务逻辑服务
package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fire-disposal/health_DT_go/config"
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

type LoginResult struct {
	Token  string
	UserID int64
	Role   string
}

// 登录业务方法，统一查找、校验、生成Token
func (s *AuthService) Login(loginType string, username string, password string) (*LoginResult, error) {
	var userID int64
	var role string
	switch loginType {
	case "admin":
		admin, err := s.repo.GetAdminUserByUsername(username)
		if err != nil || admin == nil {
			return nil, errors.New("用户名或密码错误")
		}
		ok, err := s.repo.VerifyPassword(context.Background(), admin.ID, password)
		if err != nil || !ok {
			return nil, errors.New("用户名或密码错误")
		}
		userID = admin.ID
		role = admin.Role
	case "app":
		user, err := s.repo.GetAppUserByUsername(username)
		if err != nil || user == nil {
			return nil, errors.New("用户名或密码错误")
		}
		if user.PasswordHash == "" {
			return nil, errors.New("该用户未设置密码，无法使用密码登录")
		}
		ok, err := s.repo.VerifyPassword(context.Background(), user.ID, password)
		if err != nil || !ok {
			return nil, errors.New("用户名或密码错误")
		}
		userID = user.ID
		role = "app"
	default:
		return nil, errors.New("未知登录类型")
	}
	token, err := s.GenerateToken(userID, 24*time.Hour)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}
	return &LoginResult{
		Token:  token,
		UserID: userID,
		Role:   role,
	}, nil
}

// 微信登录：通过 code 换 openid，查找/注册用户，生成 Token
func (s *AuthService) LoginWithWechatCode(code string) (*LoginResult, error) {
	// 1. 请求微信官方接口换取 openid
	// 2. 查找用户，无则自动注册
	// 3. 生成 Token 返回

	// 从配置读取微信 appid/secret
	cfg, err := config.Load()
	if err != nil {
		return nil, errors.New("微信配置加载失败")
	}
	appid := cfg.Wechat.AppID
	secret := cfg.Wechat.Secret
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + code + "&grant_type=authorization_code"

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("微信接口请求失败")
	}
	defer resp.Body.Close()
	var wxResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
		return nil, errors.New("微信响应解析失败")
	}
	if wxResp.OpenID == "" {
		return nil, errors.New("微信登录失败: " + wxResp.ErrMsg)
	}

	// 2. 查找用户
	user, err := s.repo.GetAppUserByWechatOpenID(wxResp.OpenID)
	if err != nil {
		return nil, errors.New("数据库查询失败")
	}
	if user == nil {
		// 自动注册
		user = &models.AppUser{
			WechatOpenID: wxResp.OpenID,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		// 注册到数据库（无密码，仅微信登录）
		id, err := s.repo.CreateAppUser(context.Background(), user)
		if err != nil {
			return nil, errors.New("微信用户注册失败")
		}
		user.ID = id
	}

	// 3. 生成 Token
	token, err := s.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}
	return &LoginResult{
		Token:  token,
		UserID: user.ID,
		Role:   "app",
	}, nil
}
