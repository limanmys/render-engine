package sandbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/limanmys/go/helpers"
	"time"
)

//RegularLog regular extension log.
type RegularLog struct {
	UserID      string `json:"user_id"`
	ExtensionID string `json:"extension_id"`
	ServerID    string `json:"server_id"`
	IPAddress   string `json:"ip_address"`
	Display     string `json:"display"`
	View        string `json:"view"`
	LogID       string `json:"log_id"`
}

//SpecialLog special extension log.
type SpecialLog struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Title   string `json:"title"`
	LogID   string `json:"log_id"`
	Data    string `json:"data"`
}

//WriteRegularLog Write Regular Log Object
func WriteRegularLog(logObject RegularLog) {
	b, _ := json.Marshal(logObject)
	now := time.Now().Format("2006-01-02 15:04:05")
	data := fmt.Sprintf("[%s] liman_render: EXTENSION_RENDER_PAGE %s\n", now, string(b))
	data = base64.StdEncoding.EncodeToString([]byte(data))
	helpers.ExecuteCommand("echo " + data + "| base64 --decode | tee --append /liman/logs/liman.log")
}

//WriteSpecialLog Write Special extension log object.
func WriteSpecialLog(logObject SpecialLog) {
	b, _ := json.Marshal(logObject)
	now := time.Now().Format("2006-01-02 15:04:05")
	data := fmt.Sprintf("[%s] liman_render: %s\n", now, string(b))
	data = base64.StdEncoding.EncodeToString([]byte(data))
	helpers.ExecuteCommand("echo " + data + "| base64 --decode | tee --append /liman/logs/extension.log")
}
