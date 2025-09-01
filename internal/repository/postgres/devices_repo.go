// Package postgres 设备数据仓储实现
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
)

type DevicesRepository struct {
	db *sql.DB
}

func NewDevicesRepository(db *sql.DB) *DevicesRepository {
	return &DevicesRepository{db: db}
}

func (r *DevicesRepository) Create(ctx context.Context, device *models.Device) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO devices (serial_number, name, device_type, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		device.SerialNumber, device.Name, device.DeviceType, device.IsActive, device.CreatedAt, device.UpdatedAt,
	).Scan(&id)
	return id, err
}

func (r *DevicesRepository) Get(ctx context.Context, id int) (*models.Device, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, serial_number, name, device_type, is_active, created_at, updated_at FROM devices WHERE id = $1`, id)
	var device models.Device
	err := row.Scan(&device.ID, &device.SerialNumber, &device.Name, &device.DeviceType, &device.IsActive, &device.CreatedAt, &device.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DevicesRepository) Update(ctx context.Context, device *models.Device) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE devices SET serial_number=$1, name=$2, device_type=$3, is_active=$4, updated_at=$5 WHERE id=$6`,
		device.SerialNumber, device.Name, device.DeviceType, device.IsActive, device.UpdatedAt, device.ID,
	)
	return err
}

func (r *DevicesRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id=$1`, id)
	return err
}

func (r *DevicesRepository) FindAll(ctx context.Context) ([]models.Device, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, serial_number, name, device_type, is_active, created_at, updated_at FROM devices`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.SerialNumber, &d.Name, &d.DeviceType, &d.IsActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// 设备绑定健康档案
func (r *DevicesRepository) AssignDeviceToProfile(ctx context.Context, deviceID int, profileID int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO device_assignments (device_id, health_profile_id, assigned_at) VALUES ($1, $2, $3)`,
		deviceID, profileID, time.Now(),
	)
	return err
}
