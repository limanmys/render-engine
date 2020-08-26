package web

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"renderer/src/sandbox"
	"renderer/src/sqlite"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type message struct {
	Message string
	Status  int
}

func runExtensionHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	requestData := map[string]string{}
	for key, values := range r.PostForm {
		requestData[key] = values[0]
	}
	var userID string

	if r.Header.Get("liman-token") != "" {
		userID = sqlite.GetUserIDFromLimanToken(r.Header.Get("liman-token"))
		if userID == "" {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("nope1"))
			return
		}
	} else {
		userID = sqlite.GetUserIDFromToken(token)
		if userID == "" {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("nope2"))
			return
		}
	}

	var target string
	var serverID string
	var extensionID string

	if r.FormValue("widget_id") != "" {
		widget := sqlite.GetWidget(r.FormValue("widget_id"))
		if widget.Name == "" {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Widget bulunamadı"))
		}
		target = widget.Function
		serverID = widget.ServerID
		extensionID = widget.ExtensionID
	} else {
		target = r.FormValue("lmntargetFunction")
		serverID = r.FormValue("server_id")
		extensionID = r.FormValue("extension_id")
	}

	if target == "" || serverID == "" || extensionID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("nope3"))
		return
	}

	baseURL := r.FormValue("lmnbaseurl")
	locale := r.FormValue("locale")
	if locale == "" {
		locale = "tr"
	}

	logObject := sandbox.RegularLog{
		UserID:      userID,
		ExtensionID: extensionID,
		ServerID:    serverID,
		IPAddress:   r.Header.Get("X-Real-IP"),
		Display:     "true",
		View:        target,
		LogID:       uuid.New().String(),
	}

	command := sandbox.GeneratePHPCommand(target, userID, extensionID, serverID, requestData, token, baseURL, locale, logObject)

	sandbox.WriteRegularLog(logObject)

	output := executeCommand(command)
	var objmap map[string]json.RawMessage
	err := json.Unmarshal([]byte(output), &objmap)
	contentType := "text/plain"
	var status int
	if err != nil {
		status = 200
	} else {
		contentType = "application/json"
		status, err = strconv.Atoi(strings.Trim(string(objmap["status"]), "\""))
		if err != nil {
			status = 200
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(output + "\n"))
}

func extensionLogHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	requestData := map[string]string{}
	for key, values := range r.PostForm {
		requestData[key] = values[0]
	}
	var userID string

	userID = sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("nope2"))
		return
	}

	specialLog := sandbox.SpecialLog{
		UserID:  userID,
		Message: r.FormValue("message"),
		Title:   r.FormValue("title"),
		LogID:   r.FormValue("log_id"),
	}

	sandbox.WriteSpecialLog(specialLog)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}

func executeCommand(input string) string {
	cmd := exec.Command("/bin/bash", "-c", input)
	stdout, _ := cmd.Output()
	return string(stdout)
}
