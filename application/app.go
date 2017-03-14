package application

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"sync"

	"goio/logger"
	"goio/msg"
	"goio/network"
	"goio/service"
)

var (
	help_mode    = flag.Bool("h", false, "help mode")
	version_mode = flag.Bool("v", false, "version mode")
	daemon_mode  = flag.Bool("d", false, "daemon mode")
	app_config   = flag.String("c", ".", "app config")
	log_config   = flag.String("l", ".", "log config")
)

type DelegateFunc func() error
type Application interface {
	OnStart() error
	OnStop() error
	OnFinish() error
}
type GenericApplication struct {
	io_service     *service.IOService
	config_manager *ConfigManager
	logger_config  string
	acceptors      []*network.Acceptor
	on_start       DelegateFunc
	on_stop        DelegateFunc
	on_finish      DelegateFunc
}

func (app *GenericApplication) SetOnStart(fn DelegateFunc) *GenericApplication {
	app.on_start = fn
	return app
}

func (app *GenericApplication) SetOnStop(fn DelegateFunc) *GenericApplication {
	app.on_stop = fn
	return app
}

func (app *GenericApplication) SetOnFinish(fn DelegateFunc) *GenericApplication {
	app.on_finish = fn
	return app
}

func (app *GenericApplication) OnStart() error {
	if app.on_start != nil {
		return app.on_start()
	}
	return nil
}

func (app *GenericApplication) OnStop() error {
	if app.on_stop != nil {
		return app.on_stop()
	}
	return nil
}

func (app *GenericApplication) OnFinish() error {
	if app.on_finish != nil {
		return app.on_finish()
	}
	return nil
}

func (app *GenericApplication) RegisterService(name string, fn func(msg.Message)) {
	service.Instance().RegisterService(name, fn)
}

func (app *GenericApplication) Run() (err error) {
	if err = app.Start(); err != nil {
		logger.Info("Run error %v", err)
		return err
	}

	/*if help_mode || version_mode {
		return nil
	}
	*/
	if err = app.Wait(); err != nil {
		return nil
	}

	app.Stop()

	return nil
}

func (app *GenericApplication) Start() (err error) {
	/*	if err = ParseCmd(); err != nil {
			return nil
		}

		if help_mode || version_mode {
			return nil
		}
	*/
	if err = app.InitApp(); err != nil {
		return nil
	}

	if err = app.WritePidFile(); err != nil {
		return err
	}

	if err = app.StartLogger(); err != nil {
		return err
	}

	if err = app.LoadConfig(); err != nil {
		return err
	}

	if err = app.StartIOService(); err != nil {
		return err
	}

	if err = app.OnStart(); err != nil {
		return err
	}

	if err = app.StartServiceManager(); err != nil {
		return err
	}

	return nil
}

func (app *GenericApplication) Wait() (err error) {
	if err = app.RunIOService(); err != nil {
		return err
	}
	return nil
}

func (app *GenericApplication) Stop() (err error) {
	logger.Info("Application:OnStop")
	if err = app.OnStop(); err != nil {
		return err
	}

	logger.Info("Application:StopServiceManager")
	if err = app.StopServiceManager(); err != nil {
		return err
	}

	//time.Sleep(5 * time.Second)
	if err = app.StopIOService(); err != nil {
		return err
	}

	if err = app.OnFinish(); err != nil {
		return err
	}

	return nil
}

func (app *GenericApplication) ParseCmd() {

}

func (app *GenericApplication) InitApp() (err error) {
	/*if *daemon_mode {
		if err = util.Daemon(); err != nil {
			return
		}
	}
	*/
	//util.Daemon()
	return nil
}

func (app *GenericApplication) Daemon() error {
	return nil
}

func (app *GenericApplication) WritePidFile() (err error) {
	var f *os.File
	defer f.Close()
	f, err = os.OpenFile("pid", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	f.WriteString(strconv.Itoa(os.Getpid()))
	return
}

func (app *GenericApplication) StartLogger() (err error) {
	if app.logger_config == "" {
		logger.Start(logger.LogFilePath("./log"), logger.EveryHour, logger.PrintStack, logger.AlsoStdout)
	} else {

	}
	return nil
}

func (app *GenericApplication) LoadConfig() (err error) {
	/*if app.app_config == "" {

	} else {

	}*/
	app.config_manager = NewConfigManager("./conf.toml")
	logger.Info("app conf path %s", app.config_manager.conf)
	return nil
}

func (app *GenericApplication) StartIOService() (err error) {
	/*if err = app.io_service.Init(app.io_service_config); err != nil {
		return err
	}
	*/
	/*if err = app.io_service.Start(); err != nil {
		return err
	}
	*/
	app.io_service = app.GetIOService()
	app.io_service.HandleSignal()
	app.io_service.Start()
	return nil
}

var (
	io_srv *service.IOService
	once   sync.Once
)

func (app *GenericApplication) GetIOService() *service.IOService {
	once.Do(func() {
		io_srv = &service.IOService{}
	})
	return io_srv
}

func (app *GenericApplication) StartServiceManager() (err error) {
	if app.config_manager == nil {
		return errors.New("Application.StartServiceManager no config set,ignore")
	}

	var acceptor *network.Acceptor
	logger.Info("service config %v", app.config_manager.GetServicesConfig())
	for _, service_config := range app.config_manager.GetServicesConfig() {
		acceptor, err = network.Instance().CreateAcceptor(app.GetIOService(),
			service_config.Addr, service_config.Name)

		if err != nil {
			logger.Error("Application.StartServiceManager IODescriptorFactory.CreateAcceptor fail,address:%s,name:%s",
				service_config.Addr, service_config.Name)
			return err
		}

		go acceptor.Start()
		//logger.Error("Application.StartServiceManager acceptor.Start fail, address:%s,name:%s",
		//service_config.addr, service_config.name)

		//logger.Info("Application.StartServiceManager , address:%s,name:%s",
		//service_config.addr, service_config.name)

		app.acceptors = append(app.acceptors, acceptor)
	}

	return nil
}

func (app *GenericApplication) RunIOService() error {
	app.io_service.Run()
	return nil
}

func (app *GenericApplication) StopServiceManager() error {
	for _, acceptor := range app.acceptors {
		acceptor.Stop()
	}
	return nil
}

func (app *GenericApplication) StopIOService() error {
	app.io_service.CleanUp()
	return nil
}
