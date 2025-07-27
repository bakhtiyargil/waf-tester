package bootstrap

import (
	"waf-tester/config"
	"waf-tester/logger"
)

var (
	App *Application
)

type Application struct {
	Config *config.Config
	Logger logger.Logger
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = config.InitConfig("./config/config-local.yml")
	App.Logger = logger.InitAppLogger(App.Config)
}
