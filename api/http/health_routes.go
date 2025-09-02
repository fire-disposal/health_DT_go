// Package http 健康数据记录通用路由
package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterHealthDataRoutes 注册健康数据记录通用 CRUD 路由
func RegisterHealthDataRoutes(router gin.IRouter, db *sql.DB) {
	healthGroup := router.Group("/health_data")
	{
		healthGroup.POST("", createHealthDataHandler())
		healthGroup.GET("/:id", getHealthDataHandler())
		healthGroup.PUT("/:id", updateHealthDataHandler())
		healthGroup.DELETE("/:id", deleteHealthDataHandler())
	}
	// 注册告警和事件 RESTful 路由
	RegisterAlertsRoutes(router, db)
}

// @Summary 创建健康数据记录
// @Description 新增一条健康数据
// @Tags health_data
// @Accept json
// @Produce json
// @Param data body models.HealthDataRecord true "健康数据内容"
// @Success 201 {object} map[string]string
// @Router /health_data [post]
func createHealthDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.HealthDataRecord
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// TODO: 仓储层保存 req
		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	}
}

// @Summary 查询健康数据记录
// @Description 根据ID获取健康数据
// @Tags health_data
// @Produce json
// @Param id path int true "健康数据ID"
// @Success 200 {object} map[string]interface{}
// @Router /health_data/{id} [get]
func getHealthDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		// TODO: 仓储层查询
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

// @Summary 更新健康数据记录
// @Description 根据ID更新健康数据
// @Tags health_data
// @Accept json
// @Produce json
// @Param id path int true "健康数据ID"
// @Param data body models.HealthDataRecord true "健康数据内容"
// @Success 200 {object} map[string]interface{}
// @Router /health_data/{id} [put]
func updateHealthDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var req models.HealthDataRecord
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// TODO: 仓储层更新
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
	}
}

// @Summary 删除健康数据记录
// @Description 根据ID删除健康数据
// @Tags health_data
// @Produce json
// @Param id path int true "健康数据ID"
// @Success 200 {object} map[string]interface{}
// @Router /health_data/{id} [delete]
func deleteHealthDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		// TODO: 仓储层删除
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}
