/*
@title HealthDT API
@version 1.0
@description 健康数据平台 API 文档
@host localhost:8080
@BasePath /
*/

// main.go - Gin 启动示例
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/fire-disposal/health_DT_go/api/http"
	"github.com/fire-disposal/health_DT_go/config"
	"github.com/fire-disposal/health_DT_go/internal/app"
	"github.com/fire-disposal/health_DT_go/internal/app/eventbus"
	"github.com/fire-disposal/health_DT_go/internal/app/handlers"
	"github.com/fire-disposal/health_DT_go/internal/app/handlers/health"
	"github.com/fire-disposal/health_DT_go/internal/mqtt"
	"github.com/fire-disposal/health_DT_go/internal/msgpack"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	fmt.Printf("调试: cfg.Server.Port = '%v'\n", cfg.Server.Port)

	// 初始化数据库连接
	pg := cfg.Postgres
	dsn :=
		"host=" + pg.Host +
			" port=" + fmt.Sprintf("%d", pg.Port) +
			" user=" + pg.User +
			" password=" + pg.Password +
			" dbname=" + pg.DBName +
			" sslmode=" + pg.SSLMode
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 初始化事件总线与主流程
	eventBus := eventbus.NewEventBus()
	pipeline := app.NewPipeline(eventBus)

	// 注册健康数据处理器（事件驱动）
	pipeline.RegisterProcessor(&health.HeartRateHandler{})
	pipeline.RegisterProcessor(&health.BloodPressureHandler{})
	pipeline.RegisterProcessor(&health.SpO2Handler{})
	pipeline.RegisterProcessor(&health.TemperatureHandler{})

	// 启动 MQTT 客户端并监听
	go func() {
		mqttCfg := cfg.MQTT
		mqttClient := mqtt.NewMQTTClient(mqtt.ClientConfig{
			Broker:   mqttCfg.Broker,
			ClientID: mqttCfg.ClientID,
			Username: mqttCfg.Username,
			Password: mqttCfg.Password,
		})
		if err := mqttClient.Connect(); err != nil {
			log.Printf("MQTT连接失败: %v", err)
			return
		}
		if err := mqttClient.Subscribe("device/+/data/+", 0, handlers.HandleMQTTMessage(pipeline)); err != nil {
			log.Printf("MQTT订阅失败: %v", err)
			return
		}
		log.Printf("MQTT监听已启动，Broker: %s，ClientID: %s", mqttCfg.Broker, mqttCfg.ClientID)
	}()

	// 启动 Msgpack TCP 监听
	go func() {
		server := msgpack.NewMsgpackServer(handlers.HandleMsgpackPayload(pipeline))
		if err := server.Start(); err != nil {
			log.Printf("Msgpack监听启动失败: %v", err)
		}
		log.Printf("Msgpack监听已启动，端口: %d", 5858)
	}()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 注册鉴权相关路由
	http.RegisterAuthRoutes(r, db)

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.Server.Port) // 监听配置端口
}
