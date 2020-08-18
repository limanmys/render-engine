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
func GeneratePHPCommand(targetFunction string, userID string, extensionID string, serverID string, requestData map[string]string, token string, isAJAX bool) string {
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

	result["apiRoute"] = "https://liman.mcelen.dev/extensionRun"

	result["navigationRoute"] = result["apiRoute"]

	result["token"] = token

	result["locale"] = "tr"

	result["ajax"] = "true"

	soPath := "/liman/extensions/" + strings.ToLower(extension.Name) + "/liman.so"
	soCommand := ""
	if _, err := os.Stat(soPath); err == nil {
		soCommand = "-dextension=" + soPath + " "
	}

	keyPath := "/liman/keys/" + extension.ID
	content, _ := ioutil.ReadFile(keyPath)

	b, _ = json.Marshal(result)

	encryptedData := aes256.Encrypt(string(b), string(content))

	command := "sudo runuser " + strings.Replace(extension.ID, "-", "", -1) + " -c 'timeout 30 /usr/bin/php " + soCommand + "-d display_errors=on " + combinerPath + " " + keyPath + " " + encryptedData + "'"

	return command
}
