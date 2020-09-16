package main

import (
	"renderer/src/constants"
	"renderer/src/helpers"
	"renderer/src/sqlite"
	"renderer/src/web"
)

func main() {
	constants.ActiveConnections = make(map[string]constants.Connection)

	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	web.CreateWebServer()

}
