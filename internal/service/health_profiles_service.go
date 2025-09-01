// Package service 健康档案业务逻辑服务
package service

import (
	"context"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
)

type HealthProfilesService struct {
	repo *postgres.HealthProfilesRepository
}

func NewHealthProfilesService(repo *postgres.HealthProfilesRepository) *HealthProfilesService {
	return &HealthProfilesService{repo: repo}
}

func (s *HealthProfilesService) Create(ctx context.Context, profile *models.HealthProfile) (int, error) {
	return s.repo.Create(ctx, profile)
}

func (s *HealthProfilesService) Get(ctx context.Context, id int) (*models.HealthProfile, error) {
	return s.repo.Get(ctx, id)
}

func (s *HealthProfilesService) Update(ctx context.Context, profile *models.HealthProfile) error {
	return s.repo.Update(ctx, profile)
}

func (s *HealthProfilesService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *HealthProfilesService) List(ctx context.Context) ([]models.HealthProfile, error) {
	return s.repo.FindAll(ctx)
}

// 健康档案绑定设备
func (s *HealthProfilesService) AssignProfileToDevice(ctx context.Context, profileID int, deviceID int) error {
	return s.repo.AssignProfileToDevice(ctx, profileID, deviceID)
}
