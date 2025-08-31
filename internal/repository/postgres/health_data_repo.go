package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
)

type HealthDataRepository struct {
	db *sql.DB
}

func NewHealthDataRepository(db *sql.DB) *HealthDataRepository {
	return &HealthDataRepository{db: db}
}

// Create 新增健康数据记录
func (r *HealthDataRepository) Create(record *models.HealthDataRecord) (int, error) {
	query := `INSERT INTO health_data_records (health_profile_id, device_id, schema_type, recorded_at, payload, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int
	payload, _ := json.Marshal(record.Payload)
	err := r.db.QueryRow(query, record.HealthProfileID, record.DeviceID, record.SchemaType, record.RecordedAt, payload, time.Now(), time.Now()).Scan(&id)
	return id, err
}

// Get 查询健康数据记录
func (r *HealthDataRepository) Get(id int64) (*models.HealthDataRecord, error) {
	query := `SELECT id, health_profile_id, device_id, schema_type, recorded_at, payload, created_at, updated_at FROM health_data_records WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var record models.HealthDataRecord
	var payload []byte
	err := row.Scan(&record.ID, &record.HealthProfileID, &record.DeviceID, &record.SchemaType, &record.RecordedAt, &payload, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		return nil, err
	}
	record.Payload = payload
	return &record, nil
}

// Update 更新健康数据记录
func (r *HealthDataRepository) Update(id int64, record *models.HealthDataRecord) error {
	query := `UPDATE health_data_records SET health_profile_id=$1, device_id=$2, schema_type=$3, recorded_at=$4, payload=$5, updated_at=$6 WHERE id=$7`
	payload, _ := json.Marshal(record.Payload)
	_, err := r.db.Exec(query, record.HealthProfileID, record.DeviceID, record.SchemaType, record.RecordedAt, payload, time.Now(), id)
	return err
}

// Delete 删除健康数据记录
func (r *HealthDataRepository) Delete(id int64) error {
	query := `DELETE FROM health_data_records WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}
