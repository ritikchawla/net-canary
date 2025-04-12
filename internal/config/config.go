package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Banner  string `yaml:"banner,omitempty"`
}

type Config struct {
	Services struct {
		SSH   ServiceConfig `yaml:"ssh"`
		HTTP  ServiceConfig `yaml:"http"`
		FTP   ServiceConfig `yaml:"ftp"`
		MySQL ServiceConfig `yaml:"mysql"`
		Redis ServiceConfig `yaml:"redis"`
	} `yaml:"services"`

	Logging struct {
		File     string `yaml:"file"`
		Remote   string `yaml:"remote,omitempty"`
		LogLevel string `yaml:"level"`
	} `yaml:"logging"`

	Alerts struct {
		Slack    string `yaml:"slack,omitempty"`
		Webhook  string `yaml:"webhook,omitempty"`
		Email    string `yaml:"email,omitempty"`
		Paranoia bool   `yaml:"paranoia"`
	} `yaml:"alerts"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
