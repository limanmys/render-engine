package helpers

// ReadDataFromLiman Retrieve data from liman
func ReadDataFromLiman() {
	readAppKey()
	readPGData()
}

func readAppKey() {
	output, err := ExecuteCommand("cat /liman/server/.env | grep APP_KEY")
	if err != nil {
		Abort(err.Error())
	}
	AppKey = StringAfter(output, "APP_KEY=")
}

func readPGData() {
	output, err := ExecuteCommand("cat /liman/server/.env | grep DB_HOST")
	if err != nil {
		Abort(err.Error())
	}
	DBHost = StringAfter(output, "DB_HOST=")

	output, err = ExecuteCommand("cat /liman/server/.env | grep DB_PORT")
	if err != nil {
		Abort(err.Error())
	}
	DBPort = StringAfter(output, "DB_PORT=")

	output, err = ExecuteCommand("cat /liman/server/.env | grep DB_DATABASE")
	if err != nil {
		Abort(err.Error())
	}
	DBName = StringAfter(output, "DB_DATABASE=")

	output, err = ExecuteCommand("cat /liman/server/.env | grep DB_USERNAME")
	if err != nil {
		Abort(err.Error())
	}
	DBUsername = StringAfter(output, "DB_USERNAME=")

	output, err = ExecuteCommand("cat /liman/server/.env | grep DB_PASSWORD")
	if err != nil {
		Abort(err.Error())
	}
	DBPassword = StringAfter(output, "DB_PASSWORD=")
}
