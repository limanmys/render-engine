package main

import (
	"renderer/src/connector"
	"renderer/src/helpers"
	"renderer/src/recycle"
	"renderer/src/sqlite"
	"renderer/src/web"
)

func main() {
	connector.ActiveConnections = make(map[string]connector.Connection)

	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	go recycle.Start()

	web.CreateWebServer()
}
