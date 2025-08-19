package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Logger   AppLogger
	Database DatabaseConfig
}

type ServerConfig struct {
	Default DefaultServerConfig
}

type AppLogger struct {
	Level string `yaml:"level"`
}

type DefaultServerConfig struct {
	Port           string        `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	IdleTimeout    time.Duration `yaml:"idleTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
	AllowOrigins   []string      `yaml:"allowOrigins"`
}

type DatabaseConfig struct {
	Mongo MongoConfig
}

type MongoConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
