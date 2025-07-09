package main

import (
	"waf-tester/config"
	"waf-tester/logger"
	"waf-tester/server"
)

func main() {
	var c config.Config
	cnf := c.LoadConfig("./config/config-local.yml")

	lgr := logger.NewAppLogger(cnf)
	lgr.InitLogger()

	srv := server.NewServer(cnf, server.NewHandler(lgr))
	srv.Start()
}
