package http

import (
	"database/sql"
	"net/http"

	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// RegisterAlertsRoutes 注册告警相关路由
func RegisterAlertsRoutes(router gin.IRouter, db *sql.DB) {
	router.GET("/alerts", queryAlertsHandler(db))
}

/*
// @Summary 查询告警列表
// @Description 获取所有告警信息
// @Tags alerts
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /alerts [get]
*/
func queryAlertsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := postgres.NewAlertsRepository(db)
		alerts, err := repo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"alerts": alerts})
	}
}
