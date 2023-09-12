package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token   string `yaml:"token"`
	AppID   string `yaml:"appID"`
	GuildID string `yaml:"guildID"`
}

func New() (*Config, error) {
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
