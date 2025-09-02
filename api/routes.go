// 路由与Swagger统一管理
package api

import (
	"database/sql"

	healthapi "github.com/fire-disposal/health_DT_go/api/http"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 挂载所有业务路由和Swagger UI
func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// 统一API前缀
	apiV1 := r.Group("/api/v1")

	// 挂载各模块路由
	healthapi.RegisterAuthRoutes(apiV1, db)
	healthapi.RegisterAdminUsersRoutes(apiV1)
	healthapi.RegisterDevicesRoutes(apiV1, nil)
	healthapi.RegisterDeviceAssignmentsRoutes(apiV1)
	healthapi.RegisterHealthProfilesRoutes(apiV1, nil)
	healthapi.RegisterHealthDataRoutes(apiV1, db)

	// Swagger UI 挂载到 /api/v1/swagger
	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
