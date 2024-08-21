package main

import (
	"log"
	_ "net/http/pprof"
	_ "userservice/docs"
	"userservice/internal/app"
)

// @title UserService
// @version 1.0.0
// @description Geoservice API

// @host localhost:8888
// @basePath /
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	a.Shutdown()
}
