package main

import (
	"waf-tester/bootstrap"
	"waf-tester/client"
	"waf-tester/server"
	"waf-tester/service"
)

func main() {
	tstr := service.NewInjectionTester(client.NewPureHttpClient(), bootstrap.App.Logger)
	hndlr := server.NewInjectionTestHandler(tstr, bootstrap.App.Logger)
	srv := server.NewServer(bootstrap.App.Config, hndlr, bootstrap.App.Logger)
	srv.Start()
}
