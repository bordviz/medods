package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
	Database       `yaml:"database" env-required:"true"`
	HTTPServer     `yaml:"http_server" env-required:"true"`
	Auth           `yaml:"auth" env-required:"true"`
}

type Database struct {
	Host          string        `yaml:"host" env-required:"true"`
	Port          int           `yaml:"port" env-required:"true"`
	User          string        `yaml:"user" env-required:"true"`
	Password      string        `yaml:"password" env-required:"true"`
	Name          string        `yaml:"name" env-required:"true"`
	SSLMode       string        `yaml:"sslmode" env-required:"true"`
	MaxAttempts   int           `yaml:"max_attempts" env-required:"true"`
	AttempTimeout time.Duration `yaml:"attemp_timeout" env-required:"true"`
	AttempDelay   time.Duration `yaml:"attempt_delay" env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type Auth struct {
	AccessTokenLifetime  time.Duration `yaml:"access_token_lifetime" env-required:"true"`
	RefreshTokenLifetime time.Duration `yaml:"refresh_token_lifetime" env-required:"true"`
	AccessSecret         string        `yaml:"access_secret" env-required:"true"`
	RefreshSecret        string        `yaml:"refresh_secret" env-required:"true"`
}

func MustLoad(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config %s not found", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return &cfg, nil
}
