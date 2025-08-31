package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	MQTT struct {
		Broker   string
		ClientID string
		Username string
		Password string
	}
	WebSocket struct {
		Host string
		Port int
		Path string
	}
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
