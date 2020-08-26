package sandbox

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"renderer/src/sqlite"
	"strings"

	"github.com/mervick/aes-everywhere/go/aes256"
)

// GeneratePHPCommand generate command
func GeneratePHPCommand(targetFunction string, userID string, extensionID string, serverID string, requestData map[string]string, token string, baseURL string, locale string, logObject RegularLog) string {
	result := make(map[string]string)
	combinerPath := "/liman/sandbox/php/index.php"
	server, extension, settings := sqlite.GetUserData(serverID, extensionID, userID)
	user := sqlite.GetUser(userID)

	b, _ := json.Marshal(server)
	result["server"] = string(b)

	b, _ = json.Marshal(extension)
	result["extension"] = string(b)

	b, _ = json.Marshal(settings)
	result["settings"] = string(b)

	b, _ = json.Marshal(user)
	result["user"] = string(b)

	b, _ = json.Marshal(result)

	result["functionsPath"] = "/liman/extensions/" + strings.ToLower(extension.Name) + "/views/functions.php"

	result["function"] = targetFunction

	b, _ = json.Marshal(requestData)
	result["requestData"] = string(b)

	license, _ := sqlite.GetLicense(extension.ID)

	result["license"] = license.Data

	result["apiRoute"] = "/extensionRun"

	result["navigationRoute"] = "/l/" + extension.ID + "/" + server.City + "/" + server.ID

	result["token"] = token

	result["locale"] = locale

	result["log_id"] = logObject.LogID

	result["ajax"] = "true"

	result["publicPath"] = baseURL + "/eklenti/" + extension.ID + "/public/"

	b, _ = json.Marshal(sqlite.GetPermissions(userID))
	result["permissions"] = string(b)

	soPath := "/liman/extensions/" + strings.ToLower(extension.Name) + "/liman.so"
	soCommand := ""
	if _, err := os.Stat(soPath); err == nil {
		soCommand = "-dextension=" + soPath + " "
	}

	keyPath := "/liman/keys/" + extension.ID
	content, _ := ioutil.ReadFile(keyPath)

	b, _ = json.Marshal(result)

	encryptedData := aes256.Encrypt(string(b), string(content))

	command := "sudo runuser " + strings.Replace(extension.ID, "-", "", -1) + " -c 'timeout 30 /usr/bin/php " + soCommand + "-d display_errors=on " + combinerPath + " " + keyPath + " " + encryptedData + " 2&1'"

	return command
}
