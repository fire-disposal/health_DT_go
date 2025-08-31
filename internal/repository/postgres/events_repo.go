package postgres

import (
	"database/sql"

	"github.com/fire-disposal/health_DT_go/internal/models"
)

// EventsRepository 事件数据仓储
type EventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) *EventsRepository {
	return &EventsRepository{db: db}
}

// FindAll 查询全部事件（可扩展分页/筛选）
func (r *EventsRepository) FindAll() ([]models.Event, error) {
	rows, err := r.db.Query("SELECT id, event_type, health_profile_id, device_id, source_record_id, timestamp, data, metadata, created_at, updated_at FROM events LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		// 只映射部分字段，完整字段可补充
		err := rows.Scan(&e.ID, &e.EventType, &e.HealthProfileID, &e.DeviceID, &e.SourceRecordID, &e.Timestamp, &e.Data, &e.Metadata, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
