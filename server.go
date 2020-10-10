package main

import (
	"github.com/Keda87/echo-oauth2/services/application"
	"github.com/Keda87/echo-oauth2/services/config"
)

func main() {
	conf := config.GetConfig()

	app := application.New(conf)
	defer app.PreStopServer()

	app.StartServer()
}
