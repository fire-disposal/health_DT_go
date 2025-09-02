// Package http 管理员用户CRUD路由
package http

import (
	"net/http"
	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/gin-gonic/gin"
)

var adminUsersStore = make(map[int64]*models.AdminUser)
var nextAdminID int64 = 1

func RegisterAdminUsersRoutes(router gin.IRouter) {
	group := router.Group("/admin_users")
	{
		group.POST("", createAdminUserHandler())
		group.GET("/:id", getAdminUserHandler())
		group.GET("", listAdminUsersHandler())
		group.PUT("/:id", updateAdminUserHandler())
		group.DELETE("/:id", deleteAdminUserHandler())
	}
}

/*
@Summary 创建管理员用户
@Description 新增管理员用户，需提交完整信息
@Tags AdminUser
@Accept json
@Produce json
@Param body body models.AdminUser true "管理员用户信息"
@Success 201 {object} models.AdminUser "创建成功"
@Failure 400 {object} map[string]string "参数错误"
*/
func createAdminUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AdminUser
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.ID = nextAdminID
		nextAdminID++
		adminUsersStore[req.ID] = &req
		c.JSON(http.StatusCreated, req)
	}
}

/*
@Summary 获取管理员用户详情
@Description 根据ID查询管理员用户信息
@Tags AdminUser
@Produce json
@Param id path int true "管理员用户ID"
@Success 200 {object} models.AdminUser "查询成功"
@Failure 404 {object} map[string]string "未找到"
*/
func getAdminUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		user, ok := adminUsersStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

/*
@Summary 管理员用户列表
@Description 获取所有管理员用户信息
@Tags AdminUser
@Produce json
@Success 200 {array} models.AdminUser "列表成功"
*/
func listAdminUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []*models.AdminUser
		for _, u := range adminUsersStore {
			users = append(users, u)
		}
		c.JSON(http.StatusOK, users)
	}
}

/*
@Summary 更新管理员用户
@Description 根据ID更新管理员用户信息
@Tags AdminUser
@Accept json
@Produce json
@Param id path int true "管理员用户ID"
@Param body body models.AdminUser true "管理员用户信息"
@Success 200 {object} models.AdminUser "更新成功"
@Failure 400 {object} map[string]string "参数错误"
@Failure 404 {object} map[string]string "未找到"
*/
func updateAdminUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var req models.AdminUser
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, ok := adminUsersStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		req.ID = id
		adminUsersStore[id] = &req
		c.JSON(http.StatusOK, req)
	}
}

/*
@Summary 删除管理员用户
@Description 根据ID删除管理员用户
@Tags AdminUser
@Produce json
@Param id path int true "管理员用户ID"
@Success 200 {object} map[string]interface{} "删除成功"
@Failure 404 {object} map[string]string "未找到"
*/
func deleteAdminUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		_, ok := adminUsersStore[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		delete(adminUsersStore, id)
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "deleted"})
	}
}
