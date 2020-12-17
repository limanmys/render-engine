package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/postgresql"

	"github.com/mervick/aes-everywhere/go/aes256"
)

// GeneratePHPCommand generate command
func GeneratePHPCommand(targetFunction string, userID string, extensionID string, serverID string, requestData map[string]string, token string, baseURL string, locale string, logObject RegularLog) (string, error) {
	result := make(map[string]string)
	combinerPath := "/liman/sandbox/php/index.php"
	server, extension, settings := postgresql.GetUserData(serverID, extensionID, userID)
	user := postgresql.GetUser(userID)
	clientUsername, clientPassword, _, serverKey := postgresql.GetServerKey(userID, serverID)

	if clientUsername != "" && clientPassword != "" {
		settings["clientUsername"] = clientUsername
		settings["clientPassword"] = clientPassword
	}

	result["key_type"] = serverKey.Type

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

	license, _ := postgresql.GetLicense(extension.ID)

	result["license"] = license.Data

	result["apiRoute"] = "/extensionRun"

	result["navigationRoute"] = "/l/" + extension.ID + "/" + server.City + "/" + server.ID

	result["token"] = token

	result["locale"] = locale

	result["log_id"] = logObject.LogID

	result["ajax"] = "true"

	result["publicPath"] = baseURL + "/eklenti/" + extension.ID + "/public/"

	tmpPermissions := postgresql.GetFuncPermissions(userID)
	b, _ = json.Marshal(tmpPermissions)
	result["permissions"] = string(b)

	tmpVariables := postgresql.GetVariables(userID)
	b, _ = json.Marshal(tmpVariables)
	result["variables"] = string(b)

	extensionJSONFile, err := ioutil.ReadFile("/liman/extensions/" + strings.ToLower(extension.Name) + "/db.json")
	if err != nil {
		return "", errors.New(err.Error())
	}

	jsonMap := make(map[string][]map[string]string)

	_ = json.Unmarshal(extensionJSONFile, &jsonMap)

	requiredList := []string{}
	for i := 0; i < len(jsonMap["functions"]); i++ {
		if jsonMap["functions"][i]["isActive"] == "true" {
			requiredList = append(requiredList, jsonMap["functions"][i]["name"])
		}
	}

	if user.Status != 1 && !helpers.Contains(tmpPermissions, targetFunction) && helpers.Contains(requiredList, targetFunction) {
		return "", errors.New("Bu işlem için yetkiniz yok")
	}

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

	return command, nil
}

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}
