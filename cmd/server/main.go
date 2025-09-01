/*
@title HealthDT API
@version 1.0
@description 健康数据平台 API 文档
@host localhost:8080
@BasePath /
*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	healthapi "github.com/fire-disposal/health_DT_go/api/http" // 避免与标准库 http 冲突
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

// Application 应用程序结构体，统一管理所有组件
type Application struct {
	logger     *zap.Logger
	config     *config.Config
	db         *sql.DB
	router     *gin.Engine
	server     *http.Server
	pipeline   *app.Pipeline
	mqttClient *mqtt.MQTTClient
	msgpackSrv *msgpack.MsgpackServer

	// 用于优雅关闭的context
	ctx    context.Context
	cancel context.CancelFunc
}

func main() {
	// 初始化应用
	app, err := NewApplication()
	if err != nil {
		fmt.Printf("应用初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// 启动应用
	if err := app.Run(); err != nil {
		app.logger.Fatal("应用启动失败", zap.Error(err))
	}
}

// NewApplication 创建新的应用实例
func NewApplication() (*Application, error) {
	// 初始化logger
	logger := initLogger()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		logger.Error("配置加载失败", zap.Error(err))
		return nil, fmt.Errorf("配置加载失败: %w", err)
	}

	logger.Info("应用配置加载成功",
		zap.Int("server_port", cfg.Server.Port),
		zap.String("env", getEnv("ENV", "development")))

	// 初始化数据库
	db, err := initDB(cfg, logger)
	if err != nil {
		logger.Error("数据库初始化失败", zap.Error(err))
		return nil, fmt.Errorf("数据库初始化失败: %w", err)
	}

	// 初始化事件总线和数据处理管道
	eventBus := eventbus.NewEventBus()
	pipeline := app.NewPipeline(eventBus)
	registerHealthProcessors(pipeline, logger)

	// 创建应用实例
	ctx, cancel := context.WithCancel(context.Background())
	app := &Application{
		logger:   logger,
		config:   cfg,
		db:       db,
		pipeline: pipeline,
		ctx:      ctx,
		cancel:   cancel,
	}

	// 初始化HTTP路由
	app.initRouter()

	// 初始化HTTP服务器
	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      app.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return app, nil
}

// initLogger 初始化日志配置，优化编码和输出格式
func initLogger() *zap.Logger {
	// 自定义编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 大写级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // ISO8601时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 根据环境变量决定日志级别和格式
	var config zap.Config
	env := getEnv("ENV", "development")

	if env == "production" {
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json", // 生产环境使用JSON格式
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "console", // 开发环境使用控制台格式
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
		}
	}

	logger, err := config.Build(zap.AddCallerSkip(0))
	if err != nil {
		panic(fmt.Sprintf("无法初始化logger: %v", err))
	}

	return logger
}

// initRouter 初始化HTTP路由
func (app *Application) initRouter() {
	// 根据环境设置Gin模式
	if getEnv("ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 自定义中间件
	r.Use(ginLoggerMiddleware(app.logger))
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 健康检查端点
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"timestamp": time.Now().Unix(),
			"version":   "1.0",
		})
	})

	// API路由注册
	// 注册鉴权相关路由
	// 注册鉴权相关路由
	// RegisterAuthRoutes 来自 api/http/auth_routes.go，需加包前缀
	healthapi.RegisterAuthRoutes(r, app.db)

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.router = r
}

// Run 启动应用
func (app *Application) Run() error {
	app.logger.Info("正在启动应用服务器...")

	// 启动MQTT监听（异步）
	go app.startMQTT()

	// 启动Msgpack监听（异步）
	go app.startMsgpack()

	// 启动HTTP服务器（异步）
	go func() {
		app.logger.Info("HTTP服务器启动",
			zap.String("address", app.server.Addr),
			zap.String("mode", gin.Mode()))

		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("HTTP服务器启动失败", zap.Error(err))
			app.cancel() // 取消上下文，触发关闭
		}
	}()

	// 等待信号或上下文取消
	return app.waitForShutdown()
}

// waitForShutdown 等待关闭信号
func (app *Application) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		app.logger.Info("接收到关闭信号", zap.String("signal", sig.String()))
	case <-app.ctx.Done():
		app.logger.Info("应用上下文被取消")
	}

	app.logger.Info("开始优雅关闭...")
	return app.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (app *Application) gracefulShutdown() error {
	// 设置关闭超时
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := app.server.Shutdown(shutdownCtx); err != nil {
		app.logger.Error("HTTP服务器关闭失败", zap.Error(err))
		return err
	}

	app.logger.Info("应用已优雅关闭")
	return nil
}

// Close 关闭所有资源
func (app *Application) Close() {
	if app.cancel != nil {
		app.cancel()
	}

	if app.mqttClient != nil {
		app.mqttClient.Disconnect(1000)
	}

	if app.msgpackSrv != nil {
		// MsgpackServer 没有 Stop 方法，无法优雅关闭
	}

	if app.db != nil {
		app.db.Close()
	}

	if app.logger != nil {
		app.logger.Sync()
	}
}

// initDB 数据库初始化并检测可用性，优化连接配置
func initDB(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	pg := cfg.Postgres

	// 构建DSN，明确指定编码
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8 connect_timeout=10",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode,
	)

	logger.Debug("正在连接数据库",
		zap.String("host", pg.Host),
		zap.Int("port", pg.Port),
		zap.String("database", pg.DBName))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接创建失败: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(25)           // 最大打开连接数
	db.SetMaxIdleConns(5)            // 最大空闲连接数
	db.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	logger.Info("数据库连接成功")
	return db, nil
}

// registerHealthProcessors 注册健康数据处理器
func registerHealthProcessors(pipeline *app.Pipeline, logger *zap.Logger) {
	processors := []app.HealthDataProcessor{
		&health.HeartRateHandler{},
		&health.BloodPressureHandler{},
		&health.SpO2Handler{},
		&health.TemperatureHandler{},
	}

	for _, processor := range processors {
		pipeline.RegisterProcessor(processor)
	}

	logger.Info("健康数据处理器注册完成", zap.Int("count", len(processors)))
}

// startMQTT 启动MQTT监听
func (app *Application) startMQTT() {
	cfg := app.config.MQTT
	app.logger.Info("正在启动MQTT客户端...")

	app.mqttClient = mqtt.NewMQTTClient(mqtt.ClientConfig{
		Broker:   cfg.Broker,
		ClientID: cfg.ClientID,
		Username: cfg.Username,
		Password: cfg.Password,
	})

	if err := app.mqttClient.Connect(); err != nil {
		app.logger.Error("MQTT连接失败", zap.Error(err))
		return
	}

	if err := app.mqttClient.Subscribe("device/+/data/+", 0, handlers.HandleMQTTMessage(app.pipeline)); err != nil {
		app.logger.Error("MQTT订阅失败", zap.Error(err))
		return
	}

	app.logger.Info("MQTT客户端启动成功",
		zap.String("broker", cfg.Broker),
		zap.String("client_id", cfg.ClientID))
}

// startMsgpack 启动Msgpack监听
func (app *Application) startMsgpack() {
	port := app.config.Server.MsgListenerPort
	app.logger.Info("正在启动Msgpack服务器...")

	app.msgpackSrv = msgpack.NewMsgpackServer(
		handlers.HandleMsgpackPayload(app.pipeline),
		port,
	)

	if err := app.msgpackSrv.Start(); err != nil {
		app.logger.Error("Msgpack服务器启动失败", zap.Error(err))
		return
	}

	app.logger.Info("Msgpack服务器启动成功", zap.Int("port", port))
}

// ginLoggerMiddleware Gin日志中间件
func ginLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP请求",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
		)
		return ""
	})
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getEnv 获取环境变量，提供默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
