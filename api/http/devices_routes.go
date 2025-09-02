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

func RegisterDevicesRoutes(router gin.IRouter, svc *service.DevicesService) {
	devicesService = svc
	devicesGroup := router.Group("/devices")
	{
		devicesGroup.POST("", createDeviceHandler())
		devicesGroup.GET("/:id", getDeviceHandler())
		devicesGroup.PUT("/:id", updateDeviceHandler())
		devicesGroup.DELETE("/:id", deleteDeviceHandler())
		devicesGroup.GET("", listDevicesHandler())
		devicesGroup.POST("/:id/bind_profile", bindDeviceToProfileHandler())
	}
}

/*
@Summary 创建设备
@Description 新增设备，需提交设备信息
@Tags Device
@Accept json
@Produce json
@Param body body models.Device true "设备信息"
@Success 201 {object} map[string]int "创建成功，返回设备ID"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "创建失败"
*/
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

/*
@Summary 获取设备详情
@Description 根据ID查询设备信息
@Tags Device
@Produce json
@Param id path int true "设备ID"
@Success 200 {object} models.Device "查询成功"
@Failure 404 {object} map[string]string "未找到"
*/
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

/*
@Summary 更新设备信息
@Description 根据ID更新设备信息
@Tags Device
@Accept json
@Produce json
@Param id path int true "设备ID"
@Param body body models.Device true "设备信息"
@Success 200 {object} map[string]interface{} "更新成功"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "更新失败"
*/
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

/*
@Summary 删除设备
@Description 根据ID删除设备
@Tags Device
@Produce json
@Param id path int true "设备ID"
@Success 200 {object} map[string]interface{} "删除成功"
@Failure 500 {object} map[string]string "删除失败"
*/
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

/*
@Summary 设备列表
@Description 获取所有设备信息
@Tags Device
@Produce json
@Success 200 {array} models.Device "列表成功"
@Failure 500 {object} map[string]string "获取失败"
*/
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

/*
@Summary 设备绑定健康档案
@Description 将设备绑定到指定健康档案
@Tags Device
@Accept json
@Produce json
@Param id path int true "设备ID"
@Param body body object true "健康档案ID"
@Success 200 {object} map[string]interface{} "绑定成功"
@Failure 400 {object} map[string]string "参数错误"
@Failure 500 {object} map[string]string "绑定失败"
*/
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
