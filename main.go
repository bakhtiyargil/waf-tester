package main

import (
	"log"
	"waf-tester/config"
	"waf-tester/server"
)

func main() {
	var c config.Config
	c.LoadConfig("./config/config-local.yml")

	srv := server.NewServer(&c)
	err := srv.Start()
	if err != nil {
		log.Fatalf("server start error: %v", err)
	}
}
