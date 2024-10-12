package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port                 int           `yaml:"port"`
	Host                 string        `yaml:"host"`
	Timeout              time.Duration `yaml:"timeout"`
	MaxUnaryConnections  int           `yaml:"max_unary_connections"`
	MaxStreamConnections int           `yaml:"max_stream_connections"`
}

// LoadPath loads configuration from specified path and returns config instance and error
func LoadPath(configPath string) (*Config, error) {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("config file does not exist: %s", err.Error())
	}

	return &cfg, nil
}

// MustLoad fetches path, loads configuration and panics on any error
func MustLoad() *Config {
	configPath := fetchConfigPath()

	cfg, err := LoadPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

// MustLoadPath loads configuration form configPath and panics on any error
func MustLoadPath(configPath string) *Config {
	cfg, err := LoadPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
