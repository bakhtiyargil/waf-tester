package main

import (
	"time"
	"waf-tester/bootstrap"
	"waf-tester/client"
	"waf-tester/server"
	"waf-tester/service"
	"waf-tester/test/repository/mongo"
	"waf-tester/test/usecase"
)

func main() {
	clt := client.NewPureHttpClient()
	rep := mongo.NewMongoRepository(bootstrap.App.Mongo.Database("waftester"))
	dur, _ := time.ParseDuration("5m")
	uc := usecase.NewCatUseCase(rep, dur)
	tstr := service.NewInjectionTester(clt, bootstrap.App.Logger, uc)
	hndlr := server.NewInjectionTestHandler(tstr, bootstrap.App.Logger)
	srv := server.NewServer(hndlr, bootstrap.App.Config, bootstrap.App.Logger)
	srv.Start()
}
