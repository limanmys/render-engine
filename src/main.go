package main

import (
	"renderer/src/helpers"
	"renderer/src/sqlite"
	"renderer/src/web"
)

func main() {
	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	web.CreateWebServer()

}
