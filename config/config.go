package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
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

func InitConfig(file string) (c *Config) {
	fileType := filepath.Ext(file)
	if fileType != ".yaml" && fileType != ".yml" {
		log.Fatalf("unsupported file type %s, file type must be .yml or .yaml", fileType)
	}

	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error while reading file %s, %v", file, err)
	}

	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("error while parsing file %s, %v", file, err)
	}
	return c
}
