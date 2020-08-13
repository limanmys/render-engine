package web

import (
	"net/http"
	"os/exec"
	"renderer/src/sandbox"
	"renderer/src/sqlite"
)

func runExtensionHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	requestData := map[string]string{}
	for key, values := range r.PostForm {
		requestData[key] = values[0]
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("qweqwe"))
		return
	}

	target := r.FormValue("target")
	serverID := r.FormValue("server_id")
	extensionID := r.FormValue("extension_id")

	command := sandbox.GeneratePHPCommand(target, userID, extensionID, serverID, requestData, token, false)
	w.WriteHeader(http.StatusOK)
	output := executeCommand(command)
	_, _ = w.Write([]byte(output))
}

func executeCommand(input string) string {
	cmd := exec.Command("/bin/bash", "-c", input)
	stdout, stderr := cmd.Output()
	if stderr != nil {
		return stderr.Error()
	}
	return string(stdout)
}
