package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func New() (*Config, error) {
	filename := fmt.Sprintf("./conf/app.%s.yaml", os.Getenv("BACKEND_STAGE"))
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if loadGsmErr := cfg.loadFromGsm(); loadGsmErr != nil {
		return nil, loadGsmErr
	}

	return &cfg, nil
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}
