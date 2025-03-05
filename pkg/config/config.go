package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Spec ConfigSpec `yaml:"spec"`
}

type ConfigSpec struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Path     string   `yaml:"path"`
	Commands []string `yaml:"commands"`
}

func Load(configPath string) (*Config, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
