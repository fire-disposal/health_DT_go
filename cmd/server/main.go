// main.go - Gin 启动示例
package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/fire-disposal/health_DT_go/api/http"
)

func main() {
	// 初始化数据库连接
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=health_dt sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 注册登录路由
	r.POST("/api/admin/login", func(c *gin.Context) {
		http.AdminLoginHandlerWithDB(db, c.Writer, c.Request)
	})
	r.POST("/api/app/login", func(c *gin.Context) {
		http.AppLoginHandlerWithDB(db, c.Writer, c.Request)
	})

	r.Run() // 默认监听 :8080
}
