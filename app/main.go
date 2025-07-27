package main

import (
	"waf-tester/bootstrap"
	"waf-tester/server"
)

func main() {
	srv := server.NewServer(bootstrap.App.Config, server.NewHandler(bootstrap.App.Logger))
	srv.Start()
}
