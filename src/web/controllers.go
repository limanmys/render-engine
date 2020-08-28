package web

import (
	"encoding/json"
	"errors"
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

type parsedRequest struct {
	Target      string
	UserID      string
	ExtensionID string
	ServerID    string
	RequestData map[string]string
	Token       string
	BaseURL     string
	Locale      string
	LogObject   sandbox.RegularLog
}

func validateAndExtractRequest(r *http.Request) (parsedRequest, error) {
	token := r.FormValue("token")
	requestData := map[string]string{}
	for key, values := range r.PostForm {
		requestData[key] = values[0]
	}
	var userID string

	if r.Header.Get("liman-token") != "" {
		userID = sqlite.GetUserIDFromLimanToken(r.Header.Get("liman-token"))
		if userID == "" {
			return parsedRequest{}, errors.New("Not Authorized1")
		}
	} else {
		userID = sqlite.GetUserIDFromToken(token)
		if userID == "" {
			return parsedRequest{}, errors.New("Not Authorized2")
		}
	}

	var target string
	var serverID string
	var extensionID string

	if r.FormValue("widget_id") != "" {
		widget := sqlite.GetWidget(r.FormValue("widget_id"))
		if widget.Name == "" {
			return parsedRequest{}, errors.New("Widget Bulunamadı")
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
		return parsedRequest{}, errors.New("Not Authorized3")
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

	parsedRequest := parsedRequest{
		Target:      target,
		UserID:      userID,
		ExtensionID: extensionID,
		ServerID:    serverID,
		RequestData: requestData,
		Token:       token,
		BaseURL:     baseURL,
		Locale:      locale,
		LogObject:   logObject,
	}

	return parsedRequest, nil
}

func runExtensionHandler(w http.ResponseWriter, r *http.Request) {
	parsedRequest, err := validateAndExtractRequest(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	command := sandbox.GeneratePHPCommand(parsedRequest.Target, parsedRequest.UserID, parsedRequest.ExtensionID, parsedRequest.ServerID, parsedRequest.RequestData, parsedRequest.Token, parsedRequest.BaseURL, parsedRequest.Locale, parsedRequest.LogObject)

	sandbox.WriteRegularLog(parsedRequest.LogObject)

	output := executeCommand(command)
	var objmap map[string]json.RawMessage
	err = json.Unmarshal([]byte(output), &objmap)
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

func externalAPIHandler(w http.ResponseWriter, r *http.Request) {
	parsedRequest, err := validateAndExtractRequest(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	command := sandbox.GeneratePHPCommand(parsedRequest.Target, parsedRequest.UserID, parsedRequest.ExtensionID, parsedRequest.ServerID, parsedRequest.RequestData, parsedRequest.Token, parsedRequest.BaseURL, parsedRequest.Locale, parsedRequest.LogObject)

	sandbox.WriteRegularLog(parsedRequest.LogObject)

	output := executeCommand(command)
	var objmap map[string]json.RawMessage
	err = json.Unmarshal([]byte(output), &objmap)
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

func backgroundJobHandler(w http.ResponseWriter, r *http.Request) {
	parsedRequest, err := validateAndExtractRequest(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	parsedRequest.LogObject.Display = "false"
	command := sandbox.GeneratePHPCommand(parsedRequest.Target, parsedRequest.UserID, parsedRequest.ExtensionID, parsedRequest.ServerID, parsedRequest.RequestData, parsedRequest.Token, parsedRequest.BaseURL, parsedRequest.Locale, parsedRequest.LogObject)
	sandbox.WriteRegularLog(parsedRequest.LogObject)
	go executeCommand(command)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok\n"))
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
