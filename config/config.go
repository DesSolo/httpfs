package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type callBack struct {
	Kind     string `yaml:"kind"`
	FileSize struct {
		Min int64 `yaml:"min"`
		Max int64 `yaml:"max"`
	} `yaml:"file_size"`
	Exec struct {
		Command string   `yaml:"command"`
		Args    []string `yaml:"args"`
	}
}

type Config struct {
	Address string `yaml:"address"`
	Logging struct {
		Handler string `yaml:"handler"`
		Option  struct {
			AddSource bool `yaml:"add_source"`
			Level     int  `yaml:"level"`
		} `yaml:"options"`
	} `yaml:"logging"`
	Storage struct {
		Kind string `yaml:"kind"`
		File struct {
			BasePath      string `yaml:"base_path"`
			TemporaryPath string `yaml:"temporary_path"`
		} `yaml:"file"`
	} `yaml:"storage"`
	Callback struct {
		Pre  []callBack `yaml:"pre"`
		Post []callBack `yaml:"post"`
	} `yaml:"callback"`
}

func NewConfigFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file err: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal err: %w", err)
	}

	return &cfg, nil
}
