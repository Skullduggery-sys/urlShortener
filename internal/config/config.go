package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"os"
	"time"
	"urlShortener/utils"
	"urlShortener/utils/e"
)

type Config struct {
	Postgres   PostgresConfig   `yaml:"postgres"`
	HTTPServer HTTPServerConfig `yaml:"httpServer"`
	GRPCAddr   string           `yaml:"grpcAddr" validate:"required"`
}

type PostgresConfig struct {
	Login    string `yaml:"login" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required,numeric"`
	DBName   string `yaml:"dbname" validate:"required"`
	SSLMode  string `yaml:"sslMode" validate:"required"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" validate:"required"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

func MustParseConfig(configPath string) (*Config, error) {
	const fn = "internal.config.MustParseConfig"

	if !fileExists(configPath) {
		return nil, e.WrapError(fn, os.ErrNotExist)
	}

	cfg, err := readConfig(configPath)
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	if err = validator.New().Struct(cfg); err != nil {
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, e.WrapError(fn, err)
		}
		return nil, e.WrapError(fn, utils.ValidateErrors(err))
	}
	return cfg, nil
}

func readConfig(configPath string) (*Config, error) {
	const fn = "internal.config.readConfig"

	viper.SetConfigFile(configPath)
	viper.SetDefault("httpServer.timeout", time.Second*10)
	viper.SetDefault("httpServer.idleTimeout", time.Minute)

	if err := viper.ReadInConfig(); err != nil {
		return nil, e.WrapError(fn, err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, e.WrapError(fn, err)
	}

	return &cfg, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
