package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     int    `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
		TimeZone string `yaml:"timezone"`
	} `yaml:"database"`

	Application struct {
		Port          int    `yaml:"port"`
		Env           string `yaml:"env"`
		LogLevel      string `yaml:"log_level"`
		LogFile       string `yaml:"log_file"`
		LogMaxSize    int    `yaml:"log_max_size"`
		LogMaxBackups int    `yaml:"log_max_backups"`
		LogMaxAge     int    `yaml:"log_max_age"`
	} `yaml:"application"`
}

var AppConfig Config

func LoadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %v", err)
	}
}
