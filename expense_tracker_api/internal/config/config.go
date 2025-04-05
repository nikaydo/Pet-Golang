package env

import (
	"log"
	"os"
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
	file, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling: %v", err)
	}
	return config
}
