package util

import (
	"goio/logger"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

type ConfigManager struct {
	srvConf ServicesConfig
	conf    string
}

type IOServiceConf struct {
	Srvworker int
	Srvqueue  int

	Ioworker int
	Ioqueue  int

	Matrixbucket int
	Matrixsize   int
}

var (
	cm   *ConfigManager
	once sync.Once
)

func NewConfigManager(config string) *ConfigManager {
	once.Do(func() {
		cm = &ConfigManager{
			conf: config,
		}
		var conf ServicesConfig
		conf_path, _ := filepath.Abs(config)
		if _, err := toml.DecodeFile(conf_path, &conf); err != nil {
			logger.Error("DecodeFile error %v", err)
			return
		}
		cm.srvConf = conf
	})
	return cm
}

type SrvInfo map[string]serviceInfo

type HttpConfig struct {
	Localaddr  string
	Remoteaddr string
}

type ServicesConfig struct {
	Services SrvInfo
	Engine   IOServiceConf
	Http     HttpConfig
}

type serviceInfo struct {
	Addr string
	Name string
}

func GetServicesConfig() SrvInfo {
	return cm.srvConf.Services
}

func GetIOServiceConf() IOServiceConf {
	return cm.srvConf.Engine
}

func GetHttpConfig() HttpConfig {
	return cm.srvConf.Http
}

func GetConf() string {
	return cm.conf
}

func GetManager() *ConfigManager {
	return cm
}
