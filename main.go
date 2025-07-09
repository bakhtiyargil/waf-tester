package main

import (
	"waf-tester/config"
	"waf-tester/server"
)

func main() {
	var c config.Config
	c.LoadConfig("./config/config-local.yml")

	srv := server.NewServer(&c, server.NewHandler())
	srv.Start()
}
