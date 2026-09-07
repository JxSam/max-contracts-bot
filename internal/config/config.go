package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type BotConfig struct {
	Token string `yaml:"token" env-required:"true"`
	Api   string `yaml:"api" env-required:"true"`
	Debug bool   `yaml:"debug" env-default:"false"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

type SqlConfig struct {
	Path string `yaml:"path" env-required:"true"`
}

type MessageConfig struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
	Url     string `yaml:"url"`
}

type Config struct {
	Env            string         `yaml:"env" env-default:"local"`
	BotConfig      BotConfig      `yaml:"bot"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	SqlConfig      SqlConfig      `yaml:"sql"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	return LoadPath(configPath)
}

func LoadPath(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %v", err.Error())
	}

	return &cfg, nil
}
