package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"renderer/src/connector"
	"renderer/src/sqlite"
	"strconv"
	"strings"
	"time"
)

func putFileHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"server_id", "remote_path", "local_path"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := sqlite.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	flag := val.Put(request["local_path"], request["remote_path"])

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}

}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"server_id", "remote_path", "local_path"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := sqlite.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	flag := val.Get(request["local_path"], request["remote_path"])

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}
}

func runCommandHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"server_id", "command"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := sqlite.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(val.Run(request["command"])))
}

func openTunnelHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"remote_host", "remote_port", "username", "password"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	port := connector.CreateTunnel(request["remote_host"], request["remote_port"], request["username"], request["password"])
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strconv.Itoa(port)))
}

func runOutsideCommandHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"command", "connection_type", "remote_host", "remote_port", "username", "password"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var val connector.Connection
	if val2, ok := connector.ActiveConnections[request["user_id"]+request["remote_host"]+request["username"]]; ok {
		val = val2
	} else {
		res := val.CreateShellRaw(request["connection_type"], request["username"], request["password"], request["remote_host"], request["remote_port"])
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	connector.ActiveConnections[request["user_id"]+request["remote_host"]+request["username"]] = val
	output := val.Run(request["command"])
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"ip_address", "username", "password", "port", "keyType"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	flag := connector.VerifyAuth(request["username"], request["password"], request["ip_address"], request["port"], request["keyType"])

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}
}

func runScriptHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"local_path", "root", "parameters", "server_id"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	server := sqlite.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	remotePath := "/tmp/" + filepath.Base(request["local_path"])
	if server.Os == "windows" {
		letter := val.Run("$pwd.drive.name")
		remotePath = strings.TrimSpace(letter) + ":\\Windows\\Temp\\" + filepath.Base(request["local_path"]) + ".ps1"
	}
	flag := val.Put(request["local_path"], remotePath)

	if server.Os == "linux" {
		val.Run("chmod +x " + remotePath)
	}

	if flag == false {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("Cannot send the script!"))
		return
	}

	output := val.Run(remotePath + " " + request["parameters"])
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}
