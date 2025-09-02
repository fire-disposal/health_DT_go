package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port            int `mapstructure:"port"`
	MsgListenerPort int `mapstructure:"msglistener_port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"broker"`
	ClientID string `mapstructure:"client_id"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type WebSocketConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}

type WechatConfig struct {
	AppID  string `mapstructure:"appid"`
	Secret string `mapstructure:"secret"`
}

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Postgres  PostgresConfig  `mapstructure:"postgres"`
	Redis     RedisConfig     `mapstructure:"redis"`
	MQTT      MQTTConfig      `mapstructure:"mqtt"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	JWTSecret string          `mapstructure:"jwt_secret"`
	Wechat    WechatConfig    `mapstructure:"wechat"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	var c Config
	if err := v.ReadInConfig(); err == nil {
		if err := v.Unmarshal(&c); err != nil {
			return nil, fmt.Errorf("YAML 配置解析失败: %w", err)
		}
		return &c, nil
	}

	// YAML 未找到或解析失败，回退环境变量
	c = Config{
		Server: ServerConfig{
			Port:            getenvInt("PORT", 8002),
			MsgListenerPort: getenvInt("MSGLISTENER_PORT", 5858),
		},
		Postgres: PostgresConfig{
			Host:     getenv("POSTGRES_HOST", "localhost"),
			Port:     getenvInt("POSTGRES_PORT", 5432),
			User:     getenv("POSTGRES_USER", "postgres"),
			Password: getenv("POSTGRES_PASSWORD", ""),
			DBName:   getenv("POSTGRES_DBNAME", "health_dt"),
			SSLMode:  getenv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getenv("REDIS_ADDR", "localhost:6379"),
			Password: getenv("REDIS_PASSWORD", ""),
			DB:       getenvInt("REDIS_DB", 0),
		},
		MQTT: MQTTConfig{
			Broker:   getenv("MQTT_BROKER", "tcp://39.100.101.252:1883"),
			ClientID: getenv("MQTT_CLIENT_ID", "health_dt_client"),
			Username: getenv("MQTT_USERNAME", ""),
			Password: getenv("MQTT_PASSWORD", ""),
		},
		WebSocket: WebSocketConfig{
			Host: getenv("WS_HOST", "0.0.0.0"),
			Port: getenvInt("WS_PORT", 8765),
			Path: getenv("WS_PATH", "/ws/health"),
		},
		JWTSecret: getenv("JWT_SECRET", "your-secret-key"),
		Wechat: WechatConfig{
			AppID:  getenv("WECHAT_APPID", ""),
			Secret: getenv("WECHAT_SECRET", ""),
		},
	}
	return &c, nil
}

func getenv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getenvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}
