package application

import (
	"goio/logger"
	"goio/service"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ConfigManager struct {
	ServicesConfig
	conf string
}

func NewConfigManager(config string) *ConfigManager {
	cm := &ConfigManager{
		conf: config,
	}
	var conf ServicesConfig
	conf_path, _ := filepath.Abs(config)
	if _, err := toml.DecodeFile(conf_path, &conf); err != nil {
		logger.Error("DecodeFile error %v", err)
		return nil
	}
	cm.ServicesConfig = conf
	return cm
}

type SrvInfo map[string]serviceInfo

type ServicesConfig struct {
	Services SrvInfo
	Engine   service.IOServiceConfig
}

type serviceInfo struct {
	Addr string
	Name string
}

func (c *ConfigManager) GetServicesConfig() SrvInfo {
	return c.ServicesConfig.Services
}

func (c *ConfigManager) GetIOServiceConfig() service.IOServiceConfig {
	return c.ServicesConfig.Engine
}
