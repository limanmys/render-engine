package main

import (
	"github.com/limanmys/go/connector"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/recycle"
	"github.com/limanmys/go/sqlite"
	"github.com/limanmys/go/web"
)

func main() {
	connector.ActiveConnections = make(map[string]connector.Connection)

	connector.ActiveTunnels = make(map[string]connector.Tunnel)

	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	go recycle.Start()

	web.CreateWebServer()
}
