package http

import (
	"database/sql"
	"net/http"

	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// RegisterEventsRoutes 注册事件相关路由
func RegisterEventsRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/events", queryEventsHandler(db))
}

/*
// @Summary 查询事件列表
// @Description 获取所有事件信息
// @Tags events
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /events [get]
*/
func queryEventsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := postgres.NewEventsRepository(db)
		events, err := repo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"events": events})
	}
}
