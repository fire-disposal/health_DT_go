// Package http 设备绑定CRUD路由
package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/gin-gonic/gin"
)

var deviceAssignmentsStore = make(map[int]*models.DeviceAssignment)
var nextAssignmentID int = 1

func RegisterDeviceAssignmentsRoutes(r *gin.Engine) {
	group := r.Group("/device_assignments")
	{
		group.POST("", createDeviceAssignmentHandler())
		group.GET("/:id", getDeviceAssignmentHandler())
		group.GET("", listDeviceAssignmentsHandler())
		group.PUT("/:id/unassign", unassignDeviceHandler())
		group.DELETE("/:id", deleteDeviceAssignmentHandler())
	}
}

/*
@Summary 创建设备绑定
@Description 新增设备绑定关系，需提交完整信息
@Tags DeviceAssignment
@Accept json
@Produce json
@Param body body models.DeviceAssignment true "设备绑定信息"
@Success 201 {object} models.DeviceAssignment "创建成功"
@Failure 400 {object} map[string]string "参数错误"
*/
func createDeviceAssignmentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.DeviceAssignment
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.ID = nextAssignmentID
		req.AssignedAt = time.Now()
		deviceAssignmentsStore[req.ID] = &req
		nextAssignmentID++
		c.JSON(http.StatusCreated, req)
	}
}

/*
@Summary 获取设备绑定详情
@Description 根据ID查询设备绑定信息
@Tags DeviceAssignment
@Produce json
@Param id path int true "设备绑定ID"
@Success 200 {object} models.DeviceAssignment "查询成功"
@Failure 404 {object} map[string]string "未找到"
*/
func getDeviceAssignmentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		assignment, ok := deviceAssignmentsStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, assignment)
	}
}

/*
@Summary 设备绑定列表
@Description 获取所有设备绑定信息
@Tags DeviceAssignment
@Produce json
@Success 200 {array} models.DeviceAssignment "列表成功"
*/
func listDeviceAssignmentsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var assignments []*models.DeviceAssignment
		for _, a := range deviceAssignmentsStore {
			assignments = append(assignments, a)
		}
		c.JSON(http.StatusOK, assignments)
	}
}

/*
@Summary 解绑设备
@Description 解绑指定设备，设置UnassignedAt时间
@Tags DeviceAssignment
@Produce json
@Param id path int true "设备绑定ID"
@Success 200 {object} models.DeviceAssignment "解绑成功"
@Failure 404 {object} map[string]string "未找到"
*/
func unassignDeviceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		assignment, ok := deviceAssignmentsStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		now := time.Now()
		assignment.UnassignedAt = &now
		c.JSON(http.StatusOK, assignment)
	}
}

/*
@Summary 删除设备绑定
@Description 根据ID删除设备绑定关系
@Tags DeviceAssignment
@Produce json
@Param id path int true "设备绑定ID"
@Success 200 {object} map[string]interface{} "删除成功"
@Failure 404 {object} map[string]string "未找到"
*/
func deleteDeviceAssignmentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		_, ok := deviceAssignmentsStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		delete(deviceAssignmentsStore, id)
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}
