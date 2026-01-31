package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MonitorConfig struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	RulesPath   string `yaml:"rules_path"`
	SourceColor string `yaml:"source_color"`
}

type AppConfig struct {
	Monitors []MonitorConfig `yaml:"monitors"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
