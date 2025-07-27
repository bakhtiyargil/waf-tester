package bootstrap

import (
	"waf-tester/config"
	"waf-tester/logger"
	"waf-tester/mongo"
)

const (
	CfgFilePath = "./config/config-local.yml"
)

var (
	App *Application
)

type Application struct {
	Config *config.Config
	Logger logger.Logger
	Mongo  mongo.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = InitConfig(CfgFilePath)
	App.Logger = InitAppLogger()
	App.Mongo = InitMongoDatabase()
}
