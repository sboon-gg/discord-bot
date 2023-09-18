package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DbConfig struct {
	Filename string `yaml:"filename"`
}

type Config struct {
	Token   string   `yaml:"token"`
	AppID   string   `yaml:"appID"`
	GuildID string   `yaml:"guildID"`
	Db      DbConfig `yaml:"db"`
}

func New(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
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
