package main

import (
	"renderer/src/helpers"
	"renderer/src/sqlite"
	"renderer/src/web"
)

func main() {
	helpers.ReadDataFromLiman()

	sqlite.InitDB()

	// token := "4iQK5HnGWDgRLezgGjhve1fbc8X3T6Id"

	// targetFunction := "index"

	// extensionID := "64190c43-5fd6-4536-89c2-a9db82fc901b"

	// serverID := "74917cca-f478-4051-b9ea-9401d0928f43"

	// requestData := make(map[string]string)

	// userID := sqlite.GetUserIDFromToken(token)

	// sandbox.GeneratePHPCommand(targetFunction, userID, extensionID, serverID, requestData, token, false)

	web.CreateWebServer()

}
