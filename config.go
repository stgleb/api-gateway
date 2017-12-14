package api_gateway

import (
	"github.com/BurntSushi/toml"
)

var config_instance *TomlConfig

type TomlConfig struct {
	Main MainConfig
}

type MainConfig struct {
	Title     string
	ListenStr string
	Release   string
}

func GetConfig() *TomlConfig {
	return config_instance
}

func NewConfig(file string) (*TomlConfig, error) {
	config_instance = &TomlConfig{}

	if _, err := toml.DecodeFile(file, config_instance); err != nil {
		return nil, err
	}

	return config_instance, nil
}
