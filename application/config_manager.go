package application

import (
	"goio/logger"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ConfigManager struct {
	ServicesConfig
	conf string
}

func NewConfigManager(config string) *ConfigManager {
	return &ConfigManager{
		conf: config,
	}
}

type SrvInfo map[string]serviceInfo

type ServicesConfig struct {
	Services SrvInfo
}

type serviceInfo struct {
	Addr string
	Name string
}

func (c *ConfigManager) GetServicesConfig() SrvInfo {
	var conf ServicesConfig
	conf_path, _ := filepath.Abs(c.conf)
	if _, err := toml.DecodeFile(conf_path, &conf); err != nil {
		logger.Error("DecodeFile error %v", err)
		return nil
	}
	return conf.Services
}
