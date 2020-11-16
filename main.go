package main

import (
	"fmt"
	"github.com/limanmys/go/postgresql"
)

func main() {
	postgresql.InitDB()

	fmt.Println(postgresql.GetUserData("6723ba96-c7f7-44c5-b02e-20eee74f2f4d","8bcb8fbb-e058-457a-9e7d-057900bbe396","83fb773a-61fb-41d4-9258-0e3247f2660f"))

	//connector.ActiveConnections = make(map[string]connector.Connection)
	//
	//connector.ActiveTunnels = make(map[string]connector.Tunnel)
	//
	//helpers.ReadDataFromLiman()
	//
	//sqlite.InitDB()
	//
	//go recycle.Start()
	//
	//web.CreateWebServer()
}
