// 鉴权相关路由
package http

import (
	"database/sql"
	"net/http"

	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
	"github.com/fire-disposal/health_DT_go/internal/service"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 错误响应辅助函数
func errorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// RegisterAuthRoutes 注册鉴权相关路由
func RegisterAuthRoutes(r gin.IRouter, db *sql.DB) {
	authService := service.NewAuthService(postgres.NewAuthRepository(db))
	// @Summary 管理员登录
	// @Description 管理员账号密码登录
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param login body LoginRequest true "登录信息"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {string} string "认证失败"
	// @Router /api/admin/login [post]
	r.POST("/api/admin/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, "参数错误")
			return
		}
		result, err := authService.Login("admin", req.Username, req.Password)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": result.Token, "user_id": result.UserID, "role": result.Role})
	})

	// @Summary App用户登录
	// @Description App用户账号密码登录
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param login body LoginRequest true "登录信息"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {string} string "认证失败"
	// @Router /api/app/login [post]
	r.POST("/api/app/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, "参数错误")
			return
		}
		result, err := authService.Login("app", req.Username, req.Password)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": result.Token, "user_id": result.UserID, "role": result.Role})
	})

	// 注册微信登录路由
	RegisterWechatLoginRoute(r, authService)
}

type WechatLoginRequest struct {
	Code string `json:"code"`
}

// 新增微信登录路由
// POST /api/app/wechat_login
// @Summary 微信登录
// @Description 微信授权码登录
// @Tags auth
// @Accept json
// @Produce json
// @Param code body WechatLoginRequest true "微信授权码"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {string} string "认证失败"
// @Router /api/app/wechat_login [post]
func RegisterWechatLoginRoute(r gin.IRouter, authService *service.AuthService) {
	r.POST("/api/app/wechat_login", func(c *gin.Context) {
		var req WechatLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
			errorResponse(c, http.StatusBadRequest, "参数错误")
			return
		}
		result, err := authService.LoginWithWechatCode(req.Code)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": result.Token, "user_id": result.UserID, "role": result.Role})
	})
}
