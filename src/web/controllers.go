package web

import (
	"net/http"
	"os/exec"
	"renderer/src/sandbox"
	"renderer/src/sqlite"
)

func runExtensionHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("NOPE"))
	}

	target := r.FormValue("target")
	serverID := r.FormValue("serverID")
	extensionID := r.FormValue("extensionID")

	requestData := make(map[string]string)

	command := sandbox.GeneratePHPCommand(target, userID, extensionID, serverID, requestData, token, false)
	w.WriteHeader(http.StatusOK)
	// output := executeCommand(command)
	_, _ = w.Write([]byte(command))
}

func executeCommand(input string) string {
	cmd := exec.Command("/bin/bash/", "-c", input)
	stdout, stderr := cmd.Output()
	if stderr != nil {
		return stderr.Error()
	}
	return string(stdout)
}
