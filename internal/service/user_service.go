// 用户服务层
package service

import (
	"context"
	"errors"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
)

var userRepo *postgres.UserRepo

// InitUserRepo 用于初始化仓储（需在 main 或 handler 层调用）
func InitUserRepo(repo *postgres.UserRepo) {
	userRepo = repo
}

// 获取单个用户
func GetUserByID(id int64) (*models.AppUser, error) {
	if userRepo == nil {
		return nil, errors.New("UserRepo未初始化")
	}
	return userRepo.GetByID(context.Background(), id)
}

// 获取用户列表（示例：可扩展为分页）
func ListUsers() ([]*models.AppUser, error) {
	if userRepo == nil {
		return nil, errors.New("UserRepo未初始化")
	}
	// 简单实现：遍历所有ID（实际应分页/优化）
	users := []*models.AppUser{}
	// TODO: 可扩展为 SELECT * FROM app_users
	return users, nil
}

// 创建用户
func CreateUser(u *models.AppUser) (*models.AppUser, error) {
	if userRepo == nil {
		return nil, errors.New("UserRepo未初始化")
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	id, err := userRepo.Create(context.Background(), u)
	if err != nil {
		return nil, err
	}
	return userRepo.GetByID(context.Background(), id)
}

// 更新用户
func UpdateUser(u *models.AppUser) (*models.AppUser, error) {
	if userRepo == nil {
		return nil, errors.New("UserRepo未初始化")
	}
	u.UpdatedAt = time.Now()
	err := userRepo.Update(context.Background(), u)
	if err != nil {
		return nil, err
	}
	return userRepo.GetByID(context.Background(), u.ID)
}

// UserService 结构体实现用户相关服务
type UserService struct {
	Repo *postgres.UserRepo
}

func NewUserService(repo *postgres.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(ctx context.Context, user *models.AppUser, password string) error {
	if s.Repo == nil {
		return errors.New("UserRepo未初始化")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.Repo.Register(ctx, user, password)
}

func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*models.AppUser, error) {
	if s.Repo == nil {
		return nil, errors.New("UserRepo未初始化")
	}
	return s.Repo.GetByID(ctx, userID)
}
