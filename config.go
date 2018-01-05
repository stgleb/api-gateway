package api_gateway

import (
	"github.com/BurntSushi/toml"
	"time"
)

var config_instance *TomlConfig

type TomlConfig struct {
	Main             MainConfig
	Server           ServerConfig
	TokenService     TokenServiceConfig
	ServiceDiscovery ServiceDiscoveryConfig
}

type MainConfig struct {
	Title       string
	Release     string
	ServiceName string
	LogFile     string
}

type ServerConfig struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type TokenServiceConfig struct {
	ListenStr       string
	Protocol        string
	IssueTokenPath  string
	VerifyTokenPath string
	RevokeTokenPath string
}

type ServiceDiscoveryConfig struct {
	ConsulAddress     string
	ConsulPort        int
	AdvertisedAddress string
	AdvertisedPort    int
	Interval          string
	Timeout           string
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
