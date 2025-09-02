// Package http 健康档案相关路由与处理
package http

import (
	"net/http"
	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/service"
	"github.com/gin-gonic/gin"
)

var healthProfilesService *service.HealthProfilesService

func RegisterHealthProfilesRoutes(r *gin.Engine, svc *service.HealthProfilesService) {
	healthProfilesService = svc
	profilesGroup := r.Group("/health_profiles")
	{
		profilesGroup.POST("", createHealthProfileHandler())
		profilesGroup.GET("/:id", getHealthProfileHandler())
		profilesGroup.PUT("/:id", updateHealthProfileHandler())
		profilesGroup.DELETE("/:id", deleteHealthProfileHandler())
		profilesGroup.GET("", listHealthProfilesHandler())
		profilesGroup.POST("/:id/bind_device", bindProfileToDeviceHandler())
	}
}

/*
@Summary 创建健康档案
@Description 新增健康档案，需提交完整信息
@Tags HealthProfile
@Accept json
@Produce json
@Param body body models.HealthProfile true "健康档案信息"
@Success 201 {object} map[string]int "创建成功，返回档案ID"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "创建失败"
*/
func createHealthProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.HealthProfile
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := healthProfilesService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

/*
@Summary 获取健康档案详情
@Description 根据ID查询健康档案信息
@Tags HealthProfile
@Produce json
@Param id path int true "健康档案ID"
@Success 200 {object} models.HealthProfile "查询成功"
@Failure 404 {object} map[string]string "未找到"
*/
func getHealthProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		profile, err := healthProfilesService.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}

/*
@Summary 更新健康档案
@Description 根据ID更新健康档案信息
@Tags HealthProfile
@Accept json
@Produce json
@Param id path int true "健康档案ID"
@Param body body models.HealthProfile true "健康档案信息"
@Success 200 {object} map[string]interface{} "更新成功"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "更新失败"
*/
func updateHealthProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req models.HealthProfile
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.ID = id
		if err := healthProfilesService.Update(c.Request.Context(), &req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
	}
}

/*
@Summary 删除健康档案
@Description 根据ID删除健康档案
@Tags HealthProfile
@Produce json
@Param id path int true "健康档案ID"
@Success 200 {object} map[string]interface{} "删除成功"
@Failure 500 {object} map[string]string "删除失败"
*/
func deleteHealthProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := healthProfilesService.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}

/*
@Summary 健康档案列表
@Description 获取所有健康档案信息
@Tags HealthProfile
@Produce json
@Success 200 {array} models.HealthProfile "列表成功"
@Failure 500 {object} map[string]string "获取失败"
*/
func listHealthProfilesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		profiles, err := healthProfilesService.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, profiles)
	}
}

/*
@Summary 健康档案绑定设备
@Description 将健康档案绑定到指定设备
@Tags HealthProfile
@Accept json
@Produce json
@Param id path int true "健康档案ID"
@Param body body object true "设备ID"
@Success 200 {object} map[string]interface{} "绑定成功"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "绑定失败"
*/
func bindProfileToDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		profileID, _ := strconv.Atoi(c.Param("id"))
		var req struct {
			DeviceID int `json:"device_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := healthProfilesService.AssignProfileToDevice(c.Request.Context(), profileID, req.DeviceID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"profile_id": profileID, "device_id": req.DeviceID, "message": "bound"})
	}
}
