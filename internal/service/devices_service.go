// Package service 设备业务逻辑服务
package service

import (
	"context"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
)

type DevicesService struct {
	repo *postgres.DevicesRepository
}

func NewDevicesService(repo *postgres.DevicesRepository) *DevicesService {
	return &DevicesService{repo: repo}
}

func (s *DevicesService) Create(ctx context.Context, device *models.Device) (int, error) {
	return s.repo.Create(ctx, device)
}

func (s *DevicesService) Get(ctx context.Context, id int) (*models.Device, error) {
	return s.repo.Get(ctx, id)
}

func (s *DevicesService) Update(ctx context.Context, device *models.Device) error {
	return s.repo.Update(ctx, device)
}

func (s *DevicesService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *DevicesService) List(ctx context.Context) ([]models.Device, error) {
	return s.repo.FindAll(ctx)
}

// 设备绑定健康档案
func (s *DevicesService) AssignDeviceToProfile(ctx context.Context, deviceID int, profileID int) error {
	return s.repo.AssignDeviceToProfile(ctx, deviceID, profileID)
}
