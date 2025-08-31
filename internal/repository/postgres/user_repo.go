// Package postgres 用户数据仓储实现，仅限本文件
package postgres

import (
	"context"
	"database/sql"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// User 用户数据结构（仅仓储内部使用，可扩展）
type User struct {
	ID       int64
	Username string
	Phone    string
	Password string
}

// UserRepo 用户仓储结构体
type UserRepository interface {
	Create(ctx context.Context, user *models.AppUser) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.AppUser, error)
	GetByUsername(ctx context.Context, username string) (*models.AppUser, error)
	GetByPhone(ctx context.Context, phone string) (*models.AppUser, error)
	Update(ctx context.Context, user *models.AppUser) error
	Delete(ctx context.Context, id int64) error

	ExistsByUsername(ctx context.Context, username string) (bool, error)
	FindByUsername(ctx context.Context, username string) (*models.AppUser, error)
	FindByID(ctx context.Context, id int64) (*models.AppUser, error)
}

// UserRepo 用户仓储结构体
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo 创建用户仓储实例
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create 新增用户
func (r *UserRepo) Create(ctx context.Context, user *models.AppUser) (int64, error) {
	query := "INSERT INTO app_users (username, email, phone, password_hash, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id"
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Phone, user.PasswordHash, user.IsActive).Scan(&id)
	return id, err
}

// GetByID 根据ID查询用户
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.AppUser, error) {
	query := "SELECT id, username, email, phone, password_hash, is_active, last_login, wechat_openid, created_at, updated_at FROM app_users WHERE id = $1"
	user := &models.AppUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash,
		&user.IsActive, &user.LastLogin, &user.WechatOpenID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetByUsername 根据用户名查询用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.AppUser, error) {
	query := "SELECT id, username, email, phone, password_hash, is_active, last_login, wechat_openid, created_at, updated_at FROM app_users WHERE username = $1"
	user := &models.AppUser{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash,
		&user.IsActive, &user.LastLogin, &user.WechatOpenID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetByPhone 根据手机号查询用户
func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*models.AppUser, error) {
	query := "SELECT id, username, email, phone, password_hash, is_active, last_login, wechat_openid, created_at, updated_at FROM app_users WHERE phone = $1"
	user := &models.AppUser{}
	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash,
		&user.IsActive, &user.LastLogin, &user.WechatOpenID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// Update 更新用户信息
func (r *UserRepo) Update(ctx context.Context, user *models.AppUser) error {
	query := "UPDATE app_users SET username = $1, email = $2, phone = $3, password_hash = $4, is_active = $5, last_login = $6, wechat_openid = $7, updated_at = NOW() WHERE id = $8"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Phone, user.PasswordHash, user.IsActive, user.LastLogin, user.WechatOpenID, user.ID)
	return err
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM app_users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ExistsByUsername 判断用户名是否存在
func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := "SELECT COUNT(1) FROM app_users WHERE username = $1"
	var count int
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	return count > 0, err
}

// FindByUsername 查询用户（兼容接口）
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.AppUser, error) {
	return r.GetByUsername(ctx, username)
}

// FindByID 查询用户（兼容接口）
func (r *UserRepo) FindByID(ctx context.Context, id int64) (*models.AppUser, error) {
	return r.GetByID(ctx, id)
}

// 用户注册（插入 app_users 表并加密密码）
func (r *UserRepo) Register(ctx context.Context, user *models.AppUser, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO app_users (username, email, phone, password_hash, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Phone, string(hash), true, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}
