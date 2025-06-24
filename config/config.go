package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Default DefaultServerConfig
}

type DefaultServerConfig struct {
	Port           string        `yaml:"Port"`
	ReadTimeout    time.Duration `yaml:"ReadTimeout"`
	WriteTimeout   time.Duration `yaml:"WriteTimeout"`
	IdleTimeout    time.Duration `yaml:"IdleTimeout"`
	MaxHeaderBytes int           `yaml:"MaxHeaderBytes"`
}

func (c *Config) LoadConfig(file string) *Config {
	fileType := filepath.Ext(file)
	if fileType != ".yaml" && fileType != ".yml" {
		log.Fatalf("unsupported file type %s, file type must be .yml or .yaml", fileType)
	}

	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error while reading file %s, %v", file, err)
	}

	err = yaml.Unmarshal(configFile, c)
	if err != nil {
		log.Fatalf("error while parsing file %s, %v", file, err)
	}
	return c
}
