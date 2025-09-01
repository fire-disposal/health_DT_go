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

// 健康档案绑定设备
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
