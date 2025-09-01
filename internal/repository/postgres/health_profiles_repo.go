// Package postgres 健康档案数据仓储实现
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
)

type HealthProfilesRepository struct {
	db *sql.DB
}

func NewHealthProfilesRepository(db *sql.DB) *HealthProfilesRepository {
	return &HealthProfilesRepository{db: db}
}

func (r *HealthProfilesRepository) Create(ctx context.Context, profile *models.HealthProfile) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO health_profiles (user_id, name, gender, birth_date, metadata, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		profile.UserID, profile.Name, profile.Gender, profile.BirthDate, profile.Metadata, profile.CreatedAt, profile.UpdatedAt,
	).Scan(&id)
	return id, err
}

func (r *HealthProfilesRepository) Get(ctx context.Context, id int) (*models.HealthProfile, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, gender, birth_date, metadata, created_at, updated_at FROM health_profiles WHERE id = $1`, id)
	var profile models.HealthProfile
	err := row.Scan(&profile.ID, &profile.UserID, &profile.Name, &profile.Gender, &profile.BirthDate, &profile.Metadata, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *HealthProfilesRepository) Update(ctx context.Context, profile *models.HealthProfile) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE health_profiles SET user_id=$1, name=$2, gender=$3, birth_date=$4, metadata=$5, updated_at=$6 WHERE id=$7`,
		profile.UserID, profile.Name, profile.Gender, profile.BirthDate, profile.Metadata, profile.UpdatedAt, profile.ID,
	)
	return err
}

func (r *HealthProfilesRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM health_profiles WHERE id=$1`, id)
	return err
}

func (r *HealthProfilesRepository) FindAll(ctx context.Context) ([]models.HealthProfile, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, gender, birth_date, metadata, created_at, updated_at FROM health_profiles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var profiles []models.HealthProfile
	for rows.Next() {
		var p models.HealthProfile
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Gender, &p.BirthDate, &p.Metadata, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// 健康档案绑定设备
func (r *HealthProfilesRepository) AssignProfileToDevice(ctx context.Context, profileID int, deviceID int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO device_assignments (device_id, health_profile_id, assigned_at) VALUES ($1, $2, $3)`,
		deviceID, profileID, time.Now(),
	)
	return err
}
