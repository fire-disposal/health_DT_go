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

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

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
	// 初始化 zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("无法初始化 zap logger: %v", err))
	}
	defer logger.Sync()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("配置加载失败", zap.Error(err))
	}
	logger.Info("调试", zap.String("Server.Port", cfg.Server.Port))

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
		logger.Fatal("数据库连接失败", zap.Error(err))
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
			logger.Error("MQTT连接失败", zap.Error(err))
			return
		}
		if err := mqttClient.Subscribe("device/+/data/+", 0, handlers.HandleMQTTMessage(pipeline)); err != nil {
			logger.Error("MQTT订阅失败", zap.Error(err))
			return
		}
		logger.Info("MQTT监听已启动", zap.String("Broker", mqttCfg.Broker), zap.String("ClientID", mqttCfg.ClientID))
	}()

	// 启动 Msgpack TCP 监听
	go func() {
		server := msgpack.NewMsgpackServer(handlers.HandleMsgpackPayload(pipeline))
		if err := server.Start(); err != nil {
			logger.Error("Msgpack监听启动失败", zap.Error(err))
		}
		logger.Info("Msgpack监听已启动", zap.Int("端口", 5858))
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
