package env

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server   `yaml:"server"`
	Keys   Keys     `yaml:"keys"`
	TTL    TTL      `yaml:"ttl"`
	DB     Database `yaml:"database"`
}

type Database struct {
	Path string `yaml:"path"`
}

type Keys struct {
	SecretKeyForJWT     string `yaml:"secretKeyForJWT"`
	SecretKeyForRefresh string `yaml:"secretKeyForRefresh"`
	SecretForAES        string `yaml:"secretForAES"`
}

type TTL struct {
	AccessToken  time.Duration `yaml:"accessToken"`
	RefreshToken time.Duration `yaml:"refreshToken"`
	Cookie       time.Duration `yaml:"cookie"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func GetConfig() Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	configPath := filepath.Join(wd, "/config/config.yaml")
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return config
}
