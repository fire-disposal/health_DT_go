// Package http 设备相关路由与处理
package http

import (
	"net/http"
	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/service"
	"github.com/gin-gonic/gin"
)

var devicesService *service.DevicesService

func RegisterDevicesRoutes(r *gin.Engine, svc *service.DevicesService) {
	devicesService = svc
	devicesGroup := r.Group("/devices")
	{
		devicesGroup.POST("", createDeviceHandler())
		devicesGroup.GET("/:id", getDeviceHandler())
		devicesGroup.PUT("/:id", updateDeviceHandler())
		devicesGroup.DELETE("/:id", deleteDeviceHandler())
		devicesGroup.GET("", listDevicesHandler())
		devicesGroup.POST("/:id/bind_profile", bindDeviceToProfileHandler())
	}
}

func createDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Device
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := devicesService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func getDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		device, err := devicesService.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, device)
	}
}

func updateDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req models.Device
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.ID = id
		if err := devicesService.Update(c.Request.Context(), &req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
	}
}

func deleteDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := devicesService.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}

func listDevicesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		devices, err := devicesService.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, devices)
	}
}

// 设备绑定健康档案
func bindDeviceToProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID, _ := strconv.Atoi(c.Param("id"))
		var req struct {
			ProfileID int `json:"profile_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := devicesService.AssignDeviceToProfile(c.Request.Context(), deviceID, req.ProfileID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "profile_id": req.ProfileID, "message": "bound"})
	}
}
