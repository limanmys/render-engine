package main

import (
	"github.com/limanmys/go/connector"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/postgresql"
	"github.com/limanmys/go/recycle"
	"github.com/limanmys/go/web"
)

func main() {
	connector.ActiveConnections = make(map[string]connector.Connection)

	connector.ActiveTunnels = make(map[string]connector.Tunnel)

	helpers.ReadDataFromLiman()

	postgresql.InitDB()

	go recycle.Start()

	web.CreateWebServer()
}
