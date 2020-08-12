package helpers

// ReadDataFromLiman Retrieve data from liman
func ReadDataFromLiman() {
	readAppKey()
}

func readAppKey() {
	output, err := ExecuteCommand("cat /liman/server/.env | grep APP_KEY")
	if err != nil {
		Abort(err.Error())
	}
	AppKey = StringAfter(output, "APP_KEY=")
}
