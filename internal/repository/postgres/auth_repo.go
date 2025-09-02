// Package postgres 实现鉴权数据仓储
package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthRepository 定义鉴权数据仓储接口
type AuthRepository interface {
	Create(ctx context.Context, auth *models.Auth) error
	GetByID(ctx context.Context, id int64) (*models.Auth, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Auth, error)
	GetByToken(ctx context.Context, token string) (*models.Auth, error)
	Update(ctx context.Context, auth *models.Auth) error
	Delete(ctx context.Context, id int64) error

	// 新增 Token 相关接口
	CreateToken(userID int64, expireDuration time.Duration) (string, error)
	GetToken(token string) (*models.Auth, error)
	DeleteToken(token string) error

	// 密码相关接口
	SetPassword(ctx context.Context, userID int64, password string) error
	VerifyPassword(ctx context.Context, userID int64, password string) (bool, error)
	// 用户查找接口
	GetAdminUserByUsername(username string) (*models.AdminUser, error)
	GetAppUserByUsername(username string) (*models.AppUser, error)
	// 新增：通过微信 openid 查询 app_user
	GetAppUserByWechatOpenID(openid string) (*models.AppUser, error)
	// 新增：创建 app_user（用于微信自动注册）
	CreateAppUser(ctx context.Context, user *models.AppUser) (int64, error)
}

// authRepo 实现 AuthRepository
type authRepo struct {
	db *sql.DB
}

// NewAuthRepository 创建鉴权仓储实例
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) Create(ctx context.Context, auth *models.Auth) error {
	query := `INSERT INTO auth (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, auth.UserID, auth.Token, auth.ExpiresAt).Scan(&auth.ID)
}

func (r *authRepo) GetByID(ctx context.Context, id int64) (*models.Auth, error) {
	query := `SELECT id, user_id, token, expires_at FROM auth WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var a models.Auth
	if err := row.Scan(&a.ID, &a.UserID, &a.Token, &a.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *authRepo) GetByUserID(ctx context.Context, userID int64) ([]*models.Auth, error) {
	query := `SELECT id, user_id, token, expires_at FROM auth WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*models.Auth
	for rows.Next() {
		var a models.Auth
		if err := rows.Scan(&a.ID, &a.UserID, &a.Token, &a.ExpiresAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *authRepo) GetByToken(ctx context.Context, token string) (*models.Auth, error) {
	query := `SELECT id, user_id, token, expires_at FROM auth WHERE token = $1`
	row := r.db.QueryRowContext(ctx, query, token)
	var a models.Auth
	if err := row.Scan(&a.ID, &a.UserID, &a.Token, &a.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *authRepo) Update(ctx context.Context, auth *models.Auth) error {
	query := `UPDATE auth SET user_id = $1, token = $2, expires_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, auth.UserID, auth.Token, auth.ExpiresAt, auth.ID)
	return err
}

// 生成随机 token
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *authRepo) CreateToken(userID int64, expireDuration time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(expireDuration)
	query := `INSERT INTO auth (user_id, token, expires_at, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	var id int64
	err = r.db.QueryRowContext(context.Background(), query, userID, token, expiresAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *authRepo) GetToken(token string) (*models.Auth, error) {
	query := `SELECT id, user_id, token, expires_at, created_at, updated_at FROM auth WHERE token = $1`
	row := r.db.QueryRowContext(context.Background(), query, token)
	var a models.Auth
	err := row.Scan(&a.ID, &a.UserID, &a.Token, &a.ExpiresAt, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *authRepo) DeleteToken(token string) error {
	query := `DELETE FROM auth WHERE token = $1`
	_, err := r.db.ExecContext(context.Background(), query, token)
	return err
}
func (r *authRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM auth WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// SetPassword 设置用户密码（bcrypt加密存储到 app_users 表）

func (r *authRepo) SetPassword(ctx context.Context, userID int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE app_users SET password_hash = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, string(hash), userID)
	return err
}

// VerifyPassword 校验用户密码
func (r *authRepo) VerifyPassword(ctx context.Context, userID int64, password string) (bool, error) {
	query := "SELECT password_hash FROM app_users WHERE id = $1"
	var hash string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, nil
}

// 查询管理员用户
func (r *authRepo) GetAdminUserByUsername(username string) (*models.AdminUser, error) {
	query := "SELECT id, username, password, role, created_at, updated_at FROM admin_users WHERE username = $1"
	row := r.db.QueryRow(query, username)
	var admin models.AdminUser
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.Phone,
		&admin.PasswordHash,
		&admin.Role,
		&admin.IsActive,
		&admin.LastLogin,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// 通过微信 openid 查询 app_user
func (r *authRepo) GetAppUserByWechatOpenID(openid string) (*models.AppUser, error) {
	query := "SELECT id, username, email, phone, password_hash, is_active, last_login, wechat_openid, created_at, updated_at FROM app_users WHERE wechat_openid = $1"
	row := r.db.QueryRow(query, openid)
	var user models.AppUser
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.IsActive, &user.LastLogin, &user.WechatOpenID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 创建 app_user（用于微信自动注册）
func (r *authRepo) CreateAppUser(ctx context.Context, user *models.AppUser) (int64, error) {
	query := "INSERT INTO app_users (username, email, phone, password_hash, is_active, wechat_openid, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Phone, user.PasswordHash, user.IsActive, user.WechatOpenID, user.CreatedAt, user.UpdatedAt).Scan(&id)
	return id, err
}

// 查询普通用户
func (r *authRepo) GetAppUserByUsername(username string) (*models.AppUser, error) {
	query := "SELECT id, username, email, phone, password_hash, is_active, last_login, wechat_openid, created_at, updated_at FROM app_users WHERE username = $1"
	row := r.db.QueryRow(query, username)
	var user models.AppUser
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.IsActive, &user.LastLogin, &user.WechatOpenID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
