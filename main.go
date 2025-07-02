package main

import (
	"waf-tester/config"
	"waf-tester/server"
	"waf-tester/utility"
)

func main() {
	var c config.Config
	c.LoadConfig("./config/config-local.yml")

	wp := utility.NewWorkerPool(128)
	wp.Start()

	srv := server.NewServer(&c, server.NewHandler(wp))
	srv.Start()
}
