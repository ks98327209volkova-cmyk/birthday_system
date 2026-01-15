package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis   RedisConfig   `yaml:"redis"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Service ServiceConfig `yaml:"service"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

type ServiceConfig struct {
	ReminderDays []int `yaml:"reminder_days"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return &config, nil
}
