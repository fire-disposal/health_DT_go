package postgres

import (
	"database/sql"

	"github.com/fire-disposal/health_DT_go/internal/models"
)

// AlertsRepository 告警数据仓储
type AlertsRepository struct {
	db *sql.DB
}

func NewAlertsRepository(db *sql.DB) *AlertsRepository {
	return &AlertsRepository{db: db}
}

// FindAll 查询全部告警（可扩展分页/筛选）
func (r *AlertsRepository) FindAll() ([]models.Alert, error) {
	rows, err := r.db.Query("SELECT id, health_profile_id, device_id, rule_name, level, message, status, created_at, resolved_at FROM alerts LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		// 只映射部分字段，完整字段可补充
		err := rows.Scan(&a.ID, &a.HealthProfileID, &a.DeviceID, &a.RuleName, &a.Level, &a.Message, &a.Status, &a.CreatedAt, &a.ResolvedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}
