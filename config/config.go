package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Anthropic Anthropic `yaml:"anthropic"`
	YouTrack  YouTrack  `yaml:"youtrack"`
	GitHub    GitHub    `yaml:"github"`
}

type Anthropic struct {
	APIKey string `yaml:"api_key"`
}

type YouTrack struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

type GitHub struct {
	Token       string `yaml:"token"`
	DefaultRepo string `yaml:"default_repo"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".worklog")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	var cfg Config
	data, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
