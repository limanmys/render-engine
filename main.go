package main

import (
	"github.com/limanmys/render-engine/connector"
	"github.com/limanmys/render-engine/helpers"
	"github.com/limanmys/render-engine/postgresql"
	"github.com/limanmys/render-engine/recycle"
	"github.com/limanmys/render-engine/web"
)

func main() {
	connector.ActiveConnections = make(map[string]connector.Connection)

	connector.ActiveTunnels = make(map[string]connector.Tunnel)

	helpers.ReadDataFromLiman()

	postgresql.InitDB()

	go recycle.Start()

	web.CreateWebServer()
}
