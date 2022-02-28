package helpers

import (
	"fmt"
	"strings"
)
// ReadDataFromLiman Retrieve data from liman
func ReadDataFromLiman() {
	readAppKey()
	readPGData()
	readTimeout()
}

func readAppKey() {
	output, err := ExecuteCommand("cat /liman/server/.env | grep APP_KEY")
	if err != nil {
		fmt.Println("Liman Sifreleme Anahtari Okunamadi!")
		Abort(err.Error())
	}
	AppKey = StringAfter(output, "APP_KEY=")
}

func readTimeout() {
	output, err := ExecuteCommand("cat /liman/server/.env | grep EXTENSION_TIMEOUT")
	if err != nil {
		Timeout = "30"
		return
	}
	if strings.TrimSpace(output) == "" {
		Timeout = "30"
		return
	}
	Timeout = StringAfter(output, "EXTENSION_TIMEOUT=")
}

func readPGData() {
	output, _ := ExecuteCommand("cat /liman/server/.env | grep DB_HOST")

	DBHost = StringAfter(output, "DB_HOST=")

	if DBHost == "" {
		DBHost = "127.0.0.1"
	}

	output, _ = ExecuteCommand("cat /liman/server/.env | grep DB_PORT")

	DBPort = StringAfter(output, "DB_PORT=")

	if DBPort == "" {
		DBPort = "5432"
	}

	output, _ = ExecuteCommand("cat /liman/server/.env | grep DB_DATABASE")

	DBName = StringAfter(output, "DB_DATABASE=")

	if DBName == "" {
		DBName = "liman"
	}

	output, _ = ExecuteCommand("cat /liman/server/.env | grep DB_USERNAME")

	DBUsername = StringAfter(output, "DB_USERNAME=")

	if DBUsername == "" {
		DBUsername = "liman"
	}

	output, _ = ExecuteCommand("cat /liman/server/.env | grep DB_PASSWORD")

	DBPassword = StringAfter(output, "DB_PASSWORD=")
}
