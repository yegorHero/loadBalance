package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	BackendAddresses []string `yaml:"backend_addresses"`
	AlgorithmType    string   `yaml:"algorithm_type"`
	Bucket           Bucket   `yaml:"bucket"`
}

type Bucket struct {
	Rate     time.Duration `yaml:"rate"`
	Capacity int           `yaml:"capacity"`
}

func Init(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	return &config
}
