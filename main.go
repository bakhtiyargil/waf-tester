package main

import (
	"waf-tester/config"
	"waf-tester/logger"
	"waf-tester/server"
	"waf-tester/utility"
)

func main() {
	var c config.Config
	cnf := c.LoadConfig("./config/config-local.yml")

	lgr := logger.NewAppLogger(cnf)
	lgr.InitLogger()

	wp := utility.NewWorkerPool(128)
	wp.Start()

	srv := server.NewServer(cnf, server.NewHandler(wp, lgr))
	srv.Start()
}
