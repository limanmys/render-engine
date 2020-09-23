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

	connector.ActiveTunnels = make(map[string]connector.Tunnel)

	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	go recycle.Start()

	web.CreateWebServer()
}
